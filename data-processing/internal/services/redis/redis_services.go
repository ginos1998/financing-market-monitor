package redis

import (
	"errors"
	"fmt"
	"time"

	redisDb "github.com/ginos1998/financing-market-monitor/data-processing/config/redis"
	redisRepository "github.com/ginos1998/financing-market-monitor/data-processing/internal/db/redis/intra_day"
	"github.com/ginos1998/financing-market-monitor/data-processing/internal/models/dtos"
)

var timeSeriesTypes = []string{"DAILY", "WEEKLY"}

func FlushDB(redisClient *redisDb.RedisClient) error {
	err := redisRepository.FlushDB(redisClient)
	if err != nil {
		return errors.New("Error flushing redis db: " + err.Error())
	}
	return nil
}

func SetSymbolsTimeSeries(redisClient *redisDb.RedisClient, tickers []dtos.Ticker) error {
	for _, ticker := range tickers {
		for _, timeSeriesType := range timeSeriesTypes {
			symbolKey := fmt.Sprintf("%s_%s_%s", timeSeriesType, ticker.Symbol, time.Now().Format("2006-01-02"))
			closePrices := make([]float64, 0)
			if timeSeriesType == "DAILY" {
				for _, ts := range ticker.TimeSeriesDaily.TimeSeriesData {
					closePrices = append(closePrices, ts.Close)
				}
			} else {
				for _, ts := range ticker.TimeSeriesWeekly.TimeSeriesData {
					closePrices = append(closePrices, ts.Close)
				}
			}
			err := redisRepository.SetTimeSeriesData(redisClient, symbolKey, closePrices)
			if err != nil {
				return errors.New("Error setting time series data: " + err.Error())
			}
		}
	}
	return nil
}

func GetIntraDaySymbolPrices(redisClient *redisDb.RedisClient, symbol string) (string, error) {
	symbolPrices, err := redisRepository.GetIntraDaySymbolPrices(redisClient, symbol)
	if err != nil {
		return "", errors.New("Error getting intra day symbol prices: " + err.Error())
	}
	return symbolPrices, nil
}

func CreateOrUpdateSymbolPrices(redisClient *redisDb.RedisClient, symbol string, lastPrice float64, updatesPrevious bool) error {
	symbolKey := fmt.Sprintf("%s_%s", "INTRA-DAY", symbol)
	symbolPrices, err := GetIntraDaySymbolPrices(redisClient, symbolKey)
	if err != nil || symbolPrices == "" {
		err = createSymbolPrices(redisClient, symbol, lastPrice)
		if err != nil {
			return errors.New("Error creating symbol prices: " + err.Error())
		}
		return nil
	}
	err = updateSymbolPrices(redisClient, symbolKey, lastPrice, symbolPrices, updatesPrevious)
	if err != nil {
		return errors.New("Error updating symbol prices: " + err.Error())
	}
	return nil
}

func updateSymbolPrices(redisClient *redisDb.RedisClient, symbol string, lastPrice float64, symbolPrices string, updatesPrevious bool) error {
	symbolValues := dtos.IntraDayPrices{}
	err := symbolValues.FromJSON(symbolPrices)
	if err != nil {
		return err
	}
	if updatesPrevious {
		symbolValues.Previous = symbolValues.Current
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
		return errors.New("Error marshalling new symbol values: " + err.Error())
	}
	symbolKey := fmt.Sprintf("%s_%s", "INTRA-DAY", symbol)
	err = redisRepository.SetIntraDaySymbolPrices(redisClient, symbolKey, newSymbolValuesJson)
	if err != nil {
		return errors.New("Error setting new symbol values: " + err.Error())
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

func GetDailySymbolsPrices(redisClient *redisDb.RedisClient) (map[string][]float64, error) {
	symbolPricesDailyMap, err := redisRepository.GetKeysStartWith(*redisClient, "DAILY_")
	if err != nil {
		return nil, err
	}
	return symbolPricesDailyMap, nil
}
