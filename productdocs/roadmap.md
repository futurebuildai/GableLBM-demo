# GableLBM Product Roadmap

**Vision**: To build an Open Operating System for the construction supply chain that breaks the "Feature Trap" of legacy ERPs.

## Phase 1: Alpha MVP (The "Yard Shell")
**Goal**: Prove the core tech stack and "Inventory Truth".
**Target User**: Small Yard Owner / Yardmaster.
**Timeline estimate**: Q2 2026.

### Objectives
1.  **Tech Foundation**: [x] Go Backend + React Frontend + Postgres + NATS.
2.  **Core Inventory**:
    *   [x] SKU Management (Lumber specific: Grade, Species, Dim).
    *   [x] Multi-UOM Support (MBF <-> Piece <-> LF conversions).
    *   [x] Location Tracking (Bin/Row).
3.  **Basic "In/Out"**:
    *   [x] Simple Receipt (Add stock).
    *   [x] Simple Count (Adjust stock).

## Phase 2: Beta (The "Transaction Engine") - [DEMO READY MILESTONE]
**Current Focus**: Executing Fulfillment & Financials (Sprint 6/7).
**Goal**: Enable a full "Quote-to-Cash" workflow for a live counter sale.
**Target User**: Counter Sales Rep.
**Timeline estimate**: Q3 2026 (Demo Ready by end of Sprint 7).

### Objectives
1.  **Sales Order Processing**:
    *   [x] Order Engine (Backend Service & Schema).
    *   [x] Inventory Allocation (Soft Lock).
    *   [x] Inventory Fulfillment (Hard Deduct).
    *   [x] Quick Quote UI (< 3 clicks).
    *   [x] Customer Contract Pricing (Tiered Logic).
    *   [x] Credit Limit Checks.
2.  **Financials Lite**:
    *   [x] Invoice Engine (Backend & Schema).
    *   [x] Invoice UI (List & Detail).
    *   [x] Invoicing (Generate PDF).
    *   [x] Payment Collection (Stripe/Manual Entry).
    *   [x] Daily Till Reconciliation.
    *   [x] Sales Summary Reporting.
3.  **Logistics Lite**:
    *   [x] Pick Ticket Printing.
    *   [x] Basic Delivery Scheduling.

## Phase 3: General Availability (The "Ecosystem")
**Goal**: Enterprise-ready features and external integrations.
**Target User**: CFO / Ops Manager.
**Timeline estimate**: Q4 2026.

### Objectives
1.  **Integrations**:
    *   [x] General Ledger (QuickBooks/NetSuite) Sync (Foundation/Arch).
    *   [x] Vendor EDI (Purchasing - 850 Logic).
2.  **Advanced Logistics & Operations**:
    *   [x] Inventory Transfers (Multi-Location Sync).
    *   [x] Millwork Configurator (Doors/Windows).
    *   [x] Dispatch Route Optimization (Stop Reordering & Driver App).
3.  **Governance**:
    *   [x] Partner Portal (Co-op management - Sprint 13).
    *   [x] AI Governance Layer (RFC generation - Sprint 14).
    *   [x] Executive Analytics Dashboard (Real-time KPIs - Sprint 15).

## Phase 4: Open Ecosystem (The "Sovereign Dealer")
**Goal**: Enable true data sovereignty and 3rd-party innovation.
**Target User**: Tech Admin / Co-op IT Director.
**Timeline estimate**: Q1 2027.

### Objectives
1.  **Tech Admin Panel (Sprint 17-18)**:
    *   [ ] Self-service API Key Generation & Revocation.
    *   [ ] Webhook Configuration UI (No-code event triggers).
    *   [ ] System Health & Logs Dashboard.
2.  **Integration Marketplace**:
    *   [ ] "One-Click" Run Payments Integration.
    *   [ ] "One-Click" QuickBooks Online OAuth2 Sync.
    *   [ ] Zapier/Make Connector.
3.  **Data Liberation**:
    *   [ ] Full Database Export (SQL Dump) from Admin UI.
    *   [ ] Schema Viewer for easy customized reporting.
4.  **Co-op enablement**:
    *   [ ] "Fork & Brand" Documentation.
    *   [ ] Multi-tenant Architecture options for Hosted offerings.
