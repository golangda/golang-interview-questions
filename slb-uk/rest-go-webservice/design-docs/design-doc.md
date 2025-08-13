# Distributed Application Design: API Server → Kafka → Consumer Service → MySQL

> **Tech stack:** Go • Kafka • MySQL • Docker • Minikube • Makefile
> **Patterns:** REST + Async command processing • Choreography-based **SAGA** • Idempotency • Exactly-once effect (practical) via idempotent keys + transactional DB writes • End-to-end tracing

---

## 1) Executive Summary

This system decouples synchronous API requests from database side-effects using Kafka. The **API Server** accepts CRUD requests for `Message` resources, attaches a `trace_id`, and publishes *commands* to Kafka. The **Consumer Service** processes commands, performs MySQL transactions, and emits *events/acks* back to Kafka. The API Server exposes a **Results API** that reads these acks from Kafka and serves them to clients (long-poll or immediate if cached). The overall flow is coordinated using a **choreography SAGA** (success and compensating events), with strong observability (OpenTelemetry), safe retry semantics, and idempotent processing.

---

## 2) Architecture Overview

### 2.1 Component Diagram

```mermaid
flowchart LR
  subgraph Client
    A[REST Client]
  end

  subgraph API[API Server (Go)]
    A1[/REST: Message CRUD/]
    A2[(Result Cache - in memory with TTL)]
    A3[[Kafka Producer (idempotent)]]
    A4[[Kafka Consumer (acks/events)]]
    A5[[Swagger v2 docs]]
  end

  subgraph BUS[Kafka]
    T1[(messages.commands)]
    T2[(messages.events)]
    T3[(messages.acks)]  %% optional split or reuse events
  end

  subgraph CS[Consumer Service (Go)]
    C1[[Kafka Consumer (commands)]]
    C2[(Idempotency Store in DB)]
    C3[[MySQL Tx]]
    C4[[Kafka Producer (acks/events)]]
  end

  subgraph DB[MySQL]
    D1[(messages)]
    D2[(saga_log)]
    D3[(idempotency_keys)]
  end

  A -->|HTTP CRUD| A1 -->|publish command| A3 --> T1
  C1 -->|consume| T1 --> C1
  C1 --> C3
  C3 -->|insert/update/delete| D1
  C3 -->|append| D2
  C3 -->|insert| D2
  C1 --> C4 --> T2
  C1 --> C4 --> T3
  A4 -->|consume acks| T3
  A4 --> A2
  A -->|GET /operations/{trace_id}| A1 --> A2
  A1 --> A5
```

### 2.2 Topics & Message Streams

* **`messages.commands`**: API → Consumer. CRUD commands with intent.
* **`messages.events`**: Consumer → (optionally other services). Domain events (created/updated/deleted/failure).
* **`messages.acks`**: Consumer → API. Operation results/acks (status, error).

> For a minimal setup, `messages.events` and `messages.acks` can be the same topic with different `type` fields.

---

## 3) Data Model & Contracts

### 3.1 Domain Model

```go
type Message struct {
  ID      int64  `json:"id"`
  Message string `json:"message"`
  // server-generated fields
  CreatedAt time.Time `json:"created_at,omitempty"`
  UpdatedAt time.Time `json:"updated_at,omitempty"`
}
```

### 3.2 Command Envelope (API → Kafka: `messages.commands`)

```json
{
  "trace_id": "uuid-v4",
  "correlation_id": "uuid-v4",     // = trace_id unless multi-step
  "timestamp": "RFC3339",
  "command": "Create|Update|Delete|Read",
  "resource": "Message",
  "payload": {
    "id": 123,                      // optional for create
    "message": "hello world"
  },
  "metadata": {
    "api_version": "v1",
    "user_id": "optional",
    "idempotency_key": "uuid-v4"    // also set as Kafka message key
  }
}
```

* **Kafka key:** `idempotency_key` (or `payload.id` for strong ordering per resource).
* **Headers:** `trace_id`, `correlation_id`, `command`, `resource`.

### 3.3 Ack/Event Envelope (Consumer → Kafka: `messages.acks`/`messages.events`)

