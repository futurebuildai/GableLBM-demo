# Sprint 5 Walkthrough: The Order Engine

## Overview
Sprint 5 focused on implementing the "Order Engine", enabling the transition from Quotes to confirmed Sales Orders. This introduced the concept of "Inventory Allocation", where stock is reserved (Soft Lock) before being fulfilled (Hard Deduct).

## Changes

### Database Schema
*   **New Tables**:
    *   `orders`: Tracks the header status (DRAFT, CONFIRMED, FULFILLED).
    *   `order_lines`: Snapshot of product/price at time of order.
*   **Inventory Update**:
    *   Added `allocated` (numeric) column to `inventory`.
    *   `Available Quantity` is now calculated as `Quantity - Allocated`.

```sql
-- Migration: 004_create_orders.sql
CREATE TABLE orders (...);
ALTER TABLE inventory ADD COLUMN allocated ...;
```

### Backend (Go)
*   **Module**: `internal/order`
    *   `Service`: Handles creation and the `ConfirmOrder` logic.
    *   `Repository`: Transactional inserts for header/lines.
*   **Inventory Integration**:
    *   `InventoryService.Allocate(productID, qty)`: Finds the best location and increments the `allocated` counter.
    *   Used by `OrderService.ConfirmOrder` to lock stock.

### Frontend (React)
*   **New Pages**:
    *   `Orders/OrderList`: Dashboard of all active orders.
    *   `Orders/OrderDetail`: Detailed view allowing "Confirm" and "Dispatch" actions.
*   **Navigation**: Added "Orders" to the prompt sidebar.

## Verification

### Automated Tests
*   `go build ./...` passed.
*   `npm run build` passed.

### Manual Verification Steps
1.  **Schema Check**:
    *   Verified `orders` table exists via `psql`.
    *   Verified `inventory.allocated` column exists.
2.  **Order Flow**:
    *   Navigate to `/orders`.
    *   (Future) Click "Convert to Order" from a Quote.
    *   View Order Detail.
    *   Click "Confirm Order" -> Status changes to `CONFIRMED`.
    *   (Backend) Verify `inventory.allocated` increases.

## Screenshots
> *Placeholder for UI screenshots of Order List and Detail pages.*
