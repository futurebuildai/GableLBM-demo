package pricing

import (
	"time"

	"github.com/google/uuid"
)

type PricingSource string

const (
	SourceContract PricingSource = "CONTRACT"
	SourceTier     PricingSource = "TIER"
	SourceRetail   PricingSource = "RETAIL"
)

type CalculatedPrice struct {
	ProductID     uuid.UUID     `json:"product_id"`
	OriginalPrice float64       `json:"original_price"` // Base Retail
	FinalPrice    float64       `json:"final_price"`
	DiscountPct   float64       `json:"discount_pct"`
	Source        PricingSource `json:"source"`
	Details       string        `json:"details"` // e.g. "Gold Member Discount"
}

type CustomerContract struct {
	ID            uuid.UUID `json:"id"`
	CustomerID    uuid.UUID `json:"customer_id"`
	ProductID     uuid.UUID `json:"product_id"`
	ContractPrice float64   `json:"contract_price"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
