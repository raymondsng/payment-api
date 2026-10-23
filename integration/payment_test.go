
package integration

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/raymondsng/payment-api/internal/payment"
)

func TestPaymentIdempotency(t *testing.T) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		t.Skip("set DATABASE_URL to run PostgreSQL integration tests")
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(
		context.Background(), 10*time.Second,
	)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		t.Fatal(err)
	}

	schema, err := os.ReadFile(filepath.Join(
		"..", "migrations", "001_create_payments.sql",
	))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := db.ExecContext(ctx, string(schema)); err != nil {
		t.Fatal(err)
	}

	repo := payment.NewRepository(db)
	svc := payment.NewService(repo)
	handler := payment.NewHandler(svc)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /v1/payments", handler.CreatePayment)
	mux.HandleFunc("GET /v1/payments/{id}", handler.GetPayment)

	key := "integration-" + time.Now().Format("20060102150405.000000000")
	body := []byte(`{"amount_minor":1250,"currency":"SGD"}`)

	create := func(payload []byte) *httptest.ResponseRecorder {
		req := httptest.NewRequest(
			http.MethodPost, "/v1/payments", bytes.NewReader(payload),
		)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Idempotency-Key", key)

		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		return rec
	}

	first := create(body)
	if first.Code != http.StatusCreated {
		t.Fatalf("first request: got %d: %s", first.Code, first.Body)
	}

	var p1 payment.Payment
	if err := json.Unmarshal(first.Body.Bytes(), &p1); err != nil {
		t.Fatal(err)
	}

	second := create(body)
	if second.Code != http.StatusCreated {
		t.Fatalf("retry: got %d: %s", second.Code, second.Body)
	}

	var p2 payment.Payment
	if err := json.Unmarshal(second.Body.Bytes(), &p2); err != nil {
		t.Fatal(err)
	}

	if p1.ID != p2.ID {
		t.Fatalf("retry created a different payment: %s vs %s", p1.ID, p2.ID)
	}

	conflict := create([]byte(
		`{"amount_minor":9999,"currency":"SGD"}`,
	))
	if conflict.Code != http.StatusConflict {
		t.Fatalf("different payload: got %d, want 409", conflict.Code)
	}

	getReq := httptest.NewRequest(
		http.MethodGet, "/v1/payments/"+p1.ID, nil,
	)
	getRec := httptest.NewRecorder()
	mux.ServeHTTP(getRec, getReq)

	if getRec.Code != http.StatusOK {
		t.Fatalf("GET payment: got %d: %s", getRec.Code, getRec.Body)
	}
}
