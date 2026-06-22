package sdnotify

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"sync"
	"time"

	"github.com/coreos/go-systemd/v22/daemon"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/component/componentstatus"
	"go.opentelemetry.io/collector/extension"
	"go.opentelemetry.io/collector/extension/extensioncapabilities"
	"go.uber.org/zap"
)

// Compile-time interface assertions. Dependent is implemented but isn't on
// the Extension alias because it only applies when DeepHealthcheck is on --
// the service detects it via type assertion anyway.
var (
	_ Extension                       = (*sdnotify)(nil)
	_ extensioncapabilities.Dependent = (*sdnotify)(nil)
)

type Extension interface {
	extension.Extension
	extensioncapabilities.PipelineWatcher
	componentstatus.Watcher
}

// sdnotify implements extension.Extension plus PipelineWatcher,
// componentstatus.Watcher, and (conditionally) Dependent.
//
//   - READY=1     is sent from Ready() once all pipelines are up.
//   - STOPPING=1  is sent from Shutdown().
//   - WATCHDOG=1  is pinged from a ticker goroutine when systemd has set
//     WATCHDOG_USEC for our PID (auto-detected, no config needed).
//   - STATUS=...  is sent on every status change observed via the
//     healthcheckv2 gRPC Watch stream when DeepHealthcheck is on, and
//     immediately on permanent/fatal ComponentStatusChanged events.
//
// All notifications are best-effort: if NOTIFY_SOCKET is unset, SdNotify
// returns sent=false, err=nil and we just log it -- unless the user opted
// into FailIfNotSupervised.
type sdnotify struct {
	cfg    *Config
	logger *zap.Logger

	host component.Host // captured in Start so Ready can resolve siblings

	// Lifecycle goroutines. All are lazily created; nil when their feature is off.
	wdCancel    context.CancelFunc
	wdDone      chan struct{}
	watchCancel context.CancelFunc
	watchDone   chan struct{}

	// gRPC client to healthcheckv2; nil unless DeepHealthcheck is true.
	hc *healthClient

	// Last STATUS line we sent; used to deduplicate STATUS= sends across
	// both the gRPC watcher and the ComponentStatusChanged push-path.
	lastStatusMu sync.Mutex
	lastStatus   string
}

func newSDNotify(cfg *Config, logger *zap.Logger) *sdnotify {
	return &sdnotify{cfg: cfg, logger: logger}
}

// Dependencies returns the IDs of extensions this one must be started after.
// Only non-nil when DeepHealthcheck is wired up; otherwise we impose no
// ordering constraint on the collector graph.
func (s *sdnotify) Dependencies() []component.ID {
	if !s.cfg.DeepHealthcheck || s.cfg.HealthcheckV2 == (component.ID{}) {
		return nil
	}
	return []component.ID{s.cfg.HealthcheckV2}
}

func (s *sdnotify) Start(_ context.Context, host component.Host) error {
	s.host = host

	// FailIfNotSupervised: a cheap probe (empty notify never reaches systemd
	// but daemon.SdNotify still returns sent=true when NOTIFY_SOCKET is set,
	// so a passive os.Getenv check is simpler and avoids spurious datagrams).
	if s.cfg.FailIfNotSupervised && os.Getenv("NOTIFY_SOCKET") == "" {
		return fmt.Errorf("sdnotify: NOTIFY_SOCKET not set; not running under systemd")
	}

	// Watchdog auto-enables whenever systemd has set WATCHDOG_USEC for our
	// PID. SdWatchdogEnabled returns 0 when it didn't, or when WATCHDOG_PID
	// points at a different process -- both are valid "not enabled" states
	// we treat as a no-op.
	d, err := daemon.SdWatchdogEnabled(false)
	switch {
	case err != nil:
		s.logger.Debug("sdnotify: SdWatchdogEnabled returned error; watchdog disabled",
			zap.Error(err))
	case d == 0:
		s.logger.Debug("sdnotify: WATCHDOG_USEC not set; watchdog disabled")
	default:
		s.startWatchdog(d)
	}

	return nil
}

// Ready is called by the collector once all pipelines have started and the
// service is ready to receive data. This is when READY=1 belongs.
func (s *sdnotify) Ready() error {
	sent, err := daemon.SdNotify(false, daemon.SdNotifyReady)
	if err != nil {
		return fmt.Errorf("sdnotify READY=1: %w", err)
	}
	if !sent {
		s.logger.Info("sdnotify: NOTIFY_SOCKET not set; READY=1 was a no-op")
	} else {
		s.logger.Info("sdnotify: sent READY=1 to systemd")
	}

	if s.cfg.DeepHealthcheck {
		if err := s.startDeepHealthcheck(); err != nil {
			// Don't fail the whole pipeline if we can't reach healthcheckv2;
			// the basic READY=1/STOPPING=1 path still works.
			s.logger.Warn("sdnotify: deep healthcheck startup failed; continuing without it",
				zap.Error(err))
		}
	}

	return nil
}

