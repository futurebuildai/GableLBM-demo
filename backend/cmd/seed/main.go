package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

// Helper to generate random dates
func randomDate() time.Time {
	min := time.Date(2023, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	max := time.Now().Unix()
	delta := max - min
	sec := rand.Int63n(delta) + min
	return time.Unix(sec, 0)
}

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		// Fallback for local dev if not set, though ideally passed in
		dbURL = "postgresql://postgres:postgres@localhost:5432/gable_erp_db?sslmode=disable"
		fmt.Println("DATABASE_URL not set, using default: " + dbURL)
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
	demoUserID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	log.Println("Seeding started...")

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
		}
	} else if err != nil {
		log.Fatalf("Failed to query location: %v", err)
	} else {
		locationID = uuid.MustParse(locationIDStr)
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
		{"Fastener Depot", "orders@fastenerdepot.com"},
		{"Valley Insulation", "sales@valleyinsulation.com"},
		{"Apex Windows", "orders@apexwindows.com"},
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
		if err == nil {
			vendorIDs[v.Name] = uuid.MustParse(id)
		}
	}

	// 3. Products & Inventory
	// Expanded Product List
	products := []struct {
		Description string
		SKU         string
		UOM         string
		Vendor      string
		Cost        float64
		Price       float64
		Category    string
	}{
		// Dimensional Lumber
		{"2x4x8 SPF Premium", "LUM-248-PREM", "PCS", "Gable Lumber Supply", 3.50, 5.50, "Lumber"},
		{"2x4x10 SPF Premium", "LUM-2410-PREM", "PCS", "Gable Lumber Supply", 4.50, 7.25, "Lumber"},
		{"2x4x12 SPF Premium", "LUM-2412-PREM", "PCS", "Gable Lumber Supply", 5.40, 8.75, "Lumber"},
		{"2x4x92-5/8 SPF Stud", "LUM-2492-STUD", "PCS", "Gable Lumber Supply", 3.20, 5.10, "Lumber"},
		{"2x6x10 SPF Premium", "LUM-2610-PREM", "PCS", "Gable Lumber Supply", 6.00, 9.80, "Lumber"},
		{"2x6x12 SPF Premium", "LUM-2612-PREM", "PCS", "Gable Lumber Supply", 7.20, 11.75, "Lumber"},
		{"2x6x16 SPF Premium", "LUM-2616-PREM", "PCS", "Gable Lumber Supply", 9.60, 15.60, "Lumber"},
		{"2x8x10 SPF No.2", "LUM-2810-NO2", "PCS", "Gable Lumber Supply", 8.50, 13.90, "Lumber"},
		{"2x8x16 SPF No.2", "LUM-2816-NO2", "PCS", "Gable Lumber Supply", 13.60, 22.25, "Lumber"},
		{"2x10x12 SPF No.2", "LUM-21012-NO2", "PCS", "Gable Lumber Supply", 15.00, 24.50, "Lumber"},
		{"2x12x16 SPF No.2", "LUM-21216-NO2", "PCS", "Gable Lumber Supply", 28.00, 45.00, "Lumber"},
		{"4x4x8 Pressure Treated", "LUM-448-PT", "PCS", "Gable Lumber Supply", 8.00, 12.50, "Lumber"},
		{"4x4x10 Pressure Treated", "LUM-4410-PT", "PCS", "Gable Lumber Supply", 10.00, 15.75, "Lumber"},
		{"6x6x12 Pressure Treated", "LUM-6612-PT", "PCS", "Gable Lumber Supply", 35.00, 55.00, "Lumber"},

		// Sheet Goods
		{"3/4 Plywood CDX 4x8", "PLY-34-CDX", "PCS", "Gable Lumber Supply", 24.00, 38.00, "Sheet Goods"},
		{"1/2 Plywood CDX 4x8", "PLY-12-CDX", "PCS", "Gable Lumber Supply", 18.00, 28.50, "Sheet Goods"},
		{"3/4 OSB T&G 4x8", "OSB-34-TG", "PCS", "Gable Lumber Supply", 22.00, 34.00, "Sheet Goods"},
		{"1/2 OSB 4x8", "OSB-12", "PCS", "Gable Lumber Supply", 12.00, 19.50, "Sheet Goods"},
		{"1/2 Drywall Regular 4x8", "DW-12-REG", "PCS", "Gable Lumber Supply", 11.00, 16.00, "Sheet Goods"},
		{"5/8 Drywall Firecode 4x8", "DW-58-FC", "PCS", "Gable Lumber Supply", 14.00, 21.00, "Sheet Goods"},

		// Hardware & Fasteners
		{"16d Common Nails (50lb)", "NAIL-16D-50", "BOX", "Fastener Depot", 45.00, 65.00, "Hardware"},
		{"10d Common Nails (50lb)", "NAIL-10D-50", "BOX", "Fastener Depot", 45.00, 65.00, "Hardware"},
		{"3\" Deck Screws (5lb)", "SCR-DECK-3-5", "BOX", "Fastener Depot", 18.00, 29.99, "Hardware"},
		{"Joist Hanger 2x6", "HANGER-26", "PCS", "Hardware Wholesale Inc", 0.80, 1.45, "Hardware"},
		{"Joist Hanger 2x8", "HANGER-28", "PCS", "Hardware Wholesale Inc", 0.95, 1.65, "Hardware"},
		{"Joist Hanger 2x10", "HANGER-210", "PCS", "Hardware Wholesale Inc", 1.10, 1.85, "Hardware"},
		{"Hurricane Tie H1", "TIE-H1", "PCS", "Hardware Wholesale Inc", 0.65, 1.15, "Hardware"},
		{"Simpson Strong-Tie LUS28", "SIMP-LUS28", "PCS", "Hardware Wholesale Inc", 0.90, 1.50, "Hardware"},

		// Roofing
		{"Architectural Shingles (Black)", "RF-SH-BLK", "BUNDLE", "Roofing Specialists Ltd", 28.00, 42.00, "Roofing"},
		{"Architectural Shingles (Weathered Wood)", "RF-SH-WW", "BUNDLE", "Roofing Specialists Ltd", 28.00, 42.00, "Roofing"},
		{"Roofing Felt 15lb", "RF-FELT-15", "RL", "Roofing Specialists Ltd", 15.00, 22.50, "Roofing"},
		{"Ice & Water Shield", "RF-ICE-WTR", "RL", "Roofing Specialists Ltd", 65.00, 98.00, "Roofing"},
		{"Roof Edge Drip 10'", "RF-EDGE-WHT", "PCS", "Roofing Specialists Ltd", 4.50, 7.50, "Roofing"},

		// Insulation
		{"R-13 Fiberglass Batts 15x93", "INS-R13-15", "BAG", "Valley Insulation", 45.00, 68.00, "Insulation"},
		{"R-19 Fiberglass Batts 15x93", "INS-R19-15", "BAG", "Valley Insulation", 55.00, 82.00, "Insulation"},
		{"R-30 Fiberglass Batts 24x48", "INS-R30-24", "BAG", "Valley Insulation", 65.00, 98.00, "Insulation"},

		// Millwork
		{"Int Door 30x80 6-Panel Hollow", "DR-INT-3080-6P", "PCS", "Millwork Masters", 65.00, 95.00, "Millwork"},
		{"Int Door 32x80 6-Panel Hollow", "DR-INT-3280-6P", "PCS", "Millwork Masters", 65.00, 95.00, "Millwork"},
		{"Int Door 36x80 6-Panel Hollow", "DR-INT-3680-6P", "PCS", "Millwork Masters", 68.00, 99.00, "Millwork"},
		{"Ext Door 36x80 Steel 6-Panel", "DR-EXT-3680-STL", "PCS", "Millwork Masters", 180.00, 280.00, "Millwork"},
		{"Baseboard 3-1/4 MDF 16'", "MLD-BASE-MDF", "PCS", "Millwork Masters", 12.00, 19.50, "Millwork"},
		{"Casing 2-1/4 MDF 14'", "MLD-CASE-MDF", "PCS", "Millwork Masters", 8.00, 13.50, "Millwork"},
	}

	skuToID := make(map[string]uuid.UUID)
	productPrices := make(map[string]float64) // SKU -> Retail Price

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
			continue
		}

		pid := uuid.MustParse(id)
		skuToID[p.SKU] = pid
		productPrices[p.SKU] = p.Price

		// Inventory - varied quantities
		qty := 100 + rand.Intn(900)
		_, err = db.Exec(`
			INSERT INTO inventory (product_id, location_id, location, quantity)
			VALUES ($1, $2, 'MAIN', $3)
			ON CONFLICT (product_id, location_id) DO UPDATE SET quantity = $3
		`, pid, locationID, qty)
	}
	fmt.Printf("Seed: %d Products created\n", len(products))

	// 4. Customers & Projects (Expanded)
	customers := []struct {
		Name         string
		Acct         string
		Email        string
		Projects     []string
		CreditLimit  float64
		BalanceStart float64
	}{
		{"Acme Construction", "ACME-001", "billing@acme.com", []string{"Smith Residence", "Downtown Lofts", "City Park Gazebo"}, 50000, 12500},
		{"Bob's Builders", "BOB-001", "bob@bobsbuilders.com", []string{"Miller Deck", "Kitchen Remodel 123", "Garage Addition"}, 25000, 4200},
		{"DIY Homeowner", "DIY-888", "diy@gmail.com", []string{"Garden Shed"}, 5000, 0},
		{"Summit Contracting", "SUM-100", "ap@summitcontracting.com", []string{"Highland Hotel Reno", "Riverside Apts Bld A", "Riverside Apts Bld B"}, 150000, 45000},
		{"Elite Homes", "ELITE-202", "invoices@eliteframes.com", []string{"Lot 44 Oakwood", "Lot 45 Oakwood", "Lot 46 Oakwood"}, 75000, 28000},
		{"Prestige Decks", "PRES-303", "info@prestigedecks.com", []string{"Johnson Deck", "Peters Patio", "Clubhouse Veranda"}, 20000, 1500},
		{"Modern Renovations", "MOD-404", "pay@modernrenos.com", []string{"123 Main St Flip", "456 Elm St Flip"}, 30000, 8500},
		{"Classic Carpentry", "CLASS-505", "jim@classiccarpentry.com", []string{"Library Shelving", "Courthouse Trim"}, 15000, 500},
		{"Green Earth Landscapes", "GRN-606", "office@greenearth.com", []string{"Community Center Garden", "River Walk"}, 10000, 200},
		{"Structure Masters", "STR-707", "bills@structuremasters.com", []string{"Warehouse 9 Framing", "Retail Stripe Mall"}, 80000, 18000},
		{"Valley Roofing", "VAL-808", "admin@valleyroofing.com", []string{"School Roof Repair", "Church Shingle Replacement"}, 40000, 12000},
		{"Cornerstone Concrete", "CORN-909", "dispatch@cornerstone.com", []string{"Foundation Lot 8", "Driveway Smith"}, 35000, 6000},
	}

	customerIDs := make(map[string]uuid.UUID)
	projectIDs := make([]uuid.UUID, 0)
	custToProjects := make(map[uuid.UUID][]uuid.UUID)

	for _, c := range customers {
		var cid string
		err := db.QueryRow(`
			INSERT INTO customers (name, account_number, email, credit_limit, balance_due)
			VALUES ($1, $2, $3, $4, 0) -- Balance derived from invoices usually, but setting 0 base
			ON CONFLICT (account_number) DO UPDATE SET name=$1
			RETURNING id
		`, c.Name, c.Acct, c.Email, c.CreditLimit).Scan(&cid)

		if err != nil {
			log.Printf("Failed to upsert customer %s: %v", c.Name, err)
			continue
		}
		custID := uuid.MustParse(cid)
		customerIDs[c.Name] = custID
		custToProjects[custID] = []uuid.UUID{}

		// Projects
		for _, projName := range c.Projects {
			var jid string
			err := db.QueryRow(`
				INSERT INTO customer_jobs (customer_id, name, is_active)
				VALUES ($1, $2, true)
				RETURNING id
			`, custID, projName).Scan(&jid)
			if err == nil {
				pid := uuid.MustParse(jid)
				projectIDs = append(projectIDs, pid)
				custToProjects[custID] = append(custToProjects[custID], pid)
			}
		}
	}
	fmt.Printf("Seed: %d Customers created\n", len(customers))

	// 5. Orders & Invoices (Historical Data Generation)
	// Valid statuses: 'DRAFT', 'CONFIRMED', 'FULFILLED', 'CANCELLED'
	statuses := []string{"DRAFT", "CONFIRMED", "FULFILLED", "CANCELLED"}

	totalOrders := 0
	for custName, custID := range customerIDs {
		// Generate 3-8 orders per customer
		numOrders := 3 + rand.Intn(6)
		projects := custToProjects[custID]

		for i := 0; i < numOrders; i++ {
			status := statuses[rand.Intn(len(statuses))]
			// Favor fulfilled for history
			if rand.Float32() < 0.6 {
				status = "FULFILLED"
			}

			// Random Project - (Schema doesn't support job_id on orders yet, skipped)
			if len(projects) > 0 {
				_ = projects[rand.Intn(len(projects))]
			}

			// Use projID if schema supports it, for now we just needed valid ref in memory if needed.

			// Date
			orderDate := randomDate()

			// Create Order
			orderID := uuid.New()
			// Insert order header
			_, err := db.Exec(`
				INSERT INTO orders (id, customer_id, total_amount, status, created_at)
				VALUES ($1, $2, 0, $3, $4)
			`, orderID, custID, status, orderDate)

			if err != nil {
				log.Printf("Error creating order: %v", err)
				continue
			}

			// Add random lines (3-15 items)
			numLines := 3 + rand.Intn(13)
			var orderTotal float64

			for j := 0; j < numLines; j++ {
				// Pick random product
				prodIdx := rand.Intn(len(products))
				prod := products[prodIdx]
				qty := 1 + rand.Intn(50)
				price := prod.Price
				lineTotal := float64(qty) * price
				orderTotal += lineTotal

				// Insert Line
				db.Exec(`
					INSERT INTO order_lines (order_id, product_id, quantity, price_each)
					VALUES ($1, $2, $3, $4)
				`, orderID, skuToID[prod.SKU], qty, price)
			}

			// Update Order Total
			db.Exec("UPDATE orders SET total_amount = $1 WHERE id = $2", orderTotal, orderID)
			totalOrders++

			// Create Invoice if Fulfilled
			if status == "FULFILLED" {
				invID := uuid.New()
				invStatus := "ISSUED"
				if rand.Float32() < 0.7 {
					invStatus = "PAID"
				}

				dueDate := orderDate.AddDate(0, 1, 0) // Net 30

				_, err = db.Exec(`
					INSERT INTO invoices (id, order_id, customer_id, status, total_amount, due_date, created_at)
					VALUES ($1, $2, $3, $4, $5, $6, $7)
				`, invID, orderID, custID, invStatus, orderTotal, dueDate, orderDate.AddDate(0, 0, 1))

				if err == nil && invStatus == "PAID" {
					// Add Payment
					db.Exec(`
						INSERT INTO payments (invoice_id, amount, method, reference, notes)
						VALUES ($1, $2, 'CHECK', 'CHK-' || floor(random() * 10000 + 1000)::text, 'Payment in full')
					`, invID, orderTotal)
				}
			}
		}
		fmt.Printf("  Generated orders for %s\n", custName)
	}
	fmt.Printf("Seed: %d Total Orders created\n", totalOrders)

	// 6. RFCs (Governance) - More realistic
	rfcs := []struct {
		Title            string
		Status           string
		ProblemStatement string
	}{
		{"Standardize SKU Format", "draft", "Inconsistent data across yards"},
		{"Q3 Inventory Audit Procedure", "approved", "Need stricter control on lumber counts"},
		{"Vendor Onboarding Requirements", "review", "Compliance with new insurance regs"},
		{"Credit Limit Approval Workflow", "approved", "Automate approvals under $10k"},
		{"Safety Gear Mandatory List", "published", "Update per OSHA 2025 guidelines"},
		{"Returns Restocking Fee Policy", "draft", "Customer complaints on 15% fee"},
		{"Special Order Deposit Increase", "review", "Increase from 25% to 50% for non-stock"},
	}

	for _, rfc := range rfcs {
		_, err = db.Exec(`
			INSERT INTO rfcs (title, status, author_id, problem_statement, proposed_solution)
			VALUES ($1, $2, $3, $4, 'See attached document for detailed proposal.')
		`, rfc.Title, rfc.Status, demoUserID, rfc.ProblemStatement)
	}
	fmt.Println("Seed: Governance RFCs created")

	// 7. Portal User (for screenshots/demo)
	// Find Acme
	var acmeID string
	err = db.QueryRow("SELECT id FROM customers WHERE name = 'Acme Construction'").Scan(&acmeID)
	if err == nil {
		// Create Portal User
		pwHash, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
		_, err = db.Exec(`
			INSERT INTO customer_users (customer_id, email, password_hash, name, role)
			VALUES ($1, 'demo@gable.com', $2, 'Colton Demo', 'admin')
			ON CONFLICT (email) DO UPDATE SET password_hash = $2
		`, acmeID, string(pwHash))
		if err != nil {
			log.Printf("Failed to create portal user: %v", err)
		} else {
			fmt.Println("Seed: Portal User 'demo@gable.com' / 'password' created")
		}
	}

	fmt.Println("--------------------------------------------------")
	fmt.Println("DATABASE SEEDING COMPLETE FOR PROFESSIONAL DEMO")
	fmt.Println("--------------------------------------------------")
}
