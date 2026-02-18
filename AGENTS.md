# FitTracker

Fitness activity hub: workout tracking, bioimpedance, water intake, and AI-powered planning.

## Structure

- **frontend/** — React Native (Expo); see [frontend/AGENTS.md](frontend/AGENTS.md)
- **backend/** — ASP.NET Core API; see [backend/AGENTS.md](backend/AGENTS.md)

## Commands

| Target   | Command |
|----------|---------|
| Backend  | `npm run backend` or `cd backend && dotnet run --urls="http://0.0.0.0:5000"` |
| Frontend | `npm run frontend` or `cd frontend && npx expo start` |

## Conventions

- [.cursor/rules/architecture.md](.cursor/rules/architecture.md) — Project structure
- [.cursor/rules/typescript.md](.cursor/rules/typescript.md) — TypeScript (frontend)
- [.cursor/rules/react-native.md](.cursor/rules/react-native.md) — Expo, navigation, API
- [.cursor/rules/api-conventions.md](.cursor/rules/api-conventions.md) — API design patterns
