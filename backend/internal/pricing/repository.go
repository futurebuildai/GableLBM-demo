package pricing

import (
	"context"
	"fmt"
	"time"

	"github.com/gablelbm/gable/pkg/database"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Repository interface {
	GetContract(ctx context.Context, customerID, productID uuid.UUID) (*CustomerContract, error)
	CreateContract(ctx context.Context, c *CustomerContract) error
}

type PostgresRepository struct {
	db *database.DB
}

func NewRepository(db *database.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) GetContract(ctx context.Context, customerID, productID uuid.UUID) (*CustomerContract, error) {
	query := `
		SELECT id, customer_id, product_id, contract_price, created_at, updated_at
		FROM customer_contracts
		WHERE customer_id = $1 AND product_id = $2`

	var c CustomerContract
	err := r.db.Pool.QueryRow(ctx, query, customerID, productID).Scan(
		&c.ID, &c.CustomerID, &c.ProductID, &c.ContractPrice, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // No contract found
		}
		return nil, fmt.Errorf("failed to get contract: %w", err)
	}
	return &c, nil
}

func (r *PostgresRepository) CreateContract(ctx context.Context, c *CustomerContract) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	now := time.Now()
	c.CreatedAt = now
	c.UpdatedAt = now

	query := `
		INSERT INTO customer_contracts (id, customer_id, product_id, contract_price, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (customer_id, product_id) DO UPDATE 
		SET contract_price = EXCLUDED.contract_price, updated_at = EXCLUDED.updated_at`

	_, err := r.db.Pool.Exec(ctx, query,
		c.ID, c.CustomerID, c.ProductID, c.ContractPrice, c.CreatedAt, c.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create contract: %w", err)
	}
	return nil
}
