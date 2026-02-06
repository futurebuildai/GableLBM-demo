# Sprint 03: Inventory Locations & Operations

**Goal**: Enable physical inventory tracking (Where is it?) and basic stock flux (In/Out/Count).
**Prerequisites**: Sprint 02 (Product/UOM Schema) complete.

## 1. Scope

### 1.1 Location Management (The "Bin" Logic)
*   **Schema**: `locations` table (Zone, Aisle, Bay, Level, Bin).
*   **Logic**:
    *   Hierarchy support (e.g., "Yard A" -> "Row 1").
    *   Validation: Cannot delete location if inventory exists.
*   **UI**:
    *   Location List/Tree View.
    *   "Create Location" Modal.

### 1.2 Inventory Mapping
*   **Schema**: Update `inventory` table to link `location_id`.
*   **Logic**:
    *   One SKU can exist in multiple locations.
    *   Total Stock = Sum of all location records for that SKU.

### 1.3 Basic Operations ("In/Out")
*   **Receipts**: Manual entry to add stock to a specific location.
*   **Cycle Count**: Manual adjustment of stock in a specific location (with Reason Code).
*   **Audit Log**: Record *who* changed *what* and *why*.

## 2. Technical Requirements
*   **Backend**:
    *   New `LocationService`.
    *   Update `InventoryService` to be location-aware.
*   **Frontend**:
    *   `LocationManager` component.
    *   Update `InventoryTable` to show location breakdowns.
    *   `StockAdjustment` Modal.

## 3. Success Criteria
*   [ ] Can create a hierarchy of locations (Zone -> Bin).
*   [ ] Can receive 100 2x4s into "Bin A1".
*   [ ] Can move 50 2x4s from "Bin A1" to "Bin B2".
*   [ ] Can adjust count down by 2 (Damaged) and see reason code in history.
