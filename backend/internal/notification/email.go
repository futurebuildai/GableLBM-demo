package notification

import (
	"context"
	"fmt"
	"log/slog"
)

type EmailService interface {
	SendInvoice(ctx context.Context, to string, invoiceID string, pdfBytes []byte) error
}

type LogEmailService struct {
	logger *slog.Logger
}

func NewLogEmailService(logger *slog.Logger) *LogEmailService {
	return &LogEmailService{logger: logger}
}

func (s *LogEmailService) SendInvoice(ctx context.Context, to string, invoiceID string, pdfBytes []byte) error {
	s.logger.Info("MOCK EMAIL SENT",
		"to", to,
		"subject", fmt.Sprintf("Invoice #%s", invoiceID),
		"attachment_size", len(pdfBytes),
		"body", "Please find your invoice attached.",
	)
	return nil
}
