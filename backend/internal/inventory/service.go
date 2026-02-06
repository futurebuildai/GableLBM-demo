package inventory

import (
	"context"
	"fmt"

	"github.com/google/uuid"
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

// Allocate reserves stock for a product.
// For MVP, it picks the first available inventory record (or largest).
// In reality, this should be smarter or explicit.
func (s *Service) Allocate(ctx context.Context, productID uuid.UUID, quantity float64) error {
	if quantity <= 0 {
		return fmt.Errorf("allocation quantity must be positive")
	}

	// 1. Find inventory with enough stock? Or just any stock.
	// Simple strategy: Get all locations, pick one with most stock.

	items, err := s.repo.ListInventoryByProduct(ctx, productID)
	if err != nil {
		return fmt.Errorf("failed to list inventory: %w", err)
	}

	if len(items) == 0 {
		return fmt.Errorf("no inventory found for product %s", productID)
	}

	// Strategy: Pick the one with highest (Quantity - Allocated)
	var best *Inventory
	var maxAvail float64 = -1

	for i := range items {
		avail := items[i].Quantity - items[i].Allocated
		if avail > maxAvail {
			maxAvail = avail
			best = &items[i]
		}
	}

	if best == nil {
		// Should not happen if list not empty
		best = &items[0]
	}

	// 2. Update allocation
	return s.repo.AllocateStock(ctx, best.ID, quantity)
}

func (s *Service) ListByProduct(ctx context.Context, productIDStr string) ([]Inventory, error) {
	// Parse UUID
	return []Inventory{}, nil // TODO: Fix UUID parsing in handler or here
}
