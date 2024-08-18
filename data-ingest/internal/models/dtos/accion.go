package dtos

type Accion struct {
	Company string `bson:"empresa"`
	Ticker  string `bson:"ticker"`
	HasADR  bool   `bson:"has_adr"`
	Symbol  string `bson:"symbol"`
}

func NewAccion(record []string) Accion {
	return Accion{
		Company: record[0],
		Ticker:  record[1],
		HasADR:  record[2] == "S",
		Symbol:  record[3],
	}
}
