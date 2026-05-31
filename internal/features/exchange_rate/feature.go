package exchangerate

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"

	"ktravels-publicsite-api/internal/features"
	erapp "ktravels-publicsite-api/internal/features/exchange_rate/application"
	errmq "ktravels-publicsite-api/internal/features/exchange_rate/delivery/rabbitmq"
	erinfra "ktravels-publicsite-api/internal/features/exchange_rate/infrastructure"
	sharedrabbitmq "ktravels-publicsite-api/internal/shared/rabbitmq"
)

type exchangeRateFeature struct {
	db    *mongo.Database
	cache *erapp.RateCache
}

func NewFeature(db *mongo.Database, cache *erapp.RateCache) features.Feature {
	return &exchangeRateFeature{db: db, cache: cache}
}

func (f *exchangeRateFeature) Register(ctx context.Context, deps features.Deps) error {
	repo := erinfra.NewMongoRepository(f.db)

	f.cache.LoadAll(ctx, repo, deps.Log)

	handler := func(ctx context.Context, body []byte) error {
		return erapp.HandleRateUpdated(ctx, repo, f.cache, body)
	}

	if err := sharedrabbitmq.SetupAndConsume(ctx, deps.RabbitMQ, handler, errmq.RateUpdatedQueueName, errmq.RateUpdatedRoutingKey, deps.Log); err != nil {
		return fmt.Errorf("RabbitMQ.Consume.RateUpdated: %w", err)
	}

	return nil
}
