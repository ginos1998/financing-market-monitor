package cedears

import (
	"context"
	"time"

	"github.com/ginos1998/financing-market-monitor/data-processing/config/mongod"
	dtos "github.com/ginos1998/financing-market-monitor/data-processing/internal/models/dtos"
	"go.mongodb.org/mongo-driver/bson"
)

const cedearsCollectionName = "cedears"

func UpdateCedearTimeSeriesData(mongoRepository mongod.MongoRepository, cedear dtos.Cedear) error {
	tickersCollection := mongoRepository.Collections[cedearsCollectionName]

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.D{{Key: "ticker", Value: cedear.Ticker}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "time_series_dayli", Value: cedear.TimeSeriesDayli}}}}
	tickersCollection.FindOneAndUpdate(ctx, filter, update)

	return nil
}
