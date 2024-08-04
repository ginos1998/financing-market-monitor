package consumers

import (
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/ginos1998/financing-market-monitor/data-processing/config"
	
	log "github.com/sirupsen/logrus"
)

type Consumer struct {
    consumer *kafka.Consumer
}

func CrearteKafkaConsumer() (*Consumer, error) {
	kafka_server := config.GetEnvVar("KAFKA_SERVER")
	kafka_group_id := config.GetEnvVar("KAFKA_GROUP_ID")

	config := &kafka.ConfigMap{
		"bootstrap.servers": kafka_server,
		"group.id":          kafka_group_id,
		"auto.offset.reset": "earliest",
		"message.max.bytes": 3000000,
	}

	consumer, err := kafka.NewConsumer(config)

	if err == nil {
		log.Info("Kafka consumer created")
	}
	
	return &Consumer{
			consumer: consumer,
		}, err
}