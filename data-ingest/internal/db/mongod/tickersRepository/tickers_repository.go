package tickersRepository

import (
	"context"
	"errors"
	"time"

	"github.com/ginos1998/financing-market-monitor/data-ingest/config/server"
	"github.com/ginos1998/financing-market-monitor/data-ingest/internal/models/dtos"
)

const tickersCollectionName = "tickers"

func InsertTickersAll(server server.Server, tickers []dtos.Ticker) error {
	collection := server.MongoRepository.Collections[tickersCollectionName]
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	docs := make([]interface{}, len(tickers))
	for i, ticker := range tickers {
		docs[i] = ticker
	}
	res, err := collection.InsertMany(ctx, docs, nil)
	if err != nil {
		return errors.New("error inserting tickers: " + err.Error())
	}
	server.Logger.Info("Inserted ", len(res.InsertedIDs), " tickers successfully")

	return nil
}
