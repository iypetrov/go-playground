# opentelemetry-collector-extension-sdnotify

Custom OpenTelemetry Collector distribution that bundles a single custom
extension, **`sdnotify`**, which talks to systemd via `sd_notify(3)` so the
collector can be supervised as a `Type=notify` service.

## What it does

The extension implements four collector interfaces to surface readiness and
component health back to systemd:

| Interface | Purpose |
|---|---|
| `extension.Extension` | Lifecycle (`Start`/`Shutdown`). Sends `STOPPING=1` on shutdown; in watchdog mode, runs the keepalive ticker. |
| `extensioncapabilities.PipelineWatcher` | `Ready()` sends `READY=1` once every receiver/processor/exporter in every pipeline has started. `NotReady()` stops the gRPC watcher before receivers drain. |
| `extensioncapabilities.Dependent` | When deep-healthcheck mode is on, declares a dependency on the `healthcheckv2` extension so it starts first. |
| `componentstatus.Watcher` | Receives every component status change; in deep mode, permanent/fatal errors are immediately surfaced as `STATUS=component <id> failed: <err>`. |

Two modes — both configurable, both optional:

- **Shallow (default).** Only `READY=1` (from `Ready()`) and `STOPPING=1`
  (from `Shutdown`) are sent. Component status changes are logged but not
  forwarded to systemd.
- **Deep healthcheck (`deep_healthcheck: true`).** sdnotify subscribes to
  the `grpc.health.v1.Health/Watch` stream exposed by the
  [`healthcheckv2`](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/extension/healthcheckv2extension)
  extension and emits a `STATUS=<line>` notification on every aggregated
  status change. Hard failures observed via `ComponentStatusChanged` are
  pushed immediately so `systemctl status` reflects degradation within
  milliseconds.

The **watchdog** pinger auto-enables whenever systemd sets `WATCHDOG_USEC`
(typically via `WatchdogSec=` in the unit file). It sends `WATCHDOG=1` every
`WATCHDOG_USEC/2` microseconds; no configuration knob needed — if systemd
doesn't ask for a watchdog, the pinger stays off.

## Configuration

| Key | Type | Default | Notes |
|---|---|---|---|
| `fail_if_not_supervised` | bool | `false` | Fail `Start` when `NOTIFY_SOCKET` is unset. |
| `unset_environment` | bool | `false` | Pass `unsetEnv=true` to `daemon.SdNotify` so child processes don't inherit `NOTIFY_SOCKET`. |
| `deep_healthcheck` | bool | `false` | Subscribe to `healthcheckv2`'s gRPC `Health.Watch` (overall collector health, i.e. the aggregate of every component) and emit `STATUS=...`. |
| `healthcheckv2` | component.ID | `healthcheckv2` | ID of the sibling healthcheckv2 extension. Used both for `Dependencies()` and to look up the gRPC endpoint. Override only if your sibling extension uses a non-default instance name. |
| `healthcheckv2_grpc_endpoint` | string | `""` | Optional explicit override. When empty, the endpoint is read from the sibling extension's config at runtime. |

Example with deep mode wired up:

```yaml
extensions:
  healthcheckv2:
    use_v2: true
    grpc:
      endpoint: localhost:13132

  sdnotify:
    deep_healthcheck: true

service:
  extensions: [healthcheckv2, sdnotify]
  pipelines:
    logs:
      receivers: [file_log]
      exporters: [debug]
```

### A note on endpoint resolution

In deep mode, sdnotify reads the gRPC endpoint off the sibling
`healthcheckv2` extension via reflection on its config struct. The healthcheckv2
type layout is not part of its public Go API, so if a future version reshapes
its config the lookup may fail — in that case, set `healthcheckv2_grpc_endpoint`
explicitly to bypass it.

## Build & run

```sh
make build              # generate _build/ via ocb, then `go build`
make run                # build + run with config.yaml (watches ./input/)
make test               # `go test ./...` inside the extension module
make clean
```

First-time setup also requires populating the tools module:

```sh
cd internal/tools && go mod tidy
```

## Verify locally (without systemd)

Run the collector under a `socat` listener that pretends to be systemd:

```sh
# Terminal 1: fake systemd notify socket.
export NOTIFY_SOCKET=/tmp/sdn.sock
rm -f "$NOTIFY_SOCKET"
socat UNIX-RECV:$NOTIFY_SOCKET -

# Terminal 2: run the collector with the same NOTIFY_SOCKET.
export NOTIFY_SOCKET=/tmp/sdn.sock
make run
```

You should see in terminal 1:

```
READY=1
STATUS=SERVING        # one per healthcheckv2 status change (deep mode)
...
STOPPING=1
```

To exercise the deep-healthcheck push path, point a receiver/exporter at an
unreachable endpoint so it transitions to `StatusPermanentError`; a
`STATUS=component <id> failed: ...` line will appear before the next
aggregated update from healthcheckv2.

## Use under systemd

```ini
[Unit]
Description=OpenTelemetry Collector (sdnotify)
After=network-online.target
Wants=network-online.target

[Service]
Type=notify
NotifyAccess=main
ExecStart=/usr/local/bin/otelcol-sdnotify --config=/etc/otelcol/config.yaml
# Enables systemd's watchdog supervision. sdnotify auto-detects this and
# starts pinging WATCHDOG=1 at WatchdogSec/2 cadence.
WatchdogSec=30s
Restart=on-failure
RestartSec=5s

[Install]
WantedBy=multi-user.target
```

`systemctl status otelcol-sdnotify` will then show:

- `Active: active (running)` only after `Ready()` has fired (i.e. all
  pipelines are up).
- A live `Status:` line reflecting the current aggregated health (in deep
  mode), updated as healthcheckv2 reports changes.

## Versions

- OpenTelemetry Collector: `v0.154.0`
- Component / extension v1 APIs: `v1.60.0`
- `healthcheckv2extension`: `v0.154.0` (contrib)
- ocb (builder): `v0.154.0`
- Go: `1.24+`

To upgrade, bump the versions in `manifest.yml`, `internal/tools/go.mod`,
and `extension/sdnotify/go.mod` together — they must move in lockstep.
