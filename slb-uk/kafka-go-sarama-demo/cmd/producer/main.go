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
