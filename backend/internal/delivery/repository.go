package delivery

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gablelbm/gable/pkg/database"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type PostgresRepository struct {
	db *database.DB
}

func NewRepository(db *database.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

// Fleet

func (r *PostgresRepository) CreateVehicle(ctx context.Context, v *Vehicle) error {
	query := `
		INSERT INTO vehicles (name, vehicle_type, license_plate, capacity_weight_lbs)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`
	return r.db.GetExecutor(ctx).QueryRow(ctx, query,
		v.Name, v.VehicleType, v.LicensePlate, v.CapacityWeightLbs,
	).Scan(&v.ID, &v.CreatedAt, &v.UpdatedAt)
}

func (r *PostgresRepository) ListVehicles(ctx context.Context) ([]Vehicle, error) {
	query := `
		SELECT id, name, vehicle_type, license_plate, capacity_weight_lbs, created_at, updated_at
		FROM vehicles
		WHERE deleted_at IS NULL
		ORDER BY name ASC
	`
	rows, err := r.db.GetExecutor(ctx).Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var vehicles []Vehicle
	for rows.Next() {
		var v Vehicle
		if err := rows.Scan(
			&v.ID, &v.Name, &v.VehicleType, &v.LicensePlate, &v.CapacityWeightLbs, &v.CreatedAt, &v.UpdatedAt,
		); err != nil {
			return nil, err
		}
		vehicles = append(vehicles, v)
	}
	return vehicles, nil
}

func (r *PostgresRepository) GetVehicle(ctx context.Context, id uuid.UUID) (*Vehicle, error) {
	query := `
		SELECT id, name, vehicle_type, license_plate, capacity_weight_lbs, created_at, updated_at
		FROM vehicles
		WHERE id = $1 AND deleted_at IS NULL
	`
	var v Vehicle
	err := r.db.GetExecutor(ctx).QueryRow(ctx, query, id).Scan(
		&v.ID, &v.Name, &v.VehicleType, &v.LicensePlate, &v.CapacityWeightLbs, &v.CreatedAt, &v.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("vehicle not found")
		}
		return nil, err
	}
	return &v, nil
}

func (r *PostgresRepository) CreateDriver(ctx context.Context, d *Driver) error {
	query := `
		INSERT INTO drivers (name, license_number, phone_number, status)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`
	return r.db.GetExecutor(ctx).QueryRow(ctx, query,
		d.Name, d.LicenseNumber, d.PhoneNumber, d.Status,
	).Scan(&d.ID, &d.CreatedAt, &d.UpdatedAt)
}

func (r *PostgresRepository) ListDrivers(ctx context.Context) ([]Driver, error) {
	query := `
		SELECT id, name, license_number, phone_number, status, created_at, updated_at
		FROM drivers
		WHERE deleted_at IS NULL
		ORDER BY name ASC
	`
	rows, err := r.db.GetExecutor(ctx).Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var drivers []Driver
	for rows.Next() {
		var d Driver
		if err := rows.Scan(
			&d.ID, &d.Name, &d.LicenseNumber, &d.PhoneNumber, &d.Status, &d.CreatedAt, &d.UpdatedAt,
		); err != nil {
			return nil, err
		}
		drivers = append(drivers, d)
	}
	return drivers, nil
}

// Routes

func (r *PostgresRepository) CreateRoute(ctx context.Context, route *Route) error {
	query := `
		INSERT INTO delivery_routes (vehicle_id, driver_id, scheduled_date, status, notes)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`
	return r.db.GetExecutor(ctx).QueryRow(ctx, query,
		route.VehicleID, route.DriverID, route.ScheduledDate, route.Status, route.Notes,
	).Scan(&route.ID, &route.CreatedAt, &route.UpdatedAt)
}

