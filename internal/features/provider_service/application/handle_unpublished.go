package application

import (
	"context"
	"encoding/json"
	"fmt"

	"ktravels-publicsite-api/internal/features/provider_service/domain"
	sharederrors "ktravels-publicsite-api/internal/shared/errors"
)

type unpublishedData struct {
	ProviderServiceID string `json:"provider_service_id"`
	Status            int    `json:"status"`
}

func HandleUnpublished(ctx context.Context, repo domain.Repository, body []byte) error {
	var msg publishedMessage
	if err := json.Unmarshal(body, &msg); err != nil {
		return fmt.Errorf("%w: failed to parse message envelope: %w", sharederrors.ErrInvalidInput, err)
	}

	if msg.Data == "" {
		return fmt.Errorf("%w: empty data field", sharederrors.ErrInvalidInput)
	}

	var data unpublishedData
	if err := json.Unmarshal([]byte(msg.Data), &data); err != nil {
		return fmt.Errorf("%w: failed to parse unpublished data: %w", sharederrors.ErrInvalidInput, err)
	}

	if data.ProviderServiceID == "" {
		return fmt.Errorf("%w: provider_service_id is required", sharederrors.ErrInvalidInput)
	}

	if err := repo.UpdateStatus(ctx, data.ProviderServiceID, data.Status); err != nil {
		return fmt.Errorf("failed to update status for provider service %s: %w", data.ProviderServiceID, err)
	}

	return nil
}