```json
{
  "trace_id": "uuid-v4",
  "correlation_id": "uuid-v4",
  "timestamp": "RFC3339",
  "status": "SUCCESS|FAILURE",
  "event": "MessageCreated|MessageUpdated|MessageDeleted|MessageRead",
  "payload": {
    "message": { "id": 123, "message": "hello world" }
  },
  "error": {
    "code": "DB_CONFLICT|VALIDATION|NOT_FOUND|INTERNAL",
    "detail": "optional"
  }
}
```

* **Kafka key:** same as command key for easy join.
* **Headers:** mirror `trace_id`, `correlation_id`, `status`, `event`.

---

## 4) API Server (Go)

### 4.1 Responsibilities

1. **Expose REST CRUD** for `Message`.
2. **Create context + `trace_id`** for each request; propagate via Kafka headers.
3. **Publish** a command to Kafka (idempotent producer).
4. **Results API** to **communicate request results** by reading the ack stream:

   * Consume `messages.acks` in a background consumer group.
   * Cache acks by `trace_id` in an **in-memory TTL cache** (e.g., `ristretto` or `sync.Map` + timers).
   * Provide **long-polling** endpoint to wait for result with timeout (e.g., 10–20s).

> Note: For multi-replica production, replace in-memory cache with Redis or a compacted Kafka results topic re-played per trace.

### 4.2 REST Endpoints

* `POST /v1/messages` → Create (async)
* `GET  /v1/messages/{id}` → Read (async)
* `PUT  /v1/messages/{id}` → Update (async)
* `DELETE /v1/messages/{id}` → Delete (async)
* `GET  /v1/operations/{trace_id}` → Poll result/ack for that operation
* Swagger docs:

  * `GET /swagger/doc.json` (OpenAPI/Swagger v2.x JSON)
  * `GET /docs` (Swagger UI)

**Request flow (Create):**

```mermaid
sequenceDiagram
  participant Client
  participant API as API Server
  participant K as Kafka (commands)
  participant CS as Consumer Service
  participant DB as MySQL
  participant KA as Kafka (acks)

  Client->>API: POST /v1/messages {message:"hi"}
  API->>API: Generate trace_id, idempotency_key
  API->>K: Produce Command(Create, headers: trace_id,...)
  API-->>Client: 202 Accepted {"trace_id":"..."}
  CS->>K: Consume Command(Create)
  CS->>DB: Tx: INSERT message; write saga_log; commit
  CS->>KA: Produce Ack SUCCESS
  API->>KA: Background consumer receives ack
  API->>Client: GET /v1/operations/{trace_id} -> SUCCESS payload
```

### 4.3 Tracing & Context

* Use **OpenTelemetry** (`trace_id` set in request context).
* Propagate `trace_id` in:

  * HTTP response body (`202 Accepted` with `trace_id`)
  * Kafka message headers
  * Logs (structured)

### 4.4 API Response Patterns

* **On command enqueue:** `202 Accepted`
  Body: `{"trace_id":"<uuid>", "status":"PENDING"}`
  Optionally include `operation_url: "/v1/operations/<trace_id>"`.

* **Results API (polling):**

  * `200 OK` with final `{status, event, payload | error}` if present
  * `204 No Content` if still pending (or `408` if using long-poll timeout)
  * `410 Gone` if TTL expired and result not retained

### 4.5 Swagger v2.x

Minimal example (served at `/swagger/doc.json`):

```json
{
  "swagger": "2.0",
  "info": { "title": "Message API", "version": "1.0.0" },
  "basePath": "/v1",
  "paths": {
    "/messages": {
      "post": {
        "summary": "Create message (async)",
        "parameters": [{ "in":"body", "name":"body", "required":true,
          "schema":{"type":"object","properties":{"message":{"type":"string"}}}}],
        "responses": { "202":{"description":"Accepted"} }
      }
    },
    "/messages/{id}": { "get": { "responses": {"202":{"description":"Accepted"}} } },
    "/operations/{trace_id}": { "get": { "responses": {"200":{"description":"Result"}} } }
  }
}
```

> Implementation tip: use `swaggo/swag` (v2 spec) + `swag init` to autogenerate.

---

## 5) Consumer Service (Go)

