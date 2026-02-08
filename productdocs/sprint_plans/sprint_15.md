# Sprint 15: Executive Analytics Dashboard

**Goal**: Build a comprehensive, real-time Executive Dashboard that replaces placeholder data with live analytics, providing "Owner Bob" with instant visibility into business KPIs.

**Duration**: 1 Week
**Focus**: Analytics, Data Visualization, Real-Time Updates.

## 1. Context
The main Dashboard currently displays hardcoded placeholder values. With the reporting backend already in place (`/api/reports/daily-till`, `/api/reports/sales-summary`), this sprint connects the frontend to real data and extends the analytics capabilities with inventory alerts, order metrics, and trend visualization.

## 2. Objectives

### 2.1 Backend (Analytics API Extensions)
- [ ] **Dashboard Aggregate Endpoint**:
    - [ ] `GET /api/v1/dashboard/summary` (Single endpoint returning all KPIs).
    - [ ] Includes: Revenue, Order Counts, Pending Dispatch, AR Summary.
- [ ] **Inventory Alerts Endpoint**:
    - [ ] `GET /api/v1/dashboard/inventory-alerts` (Low stock, out of stock items).
- [ ] **Top Customers Endpoint**:
    - [ ] `GET /api/v1/dashboard/top-customers` (Top 5 by revenue in period).
- [ ] **Order Activity Endpoint**:
    - [ ] `GET /api/v1/dashboard/order-activity` (Recent orders, status distribution).

### 2.2 Frontend (Dashboard Components)
- [ ] **KPI Cards**: Replace hardcoded values with live API data.
    - Revenue (with trend arrow)
    - Active Orders
    - Pending Dispatch
    - Outstanding AR
- [ ] **Charts**:
    - [ ] Revenue Trend (7-day line/area chart).
    - [ ] Order Status Distribution (Pie/Donut chart).
- [ ] **Tables**:
    - [ ] Top Customers widget (sortable mini-table).
    - [ ] Inventory Alerts widget (low stock warnings).
    - [ ] Recent Orders activity feed.
- [ ] **Real-Time Refresh**:
    - [ ] Auto-refresh every 60 seconds with visual indicator.
    - [ ] Manual refresh button.

### 2.3 Services & Types
- [ ] **DashboardService.ts**: Frontend service for new endpoints.
- [ ] **Types**: `DashboardSummary`, `InventoryAlert`, `TopCustomer`, `OrderActivity`.

### 2.4 Design System Adherence
- [ ] **Industrial Dark** theme applied to all new components.
- [ ] **JetBrains Mono** for all numerical data.
- [ ] **Micro-animations** on hover states and data refresh.

## 3. Technical Considerations
- **Caching**: Extend backend caching to new dashboard endpoints (60s TTL).
- **Performance**: Single aggregate endpoint to reduce frontend API calls.
- **Charting**: Use Recharts (already compatible with React) or Chart.js.
- **Error States**: Skeleton loaders during data fetch, graceful error messages.

## 4. Success Criteria
- Dashboard displays real-time data from the database.
- All 4 KPI cards show accurate, live values.
- Revenue trend chart visualizes last 7 days of data.
- Inventory alerts prominently display low-stock items.
- Page loads in < 500ms with progressive loading states.
- UI meets "Industrial Dark" aesthetic standards.
