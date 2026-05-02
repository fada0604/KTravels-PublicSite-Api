package domain

import "context"

type Repository interface {
	Upsert(ctx context.Context, service *ProviderService) error
	UpdateStatus(ctx context.Context, id string, status int) error
}
