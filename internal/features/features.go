package features

import (
	"context"

	"ktravels-publicsite-api/internal/platform/logger"
	rabbitmq "ktravels-publicsite-api/internal/platform/rabbitmq"
)

// Deps holds the shared runtime dependencies injected into each Feature at registration time.
type Deps struct {
	RabbitMQ *rabbitmq.Client
	Log      *logger.Logger
}

// Feature is implemented by every vertical slice that registers consumers or handlers.
type Feature interface {
	Register(ctx context.Context, deps Deps) error
}
