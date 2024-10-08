package dtos

// MetaData Estructura para los metadatos
type MetaData struct {
	Information   string `json:"1. Information"`
	Symbol        string `json:"2. Symbol"`
	LastRefreshed string `json:"3. Last Refreshed"`
	OutputSize    string `json:"4. Output Size"`
	TimeZone      string `json:"5. Time Zone"`
}

// DailyData Estructura para los datos de tiempo diario
type DailyData struct {
	Open   string `json:"1. open"`
	High   string `json:"2. high"`
	Low    string `json:"3. low"`
	Close  string `json:"4. close"`
	Volume string `json:"5. volume"`
}

// APIResponse Estructura para toodo el JSON
type APIResponse struct {
	MetaData        MetaData             `json:"Meta Data"`
	TimeSeriesDaily map[string]DailyData `json:"Time Series (Daily)"`
}

type TimeSeries struct {
	Date     string  `bson:"date"`
	Open     float64 `bson:"open"`
	High     float64 `bson:"high"`
	Low      float64 `bson:"low"`
	Close    float64 `bson:"close"`
	AdjClose float64 `bson:"adjclose"`
	Volume   int64   `bson:"volume"`
}

type Data struct {
	Symbol         string       `bson:"symbol"`
	LastRefreshed  string       `bson:"lastrefreshed"`
	TimeZone       string       `bson:"timezone"`
	OutputSize     string       `bson:"outputsize"`
	TimeSeriesType string       `bson:"timeseriestype"`
	TimeSeriesData []TimeSeries `bson:"timeseriesdata"`
}

type RateLimitResponse struct {
	Information string `json:"Information"`
}

func (ts TimeSeries) New(date string, open, high, low, close, adjClose float64, volume int64) TimeSeries {
	return TimeSeries{
		Date:     date,
		Open:     open,
		High:     high,
		Low:      low,
		Close:    close,
		AdjClose: adjClose,
		Volume:   volume,
	}
}
