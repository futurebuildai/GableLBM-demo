package ap

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

// Handler handles AP HTTP endpoints.
type Handler struct {
	service *Service
}

// NewHandler creates a new AP handler.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// RegisterRoutes registers AP API routes.
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	// Vendor Invoices
	mux.HandleFunc("POST /api/ap/invoices", h.CreateVendorInvoice)
	mux.HandleFunc("GET /api/ap/invoices", h.ListVendorInvoices)
	mux.HandleFunc("GET /api/ap/invoices/{id}", h.GetVendorInvoice)
	mux.HandleFunc("POST /api/ap/invoices/{id}/approve", h.ApproveInvoice)

	// AP Payments
	mux.HandleFunc("POST /api/ap/payments", h.PayVendor)
	mux.HandleFunc("GET /api/ap/payments", h.ListPayments)

	// Aging Report
	mux.HandleFunc("GET /api/ap/aging", h.GetAgingSummary)
}

func (h *Handler) CreateVendorInvoice(w http.ResponseWriter, r *http.Request) {
	var req CreateVendorInvoiceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	inv, err := h.service.CreateVendorInvoice(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(inv)
}

func (h *Handler) GetVendorInvoice(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid invoice ID", http.StatusBadRequest)
		return
	}

	inv, err := h.service.GetVendorInvoice(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(inv)
}

func (h *Handler) ListVendorInvoices(w http.ResponseWriter, r *http.Request) {
	var vendorID *uuid.UUID
	if vid := r.URL.Query().Get("vendor_id"); vid != "" {
		parsed, err := uuid.Parse(vid)
		if err == nil {
			vendorID = &parsed
		}
	}
	status := r.URL.Query().Get("status")

	invoices, err := h.service.ListVendorInvoices(r.Context(), vendorID, status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if invoices == nil {
		invoices = []VendorInvoice{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(invoices)
}

func (h *Handler) ApproveInvoice(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid invoice ID", http.StatusBadRequest)
		return
	}

	// In production, approverID comes from auth context
	approverID := uuid.New()

	inv, err := h.service.ApproveInvoice(r.Context(), id, approverID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(inv)
}

func (h *Handler) PayVendor(w http.ResponseWriter, r *http.Request) {
	var req CreateAPPaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	pmt, err := h.service.PayVendor(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(pmt)
}

func (h *Handler) ListPayments(w http.ResponseWriter, r *http.Request) {
	var vendorID *uuid.UUID
	if vid := r.URL.Query().Get("vendor_id"); vid != "" {
		parsed, err := uuid.Parse(vid)
		if err == nil {
			vendorID = &parsed
		}
	}

	payments, err := h.service.ListPayments(r.Context(), vendorID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if payments == nil {
		payments = []APPayment{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(payments)
}

func (h *Handler) GetAgingSummary(w http.ResponseWriter, r *http.Request) {
	summary, err := h.service.GetAgingSummary(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if summary == nil {
		summary = []APAgingSummary{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(summary)
}
