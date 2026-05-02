package application

import (
	"context"
	"encoding/json"
	"fmt"

	"ktravels-publicsite-api/internal/features/provider_service/domain"
	sharederrors "ktravels-publicsite-api/internal/shared/errors"

	"github.com/rs/zerolog/log"
)

func HandleUpdated(ctx context.Context, repo domain.Repository, body []byte) error {
	log.Info().Str("handler", "HandleUpdated").Bytes("body", body).Msg("Received updated message")

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

	log.Info().Str("provider_service_id", service.ID).Int("status", service.Status).Msg("Parsed updated data")

	if service.ID == "" {
		return fmt.Errorf("%w: provider_service_id is required", sharederrors.ErrInvalidInput)
	}

	if err := repo.Upsert(ctx, &service); err != nil {
		return fmt.Errorf("failed to upsert provider service %s: %w", service.ID, err)
	}

	log.Info().Str("provider_service_id", service.ID).Msg("Successfully upserted provider service")

	return nil
}