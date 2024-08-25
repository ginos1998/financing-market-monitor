package producers

import (
	"time"

	"github.com/ginos1998/financing-market-monitor/data-ingest/config/server"
	"github.com/ginos1998/financing-market-monitor/data-ingest/internal/apis"
	cedearsRepository "github.com/ginos1998/financing-market-monitor/data-ingest/internal/db/mongod/cedears"
	cryptoRepository "github.com/ginos1998/financing-market-monitor/data-ingest/internal/db/mongod/cryptos"
	"github.com/ginos1998/financing-market-monitor/data-ingest/internal/models/dtos"
)

func (p *KafkaProducer) UpdateHistoricalStockData(useYahooAPI bool, server server.Server) {
	test := false
	alphavantageAPI, err := apis.ConfigAlphavantageAPI(server.EnvVars, test)
	if err != nil {
		server.Logger.Fatal("Error configuring Alphavantage API: ", err)
	}

	topic := server.EnvVars["KAFKA_TOPIC_TIME_SERIES_DATA"]
	if topic == "" {
		logger.Error("envvar KAFKA_TOPIC_TIME_SERIES_DATA not set")
		return
	}

	cedears, err := cedearsRepository.GetCedearsWithoutHistoricalDailyStockData(server)
	if err != nil {
		logger.Error("Error getting cedears without historical daily stock data: ", err)
		return
	}
	logger.Info("Cedears without historical daily stock data: ", len(cedears))

	if len(cedears) == 0 {
		logger.Info("No cedears without historical daily stock data")
		return
	}

	var cedearsToUpdate []dtos.Cedear
	if useYahooAPI {
		cedearsToUpdate = cedears
	} else {
		if len(cedears) > alphavantageAPI.RequestPerDay {
			cedearsToUpdate = cedears[:alphavantageAPI.RequestPerDay]
		} else {
			cedearsToUpdate = cedears
		}
	}

	logger.Info("Updating historical stock data...")
	logger.Info("Cedears to update: ", len(cedearsToUpdate))

	var cedearsNotUpdated = make([]dtos.Cedear, 0)
	var res []byte

	for idx, cedear := range cedearsToUpdate {

		if useYahooAPI {
			res, err = apis.GetDailyHistoricalStockData(cedear.Ticker, server.EnvVars)
			if err != nil {
				logger.Error("Error getting data of ", cedear.Ticker, " from Yahoo API: ", err)
				cedearsNotUpdated = append(cedearsNotUpdated, cedear)
				continue
			}
		} else {
			logger.Info("Getting data of ", cedear.Ticker, " from Alphavantage API")
			res, err = alphavantageAPI.GetTickerDailyHistoricalData(cedear.Ticker)
			if err != nil {
				logger.Error("Error getting data of ", cedear.Ticker, " from Alphavantage API: ", err)
				cedearsNotUpdated = cedearsToUpdate[idx:]
				break
			}
		}

		logger.Info("Sending data of ", cedear.Ticker, " to Kafka...")
		err := p.produceOnTopic(topic, res)
		if err != nil {
			logger.Error("Error sending data of ", cedear.Ticker, " to Kafka: ", err)
			cedearsNotUpdated = append(cedearsNotUpdated, cedear)
		}
		logger.Info("Data of ", cedear.Ticker, " sent to Kafka successfully")
		time.Sleep(1 * time.Second)
	}

	logger.Info("The process has updated ", len(cedearsToUpdate)-len(cedearsNotUpdated), " cedears")
	if len(cedearsNotUpdated) > 0 {
		var cedearsNotUpdatedTickers []string
		for _, cedear := range cedearsNotUpdated {
			cedearsNotUpdatedTickers = append(cedearsNotUpdatedTickers, cedear.Ticker)
		}
		logger.Warn("The following cedears could not be updated: ", cedearsNotUpdatedTickers)

		if useYahooAPI {
			logger.Info("Trying again with Alphavantage API...")
			p.UpdateHistoricalStockData(false, server)
		}
	}

}

func (p *KafkaProducer) UpdateCryptosTimesSeries(server server.Server) {
	topic := server.EnvVars["KAFKA_TOPIC_CRYPTO_TIME_SERIES_DATA"]
	if topic == "" {
		logger.Error("envvar KAFKA_TOPIC_CRYPTO_TIME_SERIES_DATA not set")
		return
	}

	cryptosToUpdate, err := cryptoRepository.GetCryptos(server.MongoRepository)
	if err != nil {
		logger.Error("Error getting cryptos without historical daily stock data: ", err)
		return
	}
	logger.Info("Cryptos without historical daily stock data: ", len(cryptosToUpdate))

	if len(cryptosToUpdate) == 0 {
		logger.Info("No cedears without historical daily stock data")
		return
	}

	logger.Info("Updating historical stock data...")
	logger.Info("Cryptos to update: ", len(cryptosToUpdate))

	var cryptosNotUpdated = make([]dtos.Crypto, 0)
	var res []byte

	for _, crypto := range cryptosToUpdate {
		res, err = apis.GetDailyHistoricalStockData(crypto.YahooSymbol, server.EnvVars)
		if err != nil {
			logger.Error("Error getting data of ", crypto.YahooSymbol, " from Yahoo API: ", err)
			cryptosNotUpdated = append(cryptosNotUpdated, crypto)
			continue
		}

		logger.Info("Sending data of ", crypto.YahooSymbol, " to Kafka...")
		err := p.produceOnTopic(topic, res)
		if err != nil {
			logger.Error("Error sending data of ", crypto.YahooSymbol, " to Kafka: ", err)
			cryptosNotUpdated = append(cryptosNotUpdated, crypto)
		}
		logger.Info("Data of ", crypto.YahooSymbol, " sent to Kafka successfully")
		time.Sleep(1 * time.Second)
	}

	logger.Info("The process has updated ", len(cryptosToUpdate)-len(cryptosNotUpdated), " cryptos")
	if len(cryptosNotUpdated) > 0 {
		var cedearsNotUpdatedTickers []string
		for _, crypto := range cryptosNotUpdated {
			cedearsNotUpdatedTickers = append(cedearsNotUpdatedTickers, crypto.YahooSymbol)
		}
		logger.Warn("The following cryptos could not be updated: ", cedearsNotUpdatedTickers)
	}

}
