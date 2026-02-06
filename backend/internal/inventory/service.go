package inventory

import (
	"context"
	"fmt"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// AdjustStock handles receipt (Add) or cycle count (Set/Adjust)
func (s *Service) AdjustStock(ctx context.Context, req StockAdjustmentRequest) error {
	// 1. Check if inventory record exists
	inv, err := s.repo.GetInventory(ctx, req.ProductID, req.LocationID)
	if err != nil {
		return err
	}

	if inv == nil {
		// Create new record
		newQty := req.Quantity
		if req.IsDelta {
			// If delta on non-existent record, assume base 0
			newQty = req.Quantity
		}

		inv = &Inventory{
			ProductID:  req.ProductID,
			LocationID: req.LocationID,
			Quantity:   newQty,
			Location:   "", // Legacy field empty
		}
		return s.repo.CreateInventory(ctx, inv)
	}

	// Update existing
	if req.IsDelta {
		inv.Quantity += req.Quantity
	} else {
		inv.Quantity = req.Quantity
	}

	// Prevent negative stock?
	// Depending on business rule. For now allow negative (backorder/error) or strict?
	// Let's allow negative for now to reflect reality vs system, but warn?
	// User requested "Basic In/Out".

	return s.repo.UpdateInventory(ctx, inv)
}

func (s *Service) MoveStock(ctx context.Context, req StockMovementRequest) error {
	// Transactional logic ideally
	// 1. Remove from FromLocation
	// 2. Add to ToLocation

	// Subtract from source
	err := s.AdjustStock(ctx, StockAdjustmentRequest{
		ProductID:  req.ProductID,
		LocationID: req.FromLocationID,
		Quantity:   -req.Quantity,
		IsDelta:    true,
		Reason:     "Move Out: " + req.Reason,
	})
	if err != nil {
		return fmt.Errorf("failed to remove stock from source: %w", err)
	}

	// Add to dest
	err = s.AdjustStock(ctx, StockAdjustmentRequest{
		ProductID:  req.ProductID,
		LocationID: &req.ToLocationID,
		Quantity:   req.Quantity,
		IsDelta:    true,
		Reason:     "Move In: " + req.Reason,
	})
	if err != nil {
		// Ideally rollback source... but manual rollback here:
		// This is why we need DB transactions in service layer.
		// For MVP Sprint 03, we might accept risk or implement rudimentary transaction support.
		return fmt.Errorf("failed to add stock to destination: %w", err)
	}

	return nil
}

func (s *Service) ListByProduct(ctx context.Context, productIDStr string) ([]Inventory, error) {
	// Parse UUID
	return []Inventory{}, nil // TODO: Fix UUID parsing in handler or here
}
