package dashboard

import (
	"context"
	"fmt"
	"time"

	"github.com/gablelbm/gable/pkg/database"
)

// Repository defines the data access interface for dashboard metrics.
type Repository interface {
	GetDashboardSummary(ctx context.Context) (*DashboardSummary, error)
	GetInventoryAlerts(ctx context.Context, limit int) ([]InventoryAlert, error)
	GetTopCustomers(ctx context.Context, limit int, days int) ([]TopCustomer, error)
	GetOrderActivity(ctx context.Context, limit int) (*OrderActivity, error)
	GetRevenueTrend(ctx context.Context, days int) ([]RevenueTrendPoint, error)
}

// PostgresRepository implements Repository for Postgres.
type PostgresRepository struct {
	db *database.DB
}

// NewRepository creates a new dashboard repository.
func NewRepository(db *database.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

// GetDashboardSummary returns aggregate KPIs.
func (r *PostgresRepository) GetDashboardSummary(ctx context.Context) (*DashboardSummary, error) {
	summary := &DashboardSummary{}

	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	todayEnd := todayStart.Add(24 * time.Hour)
	yesterdayStart := todayStart.AddDate(0, 0, -1)

	// Today's revenue (cents)
	queryTodayRevenue := `
		SELECT COALESCE(SUM(amount), 0)
		FROM payments
		WHERE created_at >= $1 AND created_at < $2
	`
	if err := r.db.Pool.QueryRow(ctx, queryTodayRevenue, todayStart, todayEnd).Scan(&summary.TodayRevenue); err != nil {
		return nil, fmt.Errorf("failed to query today revenue: %w", err)
	}

	// Yesterday's revenue for comparison
	var yesterdayRevenue int64
	if err := r.db.Pool.QueryRow(ctx, queryTodayRevenue, yesterdayStart, todayStart).Scan(&yesterdayRevenue); err != nil {
		return nil, fmt.Errorf("failed to query yesterday revenue: %w", err)
	}
	if yesterdayRevenue > 0 {
		summary.TodayRevenueChange = float64(summary.TodayRevenue-yesterdayRevenue) / float64(yesterdayRevenue) * 100
	}

	// Active orders (not completed/cancelled)
	queryActiveOrders := `
		SELECT COUNT(*)
		FROM orders
		WHERE status IN ('PENDING', 'CONFIRMED', 'PROCESSING', 'READY', 'ALLOCATED')
	`
	if err := r.db.Pool.QueryRow(ctx, queryActiveOrders).Scan(&summary.ActiveOrders); err != nil {
		return nil, fmt.Errorf("failed to query active orders: %w", err)
	}

	// Pending dispatch (deliveries not dispatched)
	queryPendingDispatch := `
		SELECT COUNT(*)
		FROM deliveries
		WHERE status IN ('PENDING', 'ASSIGNED')
	`
	if err := r.db.Pool.QueryRow(ctx, queryPendingDispatch).Scan(&summary.PendingDispatch); err != nil {
		return nil, fmt.Errorf("failed to query pending dispatch: %w", err)
	}

	// Outstanding AR
	queryAR := `
		SELECT COALESCE(SUM(total_amount), 0), COUNT(*)
		FROM invoices
		WHERE status IN ('UNPAID', 'PARTIAL', 'OVERDUE')
	`
	if err := r.db.Pool.QueryRow(ctx, queryAR).Scan(&summary.OutstandingAR, &summary.OutstandingARCount); err != nil {
		return nil, fmt.Errorf("failed to query outstanding AR: %w", err)
	}

	return summary, nil
}

// GetInventoryAlerts returns products with low or zero stock.
func (r *PostgresRepository) GetInventoryAlerts(ctx context.Context, limit int) ([]InventoryAlert, error) {
	query := `
		SELECT 
			p.id, p.sku, p.name, 
			COALESCE(i.quantity, 0) as current_qty,
			COALESCE(p.reorder_point, 10) as reorder_qty,
			CASE 
				WHEN COALESCE(i.quantity, 0) = 0 THEN 'OUT_OF_STOCK'
				ELSE 'LOW_STOCK'
			END as alert_type,
			COALESCE(i.location_id::text, '') as location_id
		FROM products p
		LEFT JOIN inventory i ON p.id = i.product_id
		WHERE COALESCE(i.quantity, 0) <= COALESCE(p.reorder_point, 10)
		ORDER BY COALESCE(i.quantity, 0) ASC
		LIMIT $1
	`
	rows, err := r.db.Pool.Query(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query inventory alerts: %w", err)
	}
	defer rows.Close()

	alerts := make([]InventoryAlert, 0)
	for rows.Next() {
		var a InventoryAlert
		if err := rows.Scan(&a.ProductID, &a.SKU, &a.Name, &a.CurrentQty, &a.ReorderQty, &a.AlertType, &a.LocationID); err != nil {
			return nil, fmt.Errorf("failed to scan inventory alert: %w", err)
		}
		alerts = append(alerts, a)
	}

	return alerts, nil
}

// GetTopCustomers returns top customers by revenue in the given period.
func (r *PostgresRepository) GetTopCustomers(ctx context.Context, limit int, days int) ([]TopCustomer, error) {
	cutoff := time.Now().AddDate(0, 0, -days)
	query := `
		SELECT 
			c.id, c.name,
			COALESCE(SUM(inv.total_amount), 0) as total_revenue,
			COUNT(DISTINCT o.id) as order_count
		FROM customers c
		LEFT JOIN orders o ON c.id = o.customer_id AND o.created_at >= $1
		LEFT JOIN invoices inv ON o.id = inv.order_id
		GROUP BY c.id, c.name
		HAVING COALESCE(SUM(inv.total_amount), 0) > 0
		ORDER BY total_revenue DESC
		LIMIT $2
	`
	rows, err := r.db.Pool.Query(ctx, query, cutoff, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query top customers: %w", err)
	}
	defer rows.Close()

	customers := make([]TopCustomer, 0)
	for rows.Next() {
		var c TopCustomer
		if err := rows.Scan(&c.CustomerID, &c.CustomerName, &c.TotalRevenue, &c.OrderCount); err != nil {
			return nil, fmt.Errorf("failed to scan top customer: %w", err)
		}
		customers = append(customers, c)
	}

	return customers, nil
}

// GetOrderActivity returns recent orders and status breakdown.
func (r *PostgresRepository) GetOrderActivity(ctx context.Context, limit int) (*OrderActivity, error) {
	activity := &OrderActivity{
		RecentOrders:    make([]RecentOrder, 0),
		StatusBreakdown: make(map[string]int),
	}

	// Recent orders
	queryRecent := `
		SELECT o.id, c.name, o.total_amount, o.status, o.created_at
		FROM orders o
		LEFT JOIN customers c ON o.customer_id = c.id
		ORDER BY o.created_at DESC
		LIMIT $1
	`
	rows, err := r.db.Pool.Query(ctx, queryRecent, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query recent orders: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var o RecentOrder
		var customerName *string
		if err := rows.Scan(&o.OrderID, &customerName, &o.TotalAmount, &o.Status, &o.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan recent order: %w", err)
		}
		if customerName != nil {
			o.CustomerName = *customerName
		} else {
			o.CustomerName = "Walk-In"
		}
		activity.RecentOrders = append(activity.RecentOrders, o)
	}

	// Status breakdown
	queryStatus := `
		SELECT status, COUNT(*)
		FROM orders
		WHERE created_at >= NOW() - INTERVAL '30 days'
		GROUP BY status
	`
	statusRows, err := r.db.Pool.Query(ctx, queryStatus)
	if err != nil {
		return nil, fmt.Errorf("failed to query order status breakdown: %w", err)
	}
	defer statusRows.Close()

	for statusRows.Next() {
		var status string
		var count int
		if err := statusRows.Scan(&status, &count); err != nil {
			return nil, fmt.Errorf("failed to scan status breakdown: %w", err)
		}
		activity.StatusBreakdown[status] = count
	}

	return activity, nil
}

// GetRevenueTrend returns daily revenue for the last N days.
func (r *PostgresRepository) GetRevenueTrend(ctx context.Context, days int) ([]RevenueTrendPoint, error) {
	query := `
		SELECT 
			DATE(created_at) as date,
			COALESCE(SUM(amount), 0) as revenue
		FROM payments
		WHERE created_at >= NOW() - INTERVAL '%d days'
		GROUP BY DATE(created_at)
		ORDER BY date ASC
	`
	rows, err := r.db.Pool.Query(ctx, fmt.Sprintf(query, days))
	if err != nil {
		return nil, fmt.Errorf("failed to query revenue trend: %w", err)
	}
	defer rows.Close()

	trend := make([]RevenueTrendPoint, 0)
	for rows.Next() {
		var point RevenueTrendPoint
		var dateVal time.Time
		if err := rows.Scan(&dateVal, &point.Revenue); err != nil {
			return nil, fmt.Errorf("failed to scan revenue trend point: %w", err)
		}
		point.Date = dateVal.Format("2006-01-02")
		trend = append(trend, point)
	}

	return trend, nil
}
