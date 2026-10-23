# Payment Transaction API

A full-stack payment management demo featuring a Go REST API, PostgreSQL database, and React frontend. This project demonstrates reliable payment record creation, idempotency, conflict detection, and frontend-to-backend integration.

Note: This project simulates payment transactions. It does not process real payments or connect to a payment provider.

## Features
# Backend
- REST API built with Go
- PostgreSQL persistence
- Create and retrieve payment records
- Idempotency keys to prevent duplicate payment records
- Request fingerprinting to detect conflicting retries
- HTTP 409 Conflict when an idempotency key is reused with different payment details
- Input validation and JSON error responses
- Unit and PostgreSQL integration tests
- Docker Compose for local development
- GitHub Actions CI

# Frontend
- React interface built with Vite
- Payment creation form
- Payment lookup by ID
- API response display
- User-friendly error messages
- Conflict handling for HTTP 409 responses
- Responsive layout for desktop and mobile

## Tech Stack
- Layer	Technology
- Frontend	React, JavaScript, Vite, CSS
- Backend	Go
- Database	PostgreSQL
- API	REST, JSON, HTTP
- Development	VS Code, Git, Docker Compose
- CI	GitHub Actions

## Project Structure
payment-api/
├── frontend/
│   ├── src/
│   │   ├── api/
│   │   │   └── payments.js
│   │   ├── App.jsx
│   │   ├── App.css
│   │   ├── index.css
│   │   └── main.jsx
│   ├── vite.config.js
│   └── package.json
├── integration/
├── internal/
│   └── payment/
├── migrations/
├── .github/
├── Dockerfile
├── compose.yaml
├── go.mod
├── main.go
└── README.md

## Prerequisites

Install the following before running the project:

Git
Go — use the version specified in go.mod
Node.js — an LTS release
Docker Desktop
Visual Studio Code

## Getting Started

1. Clone the repository
git clone https://github.com/raymondsng/payment-api.git
cd payment-api


Open the project in VS Code:

code .

2. Start the backend

Make sure Docker Desktop is running.

From the repository root, run:

docker compose up --build -d


Check the containers:

docker compose ps


View the API logs if necessary:

docker compose logs -f api


The API should be available at:

http://localhost:8080

Check the health endpoint:

curl -i http://localhost:8080/healthz


The health endpoint checks whether the HTTP server is running. It does not guarantee that the database is ready.

3. Start the React frontend

Open a second terminal in VS Code.

From the repository root:

cd frontend
npm install
npm run dev


Open the local URL printed by Vite, usually:

http://localhost:5173

Keep the backend and frontend running in their respective terminals.

4. Frontend API configuration

The Vite development server proxies frontend requests to the Go API.

For example:

Frontend request	Backend endpoint
/api/v1/payments	http://localhost:8080/v1/payments
/api/v1/payments/{id}	http://localhost:8080/v1/payments/{id}
/api/healthz	http://localhost:8080/healthz

The proxy removes the /api prefix before forwarding the request.

## API Documentation

# Create a payment

POST /v1/payments

Required header:

Content-Type: application/json
Idempotency-Key: unique-request-key


Example request:

{
  "amount_minor": 1250,
  "currency": "SGD"
}


amount_minor represents the monetary amount in minor units. For example, 1250 represents SGD 12.50.

Expected behavior:

A new valid request returns 201 Created.
Repeating the same request with the same idempotency key returns the existing payment.
Reusing the same key with different payment details returns 409 Conflict.
Invalid requests return an appropriate error response.

# Retrieve a payment

GET /v1/payments/{id}

Retrieves a payment by its ID.

The API returns the payment record or an appropriate error response.

# Health check

GET /healthz

Returns 200 OK when the HTTP server is running.

## Testing the API

# Test payment creation

Run this command on macOS or Linux:

curl -i -X POST http://localhost:8080/v1/payments \
  -H "Content-Type: application/json" \
  -H "Idempotency-Key: demo-payment-001" \
  -d '{"amount_minor":1250,"currency":"SGD"}'

# Test idempotency

Repeat the request with the same key and identical body. Both requests should refer to the same payment.

# Test conflict detection

Reuse the key with a different amount:

curl -i -X POST http://localhost:8080/v1/payments \
  -H "Content-Type: application/json" \
  -H "Idempotency-Key: demo-payment-001" \
  -d '{"amount_minor":2000,"currency":"SGD"}'


The backend should return 409 Conflict if the key was previously used for different payment details.

These commands test the backend directly. The frontend conflict message can be tested by making a frontend request that receives the same HTTP response.

## Running Go Tests

Run unit tests from the repository root:

go test ./...


PostgreSQL integration tests require a running PostgreSQL instance and a valid DATABASE_URL.

Example for macOS or Linux:

DATABASE_URL="postgres://app:app_password@localhost:5432/payments?sslmode=disable" go test -count=1 ./...


Use the database credentials configured in your local environment if they differ from this example.

## Troubleshooting

# Docker or backend issues

Check container status:

docker compose ps


Check backend logs:

docker compose logs -f api

# Frontend dependency issues

From the frontend directory:

npm install
npm run dev

## Frontend cannot reach the API
- Confirm the backend is available at http://localhost:8080.
- Confirm Vite is running on the frontend development port.
- Check the browser's Developer Tools and Network tab.
- Verify the proxy configuration in frontend/vite.config.js.

## Design Decisions
- Store monetary amounts as integer minor units rather than floating-point values.
- Use database uniqueness constraints to protect against duplicate idempotency keys.
- Fingerprint normalized request data to detect conflicting key reuse.
- Separate HTTP handlers, business logic, persistence, and frontend API calls.
- Use PostgreSQL integration tests to verify persistence and idempotency behavior.

## Limitations
- Payments are simulated; no actual money is transferred.
- The project does not integrate with a real payment provider.
- Authentication, authorization, and rate limiting are not currently part of the documented scope.
- The frontend is a development interface, not a production payment administration system.
- A payment history page requires an appropriate backend listing endpoint if one is not already implemented.

## Future Improvements
- Add a payment history endpoint and transaction table.
- Improve frontend validation and automated UI testing.
- Preserve idempotency keys for safe retries.
- Add authentication, authorization, and rate limiting.
- Add database readiness checks, metrics, and structured logging.
- Integrate a payment provider and implement a complete payment lifecycle.
- Add production deployment configuration and security controls.

## License

Choose and add a license before distributing this project publicly.