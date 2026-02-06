# Sprint 13: Partner Portal (Co-op Management)

## Goal
Empower "Pro" customers (Contractors/Builders) with self-service access to their account data, reducing call volume for "Counter Carl" and enabling 24/7 operations.

## Objectives

### 1. External Authentication (Security Foundation)
- **Goal**: Securely authenticate external users (Contacts) distinct from internal staff.
- **Scope**:
    - "Contact" entity becomes verify-able user.
    - JWT-based auth with separate "Partner" audience/scope.
    - Zero-Trust: APIs explicitly whitelist "Partner" role.

### 2. Partner Dashboard
- **Goal**: "At a Glance" financial and operational health for the Contractor.
- **Scope**:
    - **Current Balance**: Total exposure (AR).
    - **Active Jobs**: List of projects with open orders.
    - **Recent Activity**: Stream of recent invoices/quotes.

### 3. "Co-op" Workflow (Quote Approval)
- **Goal**: Allow contractors to approve quotes remotely.
- **Scope**:
    - View Quote Details (Line items, prices, expiration).
    - Action: "Approve & Order" -> Converts Quote to Order in backend.
    - Action: "Request Change" -> a simple note/status update.

## Technical Constraints
- **Security**: External users MUST NEVER access data belonging to another account. Strict Tenant/Account ID filtering on EVERY query.
- **UX**: Distinct visual theme (maybe "Slate Blue" accents vs "Lumber Orange") to distinguish from Internal Admin.

## Dependencies
- `AccountService` (Existing)
- `QuoteService` (Existing)
