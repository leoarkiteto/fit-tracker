# Research: PWA Activity Calendar

**Feature**: 001-pwa-activity-calendar  
**Date**: 2025-02-18

## 1. Calendar data loading (completed workouts + scheduled by month)

**Decision**: Use existing `workoutsApi.getAll(profileId)` and `completedWorkoutsApi.getAll(profileId)` for MVP. Add optional query parameters `from` and `to` (ISO date) to `GET /api/profiles/{profileId}/completed-workouts` so the PWA can request only the visible month range and avoid loading full history when navigating to arbitrary months.

**Rationale**: Client can derive "completed" days from completed-workouts (filter by date); "scheduled" for a day can be derived from workouts whose `days` include that day of week. For a single month, fetching all completed workouts is acceptable for typical usage; for "any past/future month" (no limit), a date-range API reduces payload and improves responsiveness.

**Alternatives considered**: (a) Client-only: always fetch all completed workouts — rejected for unbounded month navigation. (b) New dedicated "calendar" endpoint returning pre-aggregated month data — rejected for scope; extending existing endpoint with query params keeps the API surface smaller.

---

## 2. Portuguese month and day-of-week labels

**Decision**: Use a fixed Portuguese mapping for month names and short day labels (Dom, Seg, Ter, Qua, Qui, Sex, Sáb). Implement via a small constant map or Intl with locale `pt-BR` (e.g. `Intl.DateTimeFormat('pt-BR', { month: 'long' })` and weekday `short`), ensuring output matches mobile (e.g. "Fevereiro", "Dom", "Seg").

**Rationale**: Spec and clarifications require labels in Portuguese to match the mobile app; Intl with pt-BR is standard and avoids hand-maintained month names while guaranteeing consistency.

**Alternatives considered**: (a) Hardcoded arrays for months and weekdays — acceptable but more maintenance. (b) Browser locale — rejected; spec requires Portuguese always.

---

## 3. Tooltip and touch/keyboard behavior

**Decision**: Use a single tooltip component that shows on hover (desktop), and on first tap (touch) with dismiss on second tap on same day or tap outside. For keyboard: focusable day cells with focus showing the same tooltip and blur/move focus dismissing it. Use existing React patterns (state for "active day", optional ref for outside-click or escape key).

**Rationale**: Aligns with spec and clarifications; avoids new dependencies; consistent with common calendar UX patterns.

**Alternatives considered**: (a) Hover-only — rejected (accessibility and touch). (b) Third-party calendar library — rejected; custom component keeps bundle small and matches mobile layout exactly.

---

## 4. Where to place the calendar on the PWA

**Decision**: Add the "Calendário de Treinos" section to the existing dashboard page (`apps/pwa/src/app/(app)/dashboard/page.tsx`), above "Today's Workouts" and below the stats/water blocks, so it matches the mobile home flow (calendar then quick actions / today).

**Rationale**: Product vision and AGENTS.md describe dashboard as the main hub; mobile screenshot places the workout calendar prominently on the home screen; no separate route required for this feature.

**Alternatives considered**: Dedicated `/calendar` route — deferred; can be added later if needed without changing the spec.
