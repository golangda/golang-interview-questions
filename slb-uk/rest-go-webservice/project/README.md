# Distributed Application: API → Kafka → Consumer → MySQL (SAGA Pattern)

## Overview

This project demonstrates a distributed, event-driven system using Go, Kafka, and MySQL, orchestrated with Kubernetes and following the SAGA pattern for transactional consistency.

**Architecture:**

```
[apisvc] → Kafka → [consumersvc] → MySQL
```

* **apisvc**: Exposes RESTful CRUD APIs for messages, publishes commands to Kafka, and provides operation result queries.
* **consumersvc**: Listens to Kafka commands, performs MySQL DB operations, and sends acknowledgements back via Kafka.

## Tech Stack

* Go (apisvc & consumersvc)
* Kafka (message broker)
* MySQL (persistence)
* Docker (containerization)
* Minikube (local K8s)
* Makefile (build & deploy automation)

## Local Development

### Prerequisites

* Docker
* Minikube
* kubectl
* golangci-lint (optional for linting)

### Build & Deploy

```bash
minikube start
make dev-up
```

This will:

1. Build `apisvc`, `consumersvc`, and MySQL images.
2. Load images into Minikube.
3. Apply Kubernetes manifests.

### Port Forward API Service

```bash
kubectl port-forward svc/apisvc 8080:80
```

## API Usage

### Create Message

```bash
curl -X POST localhost:8080/v1/messages \
  -H 'Content-Type: application/json' \
  -d '{"message":"hello world"}'
# => {"trace_id":"<uuid>","status":"PENDING"}
```

### Get Operation Result

```bash
curl localhost:8080/v1/operations/<trace_id>
```

### Read / Update / Delete Message

```bash
curl localhost:8080/v1/messages/1
curl -X PUT localhost:8080/v1/messages/1 -H 'Content-Type: application/json' -d '{"message":"updated"}'
curl -X DELETE localhost:8080/v1/messages/1
```

## Kubernetes Manifests

### `k8s/apisvc.yaml`

Defines Deployment and Service for `apisvc`.

### `k8s/consumersvc.yaml`

Defines Deployment for `consumersvc`.

Ensure you delete legacy `k8s/api.yaml` and `k8s/consumer.yaml`.

## Makefile Targets

* `make build` – Compile Go binaries.
* `make docker` – Build Docker images.
* `make minikube-load` – Load images into Minikube.
* `make k8s-apply` – Apply all K8s manifests.
* `make dev-up` – Build, load, and deploy.
* `make dev-down` – Remove all K8s resources.
* `make logs-apisvc` / `make logs-consumersvc` – Stream logs.

## Paths & Images

* Service code: `cmd/apisvc/`, `cmd/consumersvc/`
* Dockerfiles: in each service directory
* Images: `apisvc:local`, `consumersvc:local`

## Notes

* The system implements checkpoints and transactions for reliability.
* The SAGA pattern ensures eventual consistency across services.
* Swagger v2.x docs are exposed via `apisvc` endpoint for API exploration.
