# Financial & Job Accounting Specification

## 1. Overview
LBM dealers operate heavily on "charge accounts" and "job costing". This module handles AR, AP, and the unique "Job Accounting" requirements of the construction industry.

## 2. Core Concepts

### 2.1. Customer Account
A master record for a business or individual.
```
Customer {
    ID              UUID
    Name            String
    PrimaryContact  String
    CreditLimit     Decimal
    CreditStatus    Enum    // "Good", "OnHold", "COD"
    DefaultPriceLevel UUID
    TaxExemptNumber String?
}
```

### 2.2. Job (Project)
A specific project under a Customer. Allows for:
- Per-job credit limits.
- Per-job pricing overrides.
- Construction lien tracking.

```
Job {
    ID              UUID
    CustomerID      UUID
    Name            String   // "Smith Residence - 123 Main St"
    Address         Address
    Status          Enum     // "Preliminary", "Open", "Closed"
    CreditLimit     Decimal
    
    // Lien Tracking
    StartDate       Date
    NoticeToOwnerDate Date?
    LienReleaseDate Date?
}
```

### 2.3. Invoice
```
Invoice {
    ID              UUID
    OrderID         UUID?    // Link to originating order
    CustomerID      UUID
    JobID           UUID?
    
    InvoiceDate     Date
    DueDate         Date
    Terms           String   // "Net 30", "2% 10 Net 30"
    
    Subtotal        Decimal
    TaxTotal        Decimal
    Total           Decimal
    
    AmountPaid      Decimal
    Balance         Decimal
    Status          Enum     // "Open", "PartiallyPaid", "Paid", "Void"
}
```

## 3. Accounts Receivable (AR) Workflow

### 3.1. Aging Buckets
Standard AR aging report.
| Bucket     | Days       |
|------------|------------|
| Current    | 0-30       |
| 30 Days    | 31-60      |
| 60 Days    | 61-90      |
| 90+ Days   | 91+        |

### 3.2. Payment Application
Payments can be applied:
1. **Auto-Apply:** Oldest invoice first (FIFO).
2. **Specific Invoice:** Customer designates which invoice(s).
3. **On Account:** Unapplied credit balance.

### 3.3. Finance Charges
Configurable per customer:
- Rate (e.g., 1.5% per month).
- Grace Period (e.g., 15 days past due).
- Minimum Charge (e.g., $5.00).

## 4. Construction Lien Tracking
Critical for LBM dealers who supply to job sites.

### 4.1. "Notice to Owner" (NTO)
- System prompts for NTO on first delivery to a new Job.
- Tracks Date Sent and Expiration.

### 4.2. Lien Release
- On final payment for a Job, system prompts for Lien Release generation.
- Conditional releases for partial payments.

## 5. Accounts Payable (AP) - Brief
While AR is the primary financial module, basic AP is needed for:
- Vendor invoice entry (matched to POs).
- Payment scheduling.
- Rebate tracking (common in LBM for volume buys).

## 6. Reporting
| Report Name             | Description                                      |
|-------------------------|--------------------------------------------------|
| AR Aging Summary        | Customer balances by aging bucket.               |
| AR Aging Detail         | Invoice-level detail by customer.                |
| Job Profitability       | Revenue - COGS by Job.                           |
| Lien Exposure Report    | Open invoices on jobs without lien release.      |
| Sales by Rep            | Sales totals/margins by salesperson.             |
