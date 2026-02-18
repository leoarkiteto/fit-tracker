import React, { useState, useEffect, useRef } from "react";
import {
  View,
  Text,
  StyleSheet,
  ScrollView,
  TouchableOpacity,
  Alert,
} from "react-native";
import { SafeAreaView } from "react-native-safe-area-context";
import { LinearGradient } from "expo-linear-gradient";
import { Ionicons } from "@expo/vector-icons";
import {
  useNavigation,
  useRoute,
  RouteProp,
  useFocusEffect,
} from "@react-navigation/native";
import { NativeStackNavigationProp } from "@react-navigation/native-stack";

import { useApp } from "../context/AppContext";
import { Button, Card } from "../components";
import { colors, spacing, borderRadius, typography } from "../theme";
import { RootStackParamList } from "../navigation/types";
import { WORKOUT_GOALS, DAYS_OF_WEEK, MUSCLE_GROUPS, Exercise } from "../types";

const formatDuration = (seconds: number): string => {
  const mins = Math.floor(seconds / 60);
  const secs = seconds % 60;
  if (mins >= 60) {
    const hours = Math.floor(mins / 60);
    const remainMins = mins % 60;
    return `${hours}h ${remainMins}min`;
  }
  return `${mins}min ${secs}s`;
};

type WorkoutDetailScreenNavigationProp =
  NativeStackNavigationProp<RootStackParamList>;
type WorkoutDetailScreenRouteProp = RouteProp<
  RootStackParamList,
  "WorkoutDetail"
>;

