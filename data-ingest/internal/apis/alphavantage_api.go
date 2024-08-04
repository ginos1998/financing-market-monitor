package apis

import (
	"errors"
	"fmt"
	"net/http"
    "encoding/json"
    "strconv"

	appCfg "github.com/ginos1998/financing-market-monitor/data-ingest/config"
    dtos "github.com/ginos1998/financing-market-monitor/data-ingest/internal/models/dtos"
    utils "github.com/ginos1998/financing-market-monitor/data-ingest/internal/utils"

    log "github.com/sirupsen/logrus"
)

type AlphavantageAPI struct {
	URI string
	APIKey string
	RequestPerDay int
}

func (av *AlphavantageAPI) ConfigAlphavantageAPI() error {
	av.URI = appCfg.GetEnvVar("ALPHAVANTAGE_URI")
	av.APIKey = appCfg.GetEnvVar("ALPHAVANTAGE_API_KEY")
	av.RequestPerDay = 20

    if av.URI == "" || av.APIKey == "" {
        return errors.New("ALPHAVANTAGE_URI or ALPHAVANTAGE_API_KEY not set")
    }

    return nil
}

func (av *AlphavantageAPI) GetTickerDailyHistoricalData(ticker string) ([]byte, error) {
    queryParams := fmt.Sprintf("?function=TIME_SERIES_DAILY&symbol=%s&outputsize=full&apikey=%s", ticker, av.APIKey)
    url := av.URI + queryParams
    fmt.Println("URL NOT USED: ", url)

    resp, err := http.Get("https://www.alphavantage.co/query?function=TIME_SERIES_DAILY&symbol=IBM&outputsize=full&apikey=demo")
    if err != nil {
        return nil, errors.New("error getting data from Alpha Vantage: " + err.Error())
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, errors.New("error getting data from Alpha Vantage, unexpected status code: " + resp.Status)
    }

    var apiResponse dtos.APIResponse
	err = json.NewDecoder(resp.Body).Decode(&apiResponse)
	if err != nil {
		log.Error("error decoding data from Alpha Vantage: " + err.Error())
	}
	var timeSeries []dtos.TimeSeries
	for date, data := range apiResponse.TimeSeriesDaily {
		open, _ := strconv.ParseFloat(data.Open, 64)
		high, _ := strconv.ParseFloat(data.High, 64)
		low, _ := strconv.ParseFloat(data.Low, 64)
		close, _ := strconv.ParseFloat(data.Close, 64)
		volume, _ := strconv.Atoi(data.Volume)
		timeSeries = append(timeSeries, dtos.TimeSeries{
			Date:   date,
			Open:   open,
			High:   high,
			Low:    low,
			Close:  close,
			Volume: volume,
		})
	}

    orderedTimeSeries := utils.OrderTimeSeriesByDateDesc(timeSeries)

	data := dtos.Data{
		Symbol: apiResponse.MetaData.Symbol,
		LastRefreshed: apiResponse.MetaData.LastRefreshed,
		TimeZone: apiResponse.MetaData.TimeZone,
		OutputSize: apiResponse.MetaData.OutputSize,
		TimeSeriesType: "Daily",
		TimeSeriesData: orderedTimeSeries,
	}

    jsonData, err := json.Marshal(data)
    if err != nil {
        return nil, errors.New("error marshalling data: " + err.Error())
    }

    return jsonData, nil
}