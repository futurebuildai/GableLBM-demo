package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gablelbm/gable/internal/config"
	"github.com/gablelbm/gable/internal/customer"
	"github.com/gablelbm/gable/internal/inventory"
	"github.com/gablelbm/gable/internal/invoice"
	"github.com/gablelbm/gable/internal/location"
	"github.com/gablelbm/gable/internal/order"
	"github.com/gablelbm/gable/internal/product"
	"github.com/gablelbm/gable/internal/quote"
	"github.com/gablelbm/gable/pkg/database"
	"github.com/gablelbm/gable/pkg/middleware"
)

func main() {
	// 1. Setup Structured Logging (JSON)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// 2. Load Config
	cfg := config.Load()
	logger.Info("Starting server...", "port", cfg.Port)

	// 3. Database Connection
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		logger.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()
	logger.Info("Connected to database")

	// 4. Initialize Auth Middleware
	// If JWKS_URL is not set, we warn but allow startup (partial zero trust or dev mode)
	// For Strict Mode, we would exit.
	var authMw *middleware.AuthMiddleware
	if cfg.JWKSURL != "" {
		logger.Info("Initializing Auth Middleware", "jwks_url", cfg.JWKSURL)
		am, err := middleware.NewAuthMiddleware(context.Background(), middleware.AuthConfig{
			JWKSURL: cfg.JWKSURL,
			Issuer:  cfg.AuthIssuer,
		}, logger)
		if err != nil {
			logger.Error("Failed to initialize Auth Middleware", "error", err)
			os.Exit(1)
		}
		authMw = am
	} else {
		logger.Warn("NO JWKS_URL SET: AUTHENTICATION IS DISABLED (Use only for local dev)")
	}

	// 5. Setup Router & Modules
	mux := http.NewServeMux()

	// Initialize Modules
	productHandler := product.NewHandler(product.NewService(product.NewRepository(db)))
	productHandler.RegisterRoutes(mux)

	locationHandler := location.NewHandler(location.NewService(location.NewRepository(db)))
	locationHandler.RegisterRoutes(mux)

	// Inventory Service needs to be shared to Order Service
	inventoryRepo := inventory.NewRepository(db)
	inventorySvc := inventory.NewService(inventoryRepo)
	inventoryHandler := inventory.NewHandler(inventorySvc)
	inventoryHandler.RegisterRoutes(mux)

	customerHandler := customer.NewHandler(customer.NewService(customer.NewRepository(db)))
	customerHandler.RegisterRoutes(mux)

	quoteHandler := quote.NewHandler(quote.NewService(quote.NewRepository(db)))
	quoteHandler.RegisterRoutes(mux)

	// Invoice Module
	invoiceRepo := invoice.NewRepository(db)
	invoiceSvc := invoice.NewService(invoiceRepo)
	invoiceHandler := invoice.NewHandler(invoiceSvc)
	invoiceHandler.RegisterRoutes(mux)

	// Order Module - injected with InventoryService and InvoiceService
	orderHandler := order.NewHandler(order.NewService(order.NewRepository(db), inventorySvc, invoiceSvc))
	orderHandler.RegisterRoutes(mux)

	// Health Check (Public?)
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		status := "ok"
		dbStatus := "connected"
		if err := db.Pool.Ping(r.Context()); err != nil {
			status = "error"
			dbStatus = "disconnected"
			logger.Error("Health check failed", "error", err)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": status, "db": dbStatus})
	})

	// 6. Wrap Middleware
	var finalHandler http.Handler = mux
	if authMw != nil {
		// Wrap with Auth
		finalHandler = authMw.Handler(mux)
	}

	// Add Logger Middleware (Access Logs)
	finalHandler = RequestLogger(logger, finalHandler)

	// 7. Start Server with Graceful Shutdown
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Port),
		Handler: finalHandler,
	}

	// Run server in goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("Server failed", "error", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal using a buffered channel
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")

	// Create a deadline to wait for
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", "error", err)
		os.Exit(1)
	}

	logger.Info("Server exiting")
}

// RequestLogger logs incoming requests
func RequestLogger(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		// Wrap writer to capture status if needed (omitted for brevity, assume 200/error handled)
		next.ServeHTTP(w, r)
		logger.Info("Request processed",
			"method", r.Method,
			"path", r.URL.Path,
			"duration", time.Since(start).String(),
			"remote_addr", r.RemoteAddr,
		)
	})
}
