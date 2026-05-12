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