### 5.1 Responsibilities

1. **Consume** from `messages.commands` (consumer group `"message-worker"`).
2. **Validate** and **route** commands (Create/Read/Update/Delete).
3. Execute **MySQL transactions** with proper **isolation** (REPEATABLE READ) and **retry** on transient errors.
4. **Idempotency & Checkpoints**:

   * Maintain `idempotency_keys` table keyed by `idempotency_key` and last processed command.
   * Commit DB transaction **before** committing offsets.
5. **Emit acks/events** to Kafka **after** a successful commit, with deduplication.

### 5.2 Tables

```sql
CREATE TABLE messages (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  message TEXT NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE saga_log (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  trace_id CHAR(36) NOT NULL,
  step VARCHAR(64) NOT NULL,           -- e.g., "CreateMessage"
  status ENUM('PENDING','SUCCESS','FAILURE') NOT NULL,
  error_code VARCHAR(64),
  error_detail TEXT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE idempotency_keys (
  idempotency_key CHAR(36) PRIMARY KEY,
  processed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  last_status ENUM('SUCCESS','FAILURE') NOT NULL,
  trace_id CHAR(36) NOT NULL
);

-- Optional: outbox for 2-phase emit if you want stronger safety
CREATE TABLE outbox (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  aggregate_key VARCHAR(128) NOT NULL,
  topic VARCHAR(128) NOT NULL,
  payload JSON NOT NULL,
  headers JSON NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  dispatched BOOLEAN DEFAULT FALSE
);
```

### 5.3 Transaction & Offset Commit Order

1. Begin **DB Tx**.
2. Apply CRUD operation with validations (lock rows for Update/Delete as needed).
3. Upsert into `idempotency_keys` (`INSERT IGNORE` then check).
4. Write `saga_log` entry.
5. **Commit DB Tx**.
6. Produce **ack/event** to Kafka (retry with backoff; dedupe using `idempotency_key` in key and an `event_id`).
7. **Commit Kafka offset** only after step 5 (DB commit).

   * If ack produce fails after DB commit, the consumer will retry produce (safe due to idempotent key).
   * If process crashes before offset commit, Kafka will re-deliver; idempotency table prevents double DB effects.

> This ordering guarantees **at-least-once** processing with **idempotent effects**, approximating exactly-once side-effects.

### 5.4 SAGA Choreography

* **Create**:

  * Command: `Create`
  * Tx: insert row
  * Ack/Event: `MessageCreated`
  * Compensating (if downstream fails later): `Delete` by id (recorded in `saga_log` for traceability)

* **Update**:

  * Command: `Update`
  * Tx: update row (with version check if you add `row_version`)
  * Event: `MessageUpdated`
  * Compensate: re-apply previous state (store snapshot or use `saga_log` prev values)

* **Delete**:

  * Command: `Delete`
  * Tx: delete row (or soft delete with `deleted_at`)
  * Event: `MessageDeleted`
  * Compensate: re-insert previous snapshot

* **Read** (optional async read-through):

  * Command: `Read`
  * Tx: SELECT row
  * Event/Ack: `MessageRead` with payload

> In single-service SAGA this looks trivial; the structure ensures the pattern generalizes when adding more services subscribe to `messages.events`.

---

## 6) Reliability & Idempotency

* **Kafka producer:** enable idempotent producer, set a stable `linger.ms` and `acks=all`.
* **Kafka consumer:** use consumer groups, max poll interval tuned for DB work.
* **Idempotency key:** required in every command; maintained in DB.
* **Deduplication:**

  * Before applying DB changes, check `idempotency_keys` for processed keys.
  * Ack producer uses the same key → duplicates collapse at the topic partition.
* **Retries & Backoff:** exponential with jitter for DB and Kafka produce.
* **Poison messages:** send to **DLQ** topic (e.g., `messages.commands.dlq`) after `N` failures with full context.

---

## 7) Observability

* **Tracing:** OpenTelemetry: span per HTTP request, per Kafka produce/consume, attach `trace_id`.
* **Metrics:** Prometheus counters/gauges/histograms:

  * `api_requests_total{verb,route,status}`
  * `kafka_messages_total{topic,type}`
  * `db_tx_duration_seconds`
  * `saga_step_total{step,status}`
