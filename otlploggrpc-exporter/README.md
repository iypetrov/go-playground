# otlploggrpc-exporter

A minimal Go application that emits OpenTelemetry log records via gRPC to a local OpenTelemetry Collector.

## Usage

```bash
# Build the binary
make build

# Start the OTel Collector (Docker) and run the binary
make run

# Clean up
make teardown
```

The app sends a log record every second to `localhost:4317` and shuts down gracefully on SIGINT/SIGTERM.

## Configuration

The OTel Collector config is in `otel-config.yaml`. The collector exposes:

- `:4317` — gRPC OTLP receiver
- `:8888` — metrics

## Metrics Architecture

There are three layers of metrics all funneled into a single Prometheus endpoint at :2021/metrics:

┌─────────────────────────────────────────────────────┐
│          HTTP :2021/metrics  (Prometheus scrape)     │
└──────────────────────┬──────────────────────────────┘
                       │
       ┌───────────────┼───────────────────┐
       │               │                   │
  ① Custom         ② gRPC client      ③ SDK self-
  CounterVec        metrics (auto)     instrumentation

---
1. Custom Prometheus Counter (metrics.go:25-34)

A classic prometheus/client_golang counter registered on the default registry via promauto:

OutputClientLogs = promauto.With(prometheus.DefaultRegisterer).NewCounterVec(...)

- Full metric name: fluentbit_example_output_client_logs_total
- Label: host
- Incremented in main.go each time a log record is forwarded to the OTLP collector.

---
2. OpenTelemetry SDK MeterProvider + Prometheus Exporter (metrics.go:41-70)

NewGlobalMetricsSetup() builds the OTEL metrics pipeline:

1. Prometheus exporter (otelexporterprom.New) — acts as an OTEL SDK metric reader that exposes OTEL-collected metrics in Prometheus format. It's attached to prometheus.DefaultRegisterer, so everything shows up on the same /metrics endpoint alongside the custom counter.
  - Namespace: output_plugin (all OTEL-sourced metrics get this prefix).
  - Translation strategy: UnderscoreEscapingWithSuffixes — converts OTEL metric names to Prometheus-compatible names.
2. MeterProvider (sdkmetric.NewMeterProvider) — created with the Prometheus reader. Then set as the global provider via otel.SetMeterProvider() so instrumentation libraries (like otelgrpc) can discover it automatically.

---
3. gRPC Client Metrics (metrics.go:80-84)

grpc.WithStatsHandler(otelgrpc.NewClientHandler(
    otelgrpc.WithMeterProvider(m.provider),
))

This returns a grpc.DialOption that instruments all gRPC client calls (the OTLP log export calls). It automatically records metrics like rpc.client.duration, rpc.client.request.size, etc. — all routed through the MeterProvider → Prometheus exporter.

Used in main.go:61:
otlploggrpc.WithDialOption(globalMetricsSetup.GRPCStatsHandler())

---
4. SDK Self-Instrumentation (metrics.go:46)

os.Setenv("OTEL_GO_X_OBSERVABILITY", "true")

This experimental flag tells the OTEL Log SDK to emit its own internal metrics (e.g., otel.sdk.log.created, otel.sdk.exporter.sent, otel.sdk.exporter.failed). Because the global MeterProvider is already wired to Prometheus, these also show up on :2021/metrics under the output_plugin_ namespace.

---
Lifecycle

- Startup: NewGlobalMetricsSetup() → HTTP server started on :2021
- Runtime: custom counter + automatic gRPC + SDK metrics all scraped from one endpoint
- Shutdown: GlobalMetricsSetup.Shutdown() (idempotent via sync.Once) flushes the meter provider; the HTTP server is also gracefully shut down on signal.

The key insight is that the Prometheus exporter bridges two worlds — native Prometheus counters (custom) coexist with OTEL-instrumented metrics (gRPC, SDK internals) on the same registry and endpoint.
