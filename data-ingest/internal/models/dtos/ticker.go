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

func NewTickerFromBYMAMarket(record []string) Ticker {
	return Ticker{
		Company:      record[0],
		OriginSymbol: record[1],
		HasADR:       record[2] == "S",
		Symbol:       record[3],
		Ratio:        "1:1",
		IsCrypto:     false,
	}
}

func NewTickerFromCEDEAR(record []string) Ticker {
	return Ticker{
		Company:      record[0],
		Symbol:       record[1],
		OriginSymbol: record[1],
		Ratio:        record[2],
		HasADR:       true,
		IsCrypto:     false,
	}
}
