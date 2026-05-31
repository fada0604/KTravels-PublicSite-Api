# Historias de Usuario — Multimoneda

## Contexto

KTravels maneja precios en USD a nivel del BackOffice. En Venezuela, la moneda legal es el bolívar (VES) pero los precios se expresan en dólares. Se requiere mostrar los precios en USD y, si el visitante lo prefiere (vía header `X-Visitor-Currency`), el equivalente en su moneda local.

A futuro, otros países/divisas podrán agregarse sin romper el modelo.

## Arquitectura objetivo

```
BO API (PostgreSQL)           Background Service
  │                              │
  │  CRUD monedas                │  Cron 2x/día
  │  Set manual de tasa          │  Verifica date_restart en PostgreSQL
  │                              │  APIs externas con fallback + email
  │                              │
  └──────────┬───────────────────┘
             │  Publican evento RabbitMQ
             │  { event_id, pair, rate, source, timestamp }
             ▼
      Public Site API (MongoDB)
         │
         ├─► exchange_rates (histórico, idempotente por event_id)
         ├─► currencies (monedas habilitadas)
         ├─► cache in-memory: map[pair]rate
         │
         ▼
      GraphQL: X-Visitor-Currency → displayPrices
```

**Decisiones de diseño:**
- Redis **no** se usa en el Public Site. La tasa viaja en el payload del evento RabbitMQ.
- Caché in-memory: 0 latencia, se actualiza solo al recibir el evento.
- Idempotencia: `event_id` UUID en cada mensaje, evita duplicados en MongoDB.
- Monedas **independientes** de los países de operación turística. Se pueden agregar divisas de visitantes sin operar en ese país.
- Solo tasas USD → local. Sin conversiones cruzadas.

---

## US1 — BO: Administrar monedas habilitadas

**Como** administrador del BackOffice  
**Quiero** una pantalla para gestionar las monedas disponibles para conversión  
**Para** controlar qué divisas se muestran en el Public Site.

### Criterios de aceptación
- CRUD de `currency` con campos: `code` (ISO 4217), `symbol`, `name`, `decimals`, `enabled`, `date_restart`.
- `date_restart` (DateTime, nullable): si tiene fecha futura, el Background Service **no** ejecutará el fetch automático hasta que `now >= date_restart`. Útil para pausar la actualización automática cuando las APIs externas están desactualizadas y se requiere ingreso manual.
- Las monedas **no** heredan de la tabla de países. Son independientes. Ejemplo: se puede agregar `INR` (Rupia india) aunque no se opere en India.
- Al crear/modificar/eliminar una moneda, se publica el evento `currency.configured` a RabbitMQ.

### RabbitMQ — `currency.configured`
```json
{
  "event_id": "019a6c5e-8f72-73b1-a44c-3e4f5d6a7b8c",
  "code": "VES",
  "symbol": "Bs.",
  "name": "Bolívar",
  "decimals": 2,
  "enabled": true,
  "date_restart": null,
  "timestamp": "2026-05-09T14:00:00Z"
}
```

---

## US2 — BO: Registrar tasa de cambio manual

**Como** administrador del BackOffice  
**Quiero** establecer manualmente una tasa de cambio USD → moneda local con fecha de vigencia  
**Para** cuando las APIs externas fallan o la tasa real cambió y las APIs aún no lo reflejan.

### Criterios de aceptación
- Pantalla con: par de monedas (from/to), tasa, fecha de vigencia.
- Guarda en PostgreSQL (`exchange_rates` con histórico). `source = "manual"`.
- Publica evento `exchange_rate.updated` a RabbitMQ.

### RabbitMQ — `exchange_rate.updated`
```json
{
  "event_id": "019a6c5e-8f72-73b1-a44c-3e4f5d6a7b8c",
  "currency_from": "USD",
  "currency_to": "VES",
  "rate": 54.32,
  "source": "manual",
  "valid_from": "2026-05-09T14:00:00Z",
  "timestamp": "2026-05-09T14:00:01Z"
}
```

