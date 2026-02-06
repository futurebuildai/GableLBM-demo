package order

import (
	"context"
	"fmt"

	"github.com/gablelbm/gable/internal/inventory"
	"github.com/google/uuid"
)

type Service struct {
	repo         Repository
	inventorySvc *inventory.Service
}

func NewService(repo Repository, inventorySvc *inventory.Service) *Service {
	return &Service{
		repo:         repo,
		inventorySvc: inventorySvc,
	}
}

func (s *Service) CreateOrder(ctx context.Context, req CreateOrderRequest) (*Order, error) {
	// 1. Validate inputs
	if req.CustomerID == uuid.Nil {
		return nil, fmt.Errorf("customer_id is required")
	}

	if len(req.Lines) == 0 {
		return nil, fmt.Errorf("order must have at least one line item")
	}

	o := &Order{
		CustomerID: req.CustomerID,
		QuoteID:    req.QuoteID,
		Status:     StatusDraft,
	}

	var total float64
	for _, l := range req.Lines {
		if l.Quantity <= 0 {
			return nil, fmt.Errorf("line quantity must be positive")
		}
		if l.PriceEach < 0 {
			return nil, fmt.Errorf("line price must be non-negative")
		}

		line := OrderLine{
			ProductID: l.ProductID,
			Quantity:  l.Quantity,
			PriceEach: l.PriceEach,
		}
		o.Lines = append(o.Lines, line)
		total += l.Quantity * l.PriceEach
	}
	o.TotalAmount = total

	// 2. Persist Order (Draft)
	if err := s.repo.CreateOrder(ctx, o); err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	return o, nil
}

func (s *Service) ConfirmOrder(ctx context.Context, id uuid.UUID) error {
	// 1. Get Order
	o, err := s.repo.GetOrder(ctx, id)
	if err != nil {
		return err
	}

	if o.Status != StatusDraft {
		return fmt.Errorf("cannot confirm order in status %s", o.Status)
	}

	// 2. Allocate Inventory for each line
	// Note: Ideally this is transactional across services, but for MVP we do best effort or rely on saga pattern (later).
	// Since we are in same DB, we could share transaction... but services are decoupled.
	// For now, iterate and allocate. If failure, we are in stick state (half allocated).
	// TODO: Fix transaction boundary.
	for _, line := range o.Lines {
		if err := s.inventorySvc.Allocate(ctx, line.ProductID, line.Quantity); err != nil {
			return fmt.Errorf("failed to allocate stock for product %s: %w", line.ProductID, err)
		}
	}

	// 3. Update Status
	if err := s.repo.UpdateStatus(ctx, id, StatusConfirmed); err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}

	return nil
}

func (s *Service) ListOrders(ctx context.Context) ([]Order, error) {
	return s.repo.ListOrders(ctx)
}

func (s *Service) GetOrder(ctx context.Context, id uuid.UUID) (*Order, error) {
	return s.repo.GetOrder(ctx, id)
}
