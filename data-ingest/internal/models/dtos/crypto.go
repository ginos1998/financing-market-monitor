package dtos

type Crypto struct {
	Symbol           string `json:"symbol"`
	Description      string `json:"description"`
	DisplaySymbol    string `json:"displaySymbol"`
	RedisKey         string `json:"redisKey"`
	TimeSeriesDaily  Data   `json:"timeSeriesDaily"`
	TimeSeriesWeekly Data   `json:"timeSeriesWeekly"`
	YahooSymbol      string `json:"yahooSymbol"`
}
