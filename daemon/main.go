// A small, complete reference for a systemd-aware Go daemon.
//
// Demonstrates:
//   - Type=notify lifecycle: READY=1 / RELOADING=1 / STOPPING=1 / STATUS=...
//   - Watchdog pings (WATCHDOG=1) gated on a real health check
//   - SIGHUP triggers a config reload (with proper RELOADING/READY bracketing)
//   - SIGTERM/SIGINT triggers graceful HTTP drain inside TimeoutStopSec
package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/coreos/go-systemd/v22/daemon"
)

const addr = ":8081"

// healthy is the single source of truth for "should we tell systemd we're
// alive?". A real service would flip this off if the DB connection died,
// the request loop stalled, etc. Kept as an atomic bool so the watchdog
// goroutine can read it without locking.
var healthy atomic.Bool

func main() {
	log.SetFlags(log.Ltime | log.Lmicroseconds)

	// 1. Bind the listener.
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("listen: %v", err)
	}
	log.Printf("listening on %s", ln.Addr())

	// 2. Build the HTTP server. Handlers stay trivial — the point is the
	//    surrounding lifecycle, not the business logic.
	// /healthz gets instead result
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("hello from systemd-aware daemon\n"))
	})
	// /slow lets you observe graceful drain: hold a request open while you
	// `systemctl stop` and watch it finish before the process exits.
	mux.HandleFunc("/slow", func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-time.After(10 * time.Second):
			_, _ = w.Write([]byte("slow done\n"))
		case <-r.Context().Done():
		}
	})
	// /switch changes the health status
	mux.HandleFunc("/switch", func(w http.ResponseWriter, r *http.Request) {
		healthy.Swap(!healthy.Load())
		_, _ = w.Write([]byte("switched to " + strconv.FormatBool(healthy.Load())))
	})
	srv := &http.Server{Handler: mux}

	// 3. Simulate startup work. With Type=notify, dependents wait for READY=1,
	//    so this delay is observable as a delayed `systemctl start` return.
	log.Printf("warming up...")
	time.Sleep(2 * time.Second)
	healthy.Store(true)

	// 4. Start serving in a goroutine so we can sit in the signal loop below.
	serveErr := make(chan error, 1)
	go func() {
		err := srv.Serve(ln)
		// Shutdown returns ErrServerClosed; that's the happy path.
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			serveErr <- err
		}
		close(serveErr)
	}()

	// 5. Tell systemd we're ready *after* the listener is up and warmup done.
	//    Outside systemd, NOTIFY_SOCKET is unset and SdNotify returns (false, nil).
	if sent, err := daemon.SdNotify(false, daemon.SdNotifyReady); err != nil {
		log.Printf("sd_notify READY failed: %v", err)
	} else if sent {
		log.Printf("notified systemd: READY=1")
		_, _ = daemon.SdNotify(false, "STATUS=Serving HTTP on "+ln.Addr().String())
	}

	// 6. Watchdog: only runs if WatchdogSec= is set in the unit file.
	stopWatchdog := startWatchdog()
	defer stopWatchdog()

	// 7. Signal loop. SIGHUP = reload, SIGTERM/SIGINT = graceful shutdown.
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)

	for {
		select {
		case err := <-serveErr:
			if err != nil {
				log.Fatalf("http serve: %v", err)
			}
			return
		case sig := <-sigs:
			switch sig {
			case syscall.SIGHUP:
				reload()
			case syscall.SIGTERM, syscall.SIGINT:
				shutdown(srv)
				return
			}
		}
	}
}

// startWatchdog spawns a goroutine that pings systemd at half the configured
// interval, but only when healthy.Load() is true. Returns a stop function.
//
// Half-interval is the conventional choice: it absorbs scheduling jitter
// without missing the deadline. Pinging unconditionally would defeat the
// watchdog's purpose — the whole point is that a deadlocked main loop stops
// the pings.
func startWatchdog() func() {
	interval, err := daemon.SdWatchdogEnabled(false)
	if err != nil {
		log.Printf("watchdog probe failed: %v", err)
		return func() {}
	}
	if interval == 0 {
		log.Printf("watchdog disabled (no WatchdogSec= in unit)")
		return func() {}
	}
	log.Printf("watchdog enabled, pinging every %s", interval/2)

	stop := make(chan struct{})
	done := make(chan struct{})
	go func() {
		defer close(done)
		t := time.NewTicker(interval / 2)
		defer t.Stop()
		for {
			select {
			case <-stop:
				return
			case <-t.C:
				if healthy.Load() {
					_, _ = daemon.SdNotify(false, daemon.SdNotifyWatchdog)
				}
			}
		}
	}()
	return func() {
		close(stop)
		<-done
	}
}

// reload simulates re-reading config. The RELOADING=1 / READY=1 bracket is
// what `Type=notify` watches for — `systemctl reload` blocks until READY=1
// comes back, so dependents see a coherent "still up" signal.
func reload() {
	log.Printf("SIGHUP received, reloading")
	if _, err := daemon.SdNotify(false, daemon.SdNotifyReloading); err != nil {
		log.Printf("sd_notify RELOADING failed: %v", err)
	}

	healthy.Store(false) // pretend we're momentarily not serving
	time.Sleep(500 * time.Millisecond)
	healthy.Store(true)

	if _, err := daemon.SdNotify(false, daemon.SdNotifyReady); err != nil {
		log.Printf("sd_notify READY failed: %v", err)
	}
	log.Printf("reload done")
}

// shutdown drains in-flight requests and tells systemd we're going away.
func shutdown(srv *http.Server) {
	log.Printf("shutdown requested, draining")
	if _, err := daemon.SdNotify(false, daemon.SdNotifyStopping); err != nil {
		log.Printf("sd_notify STOPPING failed: %v", err)
	}
	healthy.Store(false)

	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("graceful shutdown failed, forcing: %v", err)
		_ = srv.Close()
	}
	log.Printf("bye")
}
