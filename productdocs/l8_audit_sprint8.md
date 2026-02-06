# L8 Production Readiness Audit Findings (Sprint 8)

## Executive Summary
**Status**: 🔴 NOT READY FOR PRODUCTION
**Auditor**: Antigravity (L8 SRE)
**Date**: 2026-02-05

While the functional requirements for "Financials Lite" have been met in the "Happy Path", the current implementation contains critical reliability and data integrity flaws that would cause significant issues at scale or under failure conditions.

## Critical Findings (P0)

### 1. Lack of Transactional Integrity (ACID Violation)
*   **Location**: `backend/internal/payment/service.go:ProcessPayment`
*   **Issue**: The payment creation (`repo.CreatePayment`) and invoice status update (`invoiceRepo.UpdateInvoice`) occur in two separate database calls.
*   **Risk**: If the server crashes or the second call fails after the first succeeds, the system will be in an inconsistent state: Money is recorded as taken, but the Invoice remains "UNPAID".
*   **Remediation**: Wrap both operations in a single Serializable Transaction.

### 2. Financial Type Safety (`float64` vs `DECIMAL`)
*   **Location**: `backend/internal/payment/model.go`, `backend/internal/invoice/model.go`
*   **Issue**: The Go backend uses `float64` for monetary values, while the Postgres database correctly uses `DECIMAL(10,2)`.
*   **Risk**: Floating point arithmetic errors (IEEE 754) will eventually cause "off-by-a-penny" errors in financial reporting and balancing. $0.1 + $0.2 != $0.3 in float logic.
*   **Remediation**: Transition Go structs to use `int64` (cents) or a dedicated fixed-point decimal library (e.g., `shopspring/decimal`), or strictly handle rounding at boundaries. For MVP, at minimum, documentation of this technical debt is required, but for L8/Production it is blocked.

### 3. Synchronous Email Notification
*   **Location**: `backend/internal/document/handler.go:HandleEmailInvoice`
*   **Issue**: Email sending happens synchronously within the HTTP Release handler.
*   **Risk**: If the SMTP server hangs, the User Request hangs. If it fails, the user gets an error, but does the process retry? No.
*   **Remediation**: Move email dispatch to a background worker/queue (Outbox Pattern or Task Queue).

## Performance Findings (P1)

### 1. Reporting Aggregation on Read
*   **Location**: `backend/internal/reporting/repository.go`
*   **Issue**: `GetSalesSummary` runs `SUM()` and `COUNT()` aggregates over the raw `payments` and `invoices` tables every time the dashboard is loaded.
*   **Risk**: As data grows to 100k+ rows, this dashboard will become unresponsive.
*   **Remediation**: Implement Materialized Views for daily stats or pre-calculate these values during the "End of Day" process.

## Security Findings (P1)

### 1. Missing Authorization for Payments
*   **Location**: `backend/internal/payment/handler.go`
*   **Issue**: Any authenticated user can hit `/api/payments`.
*   **Risk**: A warehouse worker could technically post a payment if they knew the API endpoint.
*   **Remediation**: Enforce RBAC. Only users with `FINANCE` or `ADMIN` roles should execute `ProcessPayment`.

## Remediation Plan (Sprint 9 or Immediate Hotfix)

1.  **Transactional Fix**: Refactor `Repository` interface to support passing a `tx` (context/transaction) to link `CreatePayment` and `UpdateInvoice`.
2.  **Type Safety**: Switch to `int64` for all internal money logic (cents) or strictly validate.
3.  **Async Email**: Simple go-routine with channel for now (MVP), proper queue later.
