package application_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"ktravels-publicsite-api/internal/features/exchange_rate/application"
	"ktravels-publicsite-api/internal/features/exchange_rate/domain"
	sharederrors "ktravels-publicsite-api/internal/shared/errors"
)

type mockRepo struct {
	mock.Mock
}

func (m *mockRepo) Insert(ctx context.Context, rate *domain.ExchangeRate) error {
	return m.Called(ctx, rate).Error(0)
}

func (m *mockRepo) ExistsByEventID(ctx context.Context, eventID string) (bool, error) {
	args := m.Called(ctx, eventID)
	return args.Bool(0), args.Error(1)
}

func (m *mockRepo) FindLatestByPair(ctx context.Context) (map[string]float64, error) {
	args := m.Called(ctx)
	return args.Get(0).(map[string]float64), args.Error(1)
}

func (m *mockRepo) FindPaginated(ctx context.Context, page, pageSize int) ([]domain.ExchangeRate, int64, error) {
	args := m.Called(ctx, page, pageSize)
	return args.Get(0).([]domain.ExchangeRate), args.Get(1).(int64), args.Error(2)
}

func validEventBody(eventID string) []byte {
	payload := map[string]any{
		"event_id":      eventID,
		"currency_from": "USD",
		"currency_to":   "VES",
		"rate":          54.32,
		"source":        "manual",
		"valid_from":    time.Now().UTC().Format(time.RFC3339),
		"timestamp":     time.Now().UTC().Format(time.RFC3339),
	}
	b, _ := json.Marshal(payload)
	return b
}

func TestHandleRateUpdated_NewEvent(t *testing.T) {
	repo := new(mockRepo)
	cache := application.NewRateCache()
	body := validEventBody("evt-001")

	repo.On("Insert", mock.Anything, mock.MatchedBy(func(r *domain.ExchangeRate) bool {
		return r.EventID == "evt-001" && r.CurrencyFrom == "USD" && r.CurrencyTo == "VES"
	})).Return(nil)

	err := application.HandleRateUpdated(context.Background(), repo, cache, body)
	require.NoError(t, err)

	rate, ok := cache.Get("USD:VES")
	assert.True(t, ok)
	assert.Equal(t, 54.32, rate)
	repo.AssertExpectations(t)
}

func TestHandleRateUpdated_DuplicateEvent(t *testing.T) {
	repo := new(mockRepo)
	cache := application.NewRateCache()
	body := validEventBody("evt-dup")

	repo.On("Insert", mock.Anything, mock.Anything).
		Return(sharederrors.ErrConflict)

	err := application.HandleRateUpdated(context.Background(), repo, cache, body)
	require.NoError(t, err)

	_, ok := cache.Get("USD:VES")
	assert.False(t, ok, "cache should not be updated for duplicate event")
	repo.AssertExpectations(t)
}

func TestHandleRateUpdated_InvalidPayload(t *testing.T) {
	repo := new(mockRepo)
	cache := application.NewRateCache()

	err := application.HandleRateUpdated(context.Background(), repo, cache, []byte("not-json"))
	require.Error(t, err)
	assert.ErrorIs(t, err, sharederrors.ErrInvalidInput)
}

func TestHandleRateUpdated_MissingRequiredFields(t *testing.T) {
	repo := new(mockRepo)
	cache := application.NewRateCache()
	body, _ := json.Marshal(map[string]any{"rate": 10.0})

	err := application.HandleRateUpdated(context.Background(), repo, cache, body)
	require.Error(t, err)
	assert.ErrorIs(t, err, sharederrors.ErrInvalidInput)
}
