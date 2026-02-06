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

type CreatePaymentRequest struct {
	InvoiceID uuid.UUID     `json:"invoice_id"`
	Amount    float64       `json:"amount"`
	Method    PaymentMethod `json:"method"`
	Reference string        `json:"reference"`
	Notes     string        `json:"notes"`
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/payments", h.CreatePayment)
	mux.HandleFunc("GET /api/invoices/{id}/payments", h.GetPaymentHistory)
}

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
