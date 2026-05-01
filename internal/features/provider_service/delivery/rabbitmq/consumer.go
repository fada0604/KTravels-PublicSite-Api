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
	QUEUE_NAME    = "provider_service_published"
	ROUTING_KEY   = "ktravels.backoffice.api.routing.key"
)

type Handler func(ctx context.Context, body []byte) error

func SetupAndConsume(ctx context.Context, client *rabbitmqclient.Client, handler Handler, log *logger.Logger) error {
	if err := client.DeclareDirectExchange(EXCHANGE_NAME); err != nil {
		return fmt.Errorf("failed to declare exchange %s: %w", EXCHANGE_NAME, err)
	}

	q, err := client.DeclareQueue(QUEUE_NAME)
	if err != nil {
		return fmt.Errorf("failed to declare queue %s: %w", QUEUE_NAME, err)
	}

	if err := client.BindQueue(q.Name, EXCHANGE_NAME, ROUTING_KEY); err != nil {
		return fmt.Errorf("failed to bind queue %s: %w", QUEUE_NAME, err)
	}

	log.Info().
		Str("exchange", EXCHANGE_NAME).
		Str("queue", QUEUE_NAME).
		Str("routingKey", ROUTING_KEY).
		Msg("RabbitMQ consumer setup complete, starting consume")

	return client.Consume(ctx, QUEUE_NAME, func(msg amqp.Delivery) {
		msgCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := handler(msgCtx, msg.Body); err != nil {
			log.Error().
				Err(err).
				Str("messageId", msg.MessageId).
				Msg("failed to handle provider service published message")
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
