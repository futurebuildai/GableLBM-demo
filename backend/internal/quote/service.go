package quote

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateQuote(ctx context.Context, q *Quote) error {
	// 1. Calculate Totals
	var total float64
	for i := range q.Lines {
		line := &q.Lines[i]
		line.LineTotal = line.Quantity * line.UnitPrice
		total += line.LineTotal
	}
	q.TotalAmount = total

	// 2. Set Defaults
	if q.State == "" {
		q.State = QuoteStateDraft
	}
	if q.Source == "" {
		q.Source = "manual"
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

func (s *Service) GetAnalytics(ctx context.Context) (*QuoteAnalytics, error) {
	return s.repo.GetQuoteAnalytics(ctx)
}

func (s *Service) GetOriginalFile(ctx context.Context, id uuid.UUID) ([]byte, string, string, error) {
	return s.repo.GetOriginalFile(ctx, id)
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
