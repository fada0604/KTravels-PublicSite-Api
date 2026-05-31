package infrastructure

import (
	"time"

	"ktravels-publicsite-api/internal/features/currency/domain"
)

type currencyDocument struct {
	Code        string     `bson:"_id"`
	Symbol      string     `bson:"symbol"`
	Name        string     `bson:"name"`
	Decimals    int        `bson:"decimals"`
	Enabled     bool       `bson:"enabled"`
	DateRestart *time.Time `bson:"date_restart,omitempty"`
}

func toDocument(c *domain.Currency) *currencyDocument {
	return &currencyDocument{
		Code:        c.Code,
		Symbol:      c.Symbol,
		Name:        c.Name,
		Decimals:    c.Decimals,
		Enabled:     c.Enabled,
		DateRestart: c.DateRestart,
	}
}

func fromDocument(doc *currencyDocument) *domain.Currency {
	return &domain.Currency{
		Code:        doc.Code,
		Symbol:      doc.Symbol,
		Name:        doc.Name,
		Decimals:    doc.Decimals,
		Enabled:     doc.Enabled,
		DateRestart: doc.DateRestart,
	}
}
