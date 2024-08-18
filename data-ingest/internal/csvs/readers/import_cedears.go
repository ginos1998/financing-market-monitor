package readers

import (
	"errors"
	"github.com/ginos1998/financing-market-monitor/data-ingest/config/server"
	cedearsRepository "github.com/ginos1998/financing-market-monitor/data-ingest/internal/db/mongod/cedears"
	"github.com/ginos1998/financing-market-monitor/data-ingest/internal/models/dtos"
)

const cedearsFileName = "resources/CEDEARS_17-08-2024.csv"

func ImportCedearsFromCsv(server server.Server) error {
	server.Logger.Info("Importing cedears data from ", cedearsFileName)
	requiredHeaders := []string{"denom", "ticker", "ratio"}

	records, err := openCsvFile(cedearsFileName)
	if err != nil || len(records) == 0 {
		panic(err)
	}
	if !checkCsvHeaders(records[0], requiredHeaders) {
		return errors.New("CEDEARs csv: invalid csv headers")
	}

	var cedears []dtos.Cedear

	for idx, record := range records {
		if idx == 0 {
			continue
		}
		cedears = append(cedears, dtos.NewCedear(record))
	}
	server.Logger.Info("Cedears data read successfully")

	err = cedearsRepository.InsertAllCEDEARs(server, cedears)
	if err != nil {
		return errors.New("error inserting CEDEARs: " + err.Error())
	}
	server.Logger.Info("CEDEARs inserted successfully")

	return nil
}
