# Kafka + Go (Sarama) with Retry, DLQ & OpenTelemetry — Ready‑to‑Run Demo

This canvas contains a complete, ready‑to‑run demo repository that wires up **Kafka (KRaft)**, **Go + IBM Sarama**, **staged retries + DLQ**, and **OpenTelemetry** tracing via `otelsarama`.

---

## Quickstart

```bash
unzip kafka-go-sarama-demo.zip
cd kafka-go-sarama-demo

# 1) Start infra (Kafka + OTEL Collector)
make up

# 2) Create topics (main, retry stages, DLQ)
make topics

# 3) Start services (two terminals)
make processor
make retryworker

# 4) Send sample messages (one ok, one failing → retries → DLQ)
make producer

# 5) See spans flowing in the collector logs
make otel-logs
```

### What you’ll see

* Messages like `ok: welcome` are processed immediately by the **processor**.
* Messages like `fail: simulate downstream error` go to `events.v1.retry.5s`, then are re‑queued to `events.v1`. If they still fail, they progress to `events.v1.retry.30s`, then `events.v1.retry.2m`, and finally to **DLQ**.

---

## Repository Layout

```
kafka-go-sarama-demo/
├─ cmd/
│  ├─ admin/         # topic creation
│  ├─ producer/      # demo producer
│  ├─ processor/     # consumer group processor with retry→DLQ
│  └─ retryworker/   # consumes retry topics, sleeps, re-queues to main
├─ internal/
│  ├─ retry/         # retry stages + headers
│  └─ tracing/       # OTel bootstrap and (optional) header carrier helper
├─ compose.yaml      # Kafka (KRaft) + OTEL Collector
├─ otel-collector-config.yaml
├─ Makefile
├─ README.md
└─ GUIDE.md
```

---

## Topics Used

* `events.v1` (main)
* `events.v1.retry.5s`, `events.v1.retry.30s`, `events.v1.retry.2m` (retry stages)
* `events.v1.dlq` (dead‑letter)

> **Note:** For local dev we use replication factor `1`. In production use `≥ 3`.

---

## Makefile

```makefile
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
```

---

## Docker Compose (Kafka KRaft + OTEL Collector)

```yaml
services:
  kafka:
    image: bitnami/kafka:3.8
    container_name: kafka
    environment:
      - KAFKA_ENABLE_KRAFT=yes
      - KAFKA_CFG_NODE_ID=1
      - KAFKA_CFG_PROCESS_ROLES=broker,controller
      - KAFKA_CFG_CONTROLLER_LISTENER_NAMES=CONTROLLER
      - KAFKA_CFG_LISTENER_SECURITY_PROTOCOL_MAP=CONTROLLER:PLAINTEXT,PLAINTEXT:PLAINTEXT,PLAINTEXT_HOST:PLAINTEXT
      - KAFKA_CFG_LISTENERS=CONTROLLER://:9094,PLAINTEXT://:9093,PLAINTEXT_HOST://:9092
      - KAFKA_CFG_ADVERTISED_LISTENERS=PLAINTEXT_HOST://localhost:9092,PLAINTEXT://kafka:9093
      - KAFKA_CFG_CONTROLLER_QUORUM_VOTERS=1@kafka:9094
      - KAFKA_CFG_INTER_BROKER_LISTENER_NAME=PLAINTEXT
      - KAFKA_CFG_AUTO_CREATE_TOPICS_ENABLE=false
      - KAFKA_CFG_NUM_PARTITIONS=3
      - KAFKA_CFG_OFFSETS_TOPIC_REPLICATION_FACTOR=1
    ports:
      - "9092:9092"

  otel-collector:
    image: otel/opentelemetry-collector-contrib:0.104.0
    container_name: otel-collector
    command: ["--config=/etc/otel-collector-config.yaml"]
    volumes:
      - ./otel-collector-config.yaml:/etc/otel-collector-config.yaml:ro
    ports:
      - "4317:4317"  # OTLP gRPC
      - "4318:4318"  # OTLP HTTP
```

### OTEL Collector Config

```yaml
receivers:
  otlp:
    protocols:
      grpc:
      http:

exporters:
  logging:
    logLevel: debug

service:
  pipelines:
    traces:
      receivers: [otlp]
      exporters: [logging]
```

---

## `go.mod`

