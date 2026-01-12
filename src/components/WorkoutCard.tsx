import React from "react";
import { View, Text, StyleSheet, TouchableOpacity } from "react-native";
import { Ionicons } from "@expo/vector-icons";
import { LinearGradient } from "expo-linear-gradient";
import { Workout, WORKOUT_GOALS, DAYS_OF_WEEK } from "../types";
import { colors, spacing, borderRadius, typography, shadows } from "../theme";

interface WorkoutCardProps {
  workout: Workout;
  onPress: () => void;
  onStart?: () => void;
  variant?: "default" | "compact" | "featured";
}

export const WorkoutCard: React.FC<WorkoutCardProps> = ({
  workout,
  onPress,
  onStart,
  variant = "default",
}) => {
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
    const dayLabels = workout.days.map((d) => {
      const day = DAYS_OF_WEEK.find((dw) => dw.key === d);
      return day?.short || d;
    });
    return dayLabels.join(", ");
  };

  const totalExercises = workout.exercises.length;
  const estimatedMinutes = totalExercises * 5; // ~5 min per exercise

  if (variant === "featured") {
    return (
      <TouchableOpacity onPress={onPress} activeOpacity={0.9}>
        <LinearGradient
          colors={getGradientColors()}
          start={{ x: 0, y: 0 }}
          end={{ x: 1, y: 1 }}
          style={styles.featuredContainer}
        >
          <View style={styles.featuredContent}>
            <View>
              <Text style={styles.featuredTitle}>{workout.name}</Text>
              <Text style={styles.featuredSubtitle}>
                {totalExercises} exercícios • {estimatedMinutes} min
              </Text>
            </View>
            <View style={styles.featuredStats}>
              <View style={styles.featuredStat}>
                <Ionicons
                  name="calendar-outline"
                  size={16}
                  color="rgba(255,255,255,0.8)"
                />
                <Text style={styles.featuredStatText}>{getDaysLabel()}</Text>
              </View>
            </View>
            {onStart && (
              <TouchableOpacity onPress={onStart} style={styles.startButton}>
                <Text style={styles.startButtonText}>Iniciar</Text>
                <Ionicons
                  name="play"
                  size={18}
                  color={getGradientColors()[0]}
                />
              </TouchableOpacity>
            )}
          </View>
          <View style={styles.featuredImageContainer}>
            <Ionicons name="fitness" size={80} color="rgba(255,255,255,0.2)" />
          </View>
        </LinearGradient>
      </TouchableOpacity>
    );
  }

  if (variant === "compact") {
    return (
      <TouchableOpacity
        onPress={onPress}
        style={styles.compactContainer}
        activeOpacity={0.7}
      >
        <View
          style={[
            styles.compactIndicator,
            { backgroundColor: getGradientColors()[0] },
          ]}
        />
        <View style={styles.compactContent}>
          <Text style={styles.compactTitle} numberOfLines={1}>
            {workout.name}
          </Text>
          <Text style={styles.compactSubtitle}>
            {totalExercises} exercícios
          </Text>
        </View>
        <Ionicons name="chevron-forward" size={20} color={colors.textLight} />
      </TouchableOpacity>
    );
  }

  return (
    <TouchableOpacity
      onPress={onPress}
      style={styles.container}
      activeOpacity={0.7}
    >
      <View style={styles.header}>
        <View
          style={[
            styles.goalBadge,
            { backgroundColor: `${getGradientColors()[0]}20` },
          ]}
        >
          <Text style={[styles.goalText, { color: getGradientColors()[0] }]}>
            {goal?.label}
          </Text>
        </View>
        <Ionicons
          name="ellipsis-horizontal"
          size={20}
          color={colors.textLight}
        />
      </View>

      <Text style={styles.title}>{workout.name}</Text>
      {workout.description && (
        <Text style={styles.description} numberOfLines={2}>
          {workout.description}
        </Text>
      )}

      <View style={styles.meta}>
        <View style={styles.metaItem}>
          <Ionicons
            name="barbell-outline"
            size={18}
            color={colors.textSecondary}
          />
          <Text style={styles.metaText}>{totalExercises} exercícios</Text>
        </View>
        <View style={styles.metaItem}>
          <Ionicons
            name="time-outline"
            size={18}
            color={colors.textSecondary}
          />
          <Text style={styles.metaText}>{estimatedMinutes} min</Text>
        </View>
      </View>

      <View style={styles.footer}>
        <Text style={styles.daysText}>{getDaysLabel()}</Text>
        {onStart && (
          <TouchableOpacity onPress={onStart} style={styles.playButton}>
            <LinearGradient
              colors={getGradientColors()}
              style={styles.playButtonGradient}
            >
              <Ionicons name="play" size={16} color={colors.white} />
            </LinearGradient>
          </TouchableOpacity>
        )}
      </View>

      <View style={styles.progressBar}>
        <View
          style={[
            styles.progressFill,
            { width: "30%", backgroundColor: getGradientColors()[0] },
          ]}
        />
      </View>
    </TouchableOpacity>
  );
};

