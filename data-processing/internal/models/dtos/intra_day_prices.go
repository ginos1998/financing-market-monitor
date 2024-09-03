package dtos

import "encoding/json"

type IntraDayPrices struct {
	Low     float64 `json:"low"`
	High    float64 `json:"high"`
	Open    float64 `json:"open"`
	Current float64 `json:"current"`
}

func NewIntraDayPrices(low, high, open, current float64) *IntraDayPrices {
	return &IntraDayPrices{
		Low:     low,
		High:    high,
		Open:    open,
		Current: current,
	}
}

func (i *IntraDayPrices) ToJSON() (string, error) {
	pricesJSON, err := json.Marshal(i)
	if err != nil {
		return "", err
	}
	return string(pricesJSON), nil
}

func (i *IntraDayPrices) FromJSON(pricesJSON string) error {
	err := json.Unmarshal([]byte(pricesJSON), i)
	if err != nil {
		return err
	}
	return nil
}
