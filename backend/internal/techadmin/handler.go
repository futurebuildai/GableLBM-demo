package techadmin

import (
	"encoding/json"
	"net/http"

	"github.com/gablelbm/gable/internal/ai"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service      *Service
	aiKeyStore   *ai.KeyStore
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// WithAIKeyStore sets the AI key store for admin settings management.
func (h *Handler) WithAIKeyStore(ks *ai.KeyStore) {
	h.aiKeyStore = ks
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/admin/keys", h.CreateKey)
	mux.HandleFunc("GET /api/admin/keys", h.ListKeys)
	mux.HandleFunc("DELETE /api/admin/keys/{id}", h.RevokeKey)
	mux.HandleFunc("GET /api/admin/settings/ai", h.GetAISettings)
	mux.HandleFunc("PUT /api/admin/settings/ai", h.SaveAISettings)
	mux.HandleFunc("DELETE /api/admin/settings/ai", h.DeleteAISettings)
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

// --- AI Settings ---

type AISettingsResponse struct {
	Configured bool   `json:"configured"`
	Source     string `json:"source"` // "admin", "env", or "none"
	KeyHint   string `json:"key_hint,omitempty"` // e.g. "sk-ant-...4f2e"
}

func (h *Handler) GetAISettings(w http.ResponseWriter, r *http.Request) {
	if h.aiKeyStore == nil {
		json.NewEncoder(w).Encode(AISettingsResponse{Source: "none"})
		return
	}

	ctx := r.Context()
	key := h.aiKeyStore.Get(ctx)
	hasDB := h.aiKeyStore.HasDBOverride(ctx)

	resp := AISettingsResponse{
		Configured: key != "",
	}

	if key != "" {
		// Show a masked hint
		if len(key) > 12 {
			resp.KeyHint = key[:10] + "..." + key[len(key)-4:]
		} else {
			resp.KeyHint = "****"
		}

		if hasDB {
			resp.Source = "admin"
		} else {
			resp.Source = "env"
		}
	} else {
		resp.Source = "none"
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *Handler) SaveAISettings(w http.ResponseWriter, r *http.Request) {
	if h.aiKeyStore == nil {
		http.Error(w, "AI key store not available", http.StatusInternalServerError)
		return
	}

	var body struct {
		APIKey string `json:"api_key"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if body.APIKey == "" {
		http.Error(w, "api_key is required", http.StatusBadRequest)
		return
	}

	if err := h.aiKeyStore.Set(r.Context(), body.APIKey); err != nil {
		http.Error(w, "Failed to save API key: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "saved"})
}

func (h *Handler) DeleteAISettings(w http.ResponseWriter, r *http.Request) {
	if h.aiKeyStore == nil {
		http.Error(w, "AI key store not available", http.StatusInternalServerError)
		return
	}

	if err := h.aiKeyStore.Delete(r.Context()); err != nil {
		http.Error(w, "Failed to delete API key: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
