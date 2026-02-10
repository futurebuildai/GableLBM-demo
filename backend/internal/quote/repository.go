package quote

import (
	"context"
	"fmt"
	"time"

	"github.com/gablelbm/gable/pkg/database"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Repository interface {
	CreateQuote(ctx context.Context, q *Quote) error
	GetQuote(ctx context.Context, id uuid.UUID) (*Quote, error)
	UpdateQuote(ctx context.Context, q *Quote) error
	ListQuotes(ctx context.Context) ([]Quote, error)
	ListQuotesByCustomer(ctx context.Context, customerID uuid.UUID) ([]Quote, error)
}

type PostgresRepository struct {
	db *database.DB
}

func NewRepository(db *database.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) CreateQuote(ctx context.Context, q *Quote) error {
	if q.ID == uuid.Nil {
		q.ID = uuid.New()
	}
	if q.State == "" {
		q.State = QuoteStateDraft
	}
	now := time.Now()
	q.CreatedAt = now
	q.UpdatedAt = now

	// Start transaction
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Insert Header
	queryHeader := `
		INSERT INTO quotes (
			id, customer_id, job_id, state, total_amount, expires_at, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err = tx.Exec(ctx, queryHeader,
		q.ID, q.CustomerID, q.JobID, q.State, q.TotalAmount, q.ExpiresAt, q.CreatedAt, q.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert quote header: %w", err)
	}

	// Insert Lines
	queryLine := `
		INSERT INTO quote_lines (
			id, quote_id, product_id, sku, description, quantity, uom, unit_price, line_total, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	for i := range q.Lines {
		line := &q.Lines[i]
		if line.ID == uuid.Nil {
			line.ID = uuid.New()
		}
		line.QuoteID = q.ID
		line.CreatedAt = now
		// Recalculate line total just in case? Or rely on service.
		// Service should ensure it's correct.

		_, err = tx.Exec(ctx, queryLine,
			line.ID, line.QuoteID, line.ProductID, line.SKU, line.Description,
			line.Quantity, line.UOM, line.UnitPrice, line.LineTotal, line.CreatedAt,
		)
		if err != nil {
			return fmt.Errorf("failed to insert quote line: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *PostgresRepository) GetQuote(ctx context.Context, id uuid.UUID) (*Quote, error) {
	q := &Quote{}

	// Get Header
	queryHeader := `
		SELECT id, customer_id, job_id, state, total_amount, expires_at, created_at, updated_at
		FROM quotes
		WHERE id = $1
	`
	err := r.db.Pool.QueryRow(ctx, queryHeader, id).Scan(
		&q.ID, &q.CustomerID, &q.JobID, &q.State, &q.TotalAmount, &q.ExpiresAt, &q.CreatedAt, &q.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("quote not found")
		}
		return nil, fmt.Errorf("failed to get quote header: %w", err)
	}

	// Get Lines
	queryLines := `
		SELECT id, quote_id, product_id, sku, description, quantity, uom, unit_price, line_total, created_at
		FROM quote_lines
		WHERE quote_id = $1
		ORDER BY created_at ASC
	`
	rows, err := r.db.Pool.Query(ctx, queryLines, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get quote lines: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var l QuoteLine
		if err := rows.Scan(
			&l.ID, &l.QuoteID, &l.ProductID, &l.SKU, &l.Description,
			&l.Quantity, &l.UOM, &l.UnitPrice, &l.LineTotal, &l.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan quote line: %w", err)
		}
		q.Lines = append(q.Lines, l)
	}

	return q, nil
}

func (r *PostgresRepository) UpdateQuote(ctx context.Context, q *Quote) error {
	// Full update for now - easier to just delete lines and re-insert or use upsert?
	// For simplicity in MVP: Update header, and if lines changed we might need smarter logic.
	// But "Quick Quote" is often created and finalized.
	// Let's assume UpdateQuote updates the header totals and state.

	q.UpdatedAt = time.Now()

	query := `
		UPDATE quotes 
		SET customer_id = $2, job_id = $3, state = $4, total_amount = $5, expires_at = $6, updated_at = $7
		WHERE id = $1
	`
	_, err := r.db.Pool.Exec(ctx, query,
		q.ID, q.CustomerID, q.JobID, q.State, q.TotalAmount, q.ExpiresAt, q.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to update quote: %w", err)
	}

	// TODO: Handle lines update if needed
	return nil
}

func (r *PostgresRepository) ListQuotes(ctx context.Context) ([]Quote, error) {
	query := `
		SELECT q.id, q.customer_id, COALESCE(c.name, ''), q.job_id, q.state, q.total_amount, q.expires_at, q.created_at, q.updated_at
		FROM quotes q
		LEFT JOIN customers c ON c.id = q.customer_id
		ORDER BY q.created_at DESC
	`
	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list quotes: %w", err)
	}
	defer rows.Close()

	var quotes []Quote
	for rows.Next() {
		var q Quote
		if err := rows.Scan(
			&q.ID, &q.CustomerID, &q.CustomerName, &q.JobID, &q.State, &q.TotalAmount, &q.ExpiresAt, &q.CreatedAt, &q.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan quote: %w", err)
		}
		quotes = append(quotes, q)
	}
	return quotes, nil
}

func (r *PostgresRepository) ListQuotesByCustomer(ctx context.Context, customerID uuid.UUID) ([]Quote, error) {
	query := `
		SELECT id, customer_id, job_id, state, total_amount, expires_at, created_at, updated_at
		FROM quotes
		WHERE customer_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.Pool.Query(ctx, query, customerID)
	if err != nil {
		return nil, fmt.Errorf("failed to list quotes: %w", err)
	}
	defer rows.Close()

	var quotes []Quote
	for rows.Next() {
		var q Quote
		if err := rows.Scan(
			&q.ID, &q.CustomerID, &q.JobID, &q.State, &q.TotalAmount, &q.ExpiresAt, &q.CreatedAt, &q.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan quote: %w", err)
		}
		quotes = append(quotes, q)
	}
	return quotes, nil
}
