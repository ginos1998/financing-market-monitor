package apis

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"errors"
	"sort"
	"time"

	"github.com/ginos1998/financing-market-monitor/data-processing/internal/models/dtos"
)

type ByDateDesc []dtos.TimeSeries

func (a ByDateDesc) Len() int           { return len(a) }
func (a ByDateDesc) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByDateDesc) Less(i, j int) bool {
    dateI, _ := time.Parse("2006-01-02", a[i].Date)
    dateJ, _ := time.Parse("2006-01-02", a[j].Date)
    return dateI.After(dateJ) // Comparación descendente
}

const alphavantage_api_key = "XS1B4GCIFGSPQBNF"
const alphavantage_url = "https://www.alphavantage.co/query?function=%s&apikey=%s"
const TIME_SERIES_DAILY = "TIME_SERIES_DAILY"
const TIME_SERIES_WEEKLY = "TIME_SERIES_WEEKLY"

func GetTickerDailyHistoricalData(ticker string) (dtos.Data, error) {
    dayliParams := fmt.Sprintf("&symbol=%s&outputsize=full", ticker)
    apiUrl := fmt.Sprintf(alphavantage_url, TIME_SERIES_DAILY, alphavantage_api_key)
    dayliUrl := apiUrl + dayliParams
    resp, err := http.Get(dayliUrl)
    if err != nil {
        return dtos.Data{}, errors.New("error getting data from Alpha Vantage: " + err.Error())
    }
    defer resp.Body.Close()

    var apiResponse dtos.APIResponse
    err = json.NewDecoder(resp.Body).Decode(&apiResponse)
    if err != nil {
        return dtos.Data{}, errors.New("error decoding data from Alpha Vantage: " + err.Error())
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

	// ordet timeSeries by date
	sort.Sort(ByDateDesc(timeSeries))

	timeSeriesLimited := timeSeries[:500]

	data:= dtos.Data{
		LastRefreshed: apiResponse.MetaData.LastRefreshed,
		TimeZone: apiResponse.MetaData.TimeZone,
		OutputSize: apiResponse.MetaData.OutputSize,
		TimeSeriesType: "Daily",
		TimeSeriesData: timeSeriesLimited,
	}

    return data, nil
}

func GetTickerDayliHistoricalData(ticker string) {
	url := fmt.Sprintf(alphavantage_url, ticker, alphavantage_api_key)
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		panic(err)
	}
	
	fmt.Println(data)
}