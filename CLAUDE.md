# FitTracker

Fitness activity hub: workout tracking, bioimpedance, water intake, AI-powered planning via Ollama.

## Stack

- **Monorepo** — pnpm workspace (`pnpm@12.18.0`) + Turborepo; workspace packages are `apps/*` and `packages/*` [pnpm-workspace.yaml](pnpm-workspace.yaml:1) [turbo.json](turbo.json:1)
- **Mobile** (`@fit-tracker/mobile`) — React Native (Expo SDK 57, RN 0.86), TypeScript strict, React Navigation 7, AsyncStorage [package.json](apps/mobile/package.json)
- **API** (`@fit-tracker/api`) — Go 1.26.2, stdlib HTTP, SQLite (modernc.org/sqlite), JWT (golang-jwt/jwt/v5), langchaingo (Ollama) [go.mod](apps/api/go.mod)
- **AI** — All backends query Ollama at `http://localhost:11434` (model varies)

## Layout

| Path | Contents |
|------|----------|
| `apps/mobile/src/components/` | Reusable UI (Button, Card, Input, DaySelector, etc.) |
| `apps/mobile/src/screens/` | Screen components (Home, Workouts, Auth, Bioimpedance, Water) |
| `apps/mobile/src/navigation/` | React Navigation setup (tab, stack, auth) |
| `apps/mobile/src/services/` | API client (`api.ts`) + endpoint config (`config.ts`) |
| `apps/mobile/src/context/` | React context (AppContext, AuthContext) |
| `apps/mobile/src/storage/` | AsyncStorage helpers (profile ID, token) |
| `apps/mobile/src/types/` | Shared TypeScript types |
| `apps/mobile/src/theme/` | Colors, spacing, typography |
| `apps/mobile/src/utils/` | Utilities (`generateId.ts` — temp client-side IDs) |
| `apps/api/cmd/api/` | Composition root: database, token signer, per-slice routers [main.go](apps/api/cmd/api/main.go:1) |
| `apps/api/internal/<slice>/` | Vertical slices: `user`, `profile`, `workout`, `completedworkout`, `bioimpedance`, `water`, `aiplanning` |
| `apps/api/pkg/{apperror,auth,database,httpx}/` | Shared backend infrastructure |
| `packages/` | Shared workspace packages (empty for now) |

## Commands

Run everything from the repository root; reach a sub-project with `pnpm --filter <package>` or `go -C apps/api`, not `cd`.

| Target | Command |
|--------|---------|
| Install dependencies | `pnpm install` (root only; one `pnpm-lock.yaml` at the root) |
| All dev servers | `pnpm dev` (`turbo dev`) |
| Mobile dev | `pnpm --filter @fit-tracker/mobile dev` (`expo start --ios`) |
| Mobile Android | `pnpm --filter @fit-tracker/mobile android` |
| Mobile typecheck | `pnpm --filter @fit-tracker/mobile typecheck` (or `pnpm typecheck` for all packages) |
| API dev | `pnpm --filter @fit-tracker/api dev` (`make run`), or `make -C apps/api run` |
| API test | `go -C apps/api test ./...` (or `make -C apps/api test`) |
| API fmt | `go -C apps/api fmt ./...` (or `make -C apps/api fmt`) |
| API vet | `go -C apps/api vet ./...` (or `make -C apps/api vet`) |

## Conventions

- **Named exports** — Components, screens, navigators use named exports, re-exported via barrel `index.ts` files [components/index.ts](apps/mobile/src/components/index.ts:1) [screens/index.ts](apps/mobile/src/screens/index.ts:1)
- **Relative imports** — Imports use `../` relative paths (no path aliases) [Button.tsx](apps/mobile/src/components/Button.tsx:8)
- **TypeScript strict** — `"strict": true` in tsconfig [tsconfig.json](apps/mobile/tsconfig.json:13)
- **Vertical slices (backend)** — Every backend feature is a package under `apps/api/internal/<slice>/` holding `<slice>.go` (domain types + errors), `port.go` (outbound interfaces), `service/` (use cases, one file per case), `handler/` (`handler.go` with `Router(mux, signer, svc)` plus one file per case), `infrastructure/` (SQLite adapters) and `mock/` (in-memory test doubles) [port.go](apps/api/internal/workout/port.go:1)
- **Services depend on ports only** — A slice's `service` package imports its own domain package, never `net/http` or `database/sql`, so it is unit tested without HTTP or a database [register.go](apps/api/internal/user/service/register.go:1)
- **Errors carry their status** — Services return `*apperror.Error`; handlers render it with `apperror.Write` [apperror.go](apps/api/pkg/apperror/apperror.go:1)
- **Go standard layout** — Entry point in `cmd/`, internal packages under `internal/`, shared infrastructure under `pkg/` [main.go](apps/api/cmd/api/main.go:1)
- **API routing** — Profile-scoped resources at `/api/profiles/{profileId}/...`, auth at `/api/auth/...`, AI at `/api/ai/...` [config.ts](apps/mobile/src/services/config.ts:30)
- **SQLite** — Backend uses SQLite (`fittracker.db`, git-ignored) via modernc.org/sqlite
- **JWT auth** — Backend authenticates via JWT bearer tokens (golang-jwt/jwt/v5)
- **Ollama AI** — AI features use Ollama (default `http://localhost:11434`, model configurable)

