package main

import (


	"github.com/ginos1998/financing-market-monitor/internal/import_data"
	"github.com/ginos1998/financing-market-monitor/config/db"
)

func main() {
	mongoClient, err := db.GetMongoClient()
	if err != nil {
		panic(err)
	}

	import_data.ImportCedears(mongoClient)
}