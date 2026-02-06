package inventory

import (
	"context"
	"fmt"

	"github.com/gablelbm/gable/pkg/database"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Repository interface {
	GetInventory(ctx context.Context, productID uuid.UUID, locationID *uuid.UUID) (*Inventory, error)
	UpdateInventory(ctx context.Context, inv *Inventory) error
	CreateInventory(ctx context.Context, inv *Inventory) error
	ListInventoryByProduct(ctx context.Context, productID uuid.UUID) ([]Inventory, error)
}

type PostgresRepository struct {
	db *database.DB
}

func NewRepository(db *database.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) GetInventory(ctx context.Context, productID uuid.UUID, locationID *uuid.UUID) (*Inventory, error) {
	query := `
		SELECT id, product_id, location_id, location, quantity, updated_at
		FROM inventory
		WHERE product_id = $1 AND (($2::uuid IS NULL AND location_id IS NULL) OR location_id = $2)
	`
	// Note: The query above handles null location_id checking.
	// Simplifying assumption: location field (text) is ignored for lookup if location_id is used.

	// Actually, we should allow lookup by location_id strictly.
	// And for legacy, if location_id is null, we might check location text?
	// For this Sprint, we focus on location_id.

	var inv Inventory
	err := r.db.Pool.QueryRow(ctx, query, productID, locationID).Scan(
		&inv.ID,
		&inv.ProductID,
		&inv.LocationID,
		&inv.Location,
		&inv.Quantity,
		&inv.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // Return nil if not found, let service handle create logic
		}
		return nil, fmt.Errorf("failed to get inventory: %w", err)
	}

	return &inv, nil
}

func (r *PostgresRepository) CreateInventory(ctx context.Context, inv *Inventory) error {
	query := `
		INSERT INTO inventory (product_id, location_id, location, quantity)
		VALUES ($1, $2, $3, $4)
		RETURNING id, updated_at
	`
	// Ensure location text is populated if location_id is nil, or just empty string
	err := r.db.Pool.QueryRow(ctx, query, inv.ProductID, inv.LocationID, inv.Location, inv.Quantity).Scan(
		&inv.ID,
		&inv.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create inventory: %w", err)
	}
	return nil
}

func (r *PostgresRepository) UpdateInventory(ctx context.Context, inv *Inventory) error {
	query := `
		UPDATE inventory
		SET quantity = $1, updated_at = NOW()
		WHERE id = $2
		RETURNING updated_at
	`
	err := r.db.Pool.QueryRow(ctx, query, inv.Quantity, inv.ID).Scan(&inv.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to update inventory: %w", err)
	}
	return nil
}

func (r *PostgresRepository) ListInventoryByProduct(ctx context.Context, productID uuid.UUID) ([]Inventory, error) {
	query := `
        SELECT i.id, i.product_id, i.location_id, 
               COALESCE(l.path, i.location, '') as location_name, 
               i.quantity, i.updated_at
        FROM inventory i
        LEFT JOIN locations l ON i.location_id = l.id
        WHERE i.product_id = $1
    `
	rows, err := r.db.Pool.Query(ctx, query, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []Inventory
	for rows.Next() {
		var i Inventory
		if err := rows.Scan(&i.ID, &i.ProductID, &i.LocationID, &i.Location, &i.Quantity, &i.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	return items, nil
}
