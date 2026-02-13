package parsing

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
)

// Handler manages HTTP requests for material list parsing.
type Handler struct {
	service *Service
}

// NewHandler creates a new parsing Handler.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// RegisterRoutes adds parsing handlers to the mux.
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /parsing/upload", h.HandleUpload)
}

// HandleUpload processes a material list image upload and returns parsed items.
// Accepts multipart/form-data with a "file" field.
func (h *Handler) HandleUpload(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	// Limit upload to 10MB
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "File too large or invalid form data", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Missing 'file' field in form data", http.StatusBadRequest)
		return
	}
	defer file.Close()

	slog.Info("Parsing material list upload",
		"filename", header.Filename,
		"size_bytes", header.Size,
	)

	// Read uploaded file bytes for the base64 preview
	imageBytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Failed to read uploaded file", http.StatusInternalServerError)
		return
	}

	// Detect content type for data URI
	contentType := http.DetectContentType(imageBytes)

	// --- Simulated AI Text Extraction ---
	// In a real implementation, the image bytes would be sent to an AI vision model.
	// For the MVP, we use a hardcoded sample material list that demonstrates all
	// confidence tiers (high, medium, special order).
	sampleMaterialList := generateSampleMaterialList()

	// Extract items using rule-based parser
	extracted := h.service.ExtractItems(sampleMaterialList)

	// Match against product catalog
	items, err := h.service.MatchProducts(r.Context(), extracted)
	if err != nil {
		slog.Error("Failed to match products", "error", err)
		http.Error(w, "Failed to process material list", http.StatusInternalServerError)
		return
	}

	// Build base64 data URI for the uploaded image
	sourceImage := fmt.Sprintf("data:%s;base64,%s",
		contentType,
		base64.StdEncoding.EncodeToString(imageBytes),
	)

	resp := ParseResponse{
		Items:       items,
		SourceImage: sourceImage,
		ParseTimeMs: time.Since(start).Milliseconds(),
		ItemCount:   len(items),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// generateSampleMaterialList returns a realistic handwritten material list
// that exercises all confidence tiers when matched against a typical LBM catalog.
func generateSampleMaterialList() string {
	return `50 pcs - 2x4x8 SPF Stud
25 pcs - 2x6x12 Doug Fir #2
30 sheets - OSB 7/16 4x8
10 sheets - CDX Plywood 1/2 4x8
15 pcs - 2x10x16 Hem Fir
8 pcs - 2x12x20.
20 bags - Quikrete 80lb
4 rolls - Tyvek House Wrap
100 lf - 2x4 Pressure Treated
Custom powder-coat railing 12ft bronze
6 pcs - Simpson Strong-Tie A35
Specialty glass panel 48x72 frosted`
}
