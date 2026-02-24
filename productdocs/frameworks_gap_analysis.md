# DMSi Frameworks ERP vs GableERP — Feature Gap Analysis

> **Purpose**: Comprehensive feature-by-feature comparison of DMSi's **Frameworks ERP** (their web-based, retail-focused LBM product) against GableERP's current implementation. Identifies parity gaps, partial implementations, and Gable's advantages.
>
> **Note**: DMSi offers TWO ERP products — **Agility** (distribution/pro dealer) and **Frameworks** (retail/home centers). This document focuses specifically on **Frameworks**. For the Agility deep-dive, see [competitor_analysis_dmsi.md](file:///home/colton/Desktop/4C%20Digital_HQ/Lumber%20Digital%20Tools/GableLBM/GableERP/productdocs/competitor_analysis_dmsi.md).

---

## Executive Summary

DMSi Frameworks is a **100% web-based, end-to-end ERP** designed for LBM retail dealers and home centers. It covers POS, contractor sales, purchasing, inventory, accounting, dispatch, ecommerce, and reporting in a single browser-accessible platform.

GableERP shares the same web-native architecture advantage, but has **significant functional gaps** in accounting, POS, pricing controls, CRM, and reporting. However, GableERP leads in **AI capabilities, open APIs, millwork configurators, and yard management depth**.

| Dimension | Frameworks | GableERP | Winner |
|:---|:---|:---|:---:|
| Architecture | 100% web-based | 100% web-based | **Tie** |
| Retail POS | Full POS module | ❌ No dedicated POS | Frameworks |
| Contractor Sales / Pro Desk | Dedicated module | Quote Builder + Orders | Frameworks |
| Project Management Dashboard | Multi-job centralized view | ❌ Not built | Frameworks |
| Purchasing | Automated + buying group EDI | PO CRUD + vendor management | Frameworks |
| Inventory Control | Single-SKU multi-supplier, multi-UOM | Full CRUD, cycle count, yard layout | **Tie** |
| Accounting (GL/AR/AP) | Fully integrated suite | GL + AR only (no AP) | Frameworks |
| Dispatch & Delivery | Drag-drop routing, Google Maps, POD | Route/stop/delivery management, POD | **Tie** |
| Pricing Controls | Rebates, special buys, contract pricing | Escalator + price waterfall (10 files) | **Tie** |
| eCommerce / Service Portal | Product catalog, online ordering, CMS | Contractor portal (invoices, orders, deliveries) | Frameworks |
| Reporting & BI | Embedded analytics, dashboards, Phocas/F9 | AR aging, customer statements, dashboard | Frameworks |
| Mobile | Browser-responsive (no native app needed) | Browser-responsive + driver mobile + yard mobile | **Gable** |
| API & Integrations | Open APIs + EDI | Open REST API + EDI module | **Gable** |
| AI Capabilities | None native | Vision AI (tally + blueprints + material lists) | **Gable** |
| Millwork / Configurator | Not included (Agility feature) | Product + door configurator, blueprint verifier | **Gable** |
| Yard Management | Basic storage location tracking | Geospatial yard layout, pick queue, cycle count | **Gable** |
| Source Code | Proprietary / closed | Open source | **Gable** |
| Governance | Vendor-controlled | RFC community governance | **Gable** |

---

## Detailed Feature Comparison

### 1. Retail Point of Sale (POS) 🔴

> [!CAUTION]
> GableERP has **no dedicated POS module**. This is a critical gap for any dealer running a retail counter.

| Feature | Frameworks | GableERP | Gap |
|:---|:---:|:---:|:---|
| Retail counter sales screen | ✅ | ❌ | **MISSING** — No quick-sale POS interface |
| Easy-to-learn interface (seasonal staff) | ✅ | ❌ | **MISSING** — No simplified counter UI |
| Cash / check / card tender types | ✅ | ❌ | **MISSING** — No tender management |
| Daily cash reconciliation | ✅ | ⚠️ | [DailyTill.tsx](file:///home/colton/Desktop/4C%20Digital_HQ/Lumber%20Digital%20Tools/GableLBM/GableERP/app/src/pages/DailyTill.tsx) exists but no full POS flow |
| Receipt printing | ✅ | ❌ | **MISSING** |
| Returns / exchanges | ✅ | ❌ | **MISSING** |
| Role-based screen hiding | ✅ | ⚠️ | Admin module exists but no POS-specific customization |

**Impact**: Cannot serve walk-in retail customers without a POS. This blocks any dealer with a retail counter.

---

### 2. Contractor Sales / Pro Desk 🟡

| Feature | Frameworks | GableERP | Gap |
|:---|:---:|:---:|:---|
| Dedicated contractor sales workflow | ✅ | ⚠️ | [QuoteBuilder.tsx](file:///home/colton/Desktop/4C%20Digital_HQ/Lumber%20Digital%20Tools/GableLBM/GableERP/app/src/pages/QuoteBuilder.tsx) covers quoting but not contractor-specific flows |
| Project management dashboard | ✅ | ❌ | **MISSING** — No centralized multi-job view |
| Flexible quoting | ✅ | ✅ | QuoteBuilder + escalator pricing |
| Special / non-stock orders | ✅ | ⚠️ | Order system exists, no dedicated special-order workflow |
| Staggered delivery scheduling | ✅ | ⚠️ | Delivery module exists but no staggered scheduling UI |
| Pick ticket generation | ✅ | ✅ | [PickQueue.tsx](file:///home/colton/Desktop/4C%20Digital_HQ/Lumber%20Digital%20Tools/GableLBM/GableERP/app/src/pages/yard/PickQueue.tsx) + [PickDetail.tsx](file:///home/colton/Desktop/4C%20Digital_HQ/Lumber%20Digital%20Tools/GableLBM/GableERP/app/src/pages/yard/PickDetail.tsx) |

**Impact**: GableERP can support contractor sales through its existing quote/order flow, but lacks the unified "Pro Desk" dashboard that ties jobs, deliveries, and invoices together.

---

### 3. Purchasing & Procurement 🟡

| Feature | Frameworks | GableERP | Gap |
|:---|:---:|:---:|:---|
| Purchase order creation | ✅ | ✅ | [NewPurchaseOrder.tsx](file:///home/colton/Desktop/4C%20Digital_HQ/Lumber%20Digital%20Tools/GableLBM/GableERP/app/src/pages/purchasing/NewPurchaseOrder.tsx) |
| PO list / tracking | ✅ | ✅ | [PurchaseOrderList.tsx](file:///home/colton/Desktop/4C%20Digital_HQ/Lumber%20Digital%20Tools/GableLBM/GableERP/app/src/pages/purchasing/PurchaseOrderList.tsx) |
| Vendor management | ✅ | ✅ | [VendorList.tsx](file:///home/colton/Desktop/4C%20Digital_HQ/Lumber%20Digital%20Tools/GableLBM/GableERP/app/src/pages/purchasing/VendorList.tsx) + [VendorDetail.tsx](file:///home/colton/Desktop/4C%20Digital_HQ/Lumber%20Digital%20Tools/GableLBM/GableERP/app/src/pages/purchasing/VendorDetail.tsx) |
| Purchasing recommendations | ✅ | ❌ | **MISSING** — No automated reorder suggestions |
| Automated purchasing | ✅ | ❌ | **MISSING** — No auto-PO generation |
| Buying group EDI integration | ✅ | ⚠️ | EDI module exists (`backend/internal/edi`) but no buying group connections |
| Real-time pricing from catalogs | ✅ | ❌ | **MISSING** — No supplier catalog sync |
| PO receiving with verification | ✅ | ✅ | [ReceivePO.tsx](file:///home/colton/Desktop/4C%20Digital_HQ/Lumber%20Digital%20Tools/GableLBM/GableERP/app/src/pages/yard/ReceivePO.tsx) |

**Impact**: Core PO workflow exists. Missing the *intelligence* layer — automated suggestions, buying group integration, and real-time catalog pricing.

---

### 4. Inventory Control 🟢

| Feature | Frameworks | GableERP | Gap |
|:---|:---:|:---:|:---|
| Product catalog / SKU management | ✅ | ✅ | Backend `product` + `inventory` modules |
| Single SKU, multi-supplier | ✅ | ⚠️ | Vendor module exists but unclear multi-supplier per SKU |
| Multi-UOM support | ✅ | ✅ | Pricing module handles UOM conversions (10 files) |
| Multi-branch inventory sync | ✅ | ✅ | Location module (`backend/internal/location`) |
| Storage location tracking | ✅ | ✅ | [YardLayout.tsx](file:///home/colton/Desktop/4C%20Digital_HQ/Lumber%20Digital%20Tools/GableLBM/GableERP/app/src/pages/yard/YardLayout.tsx) — **Gable exceeds** with geospatial layout |
| Cycle counting | ✅ | ✅ | [CycleCount.tsx](file:///home/colton/Desktop/4C%20Digital_HQ/Lumber%20Digital%20Tools/GableLBM/GableERP/app/src/pages/yard/CycleCount.tsx) |
| Inventory lookup / search | ✅ | ✅ | [InventoryLookup.tsx](file:///home/colton/Desktop/4C%20Digital_HQ/Lumber%20Digital%20Tools/GableLBM/GableERP/app/src/pages/yard/InventoryLookup.tsx) |
| Dead stock / turns reporting | ✅ | ❌ | **MISSING** — No inventory analytics |
| Barcode scanning | ✅ | ❌ | **MISSING** — No native scanner integration |

**Impact**: GableERP is at near-parity on inventory CRUD. Gaps are in analytics (turns/dead stock) and barcode scanning hardware integration.

---

### 5. Accounting & Financials 🔴

> [!CAUTION]
> Frameworks ships a **complete accounting suite** (GL, AR, AP, ACH, bank reconciliation). GableERP has GL and AR but **no Accounts Payable, no bank reconciliation, no ACH**.

| Feature | Frameworks | GableERP | Gap |
|:---|:---:|:---:|:---|
| General Ledger | ✅ | ✅ | [ChartOfAccounts.tsx](file:///home/colton/Desktop/4C%20Digital_HQ/Lumber%20Digital%20Tools/GableLBM/GableERP/app/src/pages/accounting/ChartOfAccounts.tsx), [JournalEntries.tsx](file:///home/colton/Desktop/4C%20Digital_HQ/Lumber%20Digital%20Tools/GableLBM/GableERP/app/src/pages/accounting/JournalEntries.tsx), [TrialBalance.tsx](file:///home/colton/Desktop/4C%20Digital_HQ/Lumber%20Digital%20Tools/GableLBM/GableERP/app/src/pages/accounting/TrialBalance.tsx) |
| Accounts Receivable | ✅ | ✅ | Invoice module + [ARAgingReport.tsx](file:///home/colton/Desktop/4C%20Digital_HQ/Lumber%20Digital%20Tools/GableLBM/GableERP/app/src/pages/reports/ARAgingReport.tsx) + [CustomerStatementPage.tsx](file:///home/colton/Desktop/4C%20Digital_HQ/Lumber%20Digital%20Tools/GableLBM/GableERP/app/src/pages/reports/CustomerStatementPage.tsx) |
| Accounts Payable | ✅ | ❌ | **MISSING** — No vendor invoice entry, PO matching, payment scheduling |
| ACH payments | ✅ | ❌ | **MISSING** |
| Bank reconciliation | ✅ | ❌ | **MISSING** |
| Automated period closing | ✅ | ❌ | **MISSING** |
| Customized financial statements | ✅ | ❌ | **MISSING** — Trial balance only |
| Drilldown from GL summary to transactions | ✅ | ⚠️ | Journal entries page partial |
| GL account search (advanced filters) | ✅ | ⚠️ | Chart of accounts has basic search |

**Impact**: Dealers cannot "run their books" in GableERP alone. They still need QuickBooks/NetSuite for AP and bank reconciliation. This is the #1 blocker for full ERP replacement.

---

### 6. Pricing Controls 🟢

| Feature | Frameworks | GableERP | Gap |
|:---|:---:|:---:|:---|
| Contract pricing by customer | ✅ | ✅ | Pricing module (10 files including escalator logic) |
| Customer group pricing | ✅ | ✅ | Pricing waterfall |
| Job-specific pricing | ✅ | ✅ | Escalator pricing per quote |
| Rebate management | ✅ | ❌ | **MISSING** — No vendor rebate tracking |
| Special buy pricing | ✅ | ⚠️ | Can be handled via pricing rules |
| Price escalation (commodity volatility) | ✅ | ✅ | Backend `escalator_service.go` — **Gable advantage**: live market index integration |

**Impact**: GableERP is **strong** on pricing. The escalator + market index integration is a genuine advantage. Gap is rebate management.

---

### 7. Dispatch & Delivery 🟢

| Feature | Frameworks | GableERP | Gap |
|:---|:---:|:---:|:---|
| Dispatch board | ✅ | ✅ | [DispatchBoard.tsx](file:///home/colton/Desktop/4C%20Digital_HQ/Lumber%20Digital%20Tools/GableLBM/GableERP/app/src/pages/DispatchBoard.tsx) |
| Drag-and-drop route building | ✅ | ⚠️ | Dispatch board exists, unclear if drag-drop |
| Google Maps route review | ✅ | ❌ | **MISSING** — No map integration |
| Proof of delivery (signatures/photos) | ✅ | ✅ | Driver mobile pages ([DeliveryDetail.tsx](file:///home/colton/Desktop/4C%20Digital_HQ/Lumber%20Digital%20Tools/GableLBM/GableERP/app/src/pages/driver/DeliveryDetail.tsx)) |
| Route list management | ✅ | ✅ | [RouteList.tsx](file:///home/colton/Desktop/4C%20Digital_HQ/Lumber%20Digital%20Tools/GableLBM/GableERP/app/src/pages/driver/RouteList.tsx) |
| Stop-by-stop tracking | ✅ | ✅ | [StopList.tsx](file:///home/colton/Desktop/4C%20Digital_HQ/Lumber%20Digital%20Tools/GableLBM/GableERP/app/src/pages/driver/StopList.tsx) |
| Automatic customer notifications | ✅ | ❌ | **MISSING** — No SMS/email delivery status updates |
| Driver on-site invoice adjustments | ✅ | ❌ | **MISSING** |

**Impact**: GableERP has strong delivery foundations with the driver mobile experience. Missing Google Maps integration and automated customer notifications.

---

### 8. eCommerce / Service Portal 🟡

| Feature | Frameworks | GableERP | Gap |
|:---|:---:|:---:|:---|
| Customer login / dashboard | ✅ | ✅ | [PortalDashboard.tsx](file:///home/colton/Desktop/4C%20Digital_HQ/Lumber%20Digital%20Tools/GableLBM/GableERP/app/src/pages/portal/PortalDashboard.tsx) + [PortalLogin.tsx](file:///home/colton/Desktop/4C%20Digital_HQ/Lumber%20Digital%20Tools/GableLBM/GableERP/app/src/pages/portal/PortalLogin.tsx) |
| Invoice viewing / download | ✅ | ✅ | [PortalInvoices.tsx](file:///home/colton/Desktop/4C%20Digital_HQ/Lumber%20Digital%20Tools/GableLBM/GableERP/app/src/pages/portal/PortalInvoices.tsx) |
| Order history / tracking | ✅ | ✅ | [PortalOrders.tsx](file:///home/colton/Desktop/4C%20Digital_HQ/Lumber%20Digital%20Tools/GableLBM/GableERP/app/src/pages/portal/PortalOrders.tsx) |
| Delivery tracking | ✅ | ✅ | [PortalDeliveries.tsx](file:///home/colton/Desktop/4C%20Digital_HQ/Lumber%20Digital%20Tools/GableLBM/GableERP/app/src/pages/portal/PortalDeliveries.tsx) |
| Product catalog browsing | ✅ | ❌ | **MISSING** — No catalog in portal |
| Online ordering / cart | ✅ | ❌ | **MISSING** — No cart/checkout in portal |
| Product images | ✅ | ❌ | **MISSING** — No product image management |
| No-code CMS | ✅ | ❌ | **MISSING** |
| Guest user administration | ✅ | ❌ | **MISSING** — No multi-user portal access management |
| 24/7 self-service access | ✅ | ✅ | Portal is web-accessible |

**Impact**: GableERP's portal handles the "back office" (invoices, orders, deliveries) well, but lacks the "front door" (catalog browsing, online ordering). Frameworks has a more complete B2C-style experience.

---

### 9. Reporting & Business Intelligence 🟡

| Feature | Frameworks | GableERP | Gap |
|:---|:---:|:---:|:---|
| Real-time dashboards | ✅ | ✅ | [Dashboard.tsx](file:///home/colton/Desktop/4C%20Digital_HQ/Lumber%20Digital%20Tools/GableLBM/GableERP/app/src/pages/Dashboard.tsx) |
| Embedded data analytics | ✅ | ⚠️ | Dashboard with KPIs, no embedded analytics engine |
| Configurable/custom reports | ✅ | ❌ | **MISSING** — Only predefined reports |
| AR Aging | ✅ | ✅ | ARAgingReport.tsx |
| Customer statements | ✅ | ✅ | CustomerStatementPage.tsx |
| Drilldown from summary to detail | ✅ | ⚠️ | Some drilldown in dashboard |
| F9 / Excel-based reporting | ✅ | ❌ | **MISSING** — No Excel export integration |
| Third-party BI (Phocas, CyberScience) | ✅ | ❌ | **MISSING** — No BI tool integration |
| Scheduled report distribution | ✅ | ❌ | **MISSING** |

**Impact**: GableERP has essential operational dashboards but no ad-hoc reporting or BI integration. Dealers who rely on custom reports (most do) will find this limiting.

---

### 10. Integrations & Platform 🟢

| Feature | Frameworks | GableERP | Gap |
|:---|:---:|:---:|:---|
| Open APIs | ✅ | ✅ | **Gable advantage** — free, documented REST APIs |
| EDI capabilities | ✅ | ✅ | `backend/internal/edi` module |
| Payment processing | ✅ (unspecified) | ✅ | Payment module (`backend/internal/payment`) |
| Sales tax automation (Avalara) | ✅ | ❌ | **MISSING** |
| Third-party integrations | ✅ | ⚠️ | Integrations module exists (`backend/internal/integrations`) |
| Mobile-friendly / responsive | ✅ | ✅ | Browser-responsive + dedicated mobile pages |

**Impact**: GableERP's open API story is stronger than Frameworks'. The sales tax gap (Avalara) is a P0 blocker.

---

## Gap Priority Matrix

### 🔴 Critical Gaps (Must close for Frameworks competitive parity)

| # | Gap | Impact |
|:---|:---|:---|
| 1 | **Retail POS module** | Cannot serve walk-in retail customers |
| 2 | **Accounts Payable** | Cannot run dealer's full books |
| 3 | **Product catalog in portal** | Cannot compete with Frameworks' online browsing |
| 4 | **Sales tax automation (Avalara)** | Compliance blocker for any production deployment |
| 5 | **Automated purchasing recommendations** | Missing intelligence that Frameworks offers |

### 🟡 Important Gaps (Needed for competitive demos)

| # | Gap | Impact |
|:---|:---|:---|
| 6 | **Project management dashboard** | No centralized multi-job view for pro desk |
| 7 | **Online ordering / cart in portal** | Cannot match Frameworks' self-service |
| 8 | **Google Maps dispatch integration** | Dispatchers lack visual route planning |
| 9 | **Bank reconciliation & ACH** | Accounting suite incomplete |
| 10 | **Rebate management** | Missing for dealers with buying group rebates |
| 11 | **Custom / ad-hoc reporting** | Dealers need more than predefined reports |
| 12 | **Customer delivery notifications** | No SMS/email status updates |
| 13 | **Barcode scanning** | Warehouse staff need scanning workflows |
| 14 | **Guest portal user administration** | Contractors can't manage own team's access |

### 🟢 Gable Advantages (Where Gable already surpasses Frameworks)

| # | Advantage | Details |
|:---|:---|:---|
| 1 | **AI Vision capabilities** | Tally, material list parsing, blueprint verification — Frameworks has nothing |
| 2 | **Millwork / product configurator** | Frameworks doesn't include this (Agility-only); Gable has [ProductConfigurator.tsx](file:///home/colton/Desktop/4C%20Digital_HQ/Lumber%20Digital%20Tools/GableLBM/GableERP/app/src/pages/millwork/ProductConfigurator.tsx), [DoorConfigurator.tsx](file:///home/colton/Desktop/4C%20Digital_HQ/Lumber%20Digital%20Tools/GableLBM/GableERP/app/src/pages/millwork/DoorConfigurator.tsx), [BlueprintVerifier.tsx](file:///home/colton/Desktop/4C%20Digital_HQ/Lumber%20Digital%20Tools/GableLBM/GableERP/app/src/pages/millwork/BlueprintVerifier.tsx) |
| 3 | **Geospatial yard management** | Gable's yard layout with spatial awareness > Frameworks' basic location tracking |
| 4 | **Escalator + market index pricing** | Live commodity index integration > Frameworks' static escalator rules |
| 5 | **Open source / open API** | Free integration vs. proprietary lock-in |
| 6 | **RFC governance** | Community-driven roadmap vs. vendor-controlled |
| 7 | **Dedicated driver mobile experience** | Route/stop/delivery detail pages purpose-built for drivers |
| 8 | **Yard mobile experience** | Pick queue, cycle count, inventory lookup, PO receiving on mobile |

---

## Recommended Action Plan

> [!IMPORTANT]
> Against **Frameworks** specifically (the retail ERP), GableERP's #1 gap is the **Retail POS module**. Frameworks was purpose-built for retail counter sales. Without a POS, GableERP cannot credibly compete for any dealer with a retail counter operation.

### Phase 1 — Parity Essentials (Before any Frameworks-competitive demo)
1. Build a **Retail POS module** (simplified counter sales screen with tender management)
2. Build **Accounts Payable** (vendor invoice entry, PO matching, payment scheduling)
3. Integrate **Avalara** for sales tax automation
4. Add **product catalog browsing** to contractor portal

### Phase 2 — Competitive Feature Match
5. Build **Project Management Dashboard** (centralize jobs, orders, deliveries)
6. Add **online ordering / cart** to portal
7. Integrate **Google Maps** into dispatch module
8. Build **bank reconciliation** and **ACH** support
9. Add **automated purchasing recommendations**
10. Build **custom report builder** or integrate a BI tool

### Phase 3 — Lean Into Advantages
11. Market AI Vision capabilities as the differentiator Frameworks can't match
12. Position millwork configurator as "Agility feature at Frameworks price"
13. Lead with open source + open API messaging against proprietary lock-in
14. Emphasize sovereign dealer portal vs. Frameworks' platform-controlled portal

---

## Conclusion

GableERP and DMSi Frameworks share the same **web-native, browser-first** architecture, making them direct competitors. Frameworks has the advantage of **completeness** — it ships a full POS-to-GL system that a retail dealer can run their entire business on. GableERP has the advantage of **depth and innovation** — AI vision, millwork configurators, geospatial yard mgmt, and an open source model that Frameworks cannot match.

The path to competitive parity requires closing **5 critical gaps** (POS, AP, portal catalog, Avalara, purchasing intelligence). Once those are closed, GableERP's advantages in AI, configurators, yard management, and open architecture should make it the stronger overall platform for forward-looking LBM dealers.
