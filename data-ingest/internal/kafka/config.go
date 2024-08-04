package kafka

import (
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/ginos1998/financing-market-monitor/data-ingest/config"
	
	log "github.com/sirupsen/logrus"
)

func CreateKafkaProducer() (*kafka.Producer, error) {
	server := config.GetEnvVar("KAFKA_SERVER")

	config := &kafka.ConfigMap{
		"bootstrap.servers": server,
		"message.max.bytes": 3000000,
	}
	producer, err := kafka.NewProducer(config)
	if err != nil {
		return nil, err
	}

	log.Info("Kafka producer created")

	return producer, nil
}