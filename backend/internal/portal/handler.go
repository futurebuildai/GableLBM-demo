package portal

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gablelbm/gable/pkg/middleware"
	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
)

// maxBodySize is the maximum request body size (1MB) for normal endpoints.
const maxBodySize = 1 << 20

// maxQuoteBodySize is the maximum request body size for quote creation (20MB).
// Quote requests may include a base64-encoded original file (image/PDF).
const maxQuoteBodySize = 20 << 20

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

	// User Management endpoints (Sprint 34)
	mux.Handle("GET /api/portal/v1/users", authMw(http.HandlerFunc(h.HandleListUsers)))
	mux.Handle("GET /api/portal/v1/invites", authMw(http.HandlerFunc(h.HandleListInvites)))
	mux.Handle("POST /api/portal/v1/invites", authMw(http.HandlerFunc(h.HandleInviteUser)))
	mux.Handle("PUT /api/portal/v1/users/{id}/role", authMw(http.HandlerFunc(h.HandleUpdateUserRole)))
	mux.Handle("PUT /api/portal/v1/users/{id}/status", authMw(http.HandlerFunc(h.HandleUpdateUserStatus)))

	// Quick Quote endpoints (AI-powered material list parsing)
	mux.Handle("POST /api/portal/v1/parsing/upload", authMw(http.HandlerFunc(h.HandleParseUpload)))
	mux.Handle("POST /api/portal/v1/parsing/text", authMw(http.HandlerFunc(h.HandleParseText)))
	mux.Handle("POST /api/portal/v1/quotes", authMw(http.HandlerFunc(h.HandleCreateQuickQuote)))
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

// --- User Management Handlers (Sprint 34) ---

func (h *Handler) HandleListUsers(w http.ResponseWriter, r *http.Request) {
	customerID := getPortalCustomerID(r)
	users, err := h.svc.ListCustomerUsers(r.Context(), customerID)
	if err != nil {
		portalWriteError(w, "Failed to load users", err, http.StatusInternalServerError)
		return
	}
	portalWriteJSON(w, users)
}

func (h *Handler) HandleListInvites(w http.ResponseWriter, r *http.Request) {
	customerID := getPortalCustomerID(r)
	invites, err := h.svc.ListPortalInvites(r.Context(), customerID)
	if err != nil {
		portalWriteError(w, "Failed to load invites", err, http.StatusInternalServerError)
		return
	}
	portalWriteJSON(w, invites)
}

func (h *Handler) HandleInviteUser(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)
	customerID := getPortalCustomerID(r)

	var req InviteUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		portalWriteError(w, "Invalid request body", err, http.StatusBadRequest)
		return
	}

	invite, err := h.svc.InviteUser(r.Context(), customerID, req)
	if err != nil {
		portalWriteError(w, "Failed to invite user", err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	portalWriteJSON(w, invite)
}

func (h *Handler) HandleUpdateUserRole(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)
	customerID := getPortalCustomerID(r)
	idStr := r.PathValue("id")
	userID, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var req UpdateUserRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		portalWriteError(w, "Invalid request body", err, http.StatusBadRequest)
		return
	}

	if err := h.svc.UpdateUserRole(r.Context(), customerID, userID, req.Role); err != nil {
		portalWriteError(w, "Failed to update role", err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) HandleUpdateUserStatus(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)
	customerID := getPortalCustomerID(r)
	idStr := r.PathValue("id")
	userID, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var req UpdateUserStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		portalWriteError(w, "Invalid request body", err, http.StatusBadRequest)
		return
	}

	if err := h.svc.UpdateUserStatus(r.Context(), customerID, userID, req.Status); err != nil {
		portalWriteError(w, "Failed to update status", err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// --- Quick Quote Handlers ---

// HandleParseUpload processes a material list file upload and returns parsed/matched items.
func (h *Handler) HandleParseUpload(w http.ResponseWriter, r *http.Request) {
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
		portalWriteError(w, "Failed to read uploaded file", err, http.StatusInternalServerError)
		return
	}

	contentType := http.DetectContentType(fileBytes)

	// Normalize content type for spreadsheets
	filename := header.Filename
	if contentType == "application/octet-stream" || contentType == "application/zip" {
		switch {
		case strings.HasSuffix(strings.ToLower(filename), ".xlsx") || strings.HasSuffix(strings.ToLower(filename), ".xls"):
			contentType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
		case strings.HasSuffix(strings.ToLower(filename), ".csv"):
			contentType = "text/csv"
		}
	}

	// For spreadsheets, convert to text
	if strings.Contains(contentType, "spreadsheet") || strings.HasSuffix(strings.ToLower(filename), ".xlsx") {
		textContent, convErr := convertSpreadsheetToText(fileBytes)
		if convErr != nil {
			portalWriteError(w, "Failed to process spreadsheet", convErr, http.StatusBadRequest)
			return
		}
		fileBytes = []byte(textContent)
		contentType = "text/plain"
	}

	resp, err := h.svc.ParseMaterialList(r.Context(), fileBytes, contentType)
	if err != nil {
		portalWriteError(w, "Failed to parse material list", err, http.StatusInternalServerError)
		return
	}

	portalWriteJSON(w, resp)
}

// HandleParseText parses a text material list and returns parsed/matched items.
func (h *Handler) HandleParseText(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)

	var req ParseTextRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		portalWriteError(w, "Invalid request body", err, http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(req.Text) == "" {
		http.Error(w, "Text is required", http.StatusBadRequest)
		return
	}

	resp, err := h.svc.ParseMaterialText(r.Context(), req.Text)
	if err != nil {
		portalWriteError(w, "Failed to parse text", err, http.StatusInternalServerError)
		return
	}

	portalWriteJSON(w, resp)
}

// HandleCreateQuickQuote creates a draft quote from parsed material list items.
func (h *Handler) HandleCreateQuickQuote(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxQuoteBodySize)
	customerID := getPortalCustomerID(r)

	var req PortalQuoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		portalWriteError(w, "Invalid request body", err, http.StatusBadRequest)
		return
	}

	if len(req.Items) == 0 {
		http.Error(w, "At least one item is required", http.StatusBadRequest)
		return
	}

	resp, err := h.svc.CreateQuickQuote(r.Context(), customerID, req)
	if err != nil {
		portalWriteError(w, "Failed to create quote", err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	portalWriteJSON(w, resp)
}

// convertSpreadsheetToText reads an xlsx file and converts it to plain text for AI extraction.
func convertSpreadsheetToText(data []byte) (string, error) {
	f, err := excelize.OpenReader(bytes.NewReader(data))
	if err != nil {
		return "", err
	}
	defer f.Close()

	var sb strings.Builder
	for _, sheet := range f.GetSheetList() {
		rows, err := f.GetRows(sheet)
		if err != nil {
			continue
		}
		for _, row := range rows {
			line := strings.Join(row, "\t")
			if strings.TrimSpace(line) != "" {
				sb.WriteString(line)
				sb.WriteString("\n")
			}
		}
	}
	return sb.String(), nil
}
