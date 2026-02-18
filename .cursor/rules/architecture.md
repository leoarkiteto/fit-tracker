# FitTracker Architecture

## Monorepo Layout

```
apps/
├── api/          → ASP.NET Core (FitTracker.Api)
├── pwa/          → Next.js + Tailwind PWA
└── mobile/       → React Native (Expo)

libs/
├── types/        → @fittracker/types — shared TS interfaces
└── api-client/   → @fittracker/api-client — platform-agnostic API layer
```

## Shared Libraries

- **@fittracker/types**: Workout, Exercise, UserProfile, WaterIntake, etc.
- **@fittracker/api-client**: REST client used by PWA and mobile. Config in `libs/api-client/src/config.ts`.

## Sync (PWA)

PWA installs independently. From monorepo root: `npm run pwa:sync-libs` (copies types + api-client into `apps/pwa/src/lib/`).
