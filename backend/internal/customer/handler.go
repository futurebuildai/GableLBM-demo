package customer

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
	mux.HandleFunc("GET /customers", h.HandleListCustomers)
	mux.HandleFunc("POST /customers", h.HandleCreateCustomer)
	mux.HandleFunc("GET /price_levels", h.HandleListPriceLevels)
}

func (h *Handler) HandleCreateCustomer(w http.ResponseWriter, r *http.Request) {
	var c Customer
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.service.CreateCustomer(r.Context(), &c); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(c)
}

func (h *Handler) HandleListCustomers(w http.ResponseWriter, r *http.Request) {
	// Simple list for now, no query params
	customers, err := h.service.ListCustomers(r.Context())
	if err != nil {
		http.Error(w, "Failed to fetch customers", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(customers)
}

func (h *Handler) HandleListPriceLevels(w http.ResponseWriter, r *http.Request) {
	levels, err := h.service.ListPriceLevels(r.Context())
	if err != nil {
		http.Error(w, "Failed to fetch price levels", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(levels)
}
