package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	appCfg "github.com/ginos1998/financing-market-monitor/data-processing/config"
	kafkaConsumer "github.com/ginos1998/financing-market-monitor/data-processing/internal/kafka/consumers"
    mdb "github.com/ginos1998/financing-market-monitor/data-processing/internal/db/mongod"

	"github.com/sirupsen/logrus"
)

func main() {
	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
		TimestampFormat: time.RFC3339,
	})

	err := appCfg.LoadEnvVars()
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

	sigs := make(chan os.Signal, 1)
    signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	mongocli, err := mdb.CreateMongoClient()
	if err != nil {
		log.Fatal("Error creating MongoDB client: ", err)
	}
	

	consumer, err := kafkaConsumer.CrearteKafkaConsumer()
	if err != nil {
		log.Fatalf("Error creating kafka consumer: %v", err)
	}

	go func() {
        if err := consumer.InitStockMarketDataConsumer(ctx); err != nil {
            log.Fatalf("Error running consumer: %v", err)
        }
    }()
	
	hsdConsumer, err := kafkaConsumer.CrearteKafkaConsumer()
	if err != nil {
		log.Fatalf("Error creating kafka historical stock data consumer: %v", err)
	}
	go func() {
		if err := hsdConsumer.InitHistoricalStockDataConsumer(ctx, mongocli); err != nil {
			log.Fatalf("Error running consumer: %v", err)
		}
	}()

	<-sigs
    log.Println("Shutting down...")

	cancel()
	time.Sleep(5 * time.Second)
    log.Println("Consumer has shut down")
}