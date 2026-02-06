# Sprint 12: Integrations Foundation (GL & EDI)

## Goal
Bridge the gap between GableLBM's operational data and external Enterprise systems (Financials & Supply Chain).

## Objectives

### 1. General Ledger (GL) Sync
- **Goal**: Financial transactions in GableLBM must reflect in an external Accounting System.
- **Scope**:
    - Abstraction Layer (`GLAdapter`) to support multiple ERPs (QBO, NetSuite, etc.).
    - Event Trigger: `InvoiceFinalized` -> `Dr. AR / Cr. Revenue`.
    - Event Trigger: `PaymentReceived` -> `Dr. Cash / Cr. AR`.

### 2. Vendor EDI (Purchasing)
- **Goal**: Automate ordering with big suppliers (Lumber Mills/Wholesalers).
- **Scope**:
    - **X12 850**: Purchase Order Outbound. Generate a valid ANSI X12 850 string from a Gable PO.
    - **Transport**: Stub implementation for FTP/SFTP upload.

## Technical Constraints
- **Separation of Concerns**: The core domain (`Invoice`, `Order`) must NOT know about QuickBooks or EDI internals. Use Interfaces/Adapters.
- **Idempotency**: Retrying a GL sync must not create duplicate Journal Entries.

## Dependencies
- `InvoiceService` (Existing)
- `PurchaseOrderService` (Existing - Sprint 8/9?)
