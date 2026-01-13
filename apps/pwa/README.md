# FitTracker PWA

A Progressive Web Application for the FitTracker workout management system.

## Features

- **Dashboard**: View workout stats and today's scheduled workouts
- **Workouts**: Create, view, edit, and delete workouts with exercises
- **Workout Session**: Start workout timer, track exercises with rest periods
- **Body Metrics**: Track bioimpedance data (weight, body fat, muscle mass, etc.)
- **Profile**: Manage personal information and account settings
- **Authentication**: Sign in, sign up, and change password
- **PWA**: Installable, offline-capable progressive web app

## Getting Started

### Prerequisites

- Node.js 18+
- npm or yarn
- Running API server (apps/api)

### Installation

```bash
# From the monorepo root
npm install

# Or from the PWA directory
cd apps/pwa
npm install
```

### Development

```bash
# From the PWA directory
npm run dev

# Or from the monorepo root
npm run pwa
```

### Production Build

```bash
npm run build
npm run start
```

### Environment Variables

Create a `.env.local` file:

```env
NEXT_PUBLIC_API_URL=http://localhost:5000
```

## Project Structure

```
apps/pwa/
├── public/
│   ├── icons/          # PWA icons
│   ├── manifest.json   # PWA manifest
│   └── sw.js          # Service worker
├── src/
│   ├── app/           # Next.js App Router pages
│   │   ├── (app)/     # Authenticated routes
│   │   │   ├── dashboard/
│   │   │   ├── workouts/
│   │   │   ├── bioimpedance/
│   │   │   └── profile/
│   │   └── auth/      # Authentication routes
│   ├── components/    # React components
│   │   ├── layout/    # Navigation components
│   │   └── ui/        # Reusable UI components
│   └── lib/           # Utilities and contexts
├── tailwind.config.js
├── next.config.js
└── tsconfig.json
```

## Technology Stack

- **Next.js 14**: React framework with App Router
- **TypeScript**: Type safety
- **Tailwind CSS**: Styling
- **Lucide React**: Icons
- **Shared libs**: Types and API client are copied into `src/lib/` for isolated builds

## PWA Features

- **Installable**: Add to home screen on mobile and desktop
- **Offline Support**: Service worker caches pages for offline access
- **Responsive**: Works on mobile, tablet, and desktop
- **Native Feel**: Standalone display mode, custom theme color
