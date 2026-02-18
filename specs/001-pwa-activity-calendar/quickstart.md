# Quickstart: PWA Activity Calendar

**Feature**: 001-pwa-activity-calendar  
**Date**: 2025-02-18

## Prerequisites

- Node.js and npm (workspace root)
- .NET 10 SDK (for API)

## Run the stack

1. **Start the API** (from repo root):
   ```bash
   npm run api
   ```
   API runs at `http://localhost:5000` (or `http://0.0.0.0:5000`).

2. **Start the PWA** (from repo root):
   ```bash
   npm run pwa
   ```
   PWA runs at `http://localhost:3000` (or the port shown in the terminal).

3. **Open the app**: Go to `http://localhost:3000`, sign in, and navigate to the dashboard. The activity calendar section **"Calendário de Treinos"** appears on the dashboard (once implemented).

## Sync shared libs (if you change libs/types or libs/api-client)

From repo root:
```bash
npm run pwa:sync-libs
```

## Verify calendar behavior (after implementation)

- **Default view**: Current month with day-of-week headers (Dom–Sáb) and numbered days; today highlighted in blue.
- **Completed days**: Days with at least one completed workout show in green.
- **Other days**: Grey.
- **Month navigation**: Use arrows to change month; label shows Portuguese month and year (e.g. "Fevereiro 2026").
- **Tooltip**: Hover (or tap on touch) a day to see the list of activities for that day; tap again or outside to close.
- **Loading**: Skeleton or spinner while data loads.
- **Error**: On load failure, error message and retry control.

## API (optional date range)

If the API is extended with `from`/`to` query params on `GET .../completed-workouts`, the PWA can call:

```http
GET /api/profiles/{profileId}/completed-workouts?from=2026-02-01&to=2026-02-28
```

to load only February 2026 completed workouts for the calendar.
