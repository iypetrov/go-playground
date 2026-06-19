# daemon

A small, complete reference for a systemd-aware Go service. Demonstrates:

- `Type=notify` lifecycle (`READY=1`, `RELOADING=1`, `STOPPING=1`, `STATUS=...`)
- Watchdog pings (`WATCHDOG=1`) gated on a real health flag
- `SIGHUP` reload, `SIGTERM`/`SIGINT` graceful drain

## Run standalone (no systemd)

```sh
make build run
curl localhost:8081/helathz
```

`NOTIFY_SOCKET` is unset outside systemd, so the notify/watchdog calls are
no-ops — the program just runs as a plain HTTP server.

## Run under systemd (Linux host or VM)

```sh
# 0. build the binary
make build

# 1. install the binary somewhere ExecStart= can find it
install -m 0755 ./bin/main /usr/local/bin/myservice

# 2. drop the unit in place
cp main.service /etc/systemd/system/myservice.service
systemctl daemon-reload

# 3. start
systemctl enable myservice.service
```

## Things to actually try

These are the experiments that make the concepts click.

**Observe `READY=1` gating.** `main.go` sleeps 2s before notifying ready.
`systemctl start myservice` blocks for those 2s — that's `Type=notify` working.

```sh
systemctl start myservice
systemctl status myservice    # STATUS= line shows our message
```

**Watch the watchdog kill us.** Use `curl http://localhost:8081/switch` to 
make the systemd service go in unhealthy state. systemd kills the process within
`WatchdogSec` and `Restart=on-failure` brings it back.

```sh
journalctl -u myservice -f
# expect: "Watchdog timeout (limit 10s)!" then a fresh start
```

**Trigger a reload.**

```sh
systemctl reload myservice
# journal shows: SIGHUP received, reloading -> reload done
# `systemctl reload` blocks until READY=1 comes back
```

**Test graceful drain.** Hold a slow request open while stopping:

```sh
curl localhost:8081/slow &      # 10s response
systemctl stop myservice        # waits for the curl to finish (within TimeoutStopSec)
```
