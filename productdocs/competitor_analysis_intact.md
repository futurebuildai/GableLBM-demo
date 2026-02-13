# Competitor Analysis: Intact GenetiQ (Deep Dive)

## Overview
Intact GenetiQ (formerly Intact iQ) is a **cloud-native, browser-based ERP** built by Intact Software (32+ years in the industry). Marketed in the US specifically for Lumber & Building Materials dealers, it positions itself as a "future-proof" alternative to legacy ERPs. Deployment is **SaaS-only** with three subscription tiers (pricing requires a discovery call).

---

## Key Feature Sets

### 1. Platform Architecture
- **Cloud-Native**: Browser-based UI—no desktop client, no Citrix, no VPN.
- **API-First**: Designed for integrations with 3rd parties (Creditsafe, Aphix, NBG buying groups).
- **AI-Ready**: Built on a modern tech stack prepared for AI integration.
- **Personalization**: Fully customizable dashboards and user interfaces per-user.
- **Smart Alerts**: Proactive notifications that drive timely business actions.
- **Embedded Guidance**: Built-in user assistance for onboarding and daily tasks.

### 2. Building Material Supply Specific
- **Millwork Management**: Work orders for stock replenishment + "make-to-order" items.
- **Lumber Pack Tracking**: Flexible units of measure, track/trace tagged lumber packs.
- **Dispatch & Delivery**: Visual drag-and-drop scheduling, automated customer tracking links.
- **Proof of Delivery (POD) App**: Dedicated mobile POD app for drivers.
- **Trade Counter Focus**: POS optimized for high-volume trade environments (speed + margin control).

### 3. Financials & Accounting
- **Core Ledgers**: Fully integrated General, Sales, and Purchase ledgers.
- **Cash & Bank Management**: Comprehensive bank/cash management including VAT handling.
- **Fixed Asset Register**: Track and manage company assets within the ERP.
- **Credit Control**: Dedicated credit control workflows and credit request management.
- **Real-time Dashboards**: Instant visibility into margins and financial performance.

### 4. Stock Management & Warehousing
- **Inventory Optimization**: Lead time tracking, reorder points, Economic Order Quantity (EOQ).
- **Multi-Location Tracking**: Stock across multiple locations and bins in real-time.
- **Mobile Warehouse App** (Native):
  - Barcode scanning for picking, packing, receiving.
  - Stock transfers and adjustments on the fly.
  - **Offline data capture** for poor connectivity areas (critical for rural yards).
- **Integrated Procurement**: Purchase planning aligned with real-time demand.

### 5. Point of Sale & Trade Counter
- **"Haggle" Feature**: Staff can negotiate prices with customers while seeing full margin visibility—unique in the market.
- **Transaction Flexibility**: Sales, returns, and delivery requests from a single interface.
- **Margin Protection**: Configurable max-discount controls and price override rules.
- **Lost Sales Analysis**: Tracks and analyzes lost sales to improve inventory/pricing strategy.

### 6. BI & Reporting
- **GenetiQ Analyst**: Inbuilt pivot-table style data analysis tool.
- **Power BI Integration**: Direct data pipes for advanced visualization.

### 7. Nexus Workflow Engine
- **"Your Rules, Your Way"**: Highly customizable workflow automation.
- **Manage by Exception**: Automate routine tasks, alert on deviations.
- **Extensibility**: Custom fields and rules without core code changes.

---

## Known Integrations
| Integration | Purpose |
|:---|:---|
| Creditsafe | Automated credit checks on customer accounts |
| Aphix | E-commerce platform connector |
| NBG | Buying group direct integration |
| Power BI | Advanced analytics and dashboards |

---

## Strategic Observations

### Strengths (What They Do Well)
1. **True Cloud-Native**: The only LBM ERP that doesn't require Citrix/VPN for full functionality.
2. **"Haggle" POS**: Innovative margin-aware price negotiation at the counter—no competitor has this.
3. **Offline Mobile**: Their warehouse app works offline, solving the rural connectivity problem.
4. **Lost Sales Analysis**: Proactive analytics tool that others ignore.
5. **Workflow Engine (Nexus)**: Deep customization without vendor code changes.

### Weaknesses (Where Gable Wins)
1. **Proprietary "Black Box"**: Customization is through *their* proprietary tools, not open code.
2. **Pricing Opacity**: Three tiers but no public pricing—suggests high cost.
3. **US Market Entry**: Intact is Irish/UK-first; their US presence is growing but not well-established.
4. **No Open API Docs**: API-first by claim, but no public OpenAPI/Swagger spec.
5. **No AI-Native Features**: "AI-Ready" but no shipping AI features (no vision, no LLM counter).

### Gable Differentiation vs. GenetiQ
| Dimension | GenetiQ | GableERP |
|:---|:---|:---|
| Source | Proprietary (closed) | Open Source (auditable) |
| Customization | Nexus workflow engine (vendor tool) | Direct code access + RFC governance |
| AI | "AI-Ready" (future promise) | AI-Native (VelocityAI, Smart Counter) |
| Pricing | Opaque SaaS tiers | Transparent open core |
| API | API-first (gated docs) | Open API-first (public Swagger) |
| POS Innovation | "Haggle" margin tool | AI-powered material list parsing |
| US Market | Growing entrant | US-built, US-focused |

---

## Competitive Risk Assessment
**Risk Level: MEDIUM-HIGH**

GenetiQ is the most architecturally aligned competitor to GableERP. They share the cloud-native, browser-first philosophy. However, their proprietary model and opaque pricing create the exact "vendor lock-in" that Gable's Sovereignty thesis disrupts. The biggest risk is if Intact accelerates their US go-to-market with aggressive pricing before Gable reaches feature parity on financials and millwork.

**Recommended Counter-Strategy**: Emphasize open-source trust, public API documentation, and AI-native capabilities in head-to-head evaluations. Position the "Haggle" feature as a symptom of the old way (manual negotiation) vs. Gable's AI-suggested pricing.
