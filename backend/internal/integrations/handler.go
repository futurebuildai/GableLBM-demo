package integrations

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gablelbm/gable/internal/customer"
	"github.com/gablelbm/gable/internal/order"
	"github.com/gablelbm/gable/internal/pricing"
	"github.com/gablelbm/gable/internal/product"
	"github.com/gablelbm/gable/internal/quote"
	"github.com/gablelbm/gable/internal/reporting"
	"github.com/gablelbm/gable/pkg/database"
	"github.com/google/uuid"
)

type Handler struct {
	db           *database.DB
	pricingSvc   *pricing.Service
	quoteSvc     *quote.Service
	orderSvc     *order.Service
	customerSvc  *customer.Service
	productSvc   *product.Service
	reportingSvc *reporting.Service
	apiKey       string
}

func NewHandler(db *database.DB, pricingSvc *pricing.Service, quoteSvc *quote.Service, orderSvc *order.Service, customerSvc *customer.Service, productSvc *product.Service, reportingSvc *reporting.Service) *Handler {
	apiKey := os.Getenv("INTEGRATION_API_KEY")
	if apiKey == "" {
		apiKey = "fb-brain-demo-key-2026"
	}
	return &Handler{
		db:           db,
		pricingSvc:   pricingSvc,
		quoteSvc:     quoteSvc,
		orderSvc:     orderSvc,
		customerSvc:  customerSvc,
		productSvc:   productSvc,
		reportingSvc: reportingSvc,
		apiKey:       apiKey,
	}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/integration/products", h.authMiddleware(h.ListProductsByCategory))
	mux.HandleFunc("POST /api/integration/quotes/bulk-price", h.authMiddleware(h.BulkCalculatePrice))
	mux.HandleFunc("POST /api/integration/quotes", h.authMiddleware(h.CreateQuote))
	mux.HandleFunc("POST /api/integration/quotes/{id}/accept-and-convert", h.authMiddleware(h.AcceptAndConvertQuote))

	// Read-only endpoints for FB-Brain sync
	mux.HandleFunc("GET /api/integration/customers", h.authMiddleware(h.ListCustomers))
	mux.HandleFunc("GET /api/integration/customers/{id}", h.authMiddleware(h.GetCustomer))
	mux.HandleFunc("GET /api/integration/orders", h.authMiddleware(h.ListOrders))
	mux.HandleFunc("GET /api/integration/orders/{id}", h.authMiddleware(h.GetOrderDetail))
	mux.HandleFunc("GET /api/integration/invoices", h.authMiddleware(h.ListInvoices))
	mux.HandleFunc("GET /api/integration/invoices/{id}", h.authMiddleware(h.GetInvoiceDetail))
	mux.HandleFunc("GET /api/integration/deliveries", h.authMiddleware(h.ListDeliveries))
	mux.HandleFunc("GET /api/integration/statements", h.authMiddleware(h.ListStatements))
	mux.HandleFunc("GET /api/integration/payments", h.authMiddleware(h.ListPayments))

	// Write-back endpoints (Velocity → FB-Brain → GableLBM)
	mux.HandleFunc("POST /api/integration/invoices/{id}/payment", h.authMiddleware(h.RecordPayment))
	mux.HandleFunc("PUT /api/integration/invoices/{id}/status", h.authMiddleware(h.UpdateInvoiceStatus))

	// Demo lifecycle
	mux.HandleFunc("POST /api/integration/reset", h.authMiddleware(h.HandleReset))
}

func (h *Handler) authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get("X-Integration-Key")
		if key != h.apiKey {
			writeError(w, http.StatusUnauthorized, "invalid integration key")
			return
		}
		next(w, r)
	}
}

// ProductResponse is the integration-facing product model
type ProductResponse struct {
	ID       string  `json:"id"`
	SKU      string  `json:"sku"`
	Name     string  `json:"name"`
	Category string  `json:"category"`
	UOM      string  `json:"uom"`
	Price    int64   `json:"price"` // cents
}

// ListProductsByCategory returns products filtered by category and/or text search
func (h *Handler) ListProductsByCategory(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")
	query := r.URL.Query().Get("q")

	if category == "" && query == "" {
		writeError(w, http.StatusBadRequest, "category or q query parameter required")
		return
	}

	sqlQuery := `SELECT p.id, p.sku, p.description, COALESCE(p.category, ''), p.uom_primary::text, COALESCE(p.base_price, 0)
		FROM products p WHERE 1=1`
	args := []interface{}{}
	argIdx := 1

	if category != "" {
		sqlQuery += fmt.Sprintf(` AND p.category = $%d`, argIdx)
		args = append(args, category)
		argIdx++
	}
	if query != "" {
		sqlQuery += fmt.Sprintf(` AND (p.sku ILIKE $%d OR p.description ILIKE $%d)`, argIdx, argIdx)
		args = append(args, "%"+query+"%")
		argIdx++
	}
	sqlQuery += ` ORDER BY p.sku LIMIT 20`

	rows, err := h.db.Pool.Query(r.Context(), sqlQuery, args...)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	var products []ProductResponse
	for rows.Next() {
		var p ProductResponse
		var priceFloat float64
		if err := rows.Scan(&p.ID, &p.SKU, &p.Name, &p.Category, &p.UOM, &priceFloat); err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		p.Price = int64(priceFloat * 100)
		products = append(products, p)
	}

	writeJSON(w, http.StatusOK, products)
}

