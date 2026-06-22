package sdnotify

import (
	"context"
	"errors"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"testing"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/component/componentstatus"
	"go.uber.org/zap/zaptest"
)

// notifyRecorder spawns a goroutine that reads datagrams from a unix socket
// and stores them. Used by lifecycle tests to assert what we sent to systemd.
type notifyRecorder struct {
	mu   sync.Mutex
	msgs []string
}

func (r *notifyRecorder) all() []string {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make([]string, len(r.msgs))
	copy(out, r.msgs)
	return out
}

func (r *notifyRecorder) contains(want string) bool {
	for _, m := range r.all() {
		if m == want {
			return true
		}
	}
	return false
}

func (r *notifyRecorder) countPrefix(prefix string) int {
	r.mu.Lock()
	defer r.mu.Unlock()
	n := 0
	for _, m := range r.msgs {
		if len(m) >= len(prefix) && m[:len(prefix)] == prefix {
			n++
		}
	}
	return n
}

// startNotifyRecorder binds an AF_UNIX SOCK_DGRAM listener at a short path
// (macOS caps sun_path at ~104 bytes, which t.TempDir() blows past), sets
// NOTIFY_SOCKET to it, and returns a recorder that captures every datagram.
func startNotifyRecorder(t *testing.T) *notifyRecorder {
	t.Helper()
	// Use os.MkdirTemp under a short prefix instead of t.TempDir() because
	// t.TempDir() embeds the (potentially very long) test name in the path.
	dir, err := os.MkdirTemp("", "sdn")
	if err != nil {
		t.Fatalf("MkdirTemp: %v", err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(dir) })
	path := filepath.Join(dir, "n.sock")

	addr := &net.UnixAddr{Name: path, Net: "unixgram"}
	conn, err := net.ListenUnixgram("unixgram", addr)
	if err != nil {
		t.Fatalf("ListenUnixgram: %v", err)
	}
	t.Setenv("NOTIFY_SOCKET", path)

	rec := &notifyRecorder{}
	stop := make(chan struct{})
	go func() {
		buf := make([]byte, 4096)
		for {
			_ = conn.SetReadDeadline(time.Now().Add(50 * time.Millisecond))
			n, _, err := conn.ReadFromUnix(buf)
			if n > 0 {
				rec.mu.Lock()
				rec.msgs = append(rec.msgs, string(buf[:n]))
				rec.mu.Unlock()
			}
			if err != nil {
				select {
				case <-stop:
					return
				default:
				}
				var ne net.Error
				if errors.As(err, &ne) && ne.Timeout() {
					continue
				}
				return
			}
		}
	}()
	t.Cleanup(func() {
		close(stop)
		_ = conn.Close()
	})
	return rec
}

// waitFor polls fn until it returns true or the deadline fires.
func waitFor(t *testing.T, d time.Duration, fn func() bool) bool {
	t.Helper()
	deadline := time.Now().Add(d)
	for time.Now().Before(deadline) {
		if fn() {
			return true
		}
		time.Sleep(5 * time.Millisecond)
	}
	return fn()
}

func TestLifecycle_BasicReadyAndStopping(t *testing.T) {
	rec := startNotifyRecorder(t)
	ext := newSDNotify(&Config{}, zaptest.NewLogger(t))

	if err := ext.Start(context.Background(), nil); err != nil {
		t.Fatalf("Start: %v", err)
	}
	// No READY=1 from Start anymore -- it must wait for Ready().
	if rec.contains("READY=1") {
		t.Fatalf("READY=1 was sent before Ready() was called")
	}

	if err := ext.Ready(); err != nil {
		t.Fatalf("Ready: %v", err)
	}
	if !waitFor(t, time.Second, func() bool { return rec.contains("READY=1") }) {
		t.Fatalf("expected READY=1; got %v", rec.all())
	}

	if err := ext.NotReady(); err != nil {
		t.Fatalf("NotReady: %v", err)
	}
	if err := ext.Shutdown(context.Background()); err != nil {
		t.Fatalf("Shutdown: %v", err)
	}
	if !waitFor(t, time.Second, func() bool { return rec.contains("STOPPING=1") }) {
		t.Fatalf("expected STOPPING=1; got %v", rec.all())
	}
}

