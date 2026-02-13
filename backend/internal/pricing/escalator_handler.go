package pricing

import (
	"encoding/json"
	"net/http"
)

// EscalatorHandler handles HTTP requests for price escalation endpoints.
type EscalatorHandler struct {
	service *EscalatorService
}

// NewEscalatorHandler creates a new escalator handler.
func NewEscalatorHandler(s *EscalatorService) *EscalatorHandler {
	return &EscalatorHandler{service: s}
}

// RegisterRoutes registers the escalator API routes on the given mux.
func (h *EscalatorHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/pricing/calculate-escalation", h.HandleCalculateEscalation)
	mux.HandleFunc("GET /api/v1/market-indices", h.HandleListMarketIndices)
}

// HandleCalculateEscalation calculates future pricing based on escalation parameters.
func (h *EscalatorHandler) HandleCalculateEscalation(w http.ResponseWriter, r *http.Request) {
	// Cap request body size at 1MB
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var req EscalationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.BasePrice <= 0 {
		http.Error(w, "base_price must be positive", http.StatusBadRequest)
		return
	}
	if req.EscalationType == "" {
		http.Error(w, "escalation_type is required", http.StatusBadRequest)
		return
	}
	if req.EscalationType != EscalationPercentage && req.EscalationType != EscalationIndexDelta {
		http.Error(w, "escalation_type must be PERCENTAGE or INDEX_DELTA", http.StatusBadRequest)
		return
	}
	if req.EffectiveDate == "" || req.TargetDate == "" {
		http.Error(w, "effective_date and target_date are required", http.StatusBadRequest)
		return
	}

	result, err := h.service.CalculateEscalation(r.Context(), req)
	if err != nil {
		http.Error(w, "Failed to calculate escalation", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// HandleListMarketIndices returns all active market indices.
func (h *EscalatorHandler) HandleListMarketIndices(w http.ResponseWriter, r *http.Request) {
	indices, err := h.service.ListMarketIndices(r.Context())
	if err != nil {
		http.Error(w, "Failed to list market indices", http.StatusInternalServerError)
		return
	}

	if indices == nil {
		indices = []MarketIndex{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(indices)
}
