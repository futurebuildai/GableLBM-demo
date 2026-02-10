package invoice

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /invoices", h.HandleList)
	mux.HandleFunc("GET /invoices/{id}", h.HandleGet)
	mux.HandleFunc("POST /invoices/{id}/credit-memo", h.HandleCreateCreditMemo)
	mux.HandleFunc("GET /invoices/credit-memos/{customerId}", h.HandleListCreditMemos)
}

func (h *Handler) HandleList(w http.ResponseWriter, r *http.Request) {
	invoices, err := h.svc.ListInvoices(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(invoices)
}

func (h *Handler) HandleGet(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid invoice ID", http.StatusBadRequest)
		return
	}

	inv, err := h.svc.GetInvoice(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(inv)
}

type CreateCreditMemoRequest struct {
	Amount float64 `json:"amount"` // Dollars
	Reason string  `json:"reason"`
}

func (h *Handler) HandleCreateCreditMemo(w http.ResponseWriter, r *http.Request) {
	invoiceIDStr := r.PathValue("id")
	invoiceID, err := uuid.Parse(invoiceIDStr)
	if err != nil {
		http.Error(w, "invalid invoice ID", http.StatusBadRequest)
		return
	}

	var req CreateCreditMemoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Get invoice to find customer
	inv, err := h.svc.GetInvoice(r.Context(), invoiceID)
	if err != nil {
		http.Error(w, "invoice not found", http.StatusNotFound)
		return
	}

	amountCents := int64(req.Amount*100 + 0.5)

	cm, err := h.svc.CreateCreditMemo(r.Context(), inv.CustomerID, &invoiceID, amountCents, req.Reason)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Auto-apply the credit memo
	if err := h.svc.ApplyCreditMemoFull(r.Context(), cm); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(cm)
}

func (h *Handler) HandleListCreditMemos(w http.ResponseWriter, r *http.Request) {
	customerIDStr := r.PathValue("customerId")
	customerID, err := uuid.Parse(customerIDStr)
	if err != nil {
		http.Error(w, "invalid customer ID", http.StatusBadRequest)
		return
	}

	memos, err := h.svc.ListCreditMemos(r.Context(), customerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(memos)
}
