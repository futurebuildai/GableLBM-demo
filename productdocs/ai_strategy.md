# AI Strategy Specification

## 1. Vision: "AI as a Co-Pilot, not a Gimmick"
In the LBM space, margins are thin and efficiency is everything. AI features must directly solve time-sinks.

## 2. Targeted Use Cases

### 2.1. "Smart Takeoff" (Vision + LLM)
**Problem:** Contractors send PDF blueprints (plan sets) and ask for a quote. Sales reps spend hours manually counting windows, doors, and linear feet of framing.
**Solution:**
- **Input:** PDF Blueprint.
- **Process:** Computer Vision (CV) segmenting the floorplan -> Identifying "Wall Types" -> LLM mapping walls to SKUs (2x4x92 5/8 studs).
- **Output:** Draft Quote in the system ready for human review.

### 2.2. "Inventory Detective" (Anomaly Detection)
**Problem:** Shrinkage (theft/loss) is huge in lumber yards.
**Solution:**
- **Process:** Analyze "Pick vs. Invoice" patterns and "Adjustment" history.
- **Alert:** "User 'Steve' adjusts inventory down for '2x4 cedar' 3x more often than the average yard guy."

### 2.3. "Conversational Commerce" (RAG)
**Problem:** Counter staff turnover is high; new hires don't know the obscure hardware.
**Solution:**
- **Process:** Vectorize product manuals, installation guides, and code books.
- **Interface:** Chatbot for staff.
- **Query:** "Customer needs a hanger for a double 1.75 LVL beam."
- **Answer:** "Recommend Simpson Strong-Tie HUS28-2. We have 40 in Stock."

### 2.4. "Route Optimization" (Constraint Solving)
**Problem:** Flatbeds have complex loading rules (heavy on bottom, delivery order LIFO).
**Solution:**
- **Process:** Constraint logic programming + ML based on past successful loads.
- **Output:** 3D Visualization of how to load the truck for the forklift driver.

### 2.5. "Governance Orchestrator" (Natural Language CRs)
**Problem:** Non-technical industry partners struggle to write clear technical specs for software updates.
**Solution:**
- **Process:** LLM-based translation of "Dealer needs" into Gherkin (Feature) files and architectural impact reports.
- **Outcome:** Reduces the friction of community-driven software development by 90%.

## 3. Technical Implementation
- **Vector Store:** `pgvector` (Keep it in the main Postgres DB for simplicity).
- **LLM Gateway:** Abstraction layer to swap providers (OpenAI / Anthropic / Local Llama).
- **Data Privacy:** Strict PII scrubbing before sending data to cloud LLMs.
