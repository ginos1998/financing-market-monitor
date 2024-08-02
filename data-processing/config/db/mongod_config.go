package db

import (
    "context"
    "time"
	"errors"
	"fmt"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "go.mongodb.org/mongo-driver/mongo/readpref"
)

const serverURI = "mongodb://localhost:27017"

func GetMongoClient() (*mongo.Client, error) {
	fmt.Println("Connecting to MongoDB at ", serverURI)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(serverURI))
	
	if err != nil {
		return nil, errors.New("error connecting to MongoDB: " + err.Error())
	}
	
	err = client.Ping(ctx, readpref.Primary())
	
	if err != nil {
		return nil, errors.New("error pinging MongoDB: " + err.Error())
	}

	fmt.Println("Connected to MongoDB")
	
	return client, nil
}