// API Configuration
// Platform-agnostic configuration

let baseUrl = "http://localhost:5000";

export const setApiBaseUrl = (url: string) => {
  baseUrl = url;
};

export const getApiBaseUrl = () => baseUrl;

export const API_ENDPOINTS = {
  // Health
  health: "/health",

  // Auth
  register: "/api/auth/register",
  login: "/api/auth/login",
  me: "/api/auth/me",
  refresh: "/api/auth/refresh",
  changePassword: "/api/auth/change-password",

  // Profiles
  profiles: "/api/profiles",
  profile: (id: string) => `/api/profiles/${id}`,

  // Workouts
  workouts: (profileId: string) => `/api/profiles/${profileId}/workouts`,
  workout: (profileId: string, id: string) =>
    `/api/profiles/${profileId}/workouts/${id}`,
  todayWorkouts: (profileId: string) =>
    `/api/profiles/${profileId}/workouts/today`,

  // Bioimpedance
  bioimpedance: (profileId: string) =>
    `/api/profiles/${profileId}/bioimpedance`,
  bioimpedanceItem: (profileId: string, id: string) =>
    `/api/profiles/${profileId}/bioimpedance/${id}`,
  latestBioimpedance: (profileId: string) =>
    `/api/profiles/${profileId}/bioimpedance/latest`,

  // Water intake
  water: (profileId: string) => `/api/profiles/${profileId}/water`,
  waterByDate: (profileId: string, date: string) =>
    `/api/profiles/${profileId}/water?date=${date}`,
  waterEntry: (profileId: string, entryId: string) =>
    `/api/profiles/${profileId}/water/${entryId}`,

  // Completed Workouts
  completedWorkouts: (profileId: string) =>
    `/api/profiles/${profileId}/completed-workouts`,
  completedWorkoutStats: (profileId: string) =>
    `/api/profiles/${profileId}/completed-workouts/stats`,
  completedWorkout: (profileId: string, id: string) =>
    `/api/profiles/${profileId}/completed-workouts/${id}`,

  // AI Planning
  aiPlanningGenerate: "/api/ai/planning/generate",
  aiPlanningPreview: (planId: string) => `/api/ai/planning/preview/${planId}`,
  aiPlanningAccept: "/api/ai/planning/accept",
  aiPlanningStatus: "/api/ai/planning/status",
};
