package dtos

type FinnhubWsDTO struct {
	LastPrice float64 `json:"p"`
	Symbol   string  `json:"s"`
	Timestamp int64   `json:"t"`
	Volume  float64   `json:"v"`
}

type WsData struct {
	Trades []FinnhubWsDTO `json:"data"`
}


