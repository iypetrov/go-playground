# ploggrpc-exporter-component

A minimal Go application that builds OTLP log records using `pdata/plog` and ships them through the Collector's `otlpexporter` component embedded as a library — getting queue, retry, batching, and timeout for free from `exporterhelper`.

## What it does

Constructs a `plog.Logs` payload directly and pushes it every second to `localhost:4317` via `otlpexporter` (factory → `CreateLogs` → `Start` → `ConsumeLogs`). Shuts down gracefully on SIGINT/SIGTERM, draining the queue.

Differences vs. the other two variants in this repo:

- **vs. `otlploggrpc-exporter`**: no OTEL Log SDK, no `LoggerProvider`. Records are built with `plog`, not `otlplog.Record`.
- **vs. `ploggrpc-exporter`**: doesn't call `plogotlp.GRPCClient.Export` directly. Instead pushes through the full `otlpexporter` component, which wraps `exporterhelper` and gives queue + retry + batch + timeout out of the box.

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

The app exposes a single Prometheus endpoint at `http://localhost:2021/metrics`. Two sources are merged onto one explicit `*prometheus.Registry`:

1. A custom counter `fluentbit_example_output_client_logs_total{host="..."}` incremented per enqueued log record.
2. `exporterhelper`'s built-in metrics, routed there by handing our `MeterProvider` to `set.TelemetrySettings.MeterProvider`. Expect series under the `output_plugin_` namespace covering things like:
   - sent / send-failed / enqueue-failed log records
   - queue size and capacity
   - in-flight requests

   (Exact names depend on the collector version's telemetry builder — `curl http://localhost:2021/metrics | grep output_plugin_` to see what's emitted.)

There is no `otelgrpc` stats handler in this variant — the gRPC client is owned by the exporter component, not by us.
