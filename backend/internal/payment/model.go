package payment

import (
	"time"

	"github.com/google/uuid"
)

type PaymentMethod string

const (
	PaymentMethodCash    PaymentMethod = "CASH"
	PaymentMethodCard    PaymentMethod = "CARD"
	PaymentMethodCheck   PaymentMethod = "CHECK"
	PaymentMethodAccount PaymentMethod = "ACCOUNT"
)

type Payment struct {
	ID        uuid.UUID     `json:"id" db:"id"`
	InvoiceID uuid.UUID     `json:"invoice_id" db:"invoice_id"`
	Amount    int64         `json:"amount" db:"amount"` // In Cents
	Method    PaymentMethod `json:"method" db:"method"`
	Reference string        `json:"reference" db:"reference"`
	Notes     string        `json:"notes" db:"notes"`
	CreatedAt time.Time     `json:"created_at" db:"created_at"`
}