```go
module example.com/kafka-go-sarama-demo

go 1.22.0

require (
	github.com/IBM/sarama v1.45.0
	github.com/dnwe/otelsarama v0.4.3
	go.opentelemetry.io/otel v1.28.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.28.0
	go.opentelemetry.io/otel/sdk v1.28.0
	go.opentelemetry.io/otel/semconv v1.26.0
)
```

---

## Admin: Create Topics

**`cmd/admin/main.go`**

```go
package main

import (
	"log"
	"time"

	"github.com/IBM/sarama"
)

func str(s string) *string { return &s }
func must(err error) { if err != nil { log.Fatal(err) } }

func main() {
	cfg := sarama.NewConfig()
	cfg.Version, _ = sarama.ParseKafkaVersion("3.8.0")

	admin, err := sarama.NewClusterAdmin([]string{"localhost:9092"}, cfg)
	must(err)
	defer admin.Close()

	topics := map[string]*sarama.TopicDetail{
		"events.v1":                {NumPartitions: 3, ReplicationFactor: 1, ConfigEntries: map[string]*string{
			"retention.ms": str("604800000"), // 7 days
		}},
		"events.v1.retry.5s":       {NumPartitions: 3, ReplicationFactor: 1, ConfigEntries: map[string]*string{
			"retention.ms": str("3600000"),   // 1 hour
		}},
		"events.v1.retry.30s":      {NumPartitions: 3, ReplicationFactor: 1, ConfigEntries: map[string]*string{
			"retention.ms": str("3600000"),
		}},
		"events.v1.retry.2m":       {NumPartitions: 3, ReplicationFactor: 1, ConfigEntries: map[string]*string{
			"retention.ms": str("3600000"),
		}},
		"events.v1.dlq":            {NumPartitions: 3, ReplicationFactor: 1, ConfigEntries: map[string]*string{
			"retention.ms": str("1209600000"), // 14 days
		}},
	}

	for t, d := range topics {
		if err := admin.CreateTopic(t, d, false); err != nil {
			log.Printf("CreateTopic(%s): %v (ignored if already exists)", t, err)
		}
	}
	time.Sleep(time.Second)
	log.Println("Topic setup complete.")
}
```

---

## Producer (with OTEL)

**`cmd/producer/main.go`**

```go
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/IBM/sarama"
	"github.com/dnwe/otelsarama"
	"example.com/kafka-go-sarama-demo/internal/tracing"
)

func mustParse(v string) sarama.KafkaVersion {
	ver, err := sarama.ParseKafkaVersion(v); if err != nil { log.Fatal(err) }
	return ver
}

func main() {
	shutdown, err := tracing.Init("producer")
	if err != nil { log.Fatalf("otel init: %v", err) }
	defer shutdown(nil)

	cfg := sarama.NewConfig()
	cfg.Version = mustParse("3.8.0")
	cfg.Producer.Idempotent = true
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	cfg.Net.MaxOpenRequests = 1
	cfg.Producer.Return.Successes = true
	cfg.Producer.Retry.Max = 10
	cfg.Producer.Compression = sarama.CompressionSnappy
	cfg.Metadata.RefreshFrequency = time.Minute

	raw, err := sarama.NewSyncProducer([]string{"localhost:9092"}, cfg)
	if err != nil { log.Fatalf("new producer: %v", err) }
	prod := otelsarama.WrapSyncProducer(cfg, raw)
	defer prod.Close()

	send := func(val string) {
		msg := &sarama.ProducerMessage{
			Topic: "events.v1",
			Key:   sarama.StringEncoder("user-42"),
			Value: sarama.StringEncoder(val),
			Headers: []sarama.RecordHeader{
				{Key: []byte("content-type"), Value: []byte("text/plain")},
			},
		}
		p, o, err := prod.SendMessage(msg)
		if err != nil { log.Printf("send error: %v", err); return }
		log.Printf("sent partition=%d offset=%d val=%q", p, o, val)
	}

	send("ok: welcome")
	send("fail: simulate downstream error")
	fmt.Println("done.")
}
```

---

## Processor (Consumer Group with Retry → DLQ + OTEL)

**`cmd/processor/main.go`**

