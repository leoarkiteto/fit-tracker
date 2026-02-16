# API Conventions (ASP.NET Core)

## Vertical Slice Architecture

Each feature is a slice under `Features/<FeatureName>/`:

- `*Endpoints.cs` — Minimal API route handlers
- `*Dtos.cs` — Request/response models
- `*Entry.cs` / entities — Domain models (if needed)
- Services — `I*Service` + implementation

## Patterns

- Use `MapGroup` with `WithTags` for OpenAPI grouping
- Route pattern: `/api/profiles/{profileId:guid}/<resource>`
- JWT Bearer auth on protected endpoints
- DTOs use record types; camelCase JSON
