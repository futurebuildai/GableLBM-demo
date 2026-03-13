package parsing

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
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

	// Normalize content type for spreadsheets (DetectContentType may not identify xlsx properly)
	filename := header.Filename
	if contentType == "application/octet-stream" || contentType == "application/zip" {
		switch {
		case strings.HasSuffix(strings.ToLower(filename), ".xlsx") || strings.HasSuffix(strings.ToLower(filename), ".xls"):
			contentType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
		case strings.HasSuffix(strings.ToLower(filename), ".csv"):
			contentType = "text/csv"
		}
	}

	// For spreadsheets, convert to text before sending to AI
	aiContentType := contentType
	aiBytes := imageBytes
	if strings.Contains(contentType, "spreadsheet") || strings.HasSuffix(strings.ToLower(filename), ".xlsx") {
		textContent, convErr := convertSpreadsheetToText(imageBytes)
		if convErr != nil {
			slog.Error("Failed to convert spreadsheet", "error", convErr)
			http.Error(w, "Failed to process spreadsheet", http.StatusBadRequest)
			return
		}
		aiBytes = []byte(textContent)
		aiContentType = "text/plain"
	}

	// Extract items using AI (or rule-based fallback)
	extracted, extractErr := h.service.ExtractItemsWithAI(r.Context(), aiBytes, aiContentType)
	if extractErr != nil {
		slog.Error("Failed to extract items", "error", extractErr)
		http.Error(w, "Failed to extract items from file", http.StatusInternalServerError)
		return
	}

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

// convertSpreadsheetToText reads an xlsx file and converts it to plain text
// suitable for AI extraction.
func convertSpreadsheetToText(data []byte) (string, error) {
	f, err := excelize.OpenReader(bytes.NewReader(data))
	if err != nil {
		return "", fmt.Errorf("failed to open spreadsheet: %w", err)
	}
	defer f.Close()

	var sb strings.Builder
	sheets := f.GetSheetList()
	for _, sheet := range sheets {
		rows, err := f.GetRows(sheet)
		if err != nil {
			continue
		}
		for _, row := range rows {
			line := strings.Join(row, "\t")
			if strings.TrimSpace(line) != "" {
				sb.WriteString(line)
				sb.WriteString("\n")
			}
		}
	}
	return sb.String(), nil
}
