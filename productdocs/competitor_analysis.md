# Competitor Analysis & Feature Parity Matrix

## 1. Executive Summary
*To be filled after completion of all deep dives.*

## 2. Competitive Matrix (Feature Grid)
| Feature Category | Feature | Epicor BisTrack | ECI Spruce | DMSi Agility | GableLBM (Target) |
| :--- | :--- | :--- | :--- | :--- | :--- |
| **Inventory** | Multi-UOM (MBF vs PC) | ✅ Strong Native Support | | | ✅ Native |
| | Remnant Tracking | ✅ Via "Advanced UOM" | | | ✅ Native |
| | Mull Unit Logic | ✅ Via "Production" Config | | | ✅ Native |
| **Sales** | Quick Quote | ⚠️ Complaint: Slow/Complex | | | ✅ < 3 Clicks |
| | Contract Pricing | ✅ Robust "Logic Blocks" | | | ✅ Real-time |
| **Logistics** | Load Building | ✅ "Journey Planner" | | | ✅ Auto-optimize |
| | POD Mobile App | ✅ "BisTrack Delivery" | | | ✅ Native App |
| **Labor** | Yard Crew Scheduling | ✅ "Field Scheduling" Module | | | ✅ Integrated |
| **Tech** | API Openness | ✅ REST & SOAP (Hybrid) | | | ✅ Open REST/Events |
| | Cloud Native? | ❌ Heavy Client (VPN Issues) | | | ✅ Yes |

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
*Pending Research...*

## 5. Deep Dive: DMSi Agility
*Pending Research...*
