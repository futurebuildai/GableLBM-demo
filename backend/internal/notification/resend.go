package notification

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
)

// ResendEmailService implements EmailService using the Resend API.
type ResendEmailService struct {
	apiKey string
	from   string
	client *http.Client
	logger *slog.Logger
}

// NewResendEmailService creates a new Resend email service.
func NewResendEmailService(apiKey string, logger *slog.Logger) *ResendEmailService {
	return &ResendEmailService{
		apiKey: apiKey,
		from:   "Gable ERP <notifications@gablelbm.com>",
		client: &http.Client{Timeout: 10 * time.Second},
		logger: logger,
	}
}

type resendEmailRequest struct {
	From    string   `json:"from"`
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	HTML    string   `json:"html"`
}

func (s *ResendEmailService) send(ctx context.Context, to string, subject string, html string) error {
	payload := resendEmailRequest{
		From:    s.from,
		To:      []string{to},
		Subject: subject,
		HTML:    html,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal resend payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.resend.com/emails", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create resend request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+s.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("resend API call failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("resend returned status %d: %s", resp.StatusCode, string(respBody))
	}

	s.logger.Info("Email sent via Resend", "to", to, "subject", subject)
	return nil
}

func (s *ResendEmailService) SendInvoice(ctx context.Context, to string, invoiceID string, pdfBytes []byte) error {
	subject := fmt.Sprintf("Invoice #%s from Gable ERP", invoiceID)
	html := fmt.Sprintf(`<p>Please find your invoice #%s attached.</p>`, invoiceID)
	return s.send(ctx, to, subject, html)
}

func (s *ResendEmailService) SendDeliveryNotification(ctx context.Context, to string, subject string, body string) error {
	html := fmt.Sprintf(`<p>%s</p>`, body)
	return s.send(ctx, to, subject, html)
}

func (s *ResendEmailService) SendQuoteNotification(ctx context.Context, to string, quoteID string, customerName string, totalAmount float64, quoteURL string) error {
	subject := fmt.Sprintf("New Quick Quote from %s — $%.2f", customerName, totalAmount)
	html := fmt.Sprintf(`
		<div style="font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; max-width: 600px; margin: 0 auto;">
			<div style="background: #0A0B10; padding: 24px 32px; border-radius: 12px 12px 0 0;">
				<h1 style="color: #00FFA3; margin: 0; font-size: 24px;">Gable ERP</h1>
			</div>
			<div style="background: #ffffff; padding: 32px; border: 1px solid #e5e7eb; border-top: none; border-radius: 0 0 12px 12px;">
				<h2 style="color: #1a1a1a; margin-top: 0;">New Quote Request</h2>
				<p style="color: #4b5563; line-height: 1.6;">
					<strong>%s</strong> has submitted a new AI-assisted quote request via the Contractor Portal.
				</p>
				<table style="width: 100%%; border-collapse: collapse; margin: 20px 0;">
					<tr>
						<td style="padding: 8px 0; color: #6b7280; border-bottom: 1px solid #f3f4f6;">Quote ID</td>
						<td style="padding: 8px 0; color: #1a1a1a; text-align: right; border-bottom: 1px solid #f3f4f6; font-family: monospace;">%s</td>
					</tr>
					<tr>
						<td style="padding: 8px 0; color: #6b7280; border-bottom: 1px solid #f3f4f6;">Customer</td>
						<td style="padding: 8px 0; color: #1a1a1a; text-align: right; border-bottom: 1px solid #f3f4f6;">%s</td>
					</tr>
					<tr>
						<td style="padding: 8px 0; color: #6b7280;">Estimated Total</td>
						<td style="padding: 8px 0; color: #1a1a1a; text-align: right; font-weight: 600; font-size: 18px;">$%.2f</td>
					</tr>
				</table>
				<a href="%s" style="display: inline-block; background: #00FFA3; color: #0A0B10; padding: 12px 24px; border-radius: 8px; text-decoration: none; font-weight: 600; margin-top: 8px;">
					Review Quote in ERP →
				</a>
				<p style="color: #9ca3af; font-size: 12px; margin-top: 24px;">
					This quote was generated using AI-assisted material list parsing and requires review before sending to the customer.
				</p>
			</div>
		</div>
	`, customerName, quoteID, customerName, totalAmount, quoteURL)
	return s.send(ctx, to, subject, html)
}
