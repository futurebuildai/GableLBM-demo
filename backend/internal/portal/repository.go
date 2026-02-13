package portal

import (
	"context"
	"fmt"
	"time"

	"github.com/gablelbm/gable/pkg/database"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// Repository defines data access for the portal module.
type Repository struct {
	db *database.DB
}

// NewRepository creates a new portal repository.
func NewRepository(db *database.DB) *Repository {
	return &Repository{db: db}
}

// GetCustomerUserByEmail fetches a customer user by email for login.
func (r *Repository) GetCustomerUserByEmail(ctx context.Context, email string) (*CustomerUser, error) {
	query := `
		SELECT id, customer_id, email, password_hash, name, role, created_at, updated_at
		FROM customer_users
		WHERE email = $1
	`
	var u CustomerUser
	err := r.db.Pool.QueryRow(ctx, query, email).Scan(
		&u.ID, &u.CustomerID, &u.Email, &u.PasswordHash, &u.Name, &u.Role,
		&u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get customer user: %w", err)
	}
	return &u, nil
}

// GetPortalConfig fetches the first (singleton) portal config row.
func (r *Repository) GetPortalConfig(ctx context.Context) (*PortalConfig, error) {
	query := `
		SELECT id, dealer_name, logo_url, primary_color, support_email, support_phone, created_at, updated_at
		FROM portal_config
		LIMIT 1
	`
	var cfg PortalConfig
	err := r.db.Pool.QueryRow(ctx, query).Scan(
		&cfg.ID, &cfg.DealerName, &cfg.LogoURL, &cfg.PrimaryColor,
		&cfg.SupportEmail, &cfg.SupportPhone, &cfg.CreatedAt, &cfg.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			// Return a sensible default if no config exists
			return &PortalConfig{
				DealerName:   "GableLBM",
				PrimaryColor: "#00FFA3",
			}, nil
		}
		return nil, fmt.Errorf("failed to get portal config: %w", err)
	}
	return &cfg, nil
}

// GetCustomerARSummary fetches balance, credit limit, and past-due amount.
func (r *Repository) GetCustomerARSummary(ctx context.Context, customerID uuid.UUID) (balance, creditLimit, pastDue float64, err error) {
	// Balance and credit limit from customers table
	custQuery := `SELECT balance_due, credit_limit FROM customers WHERE id = $1`
	err = r.db.Pool.QueryRow(ctx, custQuery, customerID).Scan(&balance, &creditLimit)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("failed to get customer AR: %w", err)
	}

	// Past due: sum of unpaid/overdue invoices past their due date
	pastDueQuery := `
		SELECT COALESCE(SUM(total_amount), 0)
		FROM invoices
		WHERE customer_id = $1
		  AND status IN ('UNPAID', 'OVERDUE')
		  AND due_date < NOW()
	`
	var pastDueCents int64
	err = r.db.Pool.QueryRow(ctx, pastDueQuery, customerID).Scan(&pastDueCents)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("failed to get past due: %w", err)
	}
	pastDue = float64(pastDueCents) / 100.0

	return balance, creditLimit, pastDue, nil
}

