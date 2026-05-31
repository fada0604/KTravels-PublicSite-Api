package application

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"ktravels-publicsite-api/internal/features/currency/domain"
	sharederrors "ktravels-publicsite-api/internal/shared/errors"
	sharedcurrency "ktravels-publicsite-api/internal/shared/currency"
)

type currencyConfiguredEvent struct {
	EventID     string     `json:"event_id"`
	Code        string     `json:"code"`
	Symbol      string     `json:"symbol"`
	Name        string     `json:"name"`
	Decimals    int        `json:"decimals"`
	Enabled     bool       `json:"enabled"`
	DateRestart *time.Time `json:"date_restart"`
	Timestamp   time.Time  `json:"timestamp"`
}

func HandleConfigured(ctx context.Context, repo domain.Repository, cache *CurrencyCache, body []byte) error {
	var event currencyConfiguredEvent
	if err := json.Unmarshal(body, &event); err != nil {
		return fmt.Errorf("%w: failed to parse currency.configured event: %w", sharederrors.ErrInvalidInput, err)
	}

	if event.Code == "" {
		return fmt.Errorf("%w: code is required", sharederrors.ErrInvalidInput)
	}

	currency := &domain.Currency{
		Code:        event.Code,
		Symbol:      event.Symbol,
		Name:        event.Name,
		Decimals:    event.Decimals,
		Enabled:     event.Enabled,
		DateRestart: event.DateRestart,
	}

	if err := repo.Upsert(ctx, currency); err != nil {
		return fmt.Errorf("failed to upsert currency %s: %w", event.Code, err)
	}

	cache.Set(event.Code, sharedcurrency.CurrencyInfo{
		Code:        event.Code,
		Symbol:      event.Symbol,
		Name:        event.Name,
		Decimals:    event.Decimals,
		Enabled:     event.Enabled,
		DateRestart: event.DateRestart,
	})

	return nil
}
