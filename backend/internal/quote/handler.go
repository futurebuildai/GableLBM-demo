package quote

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
	mux.HandleFunc("POST /quotes", h.HandleCreateQuote)
	mux.HandleFunc("GET /quotes", h.HandleGetQuote) // /quotes?id=... or use regex path?
	// Standard mux doesn't support parameterized paths well in Go < 1.22.
	// In Go 1.22+ "GET /quotes/{id}" works.
	// Assuming Go 1.22+ based on project freshness.
	mux.HandleFunc("GET /quotes/{id}", h.HandleGetQuotePath)
}

func (h *Handler) HandleCreateQuote(w http.ResponseWriter, r *http.Request) {
	var q Quote
	if err := json.NewDecoder(r.Body).Decode(&q); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.service.CreateQuote(r.Context(), &q); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(q)
}

// HandleGetQuotePath for Go 1.22+ routing
func (h *Handler) HandleGetQuotePath(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	q, err := h.service.GetQuote(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(q)
}

// Fallback if needed, but assuming 1.22
func (h *Handler) HandleGetQuote(w http.ResponseWriter, r *http.Request) {
	// Not implemented for query param
	http.Error(w, "Use /quotes/{id}", http.StatusBadRequest)
}
