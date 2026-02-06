# Contributing to GableLBM

We build with the **Antigravity** Standard. Zero laziness. L8 Quality.

## The Workflow
1.  **Sprints**: We work in strict sprints defined in `productdocs/sprint_plans/`.
2.  **Agent Execution**: Agents must follow the `.agent/workflows/sprint_execution.md` protocol.
3.  **Human Execution**: Humans should mimic this rigor—PRs required, tests required.

## The "Gate"
No code merges to `main` without passing the **L8 Production Readiness Gate**.
See: `productdocs/process/production_readiness_gate.md`.

## Tech Stack
*   **Backend**: Go 1.25+ (Standard Lib + Chi).
*   **Frontend**: React 19 + Vite + Shadcn/UI.
*   **Style**: Tailwind CSS ("Industrial Dark").

## Principles
*   **Headless**: Logic lives in hooks/services. UI is just a render layer.
*   **Safety**: Strong types. No `any`.
*   **Speed**: < 100ms interactions.
