package location

import (
	"time"

	"github.com/google/uuid"
)

// LocationType represents the hierarchy level of a location
type LocationType string

const (
	LocTypeZone  LocationType = "ZONE"
	LocTypeAisle LocationType = "AISLE"
	LocTypeRack  LocationType = "RACK"
	LocTypeShelf LocationType = "SHELF"
	LocTypeBin   LocationType = "BIN"
	LocTypeYard  LocationType = "YARD"
)

// Location represents a physical spot in the warehouse/yard
type Location struct {
	ID          uuid.UUID    `json:"id"`
	ParentID    *uuid.UUID   `json:"parent_id,omitempty"` // Root locations have nil ParentID
	Path        string       `json:"path"`                // "Yard A/Row 1"
	Type        LocationType `json:"type"`
	Code        string       `json:"code"` // "A", "1", "B2"
	Description string       `json:"description,omitempty"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`

	// Optional: Computed children for tree view
	Children []Location `json:"children,omitempty"`
}
