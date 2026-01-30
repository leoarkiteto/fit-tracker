import AsyncStorage from "@react-native-async-storage/async-storage";
import { API_BASE_URL, API_ENDPOINTS } from "./config";

const TOKEN_KEY = "@fittracker:token";
import {
  Workout,
  BioimpedanceData,
  UserProfile,
  Exercise,
  WorkoutGoal,
  DayOfWeek,
  MuscleGroup,
  WaterIntakeEntry,
  DailyWaterSummary,
} from "../types";

// Helper para verificar se um ID é um GUID válido
const isValidGuid = (id: string): boolean => {
  const guidRegex =
    /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i;
  return guidRegex.test(id);
};

// Tipos para as respostas da API
interface ApiUserProfile {
  id: string;
  name: string;
  age: number | null;
  height: number | null;
  currentWeight: number | null;
  goalWeight: number | null;
  avatarUrl: string | null;
}

interface ApiExercise {
  id: string;
  name: string;
  muscleGroup: string;
  sets: number;
  reps: number;
  weight: number | null;
  restSeconds: number;
  notes: string | null;
}

interface ApiWorkout {
  id: string;
  name: string;
  description: string | null;
  goal: string;
  days: string[];
  exercises: ApiExercise[];
  createdAt: string;
  updatedAt: string;
  isCompleted: boolean;
  completedAt: string | null;
}

interface ApiBioimpedance {
  id: string;
  date: string;
  weight: number;
  bodyFatPercentage: number;
  muscleMass: number;
  boneMass: number;
  waterPercentage: number;
  visceralFat: number;
  bmr: number;
  metabolicAge: number;
  notes: string | null;
}

interface ApiCompletedWorkout {
  id: string;
  workoutId: string;
  completedAt: string;
  durationSeconds: number;
}

interface ApiWorkoutStats {
  totalWorkoutsCompleted: number;
  workoutsThisWeek: number;
  totalMinutesSpent: number;
}

interface ApiWaterIntakeEntry {
  id: string;
  amountMl: number;
  consumedAt: string;
  note: string | null;
}

interface ApiDailyWaterSummary {
  date: string;
  totalMl: number;
  goalMl: number;
  entries: ApiWaterIntakeEntry[];
}

// Tipos de autenticação
interface ApiUser {
  id: string;
  email: string;
  name: string;
  profileId: string | null;
}

interface ApiAuthResponse {
  token: string;
  expiresAt: string;
  user: ApiUser;
}

export interface AuthUser {
  id: string;
  email: string;
  name: string;
  profileId: string | null;
}

export interface AuthResult {
  token: string;
  expiresAt: Date;
  user: AuthUser;
}

// Token storage
let currentToken: string | null = null;

export const setAuthToken = (token: string | null) => {
  currentToken = token;
};

export const getAuthToken = () => currentToken;

// Helper para fazer requisições
async function fetchApi<T>(
  endpoint: string,
  options?: RequestInit,
  requireAuth: boolean = true
): Promise<T> {
  const url = `${API_BASE_URL}${endpoint}`;

  const headers: HeadersInit = {
    "Content-Type": "application/json",
    ...options?.headers,
  };

  // Adicionar token de autenticação se disponível
  if (requireAuth && currentToken) {
    (headers as Record<string, string>)["Authorization"] = `Bearer ${currentToken}`;
  }

  const response = await fetch(url, {
    ...options,
    headers,
  });

  if (!response.ok) {
    const errorText = await response.text();
    throw new Error(`API Error ${response.status}: ${errorText}`);
  }

  // Se for DELETE com 204, não há body
  if (response.status === 204) {
    return undefined as T;
  }

  return response.json();
}

// Conversores de API para tipos locais
const convertApiProfile = (api: ApiUserProfile): UserProfile & { id: string } => ({
  id: api.id,
  name: api.name,
  age: api.age ?? undefined,
  height: api.height ?? undefined,
  currentWeight: api.currentWeight ?? undefined,
  goalWeight: api.goalWeight ?? undefined,
  avatarUrl: api.avatarUrl ?? undefined,
});

