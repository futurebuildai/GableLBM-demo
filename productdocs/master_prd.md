# GableLBM Master Product Requirements Document (PRD)

**Status**: Draft
**Version**: 1.0

## 1. Product Vision
**GableLBM** is an Open Operating System for the LBM (Lumber & Building Materials) industry. It rejects the complexity of legacy ERPs in favor of speed, ownership, and modern UX.
*   **Manifesto Principle**: "Real-Time is the Only Time."
*   **UX Standard**: "Bloomberg Terminal" density met with consumer-grade fluidity ("Industrial Dark" aesthetic).

## 2. Core Personas
1.  **Counter Carl (Sales Rep)**: Needs speed. Lives in the "Quick Quote" screen. Hates 15-click workflows.
2.  **Yardmaster Yvonne (Ops)**: Needs rugged tools. Uses a tablet in the rain. Needs "Big Buttons" and instant sync.
3.  **Owner Bob (Admin)**: Needs visibility. Wants to see "Total AR Exposure" and "Margin" instantly.

## 3. Functional Requirements (MVP Scope)

### 3.1 Inventory Management ("The Pile")
*   **Lumber Logic**: Native handling of Dimensional Lumber (2x4x8).
*   **Multi-UOM**: Must support simultaneous tracking of:
    *   **Piece** (Count)
    *   **MBF** (Thousand Board Feet - Volume)
    *   **LF** (Linear Feet - Mouldings)
*   **Bin Logic**: Simple "Zone/Row/Bin" location pointers.
*   **Audit**: Mandatory "Reason Codes" for any manual inventory adjustment (e.g., "Broken", "Shrinkage", "Found").

### 3.2 Sales & Order Entry ("The Transaction")
*   **Speed**: A standard 5-line order must be enterable in < 60 seconds.
*   **Search**: Omnibar style (Cmd+K) search for Customer and SKU simultaneously.
*   **Pricing Waterfall**:
    1.  Job Override (Highest priority)
    2.  Customer Specific Contract
    3.  Customer Group/Tier
    4.  Base Retail
*   **Cross-Sell**: AI-driven prompts (e.g., "Customer buying Decking -> Suggest Stainless Screws").

### 3.3 Logistics & Dispatch
*   **Pick Ticket**: Auto-generated PDF printed to specific "Zone Printers" in the yard.
*   **POD (Proof of Delivery)**: Mobile web view captures:
    *   GPS Timestamp.
    *   Photo of load.
    *   Signature.

### 3.4 Financials ("Ledger Lite")
*   **Invoicing**: Auto-generate Invoice PDF upon "Dispatched" status.
*   **AR Aging**: Simple bucket view (Current, 30, 60, 90+).
*   **Job Costing**: All transactions tagged to a specific "Project/Job" for the contractor.

## 4. Non-Functional Requirements
*   **Architecture**: Modular Monolith (Go + Postgres).
*   **Frontend**: React + Vite (SPA). No page reloads.
*   **Offline Mode**: Critical for Yard/Delivery apps (Optimistic UI updates).
*   **Performance**: < 100ms API response time for core transaction loops.

## 5. Success Metrics (Alpha)
1.  **Counter Speed**: < 60 seconds average "Quote-to-Ticket" time.
2.  **Inventory Accuracy**: 99% accuracy on "Cycle Count" verify.
3.  **Uptime**: Zero "Business Hours" downtime.
