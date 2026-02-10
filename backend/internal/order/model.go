package order

import (
	"time"

	"github.com/google/uuid"
)

type OrderStatus string

const (
	StatusDraft     OrderStatus = "DRAFT"
	StatusConfirmed OrderStatus = "CONFIRMED"
	StatusFulfilled OrderStatus = "FULFILLED"
	StatusCancelled OrderStatus = "CANCELLED"
)

type Order struct {
	ID           uuid.UUID   `json:"id"`
	CustomerID   uuid.UUID   `json:"customer_id"`
	CustomerName string      `json:"customer_name,omitempty"`
	QuoteID      *uuid.UUID  `json:"quote_id,omitempty"`
	Status       OrderStatus `json:"status"`
	TotalAmount  float64     `json:"total_amount"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`

	// Relations
	Lines []OrderLine `json:"lines,omitempty"`
}

type OrderLine struct {
	ID               uuid.UUID  `json:"id"`
	OrderID          uuid.UUID  `json:"order_id"`
	ProductID        uuid.UUID  `json:"product_id"`
	ProductSKU       string     `json:"product_sku,omitempty"`
	ProductName      string     `json:"product_name,omitempty"`
	Quantity         float64    `json:"quantity"`
	PriceEach        float64    `json:"price_each"`
	IsSpecialOrder   bool       `json:"is_special_order"`
	VendorID         *uuid.UUID `json:"vendor_id,omitempty"`
	SpecialOrderCost float64    `json:"special_order_cost,omitempty"`
}

type CreateOrderRequest struct {
	CustomerID uuid.UUID          `json:"customer_id"`
	QuoteID    *uuid.UUID         `json:"quote_id"`
	Lines      []OrderLineRequest `json:"lines"`
}

type OrderLineRequest struct {
	ProductID        uuid.UUID  `json:"product_id"`
	Quantity         float64    `json:"quantity"`
	PriceEach        float64    `json:"price_each"`
	IsSpecialOrder   bool       `json:"is_special_order"`
	VendorID         *uuid.UUID `json:"vendor_id"`
	SpecialOrderCost float64    `json:"special_order_cost"`
}

type UpdateStatusRequest struct {
	Status OrderStatus `json:"status"`
}