const convertApiExercise = (api: ApiExercise): Exercise => ({
  id: api.id,
  name: api.name,
  muscleGroup: api.muscleGroup as MuscleGroup,
  sets: api.sets,
  reps: api.reps,
  weight: api.weight,
  restSeconds: api.restSeconds,
  notes: api.notes ?? undefined,
});

const convertApiWorkout = (api: ApiWorkout): Workout => ({
  id: api.id,
  name: api.name,
  description: api.description ?? undefined,
  goal: api.goal as WorkoutGoal,
  days: api.days as DayOfWeek[],
  exercises: api.exercises.map(convertApiExercise),
  createdAt: api.createdAt,
  updatedAt: api.updatedAt,
  isCompleted: api.isCompleted,
  completedAt: api.completedAt ?? undefined,
});

const convertApiWaterEntry = (api: ApiWaterIntakeEntry): WaterIntakeEntry => ({
  id: api.id,
  amountMl: api.amountMl,
  consumedAt: api.consumedAt,
  note: api.note ?? undefined,
});

const convertApiDailyWaterSummary = (
  api: ApiDailyWaterSummary
): DailyWaterSummary => ({
  date: api.date,
  totalMl: api.totalMl,
  goalMl: api.goalMl,
  entries: api.entries.map(convertApiWaterEntry),
});

const convertApiBioimpedance = (api: ApiBioimpedance): BioimpedanceData => ({
  id: api.id,
  date: api.date,
  weight: api.weight,
  bodyFatPercentage: api.bodyFatPercentage,
  muscleMass: api.muscleMass,
  boneMass: api.boneMass,
  waterPercentage: api.waterPercentage,
  visceralFat: api.visceralFat,
  bmr: api.bmr,
  metabolicAge: api.metabolicAge,
  notes: api.notes ?? undefined,
});

// ============= Profile API =============

export const profileApi = {
  getAll: async (): Promise<(UserProfile & { id: string })[]> => {
    const data = await fetchApi<ApiUserProfile[]>(API_ENDPOINTS.profiles);
    return data.map(convertApiProfile);
  },

  getById: async (id: string): Promise<UserProfile & { id: string }> => {
    const data = await fetchApi<ApiUserProfile>(API_ENDPOINTS.profile(id));
    return convertApiProfile(data);
  },

  create: async (
    profile: Omit<UserProfile, "id">
  ): Promise<UserProfile & { id: string }> => {
    const data = await fetchApi<ApiUserProfile>(API_ENDPOINTS.profiles, {
      method: "POST",
      body: JSON.stringify({
        name: profile.name,
        age: profile.age ?? null,
        height: profile.height ?? null,
        currentWeight: profile.currentWeight ?? null,
        goalWeight: profile.goalWeight ?? null,
        avatarUrl: profile.avatarUrl ?? null,
      }),
    });
    return convertApiProfile(data);
  },

  update: async (
    id: string,
    profile: UserProfile
  ): Promise<UserProfile & { id: string }> => {
    const data = await fetchApi<ApiUserProfile>(API_ENDPOINTS.profile(id), {
      method: "PUT",
      body: JSON.stringify({
        name: profile.name,
        age: profile.age ?? null,
        height: profile.height ?? null,
        currentWeight: profile.currentWeight ?? null,
        goalWeight: profile.goalWeight ?? null,
        avatarUrl: profile.avatarUrl ?? null,
      }),
    });
    return convertApiProfile(data);
  },

  delete: async (id: string): Promise<void> => {
    await fetchApi<void>(API_ENDPOINTS.profile(id), {
      method: "DELETE",
    });
  },
};

// ============= Workouts API =============

