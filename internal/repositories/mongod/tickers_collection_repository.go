package mongod

import (
	"context"
	"time"
	"errors"

	"github.com/ginos1998/financing-market-monitor/internal/models"

	"go.mongodb.org/mongo-driver/mongo"
)




func getCollection(client *mongo.Client) (*mongo.Collection, error) {
	collection := client.Database("stock_market").Collection("tickers")
	return collection, nil
}

func InsertCedear(client *mongo.Client, cedear models.Cedear) error {
	collection, err := getCollection(client)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = collection.InsertOne(ctx, cedear)
	if err != nil {
		return errors.New("Error inserting CEDEAR: " + err.Error())
	}

	return nil
}