package consumers

import (
	"context"
	"encoding/json"
	"errors"
	cryptosRepository "github.com/ginos1998/financing-market-monitor/data-processing/internal/db/mongod/cryptos"

	"github.com/ginos1998/financing-market-monitor/data-processing/config/mongod"
	"github.com/ginos1998/financing-market-monitor/data-processing/internal/models/dtos"
)

func (c *Consumer) ConsumesCryptoDailyTimeSeries(ctx context.Context, mongoRepository mongod.MongoRepository) error {
	topic := topicsMap["KAFKA_TOPIC_CRYPTO_TIME_SERIES_DATA"]
	if topic == "" {
		logger.Fatal("KAFKA_TOPIC_CRYPTO_TIME_SERIES_DATA not set")
	}

	err := c.consumer.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		return errors.New("Failed to subscribe to topic: " + topic)
	}

	logger.Info("Consumer subscribed to topic: " + topic)

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

			var apiResponse dtos.Data
			err = json.Unmarshal(msg.Value, &apiResponse)
			if err != nil {
				logger.Error("Failed to unmarshal historical stock data: ", err)
				continue
			}

			logger.Info("Message received on topic ", msg.TopicPartition, ": ", apiResponse.Symbol)
			updateCryptoTimeSeriesData(apiResponse, mongoRepository)
		}
	}
}

func updateCryptoTimeSeriesData(data dtos.Data, mongoRepository mongod.MongoRepository) {
	var crypto = dtos.Crypto{
		YahooSymbol:     data.Symbol,
		TimeSeriesDaily: data,
	}
	err := cryptosRepository.UpdateCryptoTimeSeriesData(mongoRepository, crypto)
	if err != nil {
		logger.Error("Failed to update crypto ", crypto.YahooSymbol, " time series data: ", err)
	}
	logger.Info("Crypto ", crypto.YahooSymbol, " time series data updated")
}
