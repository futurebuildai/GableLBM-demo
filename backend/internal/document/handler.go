package document

import (
	"net/http"

	"github.com/gablelbm/gable/internal/customer"
	"github.com/gablelbm/gable/internal/invoice"
	"github.com/gablelbm/gable/internal/order"
	"github.com/google/uuid"
)

type Handler struct {
	docSvc      *Service
	orderSvc    *order.Service
	invoiceSvc  *invoice.Service
	customerSvc *customer.Service
}

func NewHandler(d *Service, o *order.Service, i *invoice.Service, c *customer.Service) *Handler {
	return &Handler{docSvc: d, orderSvc: o, invoiceSvc: i, customerSvc: c}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /documents/print/invoice/{id}", h.HandlePrintInvoice)
	mux.HandleFunc("GET /documents/print/pickticket/{id}", h.HandlePrintPickTicket)
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
