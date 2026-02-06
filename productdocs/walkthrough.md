# Sprint 8 Walkthrough: Financials & Refinement

## Overview
Sprint 8 completed the "Quote-to-Cash" journey. We implemented end-to-end Payment Processing, a Financial Dashboard (Daily Till), and "Quality of Life" improvements like Invoice Emailing and Keyboard Shortcuts.

## Changes

### 1. Payment Processing
*   **Backend**: `internal/payment` module handling Payments, Invoice Status updates (`PARTIAL`, `PAID`), and Ledger updates.
*   **Database**: Added `payments` table and migration.
*   **UI**:
    *   **Payment Modal**: Allows recording payments (Check, Cash, Credit) against invoices.
    *   **Status Tracking**: Invoices automatically transition `UNPAID` -> `PARTIAL` -> `PAID`.
    *   **History**: Transactions listed on Invoice Detail.

### 2. Financial Dashboard (Daily Till)
*   **New Page**: `/reports/daily-till` (Accessible via Sidebar).
*   **Features**:
    *   **Daily Till**: Total collected today breakdown by method (Cash vs Card).
    *   **Sales Summary**: 30-day lookback at Invoiced vs Collected (Collection Rate %).
    *   **Outstanding AR**: Total Open Invoices in period.

### 3. Invoice "Quality of Life"
*   **Emailing**: Added "Email Invoice" button (`/api/invoices/{id}/email` mock integration).
*   **PDF Links**: Generated PDFs now include a clickable "Pay Now" link (mock).

### 4. Refinement (The Polish)
*   **Shortcuts**:
    *   Added global `?` shortcut to view keybinds.
    *   `Cmd+K`: Omnibar.
    *   `G+D`: Goto Dashboard.

## Verification

### Automated Tests
*   `go test ./...`: Passed (Logic & integration).
*   `npm run build`: Passed (Type safety & component integration).

### Manual Verification Flow
1.  **Payment Flow**:
    *   Created Invoice for $1,000.
    *   Paid $500 (Partial) -> Invoice Status changed to `PARTIAL`.
    *   Paid $500 (Remaining) -> Invoice Status changed to `PAID`.
2.  **Reporting**:
    *   Navigated to "Daily Till".
    *   Verified the $1,000 showed up in "Today's Collections".
3.  **Shortcuts**:
    *   Pressed `?` -> Modal appeared.
    *   Pressed `Cmd+K` -> Omnibar appeared.
