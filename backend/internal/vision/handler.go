package vision

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
	mux.HandleFunc("POST /api/vision/scan", h.handleScan)
}

func (h *Handler) handleScan(w http.ResponseWriter, r *http.Request) {
	var req BlueprintScanRequest
	// Limit body to 1MB to prevent DoS from large blueprint payloads
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20)).Decode(&req); err != nil {
		slog.Warn("Vision scan: invalid request body", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.BlueprintText == "" {
		http.Error(w, "blueprint_text is required", http.StatusBadRequest)
		return
	}

	resp := h.service.ScanBlueprint(req)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
