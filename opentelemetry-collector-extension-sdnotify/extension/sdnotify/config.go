package sdnotify

import (
	"errors"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/confmap/xconfmap"
)

var _ xconfmap.Validator = (*Config)(nil)

// Config controls how the sdnotify extension talks to systemd.
//
// The extension uses go-systemd's daemon.SdNotify, which reads the
// NOTIFY_SOCKET env var that systemd injects into spawned services.
// When NOTIFY_SOCKET is unset (e.g. running outside systemd), notifications
// are no-ops by default; set FailIfNotSupervised: true to fail Start instead.
type Config struct {
	// FailIfNotSupervised makes Start return an error when the process is
	// not running under systemd (NOTIFY_SOCKET unset). Default: false.
	FailIfNotSupervised bool `mapstructure:"fail_if_not_supervised"`

	// EnableWatchdog turns on the systemd watchdog pinger. Requires WATCHDOG_USEC
	// to be set in the environment by systemd; otherwise it's a no-op. Default false.
	EnableWatchdog bool `mapstructure:"enable_watchdog"`

	// UnsetEnvironment passes unsetEnv=true to daemon.SdNotify so NOTIFY_SOCKET
	// is not inherited by child processes. Default false (matches current
	// behaviour and go-systemd's own default).
	UnsetEnvironment bool `mapstructure:"unset_environment"`

	// DeepHealthcheck opts in to the per-component health aggregation mode.
	// When true the extension subscribes to the configured healthcheckv2
	// extension's grpc.health.v1.Health/Watch stream and reflects status
	// changes into STATUS=<text> notifications to systemd. It also pushes a
	// STATUS=... line immediately when a component reports a permanent or
	// fatal error through ComponentStatusChanged. Default false -- when false
	// the extension only sends READY=1 and STOPPING=1.
	DeepHealthcheck bool `mapstructure:"deep_healthcheck"`

	// HealthcheckV2 is the component.ID of the healthcheckv2 extension to
	// depend on and subscribe to. Required when DeepHealthcheck is true. Used
	// for Dependencies() and to resolve the gRPC endpoint at runtime from
	// host.GetExtensions().
	HealthcheckV2 component.ID `mapstructure:"healthcheckv2"`

	// HealthcheckV2GRPCEndpoint is an optional override for the gRPC endpoint
	// to dial. When empty, the extension reads it from the sibling
	// healthcheckv2 extension's own config at startup. Provide it explicitly
	// only if that lookup proves unreliable across healthcheckv2 versions.
	HealthcheckV2GRPCEndpoint string `mapstructure:"healthcheckv2_grpc_endpoint"`

	// WatchService is the gRPC health service name to watch. Empty string
	// means overall collector health; e.g. "traces" watches just the traces
	// pipeline. Default "".
	WatchService string `mapstructure:"watch_service"`
}

// Validate is called by the collector before Start.
func (c *Config) Validate() error {
	if c.DeepHealthcheck && c.HealthcheckV2 == (component.ID{}) {
		return errors.New("deep_healthcheck requires healthcheckv2 to be set")
	}
	return nil
}
