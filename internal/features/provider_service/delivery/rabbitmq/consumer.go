package rabbitmq

import (
	"context"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"

	"ktravels-publicsite-api/internal/platform/logger"
	rabbitmqclient "ktravels-publicsite-api/internal/platform/rabbitmq"
)

const (
	EXCHANGE_NAME = "ktravels.backoffice.api.exchange"
)

const (
	PublishedQueueName   = "provider_service_published"
	PublishedRoutingKey  = "provider_service_published"

	UnpublishedQueueName = "provider_service_unpublished"
	UnpublishedRoutingKey = "provider_service_unpublished"

	UpdatedQueueName   = "provider_service_updated"
	UpdatedRoutingKey  = "provider_service_updated"

	MediaUploadedQueueName  = "media_uploaded"
	MediaUploadedRoutingKey = "media_uploaded"
)

type Handler func(ctx context.Context, body []byte) error

func SetupAndConsume(ctx context.Context, client *rabbitmqclient.Client, handler Handler, queueName string, routingKey string, log *logger.Logger) error {
	if err := client.DeclareDirectExchange(EXCHANGE_NAME); err != nil {
		return fmt.Errorf("failed to declare exchange %s: %w", EXCHANGE_NAME, err)
	}

	q, err := client.DeclareQueue(queueName)
	if err != nil {
		return fmt.Errorf("failed to declare queue %s: %w", queueName, err)
	}

	if err := client.BindQueue(q.Name, EXCHANGE_NAME, routingKey); err != nil {
		return fmt.Errorf("failed to bind queue %s: %w", queueName, err)
	}

	log.Info().
		Str("exchange", EXCHANGE_NAME).
		Str("queue", queueName).
		Str("routingKey", routingKey).
		Msg("RabbitMQ consumer setup complete, starting consume")

	return client.Consume(ctx, queueName, func(msg amqp.Delivery) {
		msgCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := handler(msgCtx, msg.Body); err != nil {
			log.Error().
				Err(err).
				Str("messageId", msg.MessageId).
				Str("queue", queueName).
				Msg("failed to handle provider service message")
			if nackErr := msg.Nack(false, true); nackErr != nil {
				log.Error().Err(nackErr).Msg("failed to nack message")
			}
			return
		}

		if err := msg.Ack(false); err != nil {
			log.Error().Err(err).Msg("failed to ack message")
		}
	})
}
