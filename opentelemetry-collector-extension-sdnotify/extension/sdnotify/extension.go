package sdnotify

import (
	"context"
	"fmt"

	"github.com/coreos/go-systemd/v22/daemon"
	"go.opentelemetry.io/collector/component"
	"go.uber.org/zap"
)

// sdnotify implements extension.Extension. On Start it sends READY=1 to
// systemd; on Shutdown it sends STOPPING=1. Both are best-effort: if the
// process is not supervised by systemd (NOTIFY_SOCKET unset), SdNotify
// returns sent=false, err=nil and we just log it -- unless the user
// opted into FailIfNotSupervised.
type sdnotify struct {
	cfg    *Config
	logger *zap.Logger
}

func newSDNotify(cfg *Config, logger *zap.Logger) *sdnotify {
	return &sdnotify{cfg: cfg, logger: logger}
}

func (s *sdnotify) Start(_ context.Context, _ component.Host) error {
	sent, err := daemon.SdNotify(false, daemon.SdNotifyReady)
	if err != nil {
		return fmt.Errorf("sdnotify READY=1: %w", err)
	}
	if !sent {
		if s.cfg.FailIfNotSupervised {
			return fmt.Errorf("sdnotify: NOTIFY_SOCKET not set; not running under systemd")
		}
		s.logger.Info("sdnotify: NOTIFY_SOCKET not set; READY=1 was a no-op")
		return nil
	}
	s.logger.Info("sdnotify: sent READY=1 to systemd")
	return nil
}

func (s *sdnotify) Shutdown(_ context.Context) error {
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
