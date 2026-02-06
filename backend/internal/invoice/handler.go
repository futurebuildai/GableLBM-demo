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
		http.Error(w, err.Error(), http.StatusNotFound) // Simplify error handling
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(inv)
}
