import React from "react";
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
import { useNavigation, useRoute, RouteProp } from "@react-navigation/native";
import { NativeStackNavigationProp } from "@react-navigation/native-stack";

import { useApp } from "../context/AppContext";
import { Button, ExerciseCard, Card } from "../components";
import { colors, spacing, borderRadius, typography } from "../theme";
import { RootStackParamList } from "../navigation/types";
import { WORKOUT_GOALS, DAYS_OF_WEEK } from "../types";

type WorkoutDetailScreenNavigationProp =
  NativeStackNavigationProp<RootStackParamList>;
type WorkoutDetailScreenRouteProp = RouteProp<
  RootStackParamList,
  "WorkoutDetail"
>;

export const WorkoutDetailScreen: React.FC = () => {
  const navigation = useNavigation<WorkoutDetailScreenNavigationProp>();
  const route = useRoute<WorkoutDetailScreenRouteProp>();
  const { workouts, deleteWorkout } = useApp();

  const workout = workouts.find((w) => w.id === route.params.workoutId);

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
      ],
    );
  };

  const handleStart = () => {
    navigation.navigate("WorkoutSession", { workoutId: workout.id });
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
              <Text style={styles.heroStatValue}>{totalExercises}</Text>
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

        {/* Schedule */}
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

        {/* Exercises */}
        <View style={styles.section}>
          <Text style={styles.sectionTitle}>Exercícios</Text>
          {workout.exercises.map((exercise, index) => (
            <ExerciseCard key={exercise.id} exercise={exercise} index={index} />
          ))}
        </View>

        <View style={{ height: 100 }} />
      </ScrollView>

      {/* Start Button */}
      <View style={styles.footer}>
        <Button
          title="Iniciar Treino"
          onPress={handleStart}
          fullWidth
          icon={<Ionicons name="play" size={20} color={colors.white} />}
        />
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
});
