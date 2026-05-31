package graph

// THIS CODE WILL BE UPDATED WITH SCHEMA CHANGES. PREVIOUS IMPLEMENTATION FOR SCHEMA CHANGES WILL BE KEPT IN THE COMMENT SECTION. IMPLEMENTATION FOR UNCHANGED SCHEMA WILL BE KEPT.

import (
	"context"

	erdomain "ktravels-publicsite-api/internal/features/exchange_rate/domain"
)

type Resolver struct {
	ExchangeRateRepo erdomain.Repository
}

// Noop is the resolver for the noop field.
func (r *mutationResolver) Noop(ctx context.Context) (bool, error) {
	return false, nil
}

// Health is the resolver for the health field.
func (r *queryResolver) Health(ctx context.Context) (string, error) {
	return "ok", nil
}

// ExchangeRates is the resolver for the exchangeRates field.
func (r *queryResolver) ExchangeRates(ctx context.Context, page *int, pageSize *int) (*ExchangeRatePage, error) {
	p, ps := 1, 20
	if page != nil {
		p = *page
	}
	if pageSize != nil {
		ps = *pageSize
	}

	rates, total, err := r.ExchangeRateRepo.FindPaginated(ctx, p, ps)
	if err != nil {
		return nil, err
	}

	items := make([]*ExchangeRate, len(rates))
	for i, rate := range rates {
		items[i] = &ExchangeRate{
			ID:           rate.ID,
			CurrencyFrom: rate.CurrencyFrom,
			CurrencyTo:   rate.CurrencyTo,
			Rate:         rate.Rate,
			Source:       rate.Source,
			ValidFrom:    rate.ValidFrom,
			Timestamp:    rate.Timestamp,
		}
	}

	return &ExchangeRatePage{
		Items:       items,
		TotalCount:  int(total),
		HasNextPage: int64(p*ps) < total,
	}, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
