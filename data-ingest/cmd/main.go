package main

import (
	appCfg "github.com/ginos1998/financing-market-monitor/data-ingest/config"
	srv "github.com/ginos1998/financing-market-monitor/data-ingest/config/server"
	crons "github.com/ginos1998/financing-market-monitor/data-ingest/internal/crons/stock_data"
	kafkaProducers "github.com/ginos1998/financing-market-monitor/data-ingest/internal/kafka/producers"
)

func main() {
	doneChan := appCfg.InitSignalHandler()

	server := srv.NewServer()
	server.Logger.Info("Server configured")

	producer, err := kafkaProducers.CreateKafkaProducer(*server)
	if err != nil {
		server.Logger.Fatal("Error creating Kafka producer: ", err)
	}
	server.Logger.Info("Kafka producer created")

	crons.InitHistStockDataProducer(producer, *server)

	go producer.InitStreamStockMarketDataProducer(*server)

	server.Logger.Info("Press Ctrl+C to exit...")
	<-doneChan
	producer.FlushAndCloseKafkaProducer()
	server.Logger.Info("Exiting...")
}
