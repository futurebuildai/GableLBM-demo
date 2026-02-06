# Sprint 8: Financials Lite & Payment Integration

## Goal
Complete the "Quote-to-Cash" cycle by implementing Payment Collection, Daily Reconciliation, and advanced Invoice Management.

## Objectives
1.  **Payment Processing**:
    *   [ ] Implement Stripe Integration (or mock for now).
    *   [ ] Record Payments against Invoices (Partial/Full).
    *   [ ] Update Invoice Status logic (UNPAID -> PARTIAL -> PAID).

2.  **Financial Dashboard**:
    *   [ ] "Daily Till" View for Counter Reps.
    *   [ ] Sales Summary Report (Daily/Weekly).

3.  **Invoice Enhancements**:
    *   [ ] Email Invoice to Customer (SMTP/SendGrid).
    *   [ ] Add "Pay Now" link to Invoice PDF.

4.  **Refinement**:
    *   [ ] Optimize Omnibar performance with large datasets.
    *   [ ] Add Keyboard Shortcuts reference modal.

## Technical Constraints
*   **Security**: Ensure PCI compliance (or use hosted fields) if touching card data.
*   **Consistency**: Payment recording must use transactional integrity with Invoice updates.
