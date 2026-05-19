# ploggrpc-exporter-component

A minimal Go application that builds OTLP log records using `pdata/plog` and ships them through the Collector's `otlpexporter` component embedded as a library — getting queue, retry, batching, and timeout for free from `exporterhelper`.


|                                         | ** `otlploggrpc`**                                                                | **`ploggrpc`**                                                                             | **`otlpexporter`**                                                                                                                        |
| --------------------------------------- | --------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------ | ----------------------------------------------------------------------------------------------------------------------------------------- |
| **Top-level lib**                       | OTEL Go **SDK** (`sdklog`)                                                        | Collector **`pdata/plog`** + `plogotlp.GRPCClient`                                         | Collector **`exporter/otlpexporter`** component                                                                                           |
| **Record type**                         | `otlplog.Record` (SDK)                                                            | `plog.LogRecord` (pdata)                                                                   | `plog.LogRecord` (pdata)                                                                                                                  |
| **Wire export**                         | `otlploggrpc.Exporter` inside SDK                                                 | `plogotlp.GRPCClient.Export` direct gRPC unary                                             | `exporter.Logs.ConsumeLogs` → `exporterhelper` → `plogotlp` under the hood                                                                |
| **Batching**                            | `sdklog.BatchProcessor` (configurable: queue size, batch size, interval, timeout) | None — DIY loop, or wrap with `exporterhelper`                                             | Built in via `QueueBatchConfig` (queue + batcher merged); enabled by default                                                              |
| **Retry on failure**                    | Built into `otlploggrpc` exporter                                                 | None — DIY                                                                                 | `RetryConfig` (`configretry.BackOffConfig`), enabled by default                                                                           |
| **Bounded queue / backpressure**        | Yes (`BatchProcessor` queue)                                                      | None — DIY                                                                                 | Yes (`QueueConfig` via `configoptional.Optional[QueueBatchConfig]`)                                                                       |
| **Per-export timeout**                  | `WithExportTimeout(d)`                                                            | DIY (`context.WithTimeout` per call)                                                       | `TimeoutConfig.Timeout` (default 5s)                                                                                                      |
| **Custom record processing**            | Implement `sdklog.Processor`, chain in front of batch                             | DIY wrapper around `Export`                                                                | Implement Collector `processor.Logs`, or wrap `ConsumeLogs`                                                                               |
| **`Emit`/`Export`/`Consume` semantics** | `Logger.Emit` returns immediately (queued)                                        | `Export` blocks until server ACKs                                                          | `ConsumeLogs` returns when item is enqueued; `Shutdown(ctx)` drains                                                                       |
| **Self-instrumentation metrics**        | Yes (`OTEL_GO_X_OBSERVABILITY=true` → `otel.sdk.log.*`, `otel.sdk.exporter.*`)    | None                                                                                       | `exporterhelper` metrics if you wire `TelemetrySettings.MeterProvider`; otherwise none                                                    |
| **gRPC-client metrics**                 | `otelgrpc` stats handler via `otlploggrpc.WithDialOption`                         | `otelgrpc` stats handler on raw `*grpc.ClientConn`                                         | Goes through the exporter's own connection — to instrument it you'd configure `cfg.ClientConfig` rather than dial yourself                |
| **Insecure / TLS config**               | `otlploggrpc.WithInsecure()`                                                      | `insecure.NewCredentials()` on dial                                                        | `cfg.ClientConfig.Endpoint = "..."`, `cfg.ClientConfig.Insecure = true`                                                                   |
| **Lifecycle**                           | `LoggerProvider.Shutdown` flushes batch                                           | Manual `conn.Close()`                                                                      | `Start(ctx, host)` → `ConsumeLogs` → `Shutdown(ctx)` (drains queue)                                                                       |
| **Lines of glue code**                  | ~50 (exporter + provider + processor)                                             | ~30 (dial + client + loop)                                                                 | ~40 (factory + config + settings + nop host)                                                                                              |
| **Dep weight**                          | OTEL SDK + log SDK + otlploggrpc                                                  | pdata only + grpc-go                                                                       | pdata + collector core + exporterhelper + otlpexporter                                                                                    |
| **You write a `plog.Logs`?**            | No — SDK builds protobuf from `otlplog.Record`                                    | Yes                                                                                        | Yes                                                                                                                                       |
| **Best for**                            | Apps that already use the OTEL Go ecosystem; want a familiar `Logger.Emit()` API  | Lightweight forwarders / fluent-bit-style plugins; want full control over the wire payload | Embedding "a real collector pipeline" in your binary; want production behaviors (queue/retry/batch/timeout) without rolling them yourself |

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
