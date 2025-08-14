# SLB Microservices & Go Interview Cheat Sheet

*(Condensed, 2-page printable)*

---

## 1. Microservices Good Practices (SAGA)

**Key Practices:**

* **Bounded context**: one service = one domain + private DB.
* **Contracts first**: versioned APIs; backward-compatible changes.
* **Observability**: logs, metrics, traces with correlation IDs.
* **Resilience**: timeouts, retries, circuit breakers, bulkheads.
* **Idempotency**: unique request IDs, dedupe logic.
* **Security**: mTLS, least privilege, vault secrets.

**SAGA Patterns:**

* **Choreography** (event-driven): decentralized, services emit/react to events.
* **Orchestration**: central controller coordinates steps/compensations.

**Compensation Example:**
CreateOrder → ReserveInventory → ChargePayment
CancelOrder ← ReleaseInventory ← RefundPayment

**SAGA Orchestration Diagram:**

```
[Orchestrator] --> [Service A] --> [Service B] --> [Service C]
      |                |              |              |
      |<--Compensate---|<--Compensate-|<--Compensate-|
```

---

## 2. Clean Code Design (Go)

* **Small packages**, clear responsibility.
* **Interfaces at consumer side**, small (1–3 methods).
* **Error wrapping**: `fmt.Errorf("context: %w", err)`.
* **Pure business logic** in functions; I/O at edges.
* **Testing**: table tests, golden files, fuzz.
* **Formatting**: `go fmt`, `golangci-lint`.

---

## 3. Debugging Distributed Microservices

* **Reproduce** issue scope → correlate via `X-Request-ID`.
* **Logs**: structured JSON, include trace info.
* **Metrics**: RED (Rate, Errors, Duration).
* **Tracing**: OpenTelemetry for cross-service latency.
* **Replay**: traffic capture & staging replays.
* **Chaos testing** to simulate failures.

---

## 4. Middlewares in Go

```go
func logging(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    start := time.Now()
    next.ServeHTTP(w, r)
    log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
  })
}
```

Used for logging, auth, tracing, rate limiting.

---

## 5. Kafka Replay & Error Handling

**Replay:**

* `Seek(offset)` or by timestamp.
* Reset group offsets.
* Consume from DLQ.

**Error Handling Pattern:**

```
Main Topic --> Retry 5s --> Retry 1m --> Retry 10m --> DLQ
```

* Idempotent consumers.
* DLQ for poison messages.

---

## 6. DB Errors Midway

* **Single DB**: ACID transactions; keep short-lived.
* **Multi-service**: SAGA + outbox pattern.
* **Idempotency**: Upserts, natural keys.

---

## 7. Message Broker QoS

* **At most once**: no retries, possible loss.
* **At least once**: retries until ack, possible duplicates.
* **Exactly once**: dedupe + transactions, costly.

---

## 8. SQL vs NoSQL (CAP & ACID)

* **SQL**: strong consistency, joins, transactions.
* **NoSQL**: flexible schema, massive scale, eventual consistency.
* **CAP**:

  * CP: Consistency > Availability (etcd, ZK).
  * AP: Availability > Consistency (Dynamo, Cassandra).

---

## 9. Indexing Pros/Cons

* **+** Faster reads, constraint enforcement.
* **–** Slower writes, more storage, maintenance overhead.

---

## 10–11. Goroutine Stack Facts

* Starts \~2 KB; grows/shrinks automatically.
* Max \~1 GB → stack overflow possible with deep recursion or huge local vars.

---

## 12. Go Concurrency Patterns

**Fan-out/Fan-in**

```
[Jobs] --> Worker1 --> \
        --> Worker2 -->  --> [Results]
        --> Worker3 --> /
```

**Pipeline**: stages connected via channels.
**Worker Pools**: bounded concurrency.
**Context Cancellation**: cooperative shutdown.
**errgroup**: concurrent tasks with aggregated errors.

---

## 13. SOLID in Go

* **S**: One reason to change per package/type.
* **O**: Open to new implementations via interfaces.
* **L**: Interfaces behave as expected for all impls.
* **I**: Small, specific interfaces.
* **D**: Depend on abstractions you own; inject concretes at edges.

---

**Retry/DLQ Diagram:**

```
[Main Topic] -> [Retry 5s] -> [Retry 1m] -> [Retry 10m] -> [DLQ]
```
