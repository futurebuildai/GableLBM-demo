package partner

import (
	"encoding/json"
	"net/http"

	"github.com/gablelbm/gable/internal/customer"
	"github.com/gablelbm/gable/pkg/middleware"
	"github.com/google/uuid"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux, authMw func(http.Handler) http.Handler) {
	mux.Handle("GET /api/partner/v1/dashboard", authMw(http.HandlerFunc(h.GetDashboard)))
	mux.Handle("GET /api/partner/v1/quotes", authMw(http.HandlerFunc(h.ListQuotes)))
	mux.Handle("GET /api/partner/v1/quotes/{id}", authMw(http.HandlerFunc(h.GetQuote)))
}

func (h *Handler) GetDashboard(w http.ResponseWriter, r *http.Request) {
	cust := r.Context().Value(middleware.CustomerContextKey).(*customer.Customer)

	dto, err := h.svc.GetDashboard(r.Context(), cust.ID)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dto)
}

func (h *Handler) ListQuotes(w http.ResponseWriter, r *http.Request) {
	cust := r.Context().Value(middleware.CustomerContextKey).(*customer.Customer)

	quotes, err := h.svc.ListQuotes(r.Context(), cust.ID)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(quotes)
}

func (h *Handler) GetQuote(w http.ResponseWriter, r *http.Request) {
	cust := r.Context().Value(middleware.CustomerContextKey).(*customer.Customer)

	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	q, err := h.svc.GetQuote(r.Context(), cust.ID, id)
	if err != nil {
		http.Error(w, "Quote not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(q)
}
