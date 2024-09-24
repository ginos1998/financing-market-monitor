package main

import (
	srv "github.com/ginos1998/financing-market-monitor/data-ingest/config/server"
	//"github.com/ginos1998/financing-market-monitor/data-ingest/internal/csvs/readers"
	"github.com/ginos1998/financing-market-monitor/data-ingest/internal/apis/nasdaq"
)

func main() {
	server := srv.NewServer()
	server.Logger.Info("Server created successfully.")

	_, err := nasdaq.FindSymbolTimeSeriesData("AAPL")
	if err != nil {
		server.Logger.Error("Error getting data from Nasdaq API: ", err)
	}

	//err := readers.ImportCedearsFromCsv(*server)
	//if err != nil {
	//	server.Logger.Error("Error importing CEDEARs: ", err)
	//}
	//
	//err = readers.ImportBYMATickersFromCsv(*server)
	//if err != nil {
	//	server.Logger.Error("Error importing acciones: ", err)
	//}

	server.Logger.Info("Process finished successfully.")
}
