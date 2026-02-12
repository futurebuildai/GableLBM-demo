# Sprint 16: Sovereign Dealer Portal (B2B)

**Goal**: Launch an open-source, white-labeled portal for contractors — the highest-impact "Sovereignty" proof point.

**Duration**: 1 Week
**Focus**: B2B Portal, Authentication, Customer Self-Service.

## 1. Context
Both DMSi Agility and Intact iQ offer contractor portals, but they are SaaS-locked. GableLBM's portal will be fully owned by the dealer — code, data, and hosting. This is the single most tangible feature to show co-ops when pitching "Sovereignty."

## 2. Objectives

### 2.1 Portal Core (New App/Service)
- [ ] **Auth & Security**:
    - [ ] Implement `CustomerUser` auth (separate from internal yard staff).
    - [ ] Token-based access to customer-specific documents.
- [ ] **Contractor Dashboard**:
    - [ ] AR Balance & Statement View.
    - [ ] Project/Job listing.
    - [ ] Document Center (Invoices, POD photos).

### 2.2 Self-Service Ordering
- [ ] **Quick Re-order**:
    - [ ] "Buy Again" logic based on historical project lists.
    - [ ] Real-time price fetch from dealer core.

### 2.3 Branding & Sovereignty
- [ ] **White-Labeling System**:
    - [ ] Secondary configuration for Logo/Primary Color (for the portal only).
    - [ ] Dealer owns the deployment — no platform lock-in.

## 3. Success Criteria
- Contractor can log in and see their current balance.
- Contractor can download a PDF of a past invoice.
- Contractor can see a POD photo of their most recent delivery.
- Portal is visually distinct from the internal admin (dealer-branded).