// BulkPriceRequest is the request body for bulk pricing
type BulkPriceRequest struct {
	CustomerID string          `json:"customer_id"`
	Items      []BulkPriceItem `json:"items"`
}

type BulkPriceItem struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

type PricedItemResponse struct {
	ProductID   string `json:"product_id"`
	ProductName string `json:"product_name"`
	SKU         string `json:"sku"`
	Quantity    int    `json:"quantity"`
	UnitPrice   int64  `json:"unit_price"`  // cents
	TotalPrice  int64  `json:"total_price"` // cents
	UOM         string `json:"uom"`
}

// BulkCalculatePrice calculates prices for multiple items for a specific customer
func (h *Handler) BulkCalculatePrice(w http.ResponseWriter, r *http.Request) {
	var req BulkPriceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	customerID, err := uuid.Parse(req.CustomerID)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid customer_id")
		return
	}

	cust, err := h.customerSvc.GetCustomer(r.Context(), customerID)
	if err != nil {
		writeError(w, http.StatusNotFound, fmt.Sprintf("customer not found: %v", err))
		return
	}

	var results []PricedItemResponse
	for _, item := range req.Items {
		productID, err := uuid.Parse(item.ProductID)
		if err != nil {
			continue
		}

		prod, err := h.productSvc.GetProduct(r.Context(), productID)
		if err != nil {
			continue
		}

		calculated, err := h.pricingSvc.CalculatePriceWithQty(r.Context(), cust, productID, prod.BasePrice, float64(item.Quantity), nil)
		if err != nil {
			continue
		}

		unitPriceCents := int64(calculated.FinalPrice * 100)
		totalPriceCents := unitPriceCents * int64(item.Quantity)

		results = append(results, PricedItemResponse{
			ProductID:   item.ProductID,
			ProductName: prod.Description,
			SKU:         prod.SKU,
			Quantity:    item.Quantity,
			UnitPrice:   unitPriceCents,
			TotalPrice:  totalPriceCents,
			UOM:         string(prod.UOMPrimary),
		})
	}

	writeJSON(w, http.StatusOK, results)
}

// CreateQuoteRequest is the request body for creating a quote
type CreateQuoteRequest struct {
	CustomerID string           `json:"customer_id"`
	Lines      []QuoteLineInput `json:"lines"`
}

type QuoteLineInput struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
	UnitPrice int64  `json:"unit_price"` // cents
}

type QuoteResponse struct {
	ID         string `json:"id"`
	CustomerID string `json:"customer_id"`
	Total      int64  `json:"total"` // cents
	Status     string `json:"status"`
}

// CreateQuote creates a DRAFT quote from pre-priced line items
func (h *Handler) CreateQuote(w http.ResponseWriter, r *http.Request) {
	var req CreateQuoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	customerID, err := uuid.Parse(req.CustomerID)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid customer_id")
		return
	}

	// Build quote lines
	var lines []quote.QuoteLine
	for _, line := range req.Lines {
		productID, err := uuid.Parse(line.ProductID)
		if err != nil {
			continue
		}

		prod, err := h.productSvc.GetProduct(r.Context(), productID)
		if err != nil {
			continue
		}

		unitPriceDollars := float64(line.UnitPrice) / 100.0
		lines = append(lines, quote.QuoteLine{
			ProductID:   productID,
			SKU:         prod.SKU,
			Description: prod.Description,
			Quantity:    float64(line.Quantity),
			UOM:         prod.UOMPrimary,
			UnitPrice:   unitPriceDollars,
		})
	}

	demoCreatedBy := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	expires := time.Now().AddDate(0, 0, 30)

	q := &quote.Quote{
		CustomerID: customerID,
		State:      quote.QuoteStateDraft,
		ExpiresAt:  &expires,
		Lines:      lines,
		Source:     "integration",
	}
	// Set CreatedBy via context or field - the service will handle totals
	_ = demoCreatedBy

	if err := h.quoteSvc.CreateQuote(r.Context(), q); err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Sprintf("create quote: %v", err))
		return
	}

	totalCents := int64(q.TotalAmount * 100)

	writeJSON(w, http.StatusCreated, QuoteResponse{
		ID:         q.ID.String(),
		CustomerID: req.CustomerID,
		Total:      totalCents,
		Status:     string(q.State),
	})
}

type OrderResponse struct {
	ID      string `json:"id"`
	QuoteID string `json:"quote_id"`
	Status  string `json:"status"`
}

