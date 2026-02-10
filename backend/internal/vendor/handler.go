package vendor

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
	mux.HandleFunc("GET /vendors", h.HandleList)
	mux.HandleFunc("POST /vendors", h.HandleCreate)
	mux.HandleFunc("GET /vendors/{id}", h.HandleGet)
}

func (h *Handler) HandleList(w http.ResponseWriter, r *http.Request) {
	vendors, err := h.service.ListVendors(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(vendors)
}

func (h *Handler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	var req CreateVendorRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	v, err := h.service.CreateVendor(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(v)
}

func (h *Handler) HandleGet(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	v, err := h.service.GetVendor(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if v == nil {
		http.Error(w, "Vendor not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(v)
}
