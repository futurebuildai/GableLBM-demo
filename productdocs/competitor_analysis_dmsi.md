# Competitor Analysis: DMSi Agility (Deep Dive)

## Overview
DMSi is an **independent, family-owned** software company based in Omaha, Nebraska, operating since **1976** (49 years). Their flagship product, **Agility ERP**, is purpose-built for the building materials supply chain — from pro dealers and distributors to millwork manufacturers and hardwood processors. They have the **highest customer satisfaction in the LBM ERP market (95%, 4.8/5.0 on FeaturedCustomers)** with 3,067+ reference ratings.

> **Key Customers**: Boise Cascade, SRS Distribution, Western Materials, Bridgewell Resources, US Lumber Group, Schutte Lumber Company.

---

## Company Profile

| Attribute | Details |
|:---|:---|
| **Founded** | 1976 (Omaha, NE) |
| **Ownership** | Independent, family-owned (marketed as stability advantage) |
| **Focus** | Building materials supply chain exclusively |
| **Products** | Agility ERP, Frameworks ERP (retail), TallyExpress, Commerce Cloud, Log System |
| **Pricing** | Custom quotes; estimated $500–$1,000/mo starting range |
| **Deployment** | SaaS + On-Premise |
| **Satisfaction** | 95% (highest in class), 4.8/5.0 FeaturedCustomers |

---

## Product Suite

### 1. Agility ERP (Flagship — Distribution/Pro Dealer)
The core ERP for mid-to-large lumber dealers, distributors, and wholesalers. Handles everything from front-counter sales to back-office accounting.

### 2. Frameworks ERP (Retail/Home Centers)
A **web-based** ERP specifically for retail operations and home centers. Lighter weight than Agility, focused on retail POS and consumer-facing workflows. *This is an important distinction — DMSi offers TWO ERPs, segmented by business type.*

### 3. TallyExpress (AI Lumber Tally)
*See detailed section below.*

### 4. Commerce Cloud (B2B E-Commerce)
*See detailed section below.*

### 5. DMSi Log System (Hardwood)
Specialized software for hardwood log scaling, inventory, and vendor contract settlements. *Niche but shows depth of industry commitment.*

---

## Core Modules (Agility ERP)

### Inventory & Warehouse Management
- **Dimensional Tracking**: Species, Grade, Treatment tracked under single item codes
- **Multi-UOM**: PC, UNIT, BF, MBF, LF, M3 — best-in-class for lumber-specific units
- **Item Builder**: Maintains consistent item codes, reduces duplicates across complex catalogs
- **Real-time Multi-Branch Sync**: Live inventory across all branches
- **Storage Location Tracking**: Stock levels by specific bin/yard location
- **Mobile Warehouse Tools**:
  - Cycle counting with barcode scanning
  - Paperless picking (eliminates pick tickets)
  - Receiving with PO verification and photo capture
  - Stock movement / inter-branch transfers
- **Automated Replenishment**: Reorder calculations with customizable thresholds

### Sales & Order Management (OrderPad)
- **Quick Sales / POS**: Single-screen transaction interface for counter speed
- **Special Orders**: Non-stock items handled within the same order flow (no separate module)
- **Job-Specific Pricing**: Contract-based, volume-based, and customer-specific pricing
- **Escalator Pricing**: Automated scheduled price increases on quotes for long-term construction projects (handles commodity volatility)
- **Variable Pricing**: Set by customer, customer group, or individual job
- **Installed Sales Project Management**: Track project-based sales with labor and materials
- **Long-Lead Order Management**: Support for extended-timeline orders

### Mobile Sales App
- **Outside Sales Tool**: Reps check A/R balances, available stock, and pricing on-site
- **Quote/Order Creation**: Build quotes and submit orders from the field
- **Customer Visit Integration**: Ties into CRM for call logging

### CRM (Customer Relationship Management)
- **Industry-Specific Intelligence**: Connects accounts with sales history, A/R status, contact management
- **Activity Tracking**: Log calls, meetings, schedule follow-ups
- **Team/Division Access**: Shared intelligence across divisions for responsive service
- **Pipeline Management**: Track opportunities and sales pipeline

### Purchasing & Procurement
- **Demand Forecasting**: Predictive purchasing suggestions based on historical data
- **Suggested Purchasing**: System-recommended reorders based on stock levels and demand
- **Centralized Purchasing**: Multi-location procurement management from single interface
- **EDI (Electronic Data Interchange)**: Electronic business transactions with suppliers
- **Pricing Visibility**: Full supplier cost transparency during purchasing decisions
- **Vendor Interface Management**: Direct connections to major suppliers

### Dispatch, Logistics & Delivery
- **Route Tracking**: Real-time fleet tracking and delivery status
- **Proof of Delivery (POD) App**: Mobile app with signature capture, confirmation photos, GPS
- **Automatic Customer Notifications**: Staged → Out for Delivery → Delivered status updates
- **Driver On-Site Adjustments**: Drivers can adjust invoice quantities on-site (solves "short-ship" billing)
- **DQT Fleet Integration**: Direct integration with DQT fleet tracking platform

