package payment

import (
	"context"
	"errors"
	"testing"
	"time"
)

type fakeStore struct {
	createFn func(
		context.Context,
		CreatePaymentRequest,
		string,
	) (Payment, error)

	getFn func(context.Context, string) (Payment, error)
}

func (f *fakeStore) Create(
	ctx context.Context,
	req CreatePaymentRequest,
	hash string,
) (Payment, error) {
	return f.createFn(ctx, req, hash)
}

func (f *fakeStore) GetByID(
	ctx context.Context,
	id string,
) (Payment, error) {
	return f.getFn(ctx, id)
}

func TestCreateNormalizesCurrency(t *testing.T) {
	store := &fakeStore{
		createFn: func(
			_ context.Context,
			req CreatePaymentRequest,
			hash string,
		) (Payment, error) {
			if req.Currency != "SGD" {
				t.Fatalf("currency = %q, want SGD", req.Currency)
			}
			if req.IdempotencyKey != "order-123" {
				t.Fatalf("unexpected key: %q", req.IdempotencyKey)
			}
			if hash == "" {
				t.Fatal("request hash must not be empty")
			}
			return Payment{
				ID:          "test-id",
				AmountMinor: req.AmountMinor,
				Currency:    req.Currency,
				Status:      "succeeded",
				CreatedAt:   time.Now(),
			}, nil
		},
	}

	svc := NewService(store)
	_, err := svc.Create(context.Background(),
		CreatePaymentRequest{
			AmountMinor: 1250,
			Currency:    "sgd",
		},
		"order-123",
	)
	if err != nil {
		t.Fatal(err)
	}
}

func TestCreateRejectsInvalidInput(t *testing.T) {
	tests := []struct {
		name string
		req  CreatePaymentRequest
		key  string
		want error
	}{
		{
			name: "missing key",
			req: CreatePaymentRequest{
				AmountMinor: 100,
				Currency:    "SGD",
			},
			want: ErrMissingKey,
		},
		{
			name: "zero amount",
			req: CreatePaymentRequest{
				AmountMinor: 0,
				Currency:    "SGD",
			},
			key:  "key-1",
			want: ErrInvalidAmount,
		},
		{
			name: "invalid currency",
			req: CreatePaymentRequest{
				AmountMinor: 100,
				Currency:    "SG",
			},
			key:  "key-1",
			want: ErrInvalidCurrency,
		},
	}

	svc := NewService(&fakeStore{})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.Create(
				context.Background(), tt.req, tt.key,
			)
			if !errors.Is(err, tt.want) {
				t.Fatalf("error = %v, want %v", err, tt.want)
			}
		})
	}
}
