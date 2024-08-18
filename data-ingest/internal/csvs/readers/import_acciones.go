package readers

import (
	"errors"
	"github.com/ginos1998/financing-market-monitor/data-ingest/config/server"
	accionesRepository "github.com/ginos1998/financing-market-monitor/data-ingest/internal/db/mongod/acciones"
	"github.com/ginos1998/financing-market-monitor/data-ingest/internal/models/dtos"
)

const accionesFileName = "resources/empresas_tickers.csv"

func ImportAccionesFromCsv(server server.Server) error {
	server.Logger.Info("Importing acciones data from ", accionesFileName)
	requiredHeaders := []string{"company", "ticker", "has_adr", "symbol"}

	records, err := openCsvFile(accionesFileName)
	if err != nil || len(records) == 0 {
		panic(err)
	}
	if !checkCsvHeaders(records[0], requiredHeaders) {
		return errors.New("acciones CSV: invalid csv headers")
	}

	var acciones []dtos.Accion

	for idx, record := range records {
		if idx == 0 {
			continue
		}
		acciones = append(acciones, dtos.NewAccion(record))
	}

	if len(acciones) == 0 {
		return errors.New("acciones CSV: no records found")
	}

	server.Logger.Info("Acciones data read successfully. Found ", len(acciones), " records")

	err = accionesRepository.InsertAllAccionesArgs(server, acciones)
	if err != nil {
		return errors.New("error inserting Acciones: " + err.Error())
	}
	server.Logger.Info("Acciones inserted successfully")

	return nil
}
