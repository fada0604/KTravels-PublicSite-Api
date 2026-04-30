package domain

import "context"

type Repository interface {
	Upsert(ctx context.Context, service *ProviderService) error
}
