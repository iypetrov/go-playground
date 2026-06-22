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
//
// The systemd watchdog pinger auto-enables whenever systemd has set
// WATCHDOG_USEC and WATCHDOG_PID matches our PID; no config knob is needed.
type Config struct {
	// FailIfNotSupervised makes Start return an error when the process is
	// not running under systemd (NOTIFY_SOCKET unset). Default: false.
	FailIfNotSupervised bool `mapstructure:"fail_if_not_supervised"`

	// DeepHealthcheck opts in to the per-component health aggregation mode.
	// When true the extension subscribes to the configured healthcheckv2
	// extension's grpc.health.v1.Health/Watch stream (service="") and
	// reflects status changes into STATUS=<text> notifications to systemd.
	// It also pushes a STATUS=... line immediately when a component reports
	// a permanent or fatal error through ComponentStatusChanged. Default
	// false -- when false the extension only sends READY=1 and STOPPING=1.
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
}

// Validate is called by the collector before Start.
func (c *Config) Validate() error {
	if c.DeepHealthcheck && c.HealthcheckV2 == (component.ID{}) {
		return errors.New("deep_healthcheck requires healthcheckv2 to be set")
	}
	return nil
}
