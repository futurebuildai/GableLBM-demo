package pim

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

// Handler manages HTTP requests for PIM
type Handler struct {
	service *Service
}

// NewHandler creates a new PIM Handler
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// RegisterRoutes adds PIM handlers to the mux
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /products/{id}/detail", h.HandleGetProductDetail)
	mux.HandleFunc("GET /products/{id}/pim/content", h.HandleGetContent)
	mux.HandleFunc("PUT /products/{id}/pim/content", h.HandleUpdateContent)
	mux.HandleFunc("POST /products/{id}/pim/generate/descriptions", h.HandleGenerateDescriptions)
	mux.HandleFunc("POST /products/{id}/pim/generate/seo", h.HandleGenerateSEO)
	mux.HandleFunc("POST /products/{id}/pim/generate/image", h.HandleGenerateImage)
	mux.HandleFunc("POST /products/{id}/pim/generate/collateral", h.HandleGenerateCollateral)
	mux.HandleFunc("GET /products/{id}/pim/media", h.HandleListMedia)
	mux.HandleFunc("DELETE /products/{id}/pim/media/{mediaId}", h.HandleDeleteMedia)
	mux.HandleFunc("PATCH /products/{id}/pim/media/{mediaId}/primary", h.HandleSetPrimaryMedia)
	mux.HandleFunc("GET /products/{id}/pim/collateral", h.HandleListCollateral)
	mux.HandleFunc("DELETE /products/{id}/pim/collateral/{collateralId}", h.HandleDeleteCollateral)
}

func (h *Handler) HandleGetProductDetail(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid product id", http.StatusBadRequest)
		return
	}

	detail, err := h.service.GetProductDetail(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(detail)
}

func (h *Handler) HandleGetContent(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid product id", http.StatusBadRequest)
		return
	}

	content, err := h.service.GetContent(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if content == nil {
		content = &PIMContent{ProductID: id, Attributes: map[string]string{}, SEOKeywords: []string{}}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(content)
}

func (h *Handler) HandleUpdateContent(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid product id", http.StatusBadRequest)
		return
	}

	var req UpdateContentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	content, err := h.service.UpdateContent(r.Context(), id, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(content)
}

func (h *Handler) HandleGenerateDescriptions(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid product id", http.StatusBadRequest)
		return
	}

	var req GenerateDescriptionsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	content, err := h.service.GenerateDescriptions(r.Context(), id, req.Tone, req.Audience)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(content)
}

func (h *Handler) HandleGenerateSEO(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid product id", http.StatusBadRequest)
		return
	}

	var req GenerateSEORequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	content, err := h.service.GenerateSEO(r.Context(), id, req.TargetKeywords)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(content)
}

func (h *Handler) HandleGenerateImage(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid product id", http.StatusBadRequest)
		return
	}

	var req GenerateImageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	media, err := h.service.GenerateImage(r.Context(), id, req.Style, req.Prompt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(media)
}

func (h *Handler) HandleGenerateCollateral(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid product id", http.StatusBadRequest)
		return
	}

	var req GenerateCollateralRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	collateral, err := h.service.GenerateCollateral(r.Context(), id, req.Type, req.Tone, req.Audience)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(collateral)
}

func (h *Handler) HandleListMedia(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid product id", http.StatusBadRequest)
		return
	}

	media, err := h.service.ListMedia(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if media == nil {
		media = []PIMMedia{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(media)
}

func (h *Handler) HandleDeleteMedia(w http.ResponseWriter, r *http.Request) {
	mediaID, err := uuid.Parse(r.PathValue("mediaId"))
	if err != nil {
		http.Error(w, "invalid media id", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteMedia(r.Context(), mediaID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) HandleSetPrimaryMedia(w http.ResponseWriter, r *http.Request) {
	productID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid product id", http.StatusBadRequest)
		return
	}

	mediaID, err := uuid.Parse(r.PathValue("mediaId"))
	if err != nil {
		http.Error(w, "invalid media id", http.StatusBadRequest)
		return
	}

	if err := h.service.SetPrimaryMedia(r.Context(), productID, mediaID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) HandleListCollateral(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid product id", http.StatusBadRequest)
		return
	}

	items, err := h.service.ListCollateral(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if items == nil {
		items = []PIMCollateral{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func (h *Handler) HandleDeleteCollateral(w http.ResponseWriter, r *http.Request) {
	collateralID, err := uuid.Parse(r.PathValue("collateralId"))
	if err != nil {
		http.Error(w, "invalid collateral id", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteCollateral(r.Context(), collateralID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
