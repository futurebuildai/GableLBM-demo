# Sprint 07: The Counter Power-Up

**Goal**: Transform the transaction flow into a professional-grade "Counter Sales" experience with advanced pricing, credit checks, and document generation.

## 1. Scope

### 1.1 Quick Quote UI (< 3 Clicks)
*   **Omnibar**: A global (Cmd+K) search component that searches both Customers and Products simultaneously.
*   **Keyboard-First Flow**: Focus on tab-index and hotkeys for rapid line-item entry.
*   **Quote-to-Order**: Finalize the "Convert" button logic with proper state preservation.

### 1.2 Pricing Waterfall Logic
*   **Pricing Service**: Backend logic to calculate the "Best Price" for a customer.
*   **Hierarchy**: 
    1. Customer Contract (explicit SKU price).
    2. Customer Group/Tier discount (e.g., "Gold" = 10% off).
    3. Base Retail.
*   **UI**: Show "Discount Applied" or "Contract Price" indicator next to line items.

### 1.3 Credit & Risk
*   **Credit Limit**: Implement `CreditLimit` on Customer model.
*   **Validation**: Block or flag orders that exceed the available credit.
*   **Visuals**: "Safety Red" banner for customers on credit hold.

### 1.4 Industrial Printing (PDFs)
*   **Service**: Go-based PDF generation for:
    *   **Invoices** (Financial record).
    *   **Pick Tickets** (Yard instructions).
*   **Frontend**: Add "Print" buttons to Order and Invoice detail views.

## 2. Technical Requirements
*   **Backend**:
    *   Migration: `006_pricing_and_credit.sql`.
    *   `internal/pricing`: New service for waterfall logic.
    *   `internal/document`: New service for PDF generation.
*   **Frontend**:
    *   `Omnibar` component.
    *   Update `QuoteBuilder` with pricing logic visibility.
    *   `CustomerDetail`: Edit credit limits.

## 3. Success Criteria
*   [ ] Can enter a 3-item quote using only the keyboard in < 30 seconds.
*   [ ] Pricing correctly switches when a "Tier 1" customer is selected.
*   [ ] System blocks fulfillment if customer is over credit limit.
*   [ ] Professional PDF generated for both Invoice and Pick Ticket.
