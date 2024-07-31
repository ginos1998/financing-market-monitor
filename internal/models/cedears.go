package models

type Cedear struct {
	Denom string
	Ticker string
	Ratio string
}

func NewCedear(record []string) Cedear {
	return Cedear{
		Denom: record[0],
		Ticker: record[1],
		Ratio: record[2],
	}
}
