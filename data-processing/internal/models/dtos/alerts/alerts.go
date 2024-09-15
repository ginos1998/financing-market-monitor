package alerts

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Alert struct {
	Id         primitive.ObjectID `bson:"_id"`
	Symbol     string             `bson:"symbol"`
	Price      float64            `bson:"price"`
	Type       string             `bson:"type"`
	Trigger    string             `bson:"trigger"`
	Name       string             `bson:"name"`
	Message    string             `bson:"message"`
	Expiration string             `bson:"expiration"`
	Active     bool               `bson:"active"`
}

func (a *Alert) ToString() string {
	return fmt.Sprintf("Alert: {Id: %s, Symbol: %s, Price: %f, Type: %s, Name: %s, Message: %s}",
		a.Id.Hex(), a.Symbol, a.Price, a.Type, a.Name, a.Message)
}
