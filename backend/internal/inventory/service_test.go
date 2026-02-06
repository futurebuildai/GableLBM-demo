package inventory

import (
	"context"
	"testing"

	"github.com/google/uuid"
)

// MockRepository implements Repository for testing
type MockRepository struct {
	invs map[string]*Inventory
}

func (m *MockRepository) GetInventory(ctx context.Context, productID uuid.UUID, locationID *uuid.UUID) (*Inventory, error) {
	if locationID == nil {
		return nil, nil // Not needed for transfer test logic which specifies location
	}
	key := productID.String() + ":" + locationID.String()
	if i, ok := m.invs[key]; ok {
		return i, nil
	}
	return nil, nil
}

func (m *MockRepository) CreateInventory(ctx context.Context, inv *Inventory) error {
	key := inv.ProductID.String() + ":" + inv.LocationID.String()
	m.invs[key] = inv
	return nil
}

func (m *MockRepository) UpdateInventory(ctx context.Context, inv *Inventory) error {
	key := inv.ProductID.String() + ":" + inv.LocationID.String()
	m.invs[key] = inv
	return nil
}

func (m *MockRepository) ExecuteInTx(ctx context.Context, fn func(context.Context) error) error {
	// Just run function directly (mock transaction)
	return fn(ctx)
}

// Stubs
func (m *MockRepository) ListInventoryByProduct(ctx context.Context, productID uuid.UUID) ([]Inventory, error) {
	return nil, nil
}
func (m *MockRepository) AllocateStock(ctx context.Context, inventoryID uuid.UUID, delta float64) error {
	return nil
}
func (m *MockRepository) FulfillStock(ctx context.Context, inventoryID uuid.UUID, delta float64) error {
	return nil
}

func TestTransferStock(t *testing.T) {
	repo := &MockRepository{
		invs: make(map[string]*Inventory),
	}
	svc := NewService(repo)

	prodID := uuid.New()
	loc1 := uuid.New()
	loc2 := uuid.New()

	// Setup initial stock in loc1
	repo.invs[prodID.String()+":"+loc1.String()] = &Inventory{
		ProductID:  prodID,
		LocationID: &loc1,
		Quantity:   100,
	}

	// Move 50 from loc1 to loc2
	err := svc.MoveStock(context.Background(), StockMovementRequest{
		ProductID:      prodID,
		FromLocationID: &loc1,
		ToLocationID:   loc2,
		Quantity:       50,
		Reason:         "Test",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify loc1 has 50
	i1 := repo.invs[prodID.String()+":"+loc1.String()]
	if i1.Quantity != 50 {
		t.Errorf("expected loc1 to have 50, got %f", i1.Quantity)
	}

	// Verify loc2 has 50
	i2 := repo.invs[prodID.String()+":"+loc2.String()]
	if i2 == nil {
		t.Fatal("expected loc2 to be created")
	}
	if i2.Quantity != 50 {
		t.Errorf("expected loc2 to have 50, got %f", i2.Quantity)
	}
}
