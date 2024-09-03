package producers

import (
	"errors"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/ginos1998/financing-market-monitor/data-ingest/config/server"
	"github.com/sirupsen/logrus"
)

type KafkaProducer struct {
	producer *kafka.Producer
}

var logger logrus.Logger

func CreateKafkaProducer(server server.Server) (*KafkaProducer, error) {
	logger = *server.Logger
	kafkaServer := server.EnvVars["KAFKA_SERVER"]
	if kafkaServer == "" {
		return nil, errors.New("KAFKA_SERVER environment variable not set")
	}

	config := &kafka.ConfigMap{
		"bootstrap.servers": kafkaServer,
		"message.max.bytes": 3000000,
	}

	producer, err := kafka.NewProducer(config)
	if err != nil {
		return nil, err
	}

	return &KafkaProducer{producer: producer}, nil
}