* **Logging:** Structured JSON logs with `trace_id`, `correlation_id`, `command`, `event`, `status`.

---

## 8) Security & Validation

* **Input validation** at API boundaries (size limits, required fields).
* **Auth** (if required): JWT/OIDC; include `user_id` in metadata.
* **DB**: least-privileged user; parameterized queries.
* **Kafka**: SASL/SSL when available; per-topic ACLs.

---

## 9) Config & Environment

Environment variables (12-factor):

```
API_HTTP_ADDR=:8080
KAFKA_BROKERS=PLAINTEXT://kafka:9092
KAFKA_TOPIC_COMMANDS=messages.commands
KAFKA_TOPIC_ACKS=messages.acks
KAFKA_GROUP_ACKS=api-acks-consumer
MYSQL_DSN=user:pass@tcp(mysql:3306)/app?parseTime=true
RESULT_CACHE_TTL=120s
RESULT_POLL_TIMEOUT=15s
```

---

## 10) Docker, Minikube, Makefile

### 10.1 Dockerfiles (sketch)

**API Server** – `Dockerfile`:

```dockerfile
# build
FROM golang:1.22 AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/api ./cmd/api

# runtime
FROM gcr.io/distroless/base-debian12
COPY --from=build /out/api /api
EXPOSE 8080
ENTRYPOINT ["/api"]
```

**Consumer Service** – `Dockerfile`:

```dockerfile
FROM golang:1.22 AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/consumer ./cmd/consumer

FROM gcr.io/distroless/base-debian12
COPY --from=build /out/consumer /consumer
ENTRYPOINT ["/consumer"]
```

### 10.2 Kubernetes (Minikube) Manifests (sketch)

`k8s/api.yaml`:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata: { name: api }
spec:
  replicas: 1
  selector: { matchLabels: { app: api } }
  template:
    metadata: { labels: { app: api } }
    spec:
      containers:
        - name: api
          image: api:local
          imagePullPolicy: IfNotPresent
          env:
            - { name: API_HTTP_ADDR, value: ":8080" }
            - { name: KAFKA_BROKERS, value: "kafka:9092" }
            - { name: KAFKA_TOPIC_COMMANDS, value: "messages.commands" }
            - { name: KAFKA_TOPIC_ACKS, value: "messages.acks" }
            - { name: MYSQL_DSN, valueFrom: { secretKeyRef: { name: app-db, key: dsn } } }
          ports:
            - containerPort: 8080
---
apiVersion: v1
kind: Service
metadata: { name: api }
spec:
  type: NodePort
  selector: { app: api }
  ports:
    - port: 80
      targetPort: 8080
```

`k8s/consumer.yaml`:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata: { name: consumer }
spec:
  replicas: 1
  selector: { matchLabels: { app: consumer } }
  template:
    metadata: { labels: { app: consumer } }
    spec:
      containers:
        - name: consumer
          image: consumer:local
          imagePullPolicy: IfNotPresent
          env:
            - { name: KAFKA_BROKERS, value: "kafka:9092" }
            - { name: KAFKA_TOPIC_COMMANDS, value: "messages.commands" }
            - { name: KAFKA_TOPIC_ACKS, value: "messages.acks" }
            - { name: MYSQL_DSN, valueFrom: { secretKeyRef: { name: app-db, key: dsn } } }
```

`k8s/mysql.yaml` (dev):

```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata: { name: mysql-pvc }
spec:
  accessModes: ["ReadWriteOnce"]
  resources: { requests: { storage: 2Gi } }
---
apiVersion: apps/v1
kind: Deployment
metadata: { name: mysql }
spec:
  selector: { matchLabels: { app: mysql } }
  template:
    metadata: { labels: { app: mysql } }
    spec:
      containers:
        - name: mysql
          image: mysql:8.0
          env:
            - { name: MYSQL_ROOT_PASSWORD, valueFrom: { secretKeyRef:{name: app-db, key: rootpw} } }
            - { name: MYSQL_DATABASE, value: "app" }
          ports: [ { containerPort: 3306 } ]
          volumeMounts:
            - { name: data, mountPath: /var/lib/mysql }
      volumes:
        - name: data
          persistentVolumeClaim: { claimName: mysql-pvc }
---
apiVersion: v1
kind: Service
metadata: { name: mysql }
spec:
  selector: { app: mysql }
  ports: [ { port: 3306 } ]
```

