package techadmin

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/admin/keys", h.CreateKey)
	mux.HandleFunc("GET /api/admin/keys", h.ListKeys)
	mux.HandleFunc("DELETE /api/admin/keys/{id}", h.RevokeKey)
}

type CreateKeyRequest struct {
	Name   string   `json:"name"`
	Scopes []string `json:"scopes"`
}

type CreateKeyResponse struct {
	APIKey string  `json:"api_key"` // The raw key, shown once
	Key    *APIKey `json:"key"`
}

func (h *Handler) CreateKey(w http.ResponseWriter, r *http.Request) {
	var req CreateKeyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	rawKey, apiKey, err := h.service.GenerateKey(r.Context(), req.Name, req.Scopes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(CreateKeyResponse{
		APIKey: rawKey,
		Key:    apiKey,
	})
}

func (h *Handler) ListKeys(w http.ResponseWriter, r *http.Request) {
	keys, err := h.service.ListKeys(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(keys)
}

func (h *Handler) RevokeKey(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	if err := h.service.RevokeKey(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
