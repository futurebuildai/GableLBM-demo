# Sprint 14 Task List

## Phase 1: Backend Core & Schema
- [x] Create `rfcs` table migration <!-- id: 1 -->
- [x] Define `RFC` model and Repository <!-- id: 2 -->
- [x] Setup `governance` package and Service stub <!-- id: 3 -->
- [x] Implement `CreateRFC` and `GetRFCs` core logic <!-- id: 4 -->

## Phase 2: AI Generation Logic
- [x] Define `AIProvider` interface <!-- id: 5 -->
- [x] Implement `TemplateAIProvider` (Mock/Template based) <!-- id: 6 -->
- [x] Wire `GenerateRFC` endpoint using `AIProvider` <!-- id: 7 -->

## Phase 3: Frontend Governance Portal
- [x] Create `GovernanceLayout` (Sidebar/Header) <!-- id: 8 -->
- [x] Implement `RFCDashboard` (List View) <!-- id: 9 -->
- [x] Implement `NewRFCWizard` (Form Input) <!-- id: 10 -->
- [x] Implement `RFCDetailView` (Markdown render/edit) <!-- id: 11 -->

## Phase 4: Integration & Verification
- [x] Connect Frontend to Backend `governance` API <!-- id: 12 -->
- [x] Verify RFC Generation flow (End-to-End) <!-- id: 13 -->
- [x] L8 Audit Check (Security/Architecture) <!-- id: 14 -->
