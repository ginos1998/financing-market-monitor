package consumers

import (
	"errors"
	"github.com/sirupsen/logrus"

	"github.com/ginos1998/financing-market-monitor/data-processing/config/server"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type Consumer struct {
	consumer *kafka.Consumer
}

var logger *logrus.Logger
var topicsMap map[string]string

func CrearteKafkaConsumer(server server.Server) (*Consumer, error) {
	loadKafkaTopics(server.EnvVars)
	kafkaServer := server.EnvVars["KAFKA_SERVER"]
	kafkaGroupId := server.EnvVars["KAFKA_GROUP_ID"]
	if kafkaServer == "" || kafkaGroupId == "" {
		return nil, errors.New("KAFKA_SERVER or KAFKA_GROUP_ID environment variable not set")
	}

	logger = server.Logger

	config := &kafka.ConfigMap{
		"bootstrap.servers": kafkaServer,
		"group.id":          kafkaGroupId,
		"auto.offset.reset": "earliest",
		"message.max.bytes": 3000000,
	}

	consumer, err := kafka.NewConsumer(config)

	return &Consumer{
		consumer: consumer,
	}, err
}

func (c *Consumer) Close() {
	err := c.consumer.Close()
	if err != nil {
		logger.Error("Error closing consumer: ", err)
		return
	}
}

func loadKafkaTopics(envVars map[string]string) {
	topicsMap = map[string]string{
		"KAFKA_TOPIC_STOCK_MARKET_DATA":     envVars["KAFKA_TOPIC_STOCK_MARKET_DATA"],
		"KAFKA_TOPIC_HIST_DAILY_STOCK_DATA": envVars["KAFKA_TOPIC_HIST_DAILY_STOCK_DATA"],
	}
}
