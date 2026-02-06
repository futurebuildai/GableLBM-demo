package purchase_order

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// CreateFromSOLine creates or updates a DRAFT PO for the vendor of the special order item
func (s *Service) CreateFromSOLine(ctx context.Context, soLineId uuid.UUID, vendorId *uuid.UUID, description string, qty float64, cost float64) error {
	// 1. Find existing Draft PO for this Vendor
	var po *PurchaseOrder
	var err error

	if vendorId != nil {
		po, err = s.repo.GetDraftPOByVendor(ctx, vendorId)
	}

	// 2. If no draft found (or error), create new PO
	if po == nil || err != nil {
		newPO := &PurchaseOrder{
			ID:       uuid.New(),
			VendorID: vendorId,
			Status:   StatusDraft,
		}
		if err := s.repo.CreatePO(ctx, newPO); err != nil {
			return fmt.Errorf("failed to create PO: %w", err)
		}
		po = newPO
	}

	// 3. Add Line Item
	line := &PurchaseOrderLine{
		ID:             uuid.New(),
		POID:           po.ID,
		Description:    description,
		Quantity:       qty,
		Cost:           cost,
		LinkedSOLineID: &soLineId,
	}

	if err := s.repo.AddPOLine(ctx, line); err != nil {
		return fmt.Errorf("failed to add PO line: %w", err)
	}

	return nil
}