// AcceptAndConvertQuote accepts a quote and converts it to an order
func (h *Handler) AcceptAndConvertQuote(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	quoteID, err := uuid.Parse(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid quote id")
		return
	}

	ctx := r.Context()

	// 1. Accept the quote
	if err := h.quoteSvc.UpdateState(ctx, quoteID, quote.QuoteStateAccepted); err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Sprintf("accept quote: %v", err))
		return
	}

	// 2. Get the quote to build order
	q, err := h.quoteSvc.GetQuote(ctx, quoteID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Sprintf("get quote: %v", err))
		return
	}

	// 3. Convert to order
	var orderLines []order.OrderLineRequest
	for _, ql := range q.Lines {
		orderLines = append(orderLines, order.OrderLineRequest{
			ProductID: ql.ProductID,
			Quantity:  ql.Quantity,
			PriceEach: ql.UnitPrice,
		})
	}

	o, err := h.orderSvc.CreateOrder(ctx, order.CreateOrderRequest{
		CustomerID: q.CustomerID,
		QuoteID:    &quoteID,
		Lines:      orderLines,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Sprintf("create order: %v", err))
		return
	}

	// Order stays DRAFT for dealer review in the ERP UI.

	writeJSON(w, http.StatusOK, OrderResponse{
		ID:      o.ID.String(),
		QuoteID: quoteID.String(),
		Status:  string(o.Status),
	})
}

func (h *Handler) confirmOrder(ctx context.Context, orderID uuid.UUID) error {
	return h.orderSvc.ConfirmOrder(ctx, orderID)
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

// ── Integration Read Endpoints (FB-Brain sync) ─────────────────────────────

// IntegrationCustomerResponse is the integration-facing customer model
type IntegrationCustomerResponse struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	AccountNumber string `json:"account_number"`
	Email         string `json:"email"`
	Phone         string `json:"phone"`
	Address       string `json:"address"`
	Tier          string `json:"tier"`
	CreditLimit   int64  `json:"credit_limit"` // cents
	BalanceDue    int64  `json:"balance_due"`   // cents
	IsActive      bool   `json:"is_active"`
}

// ListCustomers returns customers with optional filters: ?active=true, ?account_number=X
func (h *Handler) ListCustomers(w http.ResponseWriter, r *http.Request) {
	query := `SELECT id, name, account_number, COALESCE(email, ''), COALESCE(phone, ''),
		COALESCE(address, ''), COALESCE(tier::text, 'RETAIL'), COALESCE(credit_limit, 0), COALESCE(balance_due, 0), is_active
		FROM customers WHERE 1=1`
	args := []interface{}{}
	argIdx := 1

	if v := r.URL.Query().Get("active"); v != "" {
		query += fmt.Sprintf(` AND is_active = $%d`, argIdx)
		args = append(args, v == "true")
		argIdx++
	}
	if v := r.URL.Query().Get("account_number"); v != "" {
		query += fmt.Sprintf(` AND account_number = $%d`, argIdx)
		args = append(args, v)
		argIdx++
	}
	query += ` ORDER BY name LIMIT 200`

	rows, err := h.db.Pool.Query(r.Context(), query, args...)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	var customers []IntegrationCustomerResponse
	for rows.Next() {
		var c IntegrationCustomerResponse
		var creditFloat, balanceFloat float64
		if err := rows.Scan(&c.ID, &c.Name, &c.AccountNumber, &c.Email, &c.Phone,
			&c.Address, &c.Tier, &creditFloat, &balanceFloat, &c.IsActive); err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		c.CreditLimit = int64(creditFloat * 100)
		c.BalanceDue = int64(balanceFloat * 100)
		customers = append(customers, c)
	}

	writeJSON(w, http.StatusOK, customers)
}

// GetCustomer returns a single customer by ID
func (h *Handler) GetCustomer(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid customer id")
		return
	}

	query := `SELECT id, name, account_number, COALESCE(email, ''), COALESCE(phone, ''),
		COALESCE(address, ''), COALESCE(tier::text, 'RETAIL'), COALESCE(credit_limit, 0), COALESCE(balance_due, 0), is_active
		FROM customers WHERE id = $1`

	var c IntegrationCustomerResponse
	var creditFloat, balanceFloat float64
	err = h.db.Pool.QueryRow(r.Context(), query, id).Scan(
		&c.ID, &c.Name, &c.AccountNumber, &c.Email, &c.Phone,
		&c.Address, &c.Tier, &creditFloat, &balanceFloat, &c.IsActive)
	if err != nil {
		writeError(w, http.StatusNotFound, "customer not found")
		return
	}
	c.CreditLimit = int64(creditFloat * 100)
	c.BalanceDue = int64(balanceFloat * 100)

	writeJSON(w, http.StatusOK, c)
}

// IntegrationOrderResponse is the integration-facing order model
type IntegrationOrderResponse struct {
	ID                    string                       `json:"id"`
	CustomerID            string                       `json:"customer_id"`
	CustomerAccountNumber string                       `json:"customer_account_number"`
	QuoteID               *string                      `json:"quote_id,omitempty"`
	Status                string                       `json:"status"`
	TotalAmount           int64                        `json:"total_amount"` // cents
	CreatedAt             string                       `json:"created_at"`
	Lines                 []IntegrationOrderLineResult `json:"lines,omitempty"`
}

type IntegrationOrderLineResult struct {
	ProductID   string  `json:"product_id"`
	SKU         string  `json:"sku"`
	Description string  `json:"description"`
	Quantity    float64 `json:"quantity"`
	PriceEach   int64   `json:"price_each"` // cents
	UOM         string  `json:"uom"`
}

