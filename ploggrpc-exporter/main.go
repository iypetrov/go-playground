package main

import (
	"context"
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
	// metricsSetup, err := NewMetricsSetup(registry)
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
	srv := &http.Server{
		Addr:    ":2021",
		Handler: mux,
	}

	// TODO: real logic

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
				// TODO: send message
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

	// TODO: drain the exporter queue. Use a real timeout so queued items get flushed instead of dropped.
}
