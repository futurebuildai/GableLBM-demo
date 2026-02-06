package invoice

import (
	"time"

	"github.com/google/uuid"
)

type InvoiceStatus string

const (
	InvoiceStatusUnpaid  InvoiceStatus = "UNPAID"
	InvoiceStatusPartial InvoiceStatus = "PARTIAL"
	InvoiceStatusPaid    InvoiceStatus = "PAID"
	InvoiceStatusVoid    InvoiceStatus = "VOID"
	InvoiceStatusOverdue InvoiceStatus = "OVERDUE"
)

type Invoice struct {
	ID          uuid.UUID     `json:"id" db:"id"`
	OrderID     uuid.UUID     `json:"order_id" db:"order_id"`
	CustomerID  uuid.UUID     `json:"customer_id" db:"customer_id"`
	Status      InvoiceStatus `json:"status" db:"status"`
	TotalAmount int64         `json:"total_amount" db:"total_amount"` // Cents
	DueDate     *time.Time    `json:"due_date" db:"due_date"`
	PaidAt      *time.Time    `json:"paid_at" db:"paid_at"`
	CreatedAt   time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at" db:"updated_at"`

	// Relations
	Lines []InvoiceLine `json:"lines,omitempty" db:"-"`
}

type InvoiceLine struct {
	ID        uuid.UUID `json:"id" db:"id"`
	InvoiceID uuid.UUID `json:"invoice_id" db:"invoice_id"`
	ProductID uuid.UUID `json:"product_id" db:"product_id"`
	Quantity  float64   `json:"quantity" db:"quantity"`
	PriceEach int64     `json:"price_each" db:"price_each"` // Cents
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
