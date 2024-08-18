package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	srv "github.com/ginos1998/financing-market-monitor/data-processing/config/server"
	indicatorsCron "github.com/ginos1998/financing-market-monitor/data-processing/internal/crons/indicators"
	kafkaConsumer "github.com/ginos1998/financing-market-monitor/data-processing/internal/kafka/consumers"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	server := srv.NewServer()
	server.Logger.Info("Server configured")

	err := indicatorsCron.CalculateMovingAverages(server)
	if err != nil {
		server.Logger.Fatalf("Error starting cron: %v", err)
	}

	// consumes stock market data
	consumer, err := kafkaConsumer.CrearteKafkaConsumer(*server)
	if err != nil {
		server.Logger.Fatalf("Error creating kafka consumer: %v", err)
	}
	go func() {
		if err := consumer.InitStockMarketDataConsumer(ctx, server.RedisClient); err != nil {
			server.Logger.Fatalf("Error running consumer: %v", err)
		}
	}()

	// consumes times series data
	hsdConsumer, err := kafkaConsumer.CrearteKafkaConsumer(*server)
	if err != nil {
		server.Logger.Fatalf("Error creating kafka historical stock data consumer: %v", err)
	}
	go func() {
		if err := hsdConsumer.InitHistoricalStockDataConsumer(ctx, server.MongoRepository); err != nil {
			server.Logger.Fatalf("Error running consumer: %v", err)
		}
	}()

	server.Logger.Info("Server started")
	<-sigs
	server.Logger.Info("Shutting down...")

	cancel()
	time.Sleep(5 * time.Second)
	server.Logger.Info("Consumer has shut down")
}
