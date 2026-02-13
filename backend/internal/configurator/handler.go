package configurator

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
	mux.HandleFunc("GET /api/configurator/rules", h.handleGetRules)
	mux.HandleFunc("POST /api/configurator/validate", h.handleValidate)
	mux.HandleFunc("POST /api/configurator/build-sku", h.handleBuildSKU)
	mux.HandleFunc("GET /api/configurator/options", h.handleGetOptions)
	mux.HandleFunc("GET /api/configurator/presets", h.handleGetPresets)
}

func (h *Handler) handleGetRules(w http.ResponseWriter, r *http.Request) {
	rules, err := h.service.GetAllRules(r.Context())
	if err != nil {
		slog.Error("Failed to fetch configurator rules", "error", err)
		http.Error(w, "Failed to fetch configurator rules", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rules)
}

func (h *Handler) handleValidate(w http.ResponseWriter, r *http.Request) {
	var req ValidateConfigRequest
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<16)).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if len(req.Selections) == 0 {
		http.Error(w, "Selections map is required", http.StatusBadRequest)
		return
	}

	resp, err := h.service.ValidateConfig(r.Context(), req)
	if err != nil {
		slog.Error("Validation failed", "error", err)
		http.Error(w, "Internal validation error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *Handler) handleBuildSKU(w http.ResponseWriter, r *http.Request) {
	var req BuildSKURequest
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<16)).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.ProductType == "" || len(req.Selections) == 0 {
		http.Error(w, "ProductType and Selections are required", http.StatusBadRequest)
		return
	}

	resp, err := h.service.BuildSKU(r.Context(), req)
	if err != nil {
		// BuildSKU returns user-facing validation errors, safe to surface
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *Handler) handleGetOptions(w http.ResponseWriter, r *http.Request) {
	attributeType := r.URL.Query().Get("attribute_type")
	if attributeType == "" {
		http.Error(w, "attribute_type query parameter is required", http.StatusBadRequest)
		return
	}

	// Parse optional selections from query params
	selections := make(map[string]string)
	for key, values := range r.URL.Query() {
		if key != "attribute_type" && len(values) > 0 {
			selections[key] = values[0]
		}
	}

	req := AvailableOptionsRequest{
		AttributeType: attributeType,
		Selections:    selections,
	}

	options, err := h.service.GetAvailableOptions(r.Context(), req)
	if err != nil {
		slog.Error("Failed to fetch options", "error", err)
		http.Error(w, "Failed to fetch options", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(options)
}

func (h *Handler) handleGetPresets(w http.ResponseWriter, r *http.Request) {
	productType := r.URL.Query().Get("product_type")

	presets, err := h.service.GetPresets(r.Context(), productType)
	if err != nil {
		slog.Error("Failed to fetch presets", "error", err)
		http.Error(w, "Failed to fetch presets", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(presets)
}
