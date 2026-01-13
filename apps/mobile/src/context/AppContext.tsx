import React, {
  createContext,
  useContext,
  useState,
  useEffect,
  ReactNode,
  useCallback,
} from "react";
import { Workout, BioimpedanceData, UserProfile } from "../types";
import {
  profileApi,
  workoutsApi,
  bioimpedanceApi,
  completedWorkoutsApi,
  CompletedWorkoutRecord,
} from "../services";
import { useAuth } from "./AuthContext";

interface AppContextType {
  // Profile
  profile: UserProfile;
  profileId: string | null;
  updateProfile: (profile: UserProfile) => Promise<void>;

  // Workouts
  workouts: Workout[];
  addWorkout: (
    workout: Omit<Workout, "id" | "createdAt" | "updatedAt">
  ) => Promise<void>;
  updateWorkout: (workout: Workout) => Promise<void>;
  deleteWorkout: (id: string) => Promise<void>;
  completeWorkout: (id: string, durationSeconds: number) => Promise<void>;

  // Bioimpedance
  bioimpedanceHistory: BioimpedanceData[];
  addBioimpedance: (data: Omit<BioimpedanceData, "id">) => Promise<void>;
  updateBioimpedance: (data: BioimpedanceData) => Promise<void>;
  deleteBioimpedance: (id: string) => Promise<void>;

  // Completed workouts
  completedWorkouts: CompletedWorkoutRecord[];
  isWorkoutCompletedToday: (workoutId: string) => boolean;

  // Stats
  totalWorkoutsCompleted: number;
  workoutsThisWeek: number;
  totalMinutesSpent: number;

  // Loading state
  isLoading: boolean;

  // Error state
  error: string | null;
  clearError: () => void;

  // Refresh data
  refreshData: () => Promise<void>;
}

const AppContext = createContext<AppContextType | undefined>(undefined);

