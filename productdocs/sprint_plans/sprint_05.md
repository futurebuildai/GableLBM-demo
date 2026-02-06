# Sprint 05: The Order Engine

**Goal**: Enable the conversion of Quotes into Sales Orders and track Inventory Commitment.

## 1. Scope

### 1.1 Order Management
*   **Schema**: `orders`, `order_lines`.
*   **State Machine**:
    *   `DRAFT`: Being built.
    *   `CONFIRMED`: Stock allocated, ready for picking.
    *   `FULFILLED`: Shipped/Picked up (deduct from OnHand).
    *   `CANCELLED`: Void.
*   **Logic**:
    *   `ConvertQuoteToOrder(quoteID)`: Copies lines, locks pricing.
    *   `CommitStock(orderID)`: Increases `Allocated` count for SKUs.

### 1.2 Inventory Updates (Allocation)
*   Update `Inventory` model to track `Allocated` quantity.
*   Formula: `Available` = `OnHand` - `Allocated`.

## 2. Technical Requirements
*   **Backend**:
    *   Migration: `004_create_orders.sql`.
    *   Migration: Update `inventory` with `allocated` column.
    *   `OrderService`: CRUD, State Transitions.
    *   `InventoryService`: Add `Allocate(sku, qty)` method.
*   **Frontend**:
    *   `OrderList` page.
    *   `OrderDetail` page.
    *   "Convert to Order" Action on Quote.

## 3. Success Criteria
*   [ ] Can convert a Quote to an Order.
*   [ ] Confirming an Order reduces "Available" stock for that SKU.
*   [ ] Order appears in "Active Orders" list.