// ListOrders returns orders with optional filters: ?customer_id=X, ?status=X, ?since=ISO8601
func (h *Handler) ListOrders(w http.ResponseWriter, r *http.Request) {
	query := `SELECT o.id, o.customer_id, c.account_number, o.quote_id, o.status, o.total_amount, o.created_at
		FROM orders o
		LEFT JOIN customers c ON c.id = o.customer_id
		WHERE 1=1`
	args := []interface{}{}
	argIdx := 1

	if v := r.URL.Query().Get("customer_id"); v != "" {
		cid, err := uuid.Parse(v)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid customer_id")
			return
		}
		query += fmt.Sprintf(` AND o.customer_id = $%d`, argIdx)
		args = append(args, cid)
		argIdx++
	}
	if v := r.URL.Query().Get("status"); v != "" {
		query += fmt.Sprintf(` AND o.status = $%d`, argIdx)
		args = append(args, v)
		argIdx++
	}
	if v := r.URL.Query().Get("since"); v != "" {
		since, err := time.Parse(time.RFC3339, v)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid since (use ISO8601)")
			return
		}
		query += fmt.Sprintf(` AND o.created_at >= $%d`, argIdx)
		args = append(args, since)
		argIdx++
	}
	query += ` ORDER BY o.created_at DESC LIMIT 200`

	rows, err := h.db.Pool.Query(r.Context(), query, args...)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	var orders []IntegrationOrderResponse
	for rows.Next() {
		var o IntegrationOrderResponse
		var quoteID *uuid.UUID
		var totalFloat float64
		var createdAt time.Time
		if err := rows.Scan(&o.ID, &o.CustomerID, &o.CustomerAccountNumber, &quoteID, &o.Status, &totalFloat, &createdAt); err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		o.TotalAmount = int64(totalFloat * 100)
		o.CreatedAt = createdAt.Format(time.RFC3339)
		if quoteID != nil {
			s := quoteID.String()
			o.QuoteID = &s
		}
		orders = append(orders, o)
	}

	writeJSON(w, http.StatusOK, orders)
}

// GetOrderDetail returns a single order with its line items
func (h *Handler) GetOrderDetail(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid order id")
		return
	}

	// Fetch order header
	query := `SELECT o.id, o.customer_id, c.account_number, o.quote_id, o.status, o.total_amount, o.created_at
		FROM orders o
		LEFT JOIN customers c ON c.id = o.customer_id
		WHERE o.id = $1`

	var o IntegrationOrderResponse
	var quoteID *uuid.UUID
	var totalFloat float64
	var createdAt time.Time
	err = h.db.Pool.QueryRow(r.Context(), query, id).Scan(
		&o.ID, &o.CustomerID, &o.CustomerAccountNumber, &quoteID, &o.Status, &totalFloat, &createdAt)
	if err != nil {
		writeError(w, http.StatusNotFound, "order not found")
		return
	}
	o.TotalAmount = int64(totalFloat * 100)
	o.CreatedAt = createdAt.Format(time.RFC3339)
	if quoteID != nil {
		s := quoteID.String()
		o.QuoteID = &s
	}

	// Fetch lines
	linesQuery := `SELECT ol.product_id, COALESCE(p.sku, ''), COALESCE(p.description, ''),
		ol.quantity, ol.price_each, COALESCE(p.uom_primary::text, 'EACH')
		FROM order_lines ol
		LEFT JOIN products p ON p.id = ol.product_id
		WHERE ol.order_id = $1`

	rows, err := h.db.Pool.Query(r.Context(), linesQuery, id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	for rows.Next() {
		var l IntegrationOrderLineResult
		var priceFloat float64
		if err := rows.Scan(&l.ProductID, &l.SKU, &l.Description, &l.Quantity, &priceFloat, &l.UOM); err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		l.PriceEach = int64(priceFloat * 100)
		o.Lines = append(o.Lines, l)
	}

	writeJSON(w, http.StatusOK, o)
}

// IntegrationInvoiceResponse is the integration-facing invoice model
type IntegrationInvoiceResponse struct {
	ID                    string                         `json:"id"`
	OrderID               string                         `json:"order_id"`
	CustomerID            string                         `json:"customer_id"`
	CustomerAccountNumber string                         `json:"customer_account_number"`
	Status                string                         `json:"status"`
	Subtotal              int64                          `json:"subtotal"`     // cents
	TaxRate               float64                        `json:"tax_rate"`     // e.g. 800 = 8%
	TaxAmount             int64                          `json:"tax_amount"`   // cents
	TotalAmount           int64                          `json:"total_amount"` // cents
	PaymentTerms          string                         `json:"payment_terms"`
	DueDate               *string                        `json:"due_date,omitempty"`
	CreatedAt             string                         `json:"created_at"`
	Lines                 []IntegrationInvoiceLineResult `json:"lines,omitempty"`
}

type IntegrationInvoiceLineResult struct {
	ProductID   string  `json:"product_id"`
	SKU         string  `json:"sku"`
	Description string  `json:"description"`
	Quantity    float64 `json:"quantity"`
	PriceEach   int64   `json:"price_each"` // cents
}

