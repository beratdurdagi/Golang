package types

type ClientResponse struct {
	Ticker string  `json:"ticker"`
	Price  float64 `json:"price"`
}
