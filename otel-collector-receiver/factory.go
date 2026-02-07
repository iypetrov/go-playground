package main

import (
	"context"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/receiver"
)

// NewFactory creates a factory for otel-collector-receiver receiver.
func NewFactory() receiver.Factory {
	return receiver.NewFactory(
		component.MustNewType("otel-collector-receiver"),
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
				logger := params.Logger
				traceCfg := baseCfg.(*Config)
				traceRcvr := &traceReceiver{
					logger:       logger,
					nextConsumer: consumer,
					config:       traceCfg,
				}

				return traceRcvr, nil
			},
			component.StabilityLevelAlpha,
		),
	)
}
