package dtos

type TimeSeries struct {
	Date   string  `bson:"date"`
	Open   float64 `bson:"open"`
	High   float64 `bson:"high"`
	Low    float64 `bson:"low"`
	Close  float64 `bson:"close"`
	Volume int     `bson:"volume"`
}