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
