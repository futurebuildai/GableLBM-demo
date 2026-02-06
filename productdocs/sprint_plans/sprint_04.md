# Sprint 04: The Transaction Engine (Part 1 - Quotes)

**Goal**: Establish the Customer/Order schema and enable the creation of a basic "Quick Quote".
**Prerequisites**: Sprint 03 (Inventory) complete.

## 1. Scope

### 1.1 Customer Foundation
*   **Schema**: `customers`, `customer_jobs`, `price_levels`.
*   **Seeding**: Default "Retail" price level.
*   **Logic**: Simple Create/Read for customers.

### 1.2 Quote Management
*   **Schema**: `quotes`, `quote_lines`.
*   **Logic**:
    *   Create Quote for Customer.
    *   Add Line Item (Product + Qty).
    *   Calculate Totals (Price * Qty).
*   **UI**:
    *   "Quick Quote" Dashboard Widget.
    *   Full Quote Builder Page.

## 2. Technical Requirements
*   **Backend**:
    *   `CustomerService` (CRUD).
    *   `QuoteService` (Create, AddLine, Finalize).
*   **Frontend**:
    *   `CustomerSelect` Component (Async Search).
    *   `QuoteBuilder` Component (Master-Detail view).

## 3. Success Criteria
*   [ ] Can create a Customer ("Bob the Builder").
*   [ ] Can create a Quote for Bob.
*   [ ] Can add 100 2x4s to the Quote.
*   [ ] Can see the Total Price calculated correctly.
