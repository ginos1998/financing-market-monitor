package producers

import (
	"time"

	"github.com/ginos1998/financing-market-monitor/data-ingest/config/server"
	"github.com/ginos1998/financing-market-monitor/data-ingest/internal/apis/nasdaq"
	tickersRepository "github.com/ginos1998/financing-market-monitor/data-ingest/internal/db/mongod/tickers"
	"github.com/ginos1998/financing-market-monitor/data-ingest/internal/models/dtos"
)

var periods = []string{"1d"} // nasdaqApi only provides daily time series, "1wk"}

func (p *KafkaProducer) UpdateTickersTimeSeries(server server.Server) {
	topic := server.EnvVars["KAFKA_TOPIC_TIME_SERIES_DATA"]
	if topic == "" {
		logger.Error("envvar KAFKA_TOPIC_TIME_SERIES_DATA not set")
		return
	}
	tickersToUpdate, err := tickersRepository.GetTickersWithoutTimeSeries(server)
	if err != nil {
		logger.Error("Error getting tickers without time series: ", err)
		return
	}
	logger.Info("tickers without time series: ", len(tickersToUpdate))

	if len(tickersToUpdate) == 0 {
		logger.Info("No tickers without time series")
		return
	}

	logger.Info("Updating time series data...")
	logger.Info("Tickers to update: ", len(tickersToUpdate))

	var tickersNotUpdated = make([]dtos.Ticker, 0)
	var res []byte

	for _, ticker := range tickersToUpdate {
		for _, period := range periods {
			if period == "1d" && ticker.TimeSeriesDaily.TimeSeriesData != nil {
				continue
			}
			if period == "1wk" && ticker.TimeSeriesWeekly.TimeSeriesData != nil {
				continue
			}
			res, err = nasdaq.FindSymbolTimeSeriesData(ticker.Symbol, ticker.AssetClass, server.EnvVars)
			if err != nil {
				logger.Error("Error getting ", period, " data of ", ticker.Symbol, " from Yahoo API: ", err)
				tickersNotUpdated = append(tickersNotUpdated, ticker)
				continue
			}

			logger.Info("Sending ", period, " data of ", ticker.Symbol, " to Kafka...")
			err := p.produceOnTopic(topic, res)
			if err != nil {
				logger.Error("Error sending ", period, " data of ", ticker.Symbol, " to Kafka: ", err)
				tickersNotUpdated = append(tickersNotUpdated, ticker)
			}
			logger.Info("Data of ", ticker.Symbol, " sent to Kafka successfully")
			time.Sleep(1 * time.Second)
		}
	}

	logger.Info("The process has updated ", len(tickersToUpdate)-len(tickersNotUpdated), " tickers")
	if len(tickersNotUpdated) > 0 {
		var tickersNotUpdatedSymbols []string
		for _, ticker := range tickersNotUpdated {
			tickersNotUpdatedSymbols = append(tickersNotUpdatedSymbols, ticker.Symbol)
		}
		logger.Warn("The following cedears could not be updated: ", tickersNotUpdatedSymbols)
	}

}
