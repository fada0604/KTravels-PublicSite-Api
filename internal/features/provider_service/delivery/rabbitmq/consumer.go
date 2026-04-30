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
	exchangeName = "ktravels.backoffice.api.exchange"
	queueName    = "provider_service_published"
	routingKey   = "ktravels.backoffice.api.routing.key"
)

type Handler func(ctx context.Context, body []byte) error

func SetupAndConsume(ctx context.Context, client *rabbitmqclient.Client, handler Handler, log *logger.Logger) error {
	if err := client.ExchangeDeclarePassive(exchangeName); err != nil {
		return fmt.Errorf("exchange %s does not exist or is unreachable: %w", exchangeName, err)
	}

	q, err := client.DeclareQueue(queueName)
	if err != nil {
		return fmt.Errorf("failed to declare queue %s: %w", queueName, err)
	}

	if err := client.BindQueue(q.Name, exchangeName, routingKey); err != nil {
		return fmt.Errorf("failed to bind queue %s: %w", queueName, err)
	}

	log.Info().
		Str("exchange", exchangeName).
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
