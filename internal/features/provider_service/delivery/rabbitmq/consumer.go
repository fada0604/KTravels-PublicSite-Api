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
	DLX_EXCHANGE  = "ktravels.publicsite.api.dlx"
)

const (
	PublishedQueueName  = "provider_service_published"
	PublishedRoutingKey = "provider_service_published"

	UnpublishedQueueName  = "provider_service_unpublished"
	UnpublishedRoutingKey = "provider_service_unpublished"

	UpdatedQueueName  = "provider_service_updated"
	UpdatedRoutingKey = "provider_service_updated"

	MediaUploadedQueueName  = "media_uploaded"
	MediaUploadedRoutingKey = "media_uploaded"

	DeletedQueueName  = "provider_service_deleted"
	DeletedRoutingKey = "provider_service_deleted"
)

type Handler func(ctx context.Context, body []byte) error

func SetupAndConsume(ctx context.Context, client *rabbitmqclient.Client, handler Handler, queueName string, routingKey string, log *logger.Logger) error {
	if err := client.DeclareDirectExchange(EXCHANGE_NAME); err != nil {
		return fmt.Errorf("failed to declare exchange %s: %w", EXCHANGE_NAME, err)
	}

	if err := client.DeclareDirectExchange(DLX_EXCHANGE); err != nil {
		return fmt.Errorf("failed to declare DLX exchange %s: %w", DLX_EXCHANGE, err)
	}

	dlqName := queueName + ".dlq"
	dlq, err := client.DeclareQueue(dlqName)
	if err != nil {
		return fmt.Errorf("failed to declare DLQ %s: %w", dlqName, err)
	}
	if err := client.BindQueue(dlq.Name, DLX_EXCHANGE, queueName); err != nil {
		return fmt.Errorf("failed to bind DLQ %s to DLX: %w", dlqName, err)
	}

	q, err := client.DeclareQueueWithDLX(queueName, DLX_EXCHANGE)
	if err != nil {
		return fmt.Errorf("failed to declare queue %s: %w", queueName, err)
	}

	if err := client.BindQueue(q.Name, EXCHANGE_NAME, routingKey); err != nil {
		return fmt.Errorf("failed to bind queue %s: %w", queueName, err)
	}

	return client.Consume(ctx, queueName, func(msg amqp.Delivery) {
		msgCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := handler(msgCtx, msg.Body); err != nil {
			log.Error().
				Err(err).
				Str("messageId", msg.MessageId).
				Str("queue", queueName).
				Msg("failed to handle message, routing to DLQ")
			if nackErr := msg.Nack(false, false); nackErr != nil {
				log.Error().Err(nackErr).Msg("failed to nack message")
			}
			return
		}

		if err := msg.Ack(false); err != nil {
			log.Error().Err(err).Msg("failed to ack message")
		}
	})
}