### Financials & Accounting
- **Full Suite**: Integrated A/R, A/P, General Ledger
- **Automated Invoicing**: Billing flows directly from delivery/order completion
- **Payment Processing**: Integration with Worldpay (FIS) and Billtrust
- **F9 Reporting**: Excel-based financial reporting integration
- **Avalara Tax**: Automated sales tax calculation

### Production & Millwork
- **Work Orders**: Manage manufacturing processes, issue materials, track finished goods
- **Production Control**: Shop floor visibility and scheduling
- **Complex Configuration**: Rules-based configurator for doors, windows, molding assemblies
- **Bills of Materials**: Multi-level BOM handling for remanufactured and assembled products
- **Resource Scheduling**: Production line scheduling and capacity planning

---

## TallyExpress — AI Computer Vision Lumber Tally

> [!IMPORTANT]
> TallyExpress is **the most advanced AI-powered product in the entire LBM ERP market**. It directly competes with GableERP's planned AI vision capabilities and represents DMSi's strongest innovation moat.

### How It Works
1. User takes a photo of a lumber bundle with an Android phone (with a reference square visible)
2. AI computer vision detects and counts individual boards
3. Each board measured to **1/16 of an inch** accuracy
4. Instant tally report generated with widths, lengths, and bundle photo

### Key Specifications
| Metric | Value |
|:---|:---|
| **Speed** | 90 seconds vs. 15 minutes manual (10x faster) |
| **Accuracy** | Starts within 1%, improves to **99.8–99.9%** with use |
| **Board Measurement** | To 1/16 inch precision |
| **Platform** | Android only (phones and tablets) |
| **Languages** | 8 (English, German, Spanish, French, Italian, Romanian, Chinese, Russian) |
| **Intelligence** | Split ID (identifies splits as single board, flush boards as two) |
| **Export** | CSV, XLS, PDF; direct ERP integration |
| **Architecture** | Cloud-based; auto-uploads and syncs in real-time |

### TallyExpress Features
- **Green Light Interface**: Visual guide helps users line up the phone correctly
- **Multiple Length Handling**: Adjust for cutbacks within a bundle
- **Team Collaboration**: Multiple users per account; supervisor review/approve/archive workflow
- **Machine Learning**: Accuracy improves continuously with each use
- **Multilingual**: 8 languages with corresponding support

### competitive Implications for GableERP
- DMSi **already ships** AI computer vision — this is NOT roadmap
- GableERP's planned AI tally must match or exceed 99.8% accuracy
- The 90-second workflow is the benchmark to beat
- Consider: Can GableERP tally work on iOS too? (DMSi is Android-only)

---

## Agility Commerce Cloud — B2B E-Commerce

### Platform Overview
Purpose-built B2B e-commerce platform that integrates directly with Agility ERP. Not a generic Shopify-style bolt-on — it understands lumber-specific units, pricing, and workflows.

### Key Features
- **Industry-Specific UOM**: Supports PC, UNIT, BF, MBF, LF, M3 in the shopping experience
- **Real-time Pricing/Inventory**: Always synced with Agility ERP (no stale data)
- **B2B2B Retail Pricing**: Supports multi-tier pricing for distribution chains
- **Customer Self-Service 24/7**: Browse products, place orders, check backorders, manage invoices/payments
- **Cart Splitting**: Customers can split a cart across multiple ship-to addresses
- **Guest User Admin**: Customers manage their own team's portal access
- **Quick Reorder**: Frequently purchased products easily reordered
- **No-Code CMS**: Add pages, images, videos, and advertising without developers
- **Custom Landing Pages**: Targeted marketing content per customer segment
- **Mobile-Optimized**: Fully responsive design

### Strategic Implications
The Commerce Cloud is a **complete, integrated e-commerce solution** — not a third-party integration. This is a significant competitive moat that GableERP's Sprint 16 portal needs to match.

---

## Integration Ecosystem

| Partner | Purpose |
|:---|:---|
| **Avalara** | Automated sales tax calculation |
| **F9** | Excel-based financial reporting |
| **Billtrust** | Payment processing / AR automation |
| **Worldpay (FIS)** | Payment processing |
| **DQT** | Fleet tracking |
| **Trimble** | Maps / route optimization |
| **Saberis** | Order interface / EDI |
| **Phocas** | Data analytics / BI |
| **CyberScience** | Business intelligence |
| **Loftware** | Labeling solutions |
| **Barcodes** | Barcode hardware |
| **TierPoint** | Managed hosting |

---

## Technology & Architecture