---

## US3 — Background Service: Obtener tasa automática 2x/día

**Como** sistema  
**Quiero** un job programado que obtenga tasas USD → moneda local desde APIs externas configurables dos veces al día  
**Para** mantener las tasas actualizadas sin intervención manual y poder cambiar de APIs sin redeploy.

### A) Tabla de fuentes en BO (PostgreSQL): `exchange_rate_sources`

Las APIs no se hardcodean. Se parametrizan en base de datos para que el Background Service las lea en cada ejecución. Agregar, quitar o reordenar APIs no requiere detener ni redeployar el servicio.

| Columna | Tipo | Descripción |
|---|---|---|
| `id` | UUID | PK |
| `name` | VARCHAR | Nombre descriptivo (ej: `exchangerate.fun`) |
| `base_url` | VARCHAR | URL con placeholders: `{target}`, `{api_key}` |
| `api_key` | VARCHAR (encrypted) | API key encriptada (AES-256). Si no requiere, `null`. |
| `currency_filter` | VARCHAR, nullable | `null` = todas las monedas, `VES` = solo VES, `!VES` = excepto VES |
| `priority` | INT | Orden de intento (menor = primero) |
| `enabled` | BOOLEAN | Toggle sin borrar registro |
| `rate_json_path` | VARCHAR | JSONPath para extraer la tasa (ej: `$.rates.{target}`) |
| `timestamp_json_path` | VARCHAR, nullable | JSONPath para el timestamp de actualización de la API |
| `headers` | JSONB, nullable | Headers HTTP custom. Soporta placeholder `{api_key}` para inyectar la key desencriptada en runtime. |
| `created_at` | TIMESTAMP | — |
| `updated_at` | TIMESTAMP | — |

**Regla de seguridad**: `api_key` se almacena encriptado con AES-256. La clave de encriptación se inyecta por variable de entorno (`APP_ENCRYPTION_KEY`). Solo el Background Service desencripta en memoria al usarla. Nunca se loguea ni se expone en respuestas.

### B) Configuración inicial (4 APIs)

| # | Name | URL | Filter | JSONPath rate | Prioridad |
|---|---|---|---|---|---|
| 1 | exchangerate.fun | `https://api.exchangerate.fun/latest?base=USD` | `!VES` | `$.rates.{target}` | 1 |
| 2 | exchangerate-api.com | `https://v6.exchangerate-api.com/v6/{api_key}/latest/USD` | `!VES` | `$.conversion_rates.{target}` | 2 |
| 3 | exchangedyn.com | `https://api.exchangedyn.com/markets/quotes/usdves/bcv` | `VES` | `$.quotes.usdves.bcv` | 3 |
| 4 | bcv-api.rafnixg.dev | `https://bcv-api.rafnixg.dev/rates/` | `VES` | `$.rates.USD` | 4 |

### C) Flujo del Background Service

```
Cron 2x/día (horas configurables, ej: 8 AM y 2 PM hora Venezuela)

Para cada currency_pair activo (ej: USD→VES):

  1. Consultar date_restart de la moneda en DB del BO.
     Si now < date_restart → skip y loguear.

  2. SELECT * FROM exchange_rate_sources
     WHERE enabled = true
       AND (currency_filter IS NULL
            OR currency_filter = 'VES'
            OR currency_filter = '!VES')
     ORDER BY priority ASC

  3. Para cada source en orden:
     a. Reemplazar placeholders en URL y headers:
        - {target} → 'VES'
        - {api_key} → desencriptar api_key de la BD
     b. HTTP GET con headers si aplica
     c. Extraer rate con JSONPath (rate_json_path)
        Si falla el parseo → next source
     d. Si éxito → break

  4. Si todas las fuentes fallan → email al admin
     para que establezca la tasa manualmente (US2).

  5. Si éxito:
     a. Guardar en PostgreSQL (exchange_rates, histórico)
        con source = nombre de la API que respondió
     b. Publicar evento exchange_rate.updated a RabbitMQ
```

