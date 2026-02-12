# Sprint 17: Dynamic Pricing Engine (Escalator Logic)

**Goal**: Implement "Escalator Pricing" to handle market volatility and long-term project quoting.

**Duration**: 1 Week
**Focus**: Pricing Logic, Market Data, Backend API.

## 1. Context
Competitor research (DMSi Agility) reveals a critical need for automated price escalators. Currently, GableLBM pricing is static or manual. This sprint introduces the "Pricing Engine" that can automatically adjust quotes based on scheduled increases or external market indices.

## 2. Objectives

### 2.1 Backend (Pricing Engine)
- [ ] **Market Index Service**:
    - [ ] Create `MarketIndex` model and repository.
    - [ ] Implement mock integration for "Lumber Index" (e.g., Random Lengths).
- [ ] **Escalator Logic**:
    - [ ] Implement `PriceEscalator` service to calculate future pricing based on % increase or index delta.
    - [ ] Support "Expiration Dates" on individual quote lines.
- [ ] **API Extensions**:
    - [ ] `POST /api/v1/pricing/calculate-escalation`
    - [ ] `GET /api/v1/market-indices`

### 2.2 Frontend (Sales Order Entry)
- [ ] **Quote Review Enhancements**:
    - [ ] Add "Escalator" toggle to quote lines.
    - [ ] Visual indicator for "Index-Linked" pricing.
    - [ ] Display "Future Realized Price" based on escalator dates.

## 3. Success Criteria
- System can calculate future price of a 2x4x8 based on a 5% monthly escalator.
- Market index data can be manually updated or fetched via mock service.
- Sales reps receive a warning when a quote line's pricing is "Stale" vs current index.
