Este documento define las reglas, convenciones y contexto que todo agente de IA (Copilot, Cursor, Cline, Kiro, etc.) debe seguir al asistir en el desarrollo de esta API.

---

## 1. Contexto del Proyecto

Esta API pertenece al sitio publico de KTravels y se construira con los siguientes criterios base:

- Lenguaje principal: Go (Golang)
- Base de datos: MongoDB
- Estilo de API: GraphQL
- Arquitectura: Vertical Slice Architecture
- Flujo de ramas: Gitflow
- Idioma del codigo: ingles
- Idioma de la documentacion: espanol

La API debe exponer datos de forma eficiente para que el frontend consuma solo los campos que realmente necesita, evitando sobre-fetching y manteniendo una separacion clara entre dominio, aplicacion, infraestructura y capa de entrega GraphQL.

---

## 2. Reglas para Agentes de IA

Estas reglas son obligatorias para cualquier agente que genere, modifique o revise codigo en este repositorio.

1. Antes de escribir codigo, identificar primero el vertical slice afectado y trabajar dentro de ese contexto.
2. No crear estructuras paralelas ni capas duplicadas si ya existe un patron establecido en el slice correspondiente.
3. Mantener los resolvers de GraphQL delgados. La logica de negocio debe vivir en la capa de aplicacion o dominio.
4. Todo acceso a MongoDB debe pasar por repositorios o adaptadores de infraestructura; nunca desde resolvers directamente.
5. Todo cambio funcional debe incluir o actualizar pruebas cuando el alcance lo justifique.
6. No introducir dependencias nuevas sin necesidad clara y sin alinearlas al stack ya definido.
7. No inventar comportamiento de negocio no especificado. Si falta una regla funcional, dejarla explicita en la documentacion o asumir lo minimo indispensable.
8. Respetar Gitflow: no sugerir ni ejecutar trabajo directo sobre `main` o `develop`.
9. Todo codigo, nombres de tipos, funciones, variables y comentarios en codigo deben ir en ingles.
10. La documentacion tecnica y funcional del repositorio debe escribirse en espanol, salvo que exista una razon para hacerlo de otra forma.

---

## 3. Stack Tecnologico

| Capa          | Tecnologia                     |
| ------------- | ------------------------------ |
| API Backend   | Go (Golang)                    |
| Base de datos | MongoDB                        |
| API contract  | GraphQL                        |
| Contenedores  | Docker + Docker Compose        |
| Testing       | `testing` estandar + `testify` |
| Configuracion | Variables de entorno           |

### Notas

- Si no existe una decision previa en el repositorio, preferir enfoque schema-first para GraphQL.
- Si no existe una libreria ya adoptada, se recomienda `gqlgen` para la implementacion de GraphQL en Go.
- MongoDB debe modelarse con colecciones y documentos alineados al dominio, evitando acoplar el modelo de persistencia al contrato GraphQL.

---

## 4. Convenciones de Codigo

### Generales

- Variables: `camelCase`
- Constantes: `UPPER_SNAKE_CASE`
- Funciones publicas: `PascalCase`
- Funciones privadas: `camelCase`
- Nombres de archivos Go: `snake_case`
- Paquetes Go: `lowercase`, preferiblemente una sola palabra
- Codigo y comentarios en codigo: ingles
- Documentacion Markdown: espanol

### Go

- Structs: `PascalCase`
- Interfaces: `PascalCase` con nombre descriptivo, por ejemplo `TripRepository`, `BookingService`
- El `error` siempre debe ser el ultimo valor de retorno
- Recibir `context.Context` en operaciones de aplicacion, infraestructura y acceso a datos
- Preferir composicion sobre herencia o abstracciones innecesarias
- Mantener funciones pequenas y con una sola responsabilidad
- Evitar acoplar modelos de dominio a detalles de transporte o persistencia

### GraphQL

