# Sprint 06: Fulfillment & Financials

**Goal**: Complete the "Quote-to-Cash" lifecycle by enabling Order Fulfillment (Pick/Ship) and basic Invoicing.

## 1. Scope

### 1.1 Fulfillment Logic ("The Yard")
*   **Action**: `FulfillOrder(orderID)`
*   **State Transition**: `CONFIRMED` -> `FULFILLED`
*   **Inventory Impact**:
    *   `OnHand` decreases by Order Qty.
    *   `Allocated` decreases by Order Qty.
    *   *Constraint*: Cannot fulfill if `OnHand < Qty`.

### 1.2 Financials Lite (Invoicing)
*   **Schema**: `invoices`, `invoice_lines`.
*   **Trigger**: Auto-generate Invoice upon `FULFILLED` status.
*   **Data**: Snapshot of pricing and quantity at time of fulfillment.

### 1.3 UI Updates
*   **Order Detail**: Add "Fulfill Order" action (with confirmation).
*   **Invoice View**: Read-only view of the generated invoice.
*   **Inventory View**: Explicitly show `Allocated` vs `Available` columns.

## 2. Technical Requirements
*   **Backend**:
    *   Migration: `005_create_invoices.sql`.
    *   `OrderService.Fulfill()`: Transactional inventory update.
    *   `InvoiceService`: Create, Get.
*   **Frontend**:
    *   `OrderDetail`: Fulfill Button.
    *   `InvoiceList` / `InvoiceDetail` pages.
    *   Update `InventoryTable` to show allocation breakdown.

## 3. Success Criteria
*   [ ] Clicking "Fulfill" on an Order updates Inventory correctly (OnHand reduces, Allocated reduces).
*   [ ] Order Status changes to `FULFILLED`.
*   [ ] An Invoice is created automatically.
*   [ ] Can view the Invoice details.
