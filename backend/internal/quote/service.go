package quote

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/gablelbm/gable/internal/ai"
	"github.com/gablelbm/gable/internal/product"
	"github.com/google/uuid"
)

type Service struct {
	repo       Repository
	aiClient   *ai.Client
	productSvc *product.Service
}

func NewService(repo Repository, opts ...ServiceOption) *Service {
	return &Service{repo: repo}
}

type ServiceOption func(*Service)

func WithAI(client *ai.Client) ServiceOption {
	return func(s *Service) { s.aiClient = client }
}

func WithProducts(svc *product.Service) ServiceOption {
	return func(s *Service) { s.productSvc = svc }
}

func NewServiceWithDeps(repo Repository, aiClient *ai.Client, productSvc *product.Service) *Service {
	return &Service{repo: repo, aiClient: aiClient, productSvc: productSvc}
}

func (s *Service) CreateQuote(ctx context.Context, q *Quote) error {
	// 1. Calculate Totals
	var total float64
	for i := range q.Lines {
		line := &q.Lines[i]
		line.LineTotal = line.Quantity * line.UnitPrice
		total += line.LineTotal
	}
	// Include freight in total
	total += q.FreightAmount
	q.TotalAmount = total

	// 2. Set Defaults
	if q.State == "" {
		q.State = QuoteStateDraft
	}
	if q.Source == "" {
		q.Source = "manual"
	}
	if q.DeliveryType == "" {
		q.DeliveryType = "PICKUP"
	}
	// Clear vehicle if pickup
	if q.DeliveryType == "PICKUP" {
		q.VehicleID = nil
		q.FreightAmount = 0
	}

	return s.repo.CreateQuote(ctx, q)
}

func (s *Service) GetQuote(ctx context.Context, id uuid.UUID) (*Quote, error) {
	return s.repo.GetQuote(ctx, id)
}

func (s *Service) ListQuotes(ctx context.Context) ([]Quote, error) {
	return s.repo.ListQuotes(ctx)
}

func (s *Service) UpdateState(ctx context.Context, id uuid.UUID, state QuoteState) error {
	q, err := s.repo.GetQuote(ctx, id)
	if err != nil {
		return err
	}

	// Validate state transition
	if err := validateStateTransition(q.State, state); err != nil {
		return err
	}

	now := time.Now()
	q.State = state

	// Set lifecycle timestamp based on target state
	switch state {
	case QuoteStateSent:
		q.SentAt = &now
	case QuoteStateAccepted:
		q.AcceptedAt = &now
	case QuoteStateRejected:
		q.RejectedAt = &now
	}

	return s.repo.UpdateQuote(ctx, q)
}

func (s *Service) UpdateQuote(ctx context.Context, q *Quote) error {
	existing, err := s.repo.GetQuote(ctx, q.ID)
	if err != nil {
		return fmt.Errorf("quote not found: %w", err)
	}
	if existing.State != QuoteStateDraft {
		return fmt.Errorf("only DRAFT quotes can be edited")
	}

	// Recalculate totals
	var total float64
	for i := range q.Lines {
		line := &q.Lines[i]
		line.LineTotal = line.Quantity * line.UnitPrice
		total += line.LineTotal
	}
	// Include freight in total
	total += q.FreightAmount
	q.TotalAmount = total
	q.State = QuoteStateDraft

	// Clear vehicle if pickup
	if q.DeliveryType == "PICKUP" {
		q.VehicleID = nil
		q.FreightAmount = 0
	}

	return s.repo.UpdateQuoteWithLines(ctx, q)
}

func (s *Service) GetAnalytics(ctx context.Context) (*QuoteAnalytics, error) {
	return s.repo.GetQuoteAnalytics(ctx)
}

func (s *Service) GetOriginalFile(ctx context.Context, id uuid.UUID) ([]byte, string, string, error) {
	return s.repo.GetOriginalFile(ctx, id)
}

