package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// 1. Locations
	// Check if exists first because unique constraint with NULL parent might be tricky
	var locationIDStr string
	err = db.QueryRow("SELECT id FROM locations WHERE code = 'MAIN'").Scan(&locationIDStr)
	var locationID uuid.UUID

	if err == sql.ErrNoRows {
		locationID = uuid.New()
		_, err = db.Exec(`
			INSERT INTO locations (id, code, type, description, path)
			VALUES ($1, 'MAIN', 'YARD', 'Main Yard', 'MAIN')
		`, locationID)
		if err != nil {
			log.Printf("Failed to insert location: %v", err)
		} else {
			fmt.Println("Inserted Location: Main Yard")
		}
	} else if err != nil {
		log.Fatalf("Failed to query location: %v", err)
	} else {
		locationID = uuid.MustParse(locationIDStr)
		fmt.Println("Found Location: Main Yard")
	}

	// 2. Products
	products := []struct {
		Description string
		SKU         string
		UOM         string
	}{
		{"2x4x8 SPF Premium", "LUM-248-PREM", "PCS"},
		{"2x6x10 SPF Premium", "LUM-2610-PREM", "PCS"},
		{"3/4 Plywood CDX", "PLY-34-CDX", "SF"},
		{"16d Common Nails (50lb)", "NAIL-16D-50", "BOX"},
	}

	productIDs := make([]string, 0)

	for _, p := range products {
		var id string
		err := db.QueryRow(`
			INSERT INTO products (sku, description, uom_primary)
			VALUES ($1, $2, $3)
			ON CONFLICT (sku) DO UPDATE SET description=$2
			RETURNING id
		`, p.SKU, p.Description, p.UOM).Scan(&id)
		if err != nil {
			log.Printf("Failed to insert product %s: %v", p.Description, err)
		} else {
			productIDs = append(productIDs, id)
			fmt.Printf("Upserted Product: %s\n", p.SKU)
			// Add inventory
			// Note: location text is kept for backward compatibility if not null, providing value just in case
			// Actually just ignore for MVP seed
			_, err := db.Exec(`
				INSERT INTO inventory (product_id, location_id, location, quantity)
				VALUES ($1, $2, 'MAIN', 1000)
			`, id, locationID)
			if err != nil {
				// Ignore if duplicate or constraint fail for now, or use Upsert if key exists
				// Inventory doesn't have unique constraint on product/location in 001/002 clearly defined other than PK id.
				// Assuming multiple entries allowed or handled.
				// Likely no constraint, but good practice
				log.Printf("Failed to insert inventory for %s: %v", p.SKU, err)
			}
		}
	}

	// 3. Customers
	customerID := uuid.New()
	err = db.QueryRow(`
		INSERT INTO customers (id, name, account_number, email, phone, address, credit_limit, balance_due)
		VALUES ($1, 'Acme Construction', 'ACME-001', 'billing@acme.com', '555-0101', '456 Builder Blvd', 50000, 1200)
		ON CONFLICT (account_number) DO UPDATE SET name='Acme Construction'
		RETURNING id
	`, customerID).Scan(&customerID)
	if err != nil {
		log.Printf("Failed to upsert customer: %v", err)
		// Try to fetch existing
		db.QueryRow("SELECT id FROM customers WHERE account_number='ACME-001'").Scan(&customerID)
	} else {
		fmt.Println("Upserted Customer: Acme Construction")
	}

	// 4. Fleet
	driverID := uuid.New()
	// Drivers might not have unique constraint on name, so check first
	err = db.QueryRow("SELECT id FROM drivers WHERE name='Colton (Dev)'").Scan(&driverID)
	if err == sql.ErrNoRows {
		driverID = uuid.New()
		_, err = db.Exec(`
			INSERT INTO drivers (id, name, status, license_number) 
			VALUES ($1, 'Colton (Dev)', 'ACTIVE', 'DL-12345')
		`, driverID)
		if err != nil {
			log.Printf("Failed to insert driver: %v", err)
		} else {
			fmt.Println("Inserted Driver: Colton (Dev)")
		}
	}

	vehicleID := uuid.New()
	// Vehicles Check
	err = db.QueryRow("SELECT id FROM vehicles WHERE name='Truck 1'").Scan(&vehicleID)
	if err == sql.ErrNoRows {
		vehicleID = uuid.New()
		_, err = db.Exec(`
			INSERT INTO vehicles (id, name, vehicle_type, license_plate)
			VALUES ($1, 'Truck 1', 'FLATBED', 'LBM-TRK-01')
		`, vehicleID)
		if err != nil {
			log.Printf("Failed to insert vehicle: %v", err)
		} else {
			fmt.Println("Inserted Vehicle: Truck 1")
		}
	}

	// 5. Orders & Deliveries
	// Always create a new order/route for "Today"
	orderID := uuid.New()
	_, err = db.Exec(`
		INSERT INTO orders (id, customer_id, total_amount, status)
		VALUES ($1, $2, 1500.00, 'CONFIRMED')
	`, orderID, customerID)
	if err != nil {
		log.Printf("Failed to insert order: %v", err)
	} else {
		fmt.Println("Created Order")
	}

	if len(productIDs) > 0 {
		_, err := db.Exec(`
			INSERT INTO order_lines (order_id, product_id, quantity, price_each)
			VALUES ($1, $2, 10, 5.50)
		`, orderID, productIDs[0])
		if err != nil {
			log.Printf("Failed to insert order line: %v", err)
		}
	}

	// Create Route
	routeID := uuid.New()
	_, err = db.Exec(`
		INSERT INTO delivery_routes (id, driver_id, vehicle_id, scheduled_date, status)
		VALUES ($1, $2, $3, CURRENT_DATE, 'SCHEDULED')
	`, routeID, driverID, vehicleID)
	if err != nil {
		log.Printf("Failed to insert route: %v", err)
	} else {
		fmt.Println("Created Route for Today")
	}

	// Assign Order to Route
	_, err = db.Exec(`
		INSERT INTO deliveries (id, route_id, order_id, status, stop_sequence, delivery_instructions)
		VALUES (uuid_generate_v4(), $1, $2, 'PENDING', 1, 'Leave at the back gate code 1234')
	`, routeID, orderID)
	if err != nil {
		log.Printf("Failed to insert delivery: %v", err)
	} else {
		fmt.Println("Created Delivery")
	}

	fmt.Println("Database seeded successfully!")
}
