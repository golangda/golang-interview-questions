package common

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type Event struct {
	SagaID        string                 `json:"saga_id"`
	Step          int                    `json:"step"`
	SchemaVersion int                    `json:"schema_version"`
	Ts            time.Time              `json:"ts"`
	Payload       map[string]any         `json:"payload"`
}

var (
	StepLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{Name: "saga_step_latency_seconds", Help: "latency per step", Buckets: []float64{.01, .05, .1, .25, .5, 1, 2, 5}},
		[]string{"step"},
	)
	RetriesTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{Name: "saga_retries_total", Help: "retries by step/reason"},
		[]string{"step","reason"},
	)
	DLQTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{Name: "dlq_messages_total", Help: "messages sent to dlq by topic"},
		[]string{"topic"},
	)
)

func init() {
	prometheus.MustRegister(StepLatency, RetriesTotal, DLQTotal)
}

func ServeMetrics() {
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		log.Println("[metrics] listening on :8080/metrics")
		_ = http.ListenAndServe(":8080", nil)
	}()
}

func InitOTel() func(context.Context) error {
	collector := os.Getenv("JAEGER_COLLECTOR")
	if collector == "" {
		return func(context.Context) error { return nil }
	}
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(collector)))
	if err != nil {
		log.Printf("[otel] jaeger init error: %v", err)
		return func(context.Context) error { return nil }
	}
	tp := sdktrace.NewTracerProvider(sdktrace.WithBatcher(exp))
	otel.SetTracerProvider(tp)
	return tp.Shutdown
}

func NewReader(brokers, topic, group string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:  strings.Split(brokers, ","),
		Topic:    topic,
		GroupID:  group,
		MinBytes: 1,
		MaxBytes: 10e6,
	})
}

func NewWriter(brokers string) *kafka.Writer {
	return &kafka.Writer{Addr: kafka.TCP(strings.Split(brokers, ",")...), AllowAutoTopicCreation: true}
}

func MustJSON(v any) []byte { b, _ := json.Marshal(v); return b }

// Process simulates step logic. Only step 5 honors FAIL_MODE.
// Returns (nextEvent, isFatal).
func Process(step int, failMode string, evt *Event) (*Event, bool) {
	next := *evt
	next.Step = step + 1

	if step != 5 || failMode == "" || failMode == "none" {
		return &next, false
	}

	if strings.HasPrefix(failMode, "flaky:") {
		p, _ := strconv.ParseFloat(strings.TrimPrefix(failMode, "flaky:"), 64)
		if rand.Float64() < p {
			RetriesTotal.WithLabelValues(strconv.Itoa(step), "flaky").Inc()
			return evt, false // retryable
		}
		return &next, false
	}
	if failMode == "retryable" {
		RetriesTotal.WithLabelValues(strconv.Itoa(step), "timeout").Inc()
		time.Sleep(200 * time.Millisecond)
		return evt, false
	}
	if failMode == "fatal" {
		RetriesTotal.WithLabelValues(strconv.Itoa(step), "fatal").Inc()
		return evt, true
	}
	return &next, false
}

// RunStepService runs a consumer->handler->producer loop with DLQ support.
func RunStepService() error {
	ServeMetrics()
	shutdown := InitOTel()
	defer shutdown(context.Background())

	brokers := os.Getenv("KAFKA_BROKERS")
	topicIn := os.Getenv("TOPIC_IN")
	topicOut := os.Getenv("TOPIC_OUT")
	dlqTopic := os.Getenv("DLQ_TOPIC")
	group := os.Getenv("GROUP_ID")
	stepStr := os.Getenv("STEP")
	failMode := os.Getenv("FAIL_MODE")

	if brokers == "" || topicIn == "" || topicOut == "" || group == "" || stepStr == "" || dlqTopic == "" {
		return fmt.Errorf("missing required envs: KAFKA_BROKERS, TOPIC_IN, TOPIC_OUT, DLQ_TOPIC, GROUP_ID, STEP")
	}
	step, _ := strconv.Atoi(stepStr)

	reader := NewReader(brokers, topicIn, group)
	writer := NewWriter(brokers)
	defer reader.Close()

	tracer := otel.Tracer(fmt.Sprintf("saga-step-%d", step))

	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("[step%d] read error: %v", step, err)
			continue
		}
		var evt Event
		if err := json.Unmarshal(m.Value, &evt); err != nil {
			log.Printf("[step%d] bad json: %v", step, err)
			continue
		}

		ctx, span := tracer.Start(context.Background(), "handle",
			sdktrace.WithAttributes(
				attribute.String("saga_id", evt.SagaID),
				attribute.Int("step", step),
			),
		)
		t0 := time.Now()
		next, fatal := Process(step, failMode, &evt)
		StepLatency.WithLabelValues(strconv.Itoa(step)).Observe(time.Since(t0).Seconds())
		span.End()

		msg := kafka.Message{
			Key:   m.Key, // preserve per-saga ordering
			Value: MustJSON(next),
			Headers: append(m.Headers, kafka.Header{Key: "x-saga-id", Value: []byte(evt.SagaID)}),
		}

		if fatal {
			// Send to DLQ; remember original topic for replay
			msg.Topic = dlqTopic
			msg.Headers = append(msg.Headers, kafka.Header{Key: "x-original-topic", Value: []byte(topicIn)})
			if err := writer.WriteMessages(context.Background(), msg); err != nil {
				log.Printf("[step%d] dlq produce err: %v", step, err)
			}
			DLQTotal.WithLabelValues(dlqTopic).Inc()
			continue
		}

		msg.Topic = topicOut
		if err := writer.WriteMessages(ctx, msg); err != nil {
			RetriesTotal.WithLabelValues(strconv.Itoa(step), "produce_error").Inc()
			log.Printf("[step%d] produce err: %v", step, err)
			time.Sleep(time.Second)
		}
	}
}

