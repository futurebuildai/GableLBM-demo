package purchase_order

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

type Handler struct {
	service *Service
	recSvc  *RecommendationService
}

func NewHandler(service *Service, recSvc *RecommendationService) *Handler {
	return &Handler{service: service, recSvc: recSvc}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /purchase-orders", h.HandleListPOs)
	mux.HandleFunc("POST /purchase-orders", h.HandleCreatePO)
	mux.HandleFunc("GET /purchase-orders/recommendations", h.HandleGetRecommendations)
	mux.HandleFunc("GET /purchase-orders/{id}", h.HandleGetPO)
	mux.HandleFunc("POST /purchase-orders/{id}/submit", h.HandleSubmitPO)
	mux.HandleFunc("POST /purchase-orders/{id}/receive", h.HandleReceivePO)
	mux.HandleFunc("POST /purchase-orders/reorder-check", h.HandleCreateReorders)
	mux.HandleFunc("POST /purchase-orders/{id}/freight", h.HandleUploadFreight)
	mux.HandleFunc("POST /purchase-orders/{id}/freight/{freightId}/apply", h.HandleApplyFreight)
	mux.HandleFunc("GET /purchase-orders/{id}/freight", h.HandleListFreight)
}

func (h *Handler) HandleListPOs(w http.ResponseWriter, r *http.Request) {
	pos, err := h.service.ListPOs(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pos)
}

type CreatePORequest struct {
	VendorID string         `json:"vendor_id"`
	Lines    []CreatePOLine `json:"lines"`
}

type CreatePOLine struct {
	ProductID   string  `json:"product_id"`
	Description string  `json:"description"`
	Quantity    float64 `json:"quantity"`
	Cost        float64 `json:"cost"`
}

func (h *Handler) HandleCreatePO(w http.ResponseWriter, r *http.Request) {
	var req CreatePORequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	vendorID, err := uuid.Parse(req.VendorID)
	if err != nil {
		http.Error(w, "Invalid vendor_id", http.StatusBadRequest)
		return
	}

	lines := make([]CreatePOLineInput, len(req.Lines))
	for i, l := range req.Lines {
		lines[i] = CreatePOLineInput{
			ProductID:   l.ProductID,
			Description: l.Description,
			Quantity:    l.Quantity,
			Cost:        l.Cost,
		}
	}

	po, err := h.service.CreateManualPOFromHandler(r.Context(), vendorID, lines)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(po)
}

func (h *Handler) HandleGetPO(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	po, err := h.service.GetPO(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(po)
}

func (h *Handler) HandleSubmitPO(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	if err := h.service.SubmitPO(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "submitted"})
}

type ReceiveLineRequest struct {
	LineID      string  `json:"line_id"`
	QtyReceived float64 `json:"qty_received"`
	LocationID  string  `json:"location_id"`
}

type ReceivePORequest struct {
	Lines []ReceiveLineRequest `json:"lines"`
}

func (h *Handler) HandleReceivePO(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	var req ReceivePORequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	lines := make([]ReceiveLineInput, len(req.Lines))
	for i, l := range req.Lines {
		lines[i] = ReceiveLineInput{
			LineID:      l.LineID,
			QtyReceived: l.QtyReceived,
			LocationID:  l.LocationID,
		}
	}

	if err := h.service.ReceivePO(r.Context(), id, lines); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "received"})
}

func (h *Handler) HandleCreateReorders(w http.ResponseWriter, r *http.Request) {
	count, err := h.service.CreateReorders(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"count":  count,
	})
}

// HandleGetRecommendations returns AI-driven purchasing recommendations
// based on sales velocity, stock levels, and lead times.
func (h *Handler) HandleGetRecommendations(w http.ResponseWriter, r *http.Request) {
	if h.recSvc == nil {
		http.Error(w, "Recommendation service not configured", http.StatusServiceUnavailable)
		return
	}

	summary, err := h.recSvc.GenerateRecommendations(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(summary)
}

// HandleUploadFreight processes a freight invoice upload for a received PO.
func (h *Handler) HandleUploadFreight(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	poID, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid PO ID", http.StatusBadRequest)
		return
	}

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

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Failed to read uploaded file", http.StatusInternalServerError)
		return
	}

	contentType := http.DetectContentType(fileBytes)

	slog.Info("Freight invoice upload",
		"po_id", poID,
		"filename", header.Filename,
		"size_bytes", header.Size,
		"content_type", contentType,
	)

	result, err := h.service.UploadFreightInvoice(r.Context(), poID, fileBytes, contentType, header.Filename)
	if err != nil {
		slog.Error("UploadFreightInvoice failed", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// HandleApplyFreight applies a pending freight charge to product costs.
func (h *Handler) HandleApplyFreight(w http.ResponseWriter, r *http.Request) {
	poID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid PO ID", http.StatusBadRequest)
		return
	}

	freightID, err := uuid.Parse(r.PathValue("freightId"))
	if err != nil {
		http.Error(w, "Invalid freight charge ID", http.StatusBadRequest)
		return
	}

	if err := h.service.ApplyFreightCharge(r.Context(), poID, freightID); err != nil {
		slog.Error("ApplyFreightCharge failed", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "applied"})
}

// HandleListFreight returns all freight charges for a PO.
func (h *Handler) HandleListFreight(w http.ResponseWriter, r *http.Request) {
	poID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid PO ID", http.StatusBadRequest)
		return
	}

	charges, err := h.service.GetFreightCharges(r.Context(), poID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if charges == nil {
		charges = []FreightCharge{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(charges)
}
