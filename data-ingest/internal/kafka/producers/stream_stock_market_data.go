package producers

import (
	"encoding/json"
	"time"
	"fmt"
	"errors"

	appCfg "github.com/ginos1998/financing-market-monitor/data-ingest/config"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

var kafkaProducer *kafka.Producer

func InitStreamStockMarketDataProducer(producer *kafka.Producer) {
	kafkaProducer = producer

	err := streamStockMarketData()
	if err != nil {
		log.Fatal("Error streaming stock market data: ", err)
	}
}

func streamStockMarketData() error {
	topic := appCfg.GetEnvVar("KAFKA_TOPIC_STREAM_STOCK_MARKET_DATA")
	if topic == "" {
		return errors.New("KAFKA_TOPIC_STREAM_STOCK_MARKET_DATA not set")
	}

	ws := initFinnhubWebSocket()
	defer ws.Close()

	log.Info("Streaming stock market data")

	var msg interface{}
	for {
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Error("Failed to read message: ", err)
			panic(err)
		}
		jsonMsg, _ := json.Marshal(msg)

		err = kafkaProducer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          jsonMsg,
		}, nil)
		if err != nil {
			log.Error("Failed to produce message: ", err)
			break
		}

		log.Info("Message on topic ", topic, ": ", msg)

		time.Sleep(10 * time.Second)
	}

	FlushAndCloseKafkaProducer()
	return nil
}

func initFinnhubWebSocket() *websocket.Conn {
	finnhub_token := appCfg.GetEnvVar("FINNHUB_TOKEN")
	if finnhub_token == "" {
		log.Fatalf("FINNHUB_TOKEN not set")
	}

	urlStr := fmt.Sprintf("wss://ws.finnhub.io?token=%s", finnhub_token)
	
	ws, _, err := websocket.DefaultDialer.Dial(urlStr, nil)

	if err != nil {
		log.Fatalf("Failed to connect to Finnhub WebSocket: %v", err)
	}

	symbols := []string{"BINANCE:BTCUSDT"}
	for _, s := range symbols {
		msg, _ := json.Marshal(map[string]interface{}{"type": "subscribe", "symbol": s})
		ws.WriteMessage(websocket.TextMessage, msg)
	}

	log.Info("Finnhub WebSocket connected")

	return ws
}

func FlushAndCloseKafkaProducer() {
	log.Info("Flushing and closing Kafka producer...")
	kafkaProducer.Flush(11 * 1000)
	kafkaProducer.Close()
	log.Info("Kafka producer flushed and closed")
}
