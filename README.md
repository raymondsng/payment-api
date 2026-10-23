# Payment Transaction API

A simulated payment backend built with Go and PostgreSQL, focused on idempotency, database consistency, automated testing, and containerized deployment.

## Features
- REST API for creating and retrieving payments
- Idempotency keys to prevent duplicate payment records
- Request fingerprinting to detect conflicting retries
- PostgreSQL transactions and database-enforced uniqueness
- Input validation and structured JSON error responses
- Unit and PostgreSQL integration tests
- Docker Compose for local development
- GitHub Actions CI with formatting, vetting, race detection, tests, and build verification
- Health endpoint and graceful shutdown

## Tech Stack
- Go
- PostgreSQL
- database/sql with pgx
- Docker and Docker Compose
- GitHub Actions

## Requirements
- Go version specified in go.mod
- Docker Desktop or Docker Engine with Compose
- Git

## Run Locally

Start PostgreSQL and the API:

docker compose up --build -d
docker compose ps
docker compose logs -f api


API base URL: http://localhost:8080

Health check:

curl http://localhost:8080/healthz

## API Endpoints
### Create a payment

POST /v1/payments

Required header:

Idempotency-Key: unique-request-key

Example body:

{
  "amount_minor": 1250,
  "currency": "SGD"
}


A successful request returns HTTP 201 Created and a payment object.

Repeating the same request with the same key returns the existing payment. Reusing the key with different payment details returns HTTP 409 Conflict.

## Retrieve a payment

GET /v1/payments/{id}

Returns the payment record or an appropriate error response.

## Health

GET /healthz

Returns 200 OK when the HTTP server is running. This is a liveness check, not a database-readiness guarantee.

## Test

Unit tests:

go test ./...


PostgreSQL integration tests require DATABASE_URL to be set and PostgreSQL to be available.

Linux/macOS:

DATABASE_URL="postgres://app:app_password@localhost:5432/payments?sslmode=disable" go test -count=1 ./...


PowerShell:

$env:DATABASE_URL = "postgres://app:app_password@localhost:5432/payments?sslmode=disable"
go test -count=1 ./...

## Design Decisions

Store monetary amounts as integer minor units instead of floating-point values.
Use a unique database constraint on idempotency keys to handle concurrent duplicate requests.
Hash normalized request data to reject conflicting reuse of an idempotency key.
Separate HTTP handlers, business logic, and persistence.
Use PostgreSQL integration tests to verify persistence and idempotency behavior.

## Limitations

This project simulates successful payments and does not connect to a real payment provider. Authentication, authorization, rate limiting, payment-provider reconciliation, and a full payment state machine are outside the current scope.

## Future Improvements

Add a concurrency stress test for simultaneous requests using the same idempotency key.
Track whether responses are replays and return an accurate replay header.
Add database readiness checks and bounded retry behavior at startup.
Add metrics and structured request logging.
Add a payment state machine and transactional outbox for asynchronous events.
Add authentication, API rate limiting, and secret management.

## License

Choose and add a license before distributing this project publicly.