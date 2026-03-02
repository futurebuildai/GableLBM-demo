package tax

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

// Handler exposes tax-related HTTP endpoints.
type Handler struct {
	svc TaxCalculator
}

// NewHandler creates a new tax Handler.
func NewHandler(svc TaxCalculator) *Handler {
	return &Handler{svc: svc}
}

// RegisterRoutes registers tax endpoints on the given mux.
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/tax/preview", h.previewTax)
	mux.HandleFunc("GET /api/tax/exemptions/{customerID}", h.getExemptions)
	mux.HandleFunc("POST /api/tax/exemptions", h.createExemption)
	mux.HandleFunc("DELETE /api/tax/exemptions/{id}", h.deleteExemption)
}

func (h *Handler) previewTax(w http.ResponseWriter, r *http.Request) {
	var req TaxPreviewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	if len(req.Lines) == 0 {
		http.Error(w, `{"error":"at least one line item is required"}`, http.StatusBadRequest)
		return
	}

	if req.DocumentType == "" {
		req.DocumentType = "SalesInvoice"
	}

	result, err := h.svc.PreviewTax(r.Context(), &req)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (h *Handler) getExemptions(w http.ResponseWriter, r *http.Request) {
	customerIDStr := r.PathValue("customerID")
	customerID, err := uuid.Parse(customerIDStr)
	if err != nil {
		http.Error(w, `{"error":"invalid customer ID"}`, http.StatusBadRequest)
		return
	}

	exemptions, err := h.svc.GetExemptions(r.Context(), customerID)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	if exemptions == nil {
		exemptions = []TaxExemption{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(exemptions)
}

func (h *Handler) createExemption(w http.ResponseWriter, r *http.Request) {
	var req CreateExemptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	if req.CustomerID == uuid.Nil {
		http.Error(w, `{"error":"customer_id is required"}`, http.StatusBadRequest)
		return
	}
	if req.ExemptReason == "" {
		http.Error(w, `{"error":"exempt_reason is required"}`, http.StatusBadRequest)
		return
	}

	exemption, err := h.svc.SaveExemption(r.Context(), &req)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(exemption)
}

func (h *Handler) deleteExemption(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, `{"error":"invalid exemption ID"}`, http.StatusBadRequest)
		return
	}

	if err := h.svc.DeleteExemption(r.Context(), id); err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
