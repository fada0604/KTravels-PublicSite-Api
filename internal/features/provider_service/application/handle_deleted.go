package application

import (
	"context"
	"encoding/json"
	"fmt"

	"ktravels-publicsite-api/internal/features/provider_service/domain"
	sharederrors "ktravels-publicsite-api/internal/shared/errors"
)

type deletedData struct {
	ProviderServiceID string `json:"provider_service_id"`
}

func HandleDeleted(ctx context.Context, repo domain.WriteRepository, body []byte) error {
	env, err := UnmarshalEnvelope(body)
	if err != nil {
		return err
	}

	var data deletedData
	if err := json.Unmarshal([]byte(env.Data), &data); err != nil {
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
