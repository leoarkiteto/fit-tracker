// Dias da semana
export type DayOfWeek =
  | "monday"
  | "tuesday"
  | "wednesday"
  | "thursday"
  | "friday"
  | "saturday"
  | "sunday";

export const DAYS_OF_WEEK: { key: DayOfWeek; label: string; short: string }[] =
  [
    { key: "monday", label: "Segunda-feira", short: "Seg" },
    { key: "tuesday", label: "Terça-feira", short: "Ter" },
    { key: "wednesday", label: "Quarta-feira", short: "Qua" },
    { key: "thursday", label: "Quinta-feira", short: "Qui" },
    { key: "friday", label: "Sexta-feira", short: "Sex" },
    { key: "saturday", label: "Sábado", short: "Sáb" },
    { key: "sunday", label: "Domingo", short: "Dom" },
  ];

// Objetivos de treino
export type WorkoutGoal =
  | "hypertrophy"
  | "strength"
  | "endurance"
  | "weight_loss"
  | "maintenance";

export const WORKOUT_GOALS: {
  key: WorkoutGoal;
  label: string;
  icon: string;
}[] = [
  { key: "hypertrophy", label: "Hipertrofia", icon: "fitness-center" },
  { key: "strength", label: "Força", icon: "flash-on" },
  { key: "endurance", label: "Resistência", icon: "directions-run" },
  { key: "weight_loss", label: "Emagrecimento", icon: "trending-down" },
  { key: "maintenance", label: "Manutenção", icon: "balance" },
];

// Grupos musculares
export type MuscleGroup =
  | "chest"
  | "back"
  | "shoulders"
  | "biceps"
  | "triceps"
  | "legs"
  | "glutes"
  | "abs"
  | "cardio";

export const MUSCLE_GROUPS: {
  key: MuscleGroup;
  label: string;
  color: string;
}[] = [
  { key: "chest", label: "Peito", color: "#FF6B6B" },
  { key: "back", label: "Costas", color: "#4ECDC4" },
  { key: "shoulders", label: "Ombros", color: "#45B7D1" },
  { key: "biceps", label: "Bíceps", color: "#96CEB4" },
  { key: "triceps", label: "Tríceps", color: "#FFEAA7" },
  { key: "legs", label: "Pernas", color: "#DDA0DD" },
  { key: "glutes", label: "Glúteos", color: "#FF8C94" },
  { key: "abs", label: "Abdômen", color: "#87CEEB" },
  { key: "cardio", label: "Cardio", color: "#FFB347" },
];

// Exercício individual
export interface Exercise {
  id: string;
  name: string;
  muscleGroup: MuscleGroup;
  sets: number;
  reps: number;
  weight: number | null; // null = peso livre/corporal
  restSeconds: number;
  notes?: string;
}

// Treino completo
export interface Workout {
  id: string;
  name: string;
  description?: string;
  goal: WorkoutGoal;
  days: DayOfWeek[];
  exercises: Exercise[];
  createdAt: string;
  updatedAt: string;
  isCompleted?: boolean;
  completedAt?: string;
}

// Dados de Bioimpedância
export interface BioimpedanceData {
  id: string;
  date: string;
  weight: number; // kg
  bodyFatPercentage: number; // %
  muscleMass: number; // kg
  boneMass: number; // kg
  waterPercentage: number; // %
  visceralFat: number; // nível
  bmr: number; // Taxa Metabólica Basal (kcal)
  metabolicAge: number; // anos
  notes?: string;
}

// Perfil do usuário
export interface UserProfile {
  id?: string; // ID do perfil na API (undefined se ainda não sincronizado)
  name: string;
  age?: number;
  height?: number; // cm
  currentWeight?: number; // kg
  goalWeight?: number; // kg
  avatarUrl?: string;
}

// Estado da aplicação
export interface AppState {
  profile: UserProfile;
  workouts: Workout[];
  bioimpedanceHistory: BioimpedanceData[];
  completedWorkouts: { workoutId: string; completedAt: string }[];
}
