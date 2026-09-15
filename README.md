# FitTracker

A workout management application: a **React Native** mobile app and a **Go** backend (standard library HTTP + SQLite), with AI workout planning through a local **Ollama** model.

## Overview

FitTracker allows users to:
- Create workouts for each day of the week
- Define exercises, sets, reps, weight, and rest intervals
- Track actual workout executions and see weekly stats
- Record bioimpedance assessments and daily water intake

The product differentiator is the **use of Artificial Intelligence through AI Agents**, offered as premium features via subscription.

Development follows the project constitution (Clean Code, design patterns when necessary, simplicity/accessibility/speed). See `.specify/memory/constitution.md`.

## Stack

| Layer | Technology |
|-------|------------|
| Mobile | React Native 0.83 on Expo SDK 55, TypeScript (strict), React Navigation 7, AsyncStorage |
| API | Go 1.26.2, stdlib `net/http` (`ServeMux` patterns), SQLite via `modernc.org/sqlite` |
| Auth | JWT bearer tokens (`golang-jwt/jwt/v5`) + bcrypt password hashing |
| AI | `langchaingo` talking to Ollama at `http://localhost:11434` (default model `llama3.2`) |

The backend follows **vertical slice architecture**: every feature is a self-contained package under `backend/internal/<slice>/` that owns its domain, ports, use cases, HTTP handlers, SQLite adapter and test doubles — so a feature can be changed in one place.

## Project Structure

```
fittracker/
├── frontend/                      # React Native (Expo)
│   ├── src/
│   │   ├── components/            # Reusable UI
│   │   ├── context/               # AppContext, AuthContext
│   │   ├── navigation/            # tab / stack / auth navigators
│   │   ├── screens/               # Auth, Home, Workouts, Bioimpedance, Water, Planning
│   │   ├── services/              # api.ts + config.ts (endpoints)
│   │   ├── storage/               # AsyncStorage helpers (token, profile id)
│   │   ├── theme/                 # colors, spacing, typography
│   │   ├── types/                 # shared TypeScript types
│   │   └── utils/                 # helpers (temp client-side ids)
│   ├── App.tsx
│   └── package.json
│
├── backend/                       # Go
│   ├── cmd/api/main.go            # composition root: database, signer, per-slice routers
│   ├── pkg/                       # shared infrastructure
│   │   ├── apperror/              # domain error + its HTTP status
│   │   ├── auth/                  # bcrypt + JWT signing/validation
│   │   ├── database/              # SQLite connection + schema
│   │   └── httpx/                 # auth guard, body limit, JSON responses, CORS
│   ├── internal/                  # one package per feature slice
│   │   ├── user/                  # register, login, me, change password
│   │   ├── profile/               # profile read/update/delete (cascading)
│   │   ├── workout/               # workouts + exercises, today's list
│   │   ├── completedworkout/      # execution history + stats
│   │   ├── bioimpedance/          # measurements, latest reading
│   │   ├── water/                 # daily intake + goal
│   │   └── aiplanning/            # Ollama plan generation and acceptance
│   ├── Makefile
│   └── README.md
│
├── package.json                   # Convenience scripts only
└── README.md
```

Every slice repeats the same shape: `<slice>.go` (domain), `port.go` (outbound interfaces), `service/` (use cases, one file each), `handler/` (`Router()` plus one file per use case), `infrastructure/` (SQLite), and `mock/` (in-memory test doubles).

## How to Run

### Backend (API)

```bash
cd backend
export JWT_SECRET=my-secret          # required: the server refuses to start without it
go run ./cmd/api                     # or: make run
```

The API listens on `http://localhost:5000`; `GET /health` returns `200 OK`.

`backend/.env` holds `JWT_SECRET` for convenience but **nothing loads it automatically**, so export the variable (or `set -a && . ./.env && set +a`) before starting the server.

AI features are optional: without a running Ollama (`ollama serve`, model pulled with `ollama pull llama3.2`) the rest of the API works and `GET /api/ai/planning/status` reports the backend as unreachable.

### Frontend (React Native)

Install dependencies once:

```bash
cd frontend && npm install
```

Then start:

```bash
cd frontend && npx expo start
```

Or from root: `npm run frontend`, `npm run frontend:ios`, `npm run frontend:android`.

The app targets `http://localhost:5000` (iOS simulator) and `http://10.0.2.2:5000` (Android emulator) in development — see `frontend/src/services/config.ts`.

## Tests

```bash
cd backend && go test ./...          # or: make test
```

Each slice ships service unit tests that run without HTTP and without a database, wiring in the in-memory fakes from `internal/<slice>/mock`.

## API

| Area | Routes |
|------|--------|
| Health | `GET /health` |
| Auth | `POST /api/auth/register`, `POST /api/auth/login`, `GET /api/auth/me`, `PATCH /api/auth/change-password` |
| Profiles | `GET /api/profiles`, `GET /api/profiles/{id}`, `PUT /api/profiles/{id}`, `DELETE /api/profiles/{id}` |
| Workouts | `GET/POST /api/profiles/{profileId}/workouts`, `GET .../workouts/today`, `GET/PUT/DELETE .../workouts/{id}` |
| Executions | `GET/POST /api/profiles/{profileId}/completed-workouts`, `GET .../completed-workouts/stats` |
| Bioimpedance | `GET/POST /api/profiles/{profileId}/bioimpedance`, `GET .../bioimpedance/latest`, `DELETE .../bioimpedance/{id}` |
| Water | `GET/POST /api/profiles/{profileId}/water` (optional `?date=YYYY-MM-DD`), `DELETE .../water/{id}` |
| AI planning | `POST /api/ai/planning/generate`, `POST /api/ai/planning/accept`, `GET /api/ai/planning/status` |

All routes except `/health` and `GET /api/ai/planning/status` require a `Authorization: Bearer <token>` header.

## Implemented Features

- **Auth**: Sign up, login, JWT (24 h), current user, password change — no refresh token or logout endpoint yet
- **Profiles**: Read, update and delete (deleting a profile cascades to its workouts, executions, bioimpedance and water entries)
- **Workouts**: CRUD, days of the week, goals, today's list
- **Exercises**: Sets, reps, weight, rest, muscle group, notes
- **Execution**: Record a completed workout and duration; stats for total, this week and total minutes
- **Bioimpedance**: Create a measurement, history, latest reading, delete
- **Water intake**: Daily summary with a goal derived from body weight (35 ml/kg, 70 kg fallback), create and delete entries
- **AI planning**: Generate a weekly plan with Ollama, accept it (persists the workouts), and check backend status
- **Stats**: Completed workouts, this week, total minutes

## Known Gaps

- **Stale root scripts** — `npm run backend` and `npm run backend:build` in `package.json` still invoke the retired `dotnet` backend and fail; use the Go commands above.
- **No refresh token / logout** — tokens simply expire after 24 hours.
- **`/api/ai/planning/preview/{planId}` is not implemented** — `frontend/src/services/config.ts` declares it, but the backend only offers generate/accept/status.
- **No frontend tests** — only the backend slices have automated tests.
- **No CI** — nothing runs `go test`, `go vet` or `gofmt` automatically.

## License

This project is private and for exclusive use.
