# FitTracker PWA

Next.js PWA for the fitness hub: workouts, bioimpedance, water, AI planning, dashboard.

## Stack

- Next.js 14 (App Router), Tailwind CSS, React 18

## Build / Run

```bash
npm run dev
# or from root: npm run pwa
```

## Setup

When `libs/types` or `libs/api-client` change, run `npm run pwa:sync-libs` from root. See [architecture.md](../../.cursor/rules/architecture.md).

## Conventions

- [.cursor/rules/nextjs.md](../../.cursor/rules/nextjs.md) — App Router, auth, API client
- [.cursor/rules/typescript.md](../../.cursor/rules/typescript.md) — TypeScript