export const workoutsApi = {
  getAll: async (profileId: string): Promise<Workout[]> => {
    const data = await fetchApi<ApiWorkout[]>(
      API_ENDPOINTS.workouts(profileId)
    );
    return data.map(convertApiWorkout);
  },

  getById: async (profileId: string, id: string): Promise<Workout> => {
    const data = await fetchApi<ApiWorkout>(
      API_ENDPOINTS.workout(profileId, id)
    );
    return convertApiWorkout(data);
  },

  getToday: async (profileId: string): Promise<Workout[]> => {
    const data = await fetchApi<ApiWorkout[]>(
      API_ENDPOINTS.todayWorkouts(profileId)
    );
    return data.map(convertApiWorkout);
  },

  create: async (
    profileId: string,
    workout: Omit<Workout, "id" | "createdAt" | "updatedAt">
  ): Promise<Workout> => {
    const data = await fetchApi<ApiWorkout>(API_ENDPOINTS.workouts(profileId), {
      method: "POST",
      body: JSON.stringify({
        name: workout.name,
        description: workout.description ?? null,
        goal: workout.goal,
        days: workout.days,
        exercises: workout.exercises.map((e) => ({
          name: e.name,
          muscleGroup: e.muscleGroup,
          sets: e.sets,
          reps: e.reps,
          weight: e.weight,
          restSeconds: e.restSeconds,
          notes: e.notes ?? null,
        })),
      }),
    });
    return convertApiWorkout(data);
  },

  update: async (profileId: string, workout: Workout): Promise<Workout> => {
    const data = await fetchApi<ApiWorkout>(
      API_ENDPOINTS.workout(profileId, workout.id),
      {
        method: "PUT",
        body: JSON.stringify({
          name: workout.name,
          description: workout.description ?? null,
          goal: workout.goal,
          days: workout.days,
          exercises: workout.exercises.map((e) => ({
            // Se o ID não for um GUID válido (ex: ID temporário local), enviar null
            // para que a API crie um novo ID
            id: isValidGuid(e.id) ? e.id : null,
            name: e.name,
            muscleGroup: e.muscleGroup,
            sets: e.sets,
            reps: e.reps,
            weight: e.weight,
            restSeconds: e.restSeconds,
            notes: e.notes ?? null,
          })),
        }),
      }
    );
    return convertApiWorkout(data);
  },

  delete: async (profileId: string, id: string): Promise<void> => {
    await fetchApi<void>(API_ENDPOINTS.workout(profileId, id), {
      method: "DELETE",
    });
  },
};

// ============= Water Intake API =============

export const waterApi = {
  getDay: async (
    profileId: string,
    date: string
  ): Promise<DailyWaterSummary> => {
    const data = await fetchApi<ApiDailyWaterSummary>(
      API_ENDPOINTS.waterByDate(profileId, date)
    );
    return convertApiDailyWaterSummary(data);
  },

  add: async (
    profileId: string,
    amountMl: number,
    consumedAt?: string
  ): Promise<WaterIntakeEntry> => {
    const data = await fetchApi<ApiWaterIntakeEntry>(
      API_ENDPOINTS.water(profileId),
      {
        method: "POST",
        body: JSON.stringify({
          amountMl,
          consumedAt: consumedAt ?? null,
        }),
      }
    );
    return convertApiWaterEntry(data);
  },

  deleteEntry: async (
    profileId: string,
    entryId: string
  ): Promise<void> => {
    await fetchApi<void>(
      API_ENDPOINTS.waterEntry(profileId, entryId),
      { method: "DELETE" }
    );
  },
};

// ============= Bioimpedance API =============

