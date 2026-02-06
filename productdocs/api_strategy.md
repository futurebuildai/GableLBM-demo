# API Strategy Specification

## 1. Design Philosophy
- **Contract-First:** OpenAPI 3.1 spec is the source of truth. Go server code and TypeScript client code are generated from it.
- **Versioned:** All endpoints live under `/api/v1/`. Breaking changes require a `/v2/` migration path.
- **RESTful with Pragmatism:** CRUD resources are REST. Complex operations (e.g., "Convert Quote to Order") use action endpoints (`POST /quotes/{id}/convert`).

## 2. Authentication & Authorization
- **Protocol:** OAuth2 / OpenID Connect (OIDC).
- **Provider:** Self-hosted Keycloak (or compatible: Auth0, Zitadel).
- **Tokens:** JWT Access Tokens with short TTL (15 min). Refresh tokens for long sessions.
- **RBAC:** Role-Based Access Control. Roles are mapped to API scopes.

| Role              | Scopes                                      |
|-------------------|---------------------------------------------|
| `counter_rep`     | `sales:write`, `inventory:read`             |
| `yard_foreman`    | `inventory:write`, `logistics:write`        |
| `controller`      | `finance:*`, `reports:*`                    |
| `admin`           | `*` (all scopes)                            |
| `customer_portal` | `orders:read:self`, `invoices:read:self`    |

## 3. Core API Modules

### 3.1. Inventory API (`/api/v1/inventory`)
| Endpoint                          | Method | Description                               |
|-----------------------------------|--------|-------------------------------------------|
| `/products`                       | GET    | List products (paginated, filterable).    |
| `/products/{sku}`                 | GET    | Get product details by SKU.               |
| `/products/{sku}/availability`    | GET    | Real-time stock by location.              |
| `/stock-moves`                    | POST   | Create an inventory move (receive, ship). |
| `/stock-moves/{id}`               | GET    | Get move status.                          |

### 3.2. Sales API (`/api/v1/sales`)
| Endpoint                          | Method | Description                               |
|-----------------------------------|--------|-------------------------------------------|
| `/quotes`                         | POST   | Create a new quote.                       |
| `/quotes/{id}`                    | GET    | Get quote details.                        |
| `/quotes/{id}/convert`            | POST   | Convert quote to an order.                |
| `/orders`                         | GET    | List orders (by customer, job, date).     |
| `/orders/{id}`                    | GET    | Get order details.                        |
| `/orders/{id}/lines`              | POST   | Add line items to an order.               |

### 3.3. Finance API (`/api/v1/finance`)
| Endpoint                          | Method | Description                               |
|-----------------------------------|--------|-------------------------------------------|
| `/invoices`                       | GET    | List invoices (by customer, status).      |
| `/invoices/{id}`                  | GET    | Get invoice details.                      |
| `/invoices/{id}/pdf`              | GET    | Download invoice PDF.                     |
| `/payments`                       | POST   | Record a payment against an invoice.      |
| `/customers/{id}/ledger`          | GET    | Get customer AR ledger.                   |

### 3.4. Logistics API (`/api/v1/logistics`)
| Endpoint                          | Method | Description                               |
|-----------------------------------|--------|-------------------------------------------|
| `/pick-tickets`                   | GET    | List pending pick tickets.                |
| `/pick-tickets/{id}`              | PATCH  | Update pick status (picked, loaded).      |
| `/deliveries`                     | POST   | Schedule a new delivery run.              |
| `/deliveries/{id}/route`          | GET    | Get optimized route.                      |
| `/deliveries/{id}/pods`           | POST   | Upload proof of delivery (photo, sig).    |

## 4. Integration Webhooks
For external systems, we expose outbound webhooks on key events.
- **Event Format:** CloudEvents v1.0 JSON.
- **Delivery:** HTTP POST to registered endpoint with HMAC signature.

| Event Type                  | Trigger                                   |
|-----------------------------|-------------------------------------------|
| `order.created`             | A new order is placed.                    |
| `order.shipped`             | All items on an order have been shipped.  |
| `invoice.posted`            | An invoice is finalized.                  |
| `payment.received`          | A payment is applied.                     |
| `inventory.adjusted`        | A manual inventory adjustment occurs.     |

## 5. Legacy Adaptor Endpoints
To support bi-directional sync during migrations from legacy systems.
- `/api/v1/sync/bistrack` - Import/Export payloads in BisTrack JSON format.
- `/api/v1/sync/spruce` - Import payloads from Spruce SOAP (converted to JSON).
- `/api/v1/sync/dmsi` - Import payloads in DMSi JSON format.
