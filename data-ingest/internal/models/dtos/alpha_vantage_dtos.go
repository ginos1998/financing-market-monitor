package dtos

// Estructura para los metadatos
type MetaData struct {
    Information   string `json:"1. Information"`
    Symbol        string `json:"2. Symbol"`
    LastRefreshed string `json:"3. Last Refreshed"`
    OutputSize    string `json:"4. Output Size"`
    TimeZone      string `json:"5. Time Zone"`
}

// Estructura para los datos de tiempo diario
type DailyData struct {
    Open   string `json:"1. open"`
    High   string `json:"2. high"`
    Low    string `json:"3. low"`
    Close  string `json:"4. close"`
    Volume string `json:"5. volume"`
}

// Estructura para todo el JSON
type APIResponse struct {
    MetaData          MetaData              `json:"Meta Data"`
    TimeSeriesDaily   map[string]DailyData  `json:"Time Series (Daily)"`
}

type TimeSeries struct {
	Date   string  `bson:"date"`
	Open   float64 `bson:"open"`
	High   float64 `bson:"high"`
	Low    float64 `bson:"low"`
	Close  float64 `bson:"close"`
	Volume int     `bson:"volume"`
}

type Data struct {
	LastRefreshed   string      `bson:"lastrefreshed"`
	TimeZone        string      `bson:"timezone"`
	OutputSize      string      `bson:"outputsize"`
	TimeSeriesType  string      `bson:"timeseriestype"`
	TimeSeriesData  []TimeSeries `bson:"timeseriesdata"`
}