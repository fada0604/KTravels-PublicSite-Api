package domain

import "time"

type Currency struct {
	Code        string
	Symbol      string
	Name        string
	Decimals    int
	Enabled     bool
	DateRestart *time.Time
}
