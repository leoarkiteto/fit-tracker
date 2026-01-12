import React, { useState, useEffect, useRef } from "react";
import {
  View,
  Text,
  StyleSheet,
  TouchableOpacity,
  ScrollView,
  Vibration,
  Alert,
  StatusBar,
} from "react-native";
import { SafeAreaView } from "react-native-safe-area-context";
import { Ionicons } from "@expo/vector-icons";
import { useNavigation, useRoute, RouteProp } from "@react-navigation/native";
import { NativeStackNavigationProp } from "@react-navigation/native-stack";
import { LinearGradient } from "expo-linear-gradient";

import { useApp } from "../context/AppContext";
import { colors, spacing, borderRadius, typography } from "../theme";
import { RootStackParamList } from "../navigation/types";
import { Exercise, MUSCLE_GROUPS } from "../types";

type WorkoutSessionScreenNavigationProp = NativeStackNavigationProp<
  RootStackParamList,
  "WorkoutSession"
>;
type WorkoutSessionScreenRouteProp = RouteProp<
  RootStackParamList,
  "WorkoutSession"
>;

type SessionPhase = "selecting" | "exercising" | "resting" | "completed";

interface ExerciseProgress {
  exerciseId: string;
  completedSets: number;
  totalSets: number;
}

export const WorkoutSessionScreen: React.FC = () => {
  const navigation = useNavigation<WorkoutSessionScreenNavigationProp>();
  const route = useRoute<WorkoutSessionScreenRouteProp>();
  const { workouts, completeWorkout } = useApp();

  const workout = workouts.find((w) => w.id === route.params.workoutId);

  // Estados da sessão
  const [sessionPhase, setSessionPhase] = useState<SessionPhase>("selecting");
  const [workoutTimer, setWorkoutTimer] = useState(0); // Tempo total do treino em segundos
  const [restTimer, setRestTimer] = useState(0); // Tempo de descanso restante
  const [activeExercise, setActiveExercise] = useState<Exercise | null>(null);
  const [currentSet, setCurrentSet] = useState(1);
  const [exerciseProgress, setExerciseProgress] = useState<ExerciseProgress[]>(
    []
  );
  const [isWorkoutStarted, setIsWorkoutStarted] = useState(false);

  const workoutTimerRef = useRef<ReturnType<typeof setInterval> | null>(null);
  const restTimerRef = useRef<ReturnType<typeof setInterval> | null>(null);

  // Inicializar progresso dos exercícios
  useEffect(() => {
    if (workout) {
      setExerciseProgress(
        workout.exercises.map((e) => ({
          exerciseId: e.id,
          completedSets: 0,
          totalSets: e.sets,
        }))
      );
    }
  }, [workout]);

  // Timer do treino (cronômetro geral)
  useEffect(() => {
    if (isWorkoutStarted && sessionPhase !== "completed") {
      workoutTimerRef.current = setInterval(() => {
        setWorkoutTimer((prev) => prev + 1);
      }, 1000);
    }

    return () => {
      if (workoutTimerRef.current) {
        clearInterval(workoutTimerRef.current);
      }
    };
  }, [isWorkoutStarted, sessionPhase]);

  // Timer de descanso
  useEffect(() => {
    if (sessionPhase === "resting" && restTimer > 0) {
      restTimerRef.current = setInterval(() => {
        setRestTimer((prev) => {
          if (prev <= 1) {
            // Tempo de descanso acabou
            clearInterval(restTimerRef.current!);
            handleRestComplete();
            return 0;
          }
          return prev - 1;
        });
      }, 1000);
    }

    return () => {
      if (restTimerRef.current) {
        clearInterval(restTimerRef.current);
      }
    };
  }, [sessionPhase, restTimer]);

  // Vibrar quando descanso terminar
  const alertRestComplete = () => {
    Vibration.vibrate([0, 500, 200, 500, 200, 500]);
  };

  const formatTime = (seconds: number): string => {
    const mins = Math.floor(seconds / 60);
    const secs = seconds % 60;
    return `${mins.toString().padStart(2, "0")}:${secs.toString().padStart(2, "0")}`;
  };

  const handleStartWorkout = () => {
    setIsWorkoutStarted(true);
    setSessionPhase("selecting");
  };

  const handleSelectExercise = (exercise: Exercise) => {
    const progress = exerciseProgress.find((p) => p.exerciseId === exercise.id);
    if (progress && progress.completedSets >= progress.totalSets) {
      Alert.alert("Exercício Completo", "Este exercício já foi finalizado!");
      return;
    }

    setActiveExercise(exercise);
    setCurrentSet((progress?.completedSets || 0) + 1);
    setSessionPhase("exercising");
  };

  const handleCompletedRep = () => {
    if (!activeExercise) return;

    // Atualizar progresso
    setExerciseProgress((prev) =>
      prev.map((p) =>
        p.exerciseId === activeExercise.id
          ? { ...p, completedSets: p.completedSets + 1 }
          : p
      )
    );

    const progress = exerciseProgress.find(
      (p) => p.exerciseId === activeExercise.id
    );
    const newCompletedSets = (progress?.completedSets || 0) + 1;

    if (newCompletedSets >= activeExercise.sets) {
      // Exercício completo
      checkWorkoutComplete();
    } else {
      // Iniciar descanso
      setRestTimer(activeExercise.restSeconds);
      setSessionPhase("resting");
    }
  };

  const handleRestComplete = () => {
    alertRestComplete();

    if (!activeExercise) return;

    const progress = exerciseProgress.find(
      (p) => p.exerciseId === activeExercise.id
    );

    if (progress && progress.completedSets >= activeExercise.sets) {
      // Exercício completo, voltar para seleção
      checkWorkoutComplete();
    } else {
      // Próxima série
      setCurrentSet((progress?.completedSets || 0) + 1);
      setSessionPhase("exercising");
    }
  };

  const checkWorkoutComplete = () => {
    const allComplete = exerciseProgress.every(
      (p) => p.completedSets >= p.totalSets
    );

    if (allComplete) {
      setSessionPhase("completed");
    } else {
      setActiveExercise(null);
      setSessionPhase("selecting");
    }
  };

  const handleSkipRest = () => {
    if (restTimerRef.current) {
      clearInterval(restTimerRef.current);
    }
    handleRestComplete();
  };

  const handleFinishWorkout = async () => {
    if (workout) {
      await completeWorkout(workout.id);
    }
    navigation.goBack();
  };

  const handleExitWorkout = () => {
    Alert.alert(
      "Sair do Treino",
      "Tem certeza que deseja sair? O progresso será perdido.",
      [
        { text: "Continuar", style: "cancel" },
        {
          text: "Sair",
          style: "destructive",
          onPress: () => navigation.goBack(),
        },
      ]
    );
  };

  const getMuscleGroupColor = (muscleGroup: string) => {
    const group = MUSCLE_GROUPS.find((g) => g.key === muscleGroup);
    return group?.color || colors.primary;
  };

  if (!workout) {
    return (
      <SafeAreaView style={styles.container}>
        <Text style={styles.errorText}>Treino não encontrado</Text>
      </SafeAreaView>
    );
  }

  // Tela de exercício ativo
  if (sessionPhase === "exercising" && activeExercise) {
    return (
      <TouchableOpacity
        style={styles.fullScreenContainer}
        activeOpacity={1}
        onPress={handleCompletedRep}
      >
        <StatusBar barStyle="light-content" />
        <LinearGradient
          colors={[getMuscleGroupColor(activeExercise.muscleGroup), "#1a1a2e"]}
          style={styles.exerciseScreen}
        >
          {/* Timer do treino no topo */}
          <View style={styles.topTimer}>
            <Ionicons name="time-outline" size={20} color="rgba(255,255,255,0.7)" />
            <Text style={styles.topTimerText}>{formatTime(workoutTimer)}</Text>
          </View>

          <View style={styles.exerciseContent}>
            <Text style={styles.exerciseLabel}>EXERCÍCIO</Text>
            <Text style={styles.exerciseName}>{activeExercise.name}</Text>

            <View style={styles.setInfo}>
              <Text style={styles.setNumber}>Série {currentSet}</Text>
              <Text style={styles.setTotal}>de {activeExercise.sets}</Text>
            </View>

            <View style={styles.repInfo}>
              <Text style={styles.repNumber}>{activeExercise.reps}</Text>
              <Text style={styles.repLabel}>repetições</Text>
            </View>

            {activeExercise.weight && (
              <View style={styles.weightInfo}>
                <Ionicons name="barbell" size={24} color="rgba(255,255,255,0.8)" />
                <Text style={styles.weightText}>{activeExercise.weight} kg</Text>
              </View>
            )}
          </View>

          <View style={styles.tapHint}>
            <Ionicons name="hand-left-outline" size={32} color="rgba(255,255,255,0.5)" />
            <Text style={styles.tapHintText}>
              Toque na tela ao completar a série
            </Text>
          </View>

          <TouchableOpacity style={styles.exitButton} onPress={handleExitWorkout}>
            <Ionicons name="close" size={24} color="rgba(255,255,255,0.7)" />
          </TouchableOpacity>
        </LinearGradient>
      </TouchableOpacity>
    );
  }

  // Tela de descanso
  if (sessionPhase === "resting" && activeExercise) {
    const progress = (restTimer / activeExercise.restSeconds) * 100;

    return (
      <View style={styles.fullScreenContainer}>
        <StatusBar barStyle="light-content" />
        <LinearGradient
          colors={["#2d3436", "#1a1a2e"]}
          style={styles.restScreen}
        >
          {/* Timer do treino no topo */}
          <View style={styles.topTimer}>
            <Ionicons name="time-outline" size={20} color="rgba(255,255,255,0.7)" />
            <Text style={styles.topTimerText}>{formatTime(workoutTimer)}</Text>
          </View>

          <View style={styles.restContent}>
            <Text style={styles.restLabel}>DESCANSO</Text>

            <View style={styles.restTimerContainer}>
              <Text style={styles.restTimerText}>{formatTime(restTimer)}</Text>
              <View style={styles.progressBar}>
                <View
                  style={[styles.progressFill, { width: `${progress}%` }]}
                />
              </View>
            </View>

            <Text style={styles.nextSetText}>
              Próxima: Série {currentSet + 1} de {activeExercise.sets}
            </Text>

            <TouchableOpacity style={styles.skipButton} onPress={handleSkipRest}>
              <Ionicons name="play-forward" size={20} color={colors.white} />
              <Text style={styles.skipButtonText}>Pular Descanso</Text>
            </TouchableOpacity>
          </View>

          <TouchableOpacity style={styles.exitButton} onPress={handleExitWorkout}>
            <Ionicons name="close" size={24} color="rgba(255,255,255,0.7)" />
          </TouchableOpacity>
        </LinearGradient>
      </View>
    );
  }

  // Tela de treino completo
  if (sessionPhase === "completed") {
    return (
      <View style={styles.fullScreenContainer}>
        <StatusBar barStyle="light-content" />
        <LinearGradient
          colors={[colors.success, "#1a1a2e"]}
          style={styles.completedScreen}
        >
          <View style={styles.completedContent}>
            <Ionicons name="trophy" size={80} color={colors.white} />
            <Text style={styles.completedTitle}>Treino Concluído!</Text>
            <Text style={styles.completedTime}>
              Tempo total: {formatTime(workoutTimer)}
            </Text>

            <View style={styles.completedStats}>
              <View style={styles.completedStat}>
                <Text style={styles.completedStatNumber}>
                  {workout.exercises.length}
                </Text>
                <Text style={styles.completedStatLabel}>Exercícios</Text>
              </View>
              <View style={styles.completedStat}>
                <Text style={styles.completedStatNumber}>
                  {workout.exercises.reduce((acc, e) => acc + e.sets, 0)}
                </Text>
                <Text style={styles.completedStatLabel}>Séries</Text>
              </View>
            </View>

            <TouchableOpacity
              style={styles.finishButton}
              onPress={handleFinishWorkout}
            >
              <Text style={styles.finishButtonText}>Finalizar</Text>
            </TouchableOpacity>
          </View>
        </LinearGradient>
      </View>
    );
  }

  // Tela de seleção de exercícios (padrão)
  return (
    <SafeAreaView style={styles.container} edges={["top"]}>
      <StatusBar barStyle="dark-content" />
      
      {/* Header */}
      <View style={styles.header}>
        <TouchableOpacity onPress={handleExitWorkout} style={styles.backButton}>
          <Ionicons name="arrow-back" size={24} color={colors.textPrimary} />
        </TouchableOpacity>
        <Text style={styles.headerTitle} numberOfLines={1}>
          {workout.name}
        </Text>
        <View style={styles.timerBadge}>
          <Ionicons name="time-outline" size={16} color={colors.primary} />
          <Text style={styles.timerBadgeText}>{formatTime(workoutTimer)}</Text>
        </View>
      </View>

      {/* Botão Iniciar Treino */}
      {!isWorkoutStarted && (
        <View style={styles.startContainer}>
          <TouchableOpacity
            style={styles.startButton}
            onPress={handleStartWorkout}
          >
            <LinearGradient
              colors={[colors.primary, colors.primaryLight]}
              style={styles.startButtonGradient}
            >
              <Ionicons name="play" size={32} color={colors.white} />
              <Text style={styles.startButtonText}>Iniciar Treino</Text>
            </LinearGradient>
          </TouchableOpacity>
        </View>
      )}

      {/* Lista de exercícios */}
      {isWorkoutStarted && (
        <>
          <Text style={styles.sectionTitle}>Selecione um exercício</Text>
          <ScrollView
            style={styles.exerciseList}
            contentContainerStyle={styles.exerciseListContent}
          >
            {workout.exercises.map((exercise, index) => {
              const progress = exerciseProgress.find(
                (p) => p.exerciseId === exercise.id
              );
              const isCompleted =
                progress && progress.completedSets >= progress.totalSets;
              const muscleGroup = MUSCLE_GROUPS.find(
                (g) => g.key === exercise.muscleGroup
              );

              return (
                <TouchableOpacity
                  key={exercise.id}
                  style={[
                    styles.exerciseItem,
                    isCompleted && styles.exerciseItemCompleted,
                  ]}
                  onPress={() => handleSelectExercise(exercise)}
                  disabled={isCompleted}
                >
                  <View
                    style={[
                      styles.exerciseNumber,
                      { backgroundColor: muscleGroup?.color || colors.primary },
                      isCompleted && styles.exerciseNumberCompleted,
                    ]}
                  >
                    {isCompleted ? (
                      <Ionicons name="checkmark" size={20} color={colors.white} />
                    ) : (
                      <Text style={styles.exerciseNumberText}>{index + 1}</Text>
                    )}
                  </View>

                  <View style={styles.exerciseItemContent}>
                    <Text
                      style={[
                        styles.exerciseItemName,
                        isCompleted && styles.exerciseItemNameCompleted,
                      ]}
                    >
                      {exercise.name}
                    </Text>
                    <Text style={styles.exerciseItemDetails}>
                      {progress?.completedSets || 0}/{exercise.sets} séries •{" "}
                      {exercise.reps} reps
                      {exercise.weight && ` • ${exercise.weight}kg`}
                    </Text>
                  </View>

                  {!isCompleted && (
                    <Ionicons
                      name="chevron-forward"
                      size={20}
                      color={colors.textLight}
                    />
                  )}
                </TouchableOpacity>
              );
            })}
          </ScrollView>
        </>
      )}
    </SafeAreaView>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: colors.background,
  },
  fullScreenContainer: {
    flex: 1,
  },
  header: {
    flexDirection: "row",
    alignItems: "center",
    paddingHorizontal: spacing.md,
    paddingVertical: spacing.md,
    borderBottomWidth: 1,
    borderBottomColor: colors.border,
    backgroundColor: colors.surface,
  },
  backButton: {
    width: 40,
    height: 40,
    alignItems: "center",
    justifyContent: "center",
  },
  headerTitle: {
    flex: 1,
    fontSize: typography.lg,
    fontWeight: typography.bold,
    color: colors.textPrimary,
    marginHorizontal: spacing.sm,
  },
  timerBadge: {
    flexDirection: "row",
    alignItems: "center",
    gap: spacing.xs,
    backgroundColor: `${colors.primary}15`,
    paddingHorizontal: spacing.md,
    paddingVertical: spacing.sm,
    borderRadius: borderRadius.full,
  },
  timerBadgeText: {
    fontSize: typography.md,
    fontWeight: typography.semibold,
    color: colors.primary,
  },
  startContainer: {
    flex: 1,
    alignItems: "center",
    justifyContent: "center",
    padding: spacing.xl,
  },
  startButton: {
    borderRadius: borderRadius.xl,
    overflow: "hidden",
  },
  startButtonGradient: {
    flexDirection: "row",
    alignItems: "center",
    gap: spacing.md,
    paddingVertical: spacing.lg,
    paddingHorizontal: spacing.xxl,
  },
  startButtonText: {
    fontSize: typography.xl,
    fontWeight: typography.bold,
    color: colors.white,
  },
  sectionTitle: {
    fontSize: typography.md,
    fontWeight: typography.semibold,
    color: colors.textSecondary,
    paddingHorizontal: spacing.lg,
    paddingTop: spacing.lg,
    paddingBottom: spacing.md,
  },
  exerciseList: {
    flex: 1,
  },
  exerciseListContent: {
    paddingHorizontal: spacing.lg,
    paddingBottom: spacing.xxl,
  },
  exerciseItem: {
    flexDirection: "row",
    alignItems: "center",
    backgroundColor: colors.surface,
    borderRadius: borderRadius.lg,
    padding: spacing.md,
    marginBottom: spacing.sm,
  },
  exerciseItemCompleted: {
    opacity: 0.6,
  },
  exerciseNumber: {
    width: 40,
    height: 40,
    borderRadius: borderRadius.full,
    alignItems: "center",
    justifyContent: "center",
  },
  exerciseNumberCompleted: {
    backgroundColor: colors.success,
  },
  exerciseNumberText: {
    fontSize: typography.md,
    fontWeight: typography.bold,
    color: colors.white,
  },
  exerciseItemContent: {
    flex: 1,
    marginLeft: spacing.md,
  },
  exerciseItemName: {
    fontSize: typography.md,
    fontWeight: typography.semibold,
    color: colors.textPrimary,
  },
  exerciseItemNameCompleted: {
    textDecorationLine: "line-through",
  },
  exerciseItemDetails: {
    fontSize: typography.sm,
    color: colors.textSecondary,
    marginTop: spacing.xs,
  },
  // Exercise Screen
  exerciseScreen: {
    flex: 1,
    justifyContent: "center",
    alignItems: "center",
  },
  topTimer: {
    position: "absolute",
    top: 60,
    flexDirection: "row",
    alignItems: "center",
    gap: spacing.xs,
  },
  topTimerText: {
    fontSize: typography.md,
    color: "rgba(255,255,255,0.7)",
    fontWeight: typography.medium,
  },
  exerciseContent: {
    alignItems: "center",
  },
  exerciseLabel: {
    fontSize: typography.sm,
    color: "rgba(255,255,255,0.6)",
    letterSpacing: 2,
    marginBottom: spacing.sm,
  },
  exerciseName: {
    fontSize: typography.xxxl,
    fontWeight: typography.bold,
    color: colors.white,
    textAlign: "center",
    marginBottom: spacing.xl,
  },
  setInfo: {
    flexDirection: "row",
    alignItems: "baseline",
    gap: spacing.sm,
    marginBottom: spacing.lg,
  },
  setNumber: {
    fontSize: typography.xxl,
    fontWeight: typography.bold,
    color: colors.white,
  },
  setTotal: {
    fontSize: typography.lg,
    color: "rgba(255,255,255,0.7)",
  },
  repInfo: {
    alignItems: "center",
    marginBottom: spacing.lg,
  },
  repNumber: {
    fontSize: 80,
    fontWeight: typography.bold,
    color: colors.white,
  },
  repLabel: {
    fontSize: typography.lg,
    color: "rgba(255,255,255,0.7)",
  },
  weightInfo: {
    flexDirection: "row",
    alignItems: "center",
    gap: spacing.sm,
  },
  weightText: {
    fontSize: typography.xl,
    color: "rgba(255,255,255,0.8)",
    fontWeight: typography.semibold,
  },
  tapHint: {
    position: "absolute",
    bottom: 80,
    alignItems: "center",
    gap: spacing.sm,
  },
  tapHintText: {
    fontSize: typography.md,
    color: "rgba(255,255,255,0.5)",
  },
  exitButton: {
    position: "absolute",
    top: 50,
    right: 20,
    width: 44,
    height: 44,
    borderRadius: borderRadius.full,
    backgroundColor: "rgba(255,255,255,0.1)",
    alignItems: "center",
    justifyContent: "center",
  },
  // Rest Screen
  restScreen: {
    flex: 1,
    justifyContent: "center",
    alignItems: "center",
  },
  restContent: {
    alignItems: "center",
  },
  restLabel: {
    fontSize: typography.lg,
    color: "rgba(255,255,255,0.6)",
    letterSpacing: 3,
    marginBottom: spacing.xl,
  },
  restTimerContainer: {
    alignItems: "center",
    marginBottom: spacing.xl,
  },
  restTimerText: {
    fontSize: 80,
    fontWeight: typography.bold,
    color: colors.white,
    marginBottom: spacing.lg,
  },
  progressBar: {
    width: 200,
    height: 8,
    backgroundColor: "rgba(255,255,255,0.2)",
    borderRadius: borderRadius.full,
    overflow: "hidden",
  },
  progressFill: {
    height: "100%",
    backgroundColor: colors.primary,
    borderRadius: borderRadius.full,
  },
  nextSetText: {
    fontSize: typography.md,
    color: "rgba(255,255,255,0.7)",
    marginBottom: spacing.xl,
  },
  skipButton: {
    flexDirection: "row",
    alignItems: "center",
    gap: spacing.sm,
    backgroundColor: "rgba(255,255,255,0.15)",
    paddingVertical: spacing.md,
    paddingHorizontal: spacing.lg,
    borderRadius: borderRadius.full,
  },
  skipButtonText: {
    fontSize: typography.md,
    color: colors.white,
    fontWeight: typography.medium,
  },
  // Completed Screen
  completedScreen: {
    flex: 1,
    justifyContent: "center",
    alignItems: "center",
  },
  completedContent: {
    alignItems: "center",
  },
  completedTitle: {
    fontSize: typography.xxxl,
    fontWeight: typography.bold,
    color: colors.white,
    marginTop: spacing.lg,
    marginBottom: spacing.sm,
  },
  completedTime: {
    fontSize: typography.lg,
    color: "rgba(255,255,255,0.8)",
    marginBottom: spacing.xl,
  },
  completedStats: {
    flexDirection: "row",
    gap: spacing.xxl,
    marginBottom: spacing.xxl,
  },
  completedStat: {
    alignItems: "center",
  },
  completedStatNumber: {
    fontSize: typography.xxxl,
    fontWeight: typography.bold,
    color: colors.white,
  },
  completedStatLabel: {
    fontSize: typography.md,
    color: "rgba(255,255,255,0.7)",
  },
  finishButton: {
    backgroundColor: colors.white,
    paddingVertical: spacing.md,
    paddingHorizontal: spacing.xxl,
    borderRadius: borderRadius.full,
  },
  finishButtonText: {
    fontSize: typography.lg,
    fontWeight: typography.bold,
    color: colors.success,
  },
  errorText: {
    fontSize: typography.lg,
    color: colors.error,
    textAlign: "center",
    padding: spacing.xl,
  },
});