- Tipos GraphQL: `PascalCase`
- Campos GraphQL: `camelCase`
- Inputs GraphQL: sufijo `Input`
- Payloads de mutaciones: sufijo `Payload` cuando aplique
- Queries: nombradas por intencion de lectura, no por tecnologia
- Mutations: nombradas por accion de negocio, por ejemplo `createBooking`, `cancelReservation`
- No exponer campos internos de persistencia si no son parte del contrato funcional
- Implementar paginacion en listados potencialmente grandes
- Evitar resolver N+1; usar batching o DataLoader cuando aplique

### MongoDB

- Las colecciones deben nombrarse en `snake_case` y en plural cuando tenga sentido de negocio
- Definir indices explicitamente para consultas criticas
- Evitar documentos excesivamente anidados si afectan mantenimiento o consulta
- No usar agregaciones complejas desde la capa de presentacion; encapsularlas en repositorios o servicios de consulta
- Separar claramente DTOs, modelos de dominio y modelos de persistencia si sus responsabilidades divergen

---

## 5. Principios de Arquitectura

La base del proyecto es **Vertical Slice Architecture**, con separacion clara de responsabilidades dentro de cada slice.

### Reglas base

1. Cada feature debe organizarse alrededor del caso de uso, no de la tecnologia.
2. La capa GraphQL solo orquesta entrada y salida.
3. La logica de negocio pertenece a servicios de aplicacion o dominio.
4. El acceso a MongoDB debe encapsularse en repositorios.
5. El codigo compartido debe vivir en `internal/shared` o `internal/platform`, no replicarse entre slices.
6. Usar CQRS solo cuando la complejidad lo justifique. No dividir comandos y queries de forma artificial.

### Patrones recomendados

| Patron          | Cuando usarlo                                                              |
| --------------- | -------------------------------------------------------------------------- |
| Repository      | Siempre para acceso a datos                                                |
| CQRS            | Cuando lectura y escritura tengan reglas o modelos distintos               |
| Factory         | Cuando la creacion de objetos de dominio requiera validaciones o variantes |
| Strategy        | Cuando existan reglas intercambiables de negocio                           |
| Observer/Events | Para efectos secundarios desacoplados                                      |
| DataLoader      | Para evitar N+1 en resolvers GraphQL                                       |

---

## 6. Prohibiciones

Estas reglas son obligatorias. El agente no debe generar codigo que las viole.

1. No usar `panic()` para errores de negocio.
2. No acceder a MongoDB directamente desde resolvers GraphQL.
3. No mezclar logica de negocio con logica de infraestructura.
4. No hardcodear credenciales, secrets, URIs o configuraciones por entorno.
5. No ignorar errores. Todo error debe manejarse o propagarse explicitamente.
6. No usar `interface{}` ni `any` en Go sin justificacion real.
7. No crear resolvers gigantes con validacion, acceso a datos y mapeo mezclados.
8. No hacer push directo a `main` ni a `develop`.
9. No commitear dead code, codigo comentado o archivos temporales.
10. No romper la organizacion vertical slice agregando carpetas horizontales para controllers, services o repositories globales.
11. No exponer detalles internos de MongoDB, como nombres de coleccion o estructura interna de documentos, en el contrato GraphQL.
12. No crear queries GraphQL sin limites, filtros o paginacion cuando el volumen de datos pueda crecer.

---

## 7. Estructura del Proyecto

La estructura debe seguir Vertical Slice Architecture. Cada feature agrupa su entrada GraphQL, aplicacion, dominio e infraestructura relacionada.

```text
ktravels-publicsite-api/
├── app/
│   └── api/
│       └── main.go                  # Punto de entrada
├── internal/
│   ├── platform/                    # Wiring tecnico: config, server, db, logger
│   │   ├── config/
│   │   ├── database/
│   │   ├── graphql/
│   │   └── logger/
│   ├── shared/                      # Utilidades y contratos compartidos
│   │   ├── errors/
│   │   ├── pagination/
│   │   └── utils/
│   └── features/
│       └── trips/
│           ├── domain/
│           │   ├── entity.go
│           │   ├── repository.go
│           │   └── rules.go
│           ├── application/
│           │   ├── create_trip.go
│           │   ├── get_trip.go
│           │   └── list_trips.go
│           ├── infrastructure/
│           │   ├── mongo_repository.go
│           │   └── document.go
│           └── delivery/
│               └── graphql/
│                   ├── resolver.go
│                   ├── mapper.go
│                   └── schema.graphqls
├── graph/                           # Archivos generados o wiring GraphQL si la libreria lo requiere
├── tests/
│   ├── integration/
│   └── fixtures/
├── docs/
├── Dockerfile
├── docker-compose.yml
├── .env.example
├── go.mod
└── go.sum
```

