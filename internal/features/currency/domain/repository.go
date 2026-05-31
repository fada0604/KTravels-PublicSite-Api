package domain

import "context"

type Repository interface {
	Upsert(ctx context.Context, currency *Currency) error
	FindAll(ctx context.Context) ([]Currency, error)
}
