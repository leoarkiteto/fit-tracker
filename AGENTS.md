# FitTracker

This is a **fitness activity hub** monorepo: workout tracking, bioimpedance, water intake, and AI-powered planning.

## Monorepo

- **Nx** workspace (`apps/`, `libs/`)
- **npm** workspaces for JS/TS; **dotnet** for API
- Each app has its own `AGENTS.md` — see `apps/<app>/AGENTS.md`

## Commands

| Target | Command |
|--------|---------|
| API | `npm run api` or `cd apps/api/FitTracker.Api && dotnet run --urls="http://0.0.0.0:5000"` |
| PWA | `npm run pwa` or `cd apps/pwa && npm run dev` |
| Mobile | `npm run mobile` or `cd apps/mobile && npx expo start` |

## Conventions

- [.cursor/rules/architecture.md](.cursor/rules/architecture.md) — Monorepo structure and shared libs
- [.cursor/rules/typescript.md](.cursor/rules/typescript.md) — TypeScript (PWA, mobile)
- [.cursor/rules/api-conventions.md](.cursor/rules/api-conventions.md) — API design patterns