`k8s/kafka.yaml` (dev, single-broker; use Strimzi/Bitnami for quickstart) – omitted for brevity.

**Secrets**:

```yaml
apiVersion: v1
kind: Secret
metadata: { name: app-db }
type: Opaque
stringData:
  dsn: "user:pass@tcp(mysql:3306)/app?parseTime=true"
  rootpw: "rootpassword"
```

### 10.3 Makefile (local DX)

```makefile
APP?=api
CONSUMER?=consumer

.PHONY: build docker kind-up k8s-apply minikube-load test lint dev-up dev-down

build:
	go build -o bin/api ./cmd/api
	go build -o bin/consumer ./cmd/consumer

docker:
	docker build -t api:local -f ./deploy/api.Dockerfile .
	docker build -t consumer:local -f ./deploy/consumer.Dockerfile .

minikube-load: docker
	# For Docker driver in Minikube this may be unnecessary
	minikube image load api:local
	minikube image load consumer:local

k8s-apply:
	kubectl apply -f k8s/mysql.yaml
	kubectl apply -f k8s/kafka.yaml
	kubectl apply -f k8s/api.yaml
	kubectl apply -f k8s/consumer.yaml

dev-up: build docker minikube-load k8s-apply

dev-down:
	kubectl delete -f k8s/consumer.yaml || true
	kubectl delete -f k8s/api.yaml || true
	kubectl delete -f k8s/mysql.yaml || true
	kubectl delete -f k8s/kafka.yaml || true

test:
	go test ./...

lint:
	golangci-lint run
```

---

## 11) Error Handling & DLQ Strategy

* **Validation errors (API)** → return `400`; do not enqueue.
* **Enqueue failure** → `503` with retry-after header.
* **Consumer hard failure (e.g., schema violation)** → publish to `messages.commands.dlq` with full original payload + error.
* **Transient DB/Kafka failures** → automatic retries with capped backoff.

---

## 12) Scaling & Partitioning

* **Kafka partition key**:

  * For CRUD on `Message`: use `payload.id` (ensures per-record order).
  * For Create without ID: use `idempotency_key` (random, good distribution).
* **API horizontal scale**:

  * Multiple replicas OK; Results Cache must be externalized (Redis) or rely on a **compacted results topic** and re-consume on demand.
* **Consumer scale**:

  * Scale to number of partitions. Each consumer in the group gets partitions exclusively.

---

## 13) Testing Strategy

* **Unit tests**:

  * Command serialization, validation, trace propagation.
  * DB repository with `sqlmock`.
* **Integration tests**:

  * Spin Kafka (testcontainers) + MySQL (testcontainers).
  * Golden-path CRUD SAGA + retries + idempotency.
* **Contract tests**:

  * JSON schemas for command/ack envelopes (use `gojsonschema`).
* **Load tests**:

  * k6 for API; confirm latency for `202` and mean time to ack.

---

## 14) Implementation Sketches (Go)

**API: produce command**

```go
traceID := uuid.New().String()
ctx := context.WithValue(r.Context(), ctxKeyTraceID, traceID)

idemp := uuid.New().String()
cmd := Command{ /* fill from request */ }

key := []byte(idemp)
headers := []kafka.Header{
  {Key:"trace_id", Value:[]byte(traceID)},
  {Key:"command", Value:[]byte(cmd.Command)},
  {Key:"resource", Value:[]byte("Message")},
}
payload, _ := json.Marshal(cmd)
producer.Produce(&kafka.Message{
  TopicPartition: kafka.TopicPartition{Topic: &topicCommands, Partition: kafka.PartitionAny},
  Key:   key,
  Value: payload,
  Headers: headers,
}, nil)
```

**API: results endpoint (long-poll)**

