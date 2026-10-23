package payment

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
)

var ErrIdempotencyConflict = errors.New(
	"idempotency key was already used with different payment details",
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(
	ctx context.Context,
	req CreatePaymentRequest,
	requestHash string,
) (Payment, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return Payment{}, err
	}
	defer tx.Rollback()

	var p Payment

	err = tx.QueryRowContext(ctx, `
		INSERT INTO payments (
			id, idempotency_key, request_hash,
			amount_minor, currency
		)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (idempotency_key) DO NOTHING
		RETURNING id, amount_minor, currency, status, created_at
	`,
		uuid.NewString(),
		req.IdempotencyKey,
		requestHash,
		req.AmountMinor,
		req.Currency,
	).Scan(
		&p.ID,
		&p.AmountMinor,
		&p.Currency,
		&p.Status,
		&p.CreatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		var storedHash string

		err = tx.QueryRowContext(ctx, `
			SELECT id, request_hash, amount_minor,
			       currency, status, created_at
			FROM payments
			WHERE idempotency_key = $1
		`, req.IdempotencyKey).Scan(
			&p.ID,
			&storedHash,
			&p.AmountMinor,
			&p.Currency,
			&p.Status,
			&p.CreatedAt,
		)
		if err != nil {
			return Payment{}, err
		}

		if storedHash != requestHash {
			return Payment{}, ErrIdempotencyConflict
		}
	} else if err != nil {
		return Payment{}, err
	}

	if err := tx.Commit(); err != nil {
		return Payment{}, err
	}

	return p, nil
}

func (r *Repository) GetByID(
	ctx context.Context,
	id string,
) (Payment, error) {
	var p Payment

	err := r.db.QueryRowContext(ctx, `
		SELECT id, amount_minor, currency, status, created_at
		FROM payments
		WHERE id = $1
	`, id).Scan(
		&p.ID,
		&p.AmountMinor,
		&p.Currency,
		&p.Status,
		&p.CreatedAt,
	)
	if err != nil {
		return Payment{}, err
	}

	return p, nil
}
