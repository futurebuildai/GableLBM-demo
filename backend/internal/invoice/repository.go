package invoice

import (
	"context"
	"fmt"
	"time"

	"github.com/gablelbm/gable/pkg/database"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Repository interface {
	CreateInvoice(ctx context.Context, inv *Invoice) error
	GetInvoice(ctx context.Context, id uuid.UUID) (*Invoice, error)
	ListInvoices(ctx context.Context) ([]Invoice, error)
}

type PostgresRepository struct {
	db *database.DB
}

func NewRepository(db *database.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) CreateInvoice(ctx context.Context, inv *Invoice) error {
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	if inv.ID == uuid.Nil {
		inv.ID = uuid.New()
	}
	now := time.Now()
	inv.CreatedAt = now
	inv.UpdatedAt = now
	if inv.Status == "" {
		inv.Status = InvoiceStatusUnpaid
	}

	// Insert Invoice
	queryInv := `
		INSERT INTO invoices (id, order_id, customer_id, status, total_amount, due_date, paid_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err = tx.Exec(ctx, queryInv,
		inv.ID, inv.OrderID, inv.CustomerID, inv.Status, inv.TotalAmount, inv.DueDate, inv.PaidAt, inv.CreatedAt, inv.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert invoice: %w", err)
	}

	// Insert Lines
	queryLine := `
		INSERT INTO invoice_lines (id, invoice_id, product_id, quantity, price_each, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	for i := range inv.Lines {
		line := &inv.Lines[i]
		if line.ID == uuid.Nil {
			line.ID = uuid.New()
		}
		line.InvoiceID = inv.ID
		_, err = tx.Exec(ctx, queryLine,
			line.ID, line.InvoiceID, line.ProductID, line.Quantity, line.PriceEach, now,
		)
		if err != nil {
			return fmt.Errorf("failed to insert invoice line: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *PostgresRepository) GetInvoice(ctx context.Context, id uuid.UUID) (*Invoice, error) {
	queryInv := `
		SELECT id, order_id, customer_id, status, total_amount, due_date, paid_at, created_at, updated_at
		FROM invoices WHERE id = $1
	`
	var inv Invoice
	err := r.db.Pool.QueryRow(ctx, queryInv, id).Scan(
		&inv.ID, &inv.OrderID, &inv.CustomerID, &inv.Status, &inv.TotalAmount, &inv.DueDate, &inv.PaidAt, &inv.CreatedAt, &inv.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("invoice not found")
		}
		return nil, fmt.Errorf("failed to get invoice: %w", err)
	}

	// Get Lines
	queryLines := `
		SELECT id, invoice_id, product_id, quantity, price_each, created_at
		FROM invoice_lines WHERE invoice_id = $1
	`
	rows, err := r.db.Pool.Query(ctx, queryLines, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get invoice lines: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var l InvoiceLine
		if err := rows.Scan(&l.ID, &l.InvoiceID, &l.ProductID, &l.Quantity, &l.PriceEach, &l.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan invoice line: %w", err)
		}
		inv.Lines = append(inv.Lines, l)
	}

	return &inv, nil
}

func (r *PostgresRepository) ListInvoices(ctx context.Context) ([]Invoice, error) {
	query := `
		SELECT id, order_id, customer_id, status, total_amount, due_date, paid_at, created_at, updated_at
		FROM invoices ORDER BY created_at DESC
	`
	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list invoices: %w", err)
	}
	defer rows.Close()

	var invoices []Invoice
	for rows.Next() {
		var inv Invoice
		if err := rows.Scan(
			&inv.ID, &inv.OrderID, &inv.CustomerID, &inv.Status, &inv.TotalAmount, &inv.DueDate, &inv.PaidAt, &inv.CreatedAt, &inv.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan invoice: %w", err)
		}
		invoices = append(invoices, inv)
	}
	return invoices, nil
}
