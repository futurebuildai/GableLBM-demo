package product

import (
	"time"

	"github.com/google/uuid"
)

// UOM represents the strict Unit of Measure types matching the database ENUM
type UOM string

const (
	UOM_PCS    UOM = "PCS"
	UOM_EA     UOM = "EA"
	UOM_LF     UOM = "LF"
	UOM_SF     UOM = "SF"
	UOM_BF     UOM = "BF"
	UOM_MBF    UOM = "MBF"
	UOM_SQ     UOM = "SQ"
	UOM_BOX    UOM = "BOX"
	UOM_CTN    UOM = "CTN"
	UOM_RL     UOM = "RL"
	UOM_GAL    UOM = "GAL"
	UOM_LBS    UOM = "LBS"
	UOM_BAG    UOM = "BAG"
	UOM_BUNDLE UOM = "BUNDLE"
	UOM_PAIR   UOM = "PAIR"
	UOM_SET    UOM = "SET"
)

// Product represents a catalog item
type Product struct {
	ID            uuid.UUID `json:"id"`
	SKU           string    `json:"sku"`
	Description   string    `json:"description"`
	UOMPrimary    UOM       `json:"uom_primary"`
	TotalQuantity float64   `json:"total_quantity"` // Aggregated from inventory
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