func TestLifecycle_FailIfNotSupervised(t *testing.T) {
	// Explicitly clear NOTIFY_SOCKET to simulate running outside systemd.
	t.Setenv("NOTIFY_SOCKET", "")
	ext := newSDNotify(&Config{FailIfNotSupervised: true}, zaptest.NewLogger(t))
	if err := ext.Start(context.Background(), nil); err == nil {
		t.Fatalf("expected Start to fail when NOTIFY_SOCKET is unset and FailIfNotSupervised is true")
	}
}

func TestLifecycle_WatchdogPings(t *testing.T) {
	rec := startNotifyRecorder(t)
	// Tell go-systemd's SdWatchdogEnabled the watchdog is active. The pinger
	// fires every d/2 so we'll see multiple WATCHDOG=1 within the test budget.
	t.Setenv("WATCHDOG_USEC", "200000") // 200ms -> ping every 100ms
	t.Setenv("WATCHDOG_PID", strconv.Itoa(os.Getpid()))

	ext := newSDNotify(&Config{EnableWatchdog: true}, zaptest.NewLogger(t))
	if err := ext.Start(context.Background(), nil); err != nil {
		t.Fatalf("Start: %v", err)
	}
	t.Cleanup(func() { _ = ext.Shutdown(context.Background()) })

	if !waitFor(t, 2*time.Second, func() bool { return rec.countPrefix("WATCHDOG=1") >= 2 }) {
		t.Fatalf("expected >=2 WATCHDOG=1 pings; got %v", rec.all())
	}
}

func TestLifecycle_ComponentStatusChangedShallowMode(t *testing.T) {
	// In shallow mode ComponentStatusChanged only logs; it must not send STATUS=.
	rec := startNotifyRecorder(t)
	ext := newSDNotify(&Config{}, zaptest.NewLogger(t))
	if err := ext.Start(context.Background(), nil); err != nil {
		t.Fatalf("Start: %v", err)
	}

	id := componentstatus.NewInstanceID(component.MustNewID("nop"), component.KindReceiver)
	ext.ComponentStatusChanged(id, componentstatus.NewPermanentErrorEvent(errors.New("synthetic")))

	time.Sleep(50 * time.Millisecond)
	if rec.countPrefix("STATUS=") != 0 {
		t.Fatalf("shallow mode must not send STATUS=; got %v", rec.all())
	}
	_ = ext.Shutdown(context.Background())
}

func TestLifecycle_ComponentStatusChangedDeepModeSendsStatus(t *testing.T) {
	// In deep mode, a permanent error must produce a STATUS= line immediately
	// (without needing the gRPC watcher to be wired up -- this is the
	// fast-push path).
	rec := startNotifyRecorder(t)
	ext := newSDNotify(&Config{
		DeepHealthcheck: true,
		HealthcheckV2:   component.MustNewID("healthcheckv2"),
	}, zaptest.NewLogger(t))

	id := componentstatus.NewInstanceID(component.MustNewID("nop"), component.KindReceiver)
	ext.ComponentStatusChanged(id, componentstatus.NewPermanentErrorEvent(errors.New("boom")))

	if !waitFor(t, time.Second, func() bool { return rec.countPrefix("STATUS=") >= 1 }) {
		t.Fatalf("expected STATUS= in deep mode; got %v", rec.all())
	}

	// Dedup: same event should not produce another STATUS=.
	ext.ComponentStatusChanged(id, componentstatus.NewPermanentErrorEvent(errors.New("boom")))
	time.Sleep(50 * time.Millisecond)
	if got := rec.countPrefix("STATUS="); got != 1 {
		t.Fatalf("expected 1 STATUS= after dedup, got %d (%v)", got, rec.all())
	}
}

func TestDependencies(t *testing.T) {
	t.Run("nil when DeepHealthcheck is off", func(t *testing.T) {
		ext := newSDNotify(&Config{}, zaptest.NewLogger(t))
		if got := ext.Dependencies(); got != nil {
			t.Fatalf("want nil, got %v", got)
		}
	})
	t.Run("returns healthcheckv2 ID when DeepHealthcheck is on", func(t *testing.T) {
		hcID := component.MustNewID("healthcheckv2")
		ext := newSDNotify(&Config{
			DeepHealthcheck: true,
			HealthcheckV2:   hcID,
		}, zaptest.NewLogger(t))
		got := ext.Dependencies()
		if len(got) != 1 || got[0] != hcID {
			t.Fatalf("want [%s], got %v", hcID, got)
		}
	})
}
