# Implementation Plan: PWA Activity Calendar

**Branch**: `001-pwa-activity-calendar` | **Date**: 2025-02-18 | **Spec**: [spec.md](./spec.md)  
**Input**: Feature specification from `specs/001-pwa-activity-calendar/spec.md`

## Summary

Add a monthly activity calendar to the PWA (dashboard) that matches the mobile "Calendário de Treinos": show the current month by default with day-of-week headers and a grid of days. Days are visually distinct: green when at least one activity is completed that day, blue for today, grey otherwise. Month navigation (prev/next) with no range limit. Hover (or tap on touch: first tap opens, second tap or tap outside closes) shows a tooltip with the list of activities for that day. Loading state and error-with-retry when data fails. Labels in Portuguese. Uses existing workouts and completed-workouts APIs; optional API support for date-range filtering to load calendar data efficiently.

## Technical Context

**Language/Version**: TypeScript (PWA); C# / .NET 10 (API if extending)  
**Primary Dependencies**: Next.js 14 (App Router), React 18, Tailwind CSS; @fittracker/types, @fittracker/api-client (synced into PWA)  
**Storage**: Existing API persistence (SQLite via EF Core); no new storage for this feature  
**Testing**: Jest / React Testing Library or Playwright for PWA (per workspace); API tests if new endpoint behavior added  
**Target Platform**: PWA (browser); responsive and touch-friendly  
**Project Type**: Web (PWA frontend); API remains vertical-slice  
**Performance Goals**: Calendar month loads with minimal delay; tooltip appears without perceptible lag  
**Constraints**: Same visual semantics as mobile (green/blue/grey); Portuguese labels; Tailwind-only styling  
**Scale/Scope**: Single new UI section on dashboard; optional API query params for completed-workouts by date range

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Status | Notes |
|-----------|--------|-------|
| I. Monorepo & Shared Contracts | Pass | Calendar consumes @fittracker/types and api-client; no new shared types required unless we add a calendar DTO in libs. |
| II. Vertical Slice Architecture (API) | Pass | No new feature slice required; optional query params on existing Completed Workouts slice. |
| III. Platform–App Alignment | Pass | PWA-only change; Tailwind for styling; sync libs if types/api-client change. |
| IV. AI Features as First-Class | N/A | No AI in this feature. |
| V. Observability & Conventions | Pass | Follow .cursor/rules (architecture, typescript, nextjs); error/retry UX satisfies user-facing feedback. |

**Post–Phase 1**: No change. Optional API query params do not alter shared contracts; data-model and contracts align with existing types and vertical slice.

## Project Structure

### Documentation (this feature)

```text
specs/001-pwa-activity-calendar/
├── plan.md              # This file
├── research.md          # Phase 0
├── data-model.md        # Phase 1
├── quickstart.md        # Phase 1
├── contracts/           # Phase 1 (API contract for optional date-range)
└── tasks.md             # Phase 2 (/speckit.tasks - not created by plan)
```

### Source Code (repository root)

```text
apps/pwa/
├── src/
│   ├── app/(app)/dashboard/
│   │   └── page.tsx           # Add calendar section (Calendário de Treinos)
│   ├── components/
│   │   └── (new) ActivityCalendar/
│   │       ├── ActivityCalendar.tsx    # Month grid, nav, day states, tooltip
│   │       ├── CalendarDay.tsx         # Single day cell (green/blue/grey, hover/tap)
│   │       └── DayTooltip.tsx          # Tooltip content: list of activities for day
│   └── lib/
│       └── (optional) calendar-utils.ts  # Month grid generation, PT labels, day state derivation

apps/api/FitTracker.Api/
└── Features/Workouts/
    └── CompletedWorkoutEndpoints.cs   # Optional: add from/to query params to GET
```

**Structure Decision**: Calendar is a new component tree under PWA; dashboard page gains a "Calendário de Treinos" section above or alongside "Today's Workouts". API change is limited to optional query parameters on existing GET completed-workouts.

## Complexity Tracking

No constitution violations. Leave table empty or omit if no justification needed.
