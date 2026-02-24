package matching

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

// Handler handles PO matching HTTP endpoints.
type Handler struct {
	service *Service
}

// NewHandler creates a new matching handler.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// RegisterRoutes registers matching API routes.
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/matching/run/{po_id}", h.RunMatch)
	mux.HandleFunc("GET /api/matching/results/{po_id}", h.GetMatchResult)
	mux.HandleFunc("GET /api/matching/exceptions", h.ListExceptions)
	mux.HandleFunc("GET /api/matching/config", h.GetConfig)
	mux.HandleFunc("PUT /api/matching/config", h.UpdateConfig)
}

func (h *Handler) RunMatch(w http.ResponseWriter, r *http.Request) {
	poID, err := uuid.Parse(r.PathValue("po_id"))
	if err != nil {
		http.Error(w, "Invalid PO ID", http.StatusBadRequest)
		return
	}

	result, err := h.service.RunMatch(r.Context(), poID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (h *Handler) GetMatchResult(w http.ResponseWriter, r *http.Request) {
	poID, err := uuid.Parse(r.PathValue("po_id"))
	if err != nil {
		http.Error(w, "Invalid PO ID", http.StatusBadRequest)
		return
	}

	result, err := h.service.GetMatchResult(r.Context(), poID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (h *Handler) ListExceptions(w http.ResponseWriter, r *http.Request) {
	exceptions, err := h.service.ListExceptions(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if exceptions == nil {
		exceptions = []MatchException{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(exceptions)
}

func (h *Handler) GetConfig(w http.ResponseWriter, r *http.Request) {
	cfg, err := h.service.GetConfig(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cfg)
}

func (h *Handler) UpdateConfig(w http.ResponseWriter, r *http.Request) {
	var req UpdateMatchConfigRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	cfg, err := h.service.UpdateConfig(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cfg)
}
