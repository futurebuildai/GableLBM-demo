package document

import (
	"context"
	"net/http"

	"github.com/gablelbm/gable/internal/customer"
	"github.com/gablelbm/gable/internal/invoice"
	"github.com/gablelbm/gable/internal/notification"
	"github.com/gablelbm/gable/internal/order"
	"github.com/google/uuid"
)

type Handler struct {
	docSvc      *Service
	orderSvc    *order.Service
	invoiceSvc  *invoice.Service
	customerSvc *customer.Service
	emailSvc    notification.EmailService
}

func NewHandler(d *Service, o *order.Service, i *invoice.Service, c *customer.Service, e notification.EmailService) *Handler {
	return &Handler{docSvc: d, orderSvc: o, invoiceSvc: i, customerSvc: c, emailSvc: e}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /documents/print/invoice/{id}", h.HandlePrintInvoice)
	mux.HandleFunc("GET /documents/print/pickticket/{id}", h.HandlePrintPickTicket)
	mux.HandleFunc("POST /api/invoices/{id}/email", h.HandleEmailInvoice)
}

func (h *Handler) HandlePrintInvoice(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	inv, err := h.invoiceSvc.GetInvoice(r.Context(), id)
	if err != nil {
		http.Error(w, "invoice not found", http.StatusNotFound)
		return
	}

	cust, err := h.customerSvc.GetCustomer(r.Context(), inv.CustomerID)
	if err != nil {
		http.Error(w, "customer not found", http.StatusNotFound)
		return
	}

	pdfBytes, err := h.docSvc.GenerateInvoicePDF(r.Context(), inv, cust)
	if err != nil {
		http.Error(w, "failed to generate pdf: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "inline; filename=invoice.pdf")
	w.Write(pdfBytes)
}

func (h *Handler) HandlePrintPickTicket(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	o, err := h.orderSvc.GetOrder(r.Context(), id)
	if err != nil {
		http.Error(w, "order not found", http.StatusNotFound)
		return
	}

	cust, err := h.customerSvc.GetCustomer(r.Context(), o.CustomerID)
	if err != nil {
		http.Error(w, "customer not found", http.StatusNotFound)
		return
	}

	pdfBytes, err := h.docSvc.GeneratePickTicketPDF(r.Context(), o, cust)
	if err != nil {
		http.Error(w, "failed to generate pdf: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "inline; filename=pickticket.pdf")
	w.Write(pdfBytes)
}

func (h *Handler) HandleEmailInvoice(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	inv, err := h.invoiceSvc.GetInvoice(r.Context(), id)
	if err != nil {
		http.Error(w, "invoice not found", http.StatusNotFound)
		return
	}

	cust, err := h.customerSvc.GetCustomer(r.Context(), inv.CustomerID)
	if err != nil {
		http.Error(w, "customer not found", http.StatusNotFound)
		return
	}

	pdfBytes, err := h.docSvc.GenerateInvoicePDF(r.Context(), inv, cust)
	if err != nil {
		http.Error(w, "failed to generate pdf", http.StatusInternalServerError)
		return
	}

	// In real app, look up customer email. For now, mock or use a query param
	email := "customer@example.com"
	// Async Email Dispatch
	// L8 Requirement: Do not block HTTP thread on external SMTP calls.
	go func() {
		// Create a background context or use valid timeout context
		bgCtx := context.Background()
		if err := h.emailSvc.SendInvoice(bgCtx, email, inv.ID.String(), pdfBytes); err != nil {
			// Log error (should inject logger here, but fmt for MVP)
			// fmt.Printf("Failed to send async email: %v\n", err)
			_ = err
		}
	}()

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(`{"status":"queued"}`))
}