func (r *PostgresRepository) GetRoute(ctx context.Context, id uuid.UUID) (*Route, error) {
	query := `
		SELECT r.id, r.vehicle_id, r.driver_id, r.scheduled_date, r.status, r.notes, r.created_at, r.updated_at,
		       v.name as vehicle_name, d.name as driver_name,
		       (SELECT COUNT(*) FROM deliveries WHERE route_id = r.id) as stop_count
		FROM delivery_routes r
		JOIN vehicles v ON r.vehicle_id = v.id
		JOIN drivers d ON r.driver_id = d.id
		WHERE r.id = $1
	`
	var route Route
	err := r.db.GetExecutor(ctx).QueryRow(ctx, query, id).Scan(
		&route.ID, &route.VehicleID, &route.DriverID, &route.ScheduledDate, &route.Status, &route.Notes, &route.CreatedAt, &route.UpdatedAt,
		&route.VehicleName, &route.DriverName, &route.StopCount,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("route not found")
		}
		return nil, err
	}
	return &route, nil
}

func (r *PostgresRepository) ListRoutes(ctx context.Context, date *time.Time, driverID *uuid.UUID) ([]Route, error) {
	query := `
		SELECT r.id, r.vehicle_id, r.driver_id, r.scheduled_date, r.status, r.notes, r.created_at, r.updated_at,
		       v.name as vehicle_name, d.name as driver_name,
		       (SELECT COUNT(*) FROM deliveries WHERE route_id = r.id) as stop_count
		FROM delivery_routes r
		JOIN vehicles v ON r.vehicle_id = v.id
		JOIN drivers d ON r.driver_id = d.id
	`
	args := []any{}
	whereClauses := []string{}

	if date != nil {
		args = append(args, *date)
		whereClauses = append(whereClauses, fmt.Sprintf("r.scheduled_date = $%d", len(args)))
	}
	if driverID != nil {
		args = append(args, *driverID)
		whereClauses = append(whereClauses, fmt.Sprintf("r.driver_id = $%d", len(args)))
	}

	if len(whereClauses) > 0 {
		query += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	query += " ORDER BY r.scheduled_date DESC, r.created_at DESC"

	rows, err := r.db.GetExecutor(ctx).Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var routes []Route
	for rows.Next() {
		var route Route
		if err := rows.Scan(
			&route.ID, &route.VehicleID, &route.DriverID, &route.ScheduledDate, &route.Status, &route.Notes, &route.CreatedAt, &route.UpdatedAt,
			&route.VehicleName, &route.DriverName, &route.StopCount,
		); err != nil {
			return nil, err
		}
		routes = append(routes, route)
	}
	return routes, nil
}

func (r *PostgresRepository) UpdateRouteStatus(ctx context.Context, id uuid.UUID, status RouteStatus) error {
	query := `UPDATE delivery_routes SET status = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.GetExecutor(ctx).Exec(ctx, query, status, id)
	return err
}

// Deliveries

func (r *PostgresRepository) CreateDelivery(ctx context.Context, d *Delivery) error {
	query := `
		INSERT INTO deliveries (route_id, order_id, stop_sequence, status, delivery_instructions)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`
	return r.db.GetExecutor(ctx).QueryRow(ctx, query,
		d.RouteID, d.OrderID, d.StopSequence, d.Status, d.DeliveryInstructions,
	).Scan(&d.ID, &d.CreatedAt, &d.UpdatedAt)
}

func (r *PostgresRepository) GetDelivery(ctx context.Context, id uuid.UUID) (*Delivery, error) {
	query := `
		SELECT d.id, d.route_id, d.order_id, d.stop_sequence, d.status, 
		       d.pod_proof_url, d.pod_signed_by, d.pod_timestamp, d.delivery_instructions, 
		       d.created_at, d.updated_at,
		       c.name as customer_name, o.order_number
		FROM deliveries d
		JOIN orders o ON d.order_id = o.id
		JOIN customers c ON o.customer_id = c.id
		WHERE d.id = $1
	`
	var d Delivery
	err := r.db.GetExecutor(ctx).QueryRow(ctx, query, id).Scan(
		&d.ID, &d.RouteID, &d.OrderID, &d.StopSequence, &d.Status,
		&d.PODProofURL, &d.PODSignedBy, &d.PODTimestamp, &d.DeliveryInstructions,
		&d.CreatedAt, &d.UpdatedAt,
		&d.CustomerName, &d.OrderNumber,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("delivery not found")
		}
		return nil, err
	}
	return &d, nil
}

func (r *PostgresRepository) ListDeliveriesByRoute(ctx context.Context, routeID uuid.UUID) ([]Delivery, error) {
	query := `
		SELECT d.id, d.route_id, d.order_id, d.stop_sequence, d.status, 
		       d.pod_proof_url, d.pod_signed_by, d.pod_timestamp, d.delivery_instructions, 
		       d.created_at, d.updated_at,
		       c.name as customer_name, o.order_number,
		       -- Assuming customer address is what we want here, or job site address
		       c.billing_address_line1 || ', ' || c.billing_address_city as address
		FROM deliveries d
		JOIN orders o ON d.order_id = o.id
		JOIN customers c ON o.customer_id = c.id
		WHERE d.route_id = $1
		ORDER BY d.stop_sequence ASC
	`
	rows, err := r.db.GetExecutor(ctx).Query(ctx, query, routeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var deliveries []Delivery
	for rows.Next() {
		var d Delivery
		if err := rows.Scan(
			&d.ID, &d.RouteID, &d.OrderID, &d.StopSequence, &d.Status,
			&d.PODProofURL, &d.PODSignedBy, &d.PODTimestamp, &d.DeliveryInstructions,
			&d.CreatedAt, &d.UpdatedAt,
			&d.CustomerName, &d.OrderNumber, &d.Address,
		); err != nil {
			return nil, err
		}
		deliveries = append(deliveries, d)
	}
	return deliveries, nil
}

func (r *PostgresRepository) UpdateDeliveryStatus(ctx context.Context, id uuid.UUID, status DeliveryStatus, pod *PODUpdate) error {
	if pod != nil {
		query := `
			UPDATE deliveries 
			SET status = $1, pod_proof_url = $2, pod_signed_by = $3, pod_timestamp = $4, updated_at = NOW() 
			WHERE id = $5
		`
		_, err := r.db.GetExecutor(ctx).Exec(ctx, query, status, pod.ProofURL, pod.SignedBy, pod.Time, id)
		return err
	}

	query := `UPDATE deliveries SET status = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.GetExecutor(ctx).Exec(ctx, query, status, id)
	return err
}

func (r *PostgresRepository) ReorderRouteDeliveries(ctx context.Context, routeID uuid.UUID, deliveryIDs []uuid.UUID) error {
	return r.db.RunInTx(ctx, func(ctx context.Context) error {
		for i, id := range deliveryIDs {
			query := `UPDATE deliveries SET stop_sequence = $1, updated_at = NOW() WHERE id = $2 AND route_id = $3`
			if _, err := r.db.GetExecutor(ctx).Exec(ctx, query, i+1, id, routeID); err != nil {
				return err
			}
		}
		return nil
	})
}

// Capacity

func (r *PostgresRepository) GetRouteLoadWeight(ctx context.Context, routeID uuid.UUID) (float64, error) {
	query := `
		SELECT COALESCE(SUM(
			COALESCE(p.weight_lbs, 0) * ol.quantity
		), 0)
		FROM deliveries d
		JOIN order_lines ol ON ol.order_id = d.order_id
		JOIN products p ON p.id = ol.product_id
		WHERE d.route_id = $1
	`
	var weight float64
	err := r.db.GetExecutor(ctx).QueryRow(ctx, query, routeID).Scan(&weight)
	return weight, err
}

func (r *PostgresRepository) GetOrderEstimatedWeight(ctx context.Context, orderID uuid.UUID) (float64, error) {
	query := `
		SELECT COALESCE(SUM(
			COALESCE(p.weight_lbs, 0) * ol.quantity
		), 0)
		FROM order_lines ol
		JOIN products p ON p.id = ol.product_id
		WHERE ol.order_id = $1
	`
	var weight float64
	err := r.db.GetExecutor(ctx).QueryRow(ctx, query, orderID).Scan(&weight)
	return weight, err
}
