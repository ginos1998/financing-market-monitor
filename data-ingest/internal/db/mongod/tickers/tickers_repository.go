package tickers

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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

func GetTickersWithoutTimeSeries(server server.Server) ([]dtos.Ticker, error) {
	tickersCollection := server.MongoRepository.Collections[tickersCollectionName]

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{
		"$and": []bson.M{
			{"has_adr": true},
		},
		"$or": []bson.M{
			{"time_series_daily.timeseriesdata": bson.M{"$eq": nil}},
			{"time_series_weekly.timeseriesdata": bson.M{"$eq": nil}},
		},
	}

	cursor, err := tickersCollection.Find(ctx, filter)
	if err != nil {
		return nil, errors.New("GetTickersWithoutTimeSeries error: " + err.Error())
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			server.Logger.Error("Error closing cursor: ", err)
		}
	}(cursor, ctx)

	var tickers []dtos.Ticker
	if err = cursor.All(ctx, &tickers); err != nil {
		return nil, errors.New("GetTickersWithoutTimeSeries error: " + err.Error())
	}

	return tickers, nil
}
