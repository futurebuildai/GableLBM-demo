package delivery

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// Fleet Management

func (s *Service) CreateVehicle(ctx context.Context, req CreateVehicleRequest) (*Vehicle, error) {
	v := &Vehicle{
		Name:              req.Name,
		VehicleType:       req.VehicleType,
		LicensePlate:      req.LicensePlate,
		CapacityWeightLbs: req.CapacityWeightLbs,
	}
	if err := s.repo.CreateVehicle(ctx, v); err != nil {
		return nil, err
	}
	return v, nil
}

func (s *Service) ListVehicles(ctx context.Context) ([]Vehicle, error) {
	return s.repo.ListVehicles(ctx)
}

func (s *Service) CreateDriver(ctx context.Context, req CreateDriverRequest) (*Driver, error) {
	d := &Driver{
		Name:          req.Name,
		LicenseNumber: req.LicenseNumber,
		PhoneNumber:   req.PhoneNumber,
		Status:        DriverStatusActive,
	}
	if err := s.repo.CreateDriver(ctx, d); err != nil {
		return nil, err
	}
	return d, nil
}

func (s *Service) ListDrivers(ctx context.Context) ([]Driver, error) {
	return s.repo.ListDrivers(ctx)
}

// Route Management

func (s *Service) CreateRoute(ctx context.Context, req CreateRouteRequest) (*Route, error) {
	date, err := time.Parse("2006-01-02", req.ScheduledDate)
	if err != nil {
		return nil, fmt.Errorf("invalid date format: %v", err)
	}

	route := &Route{
		VehicleID:     req.VehicleID,
		DriverID:      req.DriverID,
		ScheduledDate: date,
		Status:        RouteStatusDraft,
		Notes:         req.Notes,
	}

	if err := s.repo.CreateRoute(ctx, route); err != nil {
		return nil, err
	}
	return route, nil
}

func (s *Service) GetRoute(ctx context.Context, id uuid.UUID) (*Route, error) {
	return s.repo.GetRoute(ctx, id)
}

func (s *Service) ListRoutes(ctx context.Context, dateStr *string) ([]Route, error) {
	var date *time.Time
	if dateStr != nil && *dateStr != "" {
		parsed, err := time.Parse("2006-01-02", *dateStr)
		if err != nil {
			return nil, fmt.Errorf("invalid date format")
		}
		date = &parsed
	}
	return s.repo.ListRoutes(ctx, date)
}

func (s *Service) DispatchRoute(ctx context.Context, id uuid.UUID) error {
	// TODO: Validate driver/vehicle availability?
	return s.repo.UpdateRouteStatus(ctx, id, RouteStatusInTransit)
}

// Delivery Management

func (s *Service) AssignOrderToRoute(ctx context.Context, req AssignOrderRequest) (*Delivery, error) {
	// Verify route exists
	_, err := s.repo.GetRoute(ctx, req.RouteID)
	if err != nil {
		return nil, err
	}

	d := &Delivery{
		RouteID:              req.RouteID,
		OrderID:              req.OrderID,
		StopSequence:         req.StopSequence,
		Status:               DeliveryStatusPending,
		DeliveryInstructions: req.DeliveryInstructions,
	}

	if err := s.repo.CreateDelivery(ctx, d); err != nil {
		return nil, err
	}
	return d, nil
}

func (s *Service) ListDeliveries(ctx context.Context, routeID uuid.UUID) ([]Delivery, error) {
	return s.repo.ListDeliveriesByRoute(ctx, routeID)
}

func (s *Service) CompleteDelivery(ctx context.Context, id uuid.UUID, req UpdateDeliveryStatusRequest) error {
	var pod *PODUpdate
	if req.Status == DeliveryStatusDelivered || req.Status == DeliveryStatusPartial {
		if req.PODProofURL == nil || req.PODSignedBy == nil {
			return fmt.Errorf("POD proof required for delivery completion")
		}
		now := time.Now()
		pod = &PODUpdate{
			ProofURL: *req.PODProofURL,
			SignedBy: *req.PODSignedBy,
			Time:     now,
		}
	}

	return s.repo.UpdateDeliveryStatus(ctx, id, req.Status, pod)
}
