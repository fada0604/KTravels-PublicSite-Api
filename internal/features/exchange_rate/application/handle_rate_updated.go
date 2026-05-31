package application

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"ktravels-publicsite-api/internal/features/exchange_rate/domain"
	sharederrors "ktravels-publicsite-api/internal/shared/errors"
)

type rateUpdatedEvent struct {
	EventID      string    `json:"event_id"`
	CurrencyFrom string    `json:"currency_from"`
	CurrencyTo   string    `json:"currency_to"`
	Rate         float64   `json:"rate"`
	Source       string    `json:"source"`
	ValidFrom    time.Time `json:"valid_from"`
	Timestamp    time.Time `json:"timestamp"`
}

func HandleRateUpdated(ctx context.Context, repo domain.Repository, cache *RateCache, body []byte) error {
	var event rateUpdatedEvent
	if err := json.Unmarshal(body, &event); err != nil {
		return fmt.Errorf("%w: failed to parse exchange_rate.updated event: %w", sharederrors.ErrInvalidInput, err)
	}

	if event.EventID == "" || event.CurrencyFrom == "" || event.CurrencyTo == "" {
		return fmt.Errorf("%w: event_id, currency_from and currency_to are required", sharederrors.ErrInvalidInput)
	}

	rate := &domain.ExchangeRate{
		EventID:      event.EventID,
		CurrencyFrom: event.CurrencyFrom,
		CurrencyTo:   event.CurrencyTo,
		Rate:         event.Rate,
		Source:       event.Source,
		ValidFrom:    event.ValidFrom,
		Timestamp:    event.Timestamp,
	}

	if err := repo.Insert(ctx, rate); err != nil {
		if errors.Is(err, sharederrors.ErrConflict) {
			return nil
		}
		return fmt.Errorf("failed to insert exchange rate event %s: %w", event.EventID, err)
	}

	cache.Set(event.CurrencyFrom+":"+event.CurrencyTo, event.Rate)
	return nil
}
