package payment

import (
	"context"
	"fmt"
	"time"

	"github.com/gablelbm/gable/internal/invoice"
	"github.com/gablelbm/gable/pkg/database"
	"github.com/google/uuid"
)

type Service struct {
	db          *database.DB
	repo        Repository
	invoiceRepo invoice.Repository
}

func NewService(db *database.DB, repo Repository, invoiceRepo invoice.Repository) *Service {
	return &Service{
		db:          db,
		repo:        repo,
		invoiceRepo: invoiceRepo,
	}
}

func (s *Service) ProcessPayment(ctx context.Context, invoiceID uuid.UUID, amount float64, method PaymentMethod, ref, notes string) (*Payment, error) {
	var p *Payment
	// Convert input amount (dollars) to cents
	amountCents := int64(amount*100.0 + 0.5)

	err := s.db.RunInTx(ctx, func(ctx context.Context) error {
		// 1. Get Invoice
		inv, err := s.invoiceRepo.GetInvoice(ctx, invoiceID)
		if err != nil {
			return fmt.Errorf("invoice not found: %w", err)
		}

		// 2. Create Payment
		p = &Payment{
			InvoiceID: invoiceID,
			Amount:    amountCents,
			Method:    method,
			Reference: ref,
			Notes:     notes,
		}

		if err := s.repo.CreatePayment(ctx, p); err != nil {
			return err
		}

		// 3. Calculate Totals and Update Status
		payments, err := s.repo.GetPaymentsByInvoiceID(ctx, invoiceID)
		if err != nil {
			return fmt.Errorf("failed to get payment history: %w", err)
		}

		var totalPaid int64
		for _, pay := range payments {
			totalPaid += pay.Amount
		}

		// Update Status Logic
		// inv.TotalAmount is already int64 (Cents)
		if totalPaid >= inv.TotalAmount {
			inv.Status = invoice.InvoiceStatusPaid
			if inv.PaidAt == nil {
				now := time.Now()
				inv.PaidAt = &now
			}
		} else if totalPaid > 0 {
			inv.Status = invoice.InvoiceStatusPartial
			inv.PaidAt = nil
		} else {
			inv.Status = invoice.InvoiceStatusUnpaid
			inv.PaidAt = nil
		}

		// Always update to reflect latest state
		if err := s.invoiceRepo.UpdateInvoice(ctx, inv); err != nil {
			return fmt.Errorf("failed to update invoice status: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return p, nil
}

func (s *Service) GetHistory(ctx context.Context, invoiceID uuid.UUID) ([]Payment, error) {
	return s.repo.GetPaymentsByInvoiceID(ctx, invoiceID)
}
