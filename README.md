# pismo

REST API for the Pismo code assessment.

## Stack
- Go 1.26
- Chi router
- SQLite
- OpenTelemetry
- Jaeger
- swaggo/swag

## Run locally
Install dependencies and start the API:
```bash
make tidy
make run
```

Or directly:
```bash
go mod tidy
go run ./cmd/api
```

The API starts on `http://localhost:8080`.

Environment variables:
- `PORT` - HTTP port
- `DATABASE_URL` - SQLite DSN, default: `file:pismo.db?_pragma=foreign_keys(1)`
- `OTEL_SERVICE_NAME` - service name sent in traces, default: `pismo-api`
- `OTEL_EXPORTER_OTLP_ENDPOINT` - OTLP HTTP endpoint, example: `http://localhost:4318`
- `OTEL_EXPORTER_OTLP_INSECURE` - set to `true` for local Jaeger

## Migrations
This project currently applies the database schema automatically on startup.

When the app opens SQLite in `internal/repository/sqlite/db.go:25`, it runs the migration function in `internal/repository/sqlite/db.go:42`, which executes the schema for:
- `accounts`
- `transactions`

So in normal usage you do not need to run a separate migration command. Just start the app:
```bash
make run
```

Or:
```bash
go run ./cmd/api
```

If you want to create a fresh database, delete the local SQLite file and start the app again:
```bash
rm -f pismo.db pismo.db-shm pismo.db-wal
make run
```

If you want to inspect the resulting schema manually:
```bash
sqlite3 pismo.db ".schema"
```

## Run tests
```bash
make test
```

## Format code
```bash
make fmt
```

## Generate swagger docs
Install the CLI once:
```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

Generate docs from code annotations:
```bash
make swagger
```

Or directly:
```bash
swag init -g cmd/api/main.go -o docs
```

## Docker Compose
Start the full stack with API + Jaeger:
```bash
make compose-up
```

Stop it:
```bash
make compose-down
```

Or directly:
```bash
docker compose up --build
```

Services:
- API: `http://localhost:8080`
- Swagger UI: `http://localhost:8080/swagger/index.html`
- Jaeger UI: `http://localhost:16686`

## Podman
Build the image:
```bash
make podman-build
```

Start Jaeger + API with Podman:
```bash
make podman-up
```

Stop the Podman stack:
```bash
make podman-down
```

## API examples
### Create account
```bash
curl -i -X POST http://localhost:8080/accounts \
  -H 'Content-Type: application/json' \
  -d '{"document_number":"12345678900"}'
```

### Get account
```bash
curl -i http://localhost:8080/accounts/1
```

### Create payment transaction
```bash
curl -i -X POST http://localhost:8080/transactions \
  -H 'Content-Type: application/json' \
  -d '{"account_id":1,"operation_type_id":4,"amount":123.45}'
```

### Create purchase transaction
```bash
curl -i -X POST http://localhost:8080/transactions \
  -H 'Content-Type: application/json' \
  -d '{"account_id":1,"operation_type_id":1,"amount":50.00}'
```

For purchase, installment purchase, and withdrawal, the API stores the amount as negative. For payment, it stores the amount as positive.

## Tracing verification
1. Start the stack with `docker compose up --build`
2. Call any API route
3. Open Jaeger UI at `http://localhost:16686`
4. Search for service `pismo-api`

## Swagger/OpenAPI
The app serves generated swagger docs at:
- `GET /swagger/index.html`

When handler annotations change, regenerate docs with:
```bash
swag init -g cmd/api/main.go -o docs
```
