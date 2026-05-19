# otlploggrpc-exporter

A minimal Go application that emits OpenTelemetry log records via gRPC to a local OpenTelemetry Collector.

## What it does

Sends a log record every second to `localhost:4317` (OTLP gRPC)>

## How to run

```bash
# Build the binary
make build

# Start the OTel Collector (Docker) and run the binary
make run

# Clean up
make teardown
```

The collector config lives in `otel-config.yaml` and exposes `:4317` (OTLP gRPC) and `:8888` (collector metrics).

## Metrics

The app exposes a single Prometheus endpoint at `http://localhost:2021/metrics`. Three sources are merged onto one explicit `*prometheus.Registry`:

1. A custom counter `fluentbit_example_output_client_logs_total{host="..."}` incremented per forwarded log record.
2. gRPC client metrics from `otelgrpc` (e.g. `rpc.client.duration`).
3. OTEL Log SDK self-instrumentation (`otel.sdk.log.*`, enabled via `OTEL_GO_X_OBSERVABILITY=true`).

OTEL-sourced metrics are exposed via the Prometheus exporter under the `output_plugin_` namespace.
