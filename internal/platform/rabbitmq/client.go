package rabbitmq

import (
	"context"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Client struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	config  Config
}

type Config struct {
	HostName string
	UserName string
	Password string
	Port    int
}

func New(cfg Config) (*Client, error) {
	url := fmt.Sprintf("amqp://%s:%s@%s:%d/",
		cfg.UserName,
		cfg.Password,
		cfg.HostName,
		cfg.Port,
	)

	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	return &Client{
		conn:    conn,
		channel: ch,
		config:  cfg,
	}, nil
}

func (c *Client) Consume(ctx context.Context, queue string, handler func(amqp.Delivery)) error {
	msgs, err := c.channel.Consume(
		queue,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to register consumer: %w", err)
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-msgs:
				if !ok {
					return
				}
				handler(msg)
			}
		}
	}()

	return nil
}

func (c *Client) Listen(ctx context.Context, queue string, handler func(amqp.Delivery)) error {
	msgs, err := c.channel.Consume(
		queue,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to register consumer: %w", err)
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case msg, ok := <-msgs:
			if !ok {
				return nil
			}
			handler(msg)
		}
	}
}

func (c *Client) Publish(ctx context.Context, exchange, routingKey string, body []byte) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return c.channel.PublishWithContext(
		ctx,
		exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

func (c *Client) DeclareTopicExchange(name string) error {
	return c.channel.ExchangeDeclare(
		name,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
}

// ExchangeDeclarePassive verifies that an exchange already exists without
// attempting to create or modify it. Useful when the exchange is owned by
// another service (e.g. Backoffice) and this service only needs to bind to it.
func (c *Client) ExchangeDeclarePassive(name string) error {
	return c.channel.ExchangeDeclarePassive(
		name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
}

func (c *Client) DeclareQueue(name string) (amqp.Queue, error) {
	return c.channel.QueueDeclare(
		name,
		true,
		false,
		false,
		false,
		nil,
	)
}

func (c *Client) BindQueue(queue, exchange, routingKey string) error {
	return c.channel.QueueBind(
		queue,
		routingKey,
		exchange,
		false,
		nil,
	)
}

func (c *Client) Close() error {
	if c.channel != nil {
		c.channel.Close()
	}
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

func (c *Client) IsConnected() bool {
	return c.conn != nil && !c.conn.IsClosed()
}