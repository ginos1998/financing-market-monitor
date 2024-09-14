package nasdaq

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ginos1998/financing-market-monitor/data-ingest/internal/models/dtos"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

// APIResponse NasdaqAPIResponse Root struct que representa la respuesta completa del JSON
type APIResponse struct {
	Data    Data   `json:"data"`
	Message string `json:"message"`
	Status  Status `json:"status"`
}

// Data Struct que representa la sección "data"
type Data struct {
	Symbol       string      `json:"symbol"`
	TotalRecords int         `json:"totalRecords"`
	TradesTable  TradesTable `json:"tradesTable"`
}

// TradesTable Struct que representa la tabla de operaciones "tradesTable"
type TradesTable struct {
	AsOf    interface{} `json:"asOf"` // Puede ser null, así que lo dejamos como tipo genérico (interface{})
	Headers Headers     `json:"headers"`
	Rows    []TradeRow  `json:"rows"`
}

// Headers Struct que representa los encabezados de las columnas
type Headers struct {
	Date   string `json:"date"`
	Close  string `json:"close"`
	Volume string `json:"volume"`
	Open   string `json:"open"`
	High   string `json:"high"`
	Low    string `json:"low"`
}

// TradeRow Struct que representa cada fila de datos de operaciones (trades)
type TradeRow struct {
	Date   string `json:"date"`
	Close  string `json:"close"`
	Volume string `json:"volume"`
	Open   string `json:"open"`
	High   string `json:"high"`
	Low    string `json:"low"`
}

// Status Struct que representa el estado del API
type Status struct {
	RCode            int    `json:"rCode"`
	BCodeMessage     string `json:"bCodeMessage"`
	DeveloperMessage string `json:"developerMessage"`
}

const fromDate = "2000-01-01"
const limit = 9999

func FindSymbolTimeSeriesData(stockSymbol string) ([]byte, error) {
	currentDate := time.Now().Format("2006-01-02")
	n := rand.Intn(100)
	url := fmt.Sprintf("https://api.nasdaq.com/api/quote/%s/historical?assetclass=stocks&fromdate=%s&limit=%d&todate=%s&random=%d", stockSymbol, fromDate, limit, currentDate, n)

	log.Println("Getting historical stock data from Nasdaq API for ", stockSymbol)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.New("Error creando la solicitud HTTP GET: " + err.Error())
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36")

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.New("Error al realizar la solicitud HTTP GET a Nasdaq API: " + err.Error())
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Error cerrando el cuerpo de la respuesta HTTP: %v", err)
		}
	}(response.Body)

	var nasdaqResponse APIResponse
	log.Printf("Leyendo respuesta de Nasdaq API")
	err = json.NewDecoder(response.Body).Decode(&nasdaqResponse)
	if err != nil {
		return nil, errors.New("Error leyendo la respuesta de Nasdaq API: " + err.Error())
	}

	var timeSeries []dtos.TimeSeries
	for _, record := range nasdaqResponse.Data.TradesTable.Rows {
		timeSeries = append(timeSeries, dtos.TimeSeries{
			Date:   record.Date,
			Open:   priceStringToFloat64(record.Open),
			High:   priceStringToFloat64(record.High),
			Low:    priceStringToFloat64(record.Low),
			Close:  priceStringToFloat64(record.Close),
			Volume: strToInt(record.Volume),
		})
	}

	data := dtos.Data{
		Symbol:         stockSymbol,
		LastRefreshed:  currentDate,
		TimeZone:       "UTC-3",
		OutputSize:     "Full",
		TimeSeriesType: "1d",
		TimeSeriesData: timeSeries,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, errors.New("error marshalling data: " + err.Error())
	}

	log.Printf("%s found %d records", stockSymbol, len(timeSeries))
	return jsonData, nil
}

func priceStringToFloat64(price string) float64 {
	if price[0] == '$' {
		price = price[1:]
	}
	priceFloat, _ := strconv.ParseFloat(price, 64)
	return priceFloat
}

func strToInt(str string) int64 {
	if str == "" {
		return 0
	}
	num, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0
	}
	return num
}
