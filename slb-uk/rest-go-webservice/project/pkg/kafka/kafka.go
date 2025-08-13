package kafkahelper

import (
	"github.com/IBM/sarama"
)

func NewIdempotentProducer(brokers []string) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Idempotent = true
	config.Producer.Return.Successes = true
	config.Net.MaxOpenRequests = 1
	config.Version = sarama.V2_6_0_0

	return sarama.NewSyncProducer(brokers, config)
}
