# Prometheus Go Demo (All Standard Metrics)

A minimal Go service that exposes:
- **Go runtime metrics** (`go_*`)
- **Process metrics** (`process_*`)
- **Prometheus handler metrics** (`promhttp_*`)
- A few **app-specific metrics** (`app_*`) to make dashboards interesting

Includes a Docker Compose stack for **Prometheus (9090)** and **Grafana (3000)** with an auto-provisioned dashboard.

## Quick Start

### 1) Run the Go app

```bash
cd app
go mod tidy
go run .
# listens on :2112
```

Try hitting demo endpoints in another terminal:
```bash
curl -s localhost:2112/work > /dev/null
curl -s localhost:2112/alloc > /dev/null
curl -s localhost:2112/goroutines > /dev/null
```

### 2) Run Prometheus + Grafana

```bash
docker compose -f ops/docker-compose.yml up -d
# Prometheus: http://localhost:9090
# Grafana:    http://localhost:3000
```

Dashboard: **Go App Metrics (Prom Demo)** is auto-loaded.

## Notable Metrics

- Go runtime: `go_goroutines`, `go_threads`, `go_memstats_*`, `go_gc_duration_seconds*`
- Process: `process_cpu_seconds_total`, `process_resident_memory_bytes`, `process_open_fds`, `process_start_time_seconds`
- Prom handler: `promhttp_metric_handler_requests_total`, `promhttp_metric_handler_request_duration_seconds*`
- Build info: `go_build_info`
- App: `app_requests_total`, `app_inflight_requests`, `app_work_duration_seconds_*`

## Endpoints

- `/metrics` — Prometheus text exposition
- `/work` — CPU + sleep
- `/alloc` — temporary heap allocations
- `/goroutines` — spawns short-lived goroutines
- `/healthz` — healthcheck

## License

MIT (use freely for demos).