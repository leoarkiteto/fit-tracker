# 💪 FitTracker

A mobile workout management application built with **React Native** and **ASP.NET Core Web API** backend.

## 📱 Overview

FitTracker allows users to:
- Create workouts for each day of the week
- Define exercises, sets, reps, weight, and rest intervals
- Track actual workout executions
- Record bioimpedance assessment results

The product differentiator is the **use of Artificial Intelligence through AI Agents**, offered as premium features via subscription.

## 🏗️ Architecture - Vertical Slice

The project uses **Vertical Slice Architecture**, where each feature is organized in an isolated and self-contained manner. This facilitates maintenance, testing, and adding new features (especially AI Agents).

### Backend (API) - Vertical Slice

```
api/FitTracker.Api/
├── Features/                         # 🎯 Each feature is a "vertical slice"
│   ├── Auth/                         # Authentication
│   │   ├── User.cs                   # Model
│   │   ├── AuthDtos.cs               # DTOs
│   │   ├── AuthEndpoints.cs          # Endpoints
│   │   ├── JwtService.cs             # JWT Service
│   │   └── PasswordService.cs        # Password Service
│   │
│   ├── Profiles/                     # User profiles
│   │   ├── UserProfile.cs
│   │   ├── ProfileDtos.cs
│   │   └── ProfileEndpoints.cs
│   │
│   ├── Workouts/                     # Workouts and exercises
│   │   ├── Workout.cs
│   │   ├── Exercise.cs
│   │   ├── CompletedWorkout.cs
│   │   ├── WorkoutDtos.cs
│   │   ├── WorkoutEndpoints.cs
│   │   └── CompletedWorkoutEndpoints.cs
│   │
│   ├── Bioimpedance/                 # Bioimpedance
│   │   ├── BioimpedanceData.cs
│   │   ├── BioimpedanceDtos.cs
│   │   └── BioimpedanceEndpoints.cs
│   │
│   └── AI/                           # 🧠 Future (Premium)
│       ├── Coach/                    # Virtual Coach Agent
│       ├── Progress/                 # Analyst Agent
│       ├── Bioimpedance/             # Bioimpedance Agent
│       ├── Prevention/               # Prevention Agent
│       └── Planning/                 # Planning Agent
│
├── Shared/                           # Shared code
│   └── Data/
│       └── FitTrackerDbContext.cs
│
└── Program.cs                        # Configuration and mapping
```

### Frontend (React Native) - Current

```
src/
├── components/               # Reusable components
├── context/                  # Contexts (Auth, App)
├── navigation/               # Navigation
├── screens/                  # App screens
├── services/                 # API client
├── theme/                    # Global styles
├── types/                    # TypeScript types
└── utils/                    # Utilities
```

### Benefits of Vertical Slice Architecture

| Aspect | Benefit |
|--------|---------|
| **Cohesion** | Everything related to a feature stays together |
| **Maintenance** | Easy to find and modify code |
| **Testing** | Each feature can be tested in isolation |
| **Scalability** | New features don't affect existing ones |
| **AI Agents** | Each agent will be a separate feature |
| **Team** | Developers can work in parallel |

## 🚀 How to Run

### Backend (API)

```bash
cd api/FitTracker.Api
dotnet run --urls="http://0.0.0.0:5000"
```

The API will be available at `http://localhost:5000` with Swagger UI at the root.

### Frontend (React Native)

```bash
npx expo start
```

## 📋 TODO List - Project Progress

### ✅ Implemented Features

#### 🔐 Authentication
- [x] Login Screen (SignIn)
- [x] Registration Screen (SignUp)
- [x] JWT Authentication
- [x] Password change
- [x] Logout
- [x] Session persistence (AsyncStorage)

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
- [x] Vibration when rest ends
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
