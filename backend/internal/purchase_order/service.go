package purchase_order

import (
	"context"
	"fmt"

	"github.com/gablelbm/gable/internal/domain"
	"github.com/gablelbm/gable/internal/edi"
	"github.com/google/uuid"
)

type Service struct {
	repo *Repository
	edi  *edi.Service
}

func NewService(repo *Repository, ediSvc *edi.Service) *Service {
	return &Service{repo: repo, edi: ediSvc}
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

func (s *Service) SubmitPO(ctx context.Context, id uuid.UUID) error {
	po, err := s.repo.GetPO(ctx, id)
	if err != nil {
		return err
	}

	// Generate EDI
	if s.edi != nil {
		lines := make([]domain.POlineData, len(po.Lines))
		for i, l := range po.Lines {
			lines[i] = domain.POlineData{
				LineNumber: i + 1,
				Quantity:   l.Quantity,
				Cost:       l.Cost,
				ItemCode:   "UNKNOWN", // TODO: Add ItemCode to POLine model
			}
		}

		poData := domain.POData{
			ID:       po.ID,
			PONumber: po.ID.String(),
			VendorID: *po.VendorID,
			Lines:    lines,
		}

		if err := s.edi.SendExamplePO(ctx, poData); err != nil {
			return fmt.Errorf("failed to send EDI: %w", err)
		}
	}

	// Update Status
	po.Status = StatusSent
	return s.repo.UpdatePO(ctx, po)
}
