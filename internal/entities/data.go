package entities

type PriceData struct {
	Gold           []Item `json:"gold"`
	Currency       []Item `json:"currency"`
	Cryptocurrency []Item `json:"cryptocurrency"`
}

type Item struct {
	Date          string  `json:"date"`
	Time          string  `json:"time"`
	Symbol        string  `json:"symbol"`
	Name          string  `json:"name"`
	Price         float64 `json:"price"`
	ChangePercent float64 `json:"change_percent"`
	Unit          string  `json:"unit"`
}