export const WorkoutDetailScreen: React.FC = () => {
  const navigation = useNavigation<WorkoutDetailScreenNavigationProp>();
  const route = useRoute<WorkoutDetailScreenRouteProp>();
  const { workouts, deleteWorkout, completeWorkout, isWorkoutCompletedToday } =
    useApp();

  const workout = workouts.find((w) => w.id === route.params.workoutId);
  const alreadyCompletedToday = workout
    ? isWorkoutCompletedToday(workout.id)
    : false;

  // Estados do treino
  const [isWorkoutStarted, setIsWorkoutStarted] = useState(false);
  const [workoutTimer, setWorkoutTimer] = useState(0);
  const [completedExercises, setCompletedExercises] = useState<Set<string>>(
    new Set()
  );
  const [lastExerciseId, setLastExerciseId] = useState<string | null>(null);

  const timerRef = useRef<ReturnType<typeof setInterval> | null>(null);

  // Timer do treino
  useEffect(() => {
    if (isWorkoutStarted) {
      timerRef.current = setInterval(() => {
        setWorkoutTimer((prev) => prev + 1);
      }, 1000);
    }

    return () => {
      if (timerRef.current) {
        clearInterval(timerRef.current);
      }
    };
  }, [isWorkoutStarted]);

  // Detectar quando volta da tela de exercício e marcar como concluído
  useFocusEffect(
    React.useCallback(() => {
      if (lastExerciseId && isWorkoutStarted) {
        setCompletedExercises((prev) => new Set([...prev, lastExerciseId]));
        setLastExerciseId(null);

        // Verificar se todos os exercícios foram concluídos
        if (workout) {
          const newCompleted = new Set([...completedExercises, lastExerciseId]);
          if (newCompleted.size === workout.exercises.length) {
            handleWorkoutComplete();
          }
        }
      }
    }, [lastExerciseId, isWorkoutStarted])
  );

  const formatTime = (seconds: number): string => {
    const mins = Math.floor(seconds / 60);
    const secs = seconds % 60;
    return `${mins.toString().padStart(2, "0")}:${secs.toString().padStart(2, "0")}`;
  };

  if (!workout) {
    return (
      <SafeAreaView style={styles.container}>
        <View style={styles.errorContainer}>
          <Text style={styles.errorText}>Treino não encontrado</Text>
          <Button title="Voltar" onPress={() => navigation.goBack()} />
        </View>
      </SafeAreaView>
    );
  }

  const goal = WORKOUT_GOALS.find((g) => g.key === workout.goal);

  const getGradientColors = (): [string, string] => {
    switch (workout.goal) {
      case "hypertrophy":
        return ["#FF6B6B", "#FF8E8E"];
      case "strength":
        return ["#4ECDC4", "#6EE7DF"];
      case "endurance":
        return ["#45B7D1", "#6DD5ED"];
      case "weight_loss":
        return ["#96CEB4", "#B8E4CC"];
      case "maintenance":
        return ["#DDA0DD", "#EBC4EB"];
      default:
        return [colors.primary, colors.primaryLight];
    }
  };

  const getDaysLabel = () => {
    return workout.days
      .map((d) => {
        const day = DAYS_OF_WEEK.find((dw) => dw.key === d);
        return day?.short || d;
      })
      .join(", ");
  };

  const getMuscleGroupInfo = (muscleGroup: string) => {
    return MUSCLE_GROUPS.find((g) => g.key === muscleGroup);
  };

  const totalExercises = workout.exercises.length;
  const estimatedMinutes = totalExercises * 5;
  const totalSets = workout.exercises.reduce((acc, e) => acc + e.sets, 0);

  const handleDelete = () => {
    Alert.alert(
      "Excluir Treino",
      "Tem certeza que deseja excluir este treino? Esta ação não pode ser desfeita.",
      [
        { text: "Cancelar", style: "cancel" },
        {
          text: "Excluir",
          style: "destructive",
          onPress: async () => {
            await deleteWorkout(workout.id);
            navigation.goBack();
          },
        },
      ]
    );
  };

  const handleStartWorkout = () => {
    setIsWorkoutStarted(true);
    setWorkoutTimer(0);
    setCompletedExercises(new Set());
  };

  const handleExercisePress = (exercise: Exercise) => {
    if (!isWorkoutStarted) {
      Alert.alert(
        "Treino não iniciado",
        "Clique em 'Iniciar Treino' para começar."
      );
      return;
    }

    if (completedExercises.has(exercise.id)) {
      return; // Exercício já concluído
    }

    setLastExerciseId(exercise.id);
    navigation.navigate("ExerciseSession", {
      workoutId: workout.id,
      exerciseId: exercise.id,
    });
  };

  const handleWorkoutComplete = async () => {
    if (timerRef.current) {
      clearInterval(timerRef.current);
    }

    const duration = workoutTimer;
    await completeWorkout(workout.id, duration);

    Alert.alert(
      "🎉 Treino Concluído!",
      `Parabéns! Você completou o treino em ${formatDuration(duration)}.`,
      [
        {
          text: "Finalizar",
          onPress: () => {
            setIsWorkoutStarted(false);
            setCompletedExercises(new Set());
            setWorkoutTimer(0);
          },
        },
      ]
    );
  };

  const handleStopWorkout = () => {
    Alert.alert(
      "Parar Treino",
      "Tem certeza que deseja parar o treino? O progresso será perdido.",
      [
        { text: "Continuar", style: "cancel" },
        {
          text: "Parar",
          style: "destructive",
          onPress: () => {
            if (timerRef.current) {
              clearInterval(timerRef.current);
            }
            setIsWorkoutStarted(false);
            setCompletedExercises(new Set());
            setWorkoutTimer(0);
          },
        },
      ]
    );
  };

  return (
    <SafeAreaView style={styles.container} edges={["top"]}>
      <ScrollView showsVerticalScrollIndicator={false}>
        {/* Hero Section */}
        <LinearGradient
          colors={getGradientColors()}
          start={{ x: 0, y: 0 }}
          end={{ x: 1, y: 1 }}
          style={styles.hero}
        >
          <TouchableOpacity
            onPress={() => navigation.goBack()}
            style={styles.backButton}
          >
            <Ionicons name="arrow-back" size={24} color={colors.white} />
          </TouchableOpacity>

          {!isWorkoutStarted && (
            <View style={styles.heroActions}>
              <TouchableOpacity
                onPress={() =>
                  navigation.navigate("WorkoutForm", { workoutId: workout.id })
                }
                style={styles.heroAction}
              >
                <Ionicons name="pencil" size={20} color={colors.white} />
              </TouchableOpacity>
              <TouchableOpacity onPress={handleDelete} style={styles.heroAction}>
                <Ionicons name="trash" size={20} color={colors.white} />
              </TouchableOpacity>
            </View>
          )}

          {/* Timer Badge quando treino está ativo */}
          {isWorkoutStarted && (
            <View style={styles.activeTimerBadge}>
              <Ionicons name="time" size={18} color={colors.white} />
              <Text style={styles.activeTimerText}>
                {formatTime(workoutTimer)}
              </Text>
            </View>
          )}

          <View style={styles.heroContent}>
            <View style={styles.goalBadge}>
              <Text style={styles.goalBadgeText}>{goal?.label}</Text>
            </View>
            <Text style={styles.heroTitle}>{workout.name}</Text>
            {workout.description && (
              <Text style={styles.heroDescription}>{workout.description}</Text>
            )}
          </View>

          <View style={styles.heroStats}>
            <View style={styles.heroStat}>
              <Ionicons
                name="barbell-outline"
                size={20}
                color="rgba(255,255,255,0.8)"
              />
              <Text style={styles.heroStatValue}>
                {isWorkoutStarted
                  ? `${completedExercises.size}/${totalExercises}`
                  : totalExercises}
              </Text>
              <Text style={styles.heroStatLabel}>Exercícios</Text>
            </View>
            <View style={styles.heroStatDivider} />
            <View style={styles.heroStat}>
              <Ionicons
                name="layers-outline"
                size={20}
                color="rgba(255,255,255,0.8)"
              />
              <Text style={styles.heroStatValue}>{totalSets}</Text>
              <Text style={styles.heroStatLabel}>Séries</Text>
            </View>
            <View style={styles.heroStatDivider} />
            <View style={styles.heroStat}>
              <Ionicons
                name="time-outline"
                size={20}
                color="rgba(255,255,255,0.8)"
              />
              <Text style={styles.heroStatValue}>{estimatedMinutes}</Text>
              <Text style={styles.heroStatLabel}>Minutos</Text>
            </View>
          </View>
        </LinearGradient>

        {/* Schedule - oculto quando treino está ativo */}
        {!isWorkoutStarted && (
          <View style={styles.section}>
            <Text style={styles.sectionTitle}>Programação</Text>
            <Card style={styles.scheduleCard}>
              <Ionicons
                name="calendar-outline"
                size={24}
                color={getGradientColors()[0]}
              />
              <View style={styles.scheduleContent}>
                <Text style={styles.scheduleLabel}>Dias da Semana</Text>
                <Text style={styles.scheduleDays}>{getDaysLabel()}</Text>
              </View>
            </Card>
          </View>
        )}

        {/* Exercises */}
        <View style={styles.section}>
          <Text style={styles.sectionTitle}>
            {isWorkoutStarted
              ? "Toque em um exercício para iniciar"
              : "Exercícios"}
          </Text>

          {workout.exercises.map((exercise, index) => {
            const isCompleted = completedExercises.has(exercise.id);
            const muscleGroup = getMuscleGroupInfo(exercise.muscleGroup);

            return (
              <TouchableOpacity
                key={exercise.id}
                style={[
                  styles.exerciseCard,
                  isCompleted && styles.exerciseCardCompleted,
                ]}
                onPress={() => handleExercisePress(exercise)}
                disabled={isCompleted}
                activeOpacity={isWorkoutStarted ? 0.7 : 1}
              >
                <View
                  style={[
                    styles.exerciseNumber,
                    { backgroundColor: muscleGroup?.color || colors.primary },
                    isCompleted && styles.exerciseNumberCompleted,
                  ]}
                >
                  {isCompleted ? (
                    <Ionicons
                      name="checkmark"
                      size={20}
                      color={colors.white}
                    />
                  ) : (
                    <Text style={styles.exerciseNumberText}>{index + 1}</Text>
                  )}
                </View>

                <View style={styles.exerciseInfo}>
                  <Text
                    style={[
                      styles.exerciseName,
                      isCompleted && styles.exerciseNameCompleted,
                    ]}
                  >
                    {exercise.name}
                  </Text>
                  <Text style={styles.exerciseDetails}>
                    {muscleGroup?.label} • {exercise.sets}×{exercise.reps}
                    {exercise.weight ? ` • ${exercise.weight}kg` : ""}
                  </Text>
                </View>

                {isWorkoutStarted && !isCompleted && (
                  <Ionicons
                    name="chevron-forward"
                    size={20}
                    color={colors.textLight}
                  />
                )}

                {isCompleted && (
                  <View style={styles.completedBadge}>
                    <Text style={styles.completedBadgeText}>Concluído</Text>
                  </View>
                )}
              </TouchableOpacity>
            );
          })}
        </View>

        <View style={{ height: 100 }} />
      </ScrollView>

      {/* Footer Button */}
      <View style={styles.footer}>
        {alreadyCompletedToday && !isWorkoutStarted ? (
          <View style={styles.completedTodayBanner}>
            <Ionicons name="checkmark-circle" size={24} color={colors.success} />
            <Text style={styles.completedTodayText}>
              Treino já concluído hoje!
            </Text>
          </View>
        ) : isWorkoutStarted ? (
          <View style={styles.activeFooter}>
            <TouchableOpacity
              style={styles.stopButton}
              onPress={handleStopWorkout}
            >
              <Ionicons name="stop" size={20} color={colors.error} />
              <Text style={styles.stopButtonText}>Parar</Text>
            </TouchableOpacity>

            {completedExercises.size === totalExercises && (
              <TouchableOpacity
                style={styles.finishButton}
                onPress={handleWorkoutComplete}
              >
                <Ionicons name="trophy" size={20} color={colors.white} />
                <Text style={styles.finishButtonText}>Concluir Treino</Text>
              </TouchableOpacity>
            )}
          </View>
        ) : (
          <Button
            title="Iniciar Treino"
            onPress={handleStartWorkout}
            fullWidth
            icon={<Ionicons name="play" size={20} color={colors.white} />}
          />
        )}
      </View>
    </SafeAreaView>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: colors.background,
  },
  errorContainer: {
    flex: 1,
    alignItems: "center",
    justifyContent: "center",
    gap: spacing.lg,
  },
  errorText: {
    fontSize: typography.lg,
    color: colors.textSecondary,
  },
  hero: {
    paddingTop: spacing.md,
    paddingBottom: spacing.xl,
    paddingHorizontal: spacing.lg,
    borderBottomLeftRadius: borderRadius.xl,
    borderBottomRightRadius: borderRadius.xl,
  },
  backButton: {
    width: 40,
    height: 40,
    alignItems: "center",
    justifyContent: "center",
    backgroundColor: "rgba(255,255,255,0.2)",
    borderRadius: borderRadius.full,
  },
  heroActions: {
    position: "absolute",
    top: spacing.md,
    right: spacing.lg,
    flexDirection: "row",
    gap: spacing.sm,
  },
  heroAction: {
    width: 40,
    height: 40,
    alignItems: "center",
    justifyContent: "center",
    backgroundColor: "rgba(255,255,255,0.2)",
    borderRadius: borderRadius.full,
  },
  activeTimerBadge: {
    position: "absolute",
    top: spacing.md,
    right: spacing.lg,
    flexDirection: "row",
    alignItems: "center",
    gap: spacing.xs,
    backgroundColor: "rgba(0,0,0,0.3)",
    paddingHorizontal: spacing.md,
    paddingVertical: spacing.sm,
    borderRadius: borderRadius.full,
  },
  activeTimerText: {
    fontSize: typography.md,
    fontWeight: typography.bold,
    color: colors.white,
  },
  heroContent: {
    marginTop: spacing.lg,
    marginBottom: spacing.xl,
  },
  goalBadge: {
    alignSelf: "flex-start",
    backgroundColor: "rgba(255,255,255,0.2)",
    paddingHorizontal: spacing.md,
    paddingVertical: spacing.xs,
    borderRadius: borderRadius.full,
    marginBottom: spacing.sm,
  },
  goalBadgeText: {
    fontSize: typography.xs,
    fontWeight: typography.bold,
    color: colors.white,
  },
  heroTitle: {
    fontSize: typography.xxxl,
    fontWeight: typography.bold,
    color: colors.white,
    marginBottom: spacing.xs,
  },
  heroDescription: {
    fontSize: typography.md,
    color: "rgba(255,255,255,0.8)",
  },
  heroStats: {
    flexDirection: "row",
    justifyContent: "space-around",
    backgroundColor: "rgba(255,255,255,0.15)",
    borderRadius: borderRadius.lg,
    paddingVertical: spacing.md,
  },
  heroStat: {
    alignItems: "center",
    gap: spacing.xs,
  },
  heroStatValue: {
    fontSize: typography.xl,
    fontWeight: typography.bold,
    color: colors.white,
  },
  heroStatLabel: {
    fontSize: typography.xs,
    color: "rgba(255,255,255,0.7)",
  },
  heroStatDivider: {
    width: 1,
    backgroundColor: "rgba(255,255,255,0.2)",
  },
  section: {
    paddingHorizontal: spacing.lg,
    marginTop: spacing.lg,
  },
  sectionTitle: {
    fontSize: typography.lg,
    fontWeight: typography.bold,
    color: colors.textPrimary,
    marginBottom: spacing.md,
  },
  scheduleCard: {
    flexDirection: "row",
    alignItems: "center",
    padding: spacing.lg,
    gap: spacing.md,
  },
  scheduleContent: {
    flex: 1,
  },
  scheduleLabel: {
    fontSize: typography.xs,
    color: colors.textSecondary,
  },
  scheduleDays: {
    fontSize: typography.md,
    fontWeight: typography.semibold,
    color: colors.textPrimary,
    marginTop: spacing.xs,
  },
  exerciseCard: {
    flexDirection: "row",
    alignItems: "center",
    backgroundColor: colors.surface,
    borderRadius: borderRadius.lg,
    padding: spacing.md,
    marginBottom: spacing.sm,
  },
  exerciseCardCompleted: {
    opacity: 0.7,
    backgroundColor: `${colors.success}10`,
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
  exerciseInfo: {
    flex: 1,
    marginLeft: spacing.md,
  },
  exerciseName: {
    fontSize: typography.md,
    fontWeight: typography.semibold,
    color: colors.textPrimary,
  },
  exerciseNameCompleted: {
    textDecorationLine: "line-through",
    color: colors.textSecondary,
  },
  exerciseDetails: {
    fontSize: typography.sm,
    color: colors.textSecondary,
    marginTop: spacing.xs,
  },
  completedBadge: {
    backgroundColor: colors.success,
    paddingHorizontal: spacing.sm,
    paddingVertical: spacing.xs,
    borderRadius: borderRadius.full,
  },
  completedBadgeText: {
    fontSize: typography.xs,
    fontWeight: typography.semibold,
    color: colors.white,
  },
  footer: {
    position: "absolute",
    bottom: 0,
    left: 0,
    right: 0,
    padding: spacing.lg,
    backgroundColor: colors.surface,
    borderTopWidth: 1,
    borderTopColor: colors.border,
  },
  activeFooter: {
    flexDirection: "row",
    gap: spacing.md,
  },
  stopButton: {
    flexDirection: "row",
    alignItems: "center",
    justifyContent: "center",
    gap: spacing.sm,
    paddingVertical: spacing.md,
    paddingHorizontal: spacing.lg,
    borderRadius: borderRadius.lg,
    borderWidth: 2,
    borderColor: colors.error,
  },
  stopButtonText: {
    fontSize: typography.md,
    fontWeight: typography.semibold,
    color: colors.error,
  },
  finishButton: {
    flex: 1,
    flexDirection: "row",
    alignItems: "center",
    justifyContent: "center",
    gap: spacing.sm,
    paddingVertical: spacing.md,
    paddingHorizontal: spacing.lg,
    borderRadius: borderRadius.lg,
    backgroundColor: colors.success,
  },
  finishButtonText: {
    fontSize: typography.md,
    fontWeight: typography.semibold,
    color: colors.white,
  },
  completedTodayBanner: {
    flexDirection: "row",
    alignItems: "center",
    justifyContent: "center",
    gap: spacing.sm,
    paddingVertical: spacing.md,
    paddingHorizontal: spacing.lg,
    backgroundColor: `${colors.success}15`,
    borderRadius: borderRadius.lg,
    borderWidth: 1,
    borderColor: colors.success,
  },
  completedTodayText: {
    fontSize: typography.md,
    fontWeight: typography.semibold,
    color: colors.success,
  },
});
