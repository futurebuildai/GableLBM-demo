# Sprint 07 Task List

## Phase 1: Context & Database
- [x] Create migration `006_pricing_and_credit.sql` <!-- id: 0 -->
- [x] Apply migration locally <!-- id: 1 -->
- [x] Create migration `007_add_product_price.sql` (Fix: Missing base price) <!-- id: 20 -->
- [x] Apply migration 007 <!-- id: 21 -->

## Phase 2: Backend Pricing & Credit
- [x] Implement `internal/pricing` model and repository <!-- id: 2 -->
- [x] Implement `internal/pricing` service (Waterfall logic) <!-- id: 3 -->
- [x] Update `internal/order` service to include credit limit checks <!-- id: 4 -->
- [x] Implement `internal/customer` updates for credit limit/tier <!-- id: 5 -->

## Phase 3: Document Engine (PDFs)
- [x] Setup PDF generation library in backend <!-- id: 6 -->
- [x] Implement Pick Ticket PDF template <!-- id: 7 -->
- [x] Implement Invoice PDF template <!-- id: 8 -->
- [x] Add `/api/documents/print/{type}/{id}` endpoint <!-- id: 9 -->

## Phase 4: Frontend Counter Power-Up
- [x] Implement `Omnibar` (Global Cmd+K search) <!-- id: 10 -->
- [x] Refactor `QuoteBuilder` for keyboard-first entry <!-- id: 11 --> 
- [x] Add Pricing visibility to `QuoteBuilder` line items <!-- id: 12 -->
- [x] Add "Print" buttons to `OrderDetail` and `InvoiceDetail` <!-- id: 13 -->
- [x] Implement Credit Hold visual indicators <!-- id: 14 -->

## Phase 5: Verification & L8 Audit
- [ ] Run backend tests for pricing logic <!-- id: 15 -->
- [ ] Verify PDF generation output <!-- id: 16 -->
- [ ] Manual E2E test: Quote -> Pricing -> Order -> Credit Block -> Fulfill -> PDF <!-- id: 17 -->
- [ ] Run L8 Production Readiness Gate Audit <!-- id: 18 -->
- [ ] Update Roadmap and Walkthrough <!-- id: 19 -->