```go
deadlineCtx, cancel := context.WithTimeout(r.Context(), resultPollTimeout)
defer cancel()
for {
  if res, ok := resultCache.Get(traceID); ok {
    writeJSON(w, http.StatusOK, res)
    return
  }
  select {
  case <-deadlineCtx.Done():
    w.WriteHeader(http.StatusNoContent)
    return
  case <-time.After(200 * time.Millisecond):
  }
}
```

**Consumer: idempotent DB tx**

```go
err := withTx(db, func(tx *sql.Tx) error {
  processed, err := repo.CheckIdempotency(tx, idempKey)
  if err != nil { return err }
  if processed { already = true; return nil }

  switch cmd.Command {
  case "Create":
    id, err := repo.InsertMessage(tx, cmd.Payload.Message)
    if err != nil { return err }
    saga.Log(tx, traceID, "CreateMessage", "SUCCESS", "", "")
    return repo.MarkIdempotent(tx, idempKey, traceID, "SUCCESS")
  // Update, Delete, Read ...
  }
})
if err != nil { /* retry or DLQ */ }

if already {
  // produce same ack again (safe)
}
produceAck(traceID, idempKey, status, payload)
commitOffset()
```

---

## 15) Risks & Trade-offs

* **Result cache in memory**: simple for dev; not HA. For production, use Redis or compacted Kafka state.
* **EOS with Go clients**: Kafka’s true transactions are best in Java; we emulate exactly-once effects with **idempotency keys** + Tx order + retry rules.
* **SAGA complexity**: trivial with single DB; pattern choice future-proofs when adding more services.

---

## 16) Operational Runbook

* **Create topics**:

  * `messages.commands` (partitions: 6, replication: 1 in dev)
  * `messages.acks` (partitions: 6)
  * `messages.events` (optional)
  * `messages.commands.dlq` (partitions: 6)
* **Migrate DB**: use `golang-migrate` on startup.
* **Deploy order**: MySQL → Kafka → Consumer → API.
* **Health checks**:

  * API: `/healthz` checks Kafka producer connectivity + ack consumer lag.
  * Consumer: reports consumer lag + DB ping.

---

## 17) Example Sequence: Update with Failure & Compensation

```mermaid
sequenceDiagram
  participant API
  participant K as Kafka
  participant CS as Consumer
  participant DB as MySQL
  participant KA as Kafka(acks)

  API->>K: Command(Update id=42)
  CS->>DB: Tx: SELECT ...; UPDATE message SET ...
  Note over CS,DB: Suppose downstream (not in scope) fails later.
  CS->>KA: Ack SUCCESS (MessageUpdated)
  note over CS: Another service detects failure and issues a compensating command.
  CS->>DB: Tx: Revert to previous content (compensate)
  CS->>KA: Event CompensationApplied
```

---

## 18) Deliverables

* **Go services**:

  * `cmd/api` (REST, Swagger v2, Kafka producer, ack consumer + cache)
  * `cmd/consumer` (Kafka consumer, MySQL repository, Tx logic, ack producer)
* **Common**:

  * `pkg/contracts` (envelopes, JSON schemas)
  * `pkg/trace` (context utils)
  * `pkg/kafka` (producer/consumer wrappers)
  * `pkg/repo` (DB repos)
* **Ops**:

  * Dockerfiles, `k8s/*.yaml`, `Makefile`
  * DB migrations
  * Sample `env/.env.example`

---

## 19) Non-Goals (Current Scope)

* No multi-tenant auth flows or RBAC.
* No cross-region disaster recovery.
* No full-blown CQRS read models (async reads are minimal).

---

## 20) Future Enhancements

* **Compacted results topic** with key=`trace_id` to replace in-memory cache.
* **Redis** for HA result delivery.
* **Schema Registry** & versioned envelopes.
* **Helm chart** for parameterized deployments.
* **Canary releases** for consumer with shadow topics.

---

### Appendix A — Topic Configuration (Dev)

* `cleanup.policy=delete`
* `min.insync.replicas=1`
* `retention.ms=86400000` (1 day) for commands; longer for events/acks as needed
* `acks=all` on producers; enable idempotence

### Appendix B — Minimal SQL Migrations (v0001)

```sql
-- create tables as listed in §5.2
```

---

**End of Design** ✅