## Watch out for

- **Install only from the root with pnpm** — Dependencies are managed by the pnpm workspace; the only lockfile is the root `pnpm-lock.yaml`. Do not run `npm install` or `npx expo install` inside `apps/mobile` (that creates a second lockfile and an npm-style `node_modules`). Add a mobile dependency with `pnpm --filter @fit-tracker/mobile add <pkg>`.
- **JWT_SECRET is required** — The API reads `JWT_SECRET` from the environment and refuses to start without it; `apps/api/.env` is not loaded automatically, so export it (or use a `.env` loader) before starting the API. `turbo dev` passes `JWT_SECRET`, `GOOGLE_CLIENT_IDS`, `GITHUB_CLIENT_ID` and `GITHUB_CLIENT_SECRET` through [turbo.json](turbo.json:6).
- **Mobile has no tests** — `apps/api` ships per-slice service tests (`go -C apps/api test ./...`), but `apps/mobile/src/` has none and `apps/mobile/package.json` has no test script.
- **Mobile `lint` script does not work** — `apps/mobile/package.json` declares `"lint": "eslint ."`, but ESLint is not a dependency and there is no ESLint config, so `pnpm lint` fails for mobile. No Prettier or `.editorconfig` either. The API's `make lint` runs `golangci-lint` only if installed; `go vet` is the reliable check.
- **Go 1.26.2** — Very recent; verify build toolchain compatibility if CI or teammates use older Go.
- **`uuid` override for `xcode` is now inert** — `apps/mobile/package.json` still pins `overrides.xcode.uuid` to `^11.1.1` (added to clear GHSA-w5hq-g745-h8pq while `xcode@3.0.1` depended on vulnerable `uuid`). pnpm ignores `overrides` in a workspace package (it reads only root `pnpm.overrides` / `pnpm-workspace.yaml`), and the current lockfile resolves `xcode@4.0.0`, which no longer depends on `uuid`. If a future Expo upgrade brings back an `xcode` that depends on old `uuid`, move the override to the root.
- **`expo` must stay on the SDK 57 line** — `expo` was once set to `^46.0.21` while every `expo-*` sibling was `~55.0.x`; that mismatch pulled in the whole SDK 46 toolchain and reintroduced 12 high / 2 critical advisories. Keep `expo` on the SDK 57 line (currently `^57.0.24`) and run `pnpm --filter @fit-tracker/mobile exec expo install --check` after touching Expo deps.
- **Xcode 27 replaced `Simulator.app` with `DeviceHub.app`** — Xcode 27 dropped `Xcode.app/Contents/Developer/Applications/Simulator.app` in favour of `Contents/Applications/DeviceHub.app` (bundle ID `com.apple.dt.Devices`), so `open -a Simulator` fails with "Unable to find application named 'Simulator'" (LaunchServices also keeps a stale record for the deleted path). Use `open -a DeviceHub`. This requires `@expo/cli >= 57.0.26`, which `expo@^57.0.24` pulls in; earlier CLIs hardcode `open -a Simulator` / `id of app "Simulator"` and fail with "Can't determine id of Simulator app" or `SIMULATOR_TIMEOUT`.
- **Top-level `splash` is invalid in SDK 57** — The SDK 57 config schema allows `splash` only under `web` (PWA); the native splash is configured through the `expo-splash-screen` config plugin. Leaving a top-level `splash` key makes `expo-doctor` fail with "should NOT have additional property 'splash'".
- **Native iOS builds need working `sandbox-exec`** — Xcode's SwiftPM dependency resolution runs under `sandbox-exec`. Sandboxed shells that deny `sandbox_apply` fail the build with `xcodebuild: error: Could not resolve package dependencies: sandbox-exec: sandbox_apply: Operation not permitted`. Run `expo run:ios` from a normal terminal; this is an environment limit, not a project problem.
