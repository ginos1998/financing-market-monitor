package producers

import (
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func (p *KafkaProducer) produceOnTopic(topic string, data []byte) error {
	return p.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          data,
	}, nil)
}