// ListOrdersByCustomer fetches orders with lines for a customer.
func (r *Repository) ListOrdersByCustomer(ctx context.Context, customerID uuid.UUID) ([]PortalOrderDTO, error) {
	query := `
		SELECT id, status, total_amount, created_at
		FROM orders
		WHERE customer_id = $1
		ORDER BY created_at DESC
		LIMIT 50
	`
	rows, err := r.db.Pool.Query(ctx, query, customerID)
	if err != nil {
		return nil, fmt.Errorf("failed to list orders: %w", err)
	}
	defer rows.Close()

	orders := make([]PortalOrderDTO, 0)
	for rows.Next() {
		var o PortalOrderDTO
		if err := rows.Scan(&o.ID, &o.Status, &o.TotalAmount, &o.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}
		o.Lines = make([]PortalLineDTO, 0) // Initialize empty for JSON []
		orders = append(orders, o)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	// Fetch lines for each order
	for i := range orders {
		lines, err := r.getOrderLines(ctx, orders[i].ID)
		if err != nil {
			return nil, err
		}
		orders[i].Lines = lines
	}

	return orders, nil
}

// getOrderLines fetches line items for an order.
func (r *Repository) getOrderLines(ctx context.Context, orderID uuid.UUID) ([]PortalLineDTO, error) {
	query := `
		SELECT ol.product_id, COALESCE(p.sku, ''), COALESCE(p.description, ''), ol.quantity, ol.price_each
		FROM order_lines ol
		LEFT JOIN products p ON ol.product_id = p.id
		WHERE ol.order_id = $1
	`
	rows, err := r.db.Pool.Query(ctx, query, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to list order lines: %w", err)
	}
	defer rows.Close()

	lines := make([]PortalLineDTO, 0)
	for rows.Next() {
		var l PortalLineDTO
		if err := rows.Scan(&l.ProductID, &l.ProductSKU, &l.ProductName, &l.Quantity, &l.PriceEach); err != nil {
			return nil, fmt.Errorf("failed to scan order line: %w", err)
		}
		lines = append(lines, l)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}
	return lines, nil
}

// GetOrderByIDAndCustomer fetches a single order scoped to a customer.
func (r *Repository) GetOrderByIDAndCustomer(ctx context.Context, orderID, customerID uuid.UUID) (*PortalOrderDTO, error) {
	query := `
		SELECT id, status, total_amount, created_at
		FROM orders
		WHERE id = $1 AND customer_id = $2
	`
	var o PortalOrderDTO
	err := r.db.Pool.QueryRow(ctx, query, orderID, customerID).Scan(
		&o.ID, &o.Status, &o.TotalAmount, &o.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("order not found")
		}
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	lines, err := r.getOrderLines(ctx, orderID)
	if err != nil {
		return nil, err
	}
	o.Lines = lines
	return &o, nil
}

// ListInvoicesByCustomer fetches invoices for a customer.
func (r *Repository) ListInvoicesByCustomer(ctx context.Context, customerID uuid.UUID) ([]PortalInvoiceDTO, error) {
	query := `
		SELECT id, order_id, status, total_amount, subtotal, tax_amount, payment_terms, due_date, paid_at, created_at
		FROM invoices
		WHERE customer_id = $1
		ORDER BY created_at DESC
		LIMIT 50
	`
	rows, err := r.db.Pool.Query(ctx, query, customerID)
	if err != nil {
		return nil, fmt.Errorf("failed to list invoices: %w", err)
	}
	defer rows.Close()

	invoices := make([]PortalInvoiceDTO, 0)
	for rows.Next() {
		var inv PortalInvoiceDTO
		if err := rows.Scan(
			&inv.ID, &inv.OrderID, &inv.Status, &inv.TotalAmount,
			&inv.Subtotal, &inv.TaxAmount, &inv.PaymentTerms,
			&inv.DueDate, &inv.PaidAt, &inv.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan invoice: %w", err)
		}
		inv.Lines = make([]PortalLineDTO, 0)
		invoices = append(invoices, inv)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return invoices, nil
}

// GetInvoiceByIDAndCustomer fetches a single invoice scoped to a customer.
func (r *Repository) GetInvoiceByIDAndCustomer(ctx context.Context, invoiceID, customerID uuid.UUID) (*PortalInvoiceDTO, error) {
	query := `
		SELECT id, order_id, status, total_amount, subtotal, tax_amount, payment_terms, due_date, paid_at, created_at
		FROM invoices
		WHERE id = $1 AND customer_id = $2
	`
	var inv PortalInvoiceDTO
	err := r.db.Pool.QueryRow(ctx, query, invoiceID, customerID).Scan(
		&inv.ID, &inv.OrderID, &inv.Status, &inv.TotalAmount,
		&inv.Subtotal, &inv.TaxAmount, &inv.PaymentTerms,
		&inv.DueDate, &inv.PaidAt, &inv.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("invoice not found")
		}
		return nil, fmt.Errorf("failed to get invoice: %w", err)
	}

	// Fetch invoice lines
	lineQuery := `
		SELECT il.product_id, COALESCE(p.sku, ''), COALESCE(p.description, ''), il.quantity, il.price_each
		FROM invoice_lines il
		LEFT JOIN products p ON il.product_id = p.id
		WHERE il.invoice_id = $1
	`
	lineRows, err := r.db.Pool.Query(ctx, lineQuery, invoiceID)
	if err != nil {
		return nil, fmt.Errorf("failed to list invoice lines: %w", err)
	}
	defer lineRows.Close()

	inv.Lines = make([]PortalLineDTO, 0)
	for lineRows.Next() {
		var l PortalLineDTO
		var priceEachCents int64
		if err := lineRows.Scan(&l.ProductID, &l.ProductSKU, &l.ProductName, &l.Quantity, &priceEachCents); err != nil {
			return nil, fmt.Errorf("failed to scan invoice line: %w", err)
		}
		l.PriceEach = float64(priceEachCents) / 100.0
		inv.Lines = append(inv.Lines, l)
	}
	if err := lineRows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return &inv, nil
}

// ListDeliveriesByCustomer fetches deliveries with POD info for a customer.
func (r *Repository) ListDeliveriesByCustomer(ctx context.Context, customerID uuid.UUID) ([]PortalDeliveryDTO, error) {
	query := `
		SELECT d.id, d.order_id, d.status, d.pod_proof_url, d.pod_signed_by, d.pod_timestamp,
		       d.created_at, o.id::text
		FROM deliveries d
		JOIN orders o ON d.order_id = o.id
		WHERE o.customer_id = $1
		ORDER BY d.created_at DESC
		LIMIT 50
	`
	rows, err := r.db.Pool.Query(ctx, query, customerID)
	if err != nil {
		return nil, fmt.Errorf("failed to list deliveries: %w", err)
	}
	defer rows.Close()

	deliveries := make([]PortalDeliveryDTO, 0)
	for rows.Next() {
		var d PortalDeliveryDTO
		if err := rows.Scan(
			&d.ID, &d.OrderID, &d.Status, &d.PODProofURL, &d.PODSignedBy,
			&d.PODTimestamp, &d.CreatedAt, &d.OrderNumber,
		); err != nil {
			return nil, fmt.Errorf("failed to scan delivery: %w", err)
		}
		deliveries = append(deliveries, d)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return deliveries, nil
}

// GetDeliveryByIDAndCustomer fetches a single delivery scoped to a customer.
func (r *Repository) GetDeliveryByIDAndCustomer(ctx context.Context, deliveryID, customerID uuid.UUID) (*PortalDeliveryDTO, error) {
	query := `
		SELECT d.id, d.order_id, d.status, d.pod_proof_url, d.pod_signed_by, d.pod_timestamp,
		       d.created_at, o.id::text
		FROM deliveries d
		JOIN orders o ON d.order_id = o.id
		WHERE d.id = $1 AND o.customer_id = $2
	`
	var d PortalDeliveryDTO
	err := r.db.Pool.QueryRow(ctx, query, deliveryID, customerID).Scan(
		&d.ID, &d.OrderID, &d.Status, &d.PODProofURL, &d.PODSignedBy,
		&d.PODTimestamp, &d.CreatedAt, &d.OrderNumber,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("delivery not found")
		}
		return nil, fmt.Errorf("failed to get delivery: %w", err)
	}
	return &d, nil
}

// CreateReorder duplicates order lines from a historical order into a new DRAFT order.
// Uses a transaction to ensure atomicity — partial failures roll back cleanly.
func (r *Repository) CreateReorder(ctx context.Context, customerID, sourceOrderID uuid.UUID) (uuid.UUID, error) {
	// Verify source order belongs to customer (outside tx — read-only check)
	var count int
	err := r.db.Pool.QueryRow(ctx, `SELECT COUNT(*) FROM orders WHERE id = $1 AND customer_id = $2`, sourceOrderID, customerID).Scan(&count)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to verify order ownership: %w", err)
	}
	if count == 0 {
		return uuid.Nil, fmt.Errorf("order not found")
	}

	// Begin transaction for all mutations
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx) // no-op after commit

	newOrderID := uuid.New()
	now := time.Now()

	// 1. Create new DRAFT order
	_, err = tx.Exec(ctx, `
		INSERT INTO orders (id, customer_id, status, total_amount, created_at, updated_at)
		SELECT $1, $2, 'DRAFT',
		       COALESCE(SUM(ol.quantity * ol.price_each), 0),
		       $3, $3
		FROM order_lines ol WHERE ol.order_id = $4
	`, newOrderID, customerID, now, sourceOrderID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create reorder: %w", err)
	}

	// 2. Copy lines with fresh UUIDs and current product prices
	_, err = tx.Exec(ctx, `
		INSERT INTO order_lines (id, order_id, product_id, quantity, price_each, is_special_order, vendor_id, special_order_cost)
		SELECT gen_random_uuid(), $1, ol.product_id, ol.quantity,
		       COALESCE(p.base_price, ol.price_each),
		       ol.is_special_order, ol.vendor_id, ol.special_order_cost
		FROM order_lines ol
		LEFT JOIN products p ON ol.product_id = p.id
		WHERE ol.order_id = $2
	`, newOrderID, sourceOrderID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to copy order lines: %w", err)
	}

	// 3. Recalculate total with current prices
	_, err = tx.Exec(ctx, `
		UPDATE orders SET total_amount = (
			SELECT COALESCE(SUM(quantity * price_each), 0)
			FROM order_lines WHERE order_id = $1
		) WHERE id = $1
	`, newOrderID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to update reorder total: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return uuid.Nil, fmt.Errorf("failed to commit reorder transaction: %w", err)
	}

	return newOrderID, nil
}
