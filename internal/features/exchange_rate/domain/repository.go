package domain

import "context"

type Repository interface {
	Insert(ctx context.Context, rate *ExchangeRate) error
	ExistsByEventID(ctx context.Context, eventID string) (bool, error)
	FindLatestByPair(ctx context.Context) (map[string]float64, error)
	FindPaginated(ctx context.Context, page, pageSize int) ([]ExchangeRate, int64, error)
}
