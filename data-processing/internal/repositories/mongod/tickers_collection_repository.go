package mongod

import (
	"context"
	"errors"
	"time"

	"github.com/ginos1998/financing-market-monitor/data-processing/internal/models/dtos"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)




func getCollection(client *mongo.Client) (*mongo.Collection, error) {
	collection := client.Database("stock_market").Collection("tickers")
	return collection, nil
}

func InsertCedear(client *mongo.Client, cedear dtos.Cedear) error {
	collection, err := getCollection(client)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = collection.InsertOne(ctx, cedear)
	if err != nil {
		return errors.New("error inserting CEDEAR: " + err.Error())
	}

	return nil
}

func GetCedearByTicker(client *mongo.Client, ticker string) (dtos.Cedear, error) {
	collection, err := getCollection(client)
	if err != nil {
		return dtos.Cedear{}, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	filter := bson.D{{Key: "ticker", Value: ticker}}
	var cedear dtos.Cedear
	err = collection.FindOne(ctx, filter).Decode(&cedear)
	if err != nil {
		return dtos.Cedear{}, errors.New("error getting CEDEAR: " + err.Error())
	}

	return cedear, nil
}

func GetAllCedears(client *mongo.Client) ([]dtos.Cedear, error) {
	collection, err := getCollection(client)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, errors.New("error getting CEDEARs: " + err.Error())
	}
	defer cursor.Close(ctx)

	var cedears []dtos.Cedear
	for cursor.Next(ctx) {
		var cedear dtos.Cedear
		err := cursor.Decode(&cedear)
		if err != nil {
			return nil, errors.New("error decoding CEDEAR: " + err.Error())
		}
		cedears = append(cedears, cedear)
	}

	return cedears, nil
}

func UpdateCedearTimeSeriesData(client *mongo.Client, cedear dtos.Cedear) error {
	collection, err := getCollection(client)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.D{{Key: "ticker", Value: cedear.Ticker}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "time_series_dayli", Value: cedear.TimeSeriesDayli}}}}
	_, err = collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.New("error updating CEDEAR: " + err.Error())
	}

	return nil
}