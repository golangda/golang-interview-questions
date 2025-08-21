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
