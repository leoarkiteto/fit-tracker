# Tasks: PWA Activity Calendar

**Input**: Design documents from `specs/001-pwa-activity-calendar/`  
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/

**Tests**: Not requested in the feature specification; no test tasks included.

**Organization**: Tasks are grouped by user story so each story can be implemented and tested independently.

## Format: `[ID] [P?] [Story?] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: User story (US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

- PWA: `apps/pwa/src/` (App Router: `app/(app)/dashboard/`, components: `components/ActivityCalendar/`, lib: `lib/`)
- API: `apps/api/FitTracker.Api/Features/Workouts/`
- Shared libs: `libs/api-client/`, `libs/types/`

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Ensure directory structure for the new calendar feature

- [x] T001 Ensure directory structure for ActivityCalendar: create `apps/pwa/src/components/ActivityCalendar/` and ensure `apps/pwa/src/lib/` exists for calendar-utils

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core calendar logic and optional API support that all user stories depend on

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [x] T002 Implement calendar-utils in `apps/pwa/src/lib/calendar-utils.ts`: buildMonthView(year, month, workouts, completedRecords) returning MonthView with CalendarDay[] (date, state: completed|today|default, activities), Portuguese month name and day-of-week labels (Dom–Sáb) per data-model.md and research.md

- [ ] T003 [P] (Optional) Add optional query params `from` and `to` (ISO date) to GET completed-workouts in `apps/api/FitTracker.Api/Features/Workouts/CompletedWorkoutEndpoints.cs`; add getCompletedWorkoutsInRange(profileId, from, to) in `libs/api-client/src/index.ts` and expose in API_ENDPOINTS/config if needed; run `npm run pwa:sync-libs` from root after lib changes

**Checkpoint**: Foundation ready — calendar-utils can derive MonthView; optional API supports month-range loading for navigation

---

## Phase 3: User Story 1 - View Monthly Activity Calendar (Priority: P1) 🎯 MVP

**Goal**: User sees the current month with a grid, day-of-week headers (Portuguese), and day states: green (at least one completed), blue (today), grey (default). Section title "Calendário de Treinos". Loading state and error with retry.

**Independent Test**: Open dashboard; confirm current month with Dom–Sáb headers, one day as "today" (blue), completed days green, others grey; loading spinner then content; on network error, error message and retry.

### Implementation for User Story 1

- [x] T004 [US1] Create CalendarDay component in `apps/pwa/src/components/ActivityCalendar/CalendarDay.tsx`: accept date, state (completed|today|default), and render a single day cell with Tailwind styling (green for completed, blue for today, grey for default) and day number; no tooltip yet

- [x] T005 [US1] Create ActivityCalendar component in `apps/pwa/src/components/ActivityCalendar/ActivityCalendar.tsx`: accept workouts and completedRecords from parent (dashboard fetches and passes data); use calendar-utils buildMonthView for current month; render section title "Calendário de Treinos", day-of-week headers (Portuguese), grid of CalendarDay; show loading skeleton/spinner while data loads and error message with retry button on load failure (FR-008)

- [x] T006 [US1] Add "Calendário de Treinos" section to dashboard in `apps/pwa/src/app/(app)/dashboard/page.tsx`: fetch workouts and completed-workouts (or range for current month if T003 done), pass data to ActivityCalendar; place section above "Today's Workouts" and below stats/water; handle loading and error state for calendar data

**Checkpoint**: User Story 1 is fully functional; calendar visible on dashboard with correct states and loading/error behavior

---

## Phase 4: User Story 2 - Navigate Between Months (Priority: P2)

**Goal**: User can change the visible month with prev/next controls; month and year label (e.g. "Fevereiro 2026") is visible; same green/blue/grey rules apply for the selected month.

**Independent Test**: Click previous/next month; confirm label updates and day states are correct for that month (including "today" in blue when viewing current month).

### Implementation for User Story 2

- [x] T007 [US2] Add month navigation to ActivityCalendar in `apps/pwa/src/components/ActivityCalendar/ActivityCalendar.tsx`: state for (year, month); left/right arrow controls; display month and year label in Portuguese (e.g. "Fevereiro 2026"); same grid and day states for selected month

- [x] T008 [US2] Wire month change to data in `apps/pwa/src/components/ActivityCalendar/ActivityCalendar.tsx`: when (year, month) changes, refetch completed-workouts for that month range (use from/to if T003 done, else filter client-side) and recompute MonthView so any past/future month shows correct completed/today/grey states

**Checkpoint**: User Stories 1 and 2 work; user can navigate months and see correct labels and states

---

## Phase 5: User Story 3 - See Day Activities on Hover (Priority: P3)

