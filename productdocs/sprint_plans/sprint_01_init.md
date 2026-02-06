# Sprint 01: Initialization ("Hello World")

**Goal**: Establish the technical foundation. Get the backend, frontend, and database talking to each other.
**Exit Criteria**: A developer can run `docker-compose up` and see the App Shell with a working database connection.

## 1. Backend (`/backend`)
*   **Action**: Initialize Go Module `github.com/gablelbm/gable`.
*   **Architecture**: Set up `cmd/server`, `internal/config`, `pkg/database`.
*   **Framework**: Standard `net/http` (Go 1.22+ mux) or `Chi`.
*   **Endpoints**:
    *   `GET /health`: Returns `{"status": "ok", "db": "connected"}`.

## 2. Frontend (`/app`)
*   **Action**: Create new Vite Project (React + TS).
*   **UI Lib**: Install `shadcn/ui`, `tailwindcss`, `lucide-react`.
*   **Design System**: Port `index.css` variables from `productdocs/design_system.md`.
*   **Page**: Create a simple Dashboard Shell (Sidebar + Header + Empty Content Area).

## 3. Infrastructure
*   **Docker**: Create `docker-compose.yml`.
    *   **Postgres 16**: With persistent volume.
    *   **NATS JetStream**: For future event bus.
*   **Seed**: Basic SQL init script to create a user or check connection.
