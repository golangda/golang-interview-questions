OTEL_EXPORTER_OTLP_ENDPOINT ?= localhost:4317

.PHONY: up down restart logs otel-logs topics producer processor retryworker deps clean

up:
	docker compose -f compose.yaml up -d

down:
	docker compose -f compose.yaml down -v

restart: down up

logs:
	docker compose -f compose.yaml logs -f

otel-logs:
	docker logs -f otel-collector

topics:
	go run ./cmd/admin

producer:
	OTEL_EXPORTER_OTLP_ENDPOINT=$(OTEL_EXPORTER_OTLP_ENDPOINT) go run ./cmd/producer

processor:
	OTEL_EXPORTER_OTLP_ENDPOINT=$(OTEL_EXPORTER_OTLP_ENDPOINT) go run ./cmd/processor

retryworker:
	OTEL_EXPORTER_OTLP_ENDPOINT=$(OTEL_EXPORTER_OTLP_ENDPOINT) go run ./cmd/retryworker

deps:
	go mod tidy

clean:
	docker compose -f compose.yaml down -v --remove-orphans || true
	docker rm -f kafka otel-collector || true
	docker volume prune -f || true
