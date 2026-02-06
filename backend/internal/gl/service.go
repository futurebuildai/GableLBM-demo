package gl

import (
	"context"
	"fmt"

	"log/slog"

	"github.com/gablelbm/gable/internal/domain"
	integration "github.com/gablelbm/gable/internal/integrations/gl"
)

type Service struct {
	adapter integration.GLAdapter
	logger  *slog.Logger
}

func NewService(adapter integration.GLAdapter, logger *slog.Logger) *Service {
	return &Service{adapter: adapter, logger: logger}
}

func (s *Service) SyncInvoice(ctx context.Context, invoiceID string, amount int64) error {
	// logical mapping from Invoice to Journal Entry
	// For now, simple DR AR / CR Revenue
	entry := domain.JournalEntry{
		ReferenceID: invoiceID,
		Memo:        fmt.Sprintf("Invoice %s", invoiceID),
		Lines: []domain.JournalEntryLine{
			{
				AccountName: "Accounts Receivable",
				Debit:       amount,
				Credit:      0,
			},
			{
				AccountName: "Sales Revenue",
				Debit:       0,
				Credit:      amount,
			},
		},
	}

	id, err := s.adapter.PostJournalEntry(ctx, entry)
	if err != nil {
		return fmt.Errorf("failed to post journal entry: %w", err)
	}
	s.logger.Info("Synced Invoice to GL", "invoice_id", invoiceID, "gl_ref_id", id)
	return nil
}
