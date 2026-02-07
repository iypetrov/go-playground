package tracer

import (
	"context"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.uber.org/zap"
)

type tracerReceiver struct {
	host   component.Host
	cancel context.CancelFunc

	logger       *zap.Logger
	nextConsumer consumer.Traces
	config       *Config
}

func (tr *tracerReceiver) Start(
	ctx context.Context,
	host component.Host,
) error {
	tr.host = host
	ctx = context.Background()
	ctx, tr.cancel = context.WithCancel(ctx)

	interval, err := time.ParseDuration(tr.config.Interval)
	if err != nil {
		return err
	}
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				tr.logger.Info("I should start processing traces now!")
			case <-ctx.Done():
				return
			}
		}
	}()

	return nil
}

func (tr *tracerReceiver) Shutdown(ctx context.Context) error {
	if tr.cancel != nil {
		tr.cancel()
	}
	return nil
}
