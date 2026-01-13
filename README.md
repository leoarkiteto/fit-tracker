# 💪 FitTracker

A workout management application built as a **Monorepo** with **React Native** mobile app, **Next.js PWA**, and **ASP.NET Core Web API** backend.

## 📱 Overview

FitTracker allows users to:
- Create workouts for each day of the week
- Define exercises, sets, reps, weight, and rest intervals
- Track actual workout executions
- Record bioimpedance assessment results

The product differentiator is the **use of Artificial Intelligence through AI Agents**, offered as premium features via subscription.

## 🏗️ Architecture - Monorepo with Vertical Slice

The project uses a **Monorepo structure** with shared libraries and **Vertical Slice Architecture** for the backend.

### Project Structure

```
fittracker/
├── apps/
│   ├── mobile/                 # React Native (Expo)
│   │   ├── src/
│   │   │   ├── components/
│   │   │   ├── context/
│   │   │   ├── navigation/
│   │   │   ├── screens/
│   │   │   └── ...
│   │   ├── App.tsx
│   │   └── package.json
│   │
│   ├── pwa/                    # Next.js PWA
│   │   ├── src/
│   │   │   ├── app/           # App Router
│   │   │   ├── components/
│   │   │   └── lib/
│   │   ├── public/
│   │   └── package.json
│   │
│   └── api/                    # ASP.NET Core
│       └── FitTracker.Api/
│           ├── Features/
│           │   ├── Auth/
│           │   ├── Profiles/
│           │   ├── Workouts/
│           │   ├── Bioimpedance/
│           │   └── AI/        # Future (Premium)
│           ├── Shared/
│           └── Program.cs
│
├── libs/
│   ├── types/                  # Shared TypeScript types
│   │   └── src/index.ts
│   │
│   └── api-client/             # Shared API service layer
│       └── src/
│           ├── index.ts
│           └── config.ts
│
├── nx.json                     # Nx configuration
├── package.json                # Workspaces config
├── tsconfig.base.json          # Shared TypeScript config
├── PRODUCT_VISION.md
└── README.md
```

### Shared Libraries

| Library | Purpose |
|---------|---------|
| `@fittracker/types` | Shared TypeScript interfaces (Workout, Exercise, UserProfile, etc.) |
| `@fittracker/api-client` | Platform-agnostic API service layer (works in React Native and Next.js) |

### Backend (API) - Vertical Slice

```
apps/api/FitTracker.Api/
├── Features/                    # 🎯 Each feature is a "vertical slice"
│   ├── Auth/
│   ├── Profiles/
│   ├── Workouts/
│   ├── Bioimpedance/
│   └── AI/                      # 🧠 Future (Premium)
├── Shared/
└── Program.cs
```

### Benefits of This Architecture

| Aspect | Benefit |
|--------|---------|
| **Code Sharing** | Types and API client shared between mobile and PWA |
| **Cohesion** | Features organized as vertical slices |
| **Maintenance** | Easy to find and modify code |
| **Platform Parity** | Same features on mobile and web |
| **Scalability** | Add new apps/libs without affecting existing ones |
| **AI Agents** | Each agent will be a separate feature |

## 🚀 How to Run

### Install Dependencies

```bash
# From the monorepo root
npm install
```

### Backend (API)

```bash
npm run api
# Or directly:
cd apps/api/FitTracker.Api && dotnet run --urls="http://0.0.0.0:5000"
```

The API will be available at `http://localhost:5000` with Swagger UI.

### Mobile App (React Native)

```bash
npm run mobile
# Or directly:
cd apps/mobile && npx expo start
```

### PWA (Next.js)

```bash
# First time setup
npm run pwa:install

# Run development server
npm run pwa
# Or directly:
cd apps/pwa && npm run dev
```

The PWA will be available at `http://localhost:3000`.

> **Note**: The PWA is installed independently from the main monorepo to avoid React version conflicts between React Native and Next.js.

## 📋 TODO List - Project Progress

### ✅ Implemented Features

#### 🔐 Authentication
- [x] Login Screen (SignIn)
- [x] Registration Screen (SignUp)
- [x] JWT Authentication
- [x] Password change
- [x] Logout
- [x] Session persistence

#### 👤 User Profile
- [x] Profile creation
- [x] Profile editing (name, age, height, current weight, goal weight)
- [x] Basic statistics display

#### 🏋️ Workouts (CRUD)
- [x] Create workouts
- [x] Edit workouts
- [x] Delete workouts
- [x] Set days of the week
- [x] Set goal (hypertrophy, strength, endurance, weight loss, maintenance)
- [x] List today's workouts

