package kafka

import (
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/ginos1998/financing-market-monitor/data-processing/config"
)

func CrearteKafkaConsumer() *kafka.Consumer {
	kafka_server := config.GetEnvVar("KAFKA_SERVER")
	kafka_group_id := config.GetEnvVar("KAFKA_GROUP_ID")

	config := &kafka.ConfigMap{
		"bootstrap.servers": kafka_server,
		"group.id":          kafka_group_id,
		"auto.offset.reset": "earliest",
	}

	consumer, err := kafka.NewConsumer(config)

	if err != nil {
		panic(err)
	}

	log.Info("Kafka consumer created")
	
	return consumer
}