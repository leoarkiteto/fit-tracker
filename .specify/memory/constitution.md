<!--
  Sync Impact Report
  Version change: (none) → 1.0.0
  Rationale: Initial ratification; constitution was template-only with placeholders.
  Modified principles: N/A (all new)
  Added sections: Core Principles (5), Technology Stack & Constraints, Development Workflow, Governance
  Removed sections: None
  Templates: plan-template.md ✅ (Constitution Check / gates reference unchanged, generic)
             spec-template.md ✅ (no principle-specific mandatory sections added)
             tasks-template.md ✅ (task categorization compatible)
  Commands: .cursor/commands/*.md — no principle name references to update; constitution path correct.
  Follow-up TODOs: None
-->

# FitTracker Constitution

## Core Principles

### I. Monorepo & Shared Contracts

The project is an Nx monorepo with a single workspace. Apps (API, PWA, Mobile) and shared libraries live under `apps/` and `libs/`. Shared contracts and types are the single source of truth: `@fittracker/types` and `@fittracker/api-client` define interfaces and API usage consumed by PWA and mobile. New features that affect cross-app behavior MUST update shared libs and follow the sync discipline (e.g. `npm run pwa:sync-libs` after lib changes).

**Rationale**: Prevents drift between API, PWA, and mobile; ensures one place to change types and client behavior.

### II. Vertical Slice Architecture (API)

The backend API (ASP.NET Core) MUST be organized by vertical slices: each feature or use-case is a self-contained slice with its own request/response, handlers, and behavior. Slices MUST NOT be organized primarily by technical layer (e.g. one big “services” or “controllers” folder). New API capabilities MUST be added as new slices or extensions of existing slices with clear boundaries.

**Rationale**: Keeps features cohesive and easy to locate; reduces cross-cutting coupling and eases incremental delivery.

### III. Platform–App Alignment

PWA is the Progressive Web App (Next.js, styles with TailwindCSS). Mobile is the React Native (Expo) app. Each app MUST consume shared libraries for types and API access where applicable. PWA MUST run `npm run pwa:sync-libs` after changes to `libs/types` or `libs/api-client`. Styling in the PWA MUST use TailwindCSS; mobile follows its own design system and platform conventions.

**Rationale**: Clear ownership per platform and consistent use of shared contracts without mixing styling or platform concerns.

### IV. AI Features as First-Class Capabilities

AI-driven features (e.g. AI Planner that creates activity plans from user inputs) MUST have defined inputs and outputs, testable behavior, and documented limits. AI-generated content or recommendations MUST be validated where applicable, and fallback or degradation paths MUST be defined (e.g. when the AI is unavailable or returns invalid data). New AI features MUST be specified with clear success criteria and edge cases.

**Rationale**: Ensures AI adds predictable value and fails safely; keeps specs and tests aligned with AI behavior.

### V. Observability & Conventions

Structured logging and consistent error handling are required in the API and SHOULD be considered in PWA and mobile for critical flows. All work MUST follow the conventions in `.cursor/rules/` (architecture, TypeScript, API design). Breaking changes to shared contracts or public API MUST be documented and reflected in versioning or migration guidance.

**Rationale**: Enables debugging and operations; keeps code style and API design consistent across the monorepo.

## Technology Stack & Constraints

- **Monorepo**: Nx workspace; npm workspaces for JS/TS; .NET for API.
- **API**: ASP.NET Core, Vertical Slice Architecture; run via `npm run api` or dotnet from `apps/api/FitTracker.Api`.
- **PWA**: Next.js, TailwindCSS; run via `npm run pwa`; build via `npm run pwa:build`.
- **Mobile**: React Native (Expo); run via `npm run mobile`; platform targets via `npm run mobile:ios` / `mobile:android`.
- **Shared libs**: `libs/types` (@fittracker/types), `libs/api-client` (@fittracker/api-client); PWA sync after lib changes.

Compliance with this stack is expected unless an explicit exception is documented and justified (e.g. in a spec or architecture decision).

## Development Workflow

- Use root-level npm scripts and app-specific `AGENTS.md` (e.g. `apps/pwa/AGENTS.md`, `apps/mobile/AGENTS.md`) for run and build instructions.
- After changing `libs/types` or `libs/api-client`, run `npm run pwa:sync-libs` before relying on PWA with the new code.
- Code reviews and implementation plans MUST align with this constitution; violations MUST be justified in the plan’s Complexity Tracking (or equivalent) and should trigger a constitution amendment when they become the norm.

## Governance

- This constitution supersedes ad-hoc practices for the FitTracker monorepo. When in conflict, the constitution takes precedence.
- Amendments require updating `.specify/memory/constitution.md`, documenting the change, and updating the version and Last Amended date. Principle removals or redefinitions are MAJOR; new principles or sections are MINOR; clarifications and typos are PATCH.
- All PRs and feature plans MUST verify compliance with these principles. Complexity that contradicts a principle MUST be justified and, if sustained, considered for a constitution update.
- For day-to-day development, use the root `AGENTS.md` and `.cursor/rules/` (architecture.md, typescript.md, api-conventions.md) as the primary guidance.

**Version**: 1.0.0 | **Ratified**: 2025-02-18 | **Last Amended**: 2025-02-18
