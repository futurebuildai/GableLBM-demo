# Sprint 10: Advanced Logistics & Mobile Driver App

## Goal
Empower the fleet with a dedicated Mobile Driver App for Proof of Delivery (POD) and enable advanced inventory operations (Transfers).

## Objectives

### 1. Mobile Driver App (POD)
- **Driver "My Routes" View**: Filter routes by Driver ID.
- **Stop List**: Mobile-optimized view of the day's stops.
- **Delivery Detail**: Address, Instructions, Items.
- **POD Capture**: 
  - Signature Pad.
  - Photo Upload (Camera integration).
  - Status Updates (En Route, Arrived, Delivered).

### 2. Inventory Transfers
- **Backend Core**: Logic to move stock between Bins (Audit log consistency).
- **Frontend UI**: Transfer Modal/Page (From Bin -> To Bin).

### 3. Route Optimization (Basic)
- **Manual Reordering**: Drag-and-drop stop sequencing in Dispatch Board.
- **Driver Filtering**: Update backend to list routes by Driver.

## Technical Constraints
- **Mobile First**: Driver App must work on 375px width devices.
- **Optimistic UI**: Simple local state for POD before upload.
