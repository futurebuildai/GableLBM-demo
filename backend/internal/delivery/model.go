package delivery

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Status Enums
type VehicleType string

const (
	VehicleTypeBoxTruck VehicleType = "BOX_TRUCK"
	VehicleTypeFlatbed  VehicleType = "FLATBED"
	VehicleTypePickup   VehicleType = "PICKUP"
	VehicleTypeVan      VehicleType = "VAN"
	VehicleTypeCrane    VehicleType = "CRANE"
)

type DriverStatus string

const (
	DriverStatusActive   DriverStatus = "ACTIVE"
	DriverStatusInactive DriverStatus = "INACTIVE"
	DriverStatusOnLeave  DriverStatus = "ON_LEAVE"
)

type RouteStatus string

const (
	RouteStatusDraft     RouteStatus = "DRAFT"
	RouteStatusScheduled RouteStatus = "SCHEDULED"
	RouteStatusInTransit RouteStatus = "IN_TRANSIT"
	RouteStatusCompleted RouteStatus = "COMPLETED"
	RouteStatusCancelled RouteStatus = "CANCELLED"
)

type DeliveryStatus string

const (
	DeliveryStatusPending        DeliveryStatus = "PENDING"
	DeliveryStatusOutForDelivery DeliveryStatus = "OUT_FOR_DELIVERY"
	DeliveryStatusDelivered      DeliveryStatus = "DELIVERED"
	DeliveryStatusFailed         DeliveryStatus = "FAILED"
	DeliveryStatusPartial        DeliveryStatus = "PARTIAL"
)

// Domain Models

type Vehicle struct {
	ID                uuid.UUID   `json:"id" db:"id"`
	Name              string      `json:"name" db:"name"`
	VehicleType       VehicleType `json:"vehicle_type" db:"vehicle_type"`
	LicensePlate      string      `json:"license_plate" db:"license_plate"`
	CapacityWeightLbs *int        `json:"capacity_weight_lbs" db:"capacity_weight_lbs"`
	CreatedAt         time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time   `json:"updated_at" db:"updated_at"`
}

type Driver struct {
	ID            uuid.UUID    `json:"id" db:"id"`
	Name          string       `json:"name" db:"name"`
	LicenseNumber *string      `json:"license_number" db:"license_number"`
	Status        DriverStatus `json:"status" db:"status"`
	PhoneNumber   *string      `json:"phone_number" db:"phone_number"`
	CreatedAt     time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at" db:"updated_at"`
}

type Route struct {
	ID            uuid.UUID   `json:"id" db:"id"`
	VehicleID     uuid.UUID   `json:"vehicle_id" db:"vehicle_id"`
	DriverID      uuid.UUID   `json:"driver_id" db:"driver_id"`
	ScheduledDate time.Time   `json:"scheduled_date" db:"scheduled_date"` // YYYY-MM-DD
	Status        RouteStatus `json:"status" db:"status"`
	Notes         *string     `json:"notes" db:"notes"`
	CreatedAt     time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at" db:"updated_at"`

	// Joined fields
	VehicleName *string `json:"vehicle_name,omitempty" db:"vehicle_name"`
	DriverName  *string `json:"driver_name,omitempty" db:"driver_name"`
	StopCount   int     `json:"stop_count" db:"stop_count"`
}

type Delivery struct {
	ID           uuid.UUID      `json:"id" db:"id"`
	RouteID      uuid.UUID      `json:"route_id" db:"route_id"`
	OrderID      uuid.UUID      `json:"order_id" db:"order_id"`
	StopSequence int            `json:"stop_sequence" db:"stop_sequence"`
	Status       DeliveryStatus `json:"status" db:"status"`

	// POD
	PODProofURL  *string    `json:"pod_proof_url" db:"pod_proof_url"`
	PODSignedBy  *string    `json:"pod_signed_by" db:"pod_signed_by"`
	PODTimestamp *time.Time `json:"pod_timestamp" db:"pod_timestamp"`

	DeliveryInstructions *string `json:"delivery_instructions" db:"delivery_instructions"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`

	// Joined
	CustomerName *string `json:"customer_name,omitempty" db:"customer_name"`
	OrderNumber  *string `json:"order_number,omitempty" db:"order_number"`
	Address      *string `json:"address,omitempty" db:"address"` // From Order/Customer
}

// DTOs

type CreateVehicleRequest struct {
	Name              string      `json:"name"`
	VehicleType       VehicleType `json:"vehicle_type"`
	LicensePlate      string      `json:"license_plate"`
	CapacityWeightLbs *int        `json:"capacity_weight_lbs"`
}

type CreateDriverRequest struct {
	Name          string  `json:"name"`
	LicenseNumber *string `json:"license_number"`
	PhoneNumber   *string `json:"phone_number"`
}

type CreateRouteRequest struct {
	VehicleID     uuid.UUID `json:"vehicle_id"`
	DriverID      uuid.UUID `json:"driver_id"`
	ScheduledDate string    `json:"scheduled_date"` // "2023-10-27"
	Notes         *string   `json:"notes"`
}

type AssignOrderRequest struct {
	RouteID              uuid.UUID `json:"route_id"`
	OrderID              uuid.UUID `json:"order_id"`
	StopSequence         int       `json:"stop_sequence"`
	DeliveryInstructions *string   `json:"delivery_instructions"`
}

type UpdateDeliveryStatusRequest struct {
	Status      DeliveryStatus `json:"status"`
	PODProofURL *string        `json:"pod_proof_url"`
	PODSignedBy *string        `json:"pod_signed_by"`
}

// Interfaces

type Repository interface {
	// Fleet
	CreateVehicle(ctx context.Context, vehicle *Vehicle) error
	ListVehicles(ctx context.Context) ([]Vehicle, error)
	CreateDriver(ctx context.Context, driver *Driver) error
	ListDrivers(ctx context.Context) ([]Driver, error)

	// Routes
	CreateRoute(ctx context.Context, route *Route) error
	GetRoute(ctx context.Context, id uuid.UUID) (*Route, error)
	ListRoutes(ctx context.Context, date *time.Time) ([]Route, error)
	UpdateRouteStatus(ctx context.Context, id uuid.UUID, status RouteStatus) error

	// Deliveries
	CreateDelivery(ctx context.Context, delivery *Delivery) error
	GetDelivery(ctx context.Context, id uuid.UUID) (*Delivery, error)
	ListDeliveriesByRoute(ctx context.Context, routeID uuid.UUID) ([]Delivery, error)
	UpdateDeliveryStatus(ctx context.Context, id uuid.UUID, status DeliveryStatus, pod *PODUpdate) error
}

type PODUpdate struct {
	ProofURL string
	SignedBy string
	Time     time.Time
}
