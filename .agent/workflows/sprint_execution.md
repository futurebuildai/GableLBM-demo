---
description: Comprehensive Antigravity sprint execution workflow with L7/L8 quality gates
---

# Sprint Execution Workflow

**Trigger**: New Agent Thread for a Sprint.
**Context**: You are "Antigravity". You execute with extreme precision, modularity, and zero laziness.

## Phase 1: Context & Planning
1.  **Load State**:
    - Read `productdocs/master_prd.md`, `productdocs/design_system.md`.
    - Read `productdocs/sprint_plans/<target_sprint>.md`.
    - Read `task.md`.
2.  **Initialize Task List**:
    - Update `task.md` to break the sprint goal into granular, check-able engineering steps.
    - **Rule**: No task should be larger than 1 hour of work. Break it down.

## Phase 2: The Execution Loop (Per Task)
For each item in `task.md`:
1.  **Task Boundary**: Call `task_boundary` with specific objective.
2.  **Implementation**: Write the code.
    - *Constraint*: Follow "Headless" patterns. Use `libs/` for logic, `components/` for view.
    - *Constraint*: Apply "Industrial Dark" design tokens immediately. No unstyled HTML.
3.  **L7 Self-Review (Immediate)**:
    - Before marking done, ask: "Would a Google L7 Engineer accept this PR?"
    - Check for: Error handling, Types, Comments, Modular structure.
    - *Action*: Refactor immediately if subpar.
4.  **Verify**: Run `npm run build` or `go test` to ensure no regression.
5.  **Commit**: `git commit -m "feat: <description>"`

## Phase 3: The L8 Production Readiness Gate (Recursive)
**Trigger**: When all implementation tasks are checked.

1.  **Read Audit Standard**: Open `productdocs/process/production_readiness_gate.md`.
2.  **Execute Audit**: Systematically check the entire codebase against the 5 categories (Arch, SRE, Security, UX, Build).
3.  **Decision**:
    - **FAIL**: Use `task_boundary` ("Remediating L8 Audit Findings"). Fix the issues. **GOTO Step 2**.
    - **PASS**: Only proceed if you are 100% confident.

## Phase 4: State Update & Handoff
1.  **Documentation**:
    - Update `productdocs/roadmap.md` (Mark phase complete).
    - Create/Update `walkthrough.md` with screenshots/logs of the working feature.
2.  **Next Sprint Prep**:
    - Create placeholder `productdocs/sprint_plans/sprint_<next>.md`.
3.  **Final Push**: `git push origin <branch>`.
4.  **Notify User**: "Sprint Complete. Audit Passed 100%. Ready for Next Thread."
