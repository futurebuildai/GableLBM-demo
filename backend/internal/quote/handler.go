package quote

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
	mux.HandleFunc("POST /quotes", h.HandleCreateQuote)
	mux.HandleFunc("GET /quotes", h.HandleListQuotes)
	mux.HandleFunc("GET /quotes/{id}", h.HandleGetQuotePath)
	mux.HandleFunc("POST /quotes/{id}/convert", h.HandleConvertToOrder)
}

func (h *Handler) HandleCreateQuote(w http.ResponseWriter, r *http.Request) {
	var q Quote
	if err := json.NewDecoder(r.Body).Decode(&q); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.service.CreateQuote(r.Context(), &q); err != nil {
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
