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