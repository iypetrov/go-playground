package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/go-logr/logr"
	"github.com/iypetrov/ploggrpc-exporter/metrics"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/exporter/exportertest"
	"go.opentelemetry.io/collector/exporter/otlpexporter"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
)

var (
	pluginName              = "fluent-bit-output-plugin"
	schemaURL               = "https://opentelemetry.io/schemas/1.27.0"
	version                 = "0.1.0"
	endpoint                = "localhost:4317"
	hostname                = "some-random-vm"
	eventGenerationInterval = 1 * time.Second
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	loggerOpts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}
	loggerHandler := slog.NewTextHandler(os.Stderr, loggerOpts)
	logger := logr.FromSlogHandler(loggerHandler)

	// Metrics server config
	registry := metrics.NewRegistry()
	m := metrics.NewPluginMetrics(registry)
	metricsSetup, err := NewMetricsSetup(registry)
	if err != nil {
		panic(err)
	}
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
	srv := &http.Server{
		Addr:    ":2021",
		Handler: mux,
	}

	// TODO: real logic
	//
	// CreateDefaultConfig() already enables:
	//   - TimeoutConfig (5s)
	//   - RetryConfig   (exponential backoff)
	//   - QueueConfig   (in-memory queue + batcher merged into QueueBatchConfig)
	// So queue/retry/batch are on without any extra wiring.
	factory := otlpexporter.NewFactory()
	cfg := factory.CreateDefaultConfig().(*otlpexporter.Config)
	cfg.ClientConfig.Endpoint = endpoint
	cfg.ClientConfig.TLS.Insecure = true

	// exporter.Settings: start from the nop helper for sane defaults
	// (zap.NewNop, otel noop tracer/meter), then swap MeterProvider so
	// exporterhelper internal metrics flow into our Prometheus registry.
	otlpType := component.MustNewType("otlp")
	set := exportertest.NewNopSettings(otlpType)
	set.ID = component.NewID(otlpType)
	set.TelemetrySettings.MeterProvider = metricsSetup.Provider()

	logsExp, err := factory.CreateLogs(ctx, set, cfg)
	if err != nil {
		logger.Error(err, "failed to create otlp logs exporter")
		cancel()
		return
	}

	if err := logsExp.Start(ctx, componenttest.NewNopHost()); err != nil {
		logger.Error(err, "failed to start otlp logs exporter")
		cancel()
		return
	}

	// Graceful shutdown on SIGINT/SIGTERM
	var wg sync.WaitGroup
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// Log emitter — push plog.Logs through the component pipeline. ConsumeLogs
	// returns once the item is enqueued; the queue/batcher/retry are handled
	// by exporterhelper underneath.
	wg.Go(func() {
		ticker := time.NewTicker(eventGenerationInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				msg := map[string]any{
					"name": "John",
				}
				logs := buildLogs(msg)

				if err := logsExp.ConsumeLogs(ctx, logs); err != nil {
					logger.Error(err, "failed to enqueue logs")
					continue
				}
				logger.Info("log record enqueued", "msg", msg)
				m.OutputClientLogs.WithLabelValues(hostname).Inc()
			case <-ctx.Done():
				return
			}
		}
	})

	// Metric server
	wg.Go(func() {
		go func() {
			<-ctx.Done()

			shutdownCtx, cancel := context.WithTimeout(
				context.Background(),
				5*time.Second,
			)
			defer cancel()
			// defer metricsSetup.Shutdown(shutdownCtx)

			if err := srv.Shutdown(shutdownCtx); err != nil {
				logger.Error(err, "failed shutting down metrics server")
			}
		}()

		if err := srv.ListenAndServe(); err != nil {
			logger.Error(err, "fluent-bit-output-plugin")
		}
	})

	<-sigCh
	logger.Info("shutting down...")
	cancel()
	wg.Wait()

	// Drain the exporter queue. Use a real timeout so queued items get
	// flushed instead of dropped.
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer shutdownCancel()
	if err := logsExp.Shutdown(shutdownCtx); err != nil {
		logger.Error(err, "failed to shutdown logs exporter")
	}
}

// buildLogs constructs a plog.Logs payload with one ResourceLogs / ScopeLogs /
// LogRecord — same shape as ploggrpc-exporter's variant.
func buildLogs(msg map[string]any) plog.Logs {
	logs := plog.NewLogs()

	rl := logs.ResourceLogs().AppendEmpty()
	rl.SetSchemaUrl(schemaURL)
	rl.Resource().Attributes().PutStr("host.name", hostname)

	sl := rl.ScopeLogs().AppendEmpty()
	sl.SetSchemaUrl(schemaURL)
	sl.Scope().SetName(pluginName)
	sl.Scope().SetVersion(version)

	lr := sl.LogRecords().AppendEmpty()
	now := pcommon.NewTimestampFromTime(time.Now().UTC())
	lr.SetTimestamp(now)
	lr.SetObservedTimestamp(now)
	lr.SetSeverityNumber(plog.SeverityNumberInfo)
	lr.SetSeverityText(slog.LevelInfo.String())

	body, err := json.Marshal(msg)
	if err != nil {
		lr.Body().SetStr(fmt.Sprintf("%v", msg))
	} else {
		lr.Body().SetStr(string(body))
	}

	attrs := lr.Attributes()
	for k, v := range msg {
		switch val := v.(type) {
		case string:
			attrs.PutStr(k, val)
		case int:
			attrs.PutInt(k, int64(val))
		case int64:
			attrs.PutInt(k, val)
		case float64:
			attrs.PutDouble(k, val)
		case bool:
			attrs.PutBool(k, val)
		default:
			attrs.PutStr(k, fmt.Sprintf("%v", val))
		}
	}

	return logs
}