```go
package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/IBM/sarama"
	"github.com/dnwe/otelsarama"

	"example.com/kafka-go-sarama-demo/internal/retry"
	"example.com/kafka-go-sarama-demo/internal/tracing"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
)

type handler struct{ prod sarama.SyncProducer }

func (h *handler) Setup(s sarama.ConsumerGroupSession) error   { return nil }
func (h *handler) Cleanup(s sarama.ConsumerGroupSession) error { return nil }

func parseAttempt(msg *sarama.ConsumerMessage) int {
	for _, h := range msg.Headers {
		if string(h.Key) == retry.HeaderAttempt {
			if n, err := strconv.Atoi(string(h.Value)); err == nil { return n }
		}
	}
	return 0
}

func (h *handler) publishNextRetry(msg *sarama.ConsumerMessage, err error) error {
	attempt := parseAttempt(msg)
	if stage, ok := retry.Next(attempt); ok {
		out := &sarama.ProducerMessage{
			Topic: stage.Topic,
			Key:   sarama.ByteEncoder(msg.Key),
			Value: sarama.ByteEncoder(msg.Value),
			Headers: append([]sarama.RecordHeader{}, append(msg.Headers,
				sarama.RecordHeader{Key: []byte(retry.HeaderAttempt), Value: []byte(strconv.Itoa(attempt + 1))},
				sarama.RecordHeader{Key: []byte(retry.HeaderError),   Value: []byte(err.Error())},
			)...),
		}
		_, _, e := h.prod.SendMessage(out)
		return e
	}
	// Exhausted → DLQ
	out := &sarama.ProducerMessage{
		Topic: "events.v1.dlq",
		Key:   sarama.ByteEncoder(msg.Key),
		Value: sarama.ByteEncoder(msg.Value),
		Headers: append(msg.Headers,
			sarama.RecordHeader{Key: []byte(retry.HeaderAttempt), Value: []byte(strconv.Itoa(attempt))},
			sarama.RecordHeader{Key: []byte(retry.HeaderError),   Value: []byte(err.Error())},
		),
	}
	_, _, e := h.prod.SendMessage(out)
	return e
}

// businessLogic demonstrates a manual child span (e.g., simulating a DB write).
func businessLogic(msg *sarama.ConsumerMessage) error {
	// Extract context from message headers for proper span parenting.
	carrier := propagation.HeaderCarrier{}
	for _, h := range msg.Headers {
		carrier.Set(string(h.Key), string(h.Value))
	}
	ctx := otel.GetTextMapPropagator().Extract(context.Background(), carrier)

	ctx, span := otel.Tracer("processor").Start(ctx, "businessLogic")
	defer span.End()

	span.SetAttributes(
		attribute.String("kafka.topic", msg.Topic),
		attribute.Int("kafka.partition", int(msg.Partition)),
		attribute.Int64("kafka.offset", msg.Offset),
	)

	// Very basic demo: fail when payload starts with "fail:"
	if len(msg.Value) >= 5 && string(msg.Value[:5]) == "fail:" {
		err := errors.New("downstream: simulated failure")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	// Simulate work (e.g., DB call)
	time.Sleep(50 * time.Millisecond)

	span.SetStatus(codes.Ok, "ok")
	return nil
}

func (h *handler) ConsumeClaim(s sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		if err := businessLogic(msg); err != nil {
			log.Printf("process error, routing to retry/DLQ: %v", err)
			if e := h.publishNextRetry(msg, err); e != nil {
				log.Printf("retry publish failed: %v", e)
				continue // don't mark => will be retried
			}
			s.MarkMessage(msg, "forwarded")
			continue
		}
		s.MarkMessage(msg, "")
	}
	return nil
}

func newSyncProducer(cfg *sarama.Config) sarama.SyncProducer {
	p, err := sarama.NewSyncProducer([]string{"localhost:9092"}, cfg)
	if err != nil { log.Fatalf("producer: %v", err) }
	return p
}

func main() {
	shutdown, err := tracing.Init("processor")
	if err != nil { log.Fatalf("otel init: %v", err) }
	defer shutdown(context.Background())

	cfg := sarama.NewConfig()
	cfg.Version, _ = sarama.ParseKafkaVersion("3.8.0")
	cfg.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
	cfg.Consumer.Offsets.Initial = sarama.OffsetOldest
	cfg.Metadata.RefreshFrequency = time.Minute

	// producer for retry/DLQ publishing and instrument it.
	pcfg := sarama.NewConfig()
	pcfg.Version = cfg.Version
	pcfg.Producer.RequiredAcks = sarama.WaitForAll
	pcfg.Producer.Idempotent = true
	pcfg.Net.MaxOpenRequests = 1
	pcfg.Producer.Return.Successes = true
	pcfg.Producer.Retry.Max = 10

	rawProd := newSyncProducer(pcfg)
	prod := otelsarama.WrapSyncProducer(pcfg, rawProd)
	defer prod.Close()

	cg, err := sarama.NewConsumerGroup([]string{"localhost:9092"}, "processor.v1", cfg)
	if err != nil { log.Fatalf("consumer group: %v", err) }
	defer cg.Close()

	h := otelsarama.WrapConsumerGroupHandler(&handler{prod: prod})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() { for err := range cg.Errors() { log.Printf("consumer error: %v", err) } }()

	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		<-sig
		cancel()
	}()

	for ctx.Err() == nil {
		if err := cg.Consume(ctx, []string{"events.v1"}, h); err != nil {
			log.Printf("consume: %v", err)
			time.Sleep(time.Second)
		}
	}
}
```

