package purchase_order

import (
	"context"
	"fmt"

	"github.com/gablelbm/gable/pkg/database"
	"github.com/google/uuid"
)

type Repository struct {
	db *database.DB
}

func NewRepository(db *database.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreatePO(ctx context.Context, po *PurchaseOrder) error {
	query := `
		INSERT INTO purchase_orders (id, vendor_id, status)
		VALUES ($1, $2, $3)
		RETURNING created_at, updated_at
	`
	return r.db.Pool.QueryRow(ctx, query,
		po.ID,
		po.VendorID,
		po.Status,
	).Scan(&po.CreatedAt, &po.UpdatedAt)
}

func (r *Repository) AddPOLine(ctx context.Context, line *PurchaseOrderLine) error {
	query := `
		INSERT INTO purchase_order_lines (id, po_id, description, quantity, cost, linked_so_line_id)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.Pool.Exec(ctx, query,
		line.ID,
		line.POID,
		line.Description,
		line.Quantity,
		line.Cost,
		line.LinkedSOLineID,
	)
	return err
}

func (r *Repository) GetDraftPOByVendor(ctx context.Context, vendorID *uuid.UUID) (*PurchaseOrder, error) {
	// If vendorID is nil, finding a "generic" draft PO is risky.
	// For now, we assume Special Orders usually have a vendor.
	// If vendorID is nil, we might create a new generic PO.

	if vendorID == nil {
		return nil, fmt.Errorf("vendor_id required lookup")
	}

	query := `
		SELECT id, vendor_id, status, created_at, updated_at
		FROM purchase_orders
		WHERE vendor_id = $1 AND status = 'DRAFT'
		LIMIT 1
	`
	var po PurchaseOrder
	err := r.db.Pool.QueryRow(ctx, query, vendorID).Scan(
		&po.ID,
		&po.VendorID,
		&po.Status,
		&po.CreatedAt,
		&po.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &po, nil
}

func (r *Repository) GetPO(ctx context.Context, id uuid.UUID) (*PurchaseOrder, error) {
	// Get PO Header
	// Note: We need 'scany' or manual scan. Manual scan for now.
	query := `SELECT id, vendor_id, status, created_at, updated_at FROM purchase_orders WHERE id = $1`
	var po PurchaseOrder
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&po.ID,
		&po.VendorID,
		&po.Status,
		&po.CreatedAt,
		&po.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get PO header: %w", err)
	}

	// Get Lines
	linesQuery := `
		SELECT id, po_id, description, quantity, cost, linked_so_line_id 
		FROM purchase_order_lines 
		WHERE po_id = $1
	`
	rows, err := r.db.Pool.Query(ctx, linesQuery, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get PO lines: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var line PurchaseOrderLine
		if err := rows.Scan(
			&line.ID,
			&line.POID,
			&line.Description,
			&line.Quantity,
			&line.Cost,
			&line.LinkedSOLineID,
		); err != nil {
			return nil, fmt.Errorf("failed to scan PO line: %w", err)
		}
		po.Lines = append(po.Lines, line)
	}

	return &po, nil
}

func (r *Repository) UpdatePO(ctx context.Context, po *PurchaseOrder) error {
	query := `UPDATE purchase_orders SET status = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.Pool.Exec(ctx, query, po.Status, po.ID)
	return err
}
