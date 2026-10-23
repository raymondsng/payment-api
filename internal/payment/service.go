package payment

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
)

var (
	ErrInvalidAmount   = errors.New("amount_minor must be positive")
	ErrInvalidCurrency = errors.New("currency must be a 3-letter code")
	ErrMissingKey      = errors.New("Idempotency-Key header is required")
)

type Store interface {
	Create(
		context.Context,
		CreatePaymentRequest,
		string,
	) (Payment, error)

	GetByID(context.Context, string) (Payment, error)
}

type Service struct {
	store Store
}

func NewService(store Store) *Service {
	return &Service{store: store}
}

func (s *Service) Create(
	ctx context.Context,
	req CreatePaymentRequest,
	key string,
) (Payment, error) {
	key = strings.TrimSpace(key)
	if key == "" {
		return Payment{}, ErrMissingKey
	}
	if len(key) > 255 {
		return Payment{}, ErrMissingKey
	}
	if req.AmountMinor <= 0 {
		return Payment{}, ErrInvalidAmount
	}

	req.Currency = strings.ToUpper(strings.TrimSpace(req.Currency))
	if len(req.Currency) != 3 {
		return Payment{}, ErrInvalidCurrency
	}
	for _, ch := range req.Currency {
		if ch < 'A' || ch > 'Z' {
			return Payment{}, ErrInvalidCurrency
		}
	}

	req.IdempotencyKey = key

	// A deterministic fingerprint of the normalized business payload.
	sum := sha256.Sum256([]byte(
		fmt.Sprintf("%d:%s", req.AmountMinor, req.Currency),
	))
	requestHash := hex.EncodeToString(sum[:])

	return s.store.Create(ctx, req, requestHash)
}

func (s *Service) GetByID(
	ctx context.Context,
	id string,
) (Payment, error) {
	return s.store.GetByID(ctx, id)
}
