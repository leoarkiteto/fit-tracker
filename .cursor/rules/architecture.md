# FitTracker Architecture

## Layout

```
frontend/     → React Native (Expo); types and API client in src/types, src/services
backend/      → ASP.NET Core (FitTracker.Api); Vertical Slice in Features/
```

## Backend

- **Vertical Slice**: each feature under `Features/` (Auth, Profiles, Workouts, Bioimpedance, WaterIntake, AI).
- Shared data and infra in `Shared/`.

## Frontend

- Types and API layer live inside the app: `frontend/src/types/`, `frontend/src/services/`.
- No shared libs; the frontend is a single standalone project.
