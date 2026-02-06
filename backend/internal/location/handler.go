package location

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
	// Simple routing for now. In a real app we might use a router lib.
	mux.HandleFunc("POST /locations", h.CreateLocation)
	mux.HandleFunc("GET /locations", h.ListLocations)
}

func (h *Handler) CreateLocation(w http.ResponseWriter, r *http.Request) {
	var loc Location
	if err := json.NewDecoder(r.Body).Decode(&loc); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if err := h.service.CreateLocation(r.Context(), &loc); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(loc)
}

func (h *Handler) ListLocations(w http.ResponseWriter, r *http.Request) {
	locs, err := h.service.ListLocations(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(locs)
}
