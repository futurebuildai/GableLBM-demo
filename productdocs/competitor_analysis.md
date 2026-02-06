# Competitor Analysis & Feature Parity Matrix

## 1. Executive Summary
*To be filled after completion of all deep dives.*

## 2. Competitive Matrix (Feature Grid)
| Feature Category | Feature | Epicor BisTrack | ECI Spruce | DMSi Agility | GableLBM (Target) |
| :--- | :--- | :--- | :--- | :--- | :--- |
| **Inventory** | Multi-UOM (MBF vs PC) | ✅ Strong Native Support | ✅ Native | ✅ Best-in-Class (Millwork) | ✅ Native |
| | Remnant Tracking | ✅ Via "Advanced UOM" | ❌ Limited | ✅ Native | ✅ Native |
| | Mull Unit Logic | ✅ Via "Production" Config | ❌ No Native Support | ✅ Native (E-Catalog) | ✅ Native |
| **Sales** | Quick Quote | ⚠️ Complaint: Slow/Complex | ✅ fast (POS focused) | ✅ "OrderPad" (Strong) | ✅ < 3 Clicks |
| | Contract Pricing | ✅ Robust "Logic Blocks" | ✅ Standard | ✅ Advanced | ✅ Real-time |
| **Logistics** | Load Building | ✅ "Journey Planner" | ❌ Basic | ✅ "Journey Planner" | ✅ Auto-optimize |
| | POD Mobile App | ✅ "BisTrack Delivery" | ✅ "ProLink Delivery" | ✅ "Agility Mobile POD" | ✅ Native App |
| **Labor** | Yard Crew Scheduling | ✅ "Field Scheduling" Module | ❌ Limited | ✅ via Mobile Warehouse | ✅ Integrated |
| **Tech** | API Openness | ✅ REST & SOAP (Hybrid) | ✅ REST & SOAP | ⚠️ Private / Gated | ✅ Open REST/Events |
| | Cloud Native? | ❌ Heavy Client (VPN Issues) | ✅ Cloud Options | ❌ Heavy Client (Citrix) | ✅ Yes |

## 3. Deep Dive: Epicor BisTrack
**Market Position**: Enterprise Standard. High complexity, high power.

### 3.1 Inventory & Operations
*   **Strengths**: "Advanced UOM" handles the complexity of lumber (MBF, PC conversion). Excellent support for "Remanufacturing" (milling raw lumber into finished goods).
*   **Weaknesses**: "Mull units" and complex millwork often require the heavy "Production" module, which is overkill for simple yard ops.

### 3.2 Sales & pricing
*   **Strengths**: "Logic Blocks" allow for infinitely customizable pricing rules (e.g., "Customer A gets 5% off if they buy > 100 studs").
*   **Weaknesses**: Sales Ops is often slow. Users report 30-45s load times for forms over VPN. "Quick Quote" is a misnomer; it's a full ERP form.

### 3.3 Logistics & Labor
*   **Strengths**: "BisTrack Delivery" mobile app creates a solid audit trail (signatures, photos, geo-stamps). "Field Scheduling" allows drag-and-drop crew assignment.
*   **Weaknesses**: The scheduling interface is often separate from the main sales flow ("Journey Planner"), creating a disconnect between Sales promise dates and Dispatch reality.

### 3.4 Integration & Tech Stack
*   **Stack**: Microsoft SQL Server / .NET.
*   **API**: Good REST API coverage (modern), plus legacy SOAP.
*   **Architecture**: "Smart Client" (Fat Client) architecture. Requires heavy local install or RDP/Citrix for remote access. This is the **#1 User Complaint** (Performance over VPN).

### 3.5 Modernization Opportunities (The "Gable Advantage")
1.  **True Cloud Native**: Eliminate the "VPN Tax" (slow forms, RDP costs) with a React/Vite web frontend that runs anywhere.
2.  **Unified Scheduling**: Bring "Field Scheduling" directly into the "Sales Order" view so sales reps *know* the delivery slot is available instantly.
3.  **Speed**: Replace 45-second form loads with <100ms API calls.


## 4. Deep Dive: ECI Spruce
**Market Position**: Mid-Market Generalist. "Jack of all trades" for hardware stores and lumberyards.

### 4.1 Inventory & Operations
*   **Strengths**: Strong POS integration (best-in-class for retail hardware + lumber mix). "Automated Ordering" is solid for retail goods.
*   **Weaknesses**: Less depth in heavy "Manufacturing" or "Millwork" compared to BisTrack. "Phantom Stock" issues reported by users.

### 4.2 Sales & pricing
*   **Strengths**: "Painless" POS transactions. Excellent for "Cash & Carry" lumber sales.
*   **Weaknesses**: Quote-to-Order workflow can feel disjointed.

### 4.3 Logistics & Labor
*   **Strengths**: "ProLink" Customer Portal allows customers to schedule their own deliveries (unique feature). "Mobile Delivery" app handles signatures/photos offline.
*   **Weaknesses**: "ProLink" is an external portal, not fully native.

### 4.4 Integration & Tech Stack
*   **Stack**: .NET / SQL. Cloud-hosted options available.
*   **API**: Hybrid REST and SOAP endpoints.
*   **Complaints**: "Beta state" feeling for new features. Support is slow.

### 4.5 Modernization Opportunities (The "Gable Advantage")
1.  **Unified Portal**: Instead of a separate "ProLink", the customer portal is just a permission view of the main app.
2.  **Reliability**: Offline-first architecture (PWA) to solve the "rural internet disconnect" data loss issues reported by Spruce users.


## 5. Deep Dive: DMSi Agility
**Market Position**: The Heavyweight. Dominant in large millwork, pro-dealers, and two-step distribution.

### 5.1 Inventory & Operations
*   **Strengths**: Unmatched **Millwork** support (doors, windows, molding). "E-Catalog" modules allow specialized configuration of complex assemblies (e.g., pre-hung doors).
*   **Weaknesses**: Can be overkill for smaller yards. Complexity of "remanning" (remanufacturing) features requires significant training.

### 5.2 Sales & pricing
*   **Strengths**: **"OrderPad"** is a standout feature—allows sales reps to build orders rapidly with job-specific pricing. Highly regarded by users.
*   **Weaknesses**: Integration silos—some users report needing manual tracking for processes that fall outside the core "Happy Path".

### 5.3 Logistics & Labor
*   **Strengths**: **"Agility Mobile POD"** is excellent. Real-time updates, photos, signatures, and GPS. Drivers can adjust invoice quantities on-site (solving the "short-ship" billing nightmare).
*   **Weaknesses**: Like BisTrack, typically requires Citrix/RDP for full desktop access, meaning dispatchers are tethered to office workstations.

### 5.4 Integration & Tech Stack
*   **Stack**: Proprietary / Legacy Core.
*   **API**: Exists (SaberisConnect uses it), but documentation is private/gated. Not a true "Open API" platform.
*   **User Sentiment**: 95% Satisfaction (Highest in class). Users love the functionality but dislike the legacy deployment model.

### 5.5 Modernization Opportunities (The "Gable Advantage")
1.  **Open API First**: DMSi grants API access selectively. GableLBM will document every endpoint publicly (Swagger/OpenAPI).
2.  **Democratized Millwork**: Agility's millwork tools are expensive add-ons. GableLBM should include basic "Assembly/Kitting" logic in the core open-source tier.
3.  **No Citrix**: Browser-native performance without the latency of remote desktop protocols.

