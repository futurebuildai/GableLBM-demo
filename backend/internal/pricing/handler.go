package pricing

import (
	"encoding/json"
	"net/http"

	"github.com/gablelbm/gable/internal/customer"
	"github.com/gablelbm/gable/internal/product"
	"github.com/google/uuid"
)

type Handler struct {
	service     *Service
	customerSvc *customer.Service
	productSvc  *product.Service
}

func NewHandler(s *Service, c *customer.Service, p *product.Service) *Handler {
	return &Handler{service: s, customerSvc: c, productSvc: p}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /pricing/calculate", h.HandleCalculatePrice)
}

func (h *Handler) HandleCalculatePrice(w http.ResponseWriter, r *http.Request) {
	customerIDStr := r.URL.Query().Get("customer_id")
	productIDStr := r.URL.Query().Get("product_id")

	if customerIDStr == "" || productIDStr == "" {
		http.Error(w, "customer_id and product_id are required", http.StatusBadRequest)
		return
	}

	customerID, err := uuid.Parse(customerIDStr)
	if err != nil {
		http.Error(w, "invalid customer_id", http.StatusBadRequest)
		return
	}

	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		http.Error(w, "invalid product_id", http.StatusBadRequest)
		return
	}

	cust, err := h.customerSvc.GetCustomer(r.Context(), customerID)
	if err != nil {
		http.Error(w, "failed to get customer", http.StatusNotFound)
		return
	}

	prod, err := h.productSvc.GetProduct(r.Context(), productID)
	if err != nil {
		http.Error(w, "failed to get product", http.StatusNotFound)
		return
	}

	priceResult, err := h.service.CalculatePrice(r.Context(), cust, productID, prod.BasePrice)
	if err != nil {
		http.Error(w, "failed to calculate price: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(priceResult)
}
