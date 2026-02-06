# GableLBM Product Roadmap

**Vision**: To build an Open Operating System for the construction supply chain that breaks the "Feature Trap" of legacy ERPs.

## Phase 1: Alpha MVP (The "Yard Shell")
**Goal**: Prove the core tech stack and "Inventory Truth".
**Target User**: Small Yard Owner / Yardmaster.
**Timeline estimate**: Q2 2026.

### Objectives
1.  **Tech Foundation**: Go Backend + React Frontend + Postgres + NATS.
2.  **Core Inventory**:
    *   SKU Management (Lumber specific: Grade, Species, Dim).
    *   Multi-UOM Support (MBF <-> Piece <-> LF conversions).
    *   Location Tracking (Bin/Row).
3.  **Basic "In/Out"**:
    *   Simple Receipt (Add stock).
    *   Simple Count (Adjust stock).

## Phase 2: Beta (The "Transaction Engine")
**Goal**: Enable a full "Quote-to-Cash" workflow for a live counter sale.
**Target User**: Counter Sales Rep.
**Timeline estimate**: Q3 2026.

### Objectives
1.  **Sales Order Processing**:
    *   Quick Quote UI (< 3 clicks).
    *   Customer Contract Pricing (Tiered Logic).
    *   Credit Limit Checks.
2.  **Financials Lite**:
    *   Invoicing (Generate PDF).
    *   Payment Collection (Stripe/Manual Entry).
    *   Daily Till Reconciliation.
3.  **Logistics Lite**:
    *   Pick Ticket Printing.
    *   Basic Delivery Scheduling.

## Phase 3: General Availability (The "Ecosystem")
**Goal**: Enterprise-ready features and external integrations.
**Target User**: CFO / Ops Manager.
**Timeline estimate**: Q4 2026.

### Objectives
1.  **Integrations**:
    *   General Ledger (QuickBooks/NetSuite) Sync.
    *   Vendor EDI (Purchasing).
2.  **Advanced Components**:
    *   Millwork Configurator (Doors/Windows).
    *   Dispatch Route Optimization.
3.  **Governance**:
    *   Partner Portal (Co-op management).
    *   AI Governance Layer (RFC generation).
