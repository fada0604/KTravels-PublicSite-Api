package currency

import "time"

type CurrencyInfo struct {
	Code        string
	Symbol      string
	Name        string
	Decimals    int
	Enabled     bool
	DateRestart *time.Time
}

type RateReader interface {
	Get(pair string) (float64, bool)
}

type CurrencyReader interface {
	Get(code string) (CurrencyInfo, bool)
	GetAll() []CurrencyInfo
}

type CurrencyConverter struct {
	rates      RateReader
	currencies CurrencyReader
}

func NewCurrencyConverter(rates RateReader, currencies CurrencyReader) *CurrencyConverter {
	return &CurrencyConverter{rates: rates, currencies: currencies}
}

func (c *CurrencyConverter) Convert(amount float64, fromCode, toCode string) (float64, bool) {
	rate, ok := c.rates.Get(fromCode + ":" + toCode)
	if !ok {
		return 0, false
	}
	return amount * rate, true
}

func (c *CurrencyConverter) GetEnabled() []CurrencyInfo {
	all := c.currencies.GetAll()
	enabled := make([]CurrencyInfo, 0, len(all))
	for _, ci := range all {
		if ci.Enabled {
			enabled = append(enabled, ci)
		}
	}
	return enabled
}
