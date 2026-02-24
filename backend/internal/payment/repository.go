package payment

import (
	"context"
	"fmt"
	"time"

	"github.com/gablelbm/gable/pkg/database"
	"github.com/google/uuid"
)

type Repository interface {
	CreatePayment(ctx context.Context, p *Payment) error
	GetPaymentsByInvoiceID(ctx context.Context, invoiceID uuid.UUID) ([]Payment, error)
}

type PostgresRepository struct {
	db *database.DB
}

func NewRepository(db *database.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) CreatePayment(ctx context.Context, p *Payment) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	p.CreatedAt = time.Now()

	query := `
		INSERT INTO payments (id, invoice_id, amount, method, reference, notes, created_at,
			gateway_tx_id, gateway_status, token_id, card_last4, card_brand, auth_code)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`
	// Convert Cents (int64) to Dollars (float64) for DB
	amountFloat := float64(p.Amount) / 100.0

	_, err := r.db.GetExecutor(ctx).Exec(ctx, query,
		p.ID, p.InvoiceID, amountFloat, p.Method, p.Reference, p.Notes, p.CreatedAt,
		nullIfEmpty(p.GatewayTxID), nullIfEmpty(p.GatewayStatus), nullIfEmpty(p.TokenID),
		nullIfEmpty(p.CardLast4), nullIfEmpty(p.CardBrand), nullIfEmpty(p.AuthCode),
	)
	if err != nil {
		return fmt.Errorf("failed to create payment: %w", err)
	}
	return nil
}

func (r *PostgresRepository) GetPaymentsByInvoiceID(ctx context.Context, invoiceID uuid.UUID) ([]Payment, error) {
	query := `
		SELECT id, invoice_id, amount, method, reference, notes, created_at,
			COALESCE(gateway_tx_id, '') as gateway_tx_id,
			COALESCE(gateway_status, '') as gateway_status,
			COALESCE(card_last4, '') as card_last4,
			COALESCE(card_brand, '') as card_brand,
			COALESCE(auth_code, '') as auth_code
		FROM payments
		WHERE invoice_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.GetExecutor(ctx).Query(ctx, query, invoiceID)
	if err != nil {
		return nil, fmt.Errorf("failed to list payments: %w", err)
	}
	defer rows.Close()

	var payments []Payment
	for rows.Next() {
		var p Payment
		var amountFloat float64
		if err := rows.Scan(
			&p.ID, &p.InvoiceID, &amountFloat, &p.Method, &p.Reference, &p.Notes, &p.CreatedAt,
			&p.GatewayTxID, &p.GatewayStatus, &p.CardLast4, &p.CardBrand, &p.AuthCode,
		); err != nil {
			return nil, fmt.Errorf("failed to scan payment: %w", err)
		}
		// Convert Dollars (float64) to Cents (int64)
		p.Amount = int64(amountFloat*100.0 + 0.5)
		payments = append(payments, p)
	}
	return payments, nil
}

// nullIfEmpty returns nil for empty strings (so DB stores NULL instead of empty).
func nullIfEmpty(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}