const styles = StyleSheet.create({
  container: {
    backgroundColor: colors.surface,
    borderRadius: borderRadius.xl,
    padding: spacing.lg,
    marginBottom: spacing.md,
    ...shadows.md,
  },
  header: {
    flexDirection: "row",
    justifyContent: "space-between",
    alignItems: "center",
    marginBottom: spacing.sm,
  },
  goalBadge: {
    paddingHorizontal: spacing.sm,
    paddingVertical: spacing.xs,
    borderRadius: borderRadius.full,
  },
  goalText: {
    fontSize: typography.xs,
    fontWeight: typography.semibold,
  },
  title: {
    fontSize: typography.xl,
    fontWeight: typography.bold,
    color: colors.textPrimary,
    marginBottom: spacing.xs,
  },
  description: {
    fontSize: typography.sm,
    color: colors.textSecondary,
    marginBottom: spacing.md,
  },
  meta: {
    flexDirection: "row",
    gap: spacing.lg,
    marginBottom: spacing.md,
  },
  metaItem: {
    flexDirection: "row",
    alignItems: "center",
    gap: spacing.xs,
  },
  metaText: {
    fontSize: typography.sm,
    color: colors.textSecondary,
  },
  footer: {
    flexDirection: "row",
    justifyContent: "space-between",
    alignItems: "center",
  },
  daysText: {
    fontSize: typography.sm,
    color: colors.textLight,
  },
  playButton: {
    borderRadius: borderRadius.full,
    overflow: "hidden",
  },
  playButtonGradient: {
    width: 40,
    height: 40,
    alignItems: "center",
    justifyContent: "center",
    borderRadius: borderRadius.full,
  },
  progressBar: {
    height: 4,
    backgroundColor: colors.surfaceSecondary,
    borderRadius: borderRadius.full,
    marginTop: spacing.md,
    overflow: "hidden",
  },
  progressFill: {
    height: "100%",
    borderRadius: borderRadius.full,
  },

  // Featured variant
  featuredContainer: {
    borderRadius: borderRadius.xl,
    padding: spacing.lg,
    marginBottom: spacing.md,
    flexDirection: "row",
    overflow: "hidden",
    minHeight: 160,
  },
  featuredContent: {
    flex: 1,
    justifyContent: "space-between",
  },
  featuredTitle: {
    fontSize: typography.xxl,
    fontWeight: typography.bold,
    color: colors.white,
    marginBottom: spacing.xs,
  },
  featuredSubtitle: {
    fontSize: typography.sm,
    color: "rgba(255,255,255,0.8)",
  },
  featuredStats: {
    gap: spacing.sm,
  },
  featuredStat: {
    flexDirection: "row",
    alignItems: "center",
    gap: spacing.xs,
  },
  featuredStatText: {
    fontSize: typography.sm,
    color: "rgba(255,255,255,0.9)",
  },
  startButton: {
    flexDirection: "row",
    alignItems: "center",
    justifyContent: "center",
    backgroundColor: colors.white,
    paddingVertical: spacing.sm,
    paddingHorizontal: spacing.lg,
    borderRadius: borderRadius.full,
    gap: spacing.xs,
    alignSelf: "flex-start",
  },
  startButtonText: {
    fontSize: typography.md,
    fontWeight: typography.semibold,
  },
  featuredImageContainer: {
    position: "absolute",
    right: -20,
    bottom: -20,
    opacity: 0.3,
  },

  // Compact variant
  compactContainer: {
    flexDirection: "row",
    alignItems: "center",
    backgroundColor: colors.surface,
    borderRadius: borderRadius.lg,
    padding: spacing.md,
    marginBottom: spacing.sm,
    ...shadows.sm,
  },
  compactIndicator: {
    width: 4,
    height: 40,
    borderRadius: borderRadius.full,
    marginRight: spacing.md,
  },
  compactContent: {
    flex: 1,
  },
  compactTitle: {
    fontSize: typography.md,
    fontWeight: typography.semibold,
    color: colors.textPrimary,
  },
  compactSubtitle: {
    fontSize: typography.xs,
    color: colors.textSecondary,
    marginTop: 2,
  },
});
