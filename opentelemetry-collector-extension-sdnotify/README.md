# opentelemetry-collector-extension-sdnotify

Minimal custom OpenTelemetry Collector distribution that includes a single
custom extension, **`sdnotify`**, which sends `READY=1` / `STOPPING=1` to
systemd via `sd_notify(3)` so the collector can be supervised as a
`Type=notify` service.

## Layout

```
.
├── Makefile
├── manifest.yml                    # OCB (builder) input
├── config.yaml                     # sample collector config (filelog -> debug)
├── internal/tools/go.mod           # pins the ocb tool version
└── extension/sdnotify/             # the custom extension (its own Go module)
    ├── go.mod
    ├── factory.go
    ├── config.go
    └── extension.go
```

The compiled binary is generated under `_build/` by `ocb` and built into
`bin/otelcol-sdnotify`. Both directories are git-ignored.

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

## Verify the dir → stdout pipeline

```sh
mkdir -p input
echo "hello sdnotify" > input/test.log
make run
```

You should see the line printed by the debug exporter, plus an
`sdnotify: NOTIFY_SOCKET not set; READY=1 was a no-op` log entry (because
you're running interactively, not under systemd). Add more lines to files
in `input/` to confirm tailing.

## Use under systemd

Drop a unit like:

```ini
[Service]
Type=notify
ExecStart=/usr/local/bin/otelcol-sdnotify --config=/etc/otelcol/config.yaml
```

systemd will set `NOTIFY_SOCKET`, the extension will send `READY=1` once
all components have started, and `systemctl status` will reflect the
`active (running)` state only after the collector is actually up.

## Versions

- OpenTelemetry Collector: `v0.154.0`
- Component / extension v1 APIs: `v1.60.0`
- ocb (builder): `v0.154.0`
- Go: `1.24+`

To upgrade, bump the versions in `manifest.yml`, `internal/tools/go.mod`,
and `extension/sdnotify/go.mod` together — they must move in lockstep.
