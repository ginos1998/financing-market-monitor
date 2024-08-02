package main

import (
	appCfg 			"github.com/ginos1998/financing-market-monitor/data-ingest/config"
	appKafka 		"github.com/ginos1998/financing-market-monitor/data-ingest/internal/kafka"
	kafkaProducers 	"github.com/ginos1998/financing-market-monitor/data-ingest/internal/kafka/producers"
	log 			"github.com/sirupsen/logrus"
)

func main() {
	doneChan := appCfg.InitSignalHandler()
	err := appCfg.LoadEnvVars()
	if err != nil { 
		panic(err) 
	}
	producer, _ := appKafka.CreateKafkaProducer()
	go kafkaProducers.InitStreamStockMarketDataProducer(producer)

	//import_data.UpdateCedearTimeSeriesData(mongoClient, "AAPL")

	log.Info("Press Ctrl+C to exit...")
	<-doneChan
	kafkaProducers.FlushAndCloseKafkaProducer()
	log.Info("Exiting...")
}