# FitTracker

Fitness activity hub: workout tracking, bioimpedance, water intake, AI-powered planning via Ollama.

## Stack

- **Frontend** — React Native (Expo SDK 57, RN 0.86), TypeScript strict, React Navigation 7, AsyncStorage [package.json](frontend/package.json)
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
| `backend/cmd/api/` | Composition root: database, token signer, per-slice routers [main.go](backend/cmd/api/main.go:1) |
| `backend/internal/<slice>/` | Vertical slices: `user`, `profile`, `workout`, `completedworkout`, `bioimpedance`, `water`, `aiplanning` |
| `backend/pkg/{apperror,auth,database,httpx}/` | Shared backend infrastructure |

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
- **TypeScript strict** — `"strict": true` in tsconfig [tsconfig.json](frontend/tsconfig.json:13)
- **Vertical slices (backend)** — Every backend feature is a package under `backend/internal/<slice>/` holding `<slice>.go` (domain types + errors), `port.go` (outbound interfaces), `service/` (use cases, one file per case), `handler/` (`handler.go` with `Router(mux, signer, svc)` plus one file per case), `infrastructure/` (SQLite adapters) and `mock/` (in-memory test doubles) [port.go](backend/internal/workout/port.go:1)
- **Services depend on ports only** — A slice's `service` package imports its own domain package, never `net/http` or `database/sql`, so it is unit tested without HTTP or a database [register.go](backend/internal/user/service/register.go:1)
- **Errors carry their status** — Services return `*apperror.Error`; handlers render it with `apperror.Write` [apperror.go](backend/pkg/apperror/apperror.go:1)
- **Go standard layout** — Entry point in `cmd/`, internal packages under `internal/`, shared infrastructure under `pkg/` [main.go](backend/cmd/api/main.go:1)
- **API routing** — Profile-scoped resources at `/api/profiles/{profileId}/...`, auth at `/api/auth/...`, AI at `/api/ai/...` [config.ts](frontend/src/services/config.ts:30)
- **SQLite** — Backend uses SQLite (`fittracker.db`) via modernc.org/sqlite
- **JWT auth** — Backend authenticates via JWT bearer tokens (golang-jwt/jwt/v5)
- **Ollama AI** — AI features use Ollama (default `http://localhost:11434`, model configurable)

## Watch out for

- **Stale root scripts** — `npm run backend` and `npm run backend:build` in [package.json](package.json) reference `dotnet`, but the actual backend is Go. Use `cd backend && go run ./cmd/api` instead.
- **JWT_SECRET is required** — The backend reads `JWT_SECRET` from the environment and refuses to start without it; `backend/.env` is not loaded automatically, so export it (or use a `.env` loader) before `go run ./cmd/api`.
- **Frontend has no tests** — `backend/` ships per-slice service tests (`cd backend && go test ./...`), but `frontend/src/` still has none and `package.json` has no test script.
- **No lint/format config** — No ESLint, Prettier, or .editorconfig at project level (Go `make fmt`/`make vet` work via the Makefile).
- **Go 1.26.2** — Very recent; verify build toolchain compatibility if CI or teammates use older Go.
- **`uuid` override for `xcode`** — `frontend/package.json` pins `overrides.xcode.uuid` to `^11.1.1`, outside the `^7.0.3` that `xcode@3.0.1` declares. This clears GHSA-w5hq-g745-h8pq (uuid v3/v5/v6 buffer bounds check), which no Expo SDK fixes upstream — `@expo/config-plugins` still depends on `xcode@^3.0.1` as of SDK 57 (re-verified on the SDK 57 upgrade). `xcode` only calls `uuid.v4()` and uuid 11 ships a compatible CJS build, so it is safe. Re-check on every Expo upgrade and drop the override once `xcode` bumps `uuid` itself.
- **`expo` must stay on the SDK 57 line** — `expo` was once set to `^46.0.21` while every `expo-*` sibling was `~55.0.x`; that mismatch pulled in the whole SDK 46 toolchain and reintroduced 12 high / 2 critical advisories. Keep `expo` on the SDK 57 line (currently `^57.0.24`) and run `npx expo install --check` after touching Expo deps.
- **Xcode 27 replaced `Simulator.app` with `DeviceHub.app`** — Xcode 27 dropped `Xcode.app/Contents/Developer/Applications/Simulator.app` in favour of `Contents/Applications/DeviceHub.app` (bundle ID `com.apple.dt.Devices`), so `open -a Simulator` fails with "Unable to find application named 'Simulator'" (LaunchServices also keeps a stale record for the deleted path). Use `open -a DeviceHub`. This requires `@expo/cli >= 57.0.26`, which `expo@^57.0.24` pulls in; earlier CLIs hardcode `open -a Simulator` / `id of app "Simulator"` and fail with "Can't determine id of Simulator app" or `SIMULATOR_TIMEOUT`.
- **Top-level `splash` is invalid in SDK 57** — The SDK 57 config schema allows `splash` only under `web` (PWA); the native splash is configured through the `expo-splash-screen` config plugin. Leaving a top-level `splash` key makes `npx expo-doctor` fail with "should NOT have additional property 'splash'".
- **Native iOS builds need working `sandbox-exec`** — Xcode's SwiftPM dependency resolution runs under `sandbox-exec`. Sandboxed shells that deny `sandbox_apply` fail the build with `xcodebuild: error: Could not resolve package dependencies: sandbox-exec: sandbox_apply: Operation not permitted`. Run `npx expo run:ios` from a normal terminal; this is an environment limit, not a project problem.