// ListInvoices returns invoices with optional filters: ?customer_id=X, ?status=X, ?since=ISO8601
func (h *Handler) ListInvoices(w http.ResponseWriter, r *http.Request) {
	query := `SELECT i.id, i.order_id, i.customer_id, c.account_number, i.status,
		COALESCE(i.subtotal, i.total_amount), COALESCE(i.tax_rate, 0), COALESCE(i.tax_amount, 0),
		i.total_amount, COALESCE(i.payment_terms, 'NET30'), i.due_date, i.created_at
		FROM invoices i
		LEFT JOIN customers c ON c.id = i.customer_id
		WHERE 1=1`
	args := []interface{}{}
	argIdx := 1

	if v := r.URL.Query().Get("customer_id"); v != "" {
		cid, err := uuid.Parse(v)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid customer_id")
			return
		}
		query += fmt.Sprintf(` AND i.customer_id = $%d`, argIdx)
		args = append(args, cid)
		argIdx++
	}
	if v := r.URL.Query().Get("status"); v != "" {
		query += fmt.Sprintf(` AND i.status = $%d`, argIdx)
		args = append(args, v)
		argIdx++
	}
	if v := r.URL.Query().Get("since"); v != "" {
		since, err := time.Parse(time.RFC3339, v)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid since (use ISO8601)")
			return
		}
		query += fmt.Sprintf(` AND i.created_at >= $%d`, argIdx)
		args = append(args, since)
		argIdx++
	}
	query += ` ORDER BY i.created_at DESC LIMIT 200`

	rows, err := h.db.Pool.Query(r.Context(), query, args...)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	var invoices []IntegrationInvoiceResponse
	for rows.Next() {
		var inv IntegrationInvoiceResponse
		var subtotalFloat, taxRateFloat, taxAmountFloat, totalFloat float64
		var dueDate *time.Time
		var createdAt time.Time
		if err := rows.Scan(&inv.ID, &inv.OrderID, &inv.CustomerID, &inv.CustomerAccountNumber,
			&inv.Status, &subtotalFloat, &taxRateFloat, &taxAmountFloat, &totalFloat,
			&inv.PaymentTerms, &dueDate, &createdAt); err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		inv.Subtotal = int64(subtotalFloat * 100)
		inv.TaxRate = taxRateFloat * 10000 // 0.0825 -> 825
		inv.TaxAmount = int64(taxAmountFloat * 100)
		inv.TotalAmount = int64(totalFloat * 100)
		inv.CreatedAt = createdAt.Format(time.RFC3339)
		if dueDate != nil {
			s := dueDate.Format("2006-01-02")
			inv.DueDate = &s
		}
		invoices = append(invoices, inv)
	}

	writeJSON(w, http.StatusOK, invoices)
}

// GetInvoiceDetail returns a single invoice with its line items
func (h *Handler) GetInvoiceDetail(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid invoice id")
		return
	}

	query := `SELECT i.id, i.order_id, i.customer_id, c.account_number, i.status,
		COALESCE(i.subtotal, i.total_amount), COALESCE(i.tax_rate, 0), COALESCE(i.tax_amount, 0),
		i.total_amount, COALESCE(i.payment_terms, 'NET30'), i.due_date, i.created_at
		FROM invoices i
		LEFT JOIN customers c ON c.id = i.customer_id
		WHERE i.id = $1`

	var inv IntegrationInvoiceResponse
	var subtotalFloat, taxRateFloat, taxAmountFloat, totalFloat float64
	var dueDate *time.Time
	var createdAt time.Time
	err = h.db.Pool.QueryRow(r.Context(), query, id).Scan(
		&inv.ID, &inv.OrderID, &inv.CustomerID, &inv.CustomerAccountNumber,
		&inv.Status, &subtotalFloat, &taxRateFloat, &taxAmountFloat, &totalFloat,
		&inv.PaymentTerms, &dueDate, &createdAt)
	if err != nil {
		writeError(w, http.StatusNotFound, "invoice not found")
		return
	}
	inv.Subtotal = int64(subtotalFloat * 100)
	inv.TaxRate = taxRateFloat * 10000
	inv.TaxAmount = int64(taxAmountFloat * 100)
	inv.TotalAmount = int64(totalFloat * 100)
	inv.CreatedAt = createdAt.Format(time.RFC3339)
	if dueDate != nil {
		s := dueDate.Format("2006-01-02")
		inv.DueDate = &s
	}

	// Fetch lines
	linesQuery := `SELECT il.product_id, COALESCE(p.sku, ''), COALESCE(p.description, ''),
		il.quantity, il.price_each
		FROM invoice_lines il
		LEFT JOIN products p ON p.id = il.product_id
		WHERE il.invoice_id = $1`

	rows, err := h.db.Pool.Query(r.Context(), linesQuery, id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	for rows.Next() {
		var l IntegrationInvoiceLineResult
		var priceFloat float64
		if err := rows.Scan(&l.ProductID, &l.SKU, &l.Description, &l.Quantity, &priceFloat); err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		l.PriceEach = int64(priceFloat * 100)
		inv.Lines = append(inv.Lines, l)
	}

	writeJSON(w, http.StatusOK, inv)
}

// IntegrationDeliveryResponse is the integration-facing delivery model
type IntegrationDeliveryResponse struct {
	ID           string  `json:"id"`
	RouteID      string  `json:"route_id"`
	OrderID      string  `json:"order_id"`
	Status       string  `json:"status"`
	StopSequence int     `json:"stop_sequence"`
	PODSignedBy  *string `json:"pod_signed_by,omitempty"`
	PODProofURL  *string `json:"pod_proof_url,omitempty"`
	PODTimestamp *string `json:"pod_timestamp,omitempty"`
	CreatedAt    string  `json:"created_at"`
}

