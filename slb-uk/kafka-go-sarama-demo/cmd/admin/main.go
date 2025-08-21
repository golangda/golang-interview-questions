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
