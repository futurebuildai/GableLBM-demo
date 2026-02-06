package order

import (
	"context"
	"fmt"
	"time"

	"github.com/gablelbm/gable/pkg/database"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Repository interface {
	CreateOrder(ctx context.Context, o *Order) error
	GetOrder(ctx context.Context, id uuid.UUID) (*Order, error)
	ListOrders(ctx context.Context) ([]Order, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status OrderStatus) error
}

type PostgresRepository struct {
	db *database.DB
}

func NewRepository(db *database.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) CreateOrder(ctx context.Context, o *Order) error {
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	if o.ID == uuid.Nil {
		o.ID = uuid.New()
	}
	now := time.Now()
	o.CreatedAt = now
	o.UpdatedAt = now
	o.Status = StatusDraft // Default to draft if not set

	// Insert Order
	queryOrder := `
		INSERT INTO orders (id, customer_id, quote_id, status, total_amount, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err = tx.Exec(ctx, queryOrder,
		o.ID, o.CustomerID, o.QuoteID, o.Status, o.TotalAmount, o.CreatedAt, o.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert order: %w", err)
	}

	// Insert Lines
	queryLine := `
		INSERT INTO order_lines (id, order_id, product_id, quantity, price_each, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	for i := range o.Lines {
		line := &o.Lines[i]
		if line.ID == uuid.Nil {
			line.ID = uuid.New()
		}
		line.OrderID = o.ID
		_, err = tx.Exec(ctx, queryLine,
			line.ID, line.OrderID, line.ProductID, line.Quantity, line.PriceEach, now,
		)
		if err != nil {
			return fmt.Errorf("failed to insert order line: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *PostgresRepository) GetOrder(ctx context.Context, id uuid.UUID) (*Order, error) {
	queryOrder := `
		SELECT id, customer_id, quote_id, status, total_amount, created_at, updated_at
		FROM orders WHERE id = $1
	`
	var o Order
	err := r.db.Pool.QueryRow(ctx, queryOrder, id).Scan(
		&o.ID, &o.CustomerID, &o.QuoteID, &o.Status, &o.TotalAmount, &o.CreatedAt, &o.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("order not found")
		}
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	// Get Lines
	queryLines := `
		SELECT id, order_id, product_id, quantity, price_each
		FROM order_lines WHERE order_id = $1
	`
	rows, err := r.db.Pool.Query(ctx, queryLines, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get order lines: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var l OrderLine
		if err := rows.Scan(&l.ID, &l.OrderID, &l.ProductID, &l.Quantity, &l.PriceEach); err != nil {
			return nil, fmt.Errorf("failed to scan order line: %w", err)
		}
		o.Lines = append(o.Lines, l)
	}

	return &o, nil
}

func (r *PostgresRepository) ListOrders(ctx context.Context) ([]Order, error) {
	query := `
		SELECT id, customer_id, quote_id, status, total_amount, created_at, updated_at
		FROM orders ORDER BY created_at DESC
	`
	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list orders: %w", err)
	}
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		var o Order
		if err := rows.Scan(
			&o.ID, &o.CustomerID, &o.QuoteID, &o.Status, &o.TotalAmount, &o.CreatedAt, &o.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}
		orders = append(orders, o)
	}
	return orders, nil
}

func (r *PostgresRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status OrderStatus) error {
	query := `UPDATE orders SET status = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.Pool.Exec(ctx, query, status, id)
	if err != nil {
		return fmt.Errorf("failed to update status: %w", err)
	}
	return nil
}
