# Feature Specification: PWA Activity Calendar

**Feature Branch**: `001-pwa-activity-calendar`  
**Created**: 2025-02-18  
**Status**: Draft  
**Input**: User description: "Implement activity calendar on PWA matching mobile: monthly view, completed/today/future states (green/blue/grey), day hover tooltip with activities, consistent style with mobile."

## Clarifications

### Session 2025-02-18

- Q: When should a day be shown as "completed" (green)? → A: At least one completed activity on that day is enough for the day to show as green (Option B).
- Q: When calendar data is loading or fails to load, what should the user see? → A: Show a loading state (e.g. skeleton or spinner) while fetching; on failure show an error message with retry (Option B).
- Q: On touch devices (no hover), how should the day tooltip open? → A: First tap on a day opens the tooltip; second tap on the same day or tap outside closes it (Option A).
- Q: Should month navigation be limited to a certain range? → A: No limit for this feature; allow any past and future month (Option B).
- Q: In which language should the calendar labels be shown (month name, day abbreviations)? → A: Always Portuguese (e.g. "Fevereiro", "Dom", "Seg") to match the mobile app (Option A).

## User Scenarios & Testing *(mandatory)*

### User Story 1 - View Monthly Activity Calendar (Priority: P1)

As a user, I can open the PWA and see a monthly activity calendar that shows the current month with all my activities. Each day is visually distinct based on status: days with completed activities appear in green, today is highlighted in blue, and past or future days without completion use grey. The calendar matches the layout and style of the mobile app (section title "Calendário de Treinos", day-of-week headers, numbered grid) so the experience is consistent across devices.

**Why this priority**: This is the core value—seeing the full month at a glance with clear status.

**Independent Test**: Open the PWA home (or calendar page), confirm the current month is shown with correct day labels and that at least one day is clearly "today" (blue), completed days (if any) are green, and other days are grey.

**Acceptance Scenarios**:

1. **Given** I am on the PWA activity calendar view, **When** the page loads, **Then** I see the current month with a grid of days, day-of-week headers (e.g. Dom, Seg, Ter…), and the current day indicated with a blue highlight.
2. **Given** the calendar is displayed, **When** a day has at least one completed activity (a day is "completed" when it has one or more completed activities), **Then** that day is shown with a green indicator (e.g. green circle or background) so I can tell it apart from other days.
3. **Given** the calendar is displayed, **When** a day is in the past or future and has no completed activity, **Then** that day is shown in grey so it is visually distinct from "today" and "completed".
4. **Given** I am viewing the calendar, **When** I look at the overall layout and styling, **Then** it matches the mobile app (same section title, similar spacing, same color semantics: green = completed, blue = today, grey = default).

---

### User Story 2 - Navigate Between Months (Priority: P2)

As a user, I can move to the previous or next month using controls on the calendar (e.g. left/right arrows) so I can review or plan activities in other months. The displayed month and year are clearly visible (e.g. "Fevereiro 2026"), and the same visual rules (green/blue/grey) apply for the selected month.

**Why this priority**: Essential for planning and review but depends on the calendar being visible first.

**Independent Test**: Use the month navigation to go to previous and next month; confirm the month label updates and day states (completed/today/grey) are correct for that month.

**Acceptance Scenarios**:

1. **Given** I am on the activity calendar, **When** I use the "previous month" control, **Then** the view shows the previous month and its correct month/year label.
2. **Given** I am on the activity calendar, **When** I use the "next month" control, **Then** the view shows the next month and its correct month/year label.
3. **Given** I have navigated to another month, **When** the new month is displayed, **Then** completed days are green, today (if in that month) is blue, and other days are grey, consistent with User Story 1.

---

### User Story 3 - See Day Activities on Hover (Priority: P3)

As a user, I can hover over any day in the calendar and see a tooltip that lists the activities scheduled or completed for that day. This lets me quickly check what is planned or done without leaving the calendar view.

**Why this priority**: Improves usability after the calendar and navigation are in place.

**Independent Test**: Hover over a day that has activities; confirm a tooltip appears with the list of activities for that day. Hover over a day with no activities; confirm the tooltip reflects that (e.g. "No activities" or empty list).

**Acceptance Scenarios**:

