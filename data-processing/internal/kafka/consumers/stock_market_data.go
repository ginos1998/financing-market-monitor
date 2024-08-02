package consumers

import (
	"encoding/json"
	"fmt"

	"github.com/ginos1998/financing-market-monitor/data-processing/config"
	"github.com/ginos1998/financing-market-monitor/data-processing/internal/models/dtos"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()
var kafkaConsumer *kafka.Consumer

func InitStockMarketDataConsumer(consumer *kafka.Consumer) {
	kafkaConsumer = consumer
	topic := config.GetEnvVar("KAFKA_STOCK_MARKET_DATA_TOPIC")

	err := kafkaConsumer.SubscribeTopics([]string{topic}, nil)

	if err != nil {
		panic(err)
	}
	
	log.Info("Consumer subscribed to topic: " + topic)

	for {
		if kafkaConsumer.IsClosed() {
			log.Info("Consumer closed")
			break
		}

		msg, err := kafkaConsumer.ReadMessage(-1)
		if err != nil {
			log.Errorf("Consumer error: %v (%v)\n", err, msg)
			break
		}

		var tradesData dtos.WsData
		err = json.Unmarshal(msg.Value, &tradesData)
		if err != nil {
			log.Error("Failed to unmarshal stock market data")
			continue
		}
		if len(tradesData.Trades) == 0 {
			log.Error("No trades data found")
			continue
		}
		
		log.Info(fmt.Sprintf("Message on %s: Symbol %s LastPrice %v\n", msg.TopicPartition, tradesData.Trades[0].Symbol, tradesData.Trades[0].LastPrice))
	}
}

func CloseKafkaConsumer() {
	if !kafkaConsumer.IsClosed() {
		kafkaConsumer.Close()
		log.Info("Kafka consumer closed")
	}
}
