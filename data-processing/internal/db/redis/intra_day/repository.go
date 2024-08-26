package intra_day

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	rdb "github.com/ginos1998/financing-market-monitor/data-processing/config/redis"
)

func FlushDB(redisClient *rdb.RedisClient) error {
	ctx := context.Background()
	err := redisClient.Client.FlushDB(ctx).Err()
	if err != nil {
		return errors.New("Error flushing redis db: " + err.Error())
	}
	return nil
}

func GetIntraDaySymbolPrices(redisClient *rdb.RedisClient, symbol string) (string, error) {
	ctx := context.Background()
	prices, err := redisClient.Client.Get(ctx, symbol).Result()
	if err != nil {
		return "", errors.New("Error getting intra day symbol prices: " + err.Error())
	}
	return prices, nil
}

func SetTimeSeriesData(redisClient *rdb.RedisClient, symbolKey string, timeSeriesData []float64) error {
	ctx := context.Background()
	expiration := 8 * time.Hour
	jsonData, err := json.Marshal(timeSeriesData)
	if err != nil {
		return errors.New("Error marshalling time series data: " + err.Error())
	}
	err = redisClient.Client.Set(ctx, symbolKey, jsonData, expiration).Err()
	if err != nil {
		return errors.New("Error setting time series data: " + err.Error())
	}
	return nil
}

func SetIntraDaySymbolPrices(redisClient *rdb.RedisClient, symbolKey string, prices string) error {
	ctx := context.Background()
	err := redisClient.Client.Set(ctx, symbolKey, prices, 0).Err()
	if err != nil {
		return errors.New("Error setting intra day symbol prices: " + err.Error())
	}
	return nil
}

func GetKeysStartWith(redisClient rdb.RedisClient, prefix string) (map[string][]float64, error) {
	ctx := context.Background()
	keyPattern := prefix + "*"
	keys, err := redisClient.Client.Keys(ctx, keyPattern).Result()
	if err != nil {
		return nil, errors.New("Error getting keys: " + err.Error())
	}
	symbolPricesMap := make(map[string][]float64)
	for _, key := range keys {
		symbolPrices, err := redisClient.Client.Get(ctx, key).Result()
		if err != nil {
			return nil, errors.New("Error getting symbol prices: " + err.Error())
		}
		prices := make([]float64, 0)
		err = json.Unmarshal([]byte(symbolPrices), &prices)
		if err != nil {
			return nil, errors.New("Error parsing symbol prices: " + err.Error())
		}
		symbolPricesMap[key] = prices
	}
	return symbolPricesMap, nil

}