#### 💪 Exercises
- [x] Add exercises to workout
- [x] Set sets, reps, weight
- [x] Set rest interval
- [x] Set muscle group
- [x] Exercise notes

#### ⏱️ Workout Execution
- [x] Workout stopwatch
- [x] Rest timer between sets
- [x] Vibration when rest ends (mobile)
- [x] Mark exercise as completed
- [x] Record workout duration
- [x] Lock workout after daily completion

#### 📊 Bioimpedance
- [x] Register bioimpedance assessment
- [x] Edit assessment
- [x] Assessment history
- [x] Fields: weight, body fat %, muscle mass, bone mass, water %, visceral fat, BMR, metabolic age

#### 📅 Calendar
- [x] View days with completed workouts
- [x] Visual marking on calendar

#### 📈 Basic Statistics
- [x] Total completed workouts
- [x] Workouts this week
- [x] Total minutes spent

#### 🔧 Infrastructure
- [x] REST API with Minimal API
- [x] Entity Framework Core + SQLite
- [x] CORS configured
- [x] Swagger/OpenAPI with JWT support
- [x] Refresh Token
- [x] **Vertical Slice Architecture** (API restructured)
- [x] **Monorepo structure** (Nx-style workspaces)
- [x] **Shared libraries** (types, api-client)
- [x] **Next.js PWA** (Progressive Web App)

---

### ❌ Pending Features

#### 📝 Workout Improvements (Free)
- [ ] Subjective feedback after workout (light / ok / heavy)
- [ ] RPE (Rate of Perceived Exertion) per set
- [ ] Weight history per exercise
- [ ] Weight progression chart
- [ ] Exercise reordering (drag & drop)
- [ ] Duplicate existing workout
- [ ] Pre-defined workout templates

#### 📊 Basic Analytics (Free)
- [ ] Bioimpedance comparison between dates
- [ ] Weight progression chart
- [ ] Body composition progression chart
- [ ] Weekly volume by muscle group

#### 💳 Subscription System
- [ ] Payment system integration
- [ ] Plans screen (Free vs Premium)
- [ ] Subscription management
- [ ] Trial period

---

### 🧠 Premium Features (AI) - Not Implemented

#### 🤖 Virtual Coach Agent
- [ ] AI Chat
- [ ] Short-term memory (conversation context)
- [ ] Long-term memory (user history)
- [ ] Answer workout questions
- [ ] Suggest adjustments based on data
- [ ] Explain AI decisions

#### 📊 Progress Analyst Agent
- [ ] Weight progression analysis
- [ ] Weekly volume analysis
- [ ] Training frequency analysis
- [ ] Weekly reports in natural language
- [ ] Monthly reports
- [ ] Workout x bioimpedance correlation

#### 🧪 Bioimpedance Agent
- [ ] Automatic data interpretation
- [ ] Body changes explanation
- [ ] Relate results to training
- [ ] Data-based adjustment suggestions

#### ⚠️ Injury Prevention Agent
- [ ] Abrupt weight increase detection
- [ ] Excessive volume alert
- [ ] Insufficient rest alert
- [ ] Performance drop identification
- [ ] Preventive alerts with disclaimers

#### 🗓️ Workout Planning Agent
- [ ] Automatic weekly workout creation
- [ ] Mesocycle planning (4-8 weeks)
- [ ] Deload programming
- [ ] Adaptation by goal/availability/level

#### 🔥 Automatic Adaptive Workouts
- [ ] Automatic weight adjustment
- [ ] Rep adjustment
- [ ] Volume adjustment
- [ ] Based on actual performance and bioimpedance

#### 🔄 Smart Exercise Substitution
- [ ] Equivalent exercise suggestions
- [ ] Maintain muscle stimulus
- [ ] Consider available equipment

#### 🎯 Results Prediction (Advanced)
- [ ] Estimates based on historical data
- [ ] Evolution projection
- [ ] Present as approximations (with disclaimers)

---

### 🔜 Suggested Next Steps

1. **Phase 1 - Free Improvements**
   - Implement subjective feedback (light/ok/heavy)
   - Add progression charts
   - Weight history per exercise

2. **Phase 2 - Subscription System**
   - Integrate payment system
   - Create plans screen
   - Implement premium access control

3. **Phase 3 - Basic AI**
   - Implement Virtual Coach Agent (chat)
   - LLM integration (OpenAI/Claude)
   - Context memory

4. **Phase 4 - Advanced AI**
   - Implement remaining agents
   - Adaptive workouts
   - Predictive analytics

---

## 📄 License

This project is private and for exclusive use.

---

## 👥 Contributors

- Leo França - Lead Developer
