package consumers

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	redisDb "github.com/ginos1998/financing-market-monitor/data-processing/config/redis"
	"github.com/ginos1998/financing-market-monitor/data-processing/internal/models/dtos"
	redisServices "github.com/ginos1998/financing-market-monitor/data-processing/internal/services/redis"
)

func (c *Consumer) InitStockMarketDataConsumer(ctx context.Context, redisClient redisDb.RedisClient) error {
	topic := topicsMap["KAFKA_TOPIC_STREAM_STOCK_MARKET_DATA"]
	err := c.consumer.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		logger.Fatalf("Failed to subscribe to topic %s: %v", topic, err)
	}
	logger.Info("Consumer subscribed to topic: " + topic)
	updatesPrevious := true
	startTime := time.Now().Unix()
	for {
		select {
		case <-ctx.Done():
			logger.Println("Context done, stopping consumer...")
			return c.consumer.Close()

		default:
			msg, err := c.consumer.ReadMessage(-1)
			if err != nil {
				if ctx.Err() != nil {
					logger.Println("Context canceled, exiting...")
					return nil
				}
				logger.Errorf("Error reading message: %v (%v)\n", err, msg)
				continue
			}

			var tradesData dtos.WsData
			err = json.Unmarshal(msg.Value, &tradesData)
			if err != nil {
				logger.Error("Failed to unmarshal stock market data")
				continue
			}
			if len(tradesData.Trades) == 0 {
				logger.Warn("No trades data found")
				continue
			}
			err = redisServices.CreateOrUpdateSymbolPrices(&redisClient, tradesData.Trades[0].Symbol, tradesData.Trades[0].LastPrice, updatesPrevious)
			if err != nil {
				logger.Error("redisService: Failed to update symbol prices: ", err)
			}
			logger.Info(fmt.Sprintf("Message on %s: Symbol %s LastPrice %v\n", msg.TopicPartition, tradesData.Trades[0].Symbol, tradesData.Trades[0].LastPrice))
			updatesPrevious = false
			if time.Now().Unix()-startTime >= 30 {
				updatesPrevious = true
				startTime = time.Now().Unix()
			}
		}
	}
}