export const bioimpedanceApi = {
  getAll: async (profileId: string): Promise<BioimpedanceData[]> => {
    const data = await fetchApi<ApiBioimpedance[]>(
      API_ENDPOINTS.bioimpedance(profileId)
    );
    return data.map(convertApiBioimpedance);
  },

  getById: async (profileId: string, id: string): Promise<BioimpedanceData> => {
    const data = await fetchApi<ApiBioimpedance>(
      API_ENDPOINTS.bioimpedanceItem(profileId, id)
    );
    return convertApiBioimpedance(data);
  },

  getLatest: async (profileId: string): Promise<BioimpedanceData | null> => {
    try {
      const data = await fetchApi<ApiBioimpedance>(
        API_ENDPOINTS.latestBioimpedance(profileId)
      );
      return convertApiBioimpedance(data);
    } catch {
      return null;
    }
  },

  create: async (
    profileId: string,
    bioimpedance: Omit<BioimpedanceData, "id">
  ): Promise<BioimpedanceData> => {
    const data = await fetchApi<ApiBioimpedance>(
      API_ENDPOINTS.bioimpedance(profileId),
      {
        method: "POST",
        body: JSON.stringify({
          date: bioimpedance.date,
          weight: bioimpedance.weight,
          bodyFatPercentage: bioimpedance.bodyFatPercentage,
          muscleMass: bioimpedance.muscleMass,
          boneMass: bioimpedance.boneMass,
          waterPercentage: bioimpedance.waterPercentage,
          visceralFat: bioimpedance.visceralFat,
          bmr: bioimpedance.bmr,
          metabolicAge: bioimpedance.metabolicAge,
          notes: bioimpedance.notes ?? null,
        }),
      }
    );
    return convertApiBioimpedance(data);
  },

  update: async (
    profileId: string,
    bioimpedance: BioimpedanceData
  ): Promise<BioimpedanceData> => {
    const data = await fetchApi<ApiBioimpedance>(
      API_ENDPOINTS.bioimpedanceItem(profileId, bioimpedance.id),
      {
        method: "PUT",
        body: JSON.stringify({
          date: bioimpedance.date,
          weight: bioimpedance.weight,
          bodyFatPercentage: bioimpedance.bodyFatPercentage,
          muscleMass: bioimpedance.muscleMass,
          boneMass: bioimpedance.boneMass,
          waterPercentage: bioimpedance.waterPercentage,
          visceralFat: bioimpedance.visceralFat,
          bmr: bioimpedance.bmr,
          metabolicAge: bioimpedance.metabolicAge,
          notes: bioimpedance.notes ?? null,
        }),
      }
    );
    return convertApiBioimpedance(data);
  },

  delete: async (profileId: string, id: string): Promise<void> => {
    await fetchApi<void>(API_ENDPOINTS.bioimpedanceItem(profileId, id), {
      method: "DELETE",
    });
  },
};

// ============= Completed Workouts API =============

export interface CompletedWorkoutRecord {
  workoutId: string;
  completedAt: string;
  durationSeconds: number;
}

export const completedWorkoutsApi = {
  getAll: async (profileId: string): Promise<CompletedWorkoutRecord[]> => {
    const data = await fetchApi<ApiCompletedWorkout[]>(
      API_ENDPOINTS.completedWorkouts(profileId)
    );
    return data.map((c) => ({
      workoutId: c.workoutId,
      completedAt: c.completedAt,
      durationSeconds: c.durationSeconds,
    }));
  },

  getStats: async (profileId: string): Promise<ApiWorkoutStats> => {
    return fetchApi<ApiWorkoutStats>(
      API_ENDPOINTS.completedWorkoutStats(profileId)
    );
  },

  complete: async (
    profileId: string,
    workoutId: string,
    durationSeconds: number
  ): Promise<CompletedWorkoutRecord> => {
    const data = await fetchApi<ApiCompletedWorkout>(
      API_ENDPOINTS.completedWorkouts(profileId),
      {
        method: "POST",
        body: JSON.stringify({ workoutId, durationSeconds }),
      }
    );
    return {
      workoutId: data.workoutId,
      completedAt: data.completedAt,
      durationSeconds: data.durationSeconds,
    };
  },

  delete: async (profileId: string, id: string): Promise<void> => {
    await fetchApi<void>(API_ENDPOINTS.completedWorkout(profileId, id), {
      method: "DELETE",
    });
  },
};

// ============= Health Check =============

export const healthApi = {
  check: async (): Promise<{ status: string; timestamp: string }> => {
    return fetchApi<{ status: string; timestamp: string }>(
      API_ENDPOINTS.health,
      undefined,
      false
    );
  },
};

// ============= Auth API =============

