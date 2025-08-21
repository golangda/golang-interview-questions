# Kafka + Go (Sarama) with Retry, DLQ and OpenTelemetry

This is a ready‑to‑run demo showing:

- Go + **IBM Sarama** client
- **Consumer Group** processor
- **Staged Retry topics** + **DLQ** pattern
- **OpenTelemetry** tracing via `otelsarama` to an **OTel Collector**

## Prerequisites
- Docker 24+
- Go 1.22+

## Quickstart

```bash
# 1) Start infra (Kafka KRaft + OTEL Collector)
make up

# 2) Create topics
make topics

# 3) Start the processor and retry worker (two terminals)
make processor
make retryworker

# 4) Produce some messages
make producer

# 5) Watch OpenTelemetry spans in collector logs
make otel-logs
```

### What you’ll see
- Messages like `ok: welcome` are processed immediately by the **processor**.
- Messages like `fail: simulate downstream error` go to `events.v1.retry.5s`,
  then re-queued to `events.v1`. If they still fail, they progress to
  `events.v1.retry.30s`, then `events.v1.retry.2m`, and finally to **DLQ**.

### Topics used
- `events.v1` (main)  
- `events.v1.retry.5s`, `events.v1.retry.30s`, `events.v1.retry.2m` (retry stages)  
- `events.v1.dlq` (dead-letter)

> **Note**: For local dev we use replication factor `1`. In production use `>=3`.

## Make targets

- `make up` / `make down` – start/stop Kafka + OTEL
- `make topics` – creates main, retry, and DLQ topics
- `make processor` – runs the consumer group processor
- `make retryworker` – runs the retry worker (re-queues after a delay)
- `make producer` – sends demo messages
- `make otel-logs` – tails collector logs
- `make clean` – remove containers/volumes/images (careful)

## Structure
```
cmd/
  admin/         # topic creation
  producer/      # demo producer
  processor/     # consumer group processor with retry->DLQ
  retryworker/   # consumes retry topics, sleeps, re-queues to main
internal/
  retry/         # retry stages + headers
  tracing/       # OTel bootstrap + Kafka header propagation helper
compose.yaml     # Kafka (KRaft) + OTel Collector
otel-collector-config.yaml
```

## Notes
- The **OTLP endpoint** defaults to `localhost:4317`. You can override with `OTEL_EXPORTER_OTLP_ENDPOINT` env var.
- For Docker networking on non-Linux hosts, we expose Kafka on `localhost:9092` and also provide an internal broker listener `kafka:9093` for containers.
