package product

import (
	"context"
	"fmt"

	"github.com/gablelbm/gable/pkg/database"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// Repository defines the interface for product data access
type Repository interface {
	CreateProduct(ctx context.Context, p *Product) error
	GetProduct(ctx context.Context, id uuid.UUID) (*Product, error)
	ListProducts(ctx context.Context) ([]Product, error)
	ListBelowReorder(ctx context.Context) ([]ReorderAlert, error)
}

// PostgresRepository implements Repository using pgx
type PostgresRepository struct {
	db *database.DB
}

// NewRepository creates a new PostgresRepository
func NewRepository(db *database.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

// CreateProduct inserts a new product into the database
func (r *PostgresRepository) CreateProduct(ctx context.Context, p *Product) error {
	query := `
		INSERT INTO products (sku, description, uom_primary, base_price, vendor, upc) 
		VALUES ($1, $2, $3, $4, $5, $6) 
		RETURNING id, created_at, updated_at`

	err := r.db.Pool.QueryRow(ctx, query, p.SKU, p.Description, p.UOMPrimary, p.BasePrice, p.Vendor, p.UPC).Scan(
		&p.ID,
		&p.CreatedAt,
		&p.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create product: %w", err)
	}

	return nil
}

// GetProduct retrieves a product by its ID
func (r *PostgresRepository) GetProduct(ctx context.Context, id uuid.UUID) (*Product, error) {
	query := `
		SELECT id, sku, description, uom_primary, base_price, vendor, upc,
		       COALESCE(reorder_point, 0), COALESCE(reorder_qty, 0),
		       created_at, updated_at
		FROM products
		WHERE id = $1`

	var p Product
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&p.ID,
		&p.SKU,
		&p.Description,
		&p.UOMPrimary,
		&p.BasePrice,
		&p.Vendor,
		&p.UPC,
		&p.ReorderPoint,
		&p.ReorderQty,
		&p.CreatedAt,
		&p.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("product not found")
		}
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	return &p, nil
}

// ListProducts retrieves all products
func (r *PostgresRepository) ListProducts(ctx context.Context) ([]Product, error) {
	query := `
		SELECT p.id, p.sku, p.description, p.uom_primary, p.base_price, p.vendor, p.upc,
		       COALESCE(p.reorder_point, 0), COALESCE(p.reorder_qty, 0),
		       p.created_at, p.updated_at,
		       COALESCE(SUM(i.quantity), 0) as total_quantity,
		       COALESCE(SUM(i.allocated), 0) as total_allocated
		FROM products p
		LEFT JOIN inventory i ON p.id = i.product_id
		GROUP BY p.id
		ORDER BY p.sku ASC`

	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list products: %w", err)
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var p Product
		if err := rows.Scan(
			&p.ID,
			&p.SKU,
			&p.Description,
			&p.UOMPrimary,
			&p.BasePrice,
			&p.Vendor,
			&p.UPC,
			&p.ReorderPoint,
			&p.ReorderQty,
			&p.CreatedAt,
			&p.UpdatedAt,
			&p.TotalQuantity,
			&p.TotalAllocated,
		); err != nil {
			return nil, fmt.Errorf("failed to scan product: %w", err)
		}
		products = append(products, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row error: %w", err)
	}

	return products, nil
}

// ListBelowReorder returns products whose current stock is below their reorder point
func (r *PostgresRepository) ListBelowReorder(ctx context.Context) ([]ReorderAlert, error) {
	query := `
		SELECT p.id, p.sku, p.description, p.vendor,
		       p.reorder_point, COALESCE(p.reorder_qty, 0),
		       COALESCE(SUM(i.quantity), 0) AS current_stock,
		       p.reorder_point - COALESCE(SUM(i.quantity), 0) AS deficit
		FROM products p
		LEFT JOIN inventory i ON p.id = i.product_id
		WHERE p.reorder_point > 0
		GROUP BY p.id
		HAVING COALESCE(SUM(i.quantity), 0) < p.reorder_point
		ORDER BY (p.reorder_point - COALESCE(SUM(i.quantity), 0)) DESC`

	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list reorder alerts: %w", err)
	}
	defer rows.Close()

	var alerts []ReorderAlert
	for rows.Next() {
		var a ReorderAlert
		if err := rows.Scan(
			&a.ProductID,
			&a.SKU,
			&a.Description,
			&a.Vendor,
			&a.ReorderPoint,
			&a.ReorderQty,
			&a.CurrentStock,
			&a.Deficit,
		); err != nil {
			return nil, fmt.Errorf("failed to scan reorder alert: %w", err)
		}
		alerts = append(alerts, a)
	}

	return alerts, nil
}
