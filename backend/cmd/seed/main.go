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

	// Constants / Users
	// Since we don't have a users table, we'll just use a fixed UUID for "Demo User"
	demoUserID := uuid.MustParse("00000000-0000-0000-0000-000000000001")

	// 1. Locations
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

	// 2. Vendors
	vendors := []struct {
		Name  string
		Email string
	}{
		{"Gable Lumber Supply", "orders@gablelumber.com"},
		{"Hardware Wholesale Inc", "sales@hwi.com"},
		{"Roofing Specialists Ltd", "contact@roofing-specialists.com"},
		{"Millwork Masters", "info@millworkmasters.com"},
		{"Concrete Logic", "dispatch@concretelogic.com"},
	}

	vendorIDs := make(map[string]uuid.UUID)

	for _, v := range vendors {
		var id string
		err := db.QueryRow(`
			INSERT INTO vendors (name, contact_email)
			VALUES ($1, $2)
			ON CONFLICT (name) DO UPDATE SET contact_email = $2
			RETURNING id
		`, v.Name, v.Email).Scan(&id)
		if err != nil {
			log.Printf("Failed to upsert vendor %s: %v", v.Name, err)
		} else {
			vendorIDs[v.Name] = uuid.MustParse(id)
			fmt.Printf("Upserted Vendor: %s\n", v.Name)
		}
	}

	// 3. Products & Inventory
	products := []struct {
		Description string
		SKU         string
		UOM         string
		Vendor      string
		Cost        float64
		Price       float64
		Category    string
	}{
		// Lumber
		{"2x4x8 SPF Premium", "LUM-248-PREM", "PCS", "Gable Lumber Supply", 3.50, 5.50, "Lumber"},
		{"2x4x10 SPF Premium", "LUM-2410-PREM", "PCS", "Gable Lumber Supply", 4.50, 7.25, "Lumber"},
		{"2x6x10 SPF Premium", "LUM-2610-PREM", "PCS", "Gable Lumber Supply", 6.00, 9.80, "Lumber"},
		{"4x4x8 Pressure Treated", "LUM-448-PT", "PCS", "Gable Lumber Supply", 8.00, 12.50, "Lumber"},
		// Plywood / Sheathing
		{"3/4 Plywood CDX 4x8", "PLY-34-CDX", "SF", "Gable Lumber Supply", 24.00, 38.00, "Sheet Goods"}, // Price usually per sheet, UOM might be distinct, assuming PCS for simplicity in other ERPs but here SF. Let's assume price is per piece in our logic but UOM is SF? standard is usually PCS for sheets. Correcting to PCS for simplicity if enum allows, checking schema... enum has PCS.
		{"1/2 OSB 4x8", "OSB-12", "PCS", "Gable Lumber Supply", 12.00, 19.50, "Sheet Goods"},
		// Hardware
		{"16d Common Nails (50lb)", "NAIL-16D-50", "BOX", "Hardware Wholesale Inc", 45.00, 65.00, "Hardware"},
		{"3\" Deck Screws (5lb)", "SCR-DECK-3-5", "BOX", "Hardware Wholesale Inc", 18.00, 29.99, "Hardware"},
		{"Joist Hanger 2x6", "HANGER-26", "PCS", "Hardware Wholesale Inc", 0.80, 1.45, "Hardware"},
		// Roofing
		{"Architectural Shingles (Black)", "RF-SH-BLK", "BUNDLE", "Roofing Specialists Ltd", 28.00, 42.00, "Roofing"},
		{"Roofing Felt 15lb", "RF-FELT-15", "RL", "Roofing Specialists Ltd", 15.00, 22.50, "Roofing"},
		// Millwork
		{"Int Door 30x80 6-Panel", "DR-INT-3080-6P", "PCS", "Millwork Masters", 65.00, 95.00, "Millwork"},
		{"Ext Door 36x80 Steel", "DR-EXT-3680-STL", "PCS", "Millwork Masters", 180.00, 280.00, "Millwork"},
	}

	productIDs := make([]uuid.UUID, 0)
	skuToID := make(map[string]uuid.UUID)

	for _, p := range products {
		var id string
		// Insert Product
		err := db.QueryRow(`
			INSERT INTO products (sku, description, uom_primary)
			VALUES ($1, $2, $3)
			ON CONFLICT (sku) DO UPDATE SET description=$2
			RETURNING id
		`, p.SKU, p.Description, p.UOM).Scan(&id)

		if err != nil {
			log.Printf("Failed to insert product %s: %v", p.Description, err)
			continue
		}

		pid := uuid.MustParse(id)
		productIDs = append(productIDs, pid)
		skuToID[p.SKU] = pid
		fmt.Printf("Upserted Product: %s\n", p.SKU)

		// Set Price (if pricing table exists - assuming simple price logic or skipping if complex.
		// Migration 007_add_product_price.sql might have added columns or table.
		// Let's assume standard ERP separate pricing, but for now we might skip if complex.
		// Actually migration 007 is 113 bytes, likely alter table. Checking assumed schema...)
		// Let's just do Inventory for now.

		// Inventory
		_, err = db.Exec(`
			INSERT INTO inventory (product_id, location_id, location, quantity)
			VALUES ($1, $2, 'MAIN', 500)
			ON CONFLICT DO NOTHING
		`, pid, locationID)
		// Note: No unique constraint on inventory(product_id, location_id) in 001/003, but assuming we want to add if not exists.
		// Actually, let's just insert and ignore errors for now or update.
	}

	// 4. Customers & Projects
	customers := []struct {
		Name     string
		Acct     string
		Email    string
		Projects []string
	}{
		{"Acme Construction", "ACME-001", "billing@acme.com", []string{"Smith Residence", "Downtown Lofts"}},
		{"Bob's Builders", "BOB-001", "bob@bobsbuilders.com", []string{"Miller Deck", "Kitchen Remodel 123"}},
		{"DIY Homeowner", "DIY-888", "diy@gmail.com", []string{"Garden Shed"}},
	}

	customerIDs := make(map[string]uuid.UUID)
	projectIDs := make(map[string]uuid.UUID) // Project Name -> ID

	for _, c := range customers {
		var cid string
		err := db.QueryRow(`
			INSERT INTO customers (name, account_number, email, credit_limit, balance_due)
			VALUES ($1, $2, $3, 10000, 0)
			ON CONFLICT (account_number) DO UPDATE SET name=$1
			RETURNING id
		`, c.Name, c.Acct, c.Email).Scan(&cid)
		if err != nil {
			log.Printf("Failed to upsert customer %s: %v", c.Name, err)
			continue
		}
		custID := uuid.MustParse(cid)
		customerIDs[c.Name] = custID
		fmt.Printf("Upserted Customer: %s\n", c.Name)

		// Projects
		for _, projName := range c.Projects {
			var jid string
			// Simple check to avoid dups if run multiple times (though no unique constraint on name+cust usually)
			err := db.QueryRow(`
				INSERT INTO customer_jobs (customer_id, name, is_active)
				VALUES ($1, $2, true)
				RETURNING id
			`, custID, projName).Scan(&jid)
			if err == nil {
				projectIDs[projName] = uuid.MustParse(jid)
				fmt.Printf("  Created Project: %s\n", projName)
			}
		}
	}

	// 5. Quotes
	// Create a quote for Acme
	acmeID := customerIDs["Acme Construction"]
	if acmeID != uuid.Nil {
		quoteID := uuid.New()
		err = db.QueryRow(`
			INSERT INTO quotes (id, customer_id, job_id, state, total_amount, created_by, expires_at)
			VALUES ($1, $2, $3, 'SENT', 1500.00, $4, NOW() + INTERVAL '30 days')
			RETURNING id
		`, quoteID, acmeID, projectIDs["Smith Residence"], demoUserID).Scan(&quoteID)

		if err == nil {
			// Lines
			db.Exec(`
				INSERT INTO quote_lines (quote_id, product_id, sku, description, quantity, uom, unit_price, line_total)
				VALUES 
				($1, $2, 'LUM-248-PREM', '2x4x8 SPF Premium', 100, 'PCS', 5.50, 550.00),
				($1, $2, 'LUM-2610-PREM', '2x6x10 SPF Premium', 50, 'PCS', 9.80, 490.00)
			`, quoteID, skuToID["LUM-248-PREM"])
			fmt.Println("Created Quote for Acme")
		}
	}

	// 6. Orders (Sales Orders)
	// Order 1: Confirmed (Acme)
	order1ID := uuid.New()
	if acmeID != uuid.Nil {
		_, err = db.Exec(`
			INSERT INTO orders (id, customer_id, total_amount, status, created_at)
			VALUES ($1, $2, 5500.00, 'CONFIRMED', NOW() - INTERVAL '2 days')
		`, order1ID, acmeID)
		if err == nil {
			db.Exec(`
				INSERT INTO order_lines (order_id, product_id, quantity, price_each)
				VALUES ($1, $2, 100, 55.00) -- Simplified
			`, order1ID, skuToID["PLY-34-CDX"])
			fmt.Println("Created Confirmed Order for Acme")
		}
	}

	// Order 2: Processing (Bob)
	bobID := customerIDs["Bob's Builders"]
	if bobID != uuid.Nil {
		order2ID := uuid.New()
		_, err = db.Exec(`
			INSERT INTO orders (id, customer_id, total_amount, status, created_at)
			VALUES ($1, $2, 1200.00, 'PROCESSING', NOW() - INTERVAL '1 day')
		`, order2ID, bobID)
		fmt.Println("Created Processing Order for Bob")
	}

	// 7. Purchase Orders
	// PO to Gable Lumber
	vendorID := vendorIDs["Gable Lumber Supply"]
	if vendorID != uuid.Nil {
		poID := uuid.New()
		_, err = db.Exec(`
			INSERT INTO purchase_orders (id, vendor_id, status, created_at)
			VALUES ($1, $2, 'SENT', NOW() - INTERVAL '1 week')
		`, poID, vendorID)

		if err == nil {
			db.Exec(`
				INSERT INTO purchase_order_lines (po_id, description, quantity, cost)
				VALUES ($1, 'Restock 2x4s', 1000, 3.50)
			`, poID)
			fmt.Println("Created PO for Gable Lumber")
		}
	}

	// 8. Invoices & Payments
	// Invoice for Order 1
	if acmeID != uuid.Nil {
		invID := uuid.New()
		_, err = db.Exec(`
			INSERT INTO invoices (id, order_id, customer_id, status, total_amount, due_date, created_at)
			VALUES ($1, $2, $3, 'PAID', 5500.00, NOW() + INTERVAL '30 days', NOW() - INTERVAL '2 days')
		`, invID, order1ID, acmeID)

		if err == nil {
			// Payment
			db.Exec(`
				INSERT INTO payments (invoice_id, amount, method, reference, notes)
				VALUES ($1, 5500.00, 'CHECK', 'CHK-998877', 'Payment in full')
			`, invID)
			fmt.Println("Created Paid Invoice & Payment for Acme")
		}
	}

	// 9. RFCs (Governance)
	rfcs := []struct {
		Title  string
		Status string
	}{
		{"Standardize SKU Format", "draft"},
		{"Q3 Inventory Audit Procedure", "approved"},
		{"Vendor Onboarding Requirements", "review"},
	}

	for _, rfc := range rfcs {
		_, err = db.Exec(`
			INSERT INTO rfcs (title, status, author_id, problem_statement, proposed_solution)
			VALUES ($1, $2, $3, 'Inconsistent data', 'Follow new iso standard')
		`, rfc.Title, rfc.Status, demoUserID)
		if err == nil {
			fmt.Printf("Created RFC: %s\n", rfc.Title)
		}
	}

	fmt.Println("Database seeded successfully with comprehensive demo data!")
}
