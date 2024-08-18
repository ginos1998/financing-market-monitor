package cedears

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/mongo"
	"time"

	"github.com/ginos1998/financing-market-monitor/data-ingest/config/server"
	"github.com/ginos1998/financing-market-monitor/data-ingest/internal/models/dtos"

	"go.mongodb.org/mongo-driver/bson"
)

const cedearsCollectionName = "cedears"

func InsertAllCEDEARs(server server.Server, cedears []dtos.Cedear) error {
	collection := server.MongoRepository.Collections[cedearsCollectionName]
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	docs := make([]interface{}, len(cedears))
	for i, cedear := range cedears {
		docs[i] = cedear
	}
	_, err := collection.InsertMany(ctx, docs, nil)
	if err != nil {
		return errors.New("error inserting CEDEARs: " + err.Error())
	}

	return nil
}

func InsertCedear(server server.Server, cedear dtos.Cedear) error {
	collection := server.MongoRepository.Collections[cedearsCollectionName]

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, cedear)
	if err != nil {
		return errors.New("error inserting CEDEAR: " + err.Error())
	}

	return nil
}

func GetCedearsWithoutHistoricalDailyStockData(server server.Server) ([]dtos.Cedear, error) {
	tickersCollection := server.MongoRepository.Collections[cedearsCollectionName]

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"time_series_dayli": bson.M{"$exists": false}}

	cursor, err := tickersCollection.Find(ctx, filter)
	if err != nil {
		return nil, errors.New("GetCedearsWithoutHistoricalDailyStockData error: " + err.Error())
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			server.Logger.Error("error closing cursor: ", err)
		}
	}(cursor, ctx)

	var cedears []dtos.Cedear
	for cursor.Next(ctx) {
		var cedear dtos.Cedear
		err := cursor.Decode(&cedear)
		if err != nil {
			return nil, errors.New("GetCedearsWithoutHistoricalDailyStockData error: " + err.Error())
		}
		cedears = append(cedears, cedear)
	}

	return cedears, nil

}

// func UpdateCedearTimeSeriesData(client *mongo.Client, cedear models.Cedear) error {
// 	collection, err := getCollection(client)
// 	if err != nil {
// 		return err
// 	}

// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	filter := bson.D{{Key: "ticker", Value: cedear.Ticker}}
// 	update := bson.D{{Key: "$set", Value: bson.D{{Key: "time_series_dayli", Value: cedear.TimeSeriesDayli}}}}
// 	_, err = collection.UpdateOne(ctx, filter, update)
// 	if err != nil {
// 		return errors.New("error updating CEDEAR: " + err.Error())
// 	}

// 	return nil
// }
