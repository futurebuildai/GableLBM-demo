package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gablelbm/gable/internal/config"
	"github.com/gablelbm/gable/internal/inventory"
	"github.com/gablelbm/gable/internal/location"
	"github.com/gablelbm/gable/internal/product"
	"github.com/gablelbm/gable/pkg/database"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()
	log.Println("Connected to database")

	mux := http.NewServeMux()

	// Initialize Product Module
	productRepo := product.NewRepository(db)
	productService := product.NewService(productRepo)
	productHandler := product.NewHandler(productService)
	productHandler.RegisterRoutes(mux)

	// Initialize Location Module
	locationRepo := location.NewRepository(db)
	locationService := location.NewService(locationRepo)
	locationHandler := location.NewHandler(locationService)
	locationHandler.RegisterRoutes(mux)

	// Initialize Inventory Module
	inventoryRepo := inventory.NewRepository(db)
	inventoryService := inventory.NewService(inventoryRepo)
	inventoryHandler := inventory.NewHandler(inventoryService)
	inventoryHandler.RegisterRoutes(mux)

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		status := "ok"
		dbStatus := "connected"

		if err := db.Pool.Ping(r.Context()); err != nil {
			status = "error"
			dbStatus = "disconnected"
			log.Printf("Health check failed: %v", err)
		}

		response := map[string]string{
			"status": status,
			"db":     dbStatus,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Server starting on %s...", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
