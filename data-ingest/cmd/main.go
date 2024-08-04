package main

import (
	"time"

	appCfg 			"github.com/ginos1998/financing-market-monitor/data-ingest/config"
	mdb 			"github.com/ginos1998/financing-market-monitor/data-ingest/internal/db/mongod"
	appKafka 		"github.com/ginos1998/financing-market-monitor/data-ingest/internal/kafka"
	kafkaProducers	"github.com/ginos1998/financing-market-monitor/data-ingest/internal/kafka/producers"
	apis 			"github.com/ginos1998/financing-market-monitor/data-ingest/internal/apis"
	
	"github.com/sirupsen/logrus"
)
func main() {
	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
		TimestampFormat: time.RFC3339,
	})

	doneChan := appCfg.InitSignalHandler()
	err := appCfg.LoadEnvVars()
	if err != nil { 
		log.Fatal("Error loading environment variables: ", err)
	}

	mongocli, err := mdb.CreateMongoClient()
	if err != nil {
		log.Fatal("Error creating MongoDB client: ", err)
	}

	alphavantageApi := apis.AlphavantageAPI{}
	err = alphavantageApi.ConfigAlphavantageAPI()
	if err != nil {
		log.Fatal("Error configuring Alphavantage API: ", err)
	}

	producer, err := appKafka.CreateKafkaProducer()
	if err != nil {
		log.Fatal("Error creating Kafka producer: ", err)
	}
	
	go kafkaProducers.InitStreamStockMarketDataProducer(producer)
	go kafkaProducers.InitHistStockDataProducer(producer, &alphavantageApi, mongocli)

	log.Info("Press Ctrl+C to exit...")
	<-doneChan
	kafkaProducers.FlushAndCloseKafkaProducer()
	log.Info("Exiting...")
}