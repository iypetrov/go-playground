package main

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/otlptranslator"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	otelexporterprom "go.opentelemetry.io/otel/exporters/prometheus"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"google.golang.org/grpc"
)

type GlobalMetricsSetup struct {
	provider     *sdkmetric.MeterProvider
	shutdownOnce sync.Once
}

func NewGlobalMetricsSetup(reg prometheus.Registerer) (*GlobalMetricsSetup, error) {
	// Enable OpenTelemetry Log SDK observability (self-instrumentation) metrics.
	// This is an experimental feature that emits metrics like otel.sdk.log.created,
	// otel.sdk.exporter.* etc. The environment variable must be set before SDK initialization.
	// See: https://pkg.go.dev/go.opentelemetry.io/otel/sdk/log/internal/x
	_ = os.Setenv("OTEL_GO_X_OBSERVABILITY", "true")

	// Create Prometheus exporter using the default registry
	// This ensures OTLP metrics are exposed on the same /metrics endpoint
	// as the existing Prometheus metrics (port 2021)
	promExporter, err := otelexporterprom.New(
		otelexporterprom.WithRegisterer(reg),
		otelexporterprom.WithNamespace("output_plugin"),
		otelexporterprom.WithTranslationStrategy(otlptranslator.UnderscoreEscapingWithSuffixes),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize prometheus exporter for OTLP metrics: %w", err)
	}

	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(promExporter),
	)

	// Set as global meter provider so instrumentation libraries can discover it
	otel.SetMeterProvider(meterProvider)

	return &GlobalMetricsSetup{
		provider: meterProvider,
	}, nil
}

// Provider returns the configured OpenTelemetry meter provider.
// The provider is used for creating meters and recording metrics.
func (m *GlobalMetricsSetup) Provider() *sdkmetric.MeterProvider {
	return m.provider
}

// GRPCStatsHandler returns a gRPC dial option that enables automatic
// metrics collection for gRPC client calls.
func (m *GlobalMetricsSetup) GRPCStatsHandler() grpc.DialOption {
	return grpc.WithStatsHandler(otelgrpc.NewClientHandler(
		otelgrpc.WithMeterProvider(m.provider),
	))
}

// Shutdown gracefully shuts down the meter provider and stops metrics collection.
//
// This method is idempotent - multiple calls are safe and will only perform
// the actual shutdown once. Subsequent calls return nil immediately.
//
// The context is used to enforce a timeout on the shutdown operation.
// If the context expires before shutdown completes, the context error is returned.
//
// After shutdown, the meter provider should not be used for new metric operations.
func (m *GlobalMetricsSetup) Shutdown(ctx context.Context) error {
	var shutdownErr error

	m.shutdownOnce.Do(func() {
		if err := m.provider.Shutdown(ctx); err != nil {
			shutdownErr = fmt.Errorf("failed to shutdown meter provider: %w", err)
		}
	})

	return shutdownErr
}
