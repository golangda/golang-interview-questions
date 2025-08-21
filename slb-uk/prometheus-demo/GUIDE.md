# GUIDE.md — Deep Dive & Troubleshooting

## How standard metrics are registered

We explicitly register:
- `prometheus.NewGoCollector()` → `go_*` runtime metrics
- `prometheus.NewProcessCollector(...)` → `process_*` metrics
- `prometheus.NewBuildInfoCollector()` → `go_build_info`

The `/metrics` handler is wrapped with `promhttp.InstrumentMetricHandler(...)` so you also get:
- `promhttp_metric_handler_requests_total`
- `promhttp_metric_handler_requests_in_flight`
- `promhttp_metric_handler_request_duration_seconds_*`

## Validating Metrics

```bash
curl -s localhost:2112/metrics | grep -E '^(go_|process_|promhttp_|go_build_info|app_)' | head -n 40
```

Trigger load:
```bash
for i in {1..50}; do curl -s localhost:2112/work >/dev/null; done
curl -s localhost:2112/alloc >/dev/null
curl -s localhost:2112/goroutines >/dev/null
```

## Prometheus

Open http://localhost:9090 and try:
- `go_goroutines`
- `rate(go_gc_duration_seconds_sum[5m]) / rate(go_gc_duration_seconds_count[5m])`
- `rate(process_cpu_seconds_total[5m])`

## Grafana

Open http://localhost:3000 (anonymous access enabled). The "Go App Metrics (Prom Demo)" dashboard is preloaded. Panels include:
- Goroutines (stat)
- GC Avg Duration (stat)
- Process RSS (timeseries)
- CPU Seconds Rate (timeseries)
- Scrape Requests by Code (`promhttp_*`)
- App Requests by Path (`app_*`)
- App Work Duration Histogram (`app_work_duration_seconds_bucket`)

## Common Pitfalls

1. **No targets found in Prometheus**
   - Ensure the app is running on `:2112`.
   - The `prometheus.yml` includes both `localhost:2112` and `host.docker.internal:2112` for cross-platform. On Linux without Docker Desktop, `host.docker.internal` may not resolve—`localhost` scrape should still work from the container if you run with host networking. If needed, change targets to your host IP or run Prometheus with `--net=host` (Linux).

2. **Grafana shows empty panels**
   - Confirm datasource is healthy (Configuration → Data sources → Prometheus).
   - Verify Prometheus is scraping (Targets page).

3. **High GC or memory**
   - Hit `/alloc` repeatedly; observe `go_memstats_*`, `go_gc_duration_seconds*` rising.
   - Investigate using pprof (not included in this demo), or adjust workload.


## Extending

- Add alerts in `ops/alerts.yml` and reference via `rule_files` in `prometheus.yml`.
- Add recording rules to precompute expensive expressions.
- Add service-specific metrics in code (counters/histograms/gauges/summaries).

Happy hacking!