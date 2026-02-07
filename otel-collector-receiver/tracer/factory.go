package tracer

import (
	"context"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/receiver"
)

// NewFactory creates a factory for tracer receiver.
func NewFactory() receiver.Factory {
	return receiver.NewFactory(
		component.MustNewType("tracer"),
		func() component.Config {
			return &Config{
				Interval: (1 * time.Minute).String(),
			}
		},
		receiver.WithTraces(
			func(
				_ context.Context,
				params receiver.Settings,
				baseCfg component.Config,
				consumer consumer.Traces,
			) (receiver.Traces, error) {
				return &tracerReceiver{
					logger:       params.Logger,
					nextConsumer: consumer,
					config:       baseCfg.(*Config),
				}, nil
			},
			component.StabilityLevelAlpha,
		),
	)
}
