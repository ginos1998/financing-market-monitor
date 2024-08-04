package mongod

import (
	"time"
	"context"
	"errors"

	dtos "github.com/ginos1998/financing-market-monitor/data-ingest/internal/models/dtos"

	"go.mongodb.org/mongo-driver/bson"
)

const tickersCollectionName = "tickers"

func (c *MongoRepository) GetCedearsWithoutHistoricalDayliStockData() ([]dtos.Cedear, error){
	tickersCollection := c.Collections[tickersCollectionName]

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"time_series_dayli": bson.M{"$exists": false}}

	cursor, err := tickersCollection.Find(ctx, filter)
	if err != nil {
		return nil, errors.New("GetCedearsWithoutHistoricalDayliStockData error: " + err.Error())
	}
	defer cursor.Close(ctx)

	var cedears []dtos.Cedear
	for cursor.Next(ctx) {
		var cedear dtos.Cedear
		err := cursor.Decode(&cedear)
		if err != nil {
			return nil, errors.New("GetCedearsWithoutHistoricalDayliStockData error: " + err.Error())
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