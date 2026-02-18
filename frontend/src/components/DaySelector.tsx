import React from "react";
import { View, Text, TouchableOpacity, StyleSheet } from "react-native";
import { DayOfWeek, DAYS_OF_WEEK } from "../types";
import { colors, spacing, borderRadius, typography } from "../theme";

interface DaySelectorProps {
  selectedDays: DayOfWeek[];
  onToggleDay: (day: DayOfWeek) => void;
  label?: string;
}

export const DaySelector: React.FC<DaySelectorProps> = ({
  selectedDays,
  onToggleDay,
  label,
}) => {
  return (
    <View style={styles.container}>
      {label && <Text style={styles.label}>{label}</Text>}
      <View style={styles.daysContainer}>
        {DAYS_OF_WEEK.map((day) => {
          const isSelected = selectedDays.includes(day.key);
          return (
            <TouchableOpacity
              key={day.key}
              onPress={() => onToggleDay(day.key)}
              style={[styles.dayButton, isSelected && styles.dayButtonSelected]}
              activeOpacity={0.7}
            >
              <Text
                style={[styles.dayText, isSelected && styles.dayTextSelected]}
              >
                {day.short}
              </Text>
            </TouchableOpacity>
          );
        })}
      </View>
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
  daysContainer: {
    flexDirection: "row",
    justifyContent: "space-between",
    gap: spacing.xs,
  },
  dayButton: {
    flex: 1,
    aspectRatio: 1,
    maxWidth: 48,
    alignItems: "center",
    justifyContent: "center",
    backgroundColor: colors.surfaceSecondary,
    borderRadius: borderRadius.md,
    borderWidth: 2,
    borderColor: "transparent",
  },
  dayButtonSelected: {
    backgroundColor: colors.primary,
    borderColor: colors.primaryDark,
  },
  dayText: {
    fontSize: typography.xs,
    fontWeight: typography.medium,
    color: colors.textSecondary,
  },
  dayTextSelected: {
    color: colors.white,
    fontWeight: typography.bold,
  },
});
