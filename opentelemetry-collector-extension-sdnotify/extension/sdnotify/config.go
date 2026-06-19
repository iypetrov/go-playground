package sdnotify

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
}

// Validate is called by the collector before Start. Nothing to check yet.
func (c *Config) Validate() error { return nil }
