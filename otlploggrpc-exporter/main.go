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
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	otlplog "go.opentelemetry.io/otel/log"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	sdkresource "go.opentelemetry.io/otel/sdk/resource"

	semconv "go.opentelemetry.io/otel/semconv/v1.27.0"
)

var (
	pluginName = "fluent-bit-output-plugin"
	// Using the logs schema version that matches the SDK
	schemaURL               = "https://opentelemetry.io/schemas/1.27.0"
	version                 = "0.1.0"
	endpoint                = "localhost:4317"
	hostname                = "some-random-vm"
	batchMaxQueueSize       = 1000
	batchExportTimeout      = 15 * time.Minute
	eventGenerationInterval = 1 * time.Second
)

func main() {
	// Classic setup for Go projects - context, logger, config, etc.
	ctx, cancel := context.WithCancel(context.Background())
	loggerOpts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}
	loggerHandler := slog.NewTextHandler(os.Stderr, loggerOpts)
	logger := logr.FromSlogHandler(loggerHandler)

	// Create blocking OTLP gRPC exporter
	exporterOpts := []otlploggrpc.Option{
		otlploggrpc.WithEndpoint(endpoint),
		otlploggrpc.WithInsecure(),
	}
	exporter, err := otlploggrpc.New(ctx, exporterOpts...)
	if err != nil {
		cancel()
		logger.Info("failed to create OTLP gRPC exporter: %w", err)
	}

	// Create batch processor
	batchOpts := []sdklog.BatchProcessorOption{
		sdklog.WithMaxQueueSize(batchMaxQueueSize),
		sdklog.WithExportTimeout(batchExportTimeout),
	}
	batchProcessor := sdklog.NewBatchProcessor(exporter, batchOpts...)

	// Build resource attributes
	resourceAttrs := []attribute.KeyValue{
		semconv.HostName(hostname),
	}
	resource := sdkresource.NewWithAttributes(
		semconv.SchemaURL,
		resourceAttrs...,
	)

	// Create logger provider
	loggerProvider := sdklog.NewLoggerProvider(
		sdklog.WithResource(resource),
		sdklog.WithProcessor(batchProcessor),
	)

	// Build instrumentation scope options
	scopeOptions := []otlplog.LoggerOption{
		otlplog.WithInstrumentationVersion(version),
		otlplog.WithSchemaURL(schemaURL),
	}

	// Graceful shutdown on SIGINT/SIGTERM
	var wg sync.WaitGroup
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// Log emitter
	wg.Go(func() {
		ticker := time.NewTicker(eventGenerationInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				msg := map[string]any{
					"name": "John",
				}
				record := buildRecord(msg)
				loggerProvider.Logger(pluginName, scopeOptions...).
					Emit(ctx, record)
				logger.Info("log record emitted", "body", record.Body())
				OutputClientLogs.WithLabelValues(hostname).Inc()
			case <-ctx.Done():
				return
			}
		}
	})

	// Metric server
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	srv := &http.Server{
		Addr:    ":2021",
		Handler: mux,
	}
	wg.Go(func() {
		go func() {
			<-ctx.Done()

			shutdownCtx, cancel := context.WithTimeout(
				context.Background(),
				5*time.Second,
			)
			defer cancel()

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

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	if err := loggerProvider.Shutdown(shutdownCtx); err != nil {
		logger.Error(err, "failed to shutdown logger provider")
	}
}

func buildRecord(msg map[string]any) otlplog.Record {
	var record otlplog.Record
	body, err := json.Marshal(msg)
	if err != nil {
		record.SetBody(otlplog.StringValue(fmt.Sprintf("%v", msg)))
	} else {
		record.SetBody(otlplog.StringValue(string(body)))
	}
	for k, v := range msg {
		switch val := v.(type) {
		case string:
			record.AddAttributes(otlplog.String(k, val))
		case int:
			record.AddAttributes(otlplog.Int64(k, int64(val)))
		case int64:
			record.AddAttributes(otlplog.Int64(k, val))
		case float64:
			record.AddAttributes(otlplog.Float64(k, val))
		case bool:
			record.AddAttributes(otlplog.Bool(k, val))
		default:
			record.AddAttributes(otlplog.String(k, fmt.Sprintf("%v", val)))
		}
	}
	record.SetSeverityText(slog.LevelInfo.String())
	record.SetTimestamp(time.Now().UTC())
	return record
}
