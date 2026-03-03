package pricing

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type RebateHandler struct {
	service RebateService
}

func NewRebateHandler(s RebateService) *RebateHandler {
	return &RebateHandler{service: s}
}

func (h *RebateHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /pricing/rebates/programs", h.HandleCreateProgram)
	mux.HandleFunc("GET /pricing/rebates/programs", h.HandleListPrograms)
	mux.HandleFunc("GET /pricing/rebates/programs/{id}", h.HandleGetProgram)
	mux.HandleFunc("POST /pricing/rebates/programs/{id}/claims/calculate", h.HandleCalculateClaim)
	mux.HandleFunc("GET /pricing/rebates/programs/{id}/claims", h.HandleListClaims)
}

type createProgramRequest struct {
	Program RebateProgram `json:"program"`
	Tiers   []RebateTier  `json:"tiers"`
}

func (h *RebateHandler) HandleCreateProgram(w http.ResponseWriter, r *http.Request) {
	var req createProgramRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Program.VendorID == uuid.Nil || req.Program.Name == "" {
		http.Error(w, "vendor_id and name are required", http.StatusBadRequest)
		return
	}

	prog, err := h.service.CreateProgramWithTiers(r.Context(), &req.Program, req.Tiers)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(prog)
}

func (h *RebateHandler) HandleListPrograms(w http.ResponseWriter, r *http.Request) {
	var vendorID *uuid.UUID
	if vidStr := r.URL.Query().Get("vendor_id"); vidStr != "" {
		if vid, err := uuid.Parse(vidStr); err == nil {
			vendorID = &vid
		}
	}

	programs, err := h.service.ListPrograms(r.Context(), vendorID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(programs)
}

func (h *RebateHandler) HandleGetProgram(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	prog, err := h.service.GetProgramWithTiers(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if prog == nil {
		http.Error(w, "program not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(prog)
}

type calculateClaimRequest struct {
	PeriodStart time.Time `json:"period_start"`
	PeriodEnd   time.Time `json:"period_end"`
	MockVolume  int64     `json:"mock_volume"`
}

func (h *RebateHandler) HandleCalculateClaim(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var req calculateClaimRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	claim, err := h.service.CalculateClaim(r.Context(), id, req.PeriodStart, req.PeriodEnd, req.MockVolume)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(claim)
}

func (h *RebateHandler) HandleListClaims(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	claims, err := h.service.ListClaims(r.Context(), &id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(claims)
}
