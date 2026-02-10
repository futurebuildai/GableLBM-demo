package purchase_order

import (
	"context"
	"fmt"

	"github.com/gablelbm/gable/internal/domain"
	"github.com/gablelbm/gable/internal/edi"
	"github.com/gablelbm/gable/internal/inventory"
	"github.com/google/uuid"
)

type Service struct {
	repo         *Repository
	edi          *edi.Service
	inventorySvc *inventory.Service
}

func NewService(repo *Repository, ediSvc *edi.Service, inventorySvc *inventory.Service) *Service {
	return &Service{repo: repo, edi: ediSvc, inventorySvc: inventorySvc}
}

func (s *Service) ListPOs(ctx context.Context) ([]PurchaseOrder, error) {
	return s.repo.ListPOs(ctx)
}

func (s *Service) GetPO(ctx context.Context, id uuid.UUID) (*PurchaseOrder, error) {
	return s.repo.GetPO(ctx, id)
}

// CreateManualPO creates a new PO manually (not from a sales order special item)
func (s *Service) CreateManualPO(ctx context.Context, vendorID uuid.UUID, lines []struct {
	ProductID   string
	Description string
	Quantity    float64
	Cost        float64
}) (*PurchaseOrder, error) {
	po := &PurchaseOrder{
		ID:       uuid.New(),
		VendorID: &vendorID,
		Status:   StatusDraft,
	}

	if err := s.repo.CreatePO(ctx, po); err != nil {
		return nil, fmt.Errorf("failed to create PO: %w", err)
	}

	for _, l := range lines {
		var prodID *uuid.UUID
		if l.ProductID != "" {
			parsed, err := uuid.Parse(l.ProductID)
			if err == nil {
				prodID = &parsed
			}
		}
		line := &PurchaseOrderLine{
			ID:          uuid.New(),
			POID:        po.ID,
			ProductID:   prodID,
			Description: l.Description,
			Quantity:    l.Quantity,
			Cost:        l.Cost,
		}
		if err := s.repo.AddPOLine(ctx, line); err != nil {
			return nil, fmt.Errorf("failed to add PO line: %w", err)
		}
		po.Lines = append(po.Lines, *line)
	}

	return po, nil
}

// CreateManualPOFromHandler is a convenience wrapper using the handler's request types
func (s *Service) CreateManualPOFromHandler(ctx context.Context, vendorID uuid.UUID, lines []CreatePOLineInput) (*PurchaseOrder, error) {
	po := &PurchaseOrder{
		ID:       uuid.New(),
		VendorID: &vendorID,
		Status:   StatusDraft,
	}

	if err := s.repo.CreatePO(ctx, po); err != nil {
		return nil, fmt.Errorf("failed to create PO: %w", err)
	}

	for _, l := range lines {
		var prodID *uuid.UUID
		if l.ProductID != "" {
			parsed, err := uuid.Parse(l.ProductID)
			if err == nil {
				prodID = &parsed
			}
		}
		line := &PurchaseOrderLine{
			ID:          uuid.New(),
			POID:        po.ID,
			ProductID:   prodID,
			Description: l.Description,
			Quantity:    l.Quantity,
			Cost:        l.Cost,
		}
		if err := s.repo.AddPOLine(ctx, line); err != nil {
			return nil, fmt.Errorf("failed to add PO line: %w", err)
		}
		po.Lines = append(po.Lines, *line)
	}

	return po, nil
}

// CreatePOLineInput matches handler request shape
type CreatePOLineInput struct {
	ProductID   string
	Description string
	Quantity    float64
	Cost        float64
}

// CreateFromSOLine creates or updates a DRAFT PO for the vendor of the special order item
func (s *Service) CreateFromSOLine(ctx context.Context, soLineId uuid.UUID, vendorId *uuid.UUID, description string, qty float64, cost float64) error {
	var po *PurchaseOrder
	var err error

	if vendorId != nil {
		po, err = s.repo.GetDraftPOByVendor(ctx, vendorId)
	}

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

	if s.edi != nil {
		lines := make([]domain.POlineData, len(po.Lines))
		for i, l := range po.Lines {
			lines[i] = domain.POlineData{
				LineNumber: i + 1,
				Quantity:   l.Quantity,
				Cost:       l.Cost,
				ItemCode:   "UNKNOWN",
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

	po.Status = StatusSent
	return s.repo.UpdatePO(ctx, po)
}

// ReceivePO processes goods receipt against a PO, creating inventory entries
func (s *Service) ReceivePO(ctx context.Context, poID uuid.UUID, receivedLines []ReceiveLineInput) error {
	po, err := s.repo.GetPO(ctx, poID)
	if err != nil {
		return fmt.Errorf("PO not found: %w", err)
	}

	if po.Status != StatusSent && po.Status != StatusPartialReceive {
		return fmt.Errorf("PO must be in SENT or PARTIAL status to receive (current: %s)", po.Status)
	}

	lineMap := make(map[uuid.UUID]*PurchaseOrderLine)
	for i := range po.Lines {
		lineMap[po.Lines[i].ID] = &po.Lines[i]
	}

	allFullyReceived := true
	for _, rl := range receivedLines {
		lineID, err := uuid.Parse(rl.LineID)
		if err != nil {
			return fmt.Errorf("invalid line_id: %s", rl.LineID)
		}

		locationID, err := uuid.Parse(rl.LocationID)
		if err != nil {
			return fmt.Errorf("invalid location_id: %s", rl.LocationID)
		}

		poLine, ok := lineMap[lineID]
		if !ok {
			return fmt.Errorf("line %s not found on PO %s", rl.LineID, poID)
		}

		newQtyReceived := poLine.QtyReceived + rl.QtyReceived
		if err := s.repo.UpdateLineReceived(ctx, lineID, newQtyReceived); err != nil {
			return fmt.Errorf("failed to update received qty: %w", err)
		}

		// Create inventory if product_id is set
		if poLine.ProductID != nil && s.inventorySvc != nil {
			err := s.inventorySvc.AdjustStock(ctx, inventory.StockAdjustmentRequest{
				ProductID:  *poLine.ProductID,
				LocationID: &locationID,
				Quantity:   rl.QtyReceived,
				IsDelta:    true,
				Reason:     fmt.Sprintf("PO Receipt: %s", poID),
			})
			if err != nil {
				return fmt.Errorf("failed to create inventory for line %s: %w", rl.LineID, err)
			}
		}

		if newQtyReceived < poLine.Quantity {
			allFullyReceived = false
		}
	}

	// Check if any lines not in this receipt are still under-received
	for _, line := range po.Lines {
		found := false
		for _, rl := range receivedLines {
			if rl.LineID == line.ID.String() {
				found = true
				break
			}
		}
		if !found && line.QtyReceived < line.Quantity {
			allFullyReceived = false
		}
	}

	if allFullyReceived {
		po.Status = StatusReceived
	} else {
		po.Status = StatusPartialReceive
	}

	return s.repo.UpdatePO(ctx, po)
}

// ReceiveLineInput matches handler request shape
type ReceiveLineInput struct {
	LineID      string
	QtyReceived float64
	LocationID  string
}
