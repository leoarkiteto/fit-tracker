# FitTracker API

ASP.NET Core backend for the FitTracker fitness hub.

## Requirements

- .NET 10 SDK

## Run

```bash
dotnet run --urls="http://0.0.0.0:5000"
```

API and Swagger UI: `http://localhost:5000`

## Build

```bash
dotnet build
```

## Stack

- Minimal APIs, EF Core + SQLite
- JWT auth, Semantic Kernel + Ollama for AI planning