// RunEmitter emits StartSaga events to TOPIC_OUT.
func RunEmitter() error {
	ServeMetrics()
	shutdown := InitOTel()
	defer shutdown(context.Background())

	brokers := os.Getenv("KAFKA_BROKERS")
	topic := os.Getenv("TOPIC_OUT")
	rateMs := 1000
	if v := os.Getenv("EMIT_EVERY_MS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil { rateMs = n }
	}
	if brokers == "" || topic == "" {
		return fmt.Errorf("missing envs: KAFKA_BROKERS, TOPIC_OUT")
	}
	writer := NewWriter(brokers)

	ticker := time.NewTicker(time.Duration(rateMs) * time.Millisecond)
	defer ticker.Stop()
	for range ticker.C {
		sagaID := fmt.Sprintf("%d-%d", time.Now().UnixNano(), rand.Intn(100000))
		evt := Event{SagaID: sagaID, Step: 1, SchemaVersion: 1, Ts: time.Now(), Payload: map[string]any{"demo":"start"}}
		msg := kafka.Message{Topic: topic, Key: []byte(sagaID), Value: MustJSON(evt), Headers: []kafka.Header{{Key:"x-saga-id", Value: []byte(sagaID)}}}
		if err := writer.WriteMessages(context.Background(), msg); err != nil {
			log.Printf("[emitter] produce err: %v", err)
		}
	}
	return nil
}

// RunDLQReplayer consumes DLQ and re-emits to original topic header or REPLAY_TARGET.
func RunDLQReplayer() error {
	ServeMetrics()
	shutdown := InitOTel()
	defer shutdown(context.Background())

	brokers := os.Getenv("KAFKA_BROKERS")
	dlqTopic := os.Getenv("DLQ_TOPIC")
	group := os.Getenv("GROUP_ID")
	replayDefault := os.Getenv("REPLAY_TARGET")
	sagaFilter := os.Getenv("SAGA_ID_FILTER") // optional

	if brokers == "" || dlqTopic == "" || group == "" {
		return fmt.Errorf("missing envs: KAFKA_BROKERS, DLQ_TOPIC, GROUP_ID")
	}
	reader := NewReader(brokers, dlqTopic, group)
	writer := NewWriter(brokers)
	defer reader.Close()

	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil { log.Printf("[dlq] read err: %v", err); continue }
		var evt Event
		if err := json.Unmarshal(m.Value, &evt); err != nil { log.Printf("[dlq] bad json: %v", err); continue }

		if sagaFilter != "" && evt.SagaID != sagaFilter { continue }

		orig := replayDefault
		for _, h := range m.Headers {
			if h.Key == "x-original-topic" { orig = string(h.Value) }
		}
		if orig == "" { log.Printf("[dlq] no replay target for saga %s", evt.SagaID); continue }

		msg := kafka.Message{Topic: orig, Key: m.Key, Value: m.Value, Headers: m.Headers}
		if err := writer.WriteMessages(context.Background(), msg); err != nil {
			log.Printf("[dlq] produce err: %v", err)
		} else {
			log.Printf("[dlq] replayed saga=%s to %s", evt.SagaID, orig)
		}
	}
}
