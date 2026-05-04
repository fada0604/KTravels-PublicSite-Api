package application_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"ktravels-publicsite-api/internal/features/provider_service/application"
	"ktravels-publicsite-api/internal/features/provider_service/domain"
	sharederrors "ktravels-publicsite-api/internal/shared/errors"
)

type mockRepo struct {
	mock.Mock
}

func (m *mockRepo) Upsert(ctx context.Context, service *domain.ProviderService) error {
	args := m.Called(ctx, service)
	return args.Error(0)
}

func (m *mockRepo) UpdateStatus(ctx context.Context, id string, status int) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *mockRepo) UpdateProviderLogo(ctx context.Context, providerID string, logo *domain.Image) error {
	args := m.Called(ctx, providerID, logo)
	return args.Error(0)
}

func (m *mockRepo) UpdateUnitImages(ctx context.Context, unitID string, images []domain.Image) error {
	args := m.Called(ctx, unitID, images)
	return args.Error(0)
}

func (m *mockRepo) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func buildMessage(data *domain.ProviderService) []byte {
	dataJSON, _ := json.Marshal(data)
	msg := map[string]interface{}{
		"MessageIdentifier": "test-msg-id",
		"Name":              "ProviderServicePublished",
		"Data":              string(dataJSON),
	}
	msgJSON, _ := json.Marshal(msg)
	return msgJSON
}

func validService() *domain.ProviderService {
	return &domain.ProviderService{
		ID:           "1475355f-8fdb-4479-8ff4-f2c61555c1aa",
		ProviderID:   "44444444-4444-4444-4444-444444444426",
		Title:        "Adrastea",
		ProviderName: "Amazonas Eco Lodge",
		Category: domain.Category{
			ID: 1, Name: "Lodging",
		},
		SubCategory: domain.SubCategory{
			ID: "a3333333-3333-3333-3333-333333333333", Name: "Apartamento",
		},
		Attributes: []domain.Attribute{
			{ID: "a1", Name: "Ventilación", Value: "Ventilador"},
		},
		Destination: domain.Destination{
			ID: "d0000001-0000-0000-0000-000000000001", Name: "Los Roques",
		},
		SubcategoriesDestination: []domain.SubcategoryDestination{
			{ID: "s1", Name: "Playa Caribe"},
		},
		LabelsDestination: []domain.LabelDestination{
			{ID: "l1", Name: "Arena Blanca"},
		},
		GeoLocation: domain.GeoLocation{
			Type: "Point", Coordinates: []float64{-65.312, 10.930},
		},
		Location:         "Test location",
		DistributionType: 1,
		Description: domain.Description{
			Description:       "<p>Test</p>",
			HourOperationFrom: "07:00:00",
			HourOperationTo:   "19:00:00",
			ServicePolicies:   "<p>Policies</p>",
		},
		Units: []domain.Unit{
			{
				ID:        "u1",
				Name:      "Test Unit",
				Capacity:  1,
				Quantity:  1,
				Amenities: []interface{}{},
				Images: []domain.Image{
					{Name: "img.png", URL: "https://example.com/img.png"},
				},
				Rates: domain.Rate{
					ID:       "r1",
					RateType: 3,
					Calendar: domain.Calendar{
						ID:        "c1",
						Name:      "Temporada Baja",
						StartDate: "2026-01-01",
						EndDate:   "2026-12-31",
						ActiveDays: []int{1, 2, 3, 4, 5},
					},
					CalendarPerPersonRate: domain.CalendarPerPersonRate{
						AgeRangeRates: []domain.AgeRangeRate{
							{ID: "ar1", MinAge: 0, MaxAge: 2, Rate: 10.0, Status: true},
						},
					},
				},
			},
		},
	}
}

func TestHandlePublished_Success(t *testing.T) {
	repo := new(mockRepo)
	svc := validService()
	body := buildMessage(svc)

	repo.On("Upsert", mock.Anything, mock.MatchedBy(func(s *domain.ProviderService) bool {
		return s.ID == svc.ID && s.Title == svc.Title
	})).Return(nil)

	err := application.HandlePublished(context.Background(), repo, body)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestHandlePublished_InvalidEnvelopeJSON(t *testing.T) {
	repo := new(mockRepo)

	err := application.HandlePublished(context.Background(), repo, []byte("not-json"))
	require.Error(t, err)
	assert.ErrorIs(t, err, sharederrors.ErrInvalidInput)
}

func TestHandlePublished_EmptyData(t *testing.T) {
	repo := new(mockRepo)
	msg := map[string]interface{}{
		"MessageIdentifier": "test",
		"Name":              "ProviderServicePublished",
		"Data":              "",
	}
	body, _ := json.Marshal(msg)

	err := application.HandlePublished(context.Background(), repo, body)
	require.Error(t, err)
	assert.ErrorIs(t, err, sharederrors.ErrInvalidInput)
}

func TestHandlePublished_InvalidDataJSON(t *testing.T) {
	repo := new(mockRepo)
	msg := map[string]interface{}{
		"MessageIdentifier": "test",
		"Name":              "ProviderServicePublished",
		"Data":              "not-json",
	}
	body, _ := json.Marshal(msg)

	err := application.HandlePublished(context.Background(), repo, body)
	require.Error(t, err)
	assert.ErrorIs(t, err, sharederrors.ErrInvalidInput)
}

func TestHandlePublished_EmptyProviderServiceID(t *testing.T) {
	repo := new(mockRepo)
	svc := validService()
	svc.ID = ""
	body := buildMessage(svc)

	err := application.HandlePublished(context.Background(), repo, body)
	require.Error(t, err)
	assert.ErrorIs(t, err, sharederrors.ErrInvalidInput)
	assert.Contains(t, err.Error(), "provider_service_id")
}

func TestHandlePublished_RepoError(t *testing.T) {
	repo := new(mockRepo)
	svc := validService()
	body := buildMessage(svc)

	repo.On("Upsert", mock.Anything, mock.Anything).Return(assert.AnError)

	err := application.HandlePublished(context.Background(), repo, body)
	require.Error(t, err)
	repo.AssertExpectations(t)
}
