// Days of the week
export type DayOfWeek =
  | "monday"
  | "tuesday"
  | "wednesday"
  | "thursday"
  | "friday"
  | "saturday"
  | "sunday";

export const DAYS_OF_WEEK: { key: DayOfWeek; label: string; short: string }[] = [
  { key: "monday", label: "Monday", short: "Mon" },
  { key: "tuesday", label: "Tuesday", short: "Tue" },
  { key: "wednesday", label: "Wednesday", short: "Wed" },
  { key: "thursday", label: "Thursday", short: "Thu" },
  { key: "friday", label: "Friday", short: "Fri" },
  { key: "saturday", label: "Saturday", short: "Sat" },
  { key: "sunday", label: "Sunday", short: "Sun" },
];

// Workout goals
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
  { key: "hypertrophy", label: "Hypertrophy", icon: "fitness-center" },
  { key: "strength", label: "Strength", icon: "flash-on" },
  { key: "endurance", label: "Endurance", icon: "directions-run" },
  { key: "weight_loss", label: "Weight Loss", icon: "trending-down" },
  { key: "maintenance", label: "Maintenance", icon: "balance" },
];

// Muscle groups
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
  { key: "chest", label: "Chest", color: "#FF6B6B" },
  { key: "back", label: "Back", color: "#4ECDC4" },
  { key: "shoulders", label: "Shoulders", color: "#45B7D1" },
  { key: "biceps", label: "Biceps", color: "#96CEB4" },
  { key: "triceps", label: "Triceps", color: "#FFEAA7" },
  { key: "legs", label: "Legs", color: "#DDA0DD" },
  { key: "glutes", label: "Glutes", color: "#FF8C94" },
  { key: "abs", label: "Abs", color: "#87CEEB" },
  { key: "cardio", label: "Cardio", color: "#FFB347" },
];

// Individual exercise
export interface Exercise {
  id?: string;
  name: string;
  muscleGroup?: MuscleGroup;
  sets: number;
  reps: number;
  weight?: number; // undefined = bodyweight
  restSeconds?: number;
  notes?: string;
  order?: number;
}

// Complete workout
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

// Bioimpedance data
export interface BioimpedanceData {
  id: string;
  date: string;
  weight: number; // kg
  bodyFatPercentage: number; // %
  muscleMass: number; // kg
  boneMass: number; // kg
  waterPercentage: number; // %
  visceralFat: number; // level
  bmr: number; // Basal Metabolic Rate (kcal)
  metabolicAge: number; // years
  notes?: string;
}

// User profile
export interface UserProfile {
  id?: string;
  name: string;
  age?: number;
  height?: number; // cm
  currentWeight?: number; // kg
  goalWeight?: number; // kg
  avatarUrl?: string;
}

// Completed workout record
export interface CompletedWorkoutRecord {
  id: string;
  workoutId: string;
  completedAt: string;
  durationSeconds: number;
}

// Auth user
export interface AuthUser {
  id: string;
  name: string;
  email: string;
  profileId: string | null;
}

// Auth result
export interface AuthResult {
  token: string;
  expiresAt: string;
  user: AuthUser;
}

// Workout stats
export interface WorkoutStats {
  totalWorkoutsCompleted: number;
  workoutsThisWeek: number;
  totalMinutesSpent: number;
}

// Water intake
export interface WaterIntakeEntry {
  id: string;
  amountMl: number;
  consumedAt: string;
  note?: string;
}

export interface DailyWaterSummary {
  date: string;
  totalMl: number;
  goalMl: number;
  entries: WaterIntakeEntry[];
}

// App state
export interface AppState {
  profile: UserProfile;
  workouts: Workout[];
  bioimpedanceHistory: BioimpedanceData[];
  completedWorkouts: CompletedWorkoutRecord[];
}