### D) Ventajas del diseño parametrizado

- **Cambiar orden**: editar `priority` en DB → el servicio lee en la próxima ejecución.
- **Deshabilitar una API**: `enabled = false` → se saltea sin borrar el registro.
- **Agregar nueva API**: INSERT con URL y JSONPath → cero código. El servicio la usa en el siguiente ciclo.
- **Quitar API**: soft-delete o `enabled = false`.
- **Rotar API keys**: actualizar `api_key` en DB con el nuevo valor encriptado.
- Solo se requiere redeploy si una API requiere un mecanismo de autenticación que no se pueda parametrizar (ej: OAuth2 con refresh tokens, mTLS con certificados cliente). Para el 95% de las APIs REST con JSON, JSONPath y headers parametrizados son suficientes.

---

## US4 — Public Site: Consumir eventos de tasa de cambio y monedas

**Como** Public Site API  
**Quiero** reaccionar a los eventos `exchange_rate.updated` y `currency.configured`  
**Para** mantener mi caché en memoria actualizada y el histórico en MongoDB.

### Criterios de aceptación
- **Consumer RabbitMQ** en `internal/features/exchange_rate/delivery/rabbitmq/consumer.go`:
  - Escucha `exchange_rate.updated`.
  - **Idempotencia**: verifica `event_id` en MongoDB. Si ya existe → ACK y skip. Si no, inserta en `exchange_rates` y actualiza `map["USD:VES"]` en memoria.
- **Consumer RabbitMQ** en `internal/features/currency/delivery/rabbitmq/consumer.go`:
  - Escucha `currency.configured`.
  - Upsert en colección `currencies` (MongoDB) y actualiza caché en memoria `map[code]CurrencyInfo`.
- **Precarga al iniciar**: el PS consulta MongoDB para cargar la última tasa vigente de cada par y las monedas habilitadas a la caché en memoria. Si MongoDB no responde, arranca con caché vacía y espera el próximo evento.
- **GraphQL Query** `exchangeRates`: expone el histórico paginado (para auditoría).

### Colecciones MongoDB

```
exchange_rates
├── _id
├── event_id: "019a6c5e-..." (unique index)
├── currency_from: "USD"
├── currency_to: "VES"
├── rate: 54.32
├── source: "exchangerate.fun" | "exchangedyn.com" | "manual"
├── valid_from: ISODate
├── timestamp: ISODate

currencies
├── _id
├── code: "VES" (unique index)
├── symbol: "Bs."
├── name: "Bolívar"
├── decimals: 2
├── enabled: true
├── date_restart: null | ISODate
```

---

## US5 — Public Site: Exponer precios convertidos en GraphQL

**Como** visitante del Public Site  
**Quiero** ver los precios en mi moneda local (ej: VES) además de USD  
**Para** entender el costo sin hacer la conversión mentalmente.

### Criterios de aceptación
- El header HTTP `X-Visitor-Currency` (ej: `VES`) se extrae en el contexto GraphQL.
- Si el header está ausente o la moneda no está en la lista `enabled` → solo se retorna el precio en USD.
- Si la moneda está enabled y tiene tasa cargada → se retorna USD + conversión.
- La conversión usa la tasa en caché in-memory (sin consulta a MongoDB por request).
- **La respuesta siempre incluye USD como primer elemento** de `displayPrices`.

### Schema GraphQL

