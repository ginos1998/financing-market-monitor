package producers

import (
	appCfg 	"github.com/ginos1998/financing-market-monitor/data-ingest/config"
	mdb 	"github.com/ginos1998/financing-market-monitor/data-ingest/internal/db/mongod"
	dtos	"github.com/ginos1998/financing-market-monitor/data-ingest/internal/models/dtos"
	apis 	"github.com/ginos1998/financing-market-monitor/data-ingest/internal/apis"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
)

func InitHistStockDataProducer(producer *kafka.Producer, alphaAlphavantageAPI *apis.AlphavantageAPI, mongoClient *mdb.MongoRepository) {
	log.Info("Cron <updateHistoricalStockData> created. Schedule: every day at 7:00 AM")
	c := cron.New()
	c.AddFunc("0 7 * * *", // every day at 7:00 AM
		func() {
			updateHistoricalStockData(producer, alphaAlphavantageAPI, mongoClient)
		})
	c.Start()
}

func updateHistoricalStockData(producer *kafka.Producer, alphavantageAPI *apis.AlphavantageAPI, mongoClient *mdb.MongoRepository) {
	topic := appCfg.GetEnvVar("KAFKA_TOPIC_HIST_DAYLI_STOCK_DATA")
	if topic == "" {
		log.Error("envvar KAFKA_TOPIC_HIST_DAYLI_STOCK_DATA not set")
		return
	}

	cedears, err := mongoClient.GetCedearsWithoutHistoricalDayliStockData()
	if err != nil {
		log.Error("Error getting cedears without historical daily stock data: ", err)
		return
	}

	if len(cedears) == 0 {
		log.Info("No cedears without historical daily stock data")
		return
	}

	cedearsToUpdate := cedears[:alphavantageAPI.RequestPerDay]

	log.Info("Updating historical stock data...")
	log.Info("Cedears to update: ", len(cedearsToUpdate))

	var cedearsNotUpdated []dtos.Cedear = make([]dtos.Cedear, 0)

	for idx, cedear := range cedearsToUpdate {
		log.Info("Getting data of ", cedear.Ticker, " from Alphavantage API")
		
		res, err := alphavantageAPI.GetTickerDailyHistoricalData(cedear.Ticker)
		if err != nil {
			log.Error("Error getting data of ", cedear.Ticker, " from Alphavantage API: ", err)
			cedearsNotUpdated = cedearsToUpdate[idx:]
			break
		}
		// TODO historical WEEKLY data
		
		log.Info("Sending data of ", cedear.Ticker, " to Kafka...")
		error := producer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          res,
		}, nil)

		if error != nil {
			log.Error("Error sending data of ", cedear.Ticker, " to Kafka: ", error)
			cedearsNotUpdated = append(cedearsNotUpdated, cedear)
		}
		log.Info("Data of ", cedear.Ticker, " sent to Kafka successfully")

	}

	log.Info("The process has updated ", len(cedearsToUpdate) - len(cedearsNotUpdated), " cedears")
	if len(cedearsNotUpdated) > 0 {
		var cedearsNotUpdatedTickers []string
		for _, cedear := range cedearsNotUpdated {
			cedearsNotUpdatedTickers = append(cedearsNotUpdatedTickers, cedear.Ticker)
		}
		log.Warn("The following cedears could not be updated: ", cedearsNotUpdatedTickers)
	}
}