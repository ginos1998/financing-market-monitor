package producers

import (
	"encoding/json"
	"time"
	"fmt"
	"errors"

	appCfg "github.com/ginos1998/financing-market-monitor/data-ingest/config"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()
var kafkaProducer *kafka.Producer

func InitStreamStockMarketDataProducer(producer *kafka.Producer) {
	kafkaProducer = producer

	streamStockMarketData()
}

func streamStockMarketData() {
	topic := appCfg.GetEnvVar("KAFKA_TOPIC_STREAM_STOCK_MARKET_DATA")

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

		log.Info("Message: ", msg)

		err = kafkaProducer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          jsonMsg,
		}, nil)
		if err != nil {
			log.Error("Failed to produce message: ", err)
			break
		}

		time.Sleep(10 * time.Second)
	}

	FlushAndCloseKafkaProducer()
}

func initFinnhubWebSocket() *websocket.Conn {
	finnhub_token := appCfg.GetEnvVar("FINNHUB_TOKEN")
	urlStr := fmt.Sprintf("wss://ws.finnhub.io?token=%s", finnhub_token)
	
	ws, _, err := websocket.DefaultDialer.Dial(urlStr, nil)

	if err != nil {
		panic(errors.New("Failed to connect to Finnhub WebSocket: " + err.Error()))
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
	log.Info("Flushing and closing Kafka producer")

	kafkaProducer.Flush(11 * 1000)
	kafkaProducer.Close()
}
