package consumers

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/ginos1998/financing-market-monitor/data-processing/internal/db/mongod/tickers"

	"github.com/ginos1998/financing-market-monitor/data-processing/config/mongod"
	"github.com/ginos1998/financing-market-monitor/data-processing/internal/models/dtos"
)

func (c *Consumer) InitHistoricalStockDataConsumer(ctx context.Context, mongoRepository mongod.MongoRepository) error {
	topic := topicsMap["KAFKA_TOPIC_TIME_SERIES_DATA"]
	if topic == "" {
		logger.Fatal("KAFKA_TOPIC_TIME_SERIES_DATA not set")
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
			updateTickerTimeSeriesData(apiResponse, mongoRepository)
		}
	}
}

func updateTickerTimeSeriesData(data dtos.Data, mongoRepository mongod.MongoRepository) {
	var ticker = dtos.Ticker{
		Symbol: data.Symbol,
	}

	if data.TimeSeriesType == "1d" {
		ticker.TimeSeriesDaily = data
		err := tickers.UpdateTickerTimeSeriesDaily(mongoRepository, ticker)
		if err != nil {
			logger.Error("Failed to update ticker ", ticker.Symbol, " daily time series data: ", err)
		}
	} else if data.TimeSeriesType == "1wk" {
		ticker.TimeSeriesWeekly = data
		err := tickers.UpdateTickerTimeSeriesWeekly(mongoRepository, ticker)
		if err != nil {
			logger.Error("Failed to update ticker ", ticker.Symbol, " weekly time series data: ", err)
		}
	} else {
		logger.Error("Invalid time series type: ", data.TimeSeriesType)
		return
	}

	logger.Info("Ticker ", ticker.Symbol, " ", data.TimeSeriesType, " time series data updated successfully")
}
