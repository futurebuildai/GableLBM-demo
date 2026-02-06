package governance

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/governance/rfcs", h.HandleCreateRFC)
	mux.HandleFunc("GET /api/v1/governance/rfcs", h.HandleListRFCs)
	mux.HandleFunc("GET /api/v1/governance/rfcs/{id}", h.HandleGetRFC)
	mux.HandleFunc("PUT /api/v1/governance/rfcs/{id}", h.HandleUpdateRFC)
}

func (h *Handler) HandleCreateRFC(w http.ResponseWriter, r *http.Request) {
	var input CreateRFCInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	rfc, err := h.service.DraftRFC(r.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(rfc)
}

func (h *Handler) HandleListRFCs(w http.ResponseWriter, r *http.Request) {
	rfcs, err := h.service.ListRFCs(r.Context())
	if err != nil {
		http.Error(w, "Failed to fetch RFCs", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rfcs)
}

func (h *Handler) HandleGetRFC(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid UUID", http.StatusBadRequest)
		return
	}

	rfc, err := h.service.GetRFC(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound) // Using 404 roughly for error mostly
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rfc)
}

func (h *Handler) HandleUpdateRFC(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid UUID", http.StatusBadRequest)
		return
	}

	var input UpdateRFCInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	rfc, err := h.service.UpdateRFC(r.Context(), id, input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rfc)
}
