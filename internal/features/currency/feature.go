package currency

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"

	"ktravels-publicsite-api/internal/features"
	curapp "ktravels-publicsite-api/internal/features/currency/application"
	currmq "ktravels-publicsite-api/internal/features/currency/delivery/rabbitmq"
	curinfra "ktravels-publicsite-api/internal/features/currency/infrastructure"
	sharedrabbitmq "ktravels-publicsite-api/internal/shared/rabbitmq"
)

type currencyFeature struct {
	db    *mongo.Database
	cache *curapp.CurrencyCache
}

func NewFeature(db *mongo.Database, cache *curapp.CurrencyCache) features.Feature {
	return &currencyFeature{db: db, cache: cache}
}

func (f *currencyFeature) Register(ctx context.Context, deps features.Deps) error {
	repo := curinfra.NewMongoRepository(f.db)

	f.cache.LoadAll(ctx, repo, deps.Log)

	handler := func(ctx context.Context, body []byte) error {
		return curapp.HandleConfigured(ctx, repo, f.cache, body)
	}

	if err := sharedrabbitmq.SetupAndConsume(ctx, deps.RabbitMQ, handler, currmq.ConfiguredQueueName, currmq.ConfiguredRoutingKey, deps.Log); err != nil {
		return fmt.Errorf("RabbitMQ.Consume.CurrencyConfigured: %w", err)
	}

	return nil
}