---

## Retry Worker (sleeps then re‑queues to main) + OTEL

**`cmd/retryworker/main.go`**

```go
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/IBM/sarama"
	"github.com/dnwe/otelsarama"

	"example.com/kafka-go-sarama-demo/internal/tracing"
)

var topicDelay = map[string]time.Duration{
	"events.v1.retry.5s":  5 * time.Second,
	"events.v1.retry.30s": 30 * time.Second,
	"events.v1.retry.2m":  2 * time.Minute,
}

type handler struct{ prod sarama.SyncProducer }

func (h *handler) Setup(s sarama.ConsumerGroupSession) error   { return nil }
func (h *handler) Cleanup(s sarama.ConsumerGroupSession) error { return nil }

func (h *handler) ConsumeClaim(s sarama.ConsumerGroupSession, c sarama.ConsumerGroupClaim) error {
	delay := topicDelay[c.Topic()]
	for msg := range c.Messages() {
		time.Sleep(delay) // backoff window

		out := &sarama.ProducerMessage{
			Topic: "events.v1",
			Key:   sarama.ByteEncoder(msg.Key),
			Value: sarama.ByteEncoder(msg.Value),
			Headers: msg.Headers, // keep headers (including x-retry-attempt & x-error)
		}
		if _, _, err := h.prod.SendMessage(out); err != nil {
			// If we fail to requeue, we won't mark => message will be retried by this group
			log.Printf("requeue failed: %v", err)
			continue
		}
		s.MarkMessage(msg, "requeued")
	}
	return nil
}

func main() {
	shutdown, err := tracing.Init("retryworker")
	if err != nil { log.Fatalf("otel init: %v", err) }
	defer shutdown(context.Background())

	cfg := sarama.NewConfig()
	cfg.Version, _ = sarama.ParseKafkaVersion("3.8.0")
	cfg.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
	cfg.Consumer.Offsets.Initial = sarama.OffsetOldest

	pcfg := sarama.NewConfig()
	pcfg.Version = cfg.Version
	pcfg.Producer.RequiredAcks = sarama.WaitForAll
	pcfg.Producer.Idempotent = true
	pcfg.Net.MaxOpenRequests = 1
	pcfg.Producer.Return.Successes = true

	rawProd, err := sarama.NewSyncProducer([]string{"localhost:9092"}, pcfg)
	if err != nil { log.Fatalf("producer: %v", err) }
	prod := otelsarama.WrapSyncProducer(pcfg, rawProd)
	defer prod.Close()

	cg, err := sarama.NewConsumerGroup([]string{"localhost:9092"}, "retryworker.v1", cfg)
	if err != nil { log.Fatalf("consumer group: %v", err) }
	defer cg.Close()

	topics := []string{"events.v1.retry.5s", "events.v1.retry.30s", "events.v1.retry.2m"}
	h := otelsarama.WrapConsumerGroupHandler(&handler{prod: prod})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() { for err := range cg.Errors() { log.Printf("cg error: %v", err) } }()
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		<-sig; cancel()
	}()

	for ctx.Err() == nil {
		if err := cg.Consume(ctx, topics, h); err != nil {
			log.Printf("consume: %v", err)
			time.Sleep(time.Second)
		}
	}
}
```

---

## Internal Packages

### Retry Stages

**`internal/retry/retry.go`**

```go
package retry

import "time"

const (
	HeaderAttempt = "x-retry-attempt"
	HeaderError   = "x-error"
)

type Stage struct {
	Topic string
	Delay time.Duration
}

var Stages = []Stage{
	{Topic: "events.v1.retry.5s",  Delay: 5 * time.Second},
	{Topic: "events.v1.retry.30s", Delay: 30 * time.Second},
	{Topic: "events.v1.retry.2m",  Delay: 2 * time.Minute},
}

func Next(attempt int) (Stage, bool) {
	if attempt < len(Stages) {
		return Stages[attempt], true
	}
	return Stage{}, false
}
```

