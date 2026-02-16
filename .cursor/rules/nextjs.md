# Next.js (PWA) Conventions

- **App Router**: `src/app/` — routes, layouts, pages
- **Auth**: `@/lib/auth-context` for user/session
- **API**: `@/lib/api-client` (synced from `libs/api-client`)
- **UI**: Tailwind CSS (Mobile-first approach); shared components in `src/components/`
- **Paths**: `(app)/` group for authenticated routes; `auth/` for signin/signup
