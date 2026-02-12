package dashboard

import (
	"context"
	"sync"
	"time"
)

// cache is a type-safe, TTL-based in-memory cache.
type cache[T any] struct {
	data      T
	timestamp time.Time
	valid     bool
}

// cacheStore holds all dashboard caches with a shared mutex for simplicity.
type cacheStore struct {
	mu              sync.RWMutex
	ttl             time.Duration
	summary         cache[*DashboardSummary]
	inventoryAlerts cache[[]InventoryAlert]
	topCustomers    cache[[]TopCustomer]
	orderActivity   cache[*OrderActivity]
	revenueTrend    cache[[]RevenueTrendPoint]
}

func newCacheStore(ttl time.Duration) *cacheStore {
	return &cacheStore{ttl: ttl}
}

func getCache[T any](mu *sync.RWMutex, c *cache[T], ttl time.Duration) (T, bool) {
	mu.RLock()
	defer mu.RUnlock()
	if c.valid && time.Since(c.timestamp) < ttl {
		return c.data, true
	}
	var zero T
	return zero, false
}

func setCache[T any](mu *sync.RWMutex, c *cache[T], data T) {
	mu.Lock()
	defer mu.Unlock()
	c.data = data
	c.timestamp = time.Now()
	c.valid = true
}

// Service provides dashboard business logic with type-safe caching.
type Service struct {
	repo  Repository
	store *cacheStore
}

// NewService creates a new dashboard service.
func NewService(repo Repository) *Service {
	return &Service{
		repo:  repo,
		store: newCacheStore(60 * time.Second),
	}
}

// GetSummary returns the dashboard summary with caching.
func (s *Service) GetSummary(ctx context.Context) (*DashboardSummary, error) {
	if cached, ok := getCache(&s.store.mu, &s.store.summary, s.store.ttl); ok {
		return cached, nil
	}

	data, err := s.repo.GetDashboardSummary(ctx)
	if err != nil {
		return nil, err
	}

	setCache(&s.store.mu, &s.store.summary, data)
	return data, nil
}

// GetInventoryAlerts returns inventory alerts with caching.
func (s *Service) GetInventoryAlerts(ctx context.Context) ([]InventoryAlert, error) {
	if cached, ok := getCache(&s.store.mu, &s.store.inventoryAlerts, s.store.ttl); ok {
		return cached, nil
	}

	data, err := s.repo.GetInventoryAlerts(ctx, 10)
	if err != nil {
		return nil, err
	}

	setCache(&s.store.mu, &s.store.inventoryAlerts, data)
	return data, nil
}

// GetTopCustomers returns top customers with caching.
func (s *Service) GetTopCustomers(ctx context.Context) ([]TopCustomer, error) {
	if cached, ok := getCache(&s.store.mu, &s.store.topCustomers, s.store.ttl); ok {
		return cached, nil
	}

	data, err := s.repo.GetTopCustomers(ctx, 5, 30)
	if err != nil {
		return nil, err
	}

	setCache(&s.store.mu, &s.store.topCustomers, data)
	return data, nil
}

// GetOrderActivity returns order activity with caching.
func (s *Service) GetOrderActivity(ctx context.Context) (*OrderActivity, error) {
	if cached, ok := getCache(&s.store.mu, &s.store.orderActivity, s.store.ttl); ok {
		return cached, nil
	}

	data, err := s.repo.GetOrderActivity(ctx, 10)
	if err != nil {
		return nil, err
	}

	setCache(&s.store.mu, &s.store.orderActivity, data)
	return data, nil
}

// GetRevenueTrend returns revenue trend for chart.
func (s *Service) GetRevenueTrend(ctx context.Context) ([]RevenueTrendPoint, error) {
	if cached, ok := getCache(&s.store.mu, &s.store.revenueTrend, s.store.ttl); ok {
		return cached, nil
	}

	data, err := s.repo.GetRevenueTrend(ctx, 7)
	if err != nil {
		return nil, err
	}

	setCache(&s.store.mu, &s.store.revenueTrend, data)
	return data, nil
}
