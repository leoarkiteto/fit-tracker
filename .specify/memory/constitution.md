<!--
Sync Impact Report:
- Version change: (placeholder) → 1.0.0
- Modified principles: N/A (initial fill from template)
- Added sections: Core Principles (3), Additional Constraints, Development Workflow, Governance
- Removed sections: Principle 4 & 5 slots (template had 5; project uses 3)
- Templates: plan-template.md ✅ (Constitution Check references gates from constitution); spec-template.md ✅ (no mandatory sections changed); tasks-template.md ✅ (no new task types); checklist-template.md ✅ (no principle-specific changes); commands/*.md — not present
- Follow-up TODOs: None
-->

# FitTracker Constitution

## Core Principles

### I. Clean Code

Code MUST follow Clean Code practices: meaningful names, small functions, single responsibility, minimal duplication, and clear structure. Readability and maintainability are non-negotiable. Rationale: reduces defects and onboarding cost while keeping the codebase understandable for current and future contributors.

### II. Design Patterns When Necessary

Apply design patterns only when they clearly solve a problem or improve structure; avoid pattern-heavy or speculative abstraction. Prefer simple, direct solutions; introduce patterns when duplication, complexity, or domain rules justify them. Rationale: balances consistency and reuse without over-engineering.

### III. Simplicity, Accessibility & Speed

The project MUST stay simple (easy to reason about and change), accessible (approachable for contributors and users), and fast (responsive UX and performant execution). YAGNI applies: do not add scope or infrastructure before it is needed. Rationale: keeps FitTracker maintainable and pleasant to use and contribute to.

## Additional Constraints

- **Stack**: Frontend (React Native/Expo) and Backend (ASP.NET Core) as documented in README and AGENTS.md; no new major technologies without justification.
- **Structure**: Features live under `backend/Features/` and `frontend/`; shared rules in `.cursor/rules/` and `AGENTS.md` apply.
- **Complexity**: Any new dependency, pattern, or architectural layer MUST be justified against the principles above; prefer in-repo clarity over generic frameworks.

## Development Workflow

- **Reviews**: PRs and reviews MUST verify alignment with Clean Code, justified use of patterns, and no unnecessary complexity.
- **Refactoring**: Improve naming, structure, and duplication as part of feature work or dedicated refactors; avoid large “rewrite” PRs unless planned.
- **Guidance**: Use `AGENTS.md`, `frontend/AGENTS.md`, `backend/AGENTS.md`, and `.cursor/rules/` for day-to-day development and onboarding.

## Governance

This constitution overrides conflicting local practices. Amendments require: (1) documented rationale, (2) update of this file and version/date, (3) propagation to plan/spec/task templates and README if scope or principles change. All PRs and reviews MUST check compliance; exceptions (e.g. temporary technical debt) MUST be documented and tracked. Use the AGENTS.md files and `.cursor/rules/` for runtime development guidance.

**Version**: 1.0.0 | **Ratified**: 2025-02-18 | **Last Amended**: 2025-02-18