**Goal**: Hover (or first tap on touch) shows a tooltip with the list of activities for that day; second tap on same day or tap outside closes it; keyboard focus shows tooltip, blur dismisses. Empty day shows "Nenhuma atividade" or equivalent.

**Independent Test**: Hover/tap a day with activities → tooltip shows names; tap day with no activities → tooltip shows no activities; dismiss by pointer leave, second tap, or tap outside; keyboard focus on day shows tooltip.

### Implementation for User Story 3

- [x] T009 [US3] Create DayTooltip component in `apps/pwa/src/components/ActivityCalendar/DayTooltip.tsx`: accept list of activity names (or { name, workoutId }[]); render list or "Nenhuma atividade" when empty; style consistently with PWA (Tailwind)

- [x] T010 [US3] Add tooltip behavior in `apps/pwa/src/components/ActivityCalendar/ActivityCalendar.tsx` and `CalendarDay.tsx`: activeDay state (date or null); hover on day sets activeDay, mouse leave clears; touch: first tap sets activeDay, second tap on same day or click outside clears (use ref + outside-click or escape key); render DayTooltip for activeDay with that day's activities from MonthView

- [x] T011 [US3] Add keyboard support in `apps/pwa/src/components/ActivityCalendar/CalendarDay.tsx`: make day cell focusable (tabIndex), on focus trigger callback to set activeDay so tooltip shows; ActivityCalendar clears activeDay on blur or when focus moves to another day

**Checkpoint**: All three user stories are complete; calendar has month view, navigation, and day tooltip with hover/touch/keyboard

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Validation and clean exports

- [x] T012 Run quickstart validation per `specs/001-pwa-activity-calendar/quickstart.md`: start API and PWA, sign in, open dashboard, verify "Calendário de Treinos" section, current month, today blue, completed green, month navigation, tooltip on hover/tap, loading and error+retry states; confirm behavior satisfies success criteria SC-001–SC-004 (glance identification, month change, tooltip discovery, visual consistency with mobile)

- [x] T013 [P] Add `apps/pwa/src/components/ActivityCalendar/index.ts` exporting ActivityCalendar (and CalendarDay, DayTooltip if used by other modules) for clean imports from dashboard

---

## Dependencies & Execution Order

### Phase Dependencies

- **Phase 1 (Setup)**: No dependencies — start immediately
- **Phase 2 (Foundational)**: Depends on Phase 1 — BLOCKS all user stories
- **Phase 3 (US1)**: Depends on Phase 2 — MVP
- **Phase 4 (US2)**: Depends on Phase 3 — extends ActivityCalendar
- **Phase 5 (US3)**: Depends on Phase 3 (and Phase 4 for full nav); tooltip can build on existing ActivityCalendar/CalendarDay
- **Phase 6 (Polish)**: Depends on Phase 5 (or at least Phase 3 for minimal validation)

### User Story Dependencies

- **US1 (P1)**: After Foundational only — no dependency on US2/US3
- **US2 (P2)**: Builds on US1 (same ActivityCalendar component)
- **US3 (P3)**: Builds on US1 (CalendarDay + ActivityCalendar); independent of US2 for tooltip-only

### Within Each User Story

- US1: T004 (CalendarDay) before T005 (ActivityCalendar); T005 before T006 (dashboard)
- US2: T007 (nav UI) then T008 (data for month)
- US3: T009 (DayTooltip) then T010 (tooltip behavior) then T011 (keyboard)

### Parallel Opportunities

- T003 [P] can run in parallel with T002 (different codebases: API vs PWA lib)
- T004 and T005 are sequential (T005 uses CalendarDay)
- T012 and T013 [P] can run in parallel in Phase 6

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1 (T001)
2. Complete Phase 2 (T002; T003 optional)
3. Complete Phase 3 (T004, T005, T006)
4. **STOP and VALIDATE**: Open dashboard, confirm calendar with current month, today blue, completed green, loading/error+retry
5. Deploy or demo if ready

### Incremental Delivery

1. Phase 1 + 2 → foundation ready
2. Phase 3 → MVP (calendar on dashboard)
3. Phase 4 → month navigation
4. Phase 5 → day tooltip (hover/touch/keyboard)
5. Phase 6 → quickstart check and exports

### Parallel Team Strategy

- One developer: T001 → T002 (and optionally T003) → T004 → T005 → T006 → T007 → T008 → T009 → T010 → T011 → T012, T013
- Two developers: After T002, Dev A does T004–T006 (US1), Dev B does T003 then T007–T008 (US2); then one does T009–T011 (US3)

---

## Notes

- [P] = parallelizable where noted; [USn] = maps to spec user story for traceability
- Each user story is independently testable per Independent Test in spec
- No test tasks: spec did not request TDD or explicit test tasks
- Optional T003: if skipped, calendar still works by fetching all completed and filtering by month client-side; add T003 when optimizing for large history
