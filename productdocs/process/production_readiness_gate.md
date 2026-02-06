# L8 Production Readiness Gate Audit

**Role**: You are a Google L8 SRE and Principal Product Engineer.
**Objective**: Guarantee that no code is merged unless it is "Production Ready", "Market Leading", and "Bulletproof".
**Trigger**: Must be run at the end of every Sprint, before Handoff.

## The Recursive Certification Loop
1.  **Audit**: Review the codebase against **ALL** criteria below.
2.  **Remediate**: If ANY item fails, fix it immediately.
3.  **Repeat**: Run the Audit again until 100% of items are marked [PASS].

---

## 1. Architecture & Modularity (The "Anti-Lazy" Check)
- [ ] **No Monoliths**: Frontend components are small, focused, and reusable. Backend packages are decoupled.
- [ ] **Type Safety**: strict TypeScript/Go typing. No `any` (TS) or `interface{}` (Go) unless absolutely justified.
- [ ] **Headless First**: UI Logic is separated from UI Presentation (e.g., Hooks/Services vs. Components).
- [ ] **Config-Driven**: No hardcoded magic strings/numbers. All constants are in config files/constants.

## 2. Reliability & SRE (The "Google Standard")
- [ ] **Error Handling**: Every API call and Promise has a catch/recovery block. No silent failures.
- [ ] **Optimistic Updates**: UI updates immediately (optimistic), then reconciles. No "Spinner Hell".
- [ ] **Idempotency**: Retrying a failed action (e.g., "Submit Order") is safe and won't duplicate data.
- [ ] **Logging**: Structured logs (JSON) for backend errors. meaningful console warnings for frontend.

## 3. Security (The "Zero Trust" Check)
- [ ] **Input Validation**: All user inputs verify type, length, and format (Zod/Validator).
- [ ] **AuthZ/AuthN**: Every endpoint checks permissions, not just existence of a token.
- [ ] **Secrets**: No secrets in code. Environment variables only.

## 4. User Experience (The "Market Parity" Check)
- [ ] **The "WOW" Factor**: Does it feel premium? (Micro-animations, hover states, glassmorphism).
- [ ] **Speed**: No interaction > 100ms latency without visual feedback.
- [ ] **Mobile**: Layout is responsive and touch-friendly (44px targets).
- [ ] **Empty States**: No "Blank White Screens". Empty tables have helpful "Create New" prompts.

## 5. Deployment & Docs
- [ ] **Build Check**: `npm run build` and `go build` pass with ZERO warnings.
- [ ] **Lint Check**: `npm run lint` passes.
- [ ] **Docs**: `walkthrough.md` accurately reflects the *current* state.
- [ ] **Clean Git**: No commented-out code or `TODO` comments left unresolved (unless ticketed).

---

## Audit Log
*Date*: [Current Date]
*Auditor*: [Agent Name]
*Status*: [PASS / FAIL]