// ListDeliveries returns deliveries with optional filters: ?order_id=X, ?status=X
func (h *Handler) ListDeliveries(w http.ResponseWriter, r *http.Request) {
	query := `SELECT d.id, d.route_id, d.order_id, d.status, d.stop_sequence,
		d.pod_signed_by, d.pod_proof_url, d.pod_timestamp, d.created_at
		FROM deliveries d WHERE 1=1`
	args := []interface{}{}
	argIdx := 1

	if v := r.URL.Query().Get("order_id"); v != "" {
		oid, err := uuid.Parse(v)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid order_id")
			return
		}
		query += fmt.Sprintf(` AND d.order_id = $%d`, argIdx)
		args = append(args, oid)
		argIdx++
	}
	if v := r.URL.Query().Get("status"); v != "" {
		query += fmt.Sprintf(` AND d.status = $%d`, argIdx)
		args = append(args, v)
		argIdx++
	}
	query += ` ORDER BY d.created_at DESC LIMIT 200`

	rows, err := h.db.Pool.Query(r.Context(), query, args...)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	var deliveries []IntegrationDeliveryResponse
	for rows.Next() {
		var d IntegrationDeliveryResponse
		var podTimestamp *time.Time
		var createdAt time.Time
		if err := rows.Scan(&d.ID, &d.RouteID, &d.OrderID, &d.Status, &d.StopSequence,
			&d.PODSignedBy, &d.PODProofURL, &podTimestamp, &createdAt); err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		d.CreatedAt = createdAt.Format(time.RFC3339)
		if podTimestamp != nil {
			s := podTimestamp.Format(time.RFC3339)
			d.PODTimestamp = &s
		}
		deliveries = append(deliveries, d)
	}

	writeJSON(w, http.StatusOK, deliveries)
}

// ── Statement & Payment Read Endpoints ───────────────────────────────────────

// IntegrationStatementResponse is the integration-facing statement model
type IntegrationStatementResponse struct {
	CustomerID   string  `json:"customer_id"`
	CustomerName string  `json:"customer_name"`
	PeriodStart  string  `json:"period_start"`
	PeriodEnd    string  `json:"period_end"`
	OpenBalance  float64 `json:"open_balance"`
	CloseBalance float64 `json:"close_balance"`
}

// ListStatements returns monthly statements for a customer over the last 6 months
func (h *Handler) ListStatements(w http.ResponseWriter, r *http.Request) {
	customerID := r.URL.Query().Get("customer_id")
	if customerID == "" {
		writeError(w, http.StatusBadRequest, "customer_id required")
		return
	}
	if _, err := uuid.Parse(customerID); err != nil {
		writeError(w, http.StatusBadRequest, "invalid customer_id")
		return
	}

	var statements []IntegrationStatementResponse
	now := time.Now()

	for i := 0; i < 6; i++ {
		monthStart := time.Date(now.Year(), now.Month()-time.Month(i), 1, 0, 0, 0, 0, time.UTC)
		monthEnd := time.Date(now.Year(), now.Month()-time.Month(i)+1, 0, 0, 0, 0, 0, time.UTC) // last day of month

		stmt, err := h.reportingSvc.GetCustomerStatement(
			r.Context(), customerID,
			monthStart.Format("2006-01-02"),
			monthEnd.Format("2006-01-02"),
		)
		if err != nil {
			continue
		}

		statements = append(statements, IntegrationStatementResponse{
			CustomerID:   stmt.CustomerID,
			CustomerName: stmt.CustomerName,
			PeriodStart:  monthStart.Format("2006-01-02"),
			PeriodEnd:    monthEnd.Format("2006-01-02"),
			OpenBalance:  stmt.OpenBalance,
			CloseBalance: stmt.CloseBalance,
		})
	}

	if statements == nil {
		statements = []IntegrationStatementResponse{}
	}
	writeJSON(w, http.StatusOK, statements)
}

// IntegrationPaymentResponse is the integration-facing payment model
type IntegrationPaymentResponse struct {
	ID         string `json:"id"`
	InvoiceID  string `json:"invoice_id"`
	CustomerID string `json:"customer_id"`
	Amount     int64  `json:"amount"` // cents
	Method     string `json:"method"`
	Reference  string `json:"reference"`
	CardLast4  string `json:"card_last4,omitempty"`
	CardBrand  string `json:"card_brand,omitempty"`
	CreatedAt  string `json:"created_at"`
}

