package portal

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/gablelbm/gable/pkg/middleware"
	"github.com/google/uuid"
)

// maxBodySize is the maximum request body size (1MB).
const maxBodySize = 1 << 20

// Handler provides HTTP handlers for portal endpoints.
type Handler struct {
	svc *Service
}

// NewHandler creates a new portal handler.
func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// writeJSON encodes data as JSON and writes it to the response.
func portalWriteJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// writeError logs the internal error and returns a safe, generic message to the client.
func portalWriteError(w http.ResponseWriter, msg string, err error, status int) {
	slog.Error(msg, "error", err)
	http.Error(w, msg, status)
}

// RegisterRoutes registers all portal API routes.
// Public routes (login, config) are registered directly on the mux.
// Protected routes are wrapped with portal auth middleware.
func (h *Handler) RegisterRoutes(mux *http.ServeMux, authMw func(http.Handler) http.Handler) {
	// Public endpoints
	mux.HandleFunc("POST /api/portal/v1/login", h.HandleLogin)
	mux.HandleFunc("GET /api/portal/v1/config", h.HandleGetConfig)

	// Protected endpoints
	mux.Handle("GET /api/portal/v1/dashboard", authMw(http.HandlerFunc(h.HandleDashboard)))
	mux.Handle("GET /api/portal/v1/orders", authMw(http.HandlerFunc(h.HandleListOrders)))
	mux.Handle("GET /api/portal/v1/orders/{id}", authMw(http.HandlerFunc(h.HandleGetOrder)))
	mux.Handle("POST /api/portal/v1/orders/reorder", authMw(http.HandlerFunc(h.HandleReorder)))
	mux.Handle("GET /api/portal/v1/invoices", authMw(http.HandlerFunc(h.HandleListInvoices)))
	mux.Handle("GET /api/portal/v1/invoices/{id}", authMw(http.HandlerFunc(h.HandleGetInvoice)))
	mux.Handle("GET /api/portal/v1/deliveries", authMw(http.HandlerFunc(h.HandleListDeliveries)))
	mux.Handle("GET /api/portal/v1/deliveries/{id}", authMw(http.HandlerFunc(h.HandleGetDelivery)))

	// Catalog endpoints (Sprint 27)
	mux.Handle("GET /api/portal/v1/catalog", authMw(http.HandlerFunc(h.HandleListCatalog)))
	mux.Handle("GET /api/portal/v1/catalog/{id}", authMw(http.HandlerFunc(h.HandleGetCatalogProduct)))

	// Cart endpoints (Sprint 27)
	mux.Handle("GET /api/portal/v1/cart", authMw(http.HandlerFunc(h.HandleGetCart)))
	mux.Handle("POST /api/portal/v1/cart/items", authMw(http.HandlerFunc(h.HandleAddToCart)))
	mux.Handle("PUT /api/portal/v1/cart/items/{id}", authMw(http.HandlerFunc(h.HandleUpdateCartItem)))
	mux.Handle("DELETE /api/portal/v1/cart/items/{id}", authMw(http.HandlerFunc(h.HandleRemoveCartItem)))

	// Checkout endpoint (Sprint 27)
	mux.Handle("POST /api/portal/v1/checkout", authMw(http.HandlerFunc(h.HandleCheckout)))
}

