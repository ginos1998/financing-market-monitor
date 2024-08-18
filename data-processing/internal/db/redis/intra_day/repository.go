package intra_day

import (
	"context"

	rdb "github.com/ginos1998/financing-market-monitor/data-processing/config/redis"
)

func GetIntraDaySymbolPrices(redisClient *rdb.RedisClient, symbol string) (string, error) {
	ctx := context.Background()
	prices, err := redisClient.Client.Get(ctx, symbol).Result()
	if err != nil {
		return "", err
	}
	return prices, nil
}

func SetIntraDaySymbolPrices(redisClient *rdb.RedisClient, symbol string, prices string) error {
	ctx := context.Background()
	err := redisClient.Client.Set(ctx, symbol, prices, 0).Err()
	if err != nil {
		return err
	}
	return nil
}
