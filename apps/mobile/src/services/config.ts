// Configuração da API
// Para desenvolvimento local, use o IP da máquina ou localhost
// Em produção, altere para a URL do servidor

import { Platform } from "react-native";

// Em iOS Simulator, localhost funciona
// Em Android Emulator, use 10.0.2.2 para acessar localhost da máquina host
// Em dispositivo físico, use o IP da máquina na rede local
const getBaseUrl = () => {
  if (__DEV__) {
    // Desenvolvimento
    if (Platform.OS === "android") {
      // Android Emulator: 10.0.2.2 aponta para localhost da máquina host
      return "http://10.0.2.2:5000";
    }
    // iOS Simulator ou dispositivo: use localhost ou IP da máquina
    return "http://localhost:5000";
  }
  // Produção - altere para sua URL de produção
  return "https://api.fittracker.com";
};

export const API_BASE_URL = getBaseUrl();

// Expo substitui `process.env.EXPO_PUBLIC_*` em tempo de build; o tsconfig do
// projeto não inclui os tipos do Node, então declaramos só o que usamos.
declare const process: { env: Record<string, string | undefined> };

// Login social: apenas client IDs públicos (nenhum segredo no app).
// O client secret do GitHub fica somente no backend.
export const SOCIAL_AUTH = {
  googleWebClientId:
    process.env.EXPO_PUBLIC_GOOGLE_WEB_CLIENT_ID || undefined,
  googleIosClientId:
    process.env.EXPO_PUBLIC_GOOGLE_IOS_CLIENT_ID || undefined,
  githubClientId: process.env.EXPO_PUBLIC_GITHUB_CLIENT_ID || undefined,
};

export const API_ENDPOINTS = {
  // Health
  health: "/health",

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
