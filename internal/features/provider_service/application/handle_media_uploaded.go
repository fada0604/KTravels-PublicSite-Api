package application

import (
	"context"
	"encoding/json"
	"fmt"

	"ktravels-publicsite-api/internal/features/provider_service/domain"
	sharederrors "ktravels-publicsite-api/internal/shared/errors"
)

type mediaData struct {
	EntityAssociatedID string      `json:"entity_associated_id"`
	EntityAssociated   string      `json:"entity_associated"`
	Medias             []mediaItem `json:"medias"`
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

func HandleMediaUploaded(ctx context.Context, repo domain.WriteRepository, body []byte) error {
	env, err := UnmarshalEnvelope(body)
	if err != nil {
		return err
	}

	var data mediaData
	if err := json.Unmarshal([]byte(env.Data), &data); err != nil {
		return fmt.Errorf("%w: failed to parse media data: %w", sharederrors.ErrInvalidInput, err)
	}

	if data.EntityAssociatedID == "" {
		return fmt.Errorf("%w: entity_associated_id is required", sharederrors.ErrInvalidInput)
	}

	switch domain.EntityAssociation(data.EntityAssociated) {
	case domain.EntityAssociationProviderIdentity:
		var image *domain.Image
		if len(data.Medias) > 0 {
			first := data.Medias[0]
			image = &domain.Image{Name: first.Name, URL: first.URL}
		}
		if err := repo.UpdateProviderLogo(ctx, data.EntityAssociatedID, image); err != nil {
			return fmt.Errorf("failed to update provider logo for provider_id %s: %w", data.EntityAssociatedID, err)
		}
	case domain.EntityAssociationServiceUnit:
		var images []domain.Image
		if len(data.Medias) > 0 {
			images = make([]domain.Image, len(data.Medias))
			for i, m := range data.Medias {
				images[i] = domain.Image{Name: m.Name, URL: m.URL}
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