// ListPayments returns payments for a customer
func (h *Handler) ListPayments(w http.ResponseWriter, r *http.Request) {
	customerID := r.URL.Query().Get("customer_id")
	if customerID == "" {
		writeError(w, http.StatusBadRequest, "customer_id required")
		return
	}
	cid, err := uuid.Parse(customerID)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid customer_id")
		return
	}

	query := `
		SELECT p.id, p.invoice_id, i.customer_id, p.amount, p.method,
			   COALESCE(p.reference, ''), COALESCE(p.card_last4, ''),
			   COALESCE(p.card_brand, ''), p.created_at
		FROM payments p
		JOIN invoices i ON i.id = p.invoice_id
		WHERE i.customer_id = $1
		ORDER BY p.created_at DESC
		LIMIT 100
	`

	rows, err := h.db.Pool.Query(r.Context(), query, cid)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	var payments []IntegrationPaymentResponse
	for rows.Next() {
		var p IntegrationPaymentResponse
		var amountFloat float64
		var createdAt time.Time
		if err := rows.Scan(&p.ID, &p.InvoiceID, &p.CustomerID, &amountFloat,
			&p.Method, &p.Reference, &p.CardLast4, &p.CardBrand, &createdAt); err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		p.Amount = int64(amountFloat * 100)
		p.CreatedAt = createdAt.Format(time.RFC3339)
		payments = append(payments, p)
	}

	if payments == nil {
		payments = []IntegrationPaymentResponse{}
	}
	writeJSON(w, http.StatusOK, payments)
}

// ── Write-back Endpoints (reverse flow: Velocity → FB-Brain → GableLBM) ────

// RecordPaymentRequest is the body for POST /api/integration/invoices/{id}/payment
type RecordPaymentRequest struct {
	Amount    float64 `json:"amount"`    // dollars
	Method    string  `json:"method"`    // CASH, CARD, CHECK, ACCOUNT
	Reference string  `json:"reference"` // external ref (Velocity payment ID)
}

// RecordPayment records a payment against an invoice and updates status/balance
func (h *Handler) RecordPayment(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	invoiceID, err := uuid.Parse(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid invoice id")
		return
	}

	var req RecordPaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Amount <= 0 {
		writeError(w, http.StatusBadRequest, "amount must be positive")
		return
	}
	if req.Method == "" {
		req.Method = "ACCOUNT"
	}

	ctx := r.Context()

	// Get current invoice
	var totalFloat float64
	var currentStatus string
	var customerID uuid.UUID
	err = h.db.Pool.QueryRow(ctx,
		`SELECT total_amount, status, customer_id FROM invoices WHERE id = $1`, invoiceID,
	).Scan(&totalFloat, &currentStatus, &customerID)
	if err != nil {
		writeError(w, http.StatusNotFound, "invoice not found")
		return
	}

	// Sum existing payments
	var paidFloat float64
	_ = h.db.Pool.QueryRow(ctx,
		`SELECT COALESCE(SUM(amount), 0) FROM payments WHERE invoice_id = $1`, invoiceID,
	).Scan(&paidFloat)

	// Insert payment
	_, err = h.db.Pool.Exec(ctx,
		`INSERT INTO payments (invoice_id, amount, method, reference) VALUES ($1, $2, $3, $4)`,
		invoiceID, req.Amount, req.Method, req.Reference)
	if err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Sprintf("insert payment: %v", err))
		return
	}

	// Determine new status
	newPaid := paidFloat + req.Amount
	newStatus := currentStatus
	now := time.Now()
	var paidAt *time.Time
	if newPaid >= totalFloat {
		newStatus = "PAID"
		paidAt = &now
	} else if newPaid > 0 {
		newStatus = "PARTIAL"
	}

	// Update invoice status
	_, err = h.db.Pool.Exec(ctx,
		`UPDATE invoices SET status = $1, paid_at = $2, updated_at = NOW() WHERE id = $3`,
		newStatus, paidAt, invoiceID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Sprintf("update invoice: %v", err))
		return
	}

	// Update customer balance (reduce by payment amount)
	_, err = h.db.Pool.Exec(ctx,
		`UPDATE customers SET balance_due = balance_due - $1, updated_at = NOW() WHERE id = $2`,
		req.Amount, customerID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Sprintf("update balance: %v", err))
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"invoice_id": invoiceID.String(),
		"status":     newStatus,
		"amount":     req.Amount,
	})
}

// UpdateInvoiceStatusRequest is the body for PUT /api/integration/invoices/{id}/status
type UpdateInvoiceStatusRequest struct {
	Status string `json:"status"` // UNPAID, PARTIAL, PAID, VOID, OVERDUE
}

// UpdateInvoiceStatus directly sets the invoice status
func (h *Handler) UpdateInvoiceStatus(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	invoiceID, err := uuid.Parse(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid invoice id")
		return
	}

	var req UpdateInvoiceStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	validStatuses := map[string]bool{
		"UNPAID": true, "PARTIAL": true, "PAID": true, "VOID": true, "OVERDUE": true,
	}
	if !validStatuses[req.Status] {
		writeError(w, http.StatusBadRequest, "invalid status")
		return
	}

	var paidAt *time.Time
	if req.Status == "PAID" {
		now := time.Now()
		paidAt = &now
	}

	tag, err := h.db.Pool.Exec(r.Context(),
		`UPDATE invoices SET status = $1, paid_at = $2, updated_at = NOW() WHERE id = $3`,
		req.Status, paidAt, invoiceID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Sprintf("update invoice: %v", err))
		return
	}
	if tag.RowsAffected() == 0 {
		writeError(w, http.StatusNotFound, "invoice not found")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"invoice_id": invoiceID.String(),
		"status":     req.Status,
	})
}

