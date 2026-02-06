package payment

import (
	"context"
	"fmt"
	"time"

	"github.com/gablelbm/gable/internal/invoice"
	"github.com/google/uuid"
)

type Service struct {
	repo        Repository
	invoiceRepo invoice.Repository
}

func NewService(repo Repository, invoiceRepo invoice.Repository) *Service {
	return &Service{
		repo:        repo,
		invoiceRepo: invoiceRepo,
	}
}

func (s *Service) ProcessPayment(ctx context.Context, invoiceID uuid.UUID, amount float64, method PaymentMethod, ref, notes string) (*Payment, error) {
	// 1. Get Invoice
	inv, err := s.invoiceRepo.GetInvoice(ctx, invoiceID)
	if err != nil {
		return nil, fmt.Errorf("invoice not found: %w", err)
	}

	// 2. Create Payment
	p := &Payment{
		InvoiceID: invoiceID,
		Amount:    amount,
		Method:    method,
		Reference: ref,
		Notes:     notes,
	}

	if err := s.repo.CreatePayment(ctx, p); err != nil {
		return nil, err
	}

	// 3. Calculate Totals and Update Status
	payments, err := s.repo.GetPaymentsByInvoiceID(ctx, invoiceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment history: %w", err)
	}

	var totalPaid float64
	for _, pay := range payments {
		totalPaid += pay.Amount
	}

	// Update Status Logic
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
		return nil, fmt.Errorf("failed to update invoice status: %w", err)
	}

	return p, nil
}

func (s *Service) GetHistory(ctx context.Context, invoiceID uuid.UUID) ([]Payment, error) {
	return s.repo.GetPaymentsByInvoiceID(ctx, invoiceID)
}
