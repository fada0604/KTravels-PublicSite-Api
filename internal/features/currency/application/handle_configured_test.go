package application_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"ktravels-publicsite-api/internal/features/currency/application"
	"ktravels-publicsite-api/internal/features/currency/domain"
	sharederrors "ktravels-publicsite-api/internal/shared/errors"
)

type mockRepo struct {
	mock.Mock
}

func (m *mockRepo) Upsert(ctx context.Context, currency *domain.Currency) error {
	return m.Called(ctx, currency).Error(0)
}

func (m *mockRepo) FindAll(ctx context.Context) ([]domain.Currency, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domain.Currency), args.Error(1)
}

func validCurrencyBody(code string, enabled bool) []byte {
	payload := map[string]any{
		"event_id":     "evt-cur-001",
		"code":         code,
		"symbol":       "Bs.",
		"name":         "Bolívar",
		"decimals":     2,
		"enabled":      enabled,
		"date_restart": nil,
		"timestamp":    time.Now().UTC().Format(time.RFC3339),
	}
	b, _ := json.Marshal(payload)
	return b
}

func TestHandleConfigured_NewCurrency(t *testing.T) {
	repo := new(mockRepo)
	cache := application.NewCurrencyCache()
	body := validCurrencyBody("VES", true)

	repo.On("Upsert", mock.Anything, mock.MatchedBy(func(c *domain.Currency) bool {
		return c.Code == "VES" && c.Enabled
	})).Return(nil)

	err := application.HandleConfigured(context.Background(), repo, cache, body)
	require.NoError(t, err)

	info, ok := cache.Get("VES")
	assert.True(t, ok)
	assert.Equal(t, "VES", info.Code)
	assert.Equal(t, "Bs.", info.Symbol)
	assert.True(t, info.Enabled)
	repo.AssertExpectations(t)
}

func TestHandleConfigured_ExistingCurrencyUpdated(t *testing.T) {
	repo := new(mockRepo)
	cache := application.NewCurrencyCache()

	repo.On("Upsert", mock.Anything, mock.Anything).Return(nil)

	err := application.HandleConfigured(context.Background(), repo, cache, validCurrencyBody("VES", true))
	require.NoError(t, err)

	err = application.HandleConfigured(context.Background(), repo, cache, validCurrencyBody("VES", false))
	require.NoError(t, err)

	info, ok := cache.Get("VES")
	assert.True(t, ok)
	assert.False(t, info.Enabled, "cache should reflect updated disabled state")
	repo.AssertExpectations(t)
}

func TestHandleConfigured_InvalidPayload(t *testing.T) {
	repo := new(mockRepo)
	cache := application.NewCurrencyCache()

	err := application.HandleConfigured(context.Background(), repo, cache, []byte("not-json"))
	require.Error(t, err)
	assert.ErrorIs(t, err, sharederrors.ErrInvalidInput)
}

func TestHandleConfigured_MissingCode(t *testing.T) {
	repo := new(mockRepo)
	cache := application.NewCurrencyCache()
	body, _ := json.Marshal(map[string]any{"symbol": "Bs."})

	err := application.HandleConfigured(context.Background(), repo, cache, body)
	require.Error(t, err)
	assert.ErrorIs(t, err, sharederrors.ErrInvalidInput)
}
