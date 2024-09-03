package dtos

type Ticker struct {
	Symbol           string `bson:"symbol"`
	Company          string `bson:"company"`
	Ratio            string `bson:"ratio"`
	HasADR           bool   `bson:"has_adr"`
	OriginSymbol     string `bson:"origin_symbol"`
	TimeSeriesDaily  Data   `bson:"time_series_daily"`
	TimeSeriesWeekly Data   `bson:"time_series_weekly"`
	IsCrypto         bool   `bson:"is_crypto"`
}
