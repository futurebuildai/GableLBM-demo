package portal

import (
	"context"
	"encoding/base64"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gablelbm/gable/internal/customer"
	"github.com/gablelbm/gable/internal/inventory"
	"github.com/gablelbm/gable/internal/notification"
	"github.com/gablelbm/gable/internal/order"
	"github.com/gablelbm/gable/internal/parsing"
	"github.com/gablelbm/gable/internal/pricing"
	"github.com/gablelbm/gable/internal/product"
	"github.com/gablelbm/gable/internal/quote"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Service encapsulates portal business logic.
type Service struct {
	repo              *Repository
	jwtSecret         []byte
	logger            *slog.Logger
	pricingSvc        *pricing.Service
	customerSvc       *customer.Service
	inventorySvc      *inventory.Service
	orderSvc          *order.Service
	productSvc        *product.Service
	parsingSvc        *parsing.Service
	quoteSvc          *quote.Service
	emailSvc          notification.EmailService
	notificationEmail string
	appBaseURL        string
}

// NewService creates a new portal service.
func NewService(
	repo *Repository,
	logger *slog.Logger,
	pricingSvc *pricing.Service,
	customerSvc *customer.Service,
	inventorySvc *inventory.Service,
	orderSvc *order.Service,
	productSvc *product.Service,
	parsingSvc *parsing.Service,
	quoteSvc *quote.Service,
	emailSvc notification.EmailService,
	notificationEmail string,
	appBaseURL string,
) *Service {
	secret := os.Getenv("PORTAL_JWT_SECRET")
	if secret == "" {
		secret = "portal-dev-secret-change-in-production"
	}
	return &Service{
		repo:              repo,
		jwtSecret:         []byte(secret),
		logger:            logger,
		pricingSvc:        pricingSvc,
		customerSvc:       customerSvc,
		inventorySvc:      inventorySvc,
		orderSvc:          orderSvc,
		productSvc:        productSvc,
		parsingSvc:        parsingSvc,
		quoteSvc:          quoteSvc,
		emailSvc:          emailSvc,
		notificationEmail: notificationEmail,
		appBaseURL:        appBaseURL,
	}
}

// PortalClaims holds JWT claims for portal auth.
type PortalClaims struct {
	jwt.RegisteredClaims
	CustomerID     uuid.UUID `json:"customer_id"`
	CustomerUserID uuid.UUID `json:"customer_user_id"`
	Email          string    `json:"email"`
	Name           string    `json:"name"`
}

// Login authenticates a customer user and returns a JWT.
func (s *Service) Login(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	user, err := s.repo.GetCustomerUserByEmail(ctx, req.Email)
	if err != nil {
		s.logger.Warn("Portal login: user not found", "email", req.Email)
		return nil, fmt.Errorf("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		s.logger.Warn("Portal login: invalid password", "email", req.Email)
		return nil, fmt.Errorf("invalid credentials")
	}

	// Generate JWT
	now := time.Now()
	claims := PortalClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.ID.String(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(24 * time.Hour)),
			Issuer:    "gable-portal",
		},
		CustomerID:     user.CustomerID,
		CustomerUserID: user.ID,
		Email:          user.Email,
		Name:           user.Name,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(s.jwtSecret)
	if err != nil {
		s.logger.Error("Portal login: failed to sign JWT", "error", err)
		return nil, fmt.Errorf("authentication failed")
	}

	// Get portal config for branding
	cfg, err := s.repo.GetPortalConfig(ctx)
	if err != nil {
		s.logger.Error("Portal login: failed to get config", "error", err)
		return nil, fmt.Errorf("authentication failed")
	}

	s.logger.Info("Portal login success", "email", req.Email, "customer_id", user.CustomerID)

	return &LoginResponse{
		Token:  tokenStr,
		User:   *user,
		Config: *cfg,
	}, nil
}

// GetConfig returns portal branding config (public).
func (s *Service) GetConfig(ctx context.Context) (*PortalConfig, error) {
	return s.repo.GetPortalConfig(ctx)
}

// GetDashboard returns the contractor dashboard data.
func (s *Service) GetDashboard(ctx context.Context, customerID uuid.UUID) (*PortalDashboardDTO, error) {
	balance, creditLimit, pastDue, err := s.repo.GetCustomerARSummary(ctx, customerID)
	if err != nil {
		return nil, fmt.Errorf("failed to load dashboard: %w", err)
	}

	orders, err := s.repo.ListOrdersByCustomer(ctx, customerID)
	if err != nil {
		return nil, fmt.Errorf("failed to load recent orders: %w", err)
	}

	// Limit to 5 most recent for dashboard
	recentOrders := orders
	if len(recentOrders) > 5 {
		recentOrders = recentOrders[:5]
	}

	return &PortalDashboardDTO{
		BalanceDue:   balance,
		CreditLimit:  creditLimit,
		PastDue:      pastDue,
		RecentOrders: recentOrders,
	}, nil
}

// ListOrders returns all orders for a customer.
func (s *Service) ListOrders(ctx context.Context, customerID uuid.UUID) ([]PortalOrderDTO, error) {
	return s.repo.ListOrdersByCustomer(ctx, customerID)
}

// GetOrder returns a single order scoped to a customer.
func (s *Service) GetOrder(ctx context.Context, orderID, customerID uuid.UUID) (*PortalOrderDTO, error) {
	return s.repo.GetOrderByIDAndCustomer(ctx, orderID, customerID)
}

// ListInvoices returns all invoices for a customer.
func (s *Service) ListInvoices(ctx context.Context, customerID uuid.UUID) ([]PortalInvoiceDTO, error) {
	return s.repo.ListInvoicesByCustomer(ctx, customerID)
}

// GetInvoice returns a single invoice scoped to a customer.
func (s *Service) GetInvoice(ctx context.Context, invoiceID, customerID uuid.UUID) (*PortalInvoiceDTO, error) {
	return s.repo.GetInvoiceByIDAndCustomer(ctx, invoiceID, customerID)
}

// ListDeliveries returns all deliveries for a customer.
func (s *Service) ListDeliveries(ctx context.Context, customerID uuid.UUID) ([]PortalDeliveryDTO, error) {
	return s.repo.ListDeliveriesByCustomer(ctx, customerID)
}

// GetDelivery returns a single delivery scoped to a customer.
func (s *Service) GetDelivery(ctx context.Context, deliveryID, customerID uuid.UUID) (*PortalDeliveryDTO, error) {
	return s.repo.GetDeliveryByIDAndCustomer(ctx, deliveryID, customerID)
}

// CreateReorder duplicates a historical order as a new DRAFT.
func (s *Service) CreateReorder(ctx context.Context, customerID uuid.UUID, req ReorderRequest) (*ReorderResponse, error) {
	newOrderID, err := s.repo.CreateReorder(ctx, customerID, req.OrderID)
	if err != nil {
		return nil, fmt.Errorf("failed to create reorder: %w", err)
	}

	s.logger.Info("Portal reorder created", "customer_id", customerID, "source_order", req.OrderID, "new_order", newOrderID)

	return &ReorderResponse{
		OrderID: newOrderID,
		Message: "Order draft created successfully",
	}, nil
}

// ParseMaterialList parses an uploaded material list file through AI extraction + product matching.
// For demo purposes, auto-seeds missing standard SKUs into the product catalog so the matching demo works.
func (s *Service) ParseMaterialList(ctx context.Context, fileBytes []byte, contentType string) (*parsing.ParseResponse, error) {
	start := time.Now()

	extracted, err := s.parsingSvc.ExtractItemsWithAI(ctx, fileBytes, contentType)
	if err != nil {
		return nil, fmt.Errorf("failed to extract items: %w", err)
	}

	items, err := s.parsingSvc.MatchProducts(ctx, extracted)
	if err != nil {
		return nil, fmt.Errorf("failed to match products: %w", err)
	}

	// Demo mode: auto-seed missing standard SKUs and re-match
	needsRematch := false
	for _, item := range items {
		if item.IsSpecialOrder && item.RawText != "" {
			created := s.autoSeedProduct(ctx, item.RawText, item.UOM)
			if created {
				needsRematch = true
			}
		}
	}

	if needsRematch {
		items, err = s.parsingSvc.MatchProducts(ctx, extracted)
		if err != nil {
			return nil, fmt.Errorf("failed to re-match products: %w", err)
		}
	}

	return &parsing.ParseResponse{
		Items:       items,
		ParseTimeMs: time.Since(start).Milliseconds(),
		ItemCount:   len(items),
	}, nil
}

// autoSeedProduct attempts to auto-create a product from raw parsed text.
// Returns true if a product was created. Used in demo mode to ensure SKU matching works.
// If the product type can't be identified, it creates a Special Order SKU as a fallback.
func (s *Service) autoSeedProduct(ctx context.Context, rawText string, uom string) bool {
	// Generate a reasonable SKU and description from the raw text
	sku, description, basePrice, unitUOM := inferProductFromText(rawText, uom)

	// Fallback: create a Special Order SKU for unrecognized materials
	if sku == "" {
		sku, description, basePrice, unitUOM = buildSpecialOrderSKU(rawText, uom)
	}

	p := &product.Product{
		ID:          uuid.New(),
		SKU:         sku,
		Description: description,
		UOMPrimary:  product.UOM(unitUOM),
		BasePrice:   basePrice,
		WeightLbs:   5.0,
		ReorderPoint: 0,
		ReorderQty:   0,
		AverageUnitCost: basePrice * 0.76,
		TargetMargin:    0.24,
	}

	if err := s.productSvc.CreateProduct(ctx, p); err != nil {
		s.logger.Debug("Auto-seed product skipped (may already exist)", "sku", sku, "error", err)
		return false
	}

	s.logger.Info("Auto-seeded product for demo", "sku", sku, "description", description, "price", basePrice)
	return true
}

// buildSpecialOrderSKU creates a Special Order SKU for materials that can't be identified.
// These appear in the ERP as special-order items requiring manual pricing/sourcing.
func buildSpecialOrderSKU(rawText string, uom string) (sku, description string, price float64, unitUOM string) {
	// Clean up the raw text for the description
	cleaned := strings.TrimSpace(rawText)
	// Remove leading quantity patterns (e.g. "12 -", "50 pcs -")
	qtyPat := regexp.MustCompile(`(?i)^\s*\d+\s*(?:pcs?|ea|pieces?|each|lf|sf|sheets?|bags?|rolls?|bundles?|libras?)?\s*[-–—]?\s*`)
	cleaned = qtyPat.ReplaceAllString(cleaned, "")
	cleaned = strings.TrimSpace(cleaned)
	if cleaned == "" {
		cleaned = rawText
	}

	// Generate a short hash from raw text for unique SKU
	hash := fmt.Sprintf("%04X", simpleHash(strings.ToLower(cleaned))%0xFFFF)

	// Truncate description to something reasonable
	desc := cleaned
	if len(desc) > 60 {
		desc = desc[:60]
	}

	unitUOM = "PCS"
	if uom != "" {
		unitUOM = uom
	}

	sku = fmt.Sprintf("SPO-%s", hash)
	description = fmt.Sprintf("Special Order: %s", desc)
	price = 0.00 // Price TBD — salesperson sets during quote review
	return
}

// simpleHash produces a simple numeric hash from a string (for SKU generation).
func simpleHash(s string) uint32 {
	var h uint32 = 2166136261
	for i := 0; i < len(s); i++ {
		h ^= uint32(s[i])
		h *= 16777619
	}
	return h
}

// inferProductFromText generates a SKU, description, price, and UOM from raw material list text.
func inferProductFromText(rawText string, uom string) (sku, description string, price float64, unitUOM string) {
	text := strings.ToLower(strings.TrimSpace(rawText))
	unitUOM = "PCS"
	if uom != "" {
		unitUOM = uom
	}

	// Dimensional lumber patterns: "2x4x8", "2x6x12", "2x10x16"
	dimPattern := regexp.MustCompile(`(\d+)\s*x\s*(\d+)(?:\s*x\s*(\d+))?`)
	if matches := dimPattern.FindStringSubmatch(text); len(matches) > 0 {
		dim := strings.ReplaceAll(matches[0], " ", "")
		species := "SPF"
		if strings.Contains(text, "doug") || strings.Contains(text, "df") {
			species = "DF"
		} else if strings.Contains(text, "hem") {
			species = "HF"
		} else if strings.Contains(text, "treat") || strings.Contains(text, "pt") || strings.Contains(text, "pressure") {
			species = "PT"
		} else if strings.Contains(text, "cedar") {
			species = "WRC"
		}

		grade := "#2"
		if strings.Contains(text, "stud") {
			grade = "STUD"
		} else if strings.Contains(text, "#1") || strings.Contains(text, "select") {
			grade = "SEL"
		}

		sku = fmt.Sprintf("LBR-%s-%s-%s", strings.ToUpper(dim), species, grade)
		description = fmt.Sprintf("%s %s %s Lumber", strings.ToUpper(dim), species, grade)

		// Price based on dimensions
		w := 2.0
		h := 4.0
		l := 8.0
		if len(matches) > 1 { if v, e := strconv.ParseFloat(matches[1], 64); e == nil { w = v } }
		if len(matches) > 2 { if v, e := strconv.ParseFloat(matches[2], 64); e == nil { h = v } }
		if len(matches) > 3 && matches[3] != "" { if v, e := strconv.ParseFloat(matches[3], 64); e == nil { l = v } }
		price = (w * h * l) / 100.0 * 3.50
		if price < 3.0 { price = 3.50 }
		return
	}

	// Sheet goods
	sheetPattern := regexp.MustCompile(`(?i)(osb|cdx|plywood|drywall|sheathing|hardiboard|mdf)`)
	if matches := sheetPattern.FindStringSubmatch(text); len(matches) > 0 {
		material := strings.ToUpper(matches[1])
		thickness := "1/2"
		thickPattern := regexp.MustCompile(`(\d+/\d+)`)
		if tm := thickPattern.FindStringSubmatch(text); len(tm) > 0 {
			thickness = tm[1]
		}
		sku = fmt.Sprintf("SHT-%s-%s-4X8", material, strings.ReplaceAll(thickness, "/", ""))
		description = fmt.Sprintf("%s %s\" 4x8 Sheet", material, thickness)
		switch material {
		case "OSB": price = 28.50
		case "PLYWOOD", "CDX": price = 42.00
		case "DRYWALL": price = 12.50
		case "MDF": price = 35.00
		default: price = 30.00
		}
		unitUOM = "PCS"
		return
	}

	// Concrete/masonry
	if strings.Contains(text, "quikrete") || strings.Contains(text, "concrete") || strings.Contains(text, "cement") {
		weight := "80"
		if strings.Contains(text, "60") { weight = "60" }
		if strings.Contains(text, "50") { weight = "50" }
		sku = fmt.Sprintf("CON-QKRT-%sLB", weight)
		description = fmt.Sprintf("Quikrete Concrete Mix %slb", weight)
		price = 6.50
		unitUOM = "BAG"
		return
	}

	// House wrap
	if strings.Contains(text, "tyvek") || strings.Contains(text, "house wrap") || strings.Contains(text, "housewrap") {
		sku = "WRP-TYVEK-9X150"
		description = "Tyvek HomeWrap 9ft x 150ft Roll"
		price = 165.00
		unitUOM = "RL"
		return
	}

	// Simpson Strong-Tie
	if strings.Contains(text, "simpson") || strings.Contains(text, "strong-tie") || strings.Contains(text, "strong tie") {
		connector := "A35"
		if strings.Contains(text, "h10") { connector = "H10" }
		if strings.Contains(text, "lus") { connector = "LUS26" }
		if strings.Contains(text, "a34") { connector = "A34" }
		sku = fmt.Sprintf("HWR-SST-%s", connector)
		description = fmt.Sprintf("Simpson Strong-Tie %s Connector", connector)
		price = 2.85
		unitUOM = "EA"
		return
	}

	// Nails
	if strings.Contains(text, "nail") {
		sku = "HWR-NAIL-16D-50LB"
		description = "16d Common Nails 50lb Box"
		price = 65.00
		unitUOM = "BOX"
		return
	}

	// Screws
	if strings.Contains(text, "screw") {
		sku = "HWR-SCREW-DECK-5LB"
		description = "Deck Screws #8 x 3\" 5lb Box"
		price = 28.50
		unitUOM = "BOX"
		return
	}

	// Roofing
	if strings.Contains(text, "shingle") {
		sku = "RFG-ARCH-WGRY"
		description = "Architectural Shingles Weathered Gray"
		price = 34.50
		unitUOM = "BUNDLE"
		return
	}

	// Insulation
	if strings.Contains(text, "insulation") || strings.Contains(text, "r-") || strings.Contains(text, "fiberglass") {
		rValue := "19"
		if strings.Contains(text, "13") { rValue = "13" }
		if strings.Contains(text, "30") { rValue = "30" }
		if strings.Contains(text, "38") { rValue = "38" }
		sku = fmt.Sprintf("INS-FG-R%s", rValue)
		description = fmt.Sprintf("Fiberglass Batt Insulation R-%s", rValue)
		price = 45.00
		unitUOM = "BUNDLE"
		return
	}

	// Generic fallback - don't auto-seed if we can't identify the product type
	return "", "", 0, ""
}

// ParseMaterialText parses plain text material list through AI extraction + product matching.
func (s *Service) ParseMaterialText(ctx context.Context, text string) (*parsing.ParseResponse, error) {
	return s.ParseMaterialList(ctx, []byte(text), "text/plain")
}

// CreateQuickQuote creates a draft quote from portal-submitted line items and sends email notification.
func (s *Service) CreateQuickQuote(ctx context.Context, customerID uuid.UUID, req PortalQuoteRequest) (*PortalQuoteResponse, error) {
	// Look up customer name for the quote
	cust, err := s.customerSvc.GetCustomer(ctx, customerID)
	customerName := "Portal Customer"
	if err == nil && cust != nil {
		customerName = cust.Name
	}

	// Build quote lines
	var lines []quote.QuoteLine
	for _, item := range req.Items {
		lines = append(lines, quote.QuoteLine{
			ID:          uuid.New(),
			ProductID:   item.ProductID,
			SKU:         item.SKU,
			Description: item.Description,
			Quantity:    item.Quantity,
			UOM:         "PCS",
			UnitPrice:   item.UnitPrice,
		})
	}

	deliveryType := "PICKUP"
	if req.DeliveryMethod == "DELIVERY" {
		deliveryType = "DELIVERY"
	}

	q := &quote.Quote{
		ID:              uuid.New(),
		CustomerID:      customerID,
		CustomerName:    customerName,
		State:           quote.QuoteStateDraft,
		Source:          "ai",
		DeliveryType:    deliveryType,
		DeliveryAddress: req.DeliveryAddress,
		Lines:           lines,
	}

	// Attach original input for AI traceability
	if req.OriginalText != "" {
		filename := req.OriginalFilename
		if filename == "" {
			filename = "material-list.txt"
		}
		q.OriginalFilename = filename

		// Detect if this is a base64-encoded file upload vs plain text
		lowerName := strings.ToLower(filename)
		isFileUpload := strings.HasSuffix(lowerName, ".png") ||
			strings.HasSuffix(lowerName, ".jpg") ||
			strings.HasSuffix(lowerName, ".jpeg") ||
			strings.HasSuffix(lowerName, ".pdf") ||
			strings.HasSuffix(lowerName, ".xlsx") ||
			strings.HasSuffix(lowerName, ".xls")

		if isFileUpload {
			decoded, err := base64.StdEncoding.DecodeString(req.OriginalText)
			if err == nil {
				q.OriginalFile = decoded
				q.OriginalContentType = http.DetectContentType(decoded)
			} else {
				q.OriginalFile = []byte(req.OriginalText)
				q.OriginalContentType = "text/plain"
			}
		} else {
			q.OriginalFile = []byte(req.OriginalText)
			q.OriginalContentType = "text/plain"
		}
	}

	// Attach AI parse mapping data
	if len(req.ParseMap) > 0 {
		q.ParseMap = req.ParseMap
	}

	if err := s.quoteSvc.CreateQuote(ctx, q); err != nil {
		return nil, fmt.Errorf("failed to create quote: %w", err)
	}

	s.logger.Info("Portal quick quote created",
		"quote_id", q.ID,
		"customer_id", customerID,
		"line_count", len(lines),
		"total", q.TotalAmount,
	)

	// Send email notification (best-effort, don't fail the request)
	if s.emailSvc != nil && s.notificationEmail != "" {
		quoteURL := fmt.Sprintf("%s/erp/quotes/%s", s.appBaseURL, q.ID.String())
		if err := s.emailSvc.SendQuoteNotification(ctx, s.notificationEmail, q.ID.String(), customerName, q.TotalAmount, quoteURL); err != nil {
			s.logger.Error("Failed to send quote notification email", "error", err, "quote_id", q.ID)
		}
	}

	return &PortalQuoteResponse{
		QuoteID: q.ID,
		Message: "Quote draft created successfully",
	}, nil
}

// ParseToken parses and validates a portal JWT. Used by middleware.
func (s *Service) ParseToken(tokenStr string) (*PortalClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &PortalClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return s.jwtSecret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*PortalClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}
