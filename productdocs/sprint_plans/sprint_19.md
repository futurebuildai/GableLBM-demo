# Sprint 19: Sovereign Product Configurator

**Goal**: Build a rules-based engine for custom millwork and lumber specs.

**Duration**: 1 Week
**Focus**: Knowledge Base, Configuration Logic, AI Vision Prototype.

## 1. Context
Incumbents like DMSi have specialized "Millwork Configurators." GableLBM needs a sovereign alternative that allows dealers to define their own valid combinations of Species, Grade, and Treatment without vendor intervention. This sprint also prototypes Vision AI for blueprint verification.

## 2. Objectives

### 2.1 Configuration Engine
- [ ] **Knowledge Base Schema**:
    - [ ] Define `ConfiguratorRules` table (Attribute dependency mapping).
    - [ ] Support for "Attribute Conflict" error handling.
- [ ] **Lumber Rule Matrix**:
    - [ ] Pre-load standard industry rules (e.g., "SYP" must be "Treatable", "Douglas Fir" is "Structural").

### 2.2 AI Vision Integration (Prototype)
- [ ] **Blueprint Check**:
    - [ ] Prototype service that scans a blueprint spec (PDF) and flags configurator selections that don't match.
    - [ ] Example: "Blueprint says 10' stud, user selected 8' in configurator."

### 2.3 UI/UX
- [ ] **Interactive Configurator Wizard**:
    - [ ] Stepper-based UI for selecting product attributes.
    - [ ] Non-stock item generation based on final configuration.
    - [ ] Visual "Preview" of the configured item.

## 3. Success Criteria
- User can configure a non-stock door in < 5 steps.
- System prevents invalid lumber combinations (e.g., grade that doesn't exist for the species).
- AI vision prototype can correctly identify a dimension on a sample PDF.
