package main

import (
	srv "github.com/ginos1998/financing-market-monitor/data-ingest/config/server"
	"github.com/ginos1998/financing-market-monitor/data-ingest/internal/csvs/readers"
)

func main() {
	server := srv.NewServer()
	server.Logger.Info("Server created successfully.")

	err := readers.ImportCedearsFromCsv(*server)
	if err != nil {
		server.Logger.Error("Error importing CEDEARs: ", err)
	}

	err = readers.ImportBYMATickersFromCsv(*server)
	if err != nil {
		server.Logger.Error("Error importing acciones: ", err)
	}

	server.Logger.Info("Process finished successfully.")
}
