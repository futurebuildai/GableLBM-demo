package dashboard

import (
	"encoding/json"
	"net/http"
)

// Handler provides HTTP handlers for dashboard endpoints.
type Handler struct {
	service *Service
}

// NewHandler creates a new dashboard handler.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// RegisterRoutes registers all dashboard routes.
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/dashboard/summary", h.HandleSummary)
	mux.HandleFunc("GET /api/v1/dashboard/inventory-alerts", h.HandleInventoryAlerts)
	mux.HandleFunc("GET /api/v1/dashboard/top-customers", h.HandleTopCustomers)
	mux.HandleFunc("GET /api/v1/dashboard/order-activity", h.HandleOrderActivity)
	mux.HandleFunc("GET /api/v1/dashboard/revenue-trend", h.HandleRevenueTrend)
}

// HandleSummary returns the dashboard summary KPIs.
func (h *Handler) HandleSummary(w http.ResponseWriter, r *http.Request) {
	data, err := h.service.GetSummary(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// HandleInventoryAlerts returns products with low stock.
func (h *Handler) HandleInventoryAlerts(w http.ResponseWriter, r *http.Request) {
	data, err := h.service.GetInventoryAlerts(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// HandleTopCustomers returns top revenue customers.
func (h *Handler) HandleTopCustomers(w http.ResponseWriter, r *http.Request) {
	data, err := h.service.GetTopCustomers(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// HandleOrderActivity returns recent orders and status distribution.
func (h *Handler) HandleOrderActivity(w http.ResponseWriter, r *http.Request) {
	data, err := h.service.GetOrderActivity(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// HandleRevenueTrend returns 7-day revenue trend data.
func (h *Handler) HandleRevenueTrend(w http.ResponseWriter, r *http.Request) {
	data, err := h.service.GetRevenueTrend(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
