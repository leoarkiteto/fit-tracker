# FitTracker Mobile

React Native (Expo) app for the fitness hub: workouts, bioimpedance, water, AI planning. Workspace package `@fit-tracker/mobile` in the pnpm + Turborepo monorepo.

## Stack

- Expo SDK 57, React Native 0.86, React 19, React Navigation 7, TypeScript strict, AsyncStorage [package.json](package.json)

## Build / Run

Run from the repository root (no `cd`):

```bash
pnpm --filter @fit-tracker/mobile dev        # expo start --ios
pnpm --filter @fit-tracker/mobile android    # expo start --android
pnpm --filter @fit-tracker/mobile typecheck  # tsc --noEmit
```

Install dependencies only from the root with `pnpm install`; add one with `pnpm --filter @fit-tracker/mobile add <pkg>`. Never run `npm install` or `npx expo install` in this directory.

## Conventions

See the root [CLAUDE.md](../../CLAUDE.md): layout of `src/`, named exports via barrel `index.ts`, relative `../` imports, API endpoints in [config.ts](src/services/config.ts), and the Expo/Xcode caveats.

- No tests and no working lint: the `lint` script (`eslint .`) fails because ESLint is not installed. `typecheck` is the only automated check.
