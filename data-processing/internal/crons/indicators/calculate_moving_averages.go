package indicators

import (
	"fmt"
	"strings"

	"github.com/ginos1998/financing-market-monitor/data-processing/config/mongod"
	"github.com/ginos1998/financing-market-monitor/data-processing/config/redis"
	"github.com/ginos1998/financing-market-monitor/data-processing/config/server"
	tickersRepository "github.com/ginos1998/financing-market-monitor/data-processing/internal/db/mongod/tickers"
	redisRepository "github.com/ginos1998/financing-market-monitor/data-processing/internal/db/redis/intra_day"
	"github.com/ginos1998/financing-market-monitor/data-processing/internal/indicators"
	"github.com/ginos1998/financing-market-monitor/data-processing/internal/models/dtos"
	redisServices "github.com/ginos1998/financing-market-monitor/data-processing/internal/services/redis"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

var logger logrus.Logger

func PrepareMovingAveragesData(s *server.Server) error {
	logger = *s.Logger
	c := cron.New()
	_, err := c.AddFunc("22 21 * * *", // every day at 9:10 AM
		func() {
			manageMovingAverages(&s.RedisClient, s.MongoRepository)
		})
	if err != nil {
		return err
	}
	c.Start()
	return nil
}

func CalculateMovingAverages(s *server.Server) error {
	c := cron.New()
	_, err := c.AddFunc("@every 2m",
		func() {
			calculateWMA21(&s.RedisClient)
		})
	if err != nil {
		logger.Error("Error adding cron job: ", err)
	}
	c.Start()
	return nil
}

func manageMovingAverages(redisClient *redis.RedisClient, mongoRepository mongod.MongoRepository) {
	logger.Info("MA CRON | Preparing to calculate moving averages")
	tickers, err := tickersRepository.GetTickersWithTimeSeriesData(mongoRepository)
	logger.Info("MA CRON | Tickers with time series data: ", len(tickers))
	if err != nil {
		logger.Error("Error getting tickers: ", err)
		return
	}
	err = redisServices.FlushDB(redisClient)
	if err != nil {
		logger.Error("Error flushing redis db: ", err)
		return
	}
	err = redisServices.SetSymbolsTimeSeries(redisClient, tickers)
	if err != nil {
		logger.Error("Error setting symbols time series: ", err)
		return
	}
}

func calculateWMA21(redisClient *redis.RedisClient) {
	logger.Info("MA CRON | Calculating moving averages")
	symbolPricesDailyMap, err := redisServices.GetDailySymbolsPrices(redisClient)
	if err != nil {
		logger.Error("Error getting daily symbols prices: ", err)
		return
	}

	prices := make([]float64, 0)
	for key, dailyPrices := range symbolPricesDailyMap {
		symbol := strings.Split(key, "_")[1]
		intraDaySymbolKey := fmt.Sprintf("INTRA-DAY_%s", symbol)
		symbolPricesIntraDay, err := redisRepository.GetIntraDaySymbolPrices(redisClient, intraDaySymbolKey)
		if err != nil {
			logger.Warnf("Error getting intra day symbol prices with key %s: %s", intraDaySymbolKey, err)
			continue
		}

		for _, intraDayPriceValues := range symbolPricesIntraDay {
			intraDayPrices := dtos.IntraDayPrices{}
			err := intraDayPrices.FromJSON(string(intraDayPriceValues))
			if err != nil {
				logger.Warnf("Error unmarshalling intra day prices: %s", err)
				continue
			}
			prices = append(prices, intraDayPrices.Current)
		}
		prices = append(prices, dailyPrices...)
		wma21 := indicators.WMA(prices, 21)
		logger.Infof("Symbol: %s, WMA21: %f", symbol, wma21)
	}
}
