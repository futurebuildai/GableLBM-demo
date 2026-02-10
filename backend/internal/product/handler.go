package product

import (
	"encoding/json"
	"net/http"
)

// Handler manages HTTP requests for products
type Handler struct {
	service *Service
}

// NewHandler creates a new Product Handler
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// RegisterRoutes adds handlers to the mux
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /products", h.HandleListProducts)
	mux.HandleFunc("POST /products", h.HandleCreateProduct)
	mux.HandleFunc("GET /products/reorder-alerts", h.HandleReorderAlerts)
}

// HandleCreateProduct handles POST /products
func (h *Handler) HandleCreateProduct(w http.ResponseWriter, r *http.Request) {
	var p Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.service.CreateProduct(r.Context(), &p); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(p)
}

// HandleReorderAlerts handles GET /products/reorder-alerts
func (h *Handler) HandleReorderAlerts(w http.ResponseWriter, r *http.Request) {
	alerts, err := h.service.ListBelowReorder(r.Context())
	if err != nil {
		http.Error(w, "Failed to fetch reorder alerts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(alerts)
}

// HandleListProducts handles GET /products
func (h *Handler) HandleListProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.service.ListProducts(r.Context())
	if err != nil {
		http.Error(w, "Failed to fetch products", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}
