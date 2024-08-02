package apis

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/ginos1998/financing-market-monitor/data-ingest/config"
)

var alphavantage_api_key string
const alphavantage_url = "https://www.alphavantage.co/query?function=TIME_SERIES_DAILY&symbol=%s&outputsize=full&apikey=%s"

func GetTickerDailyHistoricalData(ticker string) (http.Response, error) {
    initApiKeyFromEnv()
    url := fmt.Sprintf(alphavantage_url, ticker, alphavantage_api_key)
    resp, err := http.Get(url)
    if err != nil {
        return http.Response{}, errors.New("error getting data from Alpha Vantage: " + err.Error())
    }
    
    return *resp, nil
    // if err != nil {
    //     return dtos.Data{}, errors.New("error getting data from Alpha Vantage: " + err.Error())
    // }
    // defer resp.Body.Close()

    // var apiResponse dtos.APIResponse
    // err = json.NewDecoder(resp.Body).Decode(&apiResponse)
    // if err != nil {
    //     return dtos.Data{}, errors.New("error decoding data from Alpha Vantage: " + err.Error())
    // }

    // var timeSeries []dtos.TimeSeries
    // for date, data := range apiResponse.TimeSeriesDaily {
    //     open, _ := strconv.ParseFloat(data.Open, 64)
    //     high, _ := strconv.ParseFloat(data.High, 64)
    //     low, _ := strconv.ParseFloat(data.Low, 64)
    //     close, _ := strconv.ParseFloat(data.Close, 64)
    //     volume, _ := strconv.Atoi(data.Volume)

    //     timeSeries = append(timeSeries, dtos.TimeSeries{
    //         Date:   date,
    //         Open:   open,
    //         High:   high,
    //         Low:    low,
    //         Close:  close,
    //         Volume: volume,
    //     })
    // }

	// data:= dtos.Data{
	// 	LastRefreshed: apiResponse.MetaData.LastRefreshed,
	// 	TimeZone: apiResponse.MetaData.TimeZone,
	// 	OutputSize: apiResponse.MetaData.OutputSize,
	// 	TimeSeriesType: "Daily",
	// 	TimeSeriesData: timeSeries,
	// }

    // return data, nil
}

func GetTickerDayliHistoricalData(ticker string) {
    initApiKeyFromEnv()
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

func initApiKeyFromEnv() {
    if alphavantage_api_key != "" {
        return
    }
    alphavantage_api_key = config.GetEnvVar("ALPHAVANTAGE_API_KEY")
}