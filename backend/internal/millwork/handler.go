package millwork

import (
	"encoding/json"
	"net/http"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/millwork/options", h.handleCreateOption)
	mux.HandleFunc("GET /api/millwork/options", h.handleGetOptions)
}

func (h *Handler) handleCreateOption(w http.ResponseWriter, r *http.Request) {
	var req CreateOptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	opt, err := h.service.CreateOption(r.Context(), req)
	if err != nil {
		http.Error(w, "Failed to create option", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(opt)
}

func (h *Handler) handleGetOptions(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")
	if category == "" {
		http.Error(w, "Category query parameter is required", http.StatusBadRequest)
		return
	}

	options, err := h.service.GetOptionsByCategory(r.Context(), category)
	if err != nil {
		http.Error(w, "Failed to fetch options", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(options)
}