```graphql
type Money {
  price: Float!
  currency: String!
  symbol: String!
}

type AgeRangeRate {
  id: String!
  groupAgeRange: Int!
  name: String!
  status: Boolean!
  minAge: Int!
  maxAge: Int!
  rate: Float!
  displayPrices: [Money!]!
}

type ExceptionRateDetail {
  id: String!
  calendarExceptionId: String!
  exceptionName: String!
  rate: Float!
  displayPrices: [Money!]!
}

type CalendarRateDetail {
  rate: Float!
  displayPrices: [Money!]!
  exceptionRates: [ExceptionRateDetail!]!
}

type AgeRangeExceptionRate {
  id: String!
  calendarExceptionId: String!
  exceptionName: String!
  ageRangeRates: [AgeRangeRate!]!
}

type CalendarPerPersonRateDetail {
  ageRangeRates: [AgeRangeRate!]!
  exceptionRates: [AgeRangeExceptionRate!]!
}

type Rate {
  id: String!
  rateType: Int!
  calendar: Calendar!
  calendarRate: CalendarRateDetail
  calendarPerPersonRate: CalendarPerPersonRateDetail
}
```

### Ejemplo de respuesta (`X-Visitor-Currency: VES`, tasa 54.32)

```json
{
  "calendarRate": {
    "rate": 500.00,
    "displayPrices": [
      { "price": 500.00, "currency": "USD", "symbol": "$" },
      { "price": 27160.00, "currency": "VES", "symbol": "Bs." }
    ],
    "exceptionRates": [
      {
        "rate": 600.00,
        "displayPrices": [
          { "price": 600.00, "currency": "USD", "symbol": "$" },
          { "price": 32592.00, "currency": "VES", "symbol": "Bs." }
        ]
      }
    ]
  }
}
```

---

## US6 — Public Site: Consulta de tasas históricas para reportes

**Como** sistema de facturación/reportes  
**Quiero** consultar la tasa de cambio vigente en una fecha específica  
**Para** calcular equivalencias retroactivas en reportes financieros y facturación histórica.

### Criterios de aceptación
- Query GraphQL:
  ```graphql
  extend type Query {
    exchangeRate(pair: String!, date: Time!): ExchangeRate
  }
  ```
- Busca en MongoDB la tasa con `currency_from + currency_to` cuyo `valid_from` sea el más cercano **anterior o igual** a la fecha consultada.
- Índice compuesto en MongoDB: `{ currency_from: 1, currency_to: 1, valid_from: -1 }`.

### Tipo GraphQL

```graphql
type ExchangeRate {
  currencyFrom: String!
  currencyTo: String!
  rate: Float!
  source: String!
  validFrom: Time!
}
```

---

## Resumen de eventos RabbitMQ

| Evento | Publicador | Consumidor | Payload clave |
|---|---|---|---|
| `currency.configured` | BO | Public Site | `code`, `symbol`, `name`, `decimals`, `enabled`, `date_restart` |
| `exchange_rate.updated` | BO / BG Service | Public Site | `event_id`, `currency_from`, `currency_to`, `rate`, `source`, `valid_from` |

**Regla de idempotencia**: todos los eventos incluyen `event_id` (UUID v7). El consumidor verifica si ya fue procesado antes de insertar.

---

## Estructura esperada en Public Site

```
internal/features/
├── exchange_rate/
│   ├── domain/
│   │   ├── entity.go           # ExchangeRate entity
│   │   └── repository.go       # Repository interface
│   ├── application/
│   │   └── handle_rate_updated.go
│   ├── infrastructure/
│   │   ├── mongo_repository.go
│   │   └── document.go
│   └── delivery/
│       ├── graphql/
│       │   ├── resolver.go
│       │   ├── mapper.go
│       │   └── schema.graphqls
│       └── rabbitmq/
│           └── consumer.go
├── currency/
│   ├── domain/
│   │   ├── entity.go           # Currency entity
│   │   └── repository.go
│   ├── application/
│   │   └── handle_configured.go
│   ├── infrastructure/
│   │   ├── mongo_repository.go
│   │   └── document.go
│   └── delivery/
│       ├── graphql/
│       │   ├── resolver.go
│       │   └── schema.graphqls
│       └── rabbitmq/
│           └── consumer.go
├── provider_service/
│   └── delivery/graphql/
│       └── resolver.go         # displayPrices resolver + X-Visitor-Currency
└── shared/
    └── currency/               # CurrencyConverter (in-memory cache + conversion logic)
        └── converter.go
```
