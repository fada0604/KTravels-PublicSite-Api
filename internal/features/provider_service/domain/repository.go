package domain

import "context"

type WriteRepository interface {
	Upsert(ctx context.Context, service *ProviderService) error
	UpdateStatus(ctx context.Context, id string, status int) error
	UpdateProviderLogo(ctx context.Context, providerID string, logo *Image) error
	UpdateUnitImages(ctx context.Context, unitID string, images []Image) error
	Delete(ctx context.Context, id string) error
}

type ReadRepository interface {
	FindByID(ctx context.Context, id string) (*ProviderService, error)
}

type Repository interface {
	WriteRepository
	ReadRepository
}
