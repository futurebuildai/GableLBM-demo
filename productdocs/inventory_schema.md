# Inventory Data Model: The LBM Core

## Overview
Lumber inventory is unique because the "Unit" changes state. You buy a "Bunk" (large pack), count it in "Board Feet" (Volume), stock it as "Pieces" (Count), and sell it in "Linear Feet" (Length).

## Core Schema (PostgreSQL / Go Structs)

### 1. Product Definition
```go
type Product struct {
    ID               uuid.UUID
    SKU              string    // "2X4-8-SPF"
    Name             string    // "2x4 8' SPF Stud Premium"
    
    // Core Logic
    BaseUOMID        uuid.UUID // Link to "Piece"
    Category         string    // "Dimensional Lumber", "Fasteners", "Millwork"
    
    // Physical Properties (Vital for Load Calcs)
    Weight           float64   // Lbs per BaseUOM
    VolumeBF         float64   // Board Feet per BaseUOM (e.g., 5.33 for a 2x4x8)
    LengthInches     float64
    WidthInches      float64
    HeightInches     float64
    
    // Financials
    AverageCost      decimal.Decimal
    ReplacementCost  decimal.Decimal // Last AP Cost
}
```

### 2. Unit of Measure (UOM) Conversions
This allows the system to speak "Lumberjack" and "Accountant" simultaneously.

```go
type UOMFamily struct {
    ID   uuid.UUID
    Name string // "Lumber Lengths"
}

type UOM struct {
    ID         uuid.UUID
    FamilyID   uuid.UUID
    Name       string          // "Piece", "Bunk", "MBF" (Thousand Board Feet)
    Type       string          // "Reference", "Smaller", "Bigger"
    Ratio      decimal.Decimal // Multiplier relative to Reference UOM
    Rounding   float64         // Precision logic
}

/* 
Example Configuration for 2x4-8':
- Reference UOM: "Piece" (Ratio 1.0)
- Purchase UOM: "MBF" (Ratio 0.00533 - inversely, 1 MBF = 187.5 pieces)
- Stock UOM: "Bunk" (Ratio 294.0 - standard mill pack size)
*/
```

### 3. Stock Quantities (Multi-Dimensional)
Where is the stuff?

```go
type StockQuant struct {
    ID           uuid.UUID
    ProductID    uuid.UUID
    LocationID   uuid.UUID // "Yard A, Row 4"
    
    // The Quantities
    QtyOnInit    decimal.Decimal // Physical count
    QtyReserved  decimal.Decimal // Allocations (Sold but not Shipped)
    
    // Logistics Support
    InboundDate  time.Time       // FIFO/LIFO support
    LotNumber    string          // Traceability
}
```

### 4. Adjustments & Moves (Double-Entry Inventory)
Every move is a transaction. No "Magic Changes".

```go
type InventoryMove struct {
    ID           uuid.UUID
    Reference    string      // "PO-5001" or "Order-9022"
    MoveType     string      // "Receipt", "Delivery", "Adjustment", "Internal"
    
    FromLocID    uuid.UUID
    ToLocID      uuid.UUID
    
    State        string      // "Draft", "Done"
    Date         time.Time
}

type InventoryMoveLine struct {
    MoveID    uuid.UUID
    ProductID uuid.UUID
    Qty       decimal.Decimal
    UOMID     uuid.UUID      // The unit used for THIS move (e.g. moved 1 "Bunk")
}
```

### 5. "Reman" (Remanufacturing) Model
Turning one product into another.

```go
type RemanOrder struct {
    ID          uuid.UUID
    WorkCenter  string      // "Saw Shop"
    
    // Input
    SourceProduct uuid.UUID // "2x12-16' Douglas Fir"
    SourceQty     decimal.Decimal
    
    // Output
    OutputProduct uuid.UUID // "2x12 Stair Tread (48 inch)"
    OutputQty     decimal.Decimal
    
    // Waste/Shrinkage
    WasteFactor   decimal.Decimal // Sawdust loss
}
```
