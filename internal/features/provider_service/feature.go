package providersvc

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"

	"ktravels-publicsite-api/internal/features"
	psapp "ktravels-publicsite-api/internal/features/provider_service/application"
	psrmq "ktravels-publicsite-api/internal/features/provider_service/delivery/rabbitmq"
	psinfra "ktravels-publicsite-api/internal/features/provider_service/infrastructure"
)

type providerServiceFeature struct {
	db *mongo.Database
}

// NewFeature returns the provider_service Feature wired with the given MongoDB database.
func NewFeature(db *mongo.Database) features.Feature {
	return &providerServiceFeature{db: db}
}

func (f *providerServiceFeature) Register(ctx context.Context, deps features.Deps) error {
	repo := psinfra.NewMongoRepository(f.db)

	type consumer struct {
		handler    psrmq.Handler
		queueName  string
		routingKey string
		op         string
	}

	consumers := []consumer{
		{
			handler:    func(ctx context.Context, body []byte) error { return psapp.HandleUpsert(ctx, repo, body) },
			queueName:  psrmq.PublishedQueueName,
			routingKey: psrmq.PublishedRoutingKey,
			op:         "RabbitMQ.Consume.Published",
		},
		{
			handler:    func(ctx context.Context, body []byte) error { return psapp.HandleUnpublished(ctx, repo, body) },
			queueName:  psrmq.UnpublishedQueueName,
			routingKey: psrmq.UnpublishedRoutingKey,
			op:         "RabbitMQ.Consume.Unpublished",
		},
		{
			handler:    func(ctx context.Context, body []byte) error { return psapp.HandleUpsert(ctx, repo, body) },
			queueName:  psrmq.UpdatedQueueName,
			routingKey: psrmq.UpdatedRoutingKey,
			op:         "RabbitMQ.Consume.Updated",
		},
		{
			handler:    func(ctx context.Context, body []byte) error { return psapp.HandleMediaUploaded(ctx, repo, body) },
			queueName:  psrmq.MediaUploadedQueueName,
			routingKey: psrmq.MediaUploadedRoutingKey,
			op:         "RabbitMQ.Consume.MediaUploaded",
		},
		{
			handler:    func(ctx context.Context, body []byte) error { return psapp.HandleDeleted(ctx, repo, body) },
			queueName:  psrmq.DeletedQueueName,
			routingKey: psrmq.DeletedRoutingKey,
			op:         "RabbitMQ.Consume.Deleted",
		},
	}

	for _, c := range consumers {
		if err := psrmq.SetupAndConsume(ctx, deps.RabbitMQ, c.handler, c.queueName, c.routingKey, deps.Log); err != nil {
			return fmt.Errorf("%s: %w", c.op, err)
		}
	}

	return nil
}
