package mongod

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoRepository struct {
	Client      *mongo.Client
	Collections map[string]*mongo.Collection
}

var dbCollections = []string{
	"tickers",
	"cedears",
	"acciones",
}

func CreateMongoClient(envVars map[string]string) (*MongoRepository, error) {
	user := envVars["MONGO_USER"]
	password := envVars["MONGO_PASSWORD"]
	host := envVars["MONGO_HOST"]
	port := envVars["MONGO_PORT"]
	database := envVars["MONGO_DB"]
	authSource := envVars["MONGO_AUTH_SOURCE"]
	mongoDbURI := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s?authSource=%s", user, password, host, port, database, authSource)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoDbURI))

	if err != nil {
		return nil, errors.New("error connecting to MongoDB: " + err.Error())
	}

	err = client.Ping(ctx, readpref.Primary())

	if err != nil {
		return nil, errors.New("error pinging MongoDB: " + err.Error())
	}

	var collections = make(map[string]*mongo.Collection)
	for _, collection := range dbCollections {
		collections[collection] = client.Database(database).Collection(collection)
	}

	return &MongoRepository{
		Client:      client,
		Collections: collections,
	}, nil
}
