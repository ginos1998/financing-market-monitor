package producers

import (
	"time"

	appCfg 	"github.com/ginos1998/financing-market-monitor/data-ingest/config"
	mdb 	"github.com/ginos1998/financing-market-monitor/data-ingest/internal/db/mongod"
	dtos	"github.com/ginos1998/financing-market-monitor/data-ingest/internal/models/dtos"
	apis 	"github.com/ginos1998/financing-market-monitor/data-ingest/internal/apis"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
)

func InitHistStockDataProducer(producer *kafka.Producer, alphaAlphavantageAPI *apis.AlphavantageAPI, mongoClient *mdb.MongoRepository) {
	log.Info("Cron <updateHistoricalStockData> created. Schedule: every day at 9:10 AM")
	useYahooAPI := true
	c := cron.New()
	c.AddFunc("10 9 * * *", // every day at 9:10 AM
		func() {
			log.Info("Cron <updateHistoricalStockData> started at ", time.Now().Format(time.RFC3339))
			updateHistoricalStockData(producer, alphaAlphavantageAPI, mongoClient, useYahooAPI)
			log.Info("Cron <updateHistoricalStockData> finished at ", time.Now().Format(time.RFC3339))
		})
	c.Start()
}

func updateHistoricalStockData(producer *kafka.Producer, alphavantageAPI *apis.AlphavantageAPI, mongoClient *mdb.MongoRepository, useYahooAPI bool) {
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
	log.Info("Cedears without historical daily stock data: ", len(cedears))

	if len(cedears) == 0 {
		log.Info("No cedears without historical daily stock data")
		return
	}

	var cedearsToUpdate []dtos.Cedear
	if useYahooAPI {
		cedearsToUpdate = cedears
	} else {
		if len(cedears) > alphavantageAPI.RequestPerDay {
			cedearsToUpdate = cedears[:alphavantageAPI.RequestPerDay]
		} else {
			cedearsToUpdate = cedears
		}
	}

	log.Info("Updating historical stock data...")
	log.Info("Cedears to update: ", len(cedearsToUpdate))

	var cedearsNotUpdated []dtos.Cedear = make([]dtos.Cedear, 0)
	var res []byte

	for idx, cedear := range cedearsToUpdate {

		if useYahooAPI {
			res, err = apis.GetDayliHistoricalStockData(cedear.Ticker)
			if err != nil {
				log.Error("Error getting data of ", cedear.Ticker, " from Yahoo API: ", err)
				cedearsNotUpdated = append(cedearsNotUpdated, cedear)
				continue
			}
		} else {
			log.Info("Getting data of ", cedear.Ticker, " from Alphavantage API")
			res, err = alphavantageAPI.GetTickerDailyHistoricalData(cedear.Ticker)
			if err != nil {
				log.Error("Error getting data of ", cedear.Ticker, " from Alphavantage API: ", err)
				cedearsNotUpdated = cedearsToUpdate[idx:]
				break
			}
		}

		log.Info("Sending data of ", cedear.Ticker, " to Kafka...")
		error := produceOnTopic(producer, topic, res)
		if error != nil {
			log.Error("Error sending data of ", cedear.Ticker, " to Kafka: ", error)
			cedearsNotUpdated = append(cedearsNotUpdated, cedear)
		}
		log.Info("Data of ", cedear.Ticker, " sent to Kafka successfully")
		time.Sleep(1 * time.Second)

	}

	log.Info("The process has updated ", len(cedearsToUpdate) - len(cedearsNotUpdated), " cedears")
	if len(cedearsNotUpdated) > 0 {
		var cedearsNotUpdatedTickers []string
		for _, cedear := range cedearsNotUpdated {
			cedearsNotUpdatedTickers = append(cedearsNotUpdatedTickers, cedear.Ticker)
		}
		log.Warn("The following cedears could not be updated: ", cedearsNotUpdatedTickers)
		
		if useYahooAPI {
			log.Info("Trying again with Alphavantage API...")
			updateHistoricalStockData(producer, alphavantageAPI, mongoClient, false)
		}
	}

}