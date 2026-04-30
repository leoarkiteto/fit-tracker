# FitTracker

Fitness activity hub: workout tracking, bioimpedance, water intake, and AI-powered planning.

## Structure

- **frontend/** — React Native (Expo); see [frontend/AGENTS.md](frontend/AGENTS.md)
- **backend/** — ASP.NET Core API; see [backend/AGENTS.md](backend/AGENTS.md)
- **api-go/** — Go (stdlib) API migration; see [api-go/README.md](api-go/README.md)
- **api-gin/** — Go (Gin + GORM) API migration; see [api-gin/README.md](api-gin/README.md)

## Commands

| Target      | Command |
|-------------|---------|
| Frontend    | `npm run frontend` or `cd frontend && npx expo start` |
| Backend     | `npm run backend` or `cd backend && dotnet run --urls="http://0.0.0.0:5000"` |
| API (Go stdlib) | `cd api-go && go run ./cmd/api/main.go` |
| API (Go Gin)    | `cd api-gin && go run ./cmd/api/main.go` |

## Conventions

- [.cursor/rules/architecture.md](.cursor/rules/architecture.md) — Project structure
- [.cursor/rules/typescript.md](.cursor/rules/typescript.md) — TypeScript (frontend)
- [.cursor/rules/react-native.md](.cursor/rules/react-native.md) — Expo, navigation, API
- [.cursor/rules/csharp.md](.cursor/rules/csharp.md) — C# conventions (backend)
- [.cursor/rules/api-conventions.md](.cursor/rules/api-conventions.md) — API design patterns