// NotReady is called before receivers are stopped during shutdown.
// We stop the gRPC watcher here so it can't race with hc.Close in Shutdown.
// STOPPING=1 itself is sent from Shutdown, not here, so it brackets the
// whole drain rather than firing while receivers are still draining data.
func (s *sdnotify) NotReady() error {
	s.stopDeepHealthcheck()
	return nil
}

func (s *sdnotify) Shutdown(_ context.Context) error {
	s.stopWatchdog()
	s.stopDeepHealthcheck() // idempotent if NotReady already ran

	sent, err := daemon.SdNotify(false, daemon.SdNotifyStopping)
	if err != nil {
		// Don't block shutdown on a notify failure -- just log it.
		s.logger.Warn("sdnotify STOPPING=1 failed", zap.Error(err))
		return nil
	}
	if sent {
		s.logger.Info("sdnotify: sent STOPPING=1 to systemd")
	}
	return nil
}

// ComponentStatusChanged is called by the collector core for every status
// change of every component (receiver/processor/exporter/connector/extension).
//
// In the default (shallow) mode it only logs. With DeepHealthcheck on, hard
// errors are also surfaced immediately as STATUS=... -- the subsequent gRPC
// Watch event from healthcheckv2 will replace this with the aggregated view.
func (s *sdnotify) ComponentStatusChanged(
	source *componentstatus.InstanceID,
	event *componentstatus.Event,
) {
	s.logger.Info("component status changed",
		zap.Stringer("kind", source.Kind()),
		zap.Stringer("component", source.ComponentID()),
		zap.Stringer("status", event.Status()),
		zap.Time("timestamp", event.Timestamp()),
		zap.Error(event.Err()), // nil unless Status is an error variant
	)

	if !s.cfg.DeepHealthcheck {
		return
	}
	switch event.Status() {
	case componentstatus.StatusPermanentError, componentstatus.StatusFatalError:
		line := fmt.Sprintf("component %s failed: %v", source.ComponentID(), event.Err())
		s.sendStatusLine(line)
	}
}

// --- watchdog ---

func (s *sdnotify) startWatchdog(interval time.Duration) {
	// Per sd_watchdog_enabled(3) the recommended cadence is interval/2.
	tickEvery := interval / 2
	if tickEvery <= 0 {
		tickEvery = interval
	}
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	s.wdCancel = cancel
	s.wdDone = done

	go func() {
		defer close(done)
		t := time.NewTicker(tickEvery)
		defer t.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				if _, err := daemon.SdNotify(false, daemon.SdNotifyWatchdog); err != nil {
					s.logger.Debug("sdnotify WATCHDOG=1 failed", zap.Error(err))
				}
			}
		}
	}()
	s.logger.Info("sdnotify: systemd watchdog enabled",
		zap.Duration("interval", interval),
		zap.Duration("ping_every", tickEvery))
}

func (s *sdnotify) stopWatchdog() {
	if s.wdCancel == nil {
		return
	}
	s.wdCancel()
	<-s.wdDone
	s.wdCancel = nil
	s.wdDone = nil
}

// --- deep healthcheck (gRPC Watch) ---

func (s *sdnotify) startDeepHealthcheck() error {
	endpoint, err := s.resolveHealthcheckV2Endpoint()
	if err != nil {
		return err
	}
	// service="" -> overall collector health (aggregate of all components).
	hc, err := dialHealthClient(endpoint, "", s.logger)
	if err != nil {
		return err
	}
	s.hc = hc

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	s.watchCancel = cancel
	s.watchDone = done

	updates := hc.Watch(ctx)

	go func() {
		defer close(done)
		for u := range updates {
			s.sendStatusLine(u.Line)
		}
	}()

	s.logger.Info("sdnotify: deep healthcheck enabled",
		zap.String("endpoint", endpoint))
	return nil
}

func (s *sdnotify) stopDeepHealthcheck() {
	if s.watchCancel != nil {
		s.watchCancel()
		<-s.watchDone
		s.watchCancel = nil
		s.watchDone = nil
	}
	if s.hc != nil {
		_ = s.hc.Close()
		s.hc = nil
	}
}

