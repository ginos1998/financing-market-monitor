package alerts

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"

	"github.com/ginos1998/financing-market-monitor/data-processing/config/mongod"
	alertsDTO "github.com/ginos1998/financing-market-monitor/data-processing/internal/models/dtos/alerts"
)

const alertsCollectionName = "alerts"

func FindActiveAlerts(mongodb mongod.MongoRepository) ([]alertsDTO.Alert, error) {
	alertsCollection := mongodb.Collections[alertsCollectionName]

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	currentTime := time.Now().Format("2006-01-02")
	filter := bson.M{"$and": []bson.M{
		{"active": true},
		{"expiration": bson.M{"$gte": currentTime}},
	}}

	cursor, err := alertsCollection.Find(ctx, filter)
	if err != nil {
		return nil, errors.New("FindActiveAlerts error: " + err.Error())
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
		}
	}(cursor, ctx)

	var alerts []alertsDTO.Alert
	if err = cursor.All(ctx, &alerts); err != nil {
		return nil, errors.New("FindActiveAlerts error: " + err.Error())
	}
	return alerts, nil
}

func DisableAlertById(mdb mongod.MongoRepository, alertId primitive.ObjectID) error {
	alertsCollection := mdb.Collections[alertsCollectionName]

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": alertId}
	update := bson.M{"$set": bson.M{"active": false}}
	_, err := alertsCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.New("DisableAlertById error: " + err.Error())
	}
	return nil
}
