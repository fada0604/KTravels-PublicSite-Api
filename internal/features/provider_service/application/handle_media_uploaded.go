package application

import (
	"context"
	"encoding/json"
	"fmt"

	"ktravels-publicsite-api/internal/features/provider_service/domain"
	sharederrors "ktravels-publicsite-api/internal/shared/errors"
)

const (
	EntityAssociatedProviderDigitalIdentity  = "ProviderDigitalIdentity"
	EntityAssociatedServiceDistributionUnit   = "ServiceDistributionUnit"
)

type mediaUploadedMessage struct {
	MessageIdentifier string `json:"MessageIdentifier"`
	Name              string `json:"Name"`
	Data              string `json:"Data"`
}

type mediaData struct {
	EntityAssociatedID string       `json:"entity_associated_id"`
	EntityAssociated   string       `json:"entity_associated"`
	Medias             []mediaItem  `json:"medias"`
}

type mediaItem struct {
	ID                 string `json:"id"`
	Name               string `json:"name"`
	UniqueName         string `json:"unique_name"`
	URL                string `json:"url"`
	Type               string `json:"type"`
	EntityAssociatedID string `json:"entity_associated_id"`
	EntityAssociated   string `json:"entity_associated"`
}

func HandleMediaUploaded(ctx context.Context, repo domain.Repository, body []byte) error {
	var msg mediaUploadedMessage
	if err := json.Unmarshal(body, &msg); err != nil {
		return fmt.Errorf("%w: failed to parse message envelope: %w", sharederrors.ErrInvalidInput, err)
	}

	if msg.Data == "" {
		return fmt.Errorf("%w: empty data field", sharederrors.ErrInvalidInput)
	}

	var data mediaData
	if err := json.Unmarshal([]byte(msg.Data), &data); err != nil {
		return fmt.Errorf("%w: failed to parse media data: %w", sharederrors.ErrInvalidInput, err)
	}

	if data.EntityAssociatedID == "" {
		return fmt.Errorf("%w: entity_associated_id is required", sharederrors.ErrInvalidInput)
	}

	switch data.EntityAssociated {
	case EntityAssociatedProviderDigitalIdentity:
		var image *domain.Image
		if len(data.Medias) > 0 {
			firstMedia := data.Medias[0]
			image = &domain.Image{
				Name: firstMedia.Name,
				URL:  firstMedia.URL,
			}
		}
		if err := repo.UpdateProviderLogo(ctx, data.EntityAssociatedID, image); err != nil {
			return fmt.Errorf("failed to update provider logo for provider_id %s: %w", data.EntityAssociatedID, err)
		}
	case EntityAssociatedServiceDistributionUnit:
		var images []domain.Image
		if len(data.Medias) > 0 {
			images = make([]domain.Image, len(data.Medias))
			for i, m := range data.Medias {
				images[i] = domain.Image{
					Name: m.Name,
					URL:  m.URL,
				}
			}
		}
		if err := repo.UpdateUnitImages(ctx, data.EntityAssociatedID, images); err != nil {
			return fmt.Errorf("failed to update unit images for unit_id %s: %w", data.EntityAssociatedID, err)
		}
	default:
		return fmt.Errorf("%w: unsupported entity_associated type: %s", sharederrors.ErrInvalidInput, data.EntityAssociated)
	}

	return nil
}