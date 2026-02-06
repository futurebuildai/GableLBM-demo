# Sprint 08 Task List

## Phase 1: Database & Backend Core
- [x] Create migration `008_payments_and_till.sql` (Payments table, Invoice update) <!-- id: 1 -->
- [x] Apply migration locally <!-- id: 2 -->
- [x] Implement `internal/payment` model and repository <!-- id: 3 -->
- [x] Implement `internal/payment` service (process payment, update invoice) <!-- id: 4 -->

## Phase 2: Payment Integration
- [x] Create `internal/payment/gateway` interface <!-- id: 5 -->
- [x] Implement Mock Gateway (for dev/testing) <!-- id: 6 -->
- [x] Add `/api/payments` endpoints (POST pay, GET history) <!-- id: 7 -->
- [x] Update `internal/invoice` to handle status changes <!-- id: 8 -->

## Phase 3: Frontend Payment UI
- [x] Create `PaymentModal` component (Amount, Method) <!-- id: 9 -->
- [x] Add "Pay" button to `InvoiceDetail` view <!-- id: 10 -->
- [x] Implement Payment History view in `InvoiceDetail` <!-- id: 11 -->
- [x] Visual handling for Partial/Paid statuses <!-- id: 12 -->

## Phase 4: Financial Dashboard (Daily Till)
- [x] Implement backend `TillService` (aggregates daily payments) <!-- id: 13 -->
- [x] Create `/api/reports/daily-till` endpoint <!-- id: 14 -->
- [x] Create `DailyTill` page in Frontend <!-- id: 15 -->
- [x] Implement Sales Summary Report capability <!-- id: 16 -->

## Phase 5: Invoice Emailing & PDF Links
- [x] Setup simple SMTP/Mockmailer in backend <!-- id: 17 -->
- [x] Add "Pay Now" link generation in PDF template <!-- id: 18 -->
- [x] Implement `/api/invoices/{id}/email` endpoint <!-- id: 19 -->
- [x] Add "Email Invoice" button in frontend <!-- id: 20 -->

## Phase 6: Refinement & Polish
- [ ] Optimize Omnibar (debounce, indexes) <!-- id: 21 -->
- [ ] Create `ShortcutsModal` component <!-- id: 22 -->
- [ ] Register global `?` keybind <!-- id: 23 -->

## Phase 7: Verification & Audit
- [ ] Run backend tests for payments <!-- id: 24 -->
- [ ] Manual E2E: Invoice -> Partial Pay -> Full Pay -> Till Check <!-- id: 25 -->
- [ ] Run L8 Production Readiness Gate Audit <!-- id: 26 -->
- [ ] Update Roadmap and Walkthrough <!-- id: 27 -->
