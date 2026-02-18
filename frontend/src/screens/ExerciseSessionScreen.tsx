import React, { useState, useEffect, useRef } from "react";
import {
  View,
  Text,
  StyleSheet,
  TouchableOpacity,
  Vibration,
  StatusBar,
} from "react-native";
import { Ionicons } from "@expo/vector-icons";
import { useNavigation, useRoute, RouteProp } from "@react-navigation/native";
import { NativeStackNavigationProp } from "@react-navigation/native-stack";
import { LinearGradient } from "expo-linear-gradient";

import { useApp } from "../context/AppContext";
import { colors, spacing, borderRadius, typography } from "../theme";
import { RootStackParamList } from "../navigation/types";
import { Exercise, MUSCLE_GROUPS } from "../types";

type ExerciseSessionScreenNavigationProp = NativeStackNavigationProp<
  RootStackParamList,
  "ExerciseSession"
>;
type ExerciseSessionScreenRouteProp = RouteProp<
  RootStackParamList,
  "ExerciseSession"
>;

type SessionPhase = "exercising" | "resting" | "completed";

export const ExerciseSessionScreen: React.FC = () => {
  const navigation = useNavigation<ExerciseSessionScreenNavigationProp>();
  const route = useRoute<ExerciseSessionScreenRouteProp>();
  const { workouts } = useApp();

  const workout = workouts.find((w) => w.id === route.params.workoutId);
  const exercise = workout?.exercises.find(
    (e) => e.id === route.params.exerciseId
  );

  const [sessionPhase, setSessionPhase] = useState<SessionPhase>("exercising");
  const [currentSet, setCurrentSet] = useState(1);
  const [restTimer, setRestTimer] = useState(0);

  const restTimerRef = useRef<ReturnType<typeof setInterval> | null>(null);

  // Timer de descanso
  useEffect(() => {
    if (sessionPhase === "resting" && restTimer > 0) {
      restTimerRef.current = setInterval(() => {
        setRestTimer((prev) => {
          if (prev <= 1) {
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

  const alertRestComplete = () => {
    Vibration.vibrate([0, 500, 200, 500, 200, 500]);
  };

  const formatTime = (seconds: number): string => {
    const mins = Math.floor(seconds / 60);
    const secs = seconds % 60;
    return `${mins.toString().padStart(2, "0")}:${secs.toString().padStart(2, "0")}`;
  };

  const handleCompletedSet = () => {
    if (!exercise) return;

    if (currentSet >= exercise.sets) {
      // Exercício completo
      setSessionPhase("completed");
      // Voltar para a tela de detalhes após 1.5 segundos
      setTimeout(() => {
        navigation.goBack();
      }, 1500);
    } else {
      // Iniciar descanso
      setRestTimer(exercise.restSeconds);
      setSessionPhase("resting");
    }
  };

  const handleRestComplete = () => {
    alertRestComplete();
    setCurrentSet((prev) => prev + 1);
    setSessionPhase("exercising");
  };

  const handleSkipRest = () => {
    if (restTimerRef.current) {
      clearInterval(restTimerRef.current);
    }
    setCurrentSet((prev) => prev + 1);
    setSessionPhase("exercising");
  };

  const handleExit = () => {
    navigation.goBack();
  };

  const getMuscleGroupColor = (muscleGroup: string) => {
    const group = MUSCLE_GROUPS.find((g) => g.key === muscleGroup);
    return group?.color || colors.primary;
  };

  if (!workout || !exercise) {
    return (
      <View style={styles.container}>
        <Text style={styles.errorText}>Exercício não encontrado</Text>
      </View>
    );
  }

  // Tela de exercício completo
  if (sessionPhase === "completed") {
    return (
      <View style={styles.fullScreenContainer}>
        <StatusBar barStyle="light-content" />
        <LinearGradient
          colors={[colors.success, "#1a1a2e"]}
          style={styles.completedScreen}
        >
          <View style={styles.completedContent}>
            <Ionicons name="checkmark-circle" size={100} color={colors.white} />
            <Text style={styles.completedTitle}>Exercício Concluído!</Text>
            <Text style={styles.completedExercise}>{exercise.name}</Text>
            <Text style={styles.completedSets}>
              {exercise.sets} séries × {exercise.reps} reps
            </Text>
          </View>
        </LinearGradient>
      </View>
    );
  }

  // Tela de descanso
  if (sessionPhase === "resting") {
    const progress = (restTimer / exercise.restSeconds) * 100;

    return (
      <View style={styles.fullScreenContainer}>
        <StatusBar barStyle="light-content" />
        <LinearGradient
          colors={["#2d3436", "#1a1a2e"]}
          style={styles.restScreen}
        >
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
              Próxima: Série {currentSet + 1} de {exercise.sets}
            </Text>

            <TouchableOpacity style={styles.skipButton} onPress={handleSkipRest}>
              <Ionicons name="play-forward" size={20} color={colors.white} />
              <Text style={styles.skipButtonText}>Pular Descanso</Text>
            </TouchableOpacity>
          </View>

          <TouchableOpacity style={styles.exitButton} onPress={handleExit}>
            <Ionicons name="close" size={24} color="rgba(255,255,255,0.7)" />
          </TouchableOpacity>
        </LinearGradient>
      </View>
    );
  }

  // Tela de exercício ativo
  return (
    <TouchableOpacity
      style={styles.fullScreenContainer}
      activeOpacity={1}
      onPress={handleCompletedSet}
    >
      <StatusBar barStyle="light-content" />
      <LinearGradient
        colors={[getMuscleGroupColor(exercise.muscleGroup), "#1a1a2e"]}
        style={styles.exerciseScreen}
      >
        <View style={styles.exerciseContent}>
          <Text style={styles.exerciseLabel}>EXERCÍCIO</Text>
          <Text style={styles.exerciseName}>{exercise.name}</Text>

          <View style={styles.setInfo}>
            <Text style={styles.setNumber}>Série {currentSet}</Text>
            <Text style={styles.setTotal}>de {exercise.sets}</Text>
          </View>

          <View style={styles.repInfo}>
            <Text style={styles.repNumber}>{exercise.reps}</Text>
            <Text style={styles.repLabel}>repetições</Text>
          </View>

          {exercise.weight && (
            <View style={styles.weightInfo}>
              <Ionicons name="barbell" size={24} color="rgba(255,255,255,0.8)" />
              <Text style={styles.weightText}>{exercise.weight} kg</Text>
            </View>
          )}
        </View>

        <View style={styles.tapHint}>
          <Ionicons name="hand-left-outline" size={32} color="rgba(255,255,255,0.5)" />
          <Text style={styles.tapHintText}>
            Toque na tela ao completar a série
          </Text>
        </View>

        <TouchableOpacity style={styles.exitButton} onPress={handleExit}>
          <Ionicons name="close" size={24} color="rgba(255,255,255,0.7)" />
        </TouchableOpacity>
      </LinearGradient>
    </TouchableOpacity>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: colors.background,
    alignItems: "center",
    justifyContent: "center",
  },
  fullScreenContainer: {
    flex: 1,
  },
  errorText: {
    fontSize: typography.lg,
    color: colors.error,
  },
  // Exercise Screen
  exerciseScreen: {
    flex: 1,
    justifyContent: "center",
    alignItems: "center",
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
    paddingHorizontal: spacing.lg,
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
    top: 60,
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
    fontSize: typography.xxl,
    fontWeight: typography.bold,
    color: colors.white,
    marginTop: spacing.lg,
  },
  completedExercise: {
    fontSize: typography.lg,
    color: "rgba(255,255,255,0.8)",
    marginTop: spacing.sm,
  },
  completedSets: {
    fontSize: typography.md,
    color: "rgba(255,255,255,0.6)",
    marginTop: spacing.xs,
  },
});
