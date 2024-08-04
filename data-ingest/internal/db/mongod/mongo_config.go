package mongod

import (
    "context"
    "time"
	"errors"

	srvCfg "github.com/ginos1998/financing-market-monitor/data-ingest/config"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "go.mongodb.org/mongo-driver/mongo/readpref"
	log "github.com/sirupsen/logrus"
)

type MongoRepository struct {
	Client *mongo.Client
	Collections map[string]*mongo.Collection
}

var db_collections = []string{
	"tickers",
}

func CreateMongoClient() (*MongoRepository, error) {
	serverURI := srvCfg.GetEnvVar("MONGO_URI")
	database := srvCfg.GetEnvVar("MONGO_DATABASE")

	if serverURI == "" || database == "" {
		return nil, errors.New("MONGO_URI or MONGO_DATABASE not set")
	}

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

	var collections map[string]*mongo.Collection = make(map[string]*mongo.Collection)
	for _, collection := range db_collections {
		collections[collection] = client.Database(database).Collection(collection)
	}

	log.Info("Connected to MongoDB, URI: ", serverURI)

	return &MongoRepository{
		Client: client,
		Collections: collections,
	}, nil
}