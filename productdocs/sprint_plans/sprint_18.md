# Sprint 18: VelocityAI ERP Integration

**Goal**: Embed the standalone VelocityAI parsing engine directly into the GableLBM Sales Order Entry flow.

**Duration**: 1 Week
**Focus**: AI Integration, Document Parsing, Sales Workflow.

## 1. Context
VelocityAI (formerly "Quick-Quote") can parse handwritten material lists into structured line items with confidence scores. Currently it operates as a standalone tool. Integrating it into the ERP's core order flow creates an AI-native sales experience that neither DMSi nor Intact can match. This is the single biggest differentiator in our stack.

## 2. Objectives

### 2.1 Backend Integration
- [ ] **VelocityAI Service Bridge**:
    - [ ] Create a `velocity` package in the Go backend that calls the VelocityAI parsing API.
    - [ ] Map VelocityAI output (parsed items) to GableLBM `OrderLine` schema.
- [ ] **Confidence Handling**:
    - [ ] Items above 90% confidence auto-populate on the quote.
    - [ ] Items below 90% are flagged for manual SKU selection with alternatives.
- [ ] **Special Order Fallback**:
    - [ ] Unmatched items generate a "Special Order" line with the raw parsed text.

### 2.2 Frontend (Sales Order Entry)
- [ ] **"Upload Material List" Button**:
    - [ ] Add upload trigger directly in the Quote/Order creation flow.
    - [ ] Progress indicator during AI parsing.
- [ ] **Results Review Panel**:
    - [ ] Side-by-side view: Original image vs. Parsed line items.
    - [ ] Inline confidence badges and "Swap SKU" dropdowns.
- [ ] **One-Click Accept**:
    - [ ] "Accept All" button to push parsed items into the active quote.

## 3. Success Criteria
- Sales rep can photograph a handwritten material list and have it parsed into a quote in < 30 seconds.
- Items with low confidence are clearly flagged and easy to correct.
- End-to-end flow: Photo → Parse → Review → Quote is functional in the browser.
