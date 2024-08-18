package main

import (
	"context"
	r "github.com/ginos1998/financing-market-monitor/data-processing/config/redis"
	srv "github.com/ginos1998/financing-market-monitor/data-processing/config/server"
	"github.com/ginos1998/financing-market-monitor/data-processing/internal/models/dtos"
)

func main() {
	server := srv.NewServer()
	server.Logger.Info("Server is running")

	rdb, err := r.NewRedisClient(server.EnvVars)
	if err != nil {
		server.Logger.Fatal(err)
	}
	server.Logger.Info("Redis client is running")
	ctx := context.Background()

	btc := dtos.IntraDayPrices{
		Low:     59256.56,
		High:    60284.99,
		Open:    59491.99,
		Current: 0,
	}
	btcJson, err := btc.ToJSON()
	if err != nil {
		server.Logger.Error(err)
	}

	err = rdb.Client.Set(ctx, "BINANCE:BTCUSDT", btcJson, 0).Err()
	if err != nil {
		server.Logger.Error(err)
	}

	value, err := rdb.Client.Get(ctx, "BINANCE:BTCUSDT").Result()
	if err != nil {
		server.Logger.Error(err)
	}

	btcusdt := dtos.IntraDayPrices{}
	err = btcusdt.FromJSON(value)
	if err != nil {
		server.Logger.Error(err)
	}
	server.Logger.Info("BINANCE:BTCUSDT: ", btcusdt)

	random, err := rdb.Client.Get(ctx, "BINANCE:ETHUSDT").Result()
	if err != nil {
		server.Logger.Fatalf("Failed to get current price from redis: %v", err)
	}
	server.Logger.Info("BINANCE:ETHUSDT: ", random)
}
