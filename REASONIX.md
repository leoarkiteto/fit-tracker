# FitTracker

Fitness activity hub: workout tracking, bioimpedance, water intake, AI-powered planning via Ollama.

## Stack

- **Frontend** — React Native (Expo SDK 55, RN 0.83), TypeScript strict, React Navigation 7, AsyncStorage [package.json](frontend/package.json)
- **Backend** — Go 1.26.2, stdlib HTTP, SQLite (modernc.org/sqlite), JWT (golang-jwt/jwt/v5), langchaingo (Ollama) [go.mod](backend/go.mod)
- **AI** — All backends query Ollama at `http://localhost:11434` (model varies)

## Layout

| Path | Contents |
|------|----------|
| `frontend/src/components/` | Reusable UI (Button, Card, Input, DaySelector, etc.) |
| `frontend/src/screens/` | Screen components (Home, Workouts, Auth, Bioimpedance, Water) |
| `frontend/src/navigation/` | React Navigation setup (tab, stack, auth) |
| `frontend/src/services/` | API client (`api.ts`) + endpoint config (`config.ts`) |
| `frontend/src/context/` | React context (AppContext, AuthContext) |
| `frontend/src/storage/` | AsyncStorage helpers (profile ID, token) |
| `frontend/src/types/` | Shared TypeScript types |
| `frontend/src/theme/` | Colors, spacing, typography |
| `frontend/src/utils/` | Utilities (`generateId.ts` — temp client-side IDs) |
| `backend/cmd/api/` | Go entry point and route registration |
| `backend/internal/{handlers,models,auth,database,ai}/` | Go stdlib API packages |

## Commands

| Target | Command |
|--------|---------|
| Frontend dev | `npm run frontend` (`cd frontend && expo start`) |
| Frontend iOS | `npm run frontend:ios` (`cd frontend && expo run:ios`) |
| Frontend Android | `npm run frontend:android` (`cd frontend && expo run:android`) |
| Backend dev | `cd backend && go run ./cmd/api` (or `make run`) |
| Backend test | `cd backend && go test ./...` (or `make test`) |
| Backend fmt | `cd backend && go fmt ./...` (or `make fmt`) |
| Backend vet | `cd backend && go vet ./...` |

## Conventions

- **Named exports** — Components, screens, navigators use named exports, re-exported via barrel `index.ts` files [components/index.ts](frontend/src/components/index.ts:1) [screens/index.ts](frontend/src/screens/index.ts:1)
- **Relative imports** — Imports use `../` relative paths (no path aliases) [Button.tsx](frontend/src/components/Button.tsx:8)
- **TypeScript strict** — `"strict": true` in tsconfig [tsconfig.json](frontend/tsconfig.json:7)
- **Go standard layout** — Entry point in `cmd/`, internal packages under `internal/` [main.go](backend/cmd/api/main.go:1)
- **API routing** — Profile-scoped resources at `/api/profiles/{profileId}/...`, auth at `/api/auth/...`, AI at `/api/ai/...` [config.ts](frontend/src/services/config.ts:30)
- **SQLite** — Backend uses SQLite (`fittracker.db`) via modernc.org/sqlite
- **JWT auth** — Backend authenticates via JWT bearer tokens (golang-jwt/jwt/v5)
- **Ollama AI** — AI features use Ollama (default `http://localhost:11434`, model configurable)

## Watch out for

- **Stale root scripts** — `npm run backend` and `npm run backend:build` in [package.json](package.json) reference `dotnet`, but the actual backend is Go. Use `cd backend && go run ./cmd/api` instead.
- **No test infrastructure** — No test files found in `frontend/src/` or `backend/`. No test scripts in `package.json`.
- **No lint/format config** — No ESLint, Prettier, or .editorconfig at project level (Go `make fmt`/`make vet` work via the Makefile).
- **Go 1.26.2** — Very recent; verify build toolchain compatibility if CI or teammates use older Go.
