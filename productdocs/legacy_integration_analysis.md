# Legacy Integration Analysis

## 1. Overview
We analyzed three major industry incumbents: **Epicor BisTrack**, **DMSi Agility**, and **ECI Spruce**. The goal was to identify critical "must-have" features for compatibility and gaps where GableLBM can innovate.

## 2. Key Findings by Vendor

### 2.1 ECI Spruce (SOAP/XML)
- **Strengths:** Robust management of "Jobs" and "Card on File" for contractors.
- **LBM Specifics:** Explicit support for `Tally` logic (`IsTally`, `TallyType`, `TallyItems`).
- **Gaps:** Antiquated SOAP interface is hard to integrate with modern web/mobile apps.

### 2.2 Epicor BisTrack (REST/OpenAPI)
- **Strengths:** Flexible `SmartView` feature allows read-only custom queries. Distinguishes between `Quantity UOM` and `Sell UOM`.
- **LBM Specifics:** "Selling UOM Barcodes" - allows barcodes to represent specific quantities (e.g., 1 Pallet vs 1 Each).
- **Gaps:** Complexity of the "SmartView" implementation often leads to performance issues in legacy deployments.

### 2.3 DMSi Agility (OpenAPI)
- **Strengths:** Massive breadth covering `Reman` (Remanufacturing) and `Dispatch`.
- **LBM Specifics:** Extensive `TallyString` support and `PickingTallyUOM`. Recognizes that the unit being picked may differ from the unit ordered.
- **Gaps:** Deeply nested JSON structures and "Data-Chunking" requirements suggest an older architecture adapted for the web.

## 3. Implications for GableLBM

### 3.1. Technical Requirements
- **Unit-Aware Barcode Engine:** Our scanning logic must handle barcodes that encode both ProductID and UOMID.
- **Shorthand Tally Parser:** We should implement an AI-powered or regex-based parser for "Yard Shorthand" (e.g., `4/4 RWL Walnut`) to match DMSi's `TallyString` capability.
- **Adaptable Pricing Engine:** Must support "Price Levels", "Price Groups", and "Contractor Class" to enable smooth migration from legacy systems.

### 3.2. Competitive Advantages
- **Unified API vs SOAP:** Replacing Spruce's XML with a clean, low-latency gRPC/REST API will immediately appeal to 3rd party developers (e.g., e-commerce providers).
- **Proactive Tally Verification:** Instead of just recording tally (as legacy systems do), GableLBM can use AI to *verify* if the entered tally physically fits on the truck/bunk being loaded.
- **"Import-First" Architecture:** We should build native "Mappers" for these three schemas. A dealer using BisTrack should be able to "sync" their data into GableLBM with one click to test-drive the UX.

## 4. Architectural Patterns
- **Adaptor Pattern:** Our `Inventory` and `Sales` modules should use an Adaptor pattern for external data sources.
- **Event-Driven Reconciliation:** When a legacy system is used alongside GableLBM (hybrid phase), we'll need an event bus to handle bi-directional state sync (e.g., `legacy.order.created` -> `gable.inventory.reserve`).
