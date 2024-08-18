package cryptos

import (
	"context"
	"errors"
	"time"

	"github.com/ginos1998/financing-market-monitor/data-processing/config/mongod"
	"github.com/ginos1998/financing-market-monitor/data-processing/internal/models/dtos"

	"go.mongodb.org/mongo-driver/bson"
)

const cryptosCollectionName = "cryptos"

func GetCryptos(mongoRepository mongod.MongoRepository) ([]dtos.Crypto, error) {
	cryptosCollection := mongoRepository.Collections[cryptosCollectionName]

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := cryptosCollection.Find(ctx, bson.D{})
	if err != nil {
		return nil, errors.New("CryptosRepository: Error getting cryptos: " + err.Error())
	}

	var cryptos []dtos.Crypto
	err = cursor.All(ctx, &cryptos)
	if err != nil {
		return nil, errors.New("CryptosRepository: Error decoding cryptos: " + err.Error())
	}

	return cryptos, nil
}

func GetCryptoBySymbol(mongoRepository mongod.MongoRepository, symbol string) (*dtos.Crypto, error) {
	cryptosCollection := mongoRepository.Collections[cryptosCollectionName]

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.D{{Key: "symbol", Value: symbol}}
	var crypto dtos.Crypto
	err := cryptosCollection.FindOne(ctx, filter).Decode(&crypto)
	if err != nil {
		return nil, errors.New("CryptosRepository: Error getting crypto " + symbol + ": " + err.Error())
	}

	return &crypto, nil
}
