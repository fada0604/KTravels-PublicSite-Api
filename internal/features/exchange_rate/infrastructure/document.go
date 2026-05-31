package infrastructure

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"ktravels-publicsite-api/internal/features/exchange_rate/domain"
)

type exchangeRateDocument struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	EventID      string             `bson:"event_id"`
	CurrencyFrom string             `bson:"currency_from"`
	CurrencyTo   string             `bson:"currency_to"`
	Rate         float64            `bson:"rate"`
	Source       string             `bson:"source"`
	ValidFrom    time.Time          `bson:"valid_from"`
	Timestamp    time.Time          `bson:"timestamp"`
}

func toDocument(r *domain.ExchangeRate) *exchangeRateDocument {
	return &exchangeRateDocument{
		EventID:      r.EventID,
		CurrencyFrom: r.CurrencyFrom,
		CurrencyTo:   r.CurrencyTo,
		Rate:         r.Rate,
		Source:       r.Source,
		ValidFrom:    r.ValidFrom,
		Timestamp:    r.Timestamp,
	}
}

func fromDocument(doc *exchangeRateDocument) *domain.ExchangeRate {
	return &domain.ExchangeRate{
		ID:           doc.ID.Hex(),
		EventID:      doc.EventID,
		CurrencyFrom: doc.CurrencyFrom,
		CurrencyTo:   doc.CurrencyTo,
		Rate:         doc.Rate,
		Source:       doc.Source,
		ValidFrom:    doc.ValidFrom,
		Timestamp:    doc.Timestamp,
	}
}
