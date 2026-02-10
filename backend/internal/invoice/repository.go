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
	UpdateInvoice(ctx context.Context, inv *Invoice) error
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
	// Convert Cents -> Dollars
	totalAmountFloat := float64(inv.TotalAmount) / 100.0

	queryInv := `
		INSERT INTO invoices (id, order_id, customer_id, status, total_amount, due_date, paid_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err = tx.Exec(ctx, queryInv,
		inv.ID, inv.OrderID, inv.CustomerID, inv.Status, totalAmountFloat, inv.DueDate, inv.PaidAt, inv.CreatedAt, inv.UpdatedAt,
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
		// Convert PriceEach (Cents -> Dollars)
		priceEachFloat := float64(line.PriceEach) / 100.0

		_, err = tx.Exec(ctx, queryLine,
			line.ID, line.InvoiceID, line.ProductID, line.Quantity, priceEachFloat, now,
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
		SELECT i.id, i.order_id, i.customer_id, COALESCE(c.name, ''), i.status, i.total_amount, i.due_date, i.paid_at, i.created_at, i.updated_at
		FROM invoices i
		LEFT JOIN customers c ON c.id = i.customer_id
		WHERE i.id = $1
	`
	var inv Invoice
	var totalAmountFloat float64
	err := r.db.GetExecutor(ctx).QueryRow(ctx, queryInv, id).Scan(
		&inv.ID, &inv.OrderID, &inv.CustomerID, &inv.CustomerName, &inv.Status, &totalAmountFloat, &inv.DueDate, &inv.PaidAt, &inv.CreatedAt, &inv.UpdatedAt,
	)
	inv.TotalAmount = int64(totalAmountFloat*100.0 + 0.5)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("invoice not found")
		}
		return nil, fmt.Errorf("failed to get invoice: %w", err)
	}

	// Get Lines with product names
	queryLines := `
		SELECT il.id, il.invoice_id, il.product_id, COALESCE(p.sku, ''), COALESCE(p.description, ''), il.quantity, il.price_each, il.created_at
		FROM invoice_lines il
		LEFT JOIN products p ON p.id = il.product_id
		WHERE il.invoice_id = $1
	`
	rows, err := r.db.GetExecutor(ctx).Query(ctx, queryLines, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get invoice lines: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var l InvoiceLine
		var priceEachFloat float64
		if err := rows.Scan(&l.ID, &l.InvoiceID, &l.ProductID, &l.ProductSKU, &l.ProductName, &l.Quantity, &priceEachFloat, &l.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan invoice line: %w", err)
		}
		l.PriceEach = int64(priceEachFloat*100.0 + 0.5)
		inv.Lines = append(inv.Lines, l)
	}

	return &inv, nil
}

func (r *PostgresRepository) ListInvoices(ctx context.Context) ([]Invoice, error) {
	query := `
		SELECT i.id, i.order_id, i.customer_id, COALESCE(c.name, ''), i.status, i.total_amount, i.due_date, i.paid_at, i.created_at, i.updated_at
		FROM invoices i
		LEFT JOIN customers c ON c.id = i.customer_id
		ORDER BY i.created_at DESC
	`
	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list invoices: %w", err)
	}
	defer rows.Close()

	var invoices []Invoice
	for rows.Next() {
		var inv Invoice
		var totalAmountFloat float64
		if err := rows.Scan(
			&inv.ID, &inv.OrderID, &inv.CustomerID, &inv.CustomerName, &inv.Status, &totalAmountFloat, &inv.DueDate, &inv.PaidAt, &inv.CreatedAt, &inv.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan invoice: %w", err)
		}
		inv.TotalAmount = int64(totalAmountFloat*100.0 + 0.5)
		invoices = append(invoices, inv)
	}
	return invoices, nil
}

func (r *PostgresRepository) UpdateInvoice(ctx context.Context, inv *Invoice) error {
	inv.UpdatedAt = time.Now()
	query := `
		UPDATE invoices
		SET status = $1, paid_at = $2, updated_at = $3
		WHERE id = $4
	`
	_, err := r.db.GetExecutor(ctx).Exec(ctx, query, inv.Status, inv.PaidAt, inv.UpdatedAt, inv.ID)
	if err != nil {
		return fmt.Errorf("failed to update invoice: %w", err)
	}
	return nil
}
