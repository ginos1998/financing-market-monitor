package readers

import (
	"errors"
	"github.com/ginos1998/financing-market-monitor/data-ingest/config/server"
	"github.com/ginos1998/financing-market-monitor/data-ingest/internal/db/mongod/tickersRepository"
	"github.com/ginos1998/financing-market-monitor/data-ingest/internal/models/dtos"
)

const bymaTickersFileName = "resources/empresas_tickers.csv"

func ImportBYMATickersFromCsv(server server.Server) error {
	server.Logger.Info("Importing BYMA tickers data from ", bymaTickersFileName)
	requiredHeaders := []string{"company", "ticker", "has_adr", "symbol"}

	records, err := openCsvFile(bymaTickersFileName)
	if err != nil || len(records) == 0 {
		panic(err)
	}
	if !checkCsvHeaders(records[0], requiredHeaders) {
		return errors.New("error on BYMA tickers CSV: invalid csv headers")
	}

	var tickers []dtos.Ticker
	for idx, record := range records {
		if idx == 0 {
			continue
		}
		tickers = append(tickers, dtos.NewTickerFromBYMAMarket(record))
	}

	if len(tickers) == 0 {
		return errors.New("error on BYMA tickers CSV: no records found")
	}

	server.Logger.Info("BYMA tickers data read successfully. Found ", len(tickers), " records")

	err = tickersRepository.InsertTickersAll(server, tickers)
	if err != nil {
		return errors.New("error inserting BYMA tickers: " + err.Error())
	}
	server.Logger.Info("BYMA tickers inserted successfully")

	return nil
}
