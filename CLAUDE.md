# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# Start all services (API + MongoDB + RabbitMQ)
docker-compose up -d

# Local build
go build -o api ./app/main.go

# Run all tests
go test ./...

# Run a single test
go test ./internal/features/provider_service/application/... -run TestHandlePublished

# Regenerate GraphQL code after schema changes
go run github.com/99designs/gqlgen generate
```

> **Important:** Run `gqlgen generate` before building whenever `graph/schema.graphqls` or any feature-level `schema.graphqls` is modified. `graph/generated.go` and `graph/model.go` are git-ignored and must be regenerated locally.

## Architecture: Vertical Slice

Each feature lives entirely within `internal/features/{feature}/` and owns all its layers:

```
internal/features/{feature}/
├── domain/         # entity.go, repository.go (interfaces)
├── application/    # use cases (handle_published.go, etc.)
├── infrastructure/ # mongo_repository.go, document.go
└── delivery/
    ├── graphql/    # resolver.go, mapper.go, schema.graphqls
    └── rabbitmq/   # consumer.go
```

**Non-negotiable rules:**
1. GraphQL resolvers must be thin — business logic belongs in `application/` or `domain/`.
2. All MongoDB access goes through repository interfaces, never directly from resolvers.
3. No horizontal global folders (`controllers/`, `services/`, `repositories/`).
4. Shared code goes in `internal/shared/` or `internal/platform/`; don't move there prematurely — wait for real reuse across slices.

## Platform Layer

`internal/platform/` wires the core infrastructure used by all features:
- `config/` — Viper-based config; uses `APP_` prefix for nested keys (e.g., `app.port` → `APP_PORT`). **Use this as the source of truth for env var names, not `.env.example`** (which is outdated and uses incorrect naming).
- `database/` — MongoDB client wrapper
- `logger/` — zerolog JSON logger
- `rabbitmq/` — AMQP consumer/publisher client
- `graphql/` — gqlgen HTTP handler stub

## Environment Variables

Canonical env var names (from `internal/platform/config/config.go:50-103`):

| Env Var | Default | Purpose |
|---------|---------|---------|
| `APP_NAME` | `ktravels-publicsite-api` | Service name |
| `APP_ENV` | `development` | Environment |
| `APP_PORT` | `8080` | HTTP port |
| `MONGODB_URI` | `mongodb://localhost:27017` | MongoDB connection |
| `MONGODB_DATABASE` | `ktravels_publicsite` | Database name |
| `RABBITMQ_HOSTNAME` | `localhost` | RabbitMQ host |
| `RABBITMQ_USERNAME` / `RABBITMQ_PASSWORD` | `guest` / `guest` | RabbitMQ credentials |
| `RABBITMQ_PORT` | `5672` | RabbitMQ AMQP port |
| `GRAPHQL_PLAYGROUND_ENABLED` | `true` | Enable playground |
| `LOGGER_LEVEL` | `info` | Log level |
| `LOGGER_FORMAT` | `json` | `json` or `console` |

## GraphQL (gqlgen, schema-first)

Schema lives at `graph/schema.graphqls` (base) and optionally at `internal/features/{feature}/delivery/graphql/schema.graphqls` per feature. The `.gqlgen.yml` config merges them. After editing any schema file, regenerate with `go run github.com/99designs/gqlgen generate`.

GraphQL naming: types/queries use `PascalCase`, fields use `camelCase`, mutation inputs use `Input` suffix, mutation response types use `Payload` suffix.

## RabbitMQ Integration (provider_service feature)

Exchange: `ktravels.backoffice.api.exchange` (direct, non-durable). The `provider_service` feature consumes these routing keys:
- `provider_service_published` → upsert to MongoDB
- `provider_service_unpublished` / `provider_service_updated` → status update
- `media_uploaded` → update provider logo or unit images
- `provider_service_deleted` → delete document

MongoDB collection: `provider_services`. Infrastructure documents are separate from domain entities (see `document.go` vs `entity.go`).

## Naming & Style

- Go files: `snake_case`; packages: `lowercase` single word
- Private symbols: `camelCase`; public types/functions: `PascalCase`; constants: `UPPER_SNAKE_CASE`
- Code and inline comments: **English**; documentation (`.md`, commits, PRs): **Spanish**
- MongoDB collections: `snake_case` plural
- Use sentinel errors from `internal/shared/errors/` (`ErrNotFound`, `ErrInvalidInput`, etc.); never `panic()` for business errors; never ignore errors with `_`

## Git & Branches

Gitflow: `main` (production) ← `develop` (integration) ← `feature/{name}`.

Commit format: `<type>(<scope>): <description>` with description in Spanish.
Types: `feat`, `fix`, `refactor`, `docs`, `test`, `chore`.
