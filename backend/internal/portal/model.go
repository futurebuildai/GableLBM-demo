package portal

import (
	"time"

	"github.com/google/uuid"
)

// CustomerUser represents a contractor/customer who can log into the portal.
type CustomerUser struct {
	ID           uuid.UUID `json:"id"`
	CustomerID   uuid.UUID `json:"customer_id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // Never serialize
	Name         string    `json:"name"`
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// PortalConfig holds white-label branding for the dealer portal.
type PortalConfig struct {
	ID           uuid.UUID `json:"id"`
	DealerName   string    `json:"dealer_name"`
	LogoURL      string    `json:"logo_url"`
	PrimaryColor string    `json:"primary_color"`
	SupportEmail string    `json:"support_email"`
	SupportPhone string    `json:"support_phone"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// --- Request/Response DTOs ---

// LoginRequest is the payload for portal login.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse is returned on successful login.
type LoginResponse struct {
	Token  string       `json:"token"`
	User   CustomerUser `json:"user"`
	Config PortalConfig `json:"config"`
}

// PortalDashboardDTO aggregates AR and activity data for the contractor.
type PortalDashboardDTO struct {
	BalanceDue   float64          `json:"balance_due"`
	CreditLimit  float64          `json:"credit_limit"`
	PastDue      float64          `json:"past_due"`
	RecentOrders []PortalOrderDTO `json:"recent_orders"`
}

// PortalOrderDTO is a customer-facing order summary.
type PortalOrderDTO struct {
	ID          uuid.UUID       `json:"id"`
	Status      string          `json:"status"`
	TotalAmount float64         `json:"total_amount"`
	CreatedAt   time.Time       `json:"created_at"`
	Lines       []PortalLineDTO `json:"lines"`
}

// PortalLineDTO is a customer-facing order/invoice line item.
type PortalLineDTO struct {
	ProductID   uuid.UUID `json:"product_id"`
	ProductSKU  string    `json:"product_sku"`
	ProductName string    `json:"product_name"`
	Quantity    float64   `json:"quantity"`
	PriceEach   float64   `json:"price_each"`
}

// PortalInvoiceDTO is a customer-facing invoice summary.
type PortalInvoiceDTO struct {
	ID           uuid.UUID       `json:"id"`
	OrderID      uuid.UUID       `json:"order_id"`
	Status       string          `json:"status"`
	TotalAmount  float64         `json:"total_amount"`
	Subtotal     float64         `json:"subtotal"`
	TaxAmount    float64         `json:"tax_amount"`
	PaymentTerms string          `json:"payment_terms"`
	DueDate      *time.Time      `json:"due_date"`
	PaidAt       *time.Time      `json:"paid_at"`
	CreatedAt    time.Time       `json:"created_at"`
	Lines        []PortalLineDTO `json:"lines"`
}

// PortalDeliveryDTO is a customer-facing delivery with POD info.
type PortalDeliveryDTO struct {
	ID           uuid.UUID  `json:"id"`
	OrderID      uuid.UUID  `json:"order_id"`
	Status       string     `json:"status"`
	PODProofURL  *string    `json:"pod_proof_url"`
	PODSignedBy  *string    `json:"pod_signed_by"`
	PODTimestamp *time.Time `json:"pod_timestamp"`
	CreatedAt    time.Time  `json:"created_at"`
	OrderNumber  *string    `json:"order_number"`
}

// ReorderRequest tells which historical order to duplicate.
type ReorderRequest struct {
	OrderID uuid.UUID `json:"order_id"`
}

// ReorderResponse confirms the new draft order.
type ReorderResponse struct {
	OrderID uuid.UUID `json:"order_id"`
	Message string    `json:"message"`
}
