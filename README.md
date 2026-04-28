# KTravels PublicSite API

API GraphQL del sitio público de [KTravels](https://ktravels.com.ve), plataforma de servicios turísticos en Venezuela. Expone datos de viajes, destinos, reservas y servicios relacionados para el frontend público.

## Stack Tecnológico

| Capa       | Tecnología                |
|------------|---------------------------|
| Lenguaje   | Go 1.22+                  |
| API        | GraphQL                   |
| Base datos | MongoDB 7.0               |
| Logs       | zerolog (JSON estructurado) |
| Contenedores | Docker + Docker Compose |

## Arquitectura

El proyecto sigue **Vertical Slice Architecture**. Cada feature (viajes, reservas, destinos) agrupa su código en slices independientes, manteniendo separation clara entre:

- **Domain**: entidades y reglas de negocio
- **Application**: casos de uso
- **Infrastructure**: repositorios, adaptadores
- **Delivery**: resolvers GraphQL

```
internal/
├── platform/           # Config, logger, db, server
├── shared/             # Errores, paginación, utilerías
└── features/
    └── {feature}/
        ├── domain/
        ├── application/
        ├── infrastructure/
        └── delivery/graphql/
```

## Requisitos Previos

- Docker y Docker Compose instalados
- Puerto 8080 disponible (API)
- Puerto 27017 disponible (MongoDB)

## Levantar la Aplicación

### 1. Clonar el repositorio

```bash
git clone <repo-url>
cd KTravels-PublicSite-Api
```

### 2. Variables de entorno

Copiar `.env.example` a `.env` (opcional para desarrollo, los defaults funcionan):

```bash
cp .env.example .env
```

### 3. Iniciar servicios

```bash
docker-compose up -d
```

Esto levanta:
- **API**: `http://localhost:8080`
- **MongoDB**: `localhost:27017`
- **GraphQL Playground**: `http://localhost/graphql` (si `GRAPHQL_PLAYGROUND_ENABLED=true`)

### 4. Verificar que funciona

```bash
# Probar salud de la API
curl http://localhost:8080/health

# Probar GraphQL Playground
open http://localhost:8080
```

### 5. Detener servicios

```bash
docker-compose down
```

Para eliminar datos de MongoDB:

```bash
docker-compose down -v
```

## Configuración

| Variable | Descripción | Default |
|----------|-------------|---------|
| `APP_ENV` | Entorno (development/production) | development |
| `APP_PORT` | Puerto de la API | 8080 |
| `MONGODB_URI` | Connection string MongoDB | mongodb://localhost:27017 |
| `MONGODB_DATABASE` | Nombre de la base de datos | ktravels_publicsite |
| `GRAPHQL_PLAYGROUND_ENABLED` | Habilitar GraphQL Playground | true |
| `GRAPHQL_INTROSPECTION_ENABLED` | Habilitar introspección | true |
| `LOGGER_LEVEL` | Nivel de logs (debug/info/warn/error) | info |
| `LOGGER_FORMAT` | Formato de logs (json/console) | json |

## Estructura de Archivos

```
.
├── app/api/main.go              # Punto de entrada
├── internal/
│   ├── platform/
│   │   ├── config/              # Carga de configuración
│   │   ├── database/            # Conexión MongoDB
│   │   ├── graphql/              # Servidor GraphQL
│   │   └── logger/              # Logging estructurado
│   ├── shared/
│   │   ├── errors/              # Errores del dominio
│   │   └── pagination/          # Paginación
│   └── features/                # Features (viajes, destinos, etc.)
├── docker-compose.yml
├── Dockerfile
├── .env.example
└── README.md
```

## Desarrollo

### Agregar una nueva feature

1. Crear el slice en `internal/features/{feature}`.
2. Definir entidades en `domain/entity.go`.
3. Definir repositorio en `domain/repository.go`.
4. Implementar casos de uso en `application/`.
5. Implementar repositorio MongoDB en `infrastructure/`.
6. Definir schema GraphQL en `delivery/graphql/schema.graphqls`.
7. Implementar resolvers en `delivery/graphql/resolver.go`.

### Ejecutar tests

```bash
go test ./...
```

### Build local

```bash
go build -o api ./app/api
```

## Licencia

MIT License - Proyecto interno de KTravels.