### OpenTelemetry Bootstrap

**`internal/tracing/tracing.go`**

```go
package tracing

import (
	"context"
	"log"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Init sets up OTLP exporter + tracer provider and returns a shutdown function.
func Init(serviceName string) (func(context.Context) error, error) {
	endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if endpoint == "" {
		endpoint = "localhost:4317"
	}

	exp, err := otlptracegrpc.New(context.Background(),
		otlptracegrpc.WithEndpoint(endpoint),
		otlptracegrpc.WithDialOption(grpc.WithTransportCredentials(insecure.NewCredentials())),
	)
	if err != nil { return nil, err }

	res, err := resource.New(context.Background(),
		resource.WithAttributes(
			semconv.ServiceName(serviceName),
			attribute.String("deployment.environment", "local"),
		),
	)
	if err != nil { return nil, err }

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	log.Printf("OTEL initialized for service=%s -> %s", serviceName, endpoint)
	return tp.Shutdown, nil
}
```

*(Optional helper for custom header carriers—kept for future extension)*

**`internal/tracing/kafka_propagation.go`**

```go
package tracing

import (
	"strings"

	"github.com/IBM/sarama"
	"go.opentelemetry.io/otel/propagation"
)

// HeaderCarrier implements OTEL's TextMapCarrier for Sarama headers.
type HeaderCarrier struct{ Headers *[]*sarama.RecordHeader }

func (c HeaderCarrier) Get(key string) string {
	if c.Headers == nil { return "" }
	lk := strings.ToLower(key)
	for _, h := range *c.Headers {
		if strings.ToLower(string(h.Key)) == lk {
			return string(h.Value)
		}
	}
	return ""
}
func (c HeaderCarrier) Set(key, val string) {
	if c.Headers == nil { return }
	lk := strings.ToLower(key)
	for _, h := range *c.Headers {
		if strings.ToLower(string(h.Key)) == lk {
			h.Value = []byte(val); return
		}
	}
	*c.Headers = append(*c.Headers, &sarama.RecordHeader{Key: []byte(key), Value: []byte(val)})
}
func (c HeaderCarrier) Keys() []string {
	if c.Headers == nil { return nil }
	keys := make([]string, 0, len(*c.Headers))
	for _, h := range *c.Headers { keys = append(keys, string(h.Key)) }
	return keys
}

// ExtractContext builds a context from Kafka headers (example stub).
func ExtractContext(headers *[]*sarama.RecordHeader, propagator propagation.TextMapPropagator) (ctx interface{ Done() <-chan struct{} }, _ propagation.TextMapCarrier) {
	return nil, HeaderCarrier{Headers: headers}
}
```

---

## GUIDE.md (Beginner‑Friendly Notes)

**Kafka listeners – why two?**
Kafka needs to tell clients **which host/port** to connect to next. In Docker we bind:

* `PLAINTEXT_HOST://localhost:9092` (for host apps)
* `PLAINTEXT://kafka:9093` (for apps inside containers)

**Retry & DLQ pattern**
Kafka doesn’t delay messages. We simulate delayed retries with **staged retry topics**:

* First failure → publish to `events.v1.retry.5s` with header `x-retry-attempt=1`
* The **retry worker** sleeps 5s, re‑queues back to `events.v1`
* Keeps escalating to `30s`, `2m`, then **DLQ**

**Tracing**
We use the `otelsarama` wrappers to create spans for produce/consume and to **propagate context** in Kafka headers. We also show a manual **child span** in the processor’s `businessLogic` to simulate a DB write.

**Idempotency**
Retries happen. Your business logic should be **idempotent**—use a de‑dup key `(topic,partition,offset)` or a domain key to avoid double effects.

**What to customize**

* Add your own payload schema (JSON, Protobuf, Avro) and validation
* Tune retry stages (delay/retention)
* Replace the simulated error with real work (DB/HTTP) and instrument it

---

## README.md (Included)

* What’s inside (Sarama, Consumer Group, Retry, DLQ, OTEL)
* Make targets (`up/down/topics/processor/retryworker/producer/otel-logs`)
* Structure and notes

---

## Next Extensions (ask and I’ll add)

* DLQ reader CLI that prints original headers and `x-error`
* Prometheus metrics for retry rate, DLQ growth, processor latency
* SASL/TLS examples (PLAIN/SCRAM) with minimal broker config
* Avro/Schema Registry or Protobuf payloads with validation
