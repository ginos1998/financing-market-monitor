package producers

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/ginos1998/financing-market-monitor/data-ingest/config/server"

	"github.com/gorilla/websocket"
)

var topic string

func (p *KafkaProducer) InitStreamStockMarketDataProducer(server server.Server) {
	topic = server.EnvVars["KAFKA_TOPIC_STREAM_STOCK_MARKET_DATA"]
	if topic == "" {
		server.Logger.Fatal("KAFKA_TOPIC_STREAM_STOCK_MARKET_DATA not set")
	}

	ws, err := initFinnhubWebSocket(server.EnvVars)
	if err != nil {
		server.Logger.Fatal("Error initializing Finnhub WebSocket: ", err)
	}

	err = p.streamStockMarketData(ws)
	if err != nil {
		server.Logger.Fatal("Error streaming stock market data: ", err)
	}
}

func (p *KafkaProducer) streamStockMarketData(ws *websocket.Conn) error {
	defer func(ws *websocket.Conn) {
		err := ws.Close()
		if err != nil {
			logger.Error("Failed to close Finnhub WebSocket: ", err)
		}
	}(ws)

	logger.Info("Streaming stock market data")

	var msg interface{}
	for {
		err := ws.ReadJSON(&msg)
		if err != nil {
			logger.Error("Finnhub Websocket: Failed to read message: ", err)
			panic(err)
		}
		jsonMsg, _ := json.Marshal(msg)

		err = p.produceOnTopic(topic, jsonMsg)
		if err != nil {
			logger.Error(topic, " | Failed to produce message: ", err)
			break
		}
		logger.Info(topic + " | Message send")

		time.Sleep(10 * time.Second)
	}

	p.FlushAndCloseKafkaProducer()
	return nil
}

func initFinnhubWebSocket(envVars map[string]string) (*websocket.Conn, error) {
	finnhubToken := envVars["FINNHUB_TOKEN"]
	if finnhubToken == "" {
		logger.Fatalf("FINNHUB_TOKEN not set")
	}

	urlStr := fmt.Sprintf("wss://ws.finnhub.io?token=%s", finnhubToken)
	ws, _, err := websocket.DefaultDialer.Dial(urlStr, nil)
	if err != nil {
		logger.Fatalf("Failed to connect to Finnhub WebSocket: %v", err)
	}

	symbols := []string{"BINANCE:BTCUSDT"}
	for _, s := range symbols {
		msg, _ := json.Marshal(map[string]interface{}{"type": "subscribe", "symbol": s})
		err := ws.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Failed to subscribe to symbol: " + s))
		}
	}
	logger.Info("Finnhub WebSocket connected")

	return ws, nil
}

func (p *KafkaProducer) FlushAndCloseKafkaProducer() {
	logger.Info("Flushing and closing Kafka producer...")
	p.producer.Flush(11 * 1000)
	p.producer.Close()
	logger.Info("Kafka producer flushed and closed")
}
