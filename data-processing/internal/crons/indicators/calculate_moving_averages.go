package indicators

import (
	"github.com/ginos1998/financing-market-monitor/data-processing/config/mongod"
	"github.com/ginos1998/financing-market-monitor/data-processing/config/redis"
	"github.com/ginos1998/financing-market-monitor/data-processing/config/server"
	cryptosRepository "github.com/ginos1998/financing-market-monitor/data-processing/internal/db/mongod/cryptos"
	redisRepository "github.com/ginos1998/financing-market-monitor/data-processing/internal/db/redis/intra_day"
	"github.com/ginos1998/financing-market-monitor/data-processing/internal/indicators"
	"github.com/ginos1998/financing-market-monitor/data-processing/internal/models/dtos"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

var logger logrus.Logger

func CalculateMovingAverages(s *server.Server) error {
	logger = *s.Logger
	c := cron.New()
	_, err := c.AddFunc("@every 1m", func() {
		calculateWMA21(&s.RedisClient, s.MongoRepository)
	})
	if err != nil {
		return err
	}
	c.Start()
	return nil
}

func calculateWMA21(redisClient *redis.RedisClient, mongoRepository mongod.MongoRepository) {
	cryptos, err := cryptosRepository.GetCryptos(mongoRepository)
	if err != nil {
		logger.Error("Error getting cryptos: ", err)
		return
	}
	for _, crypto := range cryptos {
		closeValues := make([]float64, 0)
		for _, ts := range crypto.TimeSeriesDaily.TimeSeriesData[:20] {
			closeValues = append(closeValues, ts.Close)
		}
		symbolPrices, err := redisRepository.GetIntraDaySymbolPrices(redisClient, crypto.Symbol)
		symbolValues := dtos.IntraDayPrices{}
		err = symbolValues.FromJSON(symbolPrices)
		if err != nil {
			return
		}
		closeValues = append(closeValues, symbolValues.Current)
		wma := indicators.WMA(closeValues, len(closeValues))
		logger.Infof("Symbol: %s, WMA: %f", crypto.Symbol, wma)
	}
}