export const AppProvider: React.FC<{ children: ReactNode }> = ({
  children,
}) => {
  const { user, isAuthenticated } = useAuth();
  
  const [profileId, setProfileId] = useState<string | null>(null);
  const [profile, setProfile] = useState<UserProfile>({ name: "Atleta" });
  const [workouts, setWorkouts] = useState<Workout[]>([]);
  const [bioimpedanceHistory, setBioimpedanceHistory] = useState<
    BioimpedanceData[]
  >([]);
  const [completedWorkouts, setCompletedWorkouts] = useState<
    CompletedWorkoutRecord[]
  >([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // Stats from API
  const [stats, setStats] = useState({
    totalWorkoutsCompleted: 0,
    workoutsThisWeek: 0,
    totalMinutesSpent: 0,
  });

  const clearError = useCallback(() => setError(null), []);

  // Initialize: Load profile data when user is authenticated
  useEffect(() => {
    if (isAuthenticated && user?.profileId) {
      initializeApp(user.profileId);
    } else {
      // Reset state when user logs out
      setProfileId(null);
      setProfile({ name: "Atleta" });
      setWorkouts([]);
      setBioimpedanceHistory([]);
      setCompletedWorkouts([]);
      setStats({ totalWorkoutsCompleted: 0, workoutsThisWeek: 0, totalMinutesSpent: 0 });
      setIsLoading(false);
    }
  }, [isAuthenticated, user?.profileId]);

  const initializeApp = async (authProfileId: string) => {
    try {
      setIsLoading(true);
      setError(null);

      // Load profile from authenticated user
      try {
        const loadedProfile = await profileApi.getById(authProfileId);
        setProfileId(authProfileId);
        setProfile(loadedProfile);
        await loadProfileData(authProfileId);
      } catch (err) {
        console.error("Error loading profile:", err);
        setError("Erro ao carregar perfil. Tente novamente.");
      }
    } catch (err) {
      console.error("Error initializing app:", err);
      setError(
        "Erro ao conectar com o servidor. Verifique sua conexão e tente novamente."
      );
    } finally {
      setIsLoading(false);
    }
  };

  const loadProfileData = async (id: string) => {
    try {
      const [loadedWorkouts, loadedBioimpedance, loadedCompleted, loadedStats] =
        await Promise.all([
          workoutsApi.getAll(id),
          bioimpedanceApi.getAll(id),
          completedWorkoutsApi.getAll(id),
          completedWorkoutsApi.getStats(id),
        ]);

      setWorkouts(loadedWorkouts);
      setBioimpedanceHistory(loadedBioimpedance);
      setCompletedWorkouts(loadedCompleted);
      setStats(loadedStats);
    } catch (err) {
      console.error("Error loading profile data:", err);
      throw err;
    }
  };

  const refreshData = useCallback(async () => {
    if (!profileId) return;

    try {
      setError(null);
      await loadProfileData(profileId);
    } catch (err) {
      console.error("Error refreshing data:", err);
      setError("Erro ao atualizar dados. Tente novamente.");
    }
  }, [profileId]);

  // Profile
  const updateProfile = async (newProfile: UserProfile) => {
    if (!profileId) return;

    try {
      setError(null);
      const updatedProfile = await profileApi.update(profileId, newProfile);
      setProfile(updatedProfile);
    } catch (err) {
      console.error("Error updating profile:", err);
      setError("Erro ao atualizar perfil. Tente novamente.");
      throw err;
    }
  };

  // Workouts
  const addWorkout = async (
    workoutData: Omit<Workout, "id" | "createdAt" | "updatedAt">
  ) => {
    if (!profileId) return;

    try {
      setError(null);
      const newWorkout = await workoutsApi.create(profileId, workoutData);
      setWorkouts((prev) => [newWorkout, ...prev]);
    } catch (err) {
      console.error("Error adding workout:", err);
      setError("Erro ao criar treino. Tente novamente.");
      throw err;
    }
  };

  const updateWorkout = async (workout: Workout) => {
    if (!profileId) return;

    try {
      setError(null);
      const updatedWorkout = await workoutsApi.update(profileId, workout);
      setWorkouts((prev) =>
        prev.map((w) => (w.id === workout.id ? updatedWorkout : w))
      );
    } catch (err) {
      console.error("Error updating workout:", err);
      setError("Erro ao atualizar treino. Tente novamente.");
      throw err;
    }
  };

  const deleteWorkout = async (id: string) => {
    if (!profileId) return;

    try {
      setError(null);
      await workoutsApi.delete(profileId, id);
      setWorkouts((prev) => prev.filter((w) => w.id !== id));
    } catch (err) {
      console.error("Error deleting workout:", err);
      setError("Erro ao excluir treino. Tente novamente.");
      throw err;
    }
  };

  const completeWorkout = async (id: string, durationSeconds: number) => {
    if (!profileId) return;

    try {
      setError(null);
      const completed = await completedWorkoutsApi.complete(
        profileId,
        id,
        durationSeconds
      );
      setCompletedWorkouts((prev) => [completed, ...prev]);

      // Update stats
      const newStats = await completedWorkoutsApi.getStats(profileId);
      setStats(newStats);
    } catch (err) {
      console.error("Error completing workout:", err);
      setError("Erro ao marcar treino como concluído. Tente novamente.");
      throw err;
    }
  };

  // Verifica se um treino foi completado hoje
  const isWorkoutCompletedToday = useCallback(
    (workoutId: string): boolean => {
      const today = new Date();
      today.setHours(0, 0, 0, 0);

      return completedWorkouts.some((cw) => {
        if (cw.workoutId !== workoutId) return false;
        const completedDate = new Date(cw.completedAt);
        completedDate.setHours(0, 0, 0, 0);
        return completedDate.getTime() === today.getTime();
      });
    },
    [completedWorkouts]
  );

  // Bioimpedance
  const addBioimpedance = async (data: Omit<BioimpedanceData, "id">) => {
    if (!profileId) return;

    try {
      setError(null);
      const newData = await bioimpedanceApi.create(profileId, data);
      setBioimpedanceHistory((prev) => [newData, ...prev]);
    } catch (err) {
      console.error("Error adding bioimpedance:", err);
      setError("Erro ao salvar bioimpedância. Tente novamente.");
      throw err;
    }
  };

  const updateBioimpedance = async (data: BioimpedanceData) => {
    if (!profileId) return;

    try {
      setError(null);
      const updatedData = await bioimpedanceApi.update(profileId, data);
      setBioimpedanceHistory((prev) =>
        prev.map((b) => (b.id === data.id ? updatedData : b))
      );
    } catch (err) {
      console.error("Error updating bioimpedance:", err);
      setError("Erro ao atualizar bioimpedância. Tente novamente.");
      throw err;
    }
  };

  const deleteBioimpedance = async (id: string) => {
    if (!profileId) return;

    try {
      setError(null);
      await bioimpedanceApi.delete(profileId, id);
      setBioimpedanceHistory((prev) => prev.filter((b) => b.id !== id));
    } catch (err) {
      console.error("Error deleting bioimpedance:", err);
      setError("Erro ao excluir bioimpedância. Tente novamente.");
      throw err;
    }
  };

  return (
    <AppContext.Provider
      value={{
        profile,
        profileId,
        updateProfile,
        workouts,
        addWorkout,
        updateWorkout,
        deleteWorkout,
        completeWorkout,
        bioimpedanceHistory,
        addBioimpedance,
        updateBioimpedance,
        deleteBioimpedance,
        completedWorkouts,
        isWorkoutCompletedToday,
        totalWorkoutsCompleted: stats.totalWorkoutsCompleted,
        workoutsThisWeek: stats.workoutsThisWeek,
        totalMinutesSpent: stats.totalMinutesSpent,
        isLoading,
        error,
        clearError,
        refreshData,
      }}
    >
      {children}
    </AppContext.Provider>
  );
};

export const useApp = (): AppContextType => {
  const context = useContext(AppContext);
  if (!context) {
    throw new Error("useApp must be used within an AppProvider");
  }
  return context;
};
