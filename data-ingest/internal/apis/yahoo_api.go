package apis

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/ginos1998/financing-market-monitor/data-ingest/internal/models/dtos"
)

const periodFrom = 946684800 // 01/01/2000
const yahooFinanceURL = "/download/%s?period1=%d&period2=%d&interval=%s&events=history&includeAdjustedClose=true"

func FindSymbolTimeSeriesData(stockSymbol string, period string, envvars map[string]string) ([]byte, error) {
	yahooURL := envvars["YAHOO_FINANCE_URL"]
	if yahooURL == "" {
		return nil, errors.New("variable YAHOO_FINANCE_URL not set")
	}
	p := "1d"
	if period != "" {
		p = period
	}
	periodTo := time.Now().Unix()
	url := fmt.Sprintf(yahooURL+yahooFinanceURL, stockSymbol, periodFrom, periodTo, p)

	logger.Info("Getting historical stock data from Yahoo Finance API for ", stockSymbol)

	response, err := http.Get(url)
	if err != nil {
		return nil, errors.New("Error al realizar la solicitud HTTP GET a yahooFinanceURL: " + err.Error())
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.Error("Error cerrando el cuerpo de la respuesta HTTP: ", err)
		}
	}(response.Body)

	reader := csv.NewReader(response.Body)

	records, err := reader.ReadAll()
	if err != nil {
		return nil, errors.New("Error leyendo la respuesta de yahooFinanceURL: " + err.Error())
	}

	var stockData []dtos.TimeSeries

	for i, record := range records {
		if i == 0 {
			continue
		}

		open, _ := strconv.ParseFloat(record[1], 64)
		high, _ := strconv.ParseFloat(record[2], 64)
		low, _ := strconv.ParseFloat(record[3], 64)
		close, _ := strconv.ParseFloat(record[4], 64)
		adjClose, _ := strconv.ParseFloat(record[5], 64)
		volume, _ := strconv.ParseInt(record[6], 10, 64)

		data := dtos.TimeSeries{
			Date:     record[0],
			Open:     open,
			High:     high,
			Low:      low,
			Close:    close,
			AdjClose: adjClose,
			Volume:   volume,
		}

		stockData = append(stockData, data)
	}

	var lastRefreshed = ""
	if len(stockData) > 0 {
		lastRefreshed = stockData[len(stockData)-1].Date
	} else {
		logger.Warn("No se encontraron datos para el símbolo ", stockSymbol)
		return nil, errors.New("No se encontraron datos para el símbolo " + stockSymbol)

	}

	data := dtos.Data{
		Symbol:         stockSymbol,
		LastRefreshed:  lastRefreshed,
		TimeZone:       "UTC",
		OutputSize:     "Full",
		TimeSeriesType: p,
		TimeSeriesData: stockData,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, errors.New("error marshalling data: " + err.Error())
	}

	logger.Info("Historical stock data from Yahoo Finance API for ", stockSymbol, " retrieved successfully. Found ", len(stockData), " records")
	return jsonData, nil
}
