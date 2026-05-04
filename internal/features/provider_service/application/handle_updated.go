package application

import (
	"context"
	"encoding/json"
	"fmt"

	"ktravels-publicsite-api/internal/features/provider_service/domain"
	sharederrors "ktravels-publicsite-api/internal/shared/errors"
)

func HandleUpdated(ctx context.Context, repo domain.Repository, body []byte) error {
	var msg publishedMessage
	if err := json.Unmarshal(body, &msg); err != nil {
		return fmt.Errorf("%w: failed to parse message envelope: %w", sharederrors.ErrInvalidInput, err)
	}

	if msg.Data == "" {
		return fmt.Errorf("%w: empty data field", sharederrors.ErrInvalidInput)
	}

	var service domain.ProviderService
	if err := json.Unmarshal([]byte(msg.Data), &service); err != nil {
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