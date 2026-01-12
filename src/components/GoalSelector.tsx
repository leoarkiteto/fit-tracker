import React from "react";
import {
  View,
  Text,
  TouchableOpacity,
  StyleSheet,
  ScrollView,
} from "react-native";
import { MaterialIcons } from "@expo/vector-icons";
import { WorkoutGoal, WORKOUT_GOALS } from "../types";
import { colors, spacing, borderRadius, typography, shadows } from "../theme";

interface GoalSelectorProps {
  selectedGoal: WorkoutGoal | null;
  onSelectGoal: (goal: WorkoutGoal) => void;
  label?: string;
}

export const GoalSelector: React.FC<GoalSelectorProps> = ({
  selectedGoal,
  onSelectGoal,
  label,
}) => {
  const getGoalColor = (goal: WorkoutGoal) => {
    switch (goal) {
      case "hypertrophy":
        return "#FF6B6B";
      case "strength":
        return "#4ECDC4";
      case "endurance":
        return "#45B7D1";
      case "weight_loss":
        return "#96CEB4";
      case "maintenance":
        return "#DDA0DD";
      default:
        return colors.primary;
    }
  };

  return (
    <View style={styles.container}>
      {label && <Text style={styles.label}>{label}</Text>}
      <ScrollView
        horizontal
        showsHorizontalScrollIndicator={false}
        contentContainerStyle={styles.scrollContent}
      >
        {WORKOUT_GOALS.map((goal) => {
          const isSelected = selectedGoal === goal.key;
          const goalColor = getGoalColor(goal.key);

          return (
            <TouchableOpacity
              key={goal.key}
              onPress={() => onSelectGoal(goal.key)}
              style={[
                styles.goalCard,
                isSelected && { backgroundColor: goalColor },
                isSelected && shadows.md,
              ]}
              activeOpacity={0.7}
            >
              <View
                style={[
                  styles.iconContainer,
                  {
                    backgroundColor: isSelected
                      ? "rgba(255,255,255,0.2)"
                      : `${goalColor}20`,
                  },
                ]}
              >
                <MaterialIcons
                  name={goal.icon as any}
                  size={24}
                  color={isSelected ? colors.white : goalColor}
                />
              </View>
              <Text
                style={[styles.goalText, isSelected && styles.goalTextSelected]}
                numberOfLines={1}
              >
                {goal.label}
              </Text>
            </TouchableOpacity>
          );
        })}
      </ScrollView>
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    marginBottom: spacing.md,
  },
  label: {
    fontSize: typography.sm,
    fontWeight: typography.medium,
    color: colors.textPrimary,
    marginBottom: spacing.sm,
  },
  scrollContent: {
    paddingRight: spacing.md,
    gap: spacing.sm,
  },
  goalCard: {
    alignItems: "center",
    justifyContent: "center",
    backgroundColor: colors.surfaceSecondary,
    borderRadius: borderRadius.lg,
    padding: spacing.md,
    minWidth: 100,
    gap: spacing.sm,
  },
  iconContainer: {
    width: 48,
    height: 48,
    borderRadius: borderRadius.full,
    alignItems: "center",
    justifyContent: "center",
  },
  goalText: {
    fontSize: typography.sm,
    fontWeight: typography.medium,
    color: colors.textPrimary,
  },
  goalTextSelected: {
    color: colors.white,
    fontWeight: typography.bold,
  },
});
