# AGENTS.md — KTravels PublicSite API

Reglas y contexto para agentes de IA que trabajen en este repositorio.

## Stack real (verificado contra go.mod)

| Componente   | Biblioteca / Tecnología                      |
| ------------ | -------------------------------------------- |
| Lenguaje     | Go 1.26.2                                    |
| API          | GraphQL (gqlgen schema-first)                |
| Base de datos| MongoDB (mongo-driver v1.17)                 |
| Config       | spf13/viper — env vars con prefijo `APP_`    |
| Logging      | rs/zerolog — JSON estructurado               |
| Mensajería   | RabbitMQ (amqp091-go)                        |
| Testing      | `testing` estándar + `testify`               |
| Contenedores | Docker + Docker Compose                      |

## Arranque y comandos

```bash
# Levantar todo (API + MongoDB)
docker-compose up -d

# Build local
go build -o api ./app/api

# Ejecutar todos los tests
go test ./...

# Generar código GraphQL (gqlgen schema-first)
go run github.com/99designs/gqlgen generate
```

**Orden importante:** Si modificás `graph/schema.graphqls` o algún `schema.graphqls` dentro de un feature, corré `gqlgen generate` antes de compilar. El archivo generado `graph/generated.go` está en `.gitignore`.

## Configuración y variables de entorno

Viper usa prefijo `APP_` y reemplaza `_` por `.` (notación de nested keys).
La variable `app.database.uri` se expone como `APP_DATABASE_URI`.

`.env.example` **está desactualizado** — usa nombres como `MONGODB_URI` que viper no leerá. El docker-compose usa `DATABASE_URI` sin prefijo `APP_`, que tampoco funciona. Guiarse por `internal/platform/config/config.go:50-51` como fuente de verdad.

## Arquitectura: Vertical Slice

Cada feature se organiza dentro de `internal/features/{feature}/` con esta estructura:

```
internal/features/{feature}/
├── domain/         # entity.go, repository.go (interfaces), rules.go
├── application/    # casos de uso (create, get, list, etc.)
├── infrastructure/ # mongo_repository.go, document.go (modelos de persistencia)
├── delivery/
│   ├── graphql/    # resolver.go, mapper.go, schema.graphqls
│   └── rabbitmq/   # consumer.go (cuando el feature se activa por mensajería)
```

**Reglas no negociables:**
1. Los resolvers GraphQL deben ser delgados; la lógica vive en application/domain.
2. El acceso a MongoDB siempre pasa por repositorios, nunca directo desde resolvers.
3. No crear carpetas horizontales globales (`controllers/`, `services/`, `repositories/`).
4. Código compartido va en `internal/shared/` o `internal/platform/`; no duplicar entre slices.
5. No mover lógica a `shared` prematuramente — esperar a que haya reutilización real.

## Estado actual del proyecto

El repo está en etapa de bootstrap. Lo que existe:
- `internal/platform/` — wiring técnico: config, database (MongoDB), graphql (stub), logger, rabbitmq
- `internal/shared/errors/` — errores centinela (`ErrNotFound`, `ErrInvalidInput`, etc.)
- `internal/shared/pagination/` — `PageInfo` con página 0-indexada
- `graph/schema.graphqls` — schema base con `Query.health` y `Mutation.noop`
- `internal/features/placeholder.go` — documenta la estructura esperada
- `internal/features/provider_service/` — primer feature: consume eventos `provider.service.published` vía RabbitMQ y hace upsert en MongoDB

`main.go` levanta config, logger, MongoDB, y RabbitMQ. El servidor HTTP aún no tiene handlers registrados. El consumer de RabbitMQ está activo y escucha en la cola `provider_service_published`.

## Nombres y estilo

- Archivos Go: `snake_case`
- Variables/privadas: `camelCase`
- Funciones públicas/tipos: `PascalCase`
- Constantes: `UPPER_SNAKE_CASE`
- Paquetes: `lowercase`, una palabra
- Código y comentarios en código: **inglés**
- Documentación (.md, commits, PRs): **español**

GraphQL:
- Tipos y queries: `PascalCase`
- Campos: `camelCase`
- Inputs: sufijo `Input`
- Payloads de mutación: sufijo `Payload`

MongoDB:
- Colecciones: `snake_case` plural
- Separar documentos de persistencia de modelos de dominio

## Errores

Usar los centinelas de `internal/shared/errors` como base. No usar `panic()` para errores de negocio. No ignorar errores con `_`. Mapear errores de forma consistente a GraphQL sin exponer detalles internos.

## Ramas y commits

Gitflow: `main` (producción) ← `develop` (integración) ← `feature/{name}`.

Commits: `<type>(<scope>): <description>` en español.
Tipos: `feat`, `fix`, `refactor`, `docs`, `test`, `chore`.

## Testing

Se usa `testing` estándar de Go. `testify` se recomienda cuando aporte claridad. Probar lógica de negocio sin depender de MongoDB real. Repositorios con integración cuando haya consultas relevantes.
