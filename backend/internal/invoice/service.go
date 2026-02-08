package invoice

import (
	"context"
	"fmt"

	"github.com/gablelbm/gable/internal/account"
	"github.com/gablelbm/gable/internal/gl"
	"github.com/google/uuid"
)

type Service struct {
	repo    Repository
	gl      *gl.Service
	account account.Service
}

func NewService(repo Repository, glService *gl.Service, accountService account.Service) *Service {
	return &Service{repo: repo, gl: glService, account: accountService}
}

func (s *Service) CreateInvoice(ctx context.Context, inv *Invoice) error {
	if len(inv.Lines) == 0 {
		return fmt.Errorf("invoice must have lines")
	}
	// Initial status
	if inv.Status == "" {
		inv.Status = InvoiceStatusUnpaid
	}
	return s.repo.CreateInvoice(ctx, inv)
}

func (s *Service) GetInvoice(ctx context.Context, id uuid.UUID) (*Invoice, error) {
	return s.repo.GetInvoice(ctx, id)
}

func (s *Service) ListInvoices(ctx context.Context) ([]Invoice, error) {
	return s.repo.ListInvoices(ctx)
}

func (s *Service) FinalizeInvoice(ctx context.Context, id uuid.UUID) error {
	inv, err := s.repo.GetInvoice(ctx, id)
	if err != nil {
		return err
	}
	// Post to GL
	if err := s.gl.SyncInvoice(ctx, inv.ID.String(), inv.TotalAmount); err != nil {
		return fmt.Errorf("failed to sync to GL: %w", err)
	}

	// Post to Account Ledger (Debit)
	_, err = s.account.PostTransaction(ctx, inv.CustomerID, account.TransactionTypeInvoice, inv.TotalAmount, &inv.ID, "Invoice #"+inv.ID.String()) // Ideally use invoice number
	if err != nil {
		return fmt.Errorf("failed to post to account ledger: %w", err)
	}

	return nil
}
