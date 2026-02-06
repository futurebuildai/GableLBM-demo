# Sprint 7 Walkthrough: The Counter Power-Up

## Overview
Sprint 7 transformed GableLBM from a basic order system into a capable Point-of-Sale (POS) backend. We introduced tiered pricing logic, credit limit controls, and professional PDF document generation for Invoices and Pick Tickets. The frontend also received a major "Power-Up" with the new `Omnibar` for rapid navigation and search.

## Changes

### 1. Pricing Engine (Waterfall)
*   **Logic**: `CalculatePrice(customer, product)`
*   **Waterfall**:
    1.  **Contract Price**: Specific negotiated price for Customer+Product.
    2.  **Tier Price**: Percentage discount off Base Retail (Silver 10%, Gold 15%, etc).
    3.  **Base Retail**: Default fallback.
*   **Database**: Added `tier` enum to Customers and `customer_contracts` table.

### 2. Credit Control
*   **Backend**: `OrderService` now checks `(Balance + OrderTotal) > CreditLimit` before fulfillment.
*   **Frontend**: Visual indicators for "Over Limit" and "Credit Hold" in Quote Builder and Search.

### 3. Document Engine (PDF)
*   **Library**: Integrated `maroto` for Golang PDF generation.
*   **Endpoints**:
    *   `GET /api/documents/print/invoice/{id}`
    *   `GET /api/documents/print/pickticket/{id}`
*   **UI**: Added "Print" buttons to Order and Invoice detail pages.

### 4. Frontend "Omnibar"
*   **Feature**: Global `Cmd+K` (or `Ctrl+K`) search bar.
*   **Capabilities**:
    *   Search Products instantly (SKU/Desc).
    *   Search Customers (showing Credit Status).
    *   Quick Navigation commands.

## Verification

### Automated Tests
*   `go test ./internal/pricing/...`: Validated Waterfall logic (Contract > Tier > Retail).
*   `npm run build`: Confirmed frontend type safety with new PDF buttons and Omnibar.

### Manual Verification Flow
1.  **Pricing**:
    *   Selected a "Silver Tier" customer in Quote Builder.
    *   Verified prices reflected 10% discount compared to Retail.
2.  **Credit Block**:
    *   Attempted to fulfill an order exceeding the $10,000 limit.
    *   Verified backend blocked the transaction (simulated).
    *   Observed red "Exceeds Limit" warning in UI.
3.  **Documents**:
    *   Clicked "Print Pick Ticket" on a Confirmed Order -> Opened PDF in new tab.
    *   Clicked "Download PDF" on an Invoice -> Downloaded correctly formatted Invoice.
4.  **Omnibar**:
    *   Pressed `Cmd+K`.
    *   Typed "2x4".
    *   Selected Product -> (Future: Add to Cart).
