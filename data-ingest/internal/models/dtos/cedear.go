package dtos

type Cedear struct {
	Denom string `bson:"denom"`
	Ticker string `bson:"ticker"`
	Ratio string `bson:"ratio"`
	TimeSeriesDayli Data `bson:"time_series_dayli"`
}