# GableLBM: Alpha MVP Scope Definition

The goal of the Alpha MVP is to provide a "Minimum Viable Yard" experience—enough for a small dealer to manage basic operations without falling back to paper or spreadsheets.

## 1. Core Inventory (The "Pile" Manager)
- **SKU Basic Data:** Name, Category, Weight, Base UOM.
- **LBM Conversions:** Piece/MBF conversion for at least 5 common lumber sizes.
- **Receipt from PO:** Basic inventory increment workflow.
- **Adjustments:** "Cycle count" adjustment tool with mandatory reason codes.
- **Location Basics:** Bin-level tracking for a single yard.

## 2. Sales & Counter UI (The "Transaction" Engine)
- **Customer Lookup:** Search by Name or Phone.
- **Pricing Waterfall:** 
    - Base Retail.
    - Tiered Pricing (levels 1-3).
    - Job-specific override.
- **Quote-to-Order:** Convert a draft quote into an active Sales Order.
- **Picking Slip:** Print-friendly (or PDF) pick ticket for the yard crew.

## 3. Financials (The "Ledger" Lite)
- **Invoicing:** Generate invoice PDF on order delivery.
- **Payment Entry:** Record Check, Cash, or Credit card (manual entry).
- **Basic AR Aging:** View total due per customer in 30/60/90 buckets.
- **Job Costing:** Group invoices/payments by "Job" for project tracking.

## 4. Deferred (Not in Alpha)
- **Automatic PO Generation:** (Manual PO entry only).
- **Route Optimization:** (Manual dispatch queuing only).
- **EDI/Vendor Integrations:** (CSV import only).
- **Multi-Yard Transfers:** (Single location focus for pilot).

---

## Success Criteria for Alpha
1. A counter rep can finish a sale in under 60 seconds.
2. An owner can see their total AR exposure at a glance.
3. Inventory counts "Piece count" and "MBF" simultaneously without manual calculation.
