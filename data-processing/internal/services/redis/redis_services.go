package redis

import (
	redisDb "github.com/ginos1998/financing-market-monitor/data-processing/config/redis"
	redisRepository "github.com/ginos1998/financing-market-monitor/data-processing/internal/db/redis/intra_day"
	"github.com/ginos1998/financing-market-monitor/data-processing/internal/models/dtos"
)

func CreateOrUpdateSymbolPrices(redisClient *redisDb.RedisClient, symbol string, lastPrice float64) error {
	symbolPrices, err := redisRepository.GetIntraDaySymbolPrices(redisClient, symbol)
	if err != nil {
		err = createSymbolPrices(redisClient, symbol, lastPrice)
		if err != nil {
			return err
		}
	}
	err = updateSymbolPrices(redisClient, symbol, lastPrice, symbolPrices)
	if err != nil {
		return err
	}
	return nil
}

func updateSymbolPrices(redisClient *redisDb.RedisClient, symbol string, lastPrice float64, symbolPrices string) error {
	symbolValues := dtos.IntraDayPrices{}
	err := symbolValues.FromJSON(symbolPrices)
	if err != nil {
		return err
	}
	symbolValues.Current = lastPrice
	sortSymbolPrices(&symbolValues)
	symbolValuesJson, err := symbolValues.ToJSON()
	if err != nil {
		return err
	}

	err = redisRepository.SetIntraDaySymbolPrices(redisClient, symbol, symbolValuesJson)
	if err != nil {
		return err
	}

	return nil
}

func createSymbolPrices(redisClient *redisDb.RedisClient, symbol string, lastPrice float64) error {
	newSymbolValues := dtos.NewIntraDayPrices(lastPrice, lastPrice, lastPrice, lastPrice)
	newSymbolValuesJson, err := newSymbolValues.ToJSON()
	if err != nil {
		return err
	}
	err = redisRepository.SetIntraDaySymbolPrices(redisClient, symbol, newSymbolValuesJson)
	if err != nil {
		return err
	}
	return nil
}

func sortSymbolPrices(symbolValues *dtos.IntraDayPrices) {
	if symbolValues.Current < symbolValues.Low {
		symbolValues.Low = symbolValues.Current
	}
	if symbolValues.Current > symbolValues.High {
		symbolValues.High = symbolValues.Current
	}
	if symbolValues.Open == 0 {
		symbolValues.Open = symbolValues.Current
	}
}
