# LBM ERP Feature Specification (Draft)

## 1. Inventory Management (The Core)
The lumber industry relies on complex unit conversions and tracking. The system must handle this natively, not as an afterthought.

### 1.1. Dimensional Logic & Unit of Measure (UOM)
- **Multi-UOM Support:** Every item must support buying, stocking, and selling in different units (e.g., Buy in MBF, Stock in Pieces, Sell in Linear Feet).
- **Tally Management:** Support for "Tally Assessment" on receipt. (e.g., Receiving a random width/length unit of hardwood and converting it to inventory on hand).
- **Remanufacturing (Reman):** Native workflows for "breaking bulk" (e.g., turning a bunk of 2x4x16s into Studs) or "Building Kits" (Door packs).

### 1.2. Warehouse & Yard
- **Multi-Location:** Single SKU existing in multiple physical spots (Shed A, Row 4 + Yard Overflow).
- **Lot & Serial Tracking:** Mandatory for Engineered Wood Products (EWP) and special order Windows/Doors.
- **Barcode & RFID:** Optimization for rugged handhelds in the yard.

## 2. Sales & Order Management
The sales cycle in LBM is rarely a simple "Cash and Carry". It is relationship and project-based.

### 2.1. Quoting & Estimating
- **Project-Based Quotes:** Quotes tied to a specific "Job" or "Project", not just a customer.
- **Bid Management:** Version control for quotes as blueprints change.
- **Margin Management:** Real-time margin calculation based on current replacement cost vs. average cost.

### 2.2. Contractor Pricing Engine
- **Hierarchical Pricing:**
    - Base Level (Retail)
    - Customer Level (Contractor A is Level 3)
    - Job Level (This specific Hotel Project gets special pricing on Gypsum)
    - SKU Net Pricing (Hard override on 2x4s)

### 2.3. Point of Sale (POS)
- **Fast Counter Checkout:** Designed for high speed (Scanning, Quick Lookup).
- **Account Charges:** Authorized signer verification (Signature capture, ID check).
- **Returns validation:** Original invoice checking to prevent fraud.

## 3. Logistics & Dispatch
Delivery is the lifeblood of an LBM dealer.

### 3.1. Load Building & Routing
- **Truck Capacity:** Validation for weight and physical space (can this 40' I-Joist fit on the flatbed?).
- **Route Optimization:** AI-assisted routing for multi-stop runs.

### 3.2. Driver App
- **Electronic Proof of Delivery (ePOD):** Photos of the drop site (critical for damage claims), GPS stamp, Signature.
- **Turn-by-Turn:** Integrated navigation.

## 4. Financials & Accounting
LBM construction finance is distinct from generic retail.

### 4.1. Job Accounting
- **AIA Billing:** Support for progress billing standards.
- **Lien Management:** Tracking "Notice to Owner" dates and lien releases.

### 4.2. Credit Management
- **Credit Limits:** Hard vs. Soft stops.
- **Job Accounts:** Segregated credit limits per working job.

## 5. AI-Native Functionality
How we differentiate from laggy legacy systems.

- **Computer Vision:** Inbound receiving - scan a packing slip or the lumber load itself to verify tally.
- **Generative Sales:** "Chat with your Data" - "What's the price history for 5/8 OSB for Smith Construction over the last 6 months?"
- **Predictive Purchasing:** "Weather forecast predicts heavy rain; usually roofing felt sales spike 3 days prior. Suggest ordering +20%."
