package consumers

import (
	"encoding/json"
	"fmt"
	"context"

	"github.com/ginos1998/financing-market-monitor/data-processing/config"
	"github.com/ginos1998/financing-market-monitor/data-processing/internal/models/dtos"

	log "github.com/sirupsen/logrus"
)

func (c *Consumer)InitStockMarketDataConsumer(ctx context.Context) error {
	topic := config.GetEnvVar("KAFKA_TOPIC_STOCK_MARKET_DATA")

	err := c.consumer.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		panic(err)
	}
	log.Info("Consumer subscribed to topic: " + topic)

	for {
		select {
			case <-ctx.Done():
				log.Println("Context done, stopping consumer...")
            	return c.consumer.Close()
			
			default:
				msg, err := c.consumer.ReadMessage(-1)
				if err != nil {
					if ctx.Err() != nil {
						log.Println("Context canceled, exiting...")
						return nil
					}
					log.Errorf("Error reading message: %v (%v)\n", err, msg)
					continue
				}

				var tradesData dtos.WsData
				err = json.Unmarshal(msg.Value, &tradesData)
				if err != nil {
					log.Error("Failed to unmarshal stock market data")
					continue
				}
				if len(tradesData.Trades) == 0 {
					log.Warn("No trades data found")
					continue
				}
				
				log.Info(fmt.Sprintf("Message on %s: Symbol %s LastPrice %v\n", msg.TopicPartition, tradesData.Trades[0].Symbol, tradesData.Trades[0].LastPrice))
		}
	}
}