// ── Demo Reset Endpoint ─────────────────────────────────────────────────────

// HandleReset deletes all integration-created quotes, orders, invoices, and
// payments so the demo can be run again from a clean state. It identifies
// integration records by quotes.source = 'integration' and walks the
// quote → order → invoice → payment chain.
func (h *Handler) HandleReset(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Use a transaction so the reset is atomic.
	tx, err := h.db.Pool.Begin(ctx)
	if err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Sprintf("begin tx: %v", err))
		return
	}
	defer tx.Rollback(ctx)

	// 1. Collect integration quote IDs
	var quoteIDs []uuid.UUID
	rows, err := tx.Query(ctx, `SELECT id FROM quotes WHERE source = 'integration'`)
	if err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Sprintf("find integration quotes: %v", err))
		return
	}
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			rows.Close()
			writeError(w, http.StatusInternalServerError, fmt.Sprintf("scan quote id: %v", err))
			return
		}
		quoteIDs = append(quoteIDs, id)
	}
	rows.Close()

	if len(quoteIDs) == 0 {
		_ = tx.Commit(ctx)
		writeJSON(w, http.StatusOK, map[string]interface{}{
			"reset":            true,
			"quotes_deleted":   0,
			"orders_deleted":   0,
			"invoices_deleted": 0,
			"payments_deleted": 0,
		})
		return
	}

	// 2. Collect orders that originated from those quotes
	var orderIDs []uuid.UUID
	rows, err = tx.Query(ctx, `SELECT id FROM orders WHERE quote_id = ANY($1)`, quoteIDs)
	if err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Sprintf("find orders: %v", err))
		return
	}
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			rows.Close()
			writeError(w, http.StatusInternalServerError, fmt.Sprintf("scan order id: %v", err))
			return
		}
		orderIDs = append(orderIDs, id)
	}
	rows.Close()

	var invoiceIDs []uuid.UUID
	var paymentsDeleted int64

	if len(orderIDs) > 0 {
		// 3. Collect invoices for those orders
		rows, err = tx.Query(ctx, `SELECT id FROM invoices WHERE order_id = ANY($1)`, orderIDs)
		if err != nil {
			writeError(w, http.StatusInternalServerError, fmt.Sprintf("find invoices: %v", err))
			return
		}
		for rows.Next() {
			var id uuid.UUID
			if err := rows.Scan(&id); err != nil {
				rows.Close()
				writeError(w, http.StatusInternalServerError, fmt.Sprintf("scan invoice id: %v", err))
				return
			}
			invoiceIDs = append(invoiceIDs, id)
		}
		rows.Close()

		if len(invoiceIDs) > 0 {
			// 4a. Delete payments
			tag, err := tx.Exec(ctx,
				`DELETE FROM payments WHERE invoice_id = ANY($1)`, invoiceIDs)
			if err != nil {
				writeError(w, http.StatusInternalServerError, fmt.Sprintf("delete payments: %v", err))
				return
			}
			paymentsDeleted = tag.RowsAffected()

			// 4b. Delete invoice lines
			_, err = tx.Exec(ctx,
				`DELETE FROM invoice_lines WHERE invoice_id = ANY($1)`, invoiceIDs)
			if err != nil {
				writeError(w, http.StatusInternalServerError, fmt.Sprintf("delete invoice lines: %v", err))
				return
			}

			// 4c. Delete invoices
			_, err = tx.Exec(ctx,
				`DELETE FROM invoices WHERE id = ANY($1)`, invoiceIDs)
			if err != nil {
				writeError(w, http.StatusInternalServerError, fmt.Sprintf("delete invoices: %v", err))
				return
			}
		}

		// 5a. Delete deliveries for those orders (if any exist)
		_, _ = tx.Exec(ctx,
			`DELETE FROM deliveries WHERE order_id = ANY($1)`, orderIDs)

		// 5b. Delete order lines
		_, err = tx.Exec(ctx,
			`DELETE FROM order_lines WHERE order_id = ANY($1)`, orderIDs)
		if err != nil {
			writeError(w, http.StatusInternalServerError, fmt.Sprintf("delete order lines: %v", err))
			return
		}

		// 5c. Delete orders
		_, err = tx.Exec(ctx,
			`DELETE FROM orders WHERE id = ANY($1)`, orderIDs)
		if err != nil {
			writeError(w, http.StatusInternalServerError, fmt.Sprintf("delete orders: %v", err))
			return
		}
	}

	// 6a. Delete quote lines
	_, err = tx.Exec(ctx,
		`DELETE FROM quote_lines WHERE quote_id = ANY($1)`, quoteIDs)
	if err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Sprintf("delete quote lines: %v", err))
		return
	}

	// 6b. Delete quotes
	_, err = tx.Exec(ctx,
		`DELETE FROM quotes WHERE id = ANY($1)`, quoteIDs)
	if err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Sprintf("delete quotes: %v", err))
		return
	}

	if err := tx.Commit(ctx); err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Sprintf("commit: %v", err))
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"reset":            true,
		"quotes_deleted":   len(quoteIDs),
		"orders_deleted":   len(orderIDs),
		"invoices_deleted": len(invoiceIDs),
		"payments_deleted": paymentsDeleted,
	})
}
