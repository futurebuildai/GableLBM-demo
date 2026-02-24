package payment

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	// Existing routes
	mux.HandleFunc("POST /api/payments", h.CreatePayment)
	mux.HandleFunc("GET /api/invoices/{id}/payments", h.GetPaymentHistory)

	// Run Payments gateway routes
	mux.HandleFunc("POST /api/payments/intent", h.CreatePaymentIntent)
	mux.HandleFunc("POST /api/payments/card", h.ProcessCardPayment)
	mux.HandleFunc("POST /api/payments/refund", h.ProcessRefund)
}

// CreatePayment handles non-card payments (cash, check, account).
func (h *Handler) CreatePayment(w http.ResponseWriter, r *http.Request) {
	var req CreatePaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	payment, err := h.service.ProcessPayment(r.Context(), req.InvoiceID, req.Amount, req.Method, req.Reference, req.Notes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(payment)
}

// CreatePaymentIntent returns the Run Payments public key for Runner.js tokenization.
// The frontend calls this before showing the card input form.
func (h *Handler) CreatePaymentIntent(w http.ResponseWriter, r *http.Request) {
	var req CreatePaymentIntentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	publicKey := h.service.GetPublicKey()
	if publicKey == "" {
		http.Error(w, "Payment gateway not configured", http.StatusServiceUnavailable)
		return
	}

	amountCents := int64(req.Amount*100.0 + 0.5)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(PaymentIntentResponse{
		PublicKey: publicKey,
		InvoiceID: req.InvoiceID.String(),
		Amount:    amountCents,
	})
}

// ProcessCardPayment handles tokenized card payments through Run Payments.
// Flow: Frontend tokenizes via Runner.js → sends token here → we charge via gateway.
func (h *Handler) ProcessCardPayment(w http.ResponseWriter, r *http.Request) {
	var req ProcessCardPaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.TokenID == "" {
		http.Error(w, "token_id is required", http.StatusBadRequest)
		return
	}

	payment, err := h.service.ProcessCardPayment(r.Context(), req.InvoiceID, req.TokenID, req.Amount, req.Notes)
	if err != nil {
		// Check if it's a decline vs. system error
		http.Error(w, err.Error(), http.StatusPaymentRequired)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(payment)
}

// ProcessRefund handles full or partial refunds of card payments.
func (h *Handler) ProcessRefund(w http.ResponseWriter, r *http.Request) {
	var req RefundRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	refund, err := h.service.RefundPayment(r.Context(), req.PaymentID, req.Amount, req.Reason)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(refund)
}

// GetPaymentHistory returns all payments for an invoice.
func (h *Handler) GetPaymentHistory(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid Invoice ID", http.StatusBadRequest)
		return
	}

	history, err := h.service.GetHistory(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(history)
}

// CreatePaymentRequest is the existing request struct for non-card payments.
type CreatePaymentRequest struct {
	InvoiceID uuid.UUID     `json:"invoice_id"`
	Amount    float64       `json:"amount"`
	Method    PaymentMethod `json:"method"`
	Reference string        `json:"reference"`
	Notes     string        `json:"notes"`
}
