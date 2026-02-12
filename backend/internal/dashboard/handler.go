package dashboard

import (
	"encoding/json"
	"log/slog"
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

// writeJSON encodes data as JSON and writes it to the response.
func writeJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// writeError logs the internal error and returns a safe, generic message to the client.
func writeError(w http.ResponseWriter, msg string, err error, status int) {
	slog.Error(msg, "error", err)
	http.Error(w, msg, status)
}

// HandleSummary returns the dashboard summary KPIs.
func (h *Handler) HandleSummary(w http.ResponseWriter, r *http.Request) {
	data, err := h.service.GetSummary(r.Context())
	if err != nil {
		writeError(w, "Failed to load dashboard summary", err, http.StatusInternalServerError)
		return
	}
	writeJSON(w, data)
}

// HandleInventoryAlerts returns products with low stock.
func (h *Handler) HandleInventoryAlerts(w http.ResponseWriter, r *http.Request) {
	data, err := h.service.GetInventoryAlerts(r.Context())
	if err != nil {
		writeError(w, "Failed to load inventory alerts", err, http.StatusInternalServerError)
		return
	}
	writeJSON(w, data)
}

// HandleTopCustomers returns top revenue customers.
func (h *Handler) HandleTopCustomers(w http.ResponseWriter, r *http.Request) {
	data, err := h.service.GetTopCustomers(r.Context())
	if err != nil {
		writeError(w, "Failed to load top customers", err, http.StatusInternalServerError)
		return
	}
	writeJSON(w, data)
}

// HandleOrderActivity returns recent orders and status distribution.
func (h *Handler) HandleOrderActivity(w http.ResponseWriter, r *http.Request) {
	data, err := h.service.GetOrderActivity(r.Context())
	if err != nil {
		writeError(w, "Failed to load order activity", err, http.StatusInternalServerError)
		return
	}
	writeJSON(w, data)
}

// HandleRevenueTrend returns 7-day revenue trend data.
func (h *Handler) HandleRevenueTrend(w http.ResponseWriter, r *http.Request) {
	data, err := h.service.GetRevenueTrend(r.Context())
	if err != nil {
		writeError(w, "Failed to load revenue trend", err, http.StatusInternalServerError)
		return
	}
	writeJSON(w, data)
}
