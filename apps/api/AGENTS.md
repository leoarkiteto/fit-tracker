# FitTracker API

ASP.NET Core backend for the fitness hub: auth, profiles, workouts, bioimpedance, water intake, AI planning.

## Stack

- .NET 10, Minimal APIs, EF Core + SQLite
- JWT auth, Semantic Kernel + Ollama for AI

## Build / Run

```bash
dotnet run --urls="http://0.0.0.0:5000"
# or from root: npm run api
```

## Conventions

- [.cursor/rules/api-conventions.md](../../.cursor/rules/api-conventions.md) — Vertical Slice, endpoints, DTOs
- [.cursor/rules/csharp.md](../../.cursor/rules/csharp.md) — C# style
