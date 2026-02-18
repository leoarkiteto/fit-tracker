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
| API build | `npm run api:build` |
| PWA | `npm run pwa` or `cd apps/pwa && npm run dev` |
| PWA build | `npm run pwa:build` |
| PWA sync libs | `npm run pwa:sync-libs` — run when `libs/types` or `libs/api-client` change |
| Mobile | `npm run mobile` or `cd apps/mobile && npx expo start` |
| Mobile iOS | `npm run mobile:ios` |
| Mobile Android | `npm run mobile:android` |

PWA: após alterar libs, rode `npm run pwa:sync-libs`.

## Conventions

- [.cursor/rules/architecture.md](.cursor/rules/architecture.md) — Monorepo structure and shared libs
- [.cursor/rules/typescript.md](.cursor/rules/typescript.md) — TypeScript (PWA, mobile)
- [.cursor/rules/api-conventions.md](.cursor/rules/api-conventions.md) — API design patterns
