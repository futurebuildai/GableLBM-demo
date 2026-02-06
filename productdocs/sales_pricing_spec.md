# Sales & Pricing Engine Specification

## 1. Overview
LBM pricing is notoriously complex. A single SKU can have dozens of applicable prices depending on customer, job, volume, and time of year. This spec defines the "Pricing Waterfall" that the engine will follow.

## 2. Pricing Waterfall (Priority Order)
The engine evaluates prices in this order, stopping at the first match.

```
1. [HIGHEST] Contract/Net Price (Locked for a specific Customer+SKU)
2. Job-Level Override (Specific price for a Project/Job)
3. Promotional/Sale Price (Time-bound discount)
4. Volume Break (Quantity-based discount)
5. Customer Price Level (Tiered contractor pricing: L1-L5)
6. [LOWEST] Retail/Base Price
```

## 3. Data Model

### 3.1. Price Level
Represents a pricing tier (e.g., "Retail", "Contractor Level 1", "Contractor Level 5").
```
PriceLevel {
    ID          UUID
    Name        String   // "Level 3"
    Description String   // "Standard Contractor"
    Multiplier  Decimal  // 0.85 (15% off retail)
}
```

### 3.2. Price Rule
A specific override or promotion.
```
PriceRule {
    ID             UUID
    ProductID      UUID?   // Null = applies to all products
    CategoryID     UUID?   // Null = applies to all categories
    CustomerID     UUID?   // Null = applies to all customers
    JobID          UUID?   // Null = applies to all jobs
    
    Type           Enum    // "NetPrice", "Discount%", "DiscountAmt"
    Value          Decimal
    
    MinQty         Decimal // For volume breaks
    StartDate      Date?
    EndDate        Date?
}
```

## 4. Quote & Order Flow

### 4.1. Quote States
```
Draft -> Sent -> Viewed -> Accepted/Rejected/Expired
```
- **Versioning:** Revisions create a new Quote linked to the original (`RootQuoteID`).
- **Margin Lock:** On "Sent", the margin is locked to prevent cost fluctuation during negotiation.

### 4.2. Order States
```
Pending -> Confirmed -> Picking -> Loaded -> Shipped -> Delivered -> Invoiced
```
- **Partial Shipments:** An order can have multiple shipments.
- **Backorder Logic:** Items not in stock are flagged for automatic PO generation.

## 5. Special Order Workflow
For items not in stock (Millwork, Custom Windows).

```
1. Customer requests item.
2. System queries vendor (via EDI or API) for cost/lead time.
3. Counter rep creates SO Line with vendor cost + margin.
4. System may require deposit (configurable % or amount).
5. On confirmation, system auto-generates PO to vendor.
6. On receipt, system links inbound PO to original SO line.
7. Customer is notified (SMS/Email) that item is ready.
```

## 6. Credit Check Integration
At the point of "Confirm Order", the system checks:
1. Customer Credit Limit vs. (Outstanding AR + This Order Total).
2. Job Credit Limit (if applicable).
3. If over limit, order is flagged for `controller` approval.
