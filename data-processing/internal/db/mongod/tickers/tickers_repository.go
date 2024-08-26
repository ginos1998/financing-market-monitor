package tickers

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"

	"github.com/ginos1998/financing-market-monitor/data-processing/config/mongod"
	"github.com/ginos1998/financing-market-monitor/data-processing/config/server"
	"github.com/ginos1998/financing-market-monitor/data-processing/internal/models/dtos"
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

	filter := bson.M{"$and": []bson.M{
		{"time_series_daily.timeseriesdata": bson.M{"$eq": nil}},
		{"has_adr": true},
	}}

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

func UpdateTickerTimeSeriesDaily(mongoRepository mongod.MongoRepository, ticker dtos.Ticker) error {
	tickersCollection := mongoRepository.Collections[tickersCollectionName]

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.D{{Key: "symbol", Value: ticker.Symbol}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "time_series_daily", Value: ticker.TimeSeriesDaily}}}}
	tickersCollection.FindOneAndUpdate(ctx, filter, update)

	return nil
}

func UpdateTickerTimeSeriesWeekly(mongoRepository mongod.MongoRepository, ticker dtos.Ticker) error {
	tickersCollection := mongoRepository.Collections[tickersCollectionName]

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.D{{Key: "symbol", Value: ticker.Symbol}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "time_series_weekly", Value: ticker.TimeSeriesWeekly}}}}
	tickersCollection.FindOneAndUpdate(ctx, filter, update)

	return nil
}

func GetTickersWithTimeSeriesData(repository mongod.MongoRepository) ([]dtos.Ticker, error) {
	tickersCollection := repository.Collections[tickersCollectionName]

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	timeSeriesLimit := 200

	pipeline := mongo.Pipeline{
		{{"$match", bson.D{
			{"$and", bson.A{
				bson.D{{"is_crypto", false}},
				bson.D{{"time_series_daily.timeseriesdata", bson.D{{"$ne", nil}}}},
			}},
		}}},
		{{"$project", bson.D{
			{"_id", 0},
			{"symbol", 1},
			{"time_series_daily.symbol", 1},
			{"time_series_daily.lastrefreshed", 1},
			{"time_series_daily.timeseriestype", 1},
			{"time_series_daily.timeseriesdata", bson.D{
				{"$slice", bson.A{
					bson.D{{"$reverseArray", "$time_series_daily.timeseriesdata"}},
					timeSeriesLimit,
				}},
			}},
			{"time_series_weekly.symbol", 1},
			{"time_series_weekly.lastrefreshed", 1},
			{"time_series_weekly.timeseriestype", 1},
			{"time_series_weekly.timeseriesdata", bson.D{
				{"$slice", bson.A{
					bson.D{{"$reverseArray", "$time_series_weekly.timeseriesdata"}},
					timeSeriesLimit,
				}},
			}},
		}}},
	}
	res, err := tickersCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, errors.New("GetTickersWithTimeSeriesData error: " + err.Error())
	}

	var tickers []dtos.Ticker
	if err = res.All(ctx, &tickers); err != nil {
		return nil, errors.New("GetTickersWithTimeSeriesData error: " + err.Error())
	}

	return tickers, nil
}
