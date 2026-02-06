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
		INSERT INTO payments (id, invoice_id, amount, method, reference, notes, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	// Convert Cents (int64) to Dollars (float64) for DB
	amountFloat := float64(p.Amount) / 100.0

	_, err := r.db.GetExecutor(ctx).Exec(ctx, query,
		p.ID, p.InvoiceID, amountFloat, p.Method, p.Reference, p.Notes, p.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create payment: %w", err)
	}
	return nil
}

func (r *PostgresRepository) GetPaymentsByInvoiceID(ctx context.Context, invoiceID uuid.UUID) ([]Payment, error) {
	query := `
		SELECT id, invoice_id, amount, method, reference, notes, created_at
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
		); err != nil {
			return nil, fmt.Errorf("failed to scan payment: %w", err)
		}
		// Convert Dollars (float64) to Cents (int64)
		p.Amount = int64(amountFloat*100.0 + 0.5)
		payments = append(payments, p)
	}
	return payments, nil
}