1. **Given** the calendar is displayed, **When** I hover over a day that has one or more activities, **Then** a tooltip appears showing the list of activities for that day (e.g. names or short descriptions).
2. **Given** the calendar is displayed, **When** I hover over a day that has no activities, **Then** the tooltip indicates there are no activities (or shows an empty list) so the behavior is clear.
3. **Given** a tooltip is visible, **When** I move the pointer away from the day, **Then** the tooltip is dismissed so it does not obscure the calendar.

---

### Edge Cases

- What happens when the user has no activities in the current month? The calendar still shows the month with "today" in blue and all other days in grey; no green days.
- What happens when the user navigates to a month that has no data yet? While data is loading, the user sees a loading state (e.g. skeleton or spinner); once loaded, the calendar shows the month grid with correct day states. Empty or loading states must not break the layout.
- What happens when calendar data fails to load (e.g. network or server error)? The user sees an error message and a retry control so they can try again without leaving the page.
- How does the system behave when "today" falls in a month the user has navigated away from? When viewing another month, "today" is still indicated in blue only when that month contains the current date; otherwise no blue highlight in that month.
- How is touch or non-hover interaction handled on devices without hover (e.g. touch)? On touch devices, first tap on a day opens the tooltip; second tap on the same day or tap outside closes it. Keyboard users get the same day activity list via focus (e.g. focus on day shows tooltip; blur or move focus dismisses it).

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The PWA MUST display a monthly activity calendar showing all days of the selected month in a grid with day-of-week headers.
- **FR-002**: The system MUST distinguish three visual states for calendar days: completed (green; a day is completed when it has at least one completed activity), today (blue), and other days (grey), applied consistently across the calendar.
- **FR-003**: The system MUST show the current month by default when the user opens the activity calendar.
- **FR-004**: The user MUST be able to navigate to the previous and next month with no range limit (any past or future month); the displayed month and year MUST be clearly visible.
- **FR-005**: When the user hovers over a calendar day (or taps it on touch devices; first tap opens, second tap on same day or tap outside closes), the system MUST show a tooltip listing the activities for that day (or indicate none). The empty-tooltip message MUST be in Portuguese (e.g. "Nenhuma atividade").
- **FR-006**: The calendar layout, section title, and color semantics (green/blue/grey) MUST align with the mobile app so that the experience is consistent across PWA and mobile. Calendar labels (month name, day-of-week abbreviations such as Dom, Seg, Ter) MUST be in Portuguese to match the mobile app.
- **FR-007**: The system MUST show activities that the user has completed and activities that are scheduled (past, today, or future) within the calendar data available for the selected month.
- **FR-008**: While calendar data is loading, the system MUST show a loading state (e.g. skeleton or spinner). If loading fails, the system MUST show an error message and a retry control so the user can try again.

### Key Entities

- **Activity**: A single planned or completed workout/activity. Has a date, completion status, and information that can be shown in the tooltip (e.g. name or type).
- **Calendar day**: A date in the month grid; has zero or more activities and a derived state (completed when at least one activity that day is completed, today, or default/future) used for visual styling.
- **Month view**: The currently displayed month and year; drives which calendar days and activities are shown.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: A user can open the PWA and identify which day is "today," which days have completed activities, and which days are upcoming or past, within one glance at the calendar.
- **SC-002**: A user can change the visible month and see the correct days and states for that month without errors or broken layout.
- **SC-003**: A user can discover the list of activities for any day by hovering (or equivalent interaction) and see the correct list or an explicit "no activities" indication.
- **SC-004**: Visual design of the PWA activity calendar is consistent with the mobile app (section title, color meaning, and overall structure) so that users switching between PWA and mobile recognize the same patterns.

## Assumptions

- "Activities" are the same concept as in the mobile app (e.g. workouts / "treinos"); data is available from the same or a compatible source for the PWA.
- The PWA has a designated page or area (e.g. home or a dedicated calendar route) where this calendar is shown; exact placement follows existing PWA navigation.
- Month navigation is not limited to a fixed range; users can navigate to any past or future month.
- Calendar labels (month names, day abbreviations) are always in Portuguese to match the mobile app.
- Color semantics: green = completed, blue = today, grey = default (past/future) as specified by the user; this may differ from mobile’s purple for "today" by product choice for the PWA.
- Tooltip content shows a list of activity names or short descriptions; exact fields are determined by existing activity data model.
- Accessibility: Touch: first tap opens tooltip, second tap on same day or tap outside closes it. Keyboard: focus on a day shows the tooltip; blur or moving focus dismisses it.
