# Sprint 11: Millwork Configurator & Special Orders

## Goal
Enable the sale of complex, configurable millwork products (Doors/Windows) and handle special order workflows.

## Objectives

### 1. Millwork Configuration Engine
- **Configurable SKUs**: Support "Parent SKU" with "Child Options" (e.g., Door mechanism + Frame + Glass).
- **Dynamic Pricing**: Price calculated based on selected attributes, not just a static list price.
- **Configurator UI**: A step-by-step wizard for selecting millwork options.

### 2. Special Orders
- **Non-Stock Items**: Workflow to sell items not in `inventory` (Special Order placeholders).
- **PO Generation**: Auto-create a Purchase Order for the vendor when a Special Order is sold.
- **Link**: Link SO (Sales Order) line item to PO (Purchase Order) line item.

## Technical Constraints
- **Data Structure**: Need a flexible JSONB column or EAV pattern for Millwork Attributes (Width, Height, Swing, Jamb Width).
- **Validation**: Ensure incompatible options cannot be selected (e.g., French Door + tiny width).
