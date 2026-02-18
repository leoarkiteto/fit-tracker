# API Contracts: PWA Activity Calendar

**Feature**: 001-pwa-activity-calendar  
**Date**: 2025-02-18

The calendar uses existing endpoints. One extension is specified below for efficient month-range loading.

---

## Existing endpoints used as-is

### GET /api/profiles/{profileId}/workouts

- **Purpose**: List all workouts for the profile (used to derive scheduled activities by weekday).
- **Response**: `200` — array of `WorkoutDto` (id, name, description, goal, days, exercises, createdAt, updatedAt, isCompleted, completedAt).
- **Auth**: Required (JWT).

### GET /api/profiles/{profileId}/completed-workouts

- **Purpose**: List completed workout records. Used to derive which calendar days have at least one completion (green state) and to build per-day activity lists for the tooltip.
- **Response**: `200` — array of `CompletedWorkoutDto` (id, workoutId, completedAt, durationSeconds).
- **Auth**: Required (JWT).

---

## Extended contract: optional date filter

### GET /api/profiles/{profileId}/completed-workouts?from={date}&to={date}

- **Purpose**: Same as GET without query params, but restrict results to records whose `completedAt` (date part) is within the inclusive range `[from, to]`. Enables the PWA to request only the visible month when loading the calendar for a given month.
- **Query parameters** (optional; both omitted = current behavior, return all):
  - **from**: ISO date string `YYYY-MM-DD` (inclusive start).
  - **to**: ISO date string `YYYY-MM-DD` (inclusive end).
- **Response**: `200` — array of `CompletedWorkoutDto` (unchanged shape).
- **Validation**: If `from` or `to` is present but invalid, return `400 Bad Request`. If only one of `from`/`to` is provided, treat the other as unbounded (e.g. only `from` = all records on or after that date).
- **Auth**: Required (JWT).

**Example**: `GET /api/profiles/{profileId}/completed-workouts?from=2026-02-01&to=2026-02-28` returns only completed workouts in February 2026.

---

## PWA usage

1. For the displayed month (year, month), compute `from` = first day of month, `to` = last day of month.
2. Call `GET .../completed-workouts?from=...&to=...` (or without params for MVP).
3. Call `GET .../workouts` once (or per profile) to get scheduled workouts.
4. Build `MonthView` and `CalendarDay[]` from workouts + completed records (see data-model.md).
