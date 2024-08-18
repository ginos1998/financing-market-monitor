package dtos

type Cedear struct {
	Denom           string `bson:"denom"`
	Ticker          string `bson:"ticker"`
	Ratio           string `bson:"ratio"`
	TimeSeriesDayli Data   `bson:"time_series_dayli"`
}

func NewCedear(record []string) Cedear {
	return Cedear{
		Denom:  record[0],
		Ticker: record[1],
		Ratio:  record[2],
	}
}