export const authApi = {
  register: async (
    email: string,
    password: string,
    name: string
  ): Promise<AuthResult> => {
    const data = await fetchApi<ApiAuthResponse>(
      "/api/auth/register",
      {
        method: "POST",
        body: JSON.stringify({ email, password, name }),
      },
      false
    );
    return {
      token: data.token,
      expiresAt: new Date(data.expiresAt),
      user: {
        id: data.user.id,
        email: data.user.email,
        name: data.user.name,
        profileId: data.user.profileId,
      },
    };
  },

  login: async (email: string, password: string): Promise<AuthResult> => {
    const data = await fetchApi<ApiAuthResponse>(
      "/api/auth/login",
      {
        method: "POST",
        body: JSON.stringify({ email, password }),
      },
      false
    );
    return {
      token: data.token,
      expiresAt: new Date(data.expiresAt),
      user: {
        id: data.user.id,
        email: data.user.email,
        name: data.user.name,
        profileId: data.user.profileId,
      },
    };
  },

  getCurrentUser: async (): Promise<AuthUser> => {
    const data = await fetchApi<ApiUser>("/api/auth/me");
    return {
      id: data.id,
      email: data.email,
      name: data.name,
      profileId: data.profileId,
    };
  },

  refreshToken: async (): Promise<AuthResult> => {
    const data = await fetchApi<ApiAuthResponse>("/api/auth/refresh", {
      method: "POST",
    });
    return {
      token: data.token,
      expiresAt: new Date(data.expiresAt),
      user: {
        id: data.user.id,
        email: data.user.email,
        name: data.user.name,
        profileId: data.user.profileId,
      },
    };
  },

  changePassword: async (
    currentPassword: string,
    newPassword: string
  ): Promise<void> => {
    await fetchApi("/api/auth/change-password", {
      method: "POST",
      body: JSON.stringify({ currentPassword, newPassword }),
    });
  },
};

// ============= AI Planning Types =============

export interface GeneratePlanRequest {
  userProfileId: string;
  goal: WorkoutGoal;
  overrideDaysPerWeek?: number;
  overrideDurationMinutes?: number;
  overrideEquipment?: "Gym" | "Home" | "Minimal";
  additionalNotes?: string;
}

export interface GeneratedExercise {
  name: string;
  muscleGroup: MuscleGroup;
  sets: number;
  reps: number;
  suggestedWeight?: number;
  restSeconds: number;
  notes?: string;
}

export interface GeneratedWorkout {
  name: string;
  description?: string;
  goal: WorkoutGoal;
  days: DayOfWeek[];
  exercises: GeneratedExercise[];
  estimatedDurationMinutes: number;
}

export interface GeneratedPlanResponse {
  planId: string;
  summary: string;
  rationale: string;
  workouts: GeneratedWorkout[];
  generatedAt: string;
}

export interface AcceptPlanRequest {
  planId: string;
  userProfileId: string;
  workouts: GeneratedWorkout[];
}

export interface AcceptPlanResponse {
  createdWorkoutIds: string[];
  message: string;
}

export interface AIPlanningStatus {
  available: boolean;
  provider: string;
  model: string;
  endpoint: string;
}

// ============= AI Planning API =============

export const aiPlanningApi = {
  getStatus: async (): Promise<AIPlanningStatus> => {
    return fetchApi<AIPlanningStatus>(API_ENDPOINTS.aiPlanningStatus, undefined, false);
  },

  generate: async (request: GeneratePlanRequest): Promise<GeneratedPlanResponse> => {
    return fetchApi<GeneratedPlanResponse>(API_ENDPOINTS.aiPlanningGenerate, {
      method: "POST",
      body: JSON.stringify(request),
    });
  },

  preview: async (planId: string): Promise<GeneratedPlanResponse | null> => {
    try {
      return await fetchApi<GeneratedPlanResponse>(API_ENDPOINTS.aiPlanningPreview(planId));
    } catch {
      return null;
    }
  },

  accept: async (request: AcceptPlanRequest): Promise<AcceptPlanResponse> => {
    return fetchApi<AcceptPlanResponse>(API_ENDPOINTS.aiPlanningAccept, {
      method: "POST",
      body: JSON.stringify(request),
    });
  },
};
