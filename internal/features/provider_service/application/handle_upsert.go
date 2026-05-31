package application

import (
	"context"
	"encoding/json"
	"fmt"

	"ktravels-publicsite-api/internal/features/provider_service/domain"
	sharederrors "ktravels-publicsite-api/internal/shared/errors"
)

func HandleUpsert(ctx context.Context, repo domain.WriteRepository, body []byte) error {
	env, err := UnmarshalEnvelope(body)
	if err != nil {
		return err
	}

	var service domain.ProviderService
	if err := json.Unmarshal([]byte(env.Data), &service); err != nil {
		return fmt.Errorf("%w: failed to parse provider service: %w", sharederrors.ErrInvalidInput, err)
	}

	if service.ID == "" {
		return fmt.Errorf("%w: provider_service_id is required", sharederrors.ErrInvalidInput)
	}

	if err := repo.Upsert(ctx, &service); err != nil {
		return fmt.Errorf("failed to upsert provider service %s: %w", service.ID, err)
	}

	return nil
}
