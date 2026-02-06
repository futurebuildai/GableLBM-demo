# Architecture Specification

## 1. System Principles
- **Modular Monolith:** Single deployment binary, but internally strictly decoupled modules.
- **Zero-Trust Modules:** Modules never access another module's database tables directly.
- **Event-Driven:** Side effects (e.g., updating ledger after invoice posting) occur asynchronously via events.
- **Interface-Driven Interop:** Internal services use Go interfaces to allow for "Mock" or "Legacy Service" implementations (essential for migrations).
- **Federated Governance:** The platform supports a distributed contribution model, where industry partners (Co-ops) can propose core changes via a dedicated AI-mediated portal.

## 2. Technology Stack
- **Backend:** Go (Golang) 1.24+
- **Database:** PostgreSQL 16+
    - Extensions: `pgvector` (AI), `postgis` (Geospatial/Delivery).
- **Messaging:** NATS JetStream (embedded or external) for Event Bus.
- **Frontend:**
    - **Core:** React 19 + TypeScript + Vite.
    - **UI Library:** Shadcn/UI (Tailwind).
    - **State:** TanStack Query + Zustand.

## 3. Core Data Models (Schema Concepts)

### 3.1. Inventory Module (`pkg/inventory`)
The distinct complexity of LBM inventory.

```go
type Product struct {
    ID          uuid.UUID
    SKU         string
    Description string
    UOMFamilyID uuid.UUID // Links to "Dimensional Lumber" family
    
    // Tracking
    IsSerialized bool
    IsLotTracked bool
    
    // Attributes (AI Embedding Source)
    Attributes  map[string]interface{} 
}

type UOMConversion struct {
    FromUOM string // e.g. "MBF" (Thousand Board Feet)
    ToUOM   string // e.g. "PC" (Piece)
    Factor  decimal.Decimal
}
```

### 3.2. Sales Module (`pkg/sales`)
Optimized for complex pricing hierarchies.

```go
type Quote struct {
    ID         uuid.UUID
    ProjectID  uuid.UUID // Optional link to Project
    CustomerID uuid.UUID
    ExpiresAt  time.Time
    Status     QuoteStatus // Draft, Sent, Accepted, Expired
    
    // Versioning for "Bid Management"
    Version    int 
    RootQuoteID uuid.UUID
}

type PricingRule struct {
    Priority   int
    TargetType string // "Customer", "CustomerGroup", "Job"
    TargetID   uuid.UUID
    AdjustmentType string // "Markup", "Markdown", "FixedPrice"
    Value      decimal.Decimal
}
```

## 4. Inter-Module Communication
Modules must communicate without tight coupling.

### 4.1. Synchronous (Reads)
Direct Go Interface calls.
*Example:* Sales needs to check stock.
`InventoryService.GetAvailability(sku, location)` -> returns strict struct.

### 4.2. Asynchronous (Writes / Side Effects)
NATS JetStream Subjects.
*Example:* Order is Confirmed.
1. Sales publishes `sales.order.confirmed`.
2. Inventory subscribes -> Reserves Stock.
3. Logistics subscribes -> Creates "Pick Ticket".
4. Billing subscribes -> Checks Credit Limit (if not pre-checked).

## 5. API Strategy
- **Style:** RESTful JSON.
- **Definition:** OpenAPI 3.0 (Auto-generated from Go comments/types).
- **Auth:** OAuth2 / OIDC (Keycloak integration ready).

## 6. The Partner Portal & AI Governance Layer
To align with industry co-ops, GableLBM includes a specialized governance architecture.

- **Partner Portal:** A separate web interface (React) for co-op administrators to submit requirements.
- **AI Governance Engine:** 
    - **Parser:** Converts Natural Language requests into RFC-style technical specifications.
    - **Impact Analyzer:** Uses LLMs to evaluate how a requested change affects core modules (Inventory, Sales, etc.).
    - **Backlog Orchestrator:** Queues validated and de-duplicated requests into an open-source development pipeline.
- **Federated Catalog Service:** A multi-tenant sync layer that allows co-ops to push "Master SKU Data" to all member dealer instances simultaneously.

## 6. Legacy Interop & Migration Strategy
To facilitate open-source adoption, GableLBM is designed with "Import-First" capabilities.

- **Adaptor Layer:** Every core module (Inventory, Sales, Finance) includes an `adaptors/` directory. This contains mappers for:
    - Epicor BisTrack (REST JSON)
    - ECI Spruce (SOAP XML)
    - DMSi Agility (REST JSON)
- **Sync Engine:** A dedicated `pkg/sync` module handles bi-directional data flow between GableLBM and legacy ERPs during a phase-in period.
- **Schema Mapping:** We maintain a semantic mapping layer that translates legacy terms (e.g., `TallyString`) into our core `InventoryMove` models.