// HandleLogin authenticates a contractor and returns JWT + config.
func (h *Handler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		portalWriteError(w, "Invalid request body", err, http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	resp, err := h.svc.Login(r.Context(), req)
	if err != nil {
		// Always return 401 for login failures — don't leak user existence
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	portalWriteJSON(w, resp)
}

// HandleGetConfig returns portal branding config (public).
func (h *Handler) HandleGetConfig(w http.ResponseWriter, r *http.Request) {
	cfg, err := h.svc.GetConfig(r.Context())
	if err != nil {
		portalWriteError(w, "Failed to load portal configuration", err, http.StatusInternalServerError)
		return
	}
	portalWriteJSON(w, cfg)
}

// HandleDashboard returns contractor dashboard data.
func (h *Handler) HandleDashboard(w http.ResponseWriter, r *http.Request) {
	customerID := getPortalCustomerID(r)
	data, err := h.svc.GetDashboard(r.Context(), customerID)
	if err != nil {
		portalWriteError(w, "Failed to load dashboard", err, http.StatusInternalServerError)
		return
	}
	portalWriteJSON(w, data)
}

// HandleListOrders returns order history for the customer.
func (h *Handler) HandleListOrders(w http.ResponseWriter, r *http.Request) {
	customerID := getPortalCustomerID(r)
	orders, err := h.svc.ListOrders(r.Context(), customerID)
	if err != nil {
		portalWriteError(w, "Failed to load orders", err, http.StatusInternalServerError)
		return
	}
	portalWriteJSON(w, orders)
}

// HandleGetOrder returns a single order for the customer.
func (h *Handler) HandleGetOrder(w http.ResponseWriter, r *http.Request) {
	customerID := getPortalCustomerID(r)
	idStr := r.PathValue("id")
	orderID, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	order, err := h.svc.GetOrder(r.Context(), orderID, customerID)
	if err != nil {
		portalWriteError(w, "Order not found", err, http.StatusNotFound)
		return
	}
	portalWriteJSON(w, order)
}

// HandleReorder creates a new draft order from a historical order.
func (h *Handler) HandleReorder(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)
	customerID := getPortalCustomerID(r)

	var req ReorderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		portalWriteError(w, "Invalid request body", err, http.StatusBadRequest)
		return
	}

	resp, err := h.svc.CreateReorder(r.Context(), customerID, req)
	if err != nil {
		portalWriteError(w, "Failed to create reorder", err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	portalWriteJSON(w, resp)
}

// HandleListInvoices returns invoices for the customer.
func (h *Handler) HandleListInvoices(w http.ResponseWriter, r *http.Request) {
	customerID := getPortalCustomerID(r)
	invoices, err := h.svc.ListInvoices(r.Context(), customerID)
	if err != nil {
		portalWriteError(w, "Failed to load invoices", err, http.StatusInternalServerError)
		return
	}
	portalWriteJSON(w, invoices)
}

// HandleGetInvoice returns a single invoice for the customer.
func (h *Handler) HandleGetInvoice(w http.ResponseWriter, r *http.Request) {
	customerID := getPortalCustomerID(r)
	idStr := r.PathValue("id")
	invoiceID, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid invoice ID", http.StatusBadRequest)
		return
	}

	inv, err := h.svc.GetInvoice(r.Context(), invoiceID, customerID)
	if err != nil {
		portalWriteError(w, "Invoice not found", err, http.StatusNotFound)
		return
	}
	portalWriteJSON(w, inv)
}

// HandleListDeliveries returns deliveries for the customer.
func (h *Handler) HandleListDeliveries(w http.ResponseWriter, r *http.Request) {
	customerID := getPortalCustomerID(r)
	deliveries, err := h.svc.ListDeliveries(r.Context(), customerID)
	if err != nil {
		portalWriteError(w, "Failed to load deliveries", err, http.StatusInternalServerError)
		return
	}
	portalWriteJSON(w, deliveries)
}

// HandleGetDelivery returns a single delivery for the customer.
func (h *Handler) HandleGetDelivery(w http.ResponseWriter, r *http.Request) {
	customerID := getPortalCustomerID(r)
	idStr := r.PathValue("id")
	deliveryID, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid delivery ID", http.StatusBadRequest)
		return
	}

	del, err := h.svc.GetDelivery(r.Context(), deliveryID, customerID)
	if err != nil {
		portalWriteError(w, "Delivery not found", err, http.StatusNotFound)
		return
	}
	portalWriteJSON(w, del)
}

// getPortalCustomerID extracts the customer UUID from the request context.
// The middleware guarantees this is present on protected routes.
func getPortalCustomerID(r *http.Request) uuid.UUID {
	claims, ok := r.Context().Value(middleware.PortalClaimsKey).(*middleware.PortalClaims)
	if !ok || claims == nil {
		return uuid.Nil
	}
	return claims.CustomerID
}

