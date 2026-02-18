# FitTracker

A workout management application: **React Native** mobile app and **ASP.NET Core Web API** backend.

## Overview

FitTracker allows users to:
- Create workouts for each day of the week
- Define exercises, sets, reps, weight, and rest intervals
- Track actual workout executions
- Record bioimpedance assessment results

The product differentiator is the **use of Artificial Intelligence through AI Agents**, offered as premium features via subscription.

Development follows the project constitution (Clean Code, design patterns when necessary, simplicity/accessibility/speed). See `.specify/memory/constitution.md`.

## Project Structure

```
fittracker/
├── frontend/              # React Native (Expo)
│   ├── src/
│   │   ├── components/
│   │   ├── context/
│   │   ├── navigation/
│   │   ├── screens/
│   │   ├── services/
│   │   └── types/
│   ├── App.tsx
│   └── package.json
│
├── backend/               # ASP.NET Core
│   ├── Features/
│   │   ├── Auth/
│   │   ├── Profiles/
│   │   ├── Workouts/
│   │   ├── Bioimpedance/
│   │   ├── WaterIntake/
│   │   └── AI/
│   ├── Shared/
│   └── Program.cs
│
├── package.json           # Convenience scripts only
└── README.md
```

## How to Run

### Backend (API)

```bash
cd backend && dotnet run --urls="http://0.0.0.0:5000"
```

Or from root: `npm run backend`

The API and Swagger UI are at `http://localhost:5000`.

### Frontend (React Native)

Install dependencies once:

```bash
cd frontend && npm install
```

Then start:

```bash
cd frontend && npx expo start
```

Or from root: `npm run frontend`

## Implemented Features

- **Auth**: Login, signup, JWT, refresh token, password change
- **Profiles**: Create and edit (name, age, height, weight, goal)
- **Workouts**: CRUD, days of week, goals, today’s list
- **Exercises**: Sets, reps, weight, rest, muscle group
- **Execution**: Stopwatch, rest timer, mark complete, duration
- **Bioimpedance**: Register, edit, history
- **Water intake**: Track daily intake
- **AI planning**: Semantic Kernel + Ollama for workout generation
- **Stats**: Completed workouts, this week, total minutes

## License

This project is private and for exclusive use.
