package import_data

import (
	"fmt"
	"os"
	"errors"
	"encoding/csv"

	"github.com/ginos1998/financing-market-monitor/internal/models"
	"github.com/ginos1998/financing-market-monitor/internal/repositories/mongod"
	"github.com/ginos1998/financing-market-monitor/internal/import_data/apis"
	
	"go.mongodb.org/mongo-driver/mongo"
)

const filePath = "/home/ginos/Descargas/cedears_julio_2024.csv"

func UpdateCedearTimeSeriesData(mongoClient *mongo.Client, ticker string) {
	cedear, err := mongod.GetCedearByTicker(mongoClient, ticker)
	if err != nil {
		fmt.Println("error getting cedear: ", err)
		return
	}

	data, err := apis.GetTickerDailyHistoricalData(cedear.Ticker)
	if err != nil || len(data.TimeSeriesData) == 0 {
		fmt.Println("error getting cedear data: ", err)
		return
	}

	cedear.TimeSeriesDayli = data
	err = mongod.UpdateCedearTimeSeriesData(mongoClient, cedear)
	if err != nil {
		fmt.Println("error updating cedear data: ", err)
		return
	}

	fmt.Println("Cedear data updated successfully")
}

func ImportCedears(mongoClient *mongo.Client) {
	fmt.Println("Importing cedears data from ", filePath)

	records, err := openCsvFile()
	if err != nil || len(records) == 0 {
		panic(err)
	}

	var cedears []models.Cedear

	for idx, record := range records {
		if idx == 0 {
			continue
		}
		cedear := models.NewCedear(record)
		cedears = append(cedears, cedear)
	}

	fmt.Println("Cedears data imported successfully")


	for _, cedear := range cedears {
		err := mongod.InsertCedear(mongoClient, cedear)
		if err != nil {
			fmt.Println("error inserting cedear: ", err)
			return
		}
	}

	fmt.Println("Cedears data inserted successfully")
}

func openCsvFile() ([][]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, errors.New("error opening file: " + err.Error())
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, errors.New("error reading file: " + err.Error())
	}

	if !checkCsvHeaders(records[0]) {
		return nil, errors.New("invalid CSV headers")
	}

	return records, nil
}

func checkCsvHeaders(headers []string) bool {
	if len(headers) != 3 {
		return false
	}
	if headers[0] != "denom" || headers[1] != "ticker" || headers[2] != "ratio" {
		return false
	}
	return true
}
