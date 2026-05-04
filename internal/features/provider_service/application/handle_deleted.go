package application

import (
	"context"
	"encoding/json"
	"fmt"

	"ktravels-publicsite-api/internal/features/provider_service/domain"
	sharederrors "ktravels-publicsite-api/internal/shared/errors"
)

type deletedMessage struct {
	MessageIdentifier string `json:"MessageIdentifier"`
	Name              string `json:"Name"`
	Data              string `json:"Data"`
}

type deletedData struct {
	ProviderServiceID string `json:"provider_service_id"`
}

func HandleDeleted(ctx context.Context, repo domain.Repository, body []byte) error {
	var msg deletedMessage
	if err := json.Unmarshal(body, &msg); err != nil {
		return fmt.Errorf("%w: failed to parse message envelope: %w", sharederrors.ErrInvalidInput, err)
	}

	if msg.Data == "" {
		return fmt.Errorf("%w: empty data field", sharederrors.ErrInvalidInput)
	}

	var data deletedData
	if err := json.Unmarshal([]byte(msg.Data), &data); err != nil {
		return fmt.Errorf("%w: failed to parse delete data: %w", sharederrors.ErrInvalidInput, err)
	}

	if data.ProviderServiceID == "" {
		return fmt.Errorf("%w: provider_service_id is required", sharederrors.ErrInvalidInput)
	}

	if err := repo.Delete(ctx, data.ProviderServiceID); err != nil {
		return fmt.Errorf("failed to delete provider service %s: %w", data.ProviderServiceID, err)
	}

	return nil
}