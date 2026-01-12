import React, { useState, useCallback } from "react";
import {
  View,
  Text,
  StyleSheet,
  ScrollView,
  TouchableOpacity,
  RefreshControl,
} from "react-native";
import { SafeAreaView } from "react-native-safe-area-context";
import { Ionicons } from "@expo/vector-icons";
import { useNavigation } from "@react-navigation/native";
import { NativeStackNavigationProp } from "@react-navigation/native-stack";

import { useApp } from "../context/AppContext";
import { WorkoutCard, Card } from "../components";
import { colors, spacing, borderRadius, typography } from "../theme";
import { RootStackParamList } from "../navigation/types";
import { WorkoutGoal, WORKOUT_GOALS } from "../types";

type WorkoutsScreenNavigationProp =
  NativeStackNavigationProp<RootStackParamList>;

export const WorkoutsScreen: React.FC = () => {
  const navigation = useNavigation<WorkoutsScreenNavigationProp>();
  const { workouts, completeWorkout, refreshData } = useApp();
  const [filterGoal, setFilterGoal] = useState<WorkoutGoal | null>(null);
  const [refreshing, setRefreshing] = useState(false);

  const onRefresh = useCallback(async () => {
    setRefreshing(true);
    await refreshData();
    setRefreshing(false);
  }, [refreshData]);

  const filteredWorkouts = filterGoal
    ? workouts.filter((w) => w.goal === filterGoal)
    : workouts;

  const handleStartWorkout = async (workoutId: string) => {
    await completeWorkout(workoutId);
  };

  return (
    <SafeAreaView style={styles.container} edges={["top"]}>
      {/* Header */}
      <View style={styles.header}>
        <Text style={styles.title}>Treinos</Text>
        <TouchableOpacity
          onPress={() => navigation.navigate("WorkoutForm", {})}
          style={styles.addButton}
        >
          <Ionicons name="add" size={24} color={colors.white} />
        </TouchableOpacity>
      </View>

      {/* Filters */}
      <ScrollView
        horizontal
        showsHorizontalScrollIndicator={false}
        style={styles.filtersContainer}
        contentContainerStyle={styles.filtersContent}
      >
        <TouchableOpacity
          onPress={() => setFilterGoal(null)}
          style={[styles.filterChip, !filterGoal && styles.filterChipActive]}
        >
          <Text
            style={[
              styles.filterChipText,
              !filterGoal && styles.filterChipTextActive,
            ]}
          >
            Todos
          </Text>
        </TouchableOpacity>
        {WORKOUT_GOALS.map((goal) => (
          <TouchableOpacity
            key={goal.key}
            onPress={() => setFilterGoal(goal.key)}
            style={[
              styles.filterChip,
              filterGoal === goal.key && styles.filterChipActive,
            ]}
          >
            <Text
              style={[
                styles.filterChipText,
                filterGoal === goal.key && styles.filterChipTextActive,
              ]}
            >
              {goal.label}
            </Text>
          </TouchableOpacity>
        ))}
      </ScrollView>

      {/* Workouts List */}
      <ScrollView
        style={styles.listContainer}
        contentContainerStyle={styles.listContent}
        showsVerticalScrollIndicator={false}
        refreshControl={
          <RefreshControl
            refreshing={refreshing}
            onRefresh={onRefresh}
            colors={[colors.primary]}
            tintColor={colors.primary}
          />
        }
      >
        {filteredWorkouts.length > 0 ? (
          filteredWorkouts.map((workout) => (
            <WorkoutCard
              key={workout.id}
              workout={workout}
              onPress={() =>
                navigation.navigate("WorkoutDetail", { workoutId: workout.id })
              }
              onStart={() => handleStartWorkout(workout.id)}
            />
          ))
        ) : (
          <Card style={styles.emptyCard}>
            <Ionicons
              name="barbell-outline"
              size={64}
              color={colors.textLight}
            />
            <Text style={styles.emptyTitle}>Nenhum treino encontrado</Text>
            <Text style={styles.emptyText}>
              {filterGoal
                ? "Não há treinos com esse objetivo"
                : "Comece criando seu primeiro treino"}
            </Text>
            {!filterGoal && (
              <TouchableOpacity
                onPress={() => navigation.navigate("WorkoutForm", {})}
                style={styles.createButton}
              >
                <Ionicons name="add-circle" size={20} color={colors.white} />
                <Text style={styles.createButtonText}>Criar Treino</Text>
              </TouchableOpacity>
            )}
          </Card>
        )}

        <View style={{ height: spacing.xxl }} />
      </ScrollView>
    </SafeAreaView>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: colors.background,
  },
  header: {
    flexDirection: "row",
    justifyContent: "space-between",
    alignItems: "center",
    paddingHorizontal: spacing.lg,
    paddingTop: spacing.md,
    paddingBottom: spacing.md,
  },
  title: {
    fontSize: typography.xxxl,
    fontWeight: typography.bold,
    color: colors.textPrimary,
  },
  addButton: {
    width: 44,
    height: 44,
    borderRadius: borderRadius.full,
    backgroundColor: colors.primary,
    alignItems: "center",
    justifyContent: "center",
  },
  filtersContainer: {
    maxHeight: 50,
  },
  filtersContent: {
    paddingHorizontal: spacing.lg,
    gap: spacing.sm,
  },
  filterChip: {
    paddingHorizontal: spacing.md,
    paddingVertical: spacing.sm,
    borderRadius: borderRadius.full,
    backgroundColor: colors.surfaceSecondary,
    marginRight: spacing.sm,
  },
  filterChipActive: {
    backgroundColor: colors.primary,
  },
  filterChipText: {
    fontSize: typography.sm,
    fontWeight: typography.medium,
    color: colors.textSecondary,
  },
  filterChipTextActive: {
    color: colors.white,
  },
  listContainer: {
    flex: 1,
  },
  listContent: {
    padding: spacing.lg,
  },
  emptyCard: {
    alignItems: "center",
    padding: spacing.xxl,
    gap: spacing.md,
  },
  emptyTitle: {
    fontSize: typography.lg,
    fontWeight: typography.semibold,
    color: colors.textPrimary,
  },
  emptyText: {
    fontSize: typography.md,
    color: colors.textSecondary,
    textAlign: "center",
  },
  createButton: {
    flexDirection: "row",
    alignItems: "center",
    gap: spacing.sm,
    backgroundColor: colors.primary,
    paddingVertical: spacing.md,
    paddingHorizontal: spacing.lg,
    borderRadius: borderRadius.full,
    marginTop: spacing.md,
  },
  createButtonText: {
    color: colors.white,
    fontWeight: typography.semibold,
    fontSize: typography.md,
  },
});
