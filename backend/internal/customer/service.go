package customer

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

func (s *Service) CreateCustomer(ctx context.Context, c *Customer) error {
	return s.repo.CreateCustomer(ctx, c)
}

func (s *Service) GetCustomer(ctx context.Context, id uuid.UUID) (*Customer, error) {
	return s.repo.GetCustomer(ctx, id)
}

func (s *Service) ListCustomers(ctx context.Context) ([]Customer, error) {
	return s.repo.ListCustomers(ctx)
}

func (s *Service) ListPriceLevels(ctx context.Context) ([]PriceLevel, error) {
	return s.repo.ListPriceLevels(ctx)
}

func (s *Service) UpdateBalance(ctx context.Context, id uuid.UUID, delta float64) error {
	return s.repo.UpdateBalance(ctx, id, delta)
}

// Contact management

func (s *Service) CreateContact(ctx context.Context, c *Contact) error {
	return s.repo.CreateContact(ctx, c)
}

func (s *Service) GetContact(ctx context.Context, id uuid.UUID) (*Contact, error) {
	return s.repo.GetContact(ctx, id)
}

func (s *Service) ListContactsByCustomer(ctx context.Context, customerID uuid.UUID) ([]Contact, error) {
	return s.repo.ListContactsByCustomer(ctx, customerID)
}

func (s *Service) UpdateContact(ctx context.Context, c *Contact) error {
	return s.repo.UpdateContact(ctx, c)
}

func (s *Service) DeleteContact(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteContact(ctx, id)
}
