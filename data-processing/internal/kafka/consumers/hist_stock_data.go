package consumers

import (
	"context"
	"encoding/json"
	"errors"

	appCfg "github.com/ginos1998/financing-market-monitor/data-processing/config"
	mdb "github.com/ginos1998/financing-market-monitor/data-processing/internal/db/mongod"
	dtos "github.com/ginos1998/financing-market-monitor/data-processing/internal/models/dtos"

	log "github.com/sirupsen/logrus"
)

func (c *Consumer) InitHistoricalStockDataConsumer(ctx context.Context, mongoCli *mdb.MongoRepository) error {
	topic := appCfg.GetEnvVar("KAFKA_TOPIC_HIST_DAILY_STOCK_DATA")
	if topic == "" {
		log.Fatal("KAFKA_TOPIC_HIST_DAILY_STOCK_DATA not set")
	}

	err := c.consumer.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		return errors.New("Failed to subscribe to topic: " + topic)
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

			var apiResponse dtos.Data
			err = json.Unmarshal(msg.Value, &apiResponse)
			if err != nil {
				log.Error("Failed to unmarshal historical stock data: ", err)
				continue
			}

			log.Info("Message reveived on topic ", msg.TopicPartition, ": ", apiResponse.Symbol)
			updateCedearTimeSeriesData(apiResponse, mongoCli)
		}
	}
}

func updateCedearTimeSeriesData(data dtos.Data, mongoCli *mdb.MongoRepository) {
	var cedear dtos.Cedear = dtos.Cedear{
		Ticker:          data.Symbol,
		TimeSeriesDayli: data,
	}

	err := mongoCli.UpdateCedearTimeSeriesData(cedear)
	if err != nil {
		log.Error("Failed to update cedear ", cedear.Ticker, " time series data: ", err)
	}
	log.Info("Cedear ", cedear.Ticker, " time series data updated")
}