// AIEditQuote uses Claude to interpret a natural language editing command and apply it to a draft quote.
func (s *Service) AIEditQuote(ctx context.Context, quoteID uuid.UUID, userMessage string) (*Quote, string, error) {
	if s.aiClient == nil {
		return nil, "", fmt.Errorf("AI client not configured")
	}

	// Load current quote
	q, err := s.repo.GetQuote(ctx, quoteID)
	if err != nil {
		return nil, "", fmt.Errorf("quote not found: %w", err)
	}
	if q.State != QuoteStateDraft {
		return nil, "", fmt.Errorf("only DRAFT quotes can be edited")
	}

	// Build product catalog context for "add" commands
	var catalogContext string
	if s.productSvc != nil {
		products, err := s.productSvc.ListProducts(ctx)
		if err == nil && len(products) > 0 {
			var sb strings.Builder
			sb.WriteString("\nAvailable products in catalog:\n")
			for _, p := range products {
				sb.WriteString(fmt.Sprintf("- product_id: %s, sku: %s, description: %s, uom: %s, price: %.2f\n",
					p.ID, p.SKU, p.Description, p.UOMPrimary, p.BasePrice))
			}
			catalogContext = sb.String()
		}
	}

	// Build current lines JSON
	type lineInfo struct {
		ProductID   string  `json:"product_id"`
		SKU         string  `json:"sku"`
		Description string  `json:"description"`
		Quantity    float64 `json:"quantity"`
		UOM         string  `json:"uom"`
		UnitPrice   float64 `json:"unit_price"`
		UnitCost    float64 `json:"unit_cost"`
	}
	currentLines := make([]lineInfo, len(q.Lines))
	for i, l := range q.Lines {
		currentLines[i] = lineInfo{
			ProductID:   l.ProductID.String(),
			SKU:         l.SKU,
			Description: l.Description,
			Quantity:    l.Quantity,
			UOM:         string(l.UOM),
			UnitPrice:   l.UnitPrice,
			UnitCost:    l.UnitCost,
		}
	}
	linesJSON, _ := json.MarshalIndent(currentLines, "", "  ")

	aiSystemPrompt := `You are a quote editing assistant for a lumber and building materials ERP system.
Given the current quote lines and a user command, return the updated lines as JSON.

Return ONLY valid JSON in this exact format:
{
  "lines": [
    {
      "product_id": "uuid-string",
      "sku": "SKU-CODE",
      "description": "Product Description",
      "quantity": 50,
      "uom": "EA",
      "unit_price": 5.99,
      "unit_cost": 3.50
    }
  ],
  "explanation": "Brief description of what was changed"
}

Rules:
- Keep all existing lines that aren't being modified
- For quantity changes, update only the quantity field
- For adding new items, use product_id and pricing from the catalog if available
- For removing items, simply omit them from the lines array
- unit_cost can remain the same as the original if not specified
- Output ONLY the JSON object, nothing else`

	userPrompt := fmt.Sprintf("Current quote lines:\n%s\n%s\nUser command: %s", string(linesJSON), catalogContext, userMessage)

	// Call Claude
	rawResponse, err := s.aiClient.SendMessage(ctx, aiSystemPrompt, userPrompt)
	if err != nil {
		return nil, "", fmt.Errorf("AI request failed: %w", err)
	}

	// Strip markdown fences if present
	cleaned := strings.TrimSpace(rawResponse)
	if strings.HasPrefix(cleaned, "```") {
		if idx := strings.Index(cleaned, "\n"); idx != -1 {
			cleaned = cleaned[idx+1:]
		}
		if idx := strings.LastIndex(cleaned, "```"); idx != -1 {
			cleaned = cleaned[:idx]
		}
		cleaned = strings.TrimSpace(cleaned)
	}

	// Parse AI response
	var aiResult struct {
		Lines []struct {
			ProductID   string  `json:"product_id"`
			SKU         string  `json:"sku"`
			Description string  `json:"description"`
			Quantity    float64 `json:"quantity"`
			UOM         string  `json:"uom"`
			UnitPrice   float64 `json:"unit_price"`
			UnitCost    float64 `json:"unit_cost"`
		} `json:"lines"`
		Explanation string `json:"explanation"`
	}
	if err := json.Unmarshal([]byte(cleaned), &aiResult); err != nil {
		return nil, rawResponse, fmt.Errorf("failed to parse AI response: %w", err)
	}

	// Build updated quote lines
	newLines := make([]QuoteLine, len(aiResult.Lines))
	for i, al := range aiResult.Lines {
		productID, _ := uuid.Parse(al.ProductID)
		newLines[i] = QuoteLine{
			ProductID:   productID,
			SKU:         al.SKU,
			Description: al.Description,
			Quantity:    al.Quantity,
			UOM:         product.UOM(al.UOM),
			UnitPrice:   al.UnitPrice,
			UnitCost:    al.UnitCost,
		}
	}
	q.Lines = newLines

	// Recalculate and save
	if err := s.UpdateQuote(ctx, q); err != nil {
		return nil, aiResult.Explanation, fmt.Errorf("failed to save updated quote: %w", err)
	}

	// Reload to get fresh data
	updated, err := s.repo.GetQuote(ctx, quoteID)
	if err != nil {
		return q, aiResult.Explanation, nil
	}

	return updated, aiResult.Explanation, nil
}

// validateStateTransition ensures the state change is valid.
func validateStateTransition(from, to QuoteState) error {
	allowed := map[QuoteState][]QuoteState{
		QuoteStateDraft:    {QuoteStateSent, QuoteStateAccepted, QuoteStateRejected, QuoteStateExpired},
		QuoteStateSent:     {QuoteStateAccepted, QuoteStateRejected, QuoteStateExpired},
		QuoteStateAccepted: {}, // terminal
		QuoteStateRejected: {QuoteStateDraft}, // allow re-opening
		QuoteStateExpired:  {QuoteStateDraft}, // allow re-opening
	}

	targets, ok := allowed[from]
	if !ok {
		return fmt.Errorf("unknown current state: %s", from)
	}

	for _, t := range targets {
		if t == to {
			return nil
		}
	}
	return fmt.Errorf("cannot transition from %s to %s", from, to)
}
