# Data Model: PWA Activity Calendar

**Feature**: 001-pwa-activity-calendar  
**Date**: 2025-02-18

This feature reuses existing entities from the product; it adds view-model concepts for the calendar UI only. No new persistence entities.

---

## Existing entities (from @fittracker/types and API)

### Workout

- **id**: string (UUID)
- **name**: string
- **description**: string | undefined
- **goal**: WorkoutGoal
- **days**: DayOfWeek[] — which weekdays this workout is scheduled
- **exercises**: Exercise[]
- **createdAt**, **updatedAt**: string (ISO)

Used to derive "scheduled" activities for a given calendar date (date’s weekday is in `days`).

### CompletedWorkoutRecord

- **id**: string (UUID)
- **workoutId**: string
- **completedAt**: string (ISO date-time)
- **durationSeconds**: number

Used to derive "completed" state: a calendar day is completed (green) when at least one record has `completedAt` on that date.

### Workout (for tooltip)

Tooltip shows activity names (and optionally short descriptions) for a day. Source: workouts that are either (a) scheduled that day (by `days`) or (b) completed that day (via CompletedWorkoutRecord). Display: workout **name** at minimum.

---

## View-model concepts (UI only; not persisted)

### CalendarDay

- **date**: string (ISO date, e.g. `YYYY-MM-DD`)
- **state**: `'completed' | 'today' | 'default'`
  - **completed**: at least one CompletedWorkoutRecord with completedAt on this date
  - **today**: date equals current local date
  - **default**: otherwise (grey)
- **activities**: { name: string; workoutId: string }[] — for tooltip: scheduled workouts for this weekday plus completed workouts on this date (by name, deduplicated as needed)

Derived per month from: list of Workout, list of CompletedWorkoutRecord, and the month’s date range.

### MonthView

- **year**: number
- **month**: number (1–12)
- **days**: CalendarDay[] — one entry per day in that month (including leading/trailing empty cells for grid alignment if needed by the UI)
- **monthLabel**: string — Portuguese month name (e.g. "Fevereiro")
- **dayOfWeekLabels**: string[] — Portuguese short labels (Dom, Seg, Ter, Qua, Qui, Sex, Sáb)

Used to drive the calendar grid and navigation.

---

## Validation rules (from spec)

- A day is **completed** if and only if there is at least one completed workout record for that date.
- **Today** is the current local calendar date (one per view).
- Activities shown in the tooltip for a day: all workouts that are either scheduled for that weekday (from `days`) or completed on that date; display at least the workout name.

---

## State transitions

No domain state machines. UI state only:

- **Month navigation**: change `MonthView.year` / `MonthView.month`; recompute `days` and labels.
- **Tooltip**: show/hide by active day (hover or tap/focus); content is the `activities` array for that `CalendarDay`.
