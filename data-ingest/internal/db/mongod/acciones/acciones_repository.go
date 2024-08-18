package acciones

import (
	"context"
	"errors"
	"time"

	"github.com/ginos1998/financing-market-monitor/data-ingest/config/server"
	"github.com/ginos1998/financing-market-monitor/data-ingest/internal/models/dtos"
)

const accionesCollectionName = "acciones"

func InsertAllAccionesArgs(server server.Server, acciones []dtos.Accion) error {
	collection := server.MongoRepository.Collections[accionesCollectionName]

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	docs := make([]interface{}, len(acciones))
	for i, accion := range acciones {
		docs[i] = accion
	}
	_, err := collection.InsertMany(ctx, docs, nil)
	if err != nil {
		return errors.New("error inserting Acciones: " + err.Error())
	}

	return nil
}
