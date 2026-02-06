# Sprint 5: Order Engine Tasks

## Database
- [x] Create migration `004_create_orders.sql` (Tables: orders, order_lines; Col: inventory.allocated) <!-- id: 0 -->
- [x] Apply migration locally <!-- id: 1 -->

## Backend Implementation
- [x] `internal/order/model.go` (Structs for Order, OrderLine, Status Enum) <!-- id: 2 -->
- [x] `internal/order/repository.go` (CRUD, Transaction support) <!-- id: 3 -->
- [x] `internal/inventory/service.go` (Add `AllocateStock` and `ReleaseStock` methods) <!-- id: 4 -->
- [x] `internal/order/service.go` (Business Logic: Create, ConvertQuote, UpdateStatus) <!-- id: 5 -->
- [x] `internal/order/handler.go` (HTTP Endpoints) <!-- id: 6 -->
- [x] Register Order routes in `cmd/api/main.go` <!-- id: 7 -->

## Frontend Implementation
- [x] `src/types/order.ts` (TS Interfaces) <!-- id: 8 -->
- [x] `src/services/order.service.ts` (API Client) <!-- id: 9 -->
- [x] `src/pages/orders/OrderList.tsx` (Table view of active orders) <!-- id: 10 -->
- [x] `src/pages/orders/OrderDetail.tsx` (View order details, status) <!-- id: 11 -->
- [ ] `src/pages/quotes/QuoteDetail.tsx` (Add "Convert to Order" button) <!-- id: 12 -->

## Verification
- [x] Verify Database Schema (psql) <!-- id: 13 -->
- [x] Verify API: Create Order via Curl <!-- id: 14 -->
- [x] Verify UI: Quote -> Order Flow <!-- id: 15 -->
- [x] Verify Inventory: Allocation updates "Available" count <!-- id: 16 -->
