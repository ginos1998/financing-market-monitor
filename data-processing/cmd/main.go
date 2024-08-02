package main

import (
	appCfg "github.com/ginos1998/financing-market-monitor/data-processing/config"
	appKafka "github.com/ginos1998/financing-market-monitor/data-processing/internal/kafka"
	"github.com/ginos1998/financing-market-monitor/data-processing/internal/kafka/consumers"
	log "github.com/sirupsen/logrus"
)

func main() {
	doneChann := appCfg.InitSignalHandler()
	err := appCfg.LoadEnvVars()
	if err != nil {
		panic(err)
	}

	consumer := appKafka.CrearteKafkaConsumer()
	go consumers.InitStockMarketDataConsumer(consumer)

	log.Info("Press Ctrl+C to exit...")
	<-doneChann
	consumers.CloseKafkaConsumer()
	log.Info("Exiting...")
}