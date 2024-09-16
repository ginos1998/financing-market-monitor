package integrations

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"time"

	mdb "github.com/ginos1998/financing-market-monitor/data-processing/config/mongod"
	dto "github.com/ginos1998/financing-market-monitor/data-processing/internal/models/dtos/integrations"
)

const integrationsCollectionName = "integrations"

func FindIntegrationByAppName(mongodb mdb.MongoRepository, appName string) (*dto.Integration, error) {
	integrationsCollection := mongodb.Collections[integrationsCollectionName]

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"app": appName}

	var integration dto.Integration
	err := integrationsCollection.FindOne(ctx, filter).Decode(&integration)
	if err != nil {
		return nil, errors.New("FindIntegrationByAppName error: " + err.Error())
	}
	return &integration, nil
}
