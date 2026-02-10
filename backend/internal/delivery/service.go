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

func (s *Service) ListRoutes(ctx context.Context, dateStr *string, driverID *uuid.UUID) ([]Route, error) {
	var date *time.Time
	if dateStr != nil && *dateStr != "" {
		parsed, err := time.Parse("2006-01-02", *dateStr)
		if err != nil {
			return nil, fmt.Errorf("invalid date format")
		}
		date = &parsed
	}
	return s.repo.ListRoutes(ctx, date, driverID)
}

func (s *Service) DispatchRoute(ctx context.Context, id uuid.UUID) error {
	// TODO: Validate driver/vehicle availability?
	return s.repo.UpdateRouteStatus(ctx, id, RouteStatusInTransit)
}

// Delivery Management

func (s *Service) AssignOrderToRoute(ctx context.Context, req AssignOrderRequest) (*Delivery, *CapacityWarning, error) {
	// Verify route exists and get vehicle info
	route, err := s.repo.GetRoute(ctx, req.RouteID)
	if err != nil {
		return nil, nil, err
	}

	// Vehicle capacity validation
	var warning *CapacityWarning
	vehicle, err := s.repo.GetVehicle(ctx, route.VehicleID)
	if err == nil && vehicle.CapacityWeightLbs != nil && *vehicle.CapacityWeightLbs > 0 {
		currentLoad, _ := s.repo.GetRouteLoadWeight(ctx, req.RouteID)
		orderWeight, _ := s.repo.GetOrderEstimatedWeight(ctx, req.OrderID)
		totalAfter := currentLoad + orderWeight

		if totalAfter > float64(*vehicle.CapacityWeightLbs) {
			warning = &CapacityWarning{
				VehicleCapacity: float64(*vehicle.CapacityWeightLbs),
				CurrentLoad:     currentLoad,
				OrderWeight:     orderWeight,
				TotalAfter:      totalAfter,
			}
			// Warning only — still allow assignment (soft validation)
		}
	}

	d := &Delivery{
		RouteID:              req.RouteID,
		OrderID:              req.OrderID,
		StopSequence:         req.StopSequence,
		Status:               DeliveryStatusPending,
		DeliveryInstructions: req.DeliveryInstructions,
	}

	if err := s.repo.CreateDelivery(ctx, d); err != nil {
		return nil, nil, err
	}
	return d, warning, nil
}

func (s *Service) ListDeliveries(ctx context.Context, routeID uuid.UUID) ([]Delivery, error) {
	return s.repo.ListDeliveriesByRoute(ctx, routeID)
}

func (s *Service) GetDelivery(ctx context.Context, id uuid.UUID) (*Delivery, error) {
	return s.repo.GetDelivery(ctx, id)
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

func (s *Service) ReorderStops(ctx context.Context, routeID uuid.UUID, deliveryIDs []uuid.UUID) error {
	return s.repo.ReorderRouteDeliveries(ctx, routeID, deliveryIDs)
}