// --- Catalog Handlers (Sprint 27) ---

// HandleListCatalog returns the product catalog with customer-specific pricing.
func (h *Handler) HandleListCatalog(w http.ResponseWriter, r *http.Request) {
	customerID := getPortalCustomerID(r)
	filter := CatalogFilter{
		Query:    r.URL.Query().Get("q"),
		Category: r.URL.Query().Get("category"),
		Species:  r.URL.Query().Get("species"),
		Grade:    r.URL.Query().Get("grade"),
	}

	products, err := h.svc.ListCatalog(r.Context(), customerID, filter)
	if err != nil {
		portalWriteError(w, "Failed to load catalog", err, http.StatusInternalServerError)
		return
	}
	portalWriteJSON(w, products)
}

// HandleGetCatalogProduct returns a single product with customer-specific pricing.
func (h *Handler) HandleGetCatalogProduct(w http.ResponseWriter, r *http.Request) {
	customerID := getPortalCustomerID(r)
	idStr := r.PathValue("id")
	productID, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	detail, err := h.svc.GetCatalogProduct(r.Context(), customerID, productID)
	if err != nil {
		portalWriteError(w, "Product not found", err, http.StatusNotFound)
		return
	}
	portalWriteJSON(w, detail)
}

// --- Cart Handlers (Sprint 27) ---

// HandleGetCart returns the current customer's cart.
func (h *Handler) HandleGetCart(w http.ResponseWriter, r *http.Request) {
	customerID := getPortalCustomerID(r)
	cart, err := h.svc.GetCart(r.Context(), customerID)
	if err != nil {
		portalWriteError(w, "Failed to load cart", err, http.StatusInternalServerError)
		return
	}
	portalWriteJSON(w, cart)
}

// HandleAddToCart adds an item to the customer's cart.
func (h *Handler) HandleAddToCart(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)
	customerID := getPortalCustomerID(r)

	var req AddToCartRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		portalWriteError(w, "Invalid request body", err, http.StatusBadRequest)
		return
	}

	cart, err := h.svc.AddToCart(r.Context(), customerID, req)
	if err != nil {
		portalWriteError(w, "Failed to add to cart", err, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	portalWriteJSON(w, cart)
}

// HandleUpdateCartItem updates a cart item quantity.
func (h *Handler) HandleUpdateCartItem(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)
	customerID := getPortalCustomerID(r)
	idStr := r.PathValue("id")
	itemID, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid item ID", http.StatusBadRequest)
		return
	}

	var req UpdateCartItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		portalWriteError(w, "Invalid request body", err, http.StatusBadRequest)
		return
	}

	cart, err := h.svc.UpdateCartItem(r.Context(), customerID, itemID, req)
	if err != nil {
		portalWriteError(w, "Failed to update cart item", err, http.StatusInternalServerError)
		return
	}
	portalWriteJSON(w, cart)
}

// HandleRemoveCartItem removes an item from the cart.
func (h *Handler) HandleRemoveCartItem(w http.ResponseWriter, r *http.Request) {
	customerID := getPortalCustomerID(r)
	idStr := r.PathValue("id")
	itemID, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid item ID", http.StatusBadRequest)
		return
	}

	cart, err := h.svc.RemoveCartItem(r.Context(), customerID, itemID)
	if err != nil {
		portalWriteError(w, "Failed to remove cart item", err, http.StatusInternalServerError)
		return
	}
	portalWriteJSON(w, cart)
}

// HandleCheckout places an order from the current cart.
func (h *Handler) HandleCheckout(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)
	customerID := getPortalCustomerID(r)

	var req CheckoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		portalWriteError(w, "Invalid request body", err, http.StatusBadRequest)
		return
	}

	resp, err := h.svc.Checkout(r.Context(), customerID, req)
	if err != nil {
		portalWriteError(w, "Checkout failed", err, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	portalWriteJSON(w, resp)
}