### Regla de organizacion

- Crear nuevos cambios dentro del feature existente correspondiente.
- Si el caso de uso pertenece a un feature nuevo, crear un nuevo slice en `internal/features/{feature}`.
- Evitar mover logica comun a `shared` demasiado pronto. Primero debe demostrarse reutilizacion real.

---

## 8. Flujo de Trabajo

### Proceso general

1. Levantar el requerimiento con contexto funcional claro.
2. Definir el caso de uso y el contrato GraphQL que lo representa.
3. Implementar el slice correspondiente respetando dominio, aplicacion, infraestructura y delivery.
4. Agregar pruebas unitarias e integracion segun corresponda.
5. Validar manualmente la query o mutation y sus escenarios de error.
6. Abrir Pull Request hacia la rama objetivo segun Gitflow.

### Gitflow

```text
main ----------------------------------------------- produccion
 |
 └── hotfix/fix-trip-filter ---- merge -> main + develop
 |
develop -------------------------------------------- integracion
 |
 ├── feature/trips-list -------- PR -> develop
 ├── feature/trip-detail ------- PR -> develop
 └── release/v1.0.0 ------------ merge -> main + develop
```

### Reglas de ramas

- `main`: codigo en produccion
- `develop`: rama de integracion
- `feature/{name}`: nueva funcionalidad, creada desde `develop`
- `release/v{X.Y.Z}`: preparacion para despliegue, creada desde `develop`
- `hotfix/{name}`: correccion urgente, creada desde `main`

### Flujo del desarrollador

1. Crear rama desde `develop`: `git checkout -b feature/trips-list develop`
2. Implementar el cambio respetando la arquitectura del proyecto
3. Ejecutar pruebas locales
4. Abrir PR hacia `develop`
5. Esperar code review y validacion de CI
6. Mergear solo despues de aprobacion

---

## 9. Testing

### Estrategia

- Pruebas unitarias para dominio y aplicacion
- Pruebas de integracion para repositorios MongoDB y resolvers criticos
- Pruebas de contrato para queries y mutations relevantes
- Validar casos felices, validaciones y errores de negocio

### Herramientas

- `testing` estandar de Go
- `testify` para assertions y suites cuando aporte claridad

### Reglas

1. Toda logica de negocio debe ser testeable sin depender de GraphQL ni de MongoDB real.
2. Los repositorios deben probarse con integracion cuando haya consultas, filtros, indices o agregaciones relevantes.
3. Toda mutation relevante debe tener al menos una prueba de exito y una de fallo.
4. Si se agrega paginacion, filtros o sorting, deben incluirse pruebas de comportamiento.

---

## 10. CI/CD y Despliegue

### Contenedores

- La API debe tener su propio `Dockerfile`
- El entorno local debe poder levantarse con `docker-compose.yml`
- Preferir multi-stage builds para reducir tamano de imagen
- No incluir `.env` ni secretos reales dentro de la imagen

### Pipeline minimo por PR

1. Lint
2. Build
3. Test
4. Build de imagen Docker

### Reglas

- El pipeline debe fallar si fallan pruebas o validaciones basicas
- Toda imagen debe exponer solo los puertos necesarios
- Agregar `HEALTHCHECK` cuando la imagen y el entorno lo requieran

---

## 11. Estilo de Commits y Pull Requests

### Formato de commit

```text
<type>(<scope>): <description>
```

