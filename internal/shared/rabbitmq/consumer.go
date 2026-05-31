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
	ExchangeName = "ktravels.backoffice.api.exchange"
	DLXExchange  = "ktravels.publicsite.api.dlx"
)

type Handler func(ctx context.Context, body []byte) error

func SetupAndConsume(ctx context.Context, client *rabbitmqclient.Client, handler Handler, queueName, routingKey string, log *logger.Logger) error {
	if err := client.DeclareDirectExchange(ExchangeName); err != nil {
		return fmt.Errorf("failed to declare exchange %s: %w", ExchangeName, err)
	}

	if err := client.DeclareDirectExchange(DLXExchange); err != nil {
		return fmt.Errorf("failed to declare DLX exchange %s: %w", DLXExchange, err)
	}

	dlqName := queueName + ".dlq"
	dlq, err := client.DeclareQueue(dlqName)
	if err != nil {
		return fmt.Errorf("failed to declare DLQ %s: %w", dlqName, err)
	}
	if err := client.BindQueue(dlq.Name, DLXExchange, queueName); err != nil {
		return fmt.Errorf("failed to bind DLQ %s to DLX: %w", dlqName, err)
	}

	q, err := client.DeclareQueueWithDLX(queueName, DLXExchange)
	if err != nil {
		return fmt.Errorf("failed to declare queue %s: %w", queueName, err)
	}

	if err := client.BindQueue(q.Name, ExchangeName, routingKey); err != nil {
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
