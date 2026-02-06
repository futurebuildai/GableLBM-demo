package location

import (
	"context"
	"fmt"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateLocation(ctx context.Context, loc *Location) error {
	// Business logic: Ensure Code is present
	if loc.Code == "" {
		return fmt.Errorf("location code is required")
	}

	// Validate Type logic if needed (e.g. ZONE can't have parent? or must be root)
	// For now simple pass-through with path construction logic if needed?
	// The frontend might provide the path, or we compute it.
	// For MVP, assuming frontend or caller provides correct Path for now.

	return s.repo.CreateLocation(ctx, loc)
}

func (s *Service) ListLocations(ctx context.Context) ([]Location, error) {
	return s.repo.ListLocations(ctx)
}
