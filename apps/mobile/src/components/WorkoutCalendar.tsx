import React, { useState, useMemo } from "react";
import { View, Text, StyleSheet, TouchableOpacity } from "react-native";
import { Ionicons } from "@expo/vector-icons";
import { LinearGradient } from "expo-linear-gradient";
import { colors, spacing, borderRadius, typography } from "../theme";

interface WorkoutCalendarProps {
  completedDates: string[]; // Array de datas no formato ISO (YYYY-MM-DD ou ISO string)
}

const WEEKDAYS = ["Dom", "Seg", "Ter", "Qua", "Qui", "Sex", "Sáb"];
const MONTHS = [
  "Janeiro",
  "Fevereiro",
  "Março",
  "Abril",
  "Maio",
  "Junho",
  "Julho",
  "Agosto",
  "Setembro",
  "Outubro",
  "Novembro",
  "Dezembro",
];

export const WorkoutCalendar: React.FC<WorkoutCalendarProps> = ({
  completedDates,
}) => {
  const [currentDate, setCurrentDate] = useState(new Date());

  // Converter datas completadas para Set de strings YYYY-MM-DD para busca rápida
  const completedDatesSet = useMemo(() => {
    const set = new Set<string>();
    completedDates.forEach((dateStr) => {
      const date = new Date(dateStr);
      const key = `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, "0")}-${String(date.getDate()).padStart(2, "0")}`;
      set.add(key);
    });
    return set;
  }, [completedDates]);

  const year = currentDate.getFullYear();
  const month = currentDate.getMonth();

  // Primeiro dia do mês (0 = domingo, 1 = segunda, etc.)
  const firstDayOfMonth = new Date(year, month, 1).getDay();

  // Número de dias no mês
  const daysInMonth = new Date(year, month + 1, 0).getDate();

  // Dias do mês anterior para preencher o início
  const daysInPrevMonth = new Date(year, month, 0).getDate();

  const goToPrevMonth = () => {
    setCurrentDate(new Date(year, month - 1, 1));
  };

  const goToNextMonth = () => {
    setCurrentDate(new Date(year, month + 1, 1));
  };

  const goToToday = () => {
    setCurrentDate(new Date());
  };

  const isToday = (day: number) => {
    const today = new Date();
    return (
      day === today.getDate() &&
      month === today.getMonth() &&
      year === today.getFullYear()
    );
  };

  const isCompleted = (day: number) => {
    const key = `${year}-${String(month + 1).padStart(2, "0")}-${String(day).padStart(2, "0")}`;
    return completedDatesSet.has(key);
  };

  // Contar treinos do mês atual
  const workoutsThisMonth = useMemo(() => {
    let count = 0;
    for (let day = 1; day <= daysInMonth; day++) {
      if (isCompleted(day)) count++;
    }
    return count;
  }, [completedDatesSet, year, month, daysInMonth]);

  // Gerar array de dias para renderizar
  const calendarDays = useMemo(() => {
    const days: Array<{
      day: number;
      isCurrentMonth: boolean;
      isToday: boolean;
      isCompleted: boolean;
    }> = [];

    // Dias do mês anterior
    for (let i = firstDayOfMonth - 1; i >= 0; i--) {
      days.push({
        day: daysInPrevMonth - i,
        isCurrentMonth: false,
        isToday: false,
        isCompleted: false,
      });
    }

    // Dias do mês atual
    for (let day = 1; day <= daysInMonth; day++) {
      days.push({
        day,
        isCurrentMonth: true,
        isToday: isToday(day),
        isCompleted: isCompleted(day),
      });
    }

    // Dias do próximo mês para completar a última semana
    const remainingDays = 7 - (days.length % 7);
    if (remainingDays < 7) {
      for (let day = 1; day <= remainingDays; day++) {
        days.push({
          day,
          isCurrentMonth: false,
          isToday: false,
          isCompleted: false,
        });
      }
    }

    return days;
  }, [year, month, firstDayOfMonth, daysInMonth, daysInPrevMonth]);

  return (
    <View style={styles.container}>
      {/* Header */}
      <View style={styles.header}>
        <TouchableOpacity onPress={goToPrevMonth} style={styles.navButton}>
          <Ionicons name="chevron-back" size={20} color={colors.textPrimary} />
        </TouchableOpacity>

        <TouchableOpacity onPress={goToToday} style={styles.monthContainer}>
          <Text style={styles.monthText}>
            {MONTHS[month]} {year}
          </Text>
          <View style={styles.monthBadge}>
            <Ionicons name="fitness" size={12} color={colors.white} />
            <Text style={styles.monthBadgeText}>{workoutsThisMonth}</Text>
          </View>
        </TouchableOpacity>

        <TouchableOpacity onPress={goToNextMonth} style={styles.navButton}>
          <Ionicons
            name="chevron-forward"
            size={20}
            color={colors.textPrimary}
          />
        </TouchableOpacity>
      </View>

      {/* Weekdays */}
      <View style={styles.weekdaysRow}>
        {WEEKDAYS.map((day) => (
          <View key={day} style={styles.weekdayCell}>
            <Text style={styles.weekdayText}>{day}</Text>
          </View>
        ))}
      </View>

      {/* Calendar Grid */}
      <View style={styles.calendarGrid}>
        {calendarDays.map((dayInfo, index) => (
          <View key={index} style={styles.dayCell}>
            {dayInfo.isCompleted ? (
              <LinearGradient
                colors={[colors.success, "#2ecc71"]}
                style={styles.completedDay}
              >
                <Text style={styles.completedDayText}>{dayInfo.day}</Text>
                <View style={styles.checkIcon}>
                  <Ionicons name="checkmark" size={10} color={colors.white} />
                </View>
              </LinearGradient>
            ) : dayInfo.isToday ? (
              <View style={styles.todayDay}>
                <Text style={styles.todayDayText}>{dayInfo.day}</Text>
              </View>
            ) : (
              <View
                style={[
                  styles.normalDay,
                  !dayInfo.isCurrentMonth && styles.otherMonthDay,
                ]}
              >
                <Text
                  style={[
                    styles.dayText,
                    !dayInfo.isCurrentMonth && styles.otherMonthDayText,
                  ]}
                >
                  {dayInfo.day}
                </Text>
              </View>
            )}
          </View>
        ))}
      </View>

      {/* Legend */}
      <View style={styles.legend}>
        <View style={styles.legendItem}>
          <View style={styles.legendDotCompleted} />
          <Text style={styles.legendText}>Treino concluído</Text>
        </View>
        <View style={styles.legendItem}>
          <View style={styles.legendDotToday} />
          <Text style={styles.legendText}>Hoje</Text>
        </View>
      </View>
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    backgroundColor: colors.surface,
    borderRadius: borderRadius.xl,
    padding: spacing.lg,
  },
  header: {
    flexDirection: "row",
    justifyContent: "space-between",
    alignItems: "center",
    marginBottom: spacing.md,
  },
  navButton: {
    width: 36,
    height: 36,
    borderRadius: borderRadius.full,
    backgroundColor: colors.surfaceSecondary,
    alignItems: "center",
    justifyContent: "center",
  },
  monthContainer: {
    flexDirection: "row",
    alignItems: "center",
    gap: spacing.sm,
  },
  monthText: {
    fontSize: typography.lg,
    fontWeight: typography.bold,
    color: colors.textPrimary,
  },
  monthBadge: {
    flexDirection: "row",
    alignItems: "center",
    gap: 4,
    backgroundColor: colors.success,
    paddingHorizontal: spacing.sm,
    paddingVertical: 2,
    borderRadius: borderRadius.full,
  },
  monthBadgeText: {
    fontSize: typography.xs,
    fontWeight: typography.bold,
    color: colors.white,
  },
  weekdaysRow: {
    flexDirection: "row",
    marginBottom: spacing.sm,
  },
  weekdayCell: {
    flex: 1,
    alignItems: "center",
    paddingVertical: spacing.xs,
  },
  weekdayText: {
    fontSize: typography.xs,
    fontWeight: typography.medium,
    color: colors.textLight,
  },
  calendarGrid: {
    flexDirection: "row",
    flexWrap: "wrap",
  },
  dayCell: {
    width: "14.28%", // 100% / 7 dias
    aspectRatio: 1,
    padding: 2,
    alignItems: "center",
    justifyContent: "center",
  },
  normalDay: {
    width: "100%",
    height: "100%",
    alignItems: "center",
    justifyContent: "center",
    borderRadius: borderRadius.md,
  },
  otherMonthDay: {
    opacity: 0.3,
  },
  dayText: {
    fontSize: typography.sm,
    color: colors.textPrimary,
  },
  otherMonthDayText: {
    color: colors.textLight,
  },
  todayDay: {
    width: "100%",
    height: "100%",
    alignItems: "center",
    justifyContent: "center",
    borderRadius: borderRadius.md,
    borderWidth: 2,
    borderColor: colors.primary,
    backgroundColor: `${colors.primary}10`,
  },
  todayDayText: {
    fontSize: typography.sm,
    fontWeight: typography.bold,
    color: colors.primary,
  },
  completedDay: {
    width: "100%",
    height: "100%",
    alignItems: "center",
    justifyContent: "center",
    borderRadius: borderRadius.md,
    position: "relative",
  },
  completedDayText: {
    fontSize: typography.sm,
    fontWeight: typography.bold,
    color: colors.white,
  },
  checkIcon: {
    position: "absolute",
    top: 2,
    right: 2,
    width: 14,
    height: 14,
    borderRadius: 7,
    backgroundColor: "rgba(255,255,255,0.3)",
    alignItems: "center",
    justifyContent: "center",
  },
  legend: {
    flexDirection: "row",
    justifyContent: "center",
    gap: spacing.lg,
    marginTop: spacing.md,
    paddingTop: spacing.md,
    borderTopWidth: 1,
    borderTopColor: colors.border,
  },
  legendItem: {
    flexDirection: "row",
    alignItems: "center",
    gap: spacing.xs,
  },
  legendDotCompleted: {
    width: 12,
    height: 12,
    borderRadius: 6,
    backgroundColor: colors.success,
  },
  legendDotToday: {
    width: 12,
    height: 12,
    borderRadius: 6,
    borderWidth: 2,
    borderColor: colors.primary,
    backgroundColor: `${colors.primary}20`,
  },
  legendText: {
    fontSize: typography.xs,
    color: colors.textSecondary,
  },
});