// resolveHealthcheckV2Endpoint picks the gRPC endpoint to dial. The
// explicit config override wins; otherwise we reflectively pull
// `.GRPC.ServerConfig.NetAddr.Endpoint` off the sibling healthcheckv2
// extension. Reflection avoids a hard import dependency on contrib.
func (s *sdnotify) resolveHealthcheckV2Endpoint() (string, error) {
	if s.cfg.HealthcheckV2GRPCEndpoint != "" {
		return s.cfg.HealthcheckV2GRPCEndpoint, nil
	}
	if s.host == nil {
		return "", fmt.Errorf("sdnotify: host not set; cannot resolve healthcheckv2 endpoint")
	}
	ext, ok := s.host.GetExtensions()[s.cfg.HealthcheckV2]
	if !ok || ext == nil {
		return "", fmt.Errorf("sdnotify: healthcheckv2 extension %q not found", s.cfg.HealthcheckV2)
	}
	if ep, ok := extractGRPCEndpoint(ext); ok && ep != "" {
		return ep, nil
	}
	return "", fmt.Errorf("sdnotify: could not extract gRPC endpoint from healthcheckv2 extension; "+
		"set healthcheckv2_grpc_endpoint explicitly (extension type=%T)", ext)
}

// extractGRPCEndpoint walks the extension struct looking for an exported
// "GRPC" field whose nested NetAddr/TCPAddr/Endpoint string is non-empty.
// Falls back gracefully if the contrib type layout changes.
func extractGRPCEndpoint(ext component.Component) (string, bool) {
	v := reflect.Indirect(reflect.ValueOf(ext))
	if !v.IsValid() || v.Kind() != reflect.Struct {
		return "", false
	}
	// Common candidate paths inside contrib's healthcheckv2:
	//   ext.config.GRPC.ServerConfig.NetAddr.Endpoint
	//   ext.cfg.GRPC.NetAddr.Endpoint
	// We BFS up to depth 4 looking for an "Endpoint" string field under a
	// field named "GRPC" (case-insensitive).
	return findEndpointUnderGRPC(v, 0)
}

func findEndpointUnderGRPC(v reflect.Value, depth int) (string, bool) {
	if depth > 4 || !v.IsValid() {
		return "", false
	}
	v = reflect.Indirect(v)
	if v.Kind() != reflect.Struct {
		return "", false
	}
	t := v.Type()
	// First pass: any field whose name equals "GRPC" -- recurse looking for Endpoint.
	for i := 0; i < t.NumField(); i++ {
		ft := t.Field(i)
		if !ft.IsExported() {
			continue
		}
		fv := v.Field(i)
		if eqFold(ft.Name, "GRPC") {
			if ep, ok := findEndpoint(fv, 0); ok {
				return ep, true
			}
		}
	}
	// Second pass: recurse into other struct fields.
	for i := 0; i < t.NumField(); i++ {
		ft := t.Field(i)
		if !ft.IsExported() {
			continue
		}
		fv := v.Field(i)
		if reflect.Indirect(fv).Kind() == reflect.Struct {
			if ep, ok := findEndpointUnderGRPC(fv, depth+1); ok {
				return ep, true
			}
		}
	}
	return "", false
}

func findEndpoint(v reflect.Value, depth int) (string, bool) {
	if depth > 4 || !v.IsValid() {
		return "", false
	}
	v = reflect.Indirect(v)
	if v.Kind() != reflect.Struct {
		return "", false
	}
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		ft := t.Field(i)
		if !ft.IsExported() {
			continue
		}
		fv := v.Field(i)
		if eqFold(ft.Name, "Endpoint") && fv.Kind() == reflect.String {
			if s := fv.String(); s != "" {
				return s, true
			}
		}
	}
	for i := 0; i < t.NumField(); i++ {
		ft := t.Field(i)
		if !ft.IsExported() {
			continue
		}
		fv := v.Field(i)
		if reflect.Indirect(fv).Kind() == reflect.Struct {
			if ep, ok := findEndpoint(fv, depth+1); ok {
				return ep, true
			}
		}
	}
	return "", false
}

func eqFold(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		ca, cb := a[i], b[i]
		if 'A' <= ca && ca <= 'Z' {
			ca += 'a' - 'A'
		}
		if 'A' <= cb && cb <= 'Z' {
			cb += 'a' - 'A'
		}
		if ca != cb {
			return false
		}
	}
	return true
}

// --- STATUS= sending (with dedup) ---

func (s *sdnotify) sendStatusLine(line string) {
	if line == "" {
		return
	}
	s.lastStatusMu.Lock()
	if line == s.lastStatus {
		s.lastStatusMu.Unlock()
		return
	}
	s.lastStatus = line
	s.lastStatusMu.Unlock()

	payload := "STATUS=" + line
	if _, err := daemon.SdNotify(false, payload); err != nil {
		s.logger.Debug("sdnotify STATUS= failed", zap.Error(err))
		return
	}
	s.logger.Debug("sdnotify: sent STATUS=", zap.String("line", line))
}
