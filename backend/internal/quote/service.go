package quote

import (
	"context"

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

	return s.repo.CreateQuote(ctx, q)
}

func (s *Service) GetQuote(ctx context.Context, id uuid.UUID) (*Quote, error) {
	return s.repo.GetQuote(ctx, id)
}
