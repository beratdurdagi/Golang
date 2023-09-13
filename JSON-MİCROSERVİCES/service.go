package main

import (
	"context"
	"fmt"
)

type PriceFetcher interface {
	FetchPrice(context.Context, string) (float64, error)
}

type priceFetcher struct{}

func NewPriceFetcher() *priceFetcher {
	return &priceFetcher{}
}

func (p *priceFetcher) FetchPrice(ctx context.Context, ticker string) (float64, error) {
	return MockPriceFether(ctx, ticker)

}

var priceMocks = map[string]float64{
	"BTC":  25_564.12,
	"ETH":  200.0,
	"HAQQ": 100_000.0,
}

func MockPriceFether(ctx context.Context, ticker string) (float64, error) {

	price, ok := priceMocks[ticker]

	if !ok {
		return price, fmt.Errorf("Ticker (%s) is not supported", ticker)

	}
	return price, nil

}
