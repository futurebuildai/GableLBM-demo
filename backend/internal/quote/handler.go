package quote

import (
	"encoding/base64"
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
	mux.HandleFunc("POST /quotes", h.HandleCreateQuote)
	mux.HandleFunc("GET /quotes/analytics", h.HandleGetAnalytics)
	mux.HandleFunc("GET /quotes", h.HandleListQuotes)
	mux.HandleFunc("GET /quotes/{id}", h.HandleGetQuotePath)
	mux.HandleFunc("GET /quotes/{id}/file", h.HandleDownloadOriginalFile)
	mux.HandleFunc("PUT /quotes/{id}/state", h.HandleUpdateState)
	mux.HandleFunc("POST /quotes/{id}/convert", h.HandleConvertToOrder)
}

// createQuoteRequest is the JSON payload for creating a quote.
// It mirrors Quote but accepts original_file as a base64 string.
type createQuoteRequest struct {
	Quote
	OriginalFileB64 string `json:"original_file,omitempty"`
}

func (h *Handler) HandleCreateQuote(w http.ResponseWriter, r *http.Request) {
	var req createQuoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	q := &req.Quote

	// Decode base64 original file if provided
	if req.OriginalFileB64 != "" {
		data, err := base64.StdEncoding.DecodeString(req.OriginalFileB64)
		if err == nil {
			q.OriginalFile = data
		}
	}

	if err := h.service.CreateQuote(r.Context(), q); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(q)
}

func (h *Handler) HandleGetQuotePath(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	q, err := h.service.GetQuote(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(q)
}

func (h *Handler) HandleListQuotes(w http.ResponseWriter, r *http.Request) {
	quotes, err := h.service.ListQuotes(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(quotes)
}

func (h *Handler) HandleConvertToOrder(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	q, err := h.service.GetQuote(r.Context(), id)
	if err != nil {
		http.Error(w, "Quote not found", http.StatusNotFound)
		return
	}

	if q.State != QuoteStateDraft && q.State != QuoteStateSent && q.State != QuoteStateAccepted {
		http.Error(w, "Quote cannot be converted in its current state", http.StatusBadRequest)
		return
	}

	// Build order creation payload from quote
	type OrderLinePayload struct {
		ProductID uuid.UUID `json:"product_id"`
		Quantity  float64   `json:"quantity"`
		PriceEach float64   `json:"price_each"`
	}
	type OrderPayload struct {
		CustomerID uuid.UUID          `json:"customer_id"`
		QuoteID    *uuid.UUID         `json:"quote_id"`
		Lines      []OrderLinePayload `json:"lines"`
	}

	payload := OrderPayload{
		CustomerID: q.CustomerID,
		QuoteID:    &q.ID,
	}
	for _, line := range q.Lines {
		payload.Lines = append(payload.Lines, OrderLinePayload{
			ProductID: line.ProductID,
			Quantity:  line.Quantity,
			PriceEach: line.UnitPrice,
		})
	}

	// Mark quote as ACCEPTED
	if err := h.service.UpdateState(r.Context(), id, QuoteStateAccepted); err != nil {
		http.Error(w, "Failed to update quote state: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the order payload - the frontend will POST it to /orders
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(payload)
}

func (h *Handler) HandleUpdateState(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	var body struct {
		State QuoteState `json:"state"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.service.UpdateState(r.Context(), id, body.State); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Return updated quote
	q, err := h.service.GetQuote(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(q)
}

func (h *Handler) HandleGetAnalytics(w http.ResponseWriter, r *http.Request) {
	analytics, err := h.service.GetAnalytics(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(analytics)
}

func (h *Handler) HandleDownloadOriginalFile(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	data, filename, contentType, err := h.service.GetOriginalFile(r.Context(), id)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	if len(data) == 0 {
		http.Error(w, "No original file stored for this quote", http.StatusNotFound)
		return
	}

	if filename == "" {
		filename = "original-upload"
	}

	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Disposition", "inline; filename=\""+filename+"\"")
	w.Write(data)
}
