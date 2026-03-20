# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

GableLBM is an Enterprise Operating System for the Lumber & Building Materials (LBM) industry. It's a full-stack ERP with modules for inventory, orders, quotes, invoices, purchasing, logistics, accounting, POS, a B2B dealer portal, and mobile driver/yard apps.

## Tech Stack

- **Frontend**: React 19 + TypeScript + Vite 7 + Tailwind CSS 3
- **Backend**: Go 1.24 with Chi router (`github.com/go-chi/chi/v5`)
- **Database**: PostgreSQL 16 (pgx driver, pgxpool)
- **Message Queue**: NATS with JetStream
- **AI Integrations**: Claude (Anthropic), Gemini (Google), Stability AI
- **Deployment**: Railway (Docker-based)

## Development Commands

### Frontend (`app/`)
```bash
cd app && npm run dev       # Start Vite dev server (proxies API to localhost:9091)
cd app && npm run build     # TypeScript check + Vite production build
cd app && npm run lint      # ESLint
```

### Backend (`backend/`)
```bash
cd backend && go run ./cmd/server/main.go    # Start API server (port 8080/9091)
cd backend && go run ./cmd/migrate/main.go   # Run database migrations
cd backend && go run ./cmd/seed/main.go      # Seed demo data
cd backend && go test ./...                  # Run all tests
cd backend && go test ./internal/order/...   # Run tests for a single package
```

### Infrastructure (root Makefile)
```bash
make up         # Start PostgreSQL + NATS via docker-compose
make down       # Stop containers
make logs       # Tail container logs
make pg-shell   # psql into gable_db (user: gable_user, db: gable_db, port: 5434)
```

## Architecture

### Frontend (`app/src/`)

**Routing** (defined in `App.tsx`):
- `/` — Demo landing page
- `/pos` — Retail POS terminal
- `/erp/*` — Main desktop ERP (wrapped in `AppShell` layout with sidebar)
- `/driver/*` — Mobile driver app (DriverLayout)
- `/yard/*` — Mobile yard app (YardLayout)
- `/portal/*` — B2B dealer portal (PortalLayout, no auth for demo)

**Service layer** (`src/services/`): ~33 service files. Each exports an object with async methods that call `fetch()` against `VITE_API_URL`. Pattern:
```ts
const API_BASE_URL = `${import.meta.env.VITE_API_URL || ''}/api`;
export const OrderService = {
    async listOrders(): Promise<Order[]> { /* fetch(...) */ },
};
```

**UI patterns**: Tailwind + Shadcn/UI-style components in `src/components/ui/`. Uses `cn()` from `src/lib/utils.ts` (clsx + tailwind-merge). Icons from `lucide-react`. Animations via `framer-motion`.

**State**: React Context for cross-cutting concerns (Toast). No Redux/Zustand — page-level state with hooks.

### Backend (`backend/`)

**Modular domain architecture** — each domain is a package under `internal/` with its own repository, service, and handler:
```
internal/
  order/       # repository.go, service.go, handler.go
  invoice/
  inventory/
  product/
  pricing/     # pricing, rebates, escalators
  gl/          # general ledger
  portal/      # B2B dealer portal
  pos/         # point of sale
  delivery/    # logistics & routes
  ...40+ modules
```

**Server bootstrap** (`cmd/server/main.go`): Initializes all modules with dependency injection, wires handlers to Chi router. Auth via JWT/JWKS (falls back to demo mode if `JWKS_URL` not set).

**Database**: pgxpool connection pool. `pkg/database/postgres.go` provides `RunInTx()` for transactions and `GetExecutor()` that returns the active tx or pool. Migrations are numbered SQL files in `migrations/` (001–043+).

**Key integrations** (all optional, config-gated):
- Avalara AvaTax (sales tax)
- Run Payments (card processing)
- Twilio (SMS)
- Google Maps (route optimization)
- AI key store in DB (configurable via Tech Admin panel, env fallback)

### Vite Dev Proxy

The Vite dev server proxies these paths to `http://127.0.0.1:9091`: `/api`, `/products`, `/customers`, `/vendors`, `/orders`, `/quotes`, `/invoices`, `/locations`, `/health`, `/activities`, `/contacts`, `/documents`, `/gl`, `/parsing`, `/price_levels`, `/pricing`, `/purchase-orders`, `/sales-team`, `/uploads`.

## Design System — "Industrial Dark"

- **Primary**: Gable Green `#00FFA3`
- **Background**: Deep Space `#0A0B10`
- **Surfaces**: Slate Steel `#161821`, Surface-2 `#1E2029`, Surface-3 `#252836`
- **Accent**: Blueprint Blue `#38BDF8`, Safety Red `#F43F5E`
- **Typography**: Outfit (sans), JetBrains Mono (mono)
- **Elevation**: MD3 shadow system (elevation-1/2/3)
- CSS variables defined in `app/src/index.css`, Tailwind tokens in `app/tailwind.config.js`
- All UI must use these tokens. No unstyled HTML.

## Development Standards

From CONTRIBUTING.md — the "Antigravity Standard":
- **Headless**: Logic in hooks/services, UI is a render layer
- **Strong types**: No `any` in TypeScript
- **Speed**: Target < 100ms interactions
- Sprint plans live in `productdocs/sprint_plans/`
- Agent workflow protocol: `.agent/workflows/sprint_execution.md`
- Code must pass the L8 Production Readiness Gate before merging (`productdocs/process/production_readiness_gate.md`)

## Key Reference Documents

- `productdocs/master_prd.md` — Full product requirements
- `productdocs/design_system.md` — Design tokens and patterns
- `productdocs/architecture_spec.md` — System architecture spec
- `productdocs/sprint_plans/` — Sprint-by-sprint plans
