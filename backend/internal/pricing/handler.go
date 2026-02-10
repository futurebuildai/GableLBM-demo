package pricing

import (
	"encoding/json"
	"net/http"
	"strconv"

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
	mux.HandleFunc("POST /pricing/rules", h.HandleCreateRule)
	mux.HandleFunc("GET /pricing/rules", h.HandleListRules)
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

	// Optional: quantity for volume pricing
	quantity := 1.0
	if qtyStr := r.URL.Query().Get("quantity"); qtyStr != "" {
		if q, err := strconv.ParseFloat(qtyStr, 64); err == nil && q > 0 {
			quantity = q
		}
	}

	// Optional: job_id for job-level pricing
	var jobID *uuid.UUID
	if jobIDStr := r.URL.Query().Get("job_id"); jobIDStr != "" {
		if jid, err := uuid.Parse(jobIDStr); err == nil {
			jobID = &jid
		}
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

	priceResult, err := h.service.CalculatePriceWithQty(r.Context(), cust, productID, prod.BasePrice, quantity, jobID)
	if err != nil {
		http.Error(w, "failed to calculate price: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(priceResult)
}

func (h *Handler) HandleCreateRule(w http.ResponseWriter, r *http.Request) {
	var rule PricingRule
	if err := json.NewDecoder(r.Body).Decode(&rule); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if rule.Name == "" || rule.RuleType == "" {
		http.Error(w, "name and rule_type are required", http.StatusBadRequest)
		return
	}

	if err := h.service.CreateRule(r.Context(), &rule); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(rule)
}

func (h *Handler) HandleListRules(w http.ResponseWriter, r *http.Request) {
	rules, err := h.service.ListRules(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rules)
}
