# Logistics & Dispatch Specification

## 1. Overview
Delivery is the defining differentiator for LBM dealers. Customers choose dealers based on ability to deliver on time. This module covers Yard Management (picking), Truck Scheduling, Route Optimization, and Driver Operations.

## 2. Core Entities

### 2.1. Pick Ticket
Generated when an order is ready for fulfillment.
```
PickTicket {
    ID              UUID
    OrderID         UUID
    Status          Enum    // "Pending", "Assigned", "Picking", "Picked", "Loaded"
    AssignedTo      UUID?   // Yard foreman/picker
    Priority        Enum    // "Standard", "Rush", "Will Call"
    
    Lines           []PickLine
}

PickLine {
    ProductID       UUID
    LocationHint    String  // "Yard A, Row 4, Bin 2"
    RequestedQty    Decimal
    PickedQty       Decimal
    UOM             String
}
```

### 2.2. Delivery Run
A scheduled truck route with multiple stops.
```
DeliveryRun {
    ID              UUID
    TruckID         UUID
    DriverID        UUID
    ScheduledDate   Date
    Status          Enum    // "Planned", "Loading", "EnRoute", "Completed"
    
    Stops           []DeliveryStop
}

DeliveryStop {
    Sequence        Int
    OrderID         UUID
    Address         Address
    TimeWindow      TimeWindow // e.g., "8am-10am"
    ActualArrival   Timestamp?
    Status          Enum    // "Pending", "Arrived", "Delivered", "Failed"
}
```

### 2.3. Proof of Delivery (POD)
```
ProofOfDelivery {
    DeliveryStopID  UUID
    SignatureImage  Blob?
    SignerName      String?
    GPSLat          Float
    GPSLong         Float
    Timestamp       Timestamp
    Photos          []Blob  // Photos of drop site, damage, etc.
    Notes           String?
}
```

## 3. Yard Management Workflow

### 3.1. Picking Flow
```
1. Pick Ticket created (auto from Order Confirm, or manual).
2. Ticket assigned to a Picker (or self-selected via mobile app).
3. Picker navigates to location (map view on tablet).
4. Picker scans/confirms product, enters quantity.
5. If shortage, Picker flags for backorder or substitution.
6. On completion, status -> "Picked".
```

### 3.2. Load Building
```
1. Picker stages materials in Loading Zone.
2. Dispatcher assigns Picked tickets to a Delivery Run.
3. System suggests load order (LIFO - last stop on bottom).
4. Forklift driver loads truck per suggested order.
5. On completion, status -> "Loaded", Run status -> "Loading".
```

## 4. Route Optimization

### 4.1. Constraints
- **Time Windows:** Customer-requested delivery times.
- **Truck Capacity:** Weight and cubic limits.
- **Product Restrictions:** Some items can't mix (e.g., treated lumber + drywall in rain).
- **Driver Hours:** DOT compliance for CDL drivers.

### 4.2. Optimization Engine
- **Algorithm:** Vehicle Routing Problem (VRP) solver.
- **AI Enhancement:** Learns from "Dispatcher Overrides" to improve suggestions.
- **Real-Time:** Re-routes if a truck is delayed or a customer cancels.

## 5. Driver Mobile App Features
| Feature                | Description                                        |
|------------------------|----------------------------------------------------|
| Turn-by-Turn Nav       | Integrated maps with stop sequence.                |
| Digital POD            | Capture signature, photos, GPS stamp.              |
| Customer Communication | One-tap call or text to customer.                  |
| Issue Reporting        | Report damaged goods, refused delivery, delays.    |
| Load Checklist         | Verify items before leaving yard.                  |

## 6. Will Call (Customer Pickup)
- Order marked as "Will Call" skips delivery scheduling.
- Counter notifies customer when "Ready for Pickup".
- Yard releases on customer signature (ID check configurable).

## 7. Reporting
| Report Name             | Description                                      |
|-------------------------|--------------------------------------------------|
| On-Time Delivery Rate   | % of deliveries within promised window.          |
| Truck Utilization       | Capacity used vs. available.                     |
| Driver Efficiency       | Stops per hour, miles per stop.                  |
| Failed Deliveries       | Reasons and frequency.                           |
