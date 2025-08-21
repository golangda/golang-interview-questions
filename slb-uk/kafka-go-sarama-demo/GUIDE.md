# Guide: How this works (beginner friendly)

## 1) Kafka listeners – why two?
Kafka needs to tell clients **which host/port** to connect to next. In Docker we bind:
- `PLAINTEXT_HOST://localhost:9092` (host apps)
- `PLAINTEXT://kafka:9093` (other containers)

This avoids the classic “connects then fails” pitfall.

## 2) Retry & DLQ pattern
Kafka doesn’t delay messages. We simulate delayed retries with **staged retry topics**:
- First failure → publish to `events.v1.retry.5s` with header `x-retry-attempt=1`
- The **retry worker** sleeps 5s, re-queues back to `events.v1`
- Keeps escalating to `30s`, `2m`, then **DLQ**

## 3) Tracing
We use the `otelsarama` wrappers to create spans for produce/consume and to **propagate context** in Kafka headers.
We also show a manual **child span** in the processor’s `businessLogic` to simulate a DB write.

## 4) Idempotency
Retries happen. Your business logic should be **idempotent**:
- Use a de-dup key `(topic,partition,offset)` or your message key to avoid double processing.

## 5) What to customize
- Add your own payload schema (JSON, Protobuf, Avro) and validation
- Tune retry stages (delay/retention)
- Replace the simulated business error with real work (DB/HTTP) and instrument it
