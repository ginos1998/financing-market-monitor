package producers

import (
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func produceOnTopic(producer *kafka.Producer, topic string, data []byte) error {
	return producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          data,
	}, nil)
}