package inventory

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/inventory/adjust", h.AdjustStock)
	mux.HandleFunc("POST /api/v1/inventory/transfer", h.MoveStock)
	mux.HandleFunc("GET /api/v1/inventory", h.ListInventory)
}

func (h *Handler) AdjustStock(w http.ResponseWriter, r *http.Request) {
	var req StockAdjustmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if err := h.service.AdjustStock(r.Context(), req); err != nil {
		slog.Error("AdjustStock failed", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}

func (h *Handler) MoveStock(w http.ResponseWriter, r *http.Request) {
	var req StockMovementRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if err := h.service.MoveStock(r.Context(), req); err != nil {
		slog.Error("MoveStock failed", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}

func (h *Handler) ListInventory(w http.ResponseWriter, r *http.Request) {
	prodID := r.URL.Query().Get("product_id")
	if prodID == "" {
		http.Error(w, "product_id required", http.StatusBadRequest)
		return
	}

	items, err := h.service.ListByProduct(r.Context(), prodID)
	if err != nil {
		slog.Error("ListByProduct failed", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}
