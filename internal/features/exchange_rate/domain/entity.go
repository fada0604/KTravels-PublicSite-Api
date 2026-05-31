package domain

import "time"

type ExchangeRate struct {
	ID           string
	EventID      string
	CurrencyFrom string
	CurrencyTo   string
	Rate         float64
	Source       string
	ValidFrom    time.Time
	Timestamp    time.Time
}