- `type`: `feat`, `fix`, `refactor`, `docs`, `test`, `chore`, `style`, `perf`, `ci`
- `scope`: feature o archivo principal afectado
- `description`: en espanol, clara y concreta

### Ejemplos

```text
feat(trips): agregar query para listar viajes publicados con paginacion
fix(graphql): corregir validacion de filtros vacios en busqueda de destinos
refactor(trips): separar mapper de documentos MongoDB del dominio
test(trips): agregar pruebas de integracion para repositorio de viajes
docs(agents): definir reglas de arquitectura vertical slice y gitflow
```

### Pull Requests

- El titulo debe seguir el mismo formato del commit principal
- La descripcion debe indicar que se cambio, por que y como probarlo
- Todo PR debe apuntar a la rama correcta segun Gitflow
- No mergear sin aprobacion ni sin CI en verde

---

## 12. Manejo de Errores

### Reglas generales

- En Go, siempre retornar `error` como ultimo valor
- No ignorar errores con `_`
- Agregar contexto al error cuando ayude a diagnostico, sin filtrar informacion sensible
- Diferenciar errores de validacion, negocio, not-found, conflicto e infraestructura

### GraphQL

- Los errores deben mapearse de forma consistente desde aplicacion a GraphQL
- No exponer trazas internas, queries MongoDB ni detalles sensibles en mensajes al cliente
- Usar mensajes claros para validaciones y reglas de negocio
- Los errores tecnicos deben loguearse con contexto suficiente para diagnostico

### Categorias recomendadas

- `VALIDATION_ERROR`
- `NOT_FOUND`
- `CONFLICT`
- `UNAUTHORIZED`
- `FORBIDDEN`
- `INTERNAL_ERROR`

---

## 13. Logging y Observabilidad

### Reglas

- Usar structured logging en formato JSON
- Incluir como minimo: timestamp, level, service, operation, message y contexto relevante
- Nunca loguear secretos, tokens completos ni datos sensibles
- En produccion, el nivel minimo recomendado es `info`

### Ejemplo

```json
{
  "level": "error",
  "timestamp": "2026-04-11T10:15:30Z",
  "service": "ktravels-publicsite-api",
  "operation": "Trips.List",
  "message": "failed to fetch published trips",
  "error": "mongo: context deadline exceeded"
}
```

---

## 14. Variables de Entorno y Configuracion

- Toda configuracion debe venir de variables de entorno
- El repositorio debe incluir `.env.example` sin secretos reales
- No definir valores por defecto inseguros para produccion
- Centralizar la lectura de configuracion en `internal/platform/config`

### Ejemplo de `.env.example`

```env
APP_NAME=ktravels-publicsite-api
APP_ENV=development
APP_PORT=8080
LOG_LEVEL=debug

MONGODB_URI=mongodb://localhost:27017
MONGODB_DATABASE=ktravels_publicsite

GRAPHQL_PLAYGROUND_ENABLED=true
GRAPHQL_INTROSPECTION_ENABLED=true
```

---

## 15. Buenas Practicas Especificas para GraphQL

1. Diseñar el schema desde el lenguaje del negocio, no desde la estructura de MongoDB.
2. Mantener los resolvers delgados y delegar en casos de uso.
3. Evitar campos que disparen consultas innecesarias por cada nodo.
4. Implementar paginacion, filtros y ordenamiento de forma explicita en listados.
5. No usar GraphQL como excusa para devolver estructuras ambiguas o poco tipadas.
6. Mantener separacion entre modelos GraphQL, dominio y persistencia.
7. Habilitar playground o introspection solo segun entorno y configuracion.

---

## 16. Definicion de Hecho Minima

Un cambio se considera completo cuando cumple como minimo con lo siguiente:

1. Respeta la estructura vertical slice
2. La logica de negocio no vive en resolvers
3. El acceso a datos pasa por repositorios
4. El contrato GraphQL esta alineado al caso de uso
5. Incluye pruebas razonables para el alcance
6. No introduce secretos, codigo muerto ni acoplamientos innecesarios
7. Esta listo para integrarse por Pull Request siguiendo Gitflow