| Attribute | Details |
|:---|:---|
| **Core Stack** | Proprietary / Legacy |
| **Desktop Access** | Citrix/RDP for full Agility functionality |
| **Cloud Products** | Commerce Cloud, Frameworks ERP, TallyExpress are cloud-native |
| **Mobile** | Native iOS & Android apps for warehouse, sales, delivery |
| **API** | Available but **private/gated**; actively modernizing (sawmill API initiative) |
| **API Documentation** | Not publicly available — requires partnership |
| **Deployment** | SaaS + On-Premise options |

### API Modernization Note
DMSi has published content about "Modernizing Sawmills with APIs," indicating they are actively investing in API capabilities. A university senior design project (2024-2025) developed a "Legendary Lumber Lookup Tool" as a proof-of-concept for Agility, suggesting continued innovation investment.

---

## Industry Segments Served

| Segment | Product | Depth |
|:---|:---|:---|
| **Pro Dealers / Lumberyards** | Agility ERP | ✅ Deep |
| **Distributors / Wholesalers** | Agility ERP | ✅ Deep |
| **Millwork Manufacturers** | Agility ERP + Production | ✅ Best-in-class |
| **Hardwood Processors** | Agility ERP + Log System | ✅ Specialized |
| **Truss / Pallet** | Agility ERP | ✅ Specialized |
| **Roofing / Siding** | Agility ERP | ✅ Specialized |
| **Drywall** | Agility ERP | ✅ Specialized (sqft/lf quoting) |
| **Retail / Home Centers** | Frameworks ERP | ✅ Dedicated product |

---

## User Reviews & Sentiment

### Strengths (What Users Love)
1. **"It just works for lumber"** — Industry-specific UOM, pricing, and workflows out of the box
2. **OrderPad** — Fastest order entry in the market; new salespeople adapt quickly
3. **Millwork depth** — Unmatched configurator for doors/windows/assemblies
4. **Customization** — High flexibility; adapts to specific business processes
5. **TallyExpress** — AI vision tally is a genuine wow feature
6. **POD app** — Signatures, photos, on-site adjustments solve real problems
7. **Commerce Cloud** — Purpose-built B2B e-commerce, not a bolt-on

### Weaknesses (User Complaints)
1. **Citrix/RDP requirement** — Full desktop requires remote access; dispatchers tethered to office
2. **Gated API** — Integration access is selective and expensive ($5K–$15K reported)
3. **Manual process gaps** — Some processes outside the "happy path" require manual tracking
4. **High cost of entry** — Custom pricing; enterprise-grade cost structure
5. **Training curve** — Remanufacturing/millwork features require significant training
6. **Legacy core feel** — Powerful but not "modern" in UI/UX

---

## Strategic Assessment

### Risk Level: **HIGH** 🔴

DMSi Agility is the **most functionally complete** and **highest-satisfaction** competitor in the LBM ERP market. Their combination of industry depth, TallyExpress AI innovation, and Commerce Cloud creates multiple competitive moats.

### Why Agility Is Hard to Displace
1. **95% satisfaction** — Users genuinely love the product
2. **TallyExpress** — Shipping AI product (not roadmap) with 99.8% accuracy
3. **Commerce Cloud** — Integrated B2B platform, not a 3rd party bolt-on
4. **Millwork depth** — Best configurator in the market
5. **Family ownership** — Marketed as stability; no PE exit risk

### Where GableERP Wins
| Dimension | DMSi Agility | GableERP |
|:---|:---|:---|
| **Access** | Citrix/RDP (tethered) | Browser-native (anywhere) |
| **API** | Private/Gated ($5K–$15K) | Open REST + Events (free) |
| **Millwork** | Expensive add-on module | Core "Assembly/Kitting" in open tier |
| **Pricing Logic** | Escalator pricing (strong) | Escalator + live market index integration |
| **AI Vision** | TallyExpress (Android-only, tally only) | VelocityAI (multi-modal: tally + material lists + blueprints) |
| **Source** | Proprietary | Open Source |
| **Governance** | Vendor-controlled roadmap | RFC community governance |
| **E-Commerce** | Commerce Cloud (strong) | Sovereign Dealer Portal (white-labeled, dealer-owned) |
| **Tally Platform** | Android only | Cross-platform (planned) |
| **Cost** | $500–$1K+/mo + integration fees | Open core + transparent pricing |

### Recommended Counter-Strategy
1. **Don't compete on millwork depth initially** — Compete on speed, access, and cost
2. **Position against Citrix** — "Agility's functionality, accessible from any browser"
3. **Match TallyExpress accuracy** — AI vision tally must hit 99%+ or don't ship it
4. **Extend beyond tally** — VelocityAI should parse material lists AND blueprints (TallyExpress only counts boards)
5. **Open API advantage** — Every dealer frustrated by DMSi's API fees is a Gable prospect
6. **iOS tally** — TallyExpress is Android-only; iOS support is an easy win
7. **Sovereign e-commerce** — Commerce Cloud locks dealers to DMSi; GableERP portal is dealer-owned
