# Product API — Go Clean Architecture

A RESTful Product API built with Go, Gin, PostgreSQL and clean architecture principles.

---

## Project Structure

```
product-api/
├── cmd/api/             # Entry point — main.go
├── internal/
│   ├── domain/          # Entities, interfaces, errors (no dependencies)
│   ├── usecase/         # Business logic / orchestrator
│   ├── repository/      # PostgreSQL data access
│   ├── handler/         # HTTP handlers (Gin)
│   └── component/       # E2E component tests
├── pkg/
│   ├── database/        # DB connection setup
│   └── response/        # Standard response envelope
├── docs/                # Swagger + SQL migrations
├── Dockerfile
└── docker-compose.yml
```

## Architecture Layers

```
Handler (HTTP) → Usecase (Business Logic) → Repository (Data Access) → PostgreSQL
```

All layers communicate through **interfaces** defined in `domain/`, enabling full testability via mocks.

---

## API Endpoints

| Method | Path             | Description               |
|--------|------------------|---------------------------|
| POST   | `/product`       | Create a product          |
| PATCH  | `/product/{id}`  | Partially update a product|
| GET    | `/api-docs/*`    | Swagger UI                |

### POST /product

```json
{
  "name": "Widget",
  "description": "Optional description",
  "price": 100.0,
  "sale_price": 79.99
}
```

**Response:**
```json
{
  "successful": true,
  "error_code": "",
  "data": { "data1": "<uuid>", "data2": "Widget" }
}
```

### PATCH /product/{id}

Only provided fields are updated (true partial update with `**T` pattern for nullable fields).

```json
{
  "name": "New Name",
  "sale_price": null
}
```

**Response:**
```json
{ "successful": true, "error_code": "", "data": null }
```

---

## How to Start

### Option 1: Docker Compose (recommended)

```bash
docker-compose up --build
```

API available at: `http://localhost:8080`  
Swagger docs: `http://localhost:8080/api-docs/index.html`

### Option 2: Local Go

```bash
# Start PostgreSQL, then:
export DB_HOST=localhost DB_PORT=5432 DB_USER=postgres DB_PASSWORD=postgres DB_NAME=productdb

# Apply migration
psql -U postgres -d productdb -f docs/001_create_products.up.sql

# Generate Swagger docs
swag init -g cmd/api/main.go

# Run
go run ./cmd/api
```

---

## Running Tests

### Unit tests (no DB required)

```bash
# Usecase tests (orchestrator)
go test ./internal/usecase/...

# Handler tests (service layer)
go test ./internal/handler/...
```

### Repository integration tests

```bash
export TEST_DB_DSN="host=localhost port=5432 user=postgres password=postgres dbname=productdb sslmode=disable"
go test ./internal/repository/...
```

### Component (E2E within service) tests

```bash
export TEST_DB_DSN="host=localhost port=5432 user=postgres password=postgres dbname=productdb sslmode=disable"
go test ./internal/component/...
```

### All tests

```bash
TEST_DB_DSN="..." go test ./...
```

---

## Test Coverage Map

| Test File | Layer | Type |
|---|---|---|
| `usecase/product_usecase_test.go` | Usecase | Unit (mocked repo) |
| `handler/product_handler_test.go` | Handler | Unit (mocked usecase) |
| `repository/product_repository_test.go` | Repository | Integration (real DB) |
| `component/component_test.go` | Full stack | E2E within service |

---

## Environment Variables

| Variable | Default | Description |
|---|---|---|
| `DB_HOST` | `localhost` | PostgreSQL host |
| `DB_PORT` | `5432` | PostgreSQL port |
| `DB_USER` | `postgres` | DB username |
| `DB_PASSWORD` | `postgres` | DB password |
| `DB_NAME` | `productdb` | Database name |
| `DB_SSLMODE` | `disable` | SSL mode |
| `PORT` | `8080` | HTTP server port |
| `TEST_DB_DSN` | _(none)_ | Full DSN for integration tests |

---

## Generate Swagger

```bash
go install github.com/swaggo/swag/cmd/swag@latest
swag init -g cmd/api/main.go -o docs/
```

Then visit: `http://localhost:8080/api-docs/index.html`
