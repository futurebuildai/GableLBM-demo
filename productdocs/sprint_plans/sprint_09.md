# Sprint 9: Logistics & Dispatch MVP

## Goal
Enable the scheduling, staging, and final delivery of orders, completing the Phase 2 "Transaction Engine" loop.

## Objectives
1.  **Delivery Management**:
    *   [ ] Delivery Schema (Routes, Stops, Vehicles, Drivers).
    *   [ ] Link Sales Orders to Deliveries.
    *   [ ] Dispatch Dashboard (Assign Orders to Vehicles).

2.  **Yard Operations**:
    *   [ ] Digital Pick Ticket (Mobile-friendly view).
    *   [ ] Staging Status Updates (Picking -> Staged -> Loaded).

3.  **Proof of Delivery (POD)**:
    *   [ ] Basic POD Capture (Signature, Photo upload).
    *   [ ] "Delivered" Status update.

## Technical Constraints
*   **Mobile First**: Yard views must work on small screens.
*   **Offline Capable**: (Stretch) POD should work with flaky connection (optimistic updates).
