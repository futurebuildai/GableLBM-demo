package purchase_order

import (
	"time"

	"github.com/google/uuid"
)

type PurchaseOrder struct {
	ID        uuid.UUID           `json:"id"`
	VendorID  *uuid.UUID          `json:"vendor_id,omitempty"`
	Status    string              `json:"status"`
	CreatedAt time.Time           `json:"created_at"`
	UpdatedAt time.Time           `json:"updated_at"`
	Lines     []PurchaseOrderLine `json:"lines,omitempty"`
}

type PurchaseOrderLine struct {
	ID             uuid.UUID  `json:"id"`
	POID           uuid.UUID  `json:"po_id"`
	Description    string     `json:"description"`
	Quantity       float64    `json:"quantity"`
	Cost           float64    `json:"cost"`
	LinkedSOLineID *uuid.UUID `json:"linked_so_line_id,omitempty"`
}

const (
	StatusDraft     = "DRAFT"
	StatusSent      = "SENT"
	StatusReceived  = "RECEIVED"
	StatusCancelled = "CANCELLED"
)
