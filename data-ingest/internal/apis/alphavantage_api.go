package apis

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
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
	DefaultSymbol string
}

const maxTimeSeriesData = 1500
var testing = false

func (av *AlphavantageAPI) ConfigAlphavantageAPI(test bool) error {
	if test {
		testing = true
		av.URI = "https://www.alphavantage.co/query"
		av.APIKey = "demo"
		av.RequestPerDay = 10000
		av.DefaultSymbol = "IBM"
		return nil
	}

	av.URI = appCfg.GetEnvVar("ALPHAVANTAGE_URI")
	av.APIKey = appCfg.GetEnvVar("ALPHAVANTAGE_API_KEY")
	av.RequestPerDay = 25

    if av.URI == "" || av.APIKey == "" {
        return errors.New("ALPHAVANTAGE_URI or ALPHAVANTAGE_API_KEY not set")
    }

    return nil
}

func (av *AlphavantageAPI) GetTickerDailyHistoricalData(ticker string) ([]byte, error) {
	var url string
	if testing {
		fmt.Println("URL NOT USED: ", url)
    	url = fmt.Sprintf("https://www.alphavantage.co/query?function=TIME_SERIES_DAILY&symbol=%s&outputsize=full&apikey=%s", av.DefaultSymbol, av.APIKey)
	} else {
		queryParams := fmt.Sprintf("?function=TIME_SERIES_DAILY&symbol=%s&outputsize=full&apikey=%s", ticker, av.APIKey)
    	url = av.URI + queryParams
	}

	resp, err := http.Get(url)
    if err != nil {
        return nil, errors.New("error getting data from Alpha Vantage: " + err.Error())
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, errors.New("error getting data from Alpha Vantage, unexpected status code: " + resp.Status)
    }

	data, err := decodeTickerDailyHistoricalData(resp.Body)
	if err != nil {
		return nil, errors.New("error decoding data from Alpha Vantage: " + err.Error())
	}

    jsonData, err := json.Marshal(data)
    if err != nil {
        return nil, errors.New("error marshalling data: " + err.Error())
    }

    return jsonData, nil
}

func decodeTickerDailyHistoricalData(body io.ReadCloser) (*dtos.Data, error) {
	var rateLimitResponse dtos.RateLimitResponse
    if err := json.NewDecoder(body).Decode(&rateLimitResponse); err == nil {
        if rateLimitResponse.Information != "" {
            return nil, errors.New("API rate limit reached: " + rateLimitResponse.Information)
        }
    }

    var apiResponse dtos.APIResponse
	err := json.NewDecoder(body).Decode(&apiResponse)
	if err != nil {
		log.Error("error decoding data from Alpha Vantage: " + err.Error())
		return nil, errors.New("error decoding data from Alpha Vantage: " + err.Error())
	}

	var timeSeries []dtos.TimeSeries
	for date, data := range apiResponse.TimeSeriesDaily {
		open, _ := strconv.ParseFloat(data.Open, 64)
		high, _ := strconv.ParseFloat(data.High, 64)
		low, _ := strconv.ParseFloat(data.Low, 64)
		close, _ := strconv.ParseFloat(data.Close, 64)
		volume, _ := strconv.ParseInt(data.Volume, 10, 64)
		timeSeries = append(timeSeries, dtos.TimeSeries{
			Date:   date,
			Open:   open,
			High:   high,
			Low:    low,
			Close:  close,
			Volume: volume,
		})
	}
	
	orderTruncTimeSeries := orderAndTruncateTimeSeries(timeSeries)

	return &dtos.Data{
		Symbol: apiResponse.MetaData.Symbol,
		LastRefreshed: apiResponse.MetaData.LastRefreshed,
		TimeZone: apiResponse.MetaData.TimeZone,
		OutputSize: apiResponse.MetaData.OutputSize,
		TimeSeriesType: "Daily",
		TimeSeriesData: orderTruncTimeSeries,
	}, nil
}

func orderAndTruncateTimeSeries(timeSeries []dtos.TimeSeries) []dtos.TimeSeries {
	orderedTimeSeries := utils.OrderTimeSeriesByDateDesc(timeSeries)

	if len(orderedTimeSeries) > maxTimeSeriesData {
		return orderedTimeSeries[:maxTimeSeriesData]
	} else {
		return orderedTimeSeries
	}
}