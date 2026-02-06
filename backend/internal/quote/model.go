package quote

import (
	"time"

	"github.com/gablelbm/gable/internal/product"
	"github.com/google/uuid"
)

type QuoteState string

const (
	QuoteStateDraft    QuoteState = "DRAFT"
	QuoteStateSent     QuoteState = "SENT"
	QuoteStateAccepted QuoteState = "ACCEPTED"
	QuoteStateRejected QuoteState = "REJECTED"
	QuoteStateExpired  QuoteState = "EXPIRED"
)

type Quote struct {
	ID          uuid.UUID  `json:"id"`
	CustomerID  uuid.UUID  `json:"customer_id"`
	JobID       *uuid.UUID `json:"job_id,omitempty"`
	State       QuoteState `json:"state"`
	TotalAmount float64    `json:"total_amount"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`

	Lines []QuoteLine `json:"lines,omitempty"`
}

type QuoteLine struct {
	ID          uuid.UUID   `json:"id"`
	QuoteID     uuid.UUID   `json:"quote_id"`
	ProductID   uuid.UUID   `json:"product_id"`
	SKU         string      `json:"sku"`
	Description string      `json:"description"`
	Quantity    float64     `json:"quantity"`
	UOM         product.UOM `json:"uom"`
	UnitPrice   float64     `json:"unit_price"`
	LineTotal   float64     `json:"line_total"`
	CreatedAt   time.Time   `json:"created_at"`
}
