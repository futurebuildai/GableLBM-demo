package product

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
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
	mux.HandleFunc("GET /products/{id}", h.HandleGetProduct)
	mux.HandleFunc("PATCH /products/{id}/margins", h.HandleUpdateMarginRules)
}

// HandleGetProduct handles GET /products/{id}
func (h *Handler) HandleGetProduct(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid id format", http.StatusBadRequest)
		return
	}

	p, err := h.service.GetProduct(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
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

// HandleUpdateMarginRules handles PATCH /products/{id}/margins
func (h *Handler) HandleUpdateMarginRules(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid id format", http.StatusBadRequest)
		return
	}

	var req struct {
		TargetMargin   float64 `json:"target_margin"`
		CommissionRate float64 `json:"commission_rate"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.service.UpdateMarginRules(r.Context(), id, req.TargetMargin, req.CommissionRate); err != nil {
		http.Error(w, "Failed to update margin rules", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
