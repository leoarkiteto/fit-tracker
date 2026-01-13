import React from "react";
import { View, Text, StyleSheet, TouchableOpacity } from "react-native";
import { Ionicons } from "@expo/vector-icons";
import { Exercise, MUSCLE_GROUPS } from "../types";
import { colors, spacing, borderRadius, typography, shadows } from "../theme";

interface ExerciseCardProps {
  exercise: Exercise;
  onPress?: () => void;
  onDelete?: () => void;
  index?: number;
}

export const ExerciseCard: React.FC<ExerciseCardProps> = ({
  exercise,
  onPress,
  onDelete,
  index,
}) => {
  const muscleGroup = MUSCLE_GROUPS.find((m) => m.key === exercise.muscleGroup);
  const groupColor = muscleGroup?.color || colors.primary;

  const formatWeight = () => {
    if (exercise.weight === null || exercise.weight === 0) {
      return "Peso livre";
    }
    return `${exercise.weight} kg`;
  };

  const formatRest = () => {
    if (exercise.restSeconds >= 60) {
      const minutes = Math.floor(exercise.restSeconds / 60);
      const seconds = exercise.restSeconds % 60;
      return seconds > 0 ? `${minutes}min ${seconds}s` : `${minutes}min`;
    }
    return `${exercise.restSeconds}s`;
  };

  return (
    <TouchableOpacity
      onPress={onPress}
      activeOpacity={onPress ? 0.7 : 1}
      style={styles.container}
    >
      <View style={[styles.indicator, { backgroundColor: groupColor }]} />
      <View style={styles.content}>
        <View style={styles.header}>
          <View style={styles.titleContainer}>
            {index !== undefined && (
              <View
                style={[styles.indexBadge, { backgroundColor: groupColor }]}
              >
                <Text style={styles.indexText}>{index + 1}</Text>
              </View>
            )}
            <View>
              <Text style={styles.name}>{exercise.name}</Text>
              <Text style={styles.muscleGroup}>{muscleGroup?.label}</Text>
            </View>
          </View>
          {onDelete && (
            <TouchableOpacity onPress={onDelete} style={styles.deleteButton}>
              <Ionicons name="trash-outline" size={18} color={colors.error} />
            </TouchableOpacity>
          )}
        </View>

        <View style={styles.details}>
          <View style={styles.detailItem}>
            <Ionicons
              name="repeat-outline"
              size={16}
              color={colors.textSecondary}
            />
            <Text style={styles.detailText}>
              {exercise.sets}x{exercise.reps}
            </Text>
          </View>
          <View style={styles.detailItem}>
            <Ionicons
              name="barbell-outline"
              size={16}
              color={colors.textSecondary}
            />
            <Text style={styles.detailText}>{formatWeight()}</Text>
          </View>
          <View style={styles.detailItem}>
            <Ionicons
              name="timer-outline"
              size={16}
              color={colors.textSecondary}
            />
            <Text style={styles.detailText}>{formatRest()}</Text>
          </View>
        </View>

        {exercise.notes && (
          <Text style={styles.notes} numberOfLines={2}>
            {exercise.notes}
          </Text>
        )}
      </View>
    </TouchableOpacity>
  );
};

const styles = StyleSheet.create({
  container: {
    flexDirection: "row",
    backgroundColor: colors.surface,
    borderRadius: borderRadius.lg,
    overflow: "hidden",
    marginBottom: spacing.sm,
    ...shadows.sm,
  },
  indicator: {
    width: 4,
  },
  content: {
    flex: 1,
    padding: spacing.md,
  },
  header: {
    flexDirection: "row",
    justifyContent: "space-between",
    alignItems: "flex-start",
    marginBottom: spacing.sm,
  },
  titleContainer: {
    flexDirection: "row",
    alignItems: "center",
    gap: spacing.sm,
    flex: 1,
  },
  indexBadge: {
    width: 28,
    height: 28,
    borderRadius: borderRadius.full,
    alignItems: "center",
    justifyContent: "center",
  },
  indexText: {
    color: colors.white,
    fontSize: typography.sm,
    fontWeight: typography.bold,
  },
  name: {
    fontSize: typography.md,
    fontWeight: typography.semibold,
    color: colors.textPrimary,
  },
  muscleGroup: {
    fontSize: typography.xs,
    color: colors.textSecondary,
    marginTop: 2,
  },
  deleteButton: {
    padding: spacing.xs,
  },
  details: {
    flexDirection: "row",
    gap: spacing.lg,
  },
  detailItem: {
    flexDirection: "row",
    alignItems: "center",
    gap: spacing.xs,
  },
  detailText: {
    fontSize: typography.sm,
    color: colors.textSecondary,
  },
  notes: {
    fontSize: typography.xs,
    color: colors.textLight,
    fontStyle: "italic",
    marginTop: spacing.sm,
  },
});
