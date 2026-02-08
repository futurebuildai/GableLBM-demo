package dashboard

import (
	"context"
	"sync"
	"time"
)

type cachedItem struct {
	data      interface{}
	timestamp time.Time
}

// Service provides dashboard business logic with caching.
type Service struct {
	repo       Repository
	cache      map[string]cachedItem
	cacheMutex sync.RWMutex
	cacheTTL   time.Duration
}

// NewService creates a new dashboard service.
func NewService(repo Repository) *Service {
	return &Service{
		repo:     repo,
		cache:    make(map[string]cachedItem),
		cacheTTL: 60 * time.Second,
	}
}

func (s *Service) getFromCache(key string) (interface{}, bool) {
	s.cacheMutex.RLock()
	defer s.cacheMutex.RUnlock()

	if item, ok := s.cache[key]; ok {
		if time.Since(item.timestamp) < s.cacheTTL {
			return item.data, true
		}
	}
	return nil, false
}

func (s *Service) setCache(key string, data interface{}) {
	s.cacheMutex.Lock()
	defer s.cacheMutex.Unlock()
	s.cache[key] = cachedItem{data: data, timestamp: time.Now()}
}

// GetSummary returns the dashboard summary with caching.
func (s *Service) GetSummary(ctx context.Context) (*DashboardSummary, error) {
	const key = "dashboard:summary"
	if cached, ok := s.getFromCache(key); ok {
		return cached.(*DashboardSummary), nil
	}

	data, err := s.repo.GetDashboardSummary(ctx)
	if err != nil {
		return nil, err
	}

	s.setCache(key, data)
	return data, nil
}

// GetInventoryAlerts returns inventory alerts with caching.
func (s *Service) GetInventoryAlerts(ctx context.Context) ([]InventoryAlert, error) {
	const key = "dashboard:inventory-alerts"
	if cached, ok := s.getFromCache(key); ok {
		return cached.([]InventoryAlert), nil
	}

	data, err := s.repo.GetInventoryAlerts(ctx, 10)
	if err != nil {
		return nil, err
	}

	s.setCache(key, data)
	return data, nil
}

// GetTopCustomers returns top customers with caching.
func (s *Service) GetTopCustomers(ctx context.Context) ([]TopCustomer, error) {
	const key = "dashboard:top-customers"
	if cached, ok := s.getFromCache(key); ok {
		return cached.([]TopCustomer), nil
	}

	data, err := s.repo.GetTopCustomers(ctx, 5, 30)
	if err != nil {
		return nil, err
	}

	s.setCache(key, data)
	return data, nil
}

// GetOrderActivity returns order activity with caching.
func (s *Service) GetOrderActivity(ctx context.Context) (*OrderActivity, error) {
	const key = "dashboard:order-activity"
	if cached, ok := s.getFromCache(key); ok {
		return cached.(*OrderActivity), nil
	}

	data, err := s.repo.GetOrderActivity(ctx, 10)
	if err != nil {
		return nil, err
	}

	s.setCache(key, data)
	return data, nil
}

// GetRevenueTrend returns revenue trend for chart.
func (s *Service) GetRevenueTrend(ctx context.Context) ([]RevenueTrendPoint, error) {
	const key = "dashboard:revenue-trend"
	if cached, ok := s.getFromCache(key); ok {
		return cached.([]RevenueTrendPoint), nil
	}

	data, err := s.repo.GetRevenueTrend(ctx, 7)
	if err != nil {
		return nil, err
	}

	s.setCache(key, data)
	return data, nil
}
