import React, { useCallback } from "react";
import {
  View,
  Text,
  StyleSheet,
  ScrollView,
  TouchableOpacity,
  RefreshControl,
  ActivityIndicator,
} from "react-native";
import { SafeAreaView } from "react-native-safe-area-context";
import { LinearGradient } from "expo-linear-gradient";
import { Ionicons } from "@expo/vector-icons";
import { useNavigation } from "@react-navigation/native";
import { NativeStackNavigationProp } from "@react-navigation/native-stack";

import { useApp } from "../context/AppContext";
import { WorkoutCard, Card, WorkoutCalendar } from "../components";
import { colors, spacing, borderRadius, typography, shadows } from "../theme";
import { RootStackParamList } from "../navigation/types";
import { DAYS_OF_WEEK, DayOfWeek } from "../types";
import { waterApi } from "../services";
import type { DailyWaterSummary } from "../types";

type HomeScreenNavigationProp = NativeStackNavigationProp<RootStackParamList>;

export const HomeScreen: React.FC = () => {
  const navigation = useNavigation<HomeScreenNavigationProp>();
  const {
    profile,
    profileId,
    workouts,
    totalWorkoutsCompleted,
    workoutsThisWeek,
    totalMinutesSpent,
    bioimpedanceHistory,
    completedWorkouts,
    isLoading,
    error,
    clearError,
    refreshData,
    isWorkoutCompletedToday,
  } = useApp();

  const [refreshing, setRefreshing] = React.useState(false);
  const [waterSummary, setWaterSummary] = React.useState<DailyWaterSummary | null>(null);
  const [waterAdding, setWaterAdding] = React.useState(false);

  const todayDateUtc = () => new Date().toISOString().slice(0, 10);

  const loadWater = useCallback(async () => {
    if (!profileId) return;
    try {
      const summary = await waterApi.getDay(profileId, todayDateUtc());
      setWaterSummary(summary);
    } catch (e) {
      console.error("Error loading water:", e);
    }
  }, [profileId]);

  React.useEffect(() => {
    if (profileId && !isLoading) {
      loadWater();
    }
  }, [profileId, isLoading, loadWater]);

  const addWater = useCallback(
    async (amountMl: number) => {
      if (!profileId || waterAdding) return;
      setWaterAdding(true);
      try {
        await waterApi.add(profileId, amountMl);
        await loadWater();
      } catch (e) {
        console.error("Error adding water:", e);
      } finally {
        setWaterAdding(false);
      }
    },
    [profileId, waterAdding, loadWater]
  );

  const onRefresh = useCallback(async () => {
    setRefreshing(true);
    await refreshData();
    await loadWater();
    setRefreshing(false);
  }, [refreshData, loadWater]);

  // Get today's day of week
  const today = new Date()
    .toLocaleDateString("en-US", { weekday: "long" })
    .toLowerCase() as DayOfWeek;
  const todaysWorkouts = workouts.filter((w) => w.days.includes(today));

  // Get latest bioimpedance (sorted by date descending from API)
  const latestBioimpedance =
    bioimpedanceHistory.length > 0 ? bioimpedanceHistory[0] : null;


  // Loading state
  if (isLoading) {
    return (
      <SafeAreaView style={styles.container} edges={["top"]}>
        <View style={styles.loadingContainer}>
          <ActivityIndicator size="large" color={colors.primary} />
          <Text style={styles.loadingText}>Carregando...</Text>
        </View>
      </SafeAreaView>
    );
  }

  return (
    <SafeAreaView style={styles.container} edges={["top"]}>
      <ScrollView
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
        {/* Error Banner */}
        {error && (
          <TouchableOpacity style={styles.errorBanner} onPress={clearError}>
            <Ionicons
              name="alert-circle"
              size={20}
              color={colors.white}
              style={styles.errorIcon}
            />
            <Text style={styles.errorText}>{error}</Text>
            <Ionicons name="close" size={20} color={colors.white} />
          </TouchableOpacity>
        )}

        {/* Header */}
        <View style={styles.header}>
          <View>
            <Text style={styles.greeting}>Olá, {profile.name}!</Text>
            <Text style={styles.subtitle}>Vamos verificar sua atividade</Text>
          </View>
          <TouchableOpacity style={styles.avatarContainer}>
            <LinearGradient
              colors={[colors.primary, colors.primaryLight]}
              style={styles.avatar}
            >
              <Text style={styles.avatarText}>
                {profile.name.charAt(0).toUpperCase()}
              </Text>
            </LinearGradient>
          </TouchableOpacity>
        </View>

        {/* Stats Cards */}
        <View style={styles.statsContainer}>
          <Card style={styles.statCard}>
            <View style={styles.statIconContainer}>
              <Ionicons
                name="checkmark-circle"
                size={24}
                color={colors.success}
              />
            </View>
            <Text style={styles.statNumber}>{totalWorkoutsCompleted}</Text>
            <Text style={styles.statLabel}>Treinos{"\n"}Concluídos</Text>
          </Card>

          <Card style={styles.statCard}>
            <View style={styles.statIconContainer}>
              <Ionicons name="flame" size={24} color={colors.warning} />
            </View>
            <Text style={styles.statNumber}>{workoutsThisWeek}</Text>
            <Text style={styles.statLabel}>Esta{"\n"}Semana</Text>
          </Card>

          <Card style={styles.statCard}>
            <View style={styles.statIconContainer}>
              <Ionicons name="time" size={24} color={colors.info} />
            </View>
            <Text style={styles.statNumber}>{totalMinutesSpent}</Text>
            <Text style={styles.statLabel}>Minutos{"\n"}Total</Text>
          </Card>
        </View>

        {/* Water today */}
        <View style={styles.section}>
          <View style={styles.sectionHeader}>
            <Text style={styles.sectionTitle}>Água hoje</Text>
            <TouchableOpacity onPress={() => navigation.navigate("Water")}>
              <Text style={styles.seeAll}>Ver detalhes</Text>
            </TouchableOpacity>
          </View>
          <Card style={styles.waterCard}>
            <View style={styles.waterHeader}>
              <View style={styles.waterIconContainer}>
                <Ionicons name="water" size={24} color="#0891B2" />
              </View>
              <View style={styles.waterSummary}>
                <Text style={styles.waterLabel}>Consumido / Meta</Text>
                <Text style={styles.waterValue}>
                  {(waterSummary ? waterSummary.totalMl / 1000 : 0).toFixed(1)} L /{" "}
                  {(waterSummary ? waterSummary.goalMl / 1000 : 2).toFixed(1)} L
                </Text>
              </View>
            </View>
            <View style={styles.waterProgressBar}>
              <View
                style={[
                  styles.waterProgressFill,
                  {
                    width: `${
                      waterSummary && waterSummary.goalMl > 0
                        ? Math.min(
                            100,
                            (waterSummary.totalMl / waterSummary.goalMl) * 100
                          )
                        : 0
                    }%`,
                  },
                ]}
              />
            </View>
            <View style={styles.waterQuickButtons}>
              <TouchableOpacity
                style={[
                  styles.waterQuickButton,
                  waterAdding && styles.waterQuickButtonDisabled,
                ]}
                onPress={() => addWater(250)}
                disabled={waterAdding}
              >
                <Text style={styles.waterQuickButtonText}>+250 ml</Text>
              </TouchableOpacity>
              <TouchableOpacity
                style={[
                  styles.waterQuickButton,
                  waterAdding && styles.waterQuickButtonDisabled,
                ]}
                onPress={() => addWater(500)}
                disabled={waterAdding}
              >
                <Text style={styles.waterQuickButtonText}>+500 ml</Text>
              </TouchableOpacity>
            </View>
          </Card>
        </View>

        {/* Today's Workouts */}
        <View style={styles.section}>
          <View style={styles.sectionHeader}>
            <Text style={styles.sectionTitle}>Treinos de Hoje</Text>
            <Text style={styles.sectionDay}>
              {DAYS_OF_WEEK.find((d) => d.key === today)?.label}
            </Text>
          </View>

          {todaysWorkouts.length > 0 ? (
            todaysWorkouts.map((workout, index) => (
              <WorkoutCard
                key={workout.id}
                workout={workout}
                variant={index === 0 ? "featured" : "default"}
                onPress={() =>
                  navigation.navigate("WorkoutDetail", {
                    workoutId: workout.id,
                  })
                }
                isCompletedToday={isWorkoutCompletedToday(workout.id)}
              />
            ))
          ) : (
            <Card style={styles.emptyCard}>
              <Ionicons
                name="calendar-outline"
                size={48}
                color={colors.textLight}
              />
              <Text style={styles.emptyText}>
                Nenhum treino programado para hoje
              </Text>
              <TouchableOpacity
                onPress={() => navigation.navigate("WorkoutForm", {})}
                style={styles.addButton}
              >
                <Text style={styles.addButtonText}>Criar Treino</Text>
              </TouchableOpacity>
            </Card>
          )}
        </View>

        {/* Workout Calendar */}
        <View style={styles.section}>
          <Text style={styles.sectionTitle}>Calendário de Treinos</Text>
          <WorkoutCalendar
            completedDates={completedWorkouts.map((cw) => cw.completedAt)}
          />
        </View>

        {/* Bioimpedance Summary */}
        {latestBioimpedance && (
          <View style={styles.section}>
            <View style={styles.sectionHeader}>
              <Text style={styles.sectionTitle}>Bioimpedância</Text>
              <TouchableOpacity
                onPress={() => navigation.navigate("BioimpedanceHistory")}
              >
                <Text style={styles.seeAll}>Ver Histórico</Text>
              </TouchableOpacity>
            </View>

            <Card style={styles.bioCard}>
              <View style={styles.bioHeader}>
                <View>
                  <Text style={styles.bioDate}>
                    {new Date(latestBioimpedance.date).toLocaleDateString(
                      "pt-BR"
                    )}
                  </Text>
                  <Text style={styles.bioLabel}>Última medição</Text>
                </View>
                <TouchableOpacity
                  onPress={() => navigation.navigate("BioimpedanceForm", {})}
                  style={styles.bioAddButton}
                >
                  <Ionicons name="add" size={20} color={colors.primary} />
                </TouchableOpacity>
              </View>

              <View style={styles.bioGrid}>
                <View style={styles.bioItem}>
                  <Text style={styles.bioValue}>
                    {latestBioimpedance.weight}
                  </Text>
                  <Text style={styles.bioUnit}>kg</Text>
                  <Text style={styles.bioItemLabel}>Peso</Text>
                </View>
                <View style={styles.bioItem}>
                  <Text style={styles.bioValue}>
                    {latestBioimpedance.bodyFatPercentage}
                  </Text>
                  <Text style={styles.bioUnit}>%</Text>
                  <Text style={styles.bioItemLabel}>Gordura</Text>
                </View>
                <View style={styles.bioItem}>
                  <Text style={styles.bioValue}>
                    {latestBioimpedance.muscleMass}
                  </Text>
                  <Text style={styles.bioUnit}>kg</Text>
                  <Text style={styles.bioItemLabel}>Músculo</Text>
                </View>
                <View style={styles.bioItem}>
                  <Text style={styles.bioValue}>
                    {latestBioimpedance.waterPercentage}
                  </Text>
                  <Text style={styles.bioUnit}>%</Text>
                  <Text style={styles.bioItemLabel}>Água</Text>
                </View>
              </View>
            </Card>
          </View>
        )}

        {/* Quick Actions */}
        <View style={styles.section}>
          <Text style={styles.sectionTitle}>Ações Rápidas</Text>
          <View style={styles.quickActions}>
            <TouchableOpacity
              onPress={() => navigation.navigate("WorkoutForm", {})}
              style={styles.quickAction}
            >
              <LinearGradient
                colors={["#FF6B6B", "#FF8E8E"]}
                style={styles.quickActionGradient}
              >
                <Ionicons name="add-circle" size={28} color={colors.white} />
              </LinearGradient>
              <Text style={styles.quickActionText}>Novo Treino</Text>
            </TouchableOpacity>

            <TouchableOpacity
              onPress={() => navigation.navigate("BioimpedanceForm", {})}
              style={styles.quickAction}
            >
              <LinearGradient
                colors={["#4ECDC4", "#6EE7DF"]}
                style={styles.quickActionGradient}
              >
                <Ionicons name="body" size={28} color={colors.white} />
              </LinearGradient>
              <Text style={styles.quickActionText}>Bioimpedância</Text>
            </TouchableOpacity>

            <TouchableOpacity
              onPress={() => navigation.navigate("WorkoutPlanning")}
              style={styles.quickAction}
            >
              <LinearGradient
                colors={[colors.accent, "#FBBF24"]}
                style={styles.quickActionGradient}
              >
                <Ionicons name="sparkles" size={28} color={colors.white} />
              </LinearGradient>
              <Text style={styles.quickActionText}>AI Planner</Text>
            </TouchableOpacity>

            <TouchableOpacity
              onPress={() => navigation.navigate("Water")}
              style={styles.quickAction}
            >
              <LinearGradient
                colors={["#0891B2", "#22D3EE"]}
                style={styles.quickActionGradient}
              >
                <Ionicons name="water" size={28} color={colors.white} />
              </LinearGradient>
              <Text style={styles.quickActionText}>Água</Text>
            </TouchableOpacity>

            <TouchableOpacity
              onPress={() => navigation.navigate("Profile")}
              style={styles.quickAction}
            >
              <LinearGradient
                colors={["#45B7D1", "#6DD5ED"]}
                style={styles.quickActionGradient}
              >
                <Ionicons name="person" size={28} color={colors.white} />
              </LinearGradient>
              <Text style={styles.quickActionText}>Perfil</Text>
            </TouchableOpacity>
          </View>
        </View>

        {/* Espaço extra para não ser coberto pelo navbar */}
        <View style={{ height: 100 }} />
      </ScrollView>
    </SafeAreaView>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: colors.background,
  },
  loadingContainer: {
    flex: 1,
    alignItems: "center",
    justifyContent: "center",
    gap: spacing.md,
  },
  loadingText: {
    fontSize: typography.md,
    color: colors.textSecondary,
  },
  errorBanner: {
    flexDirection: "row",
    alignItems: "center",
    backgroundColor: colors.error,
    marginHorizontal: spacing.lg,
    marginTop: spacing.sm,
    padding: spacing.md,
    borderRadius: borderRadius.md,
  },
  errorIcon: {
    marginRight: spacing.sm,
  },
  errorText: {
    flex: 1,
    fontSize: typography.sm,
    color: colors.white,
  },
  header: {
    flexDirection: "row",
    justifyContent: "space-between",
    alignItems: "center",
    paddingHorizontal: spacing.lg,
    paddingTop: spacing.md,
    paddingBottom: spacing.lg,
  },
  greeting: {
    fontSize: typography.xxl,
    fontWeight: typography.bold,
    color: colors.textPrimary,
  },
  subtitle: {
    fontSize: typography.md,
    color: colors.textSecondary,
    marginTop: spacing.xs,
  },
  avatarContainer: {
    borderRadius: borderRadius.full,
    ...shadows.md,
  },
  avatar: {
    width: 50,
    height: 50,
    borderRadius: borderRadius.full,
    alignItems: "center",
    justifyContent: "center",
  },
  avatarText: {
    fontSize: typography.xl,
    fontWeight: typography.bold,
    color: colors.white,
  },
  statsContainer: {
    flexDirection: "row",
    paddingHorizontal: spacing.lg,
    gap: spacing.sm,
    marginBottom: spacing.lg,
  },
  statCard: {
    flex: 1,
    alignItems: "center",
    padding: spacing.md,
  },
  statIconContainer: {
    marginBottom: spacing.xs,
  },
  statNumber: {
    fontSize: typography.xxl,
    fontWeight: typography.bold,
    color: colors.textPrimary,
  },
  statLabel: {
    fontSize: typography.xs,
    color: colors.textSecondary,
    textAlign: "center",
    marginTop: spacing.xs,
  },
  section: {
    paddingHorizontal: spacing.lg,
    marginBottom: spacing.lg,
  },
  sectionHeader: {
    flexDirection: "row",
    justifyContent: "space-between",
    alignItems: "center",
    marginBottom: spacing.md,
  },
  sectionTitle: {
    fontSize: typography.lg,
    fontWeight: typography.bold,
    color: colors.textPrimary,
  },
  sectionDay: {
    fontSize: typography.sm,
    color: colors.primary,
    fontWeight: typography.medium,
  },
  seeAll: {
    fontSize: typography.sm,
    color: colors.primary,
    fontWeight: typography.medium,
  },
  emptyCard: {
    alignItems: "center",
    padding: spacing.xl,
    gap: spacing.md,
  },
  emptyText: {
    fontSize: typography.md,
    color: colors.textSecondary,
    textAlign: "center",
  },
  addButton: {
    paddingHorizontal: spacing.lg,
    paddingVertical: spacing.sm,
    backgroundColor: colors.primary,
    borderRadius: borderRadius.full,
  },
  addButtonText: {
    color: colors.white,
    fontWeight: typography.semibold,
  },
  bioCard: {
    padding: spacing.lg,
  },
  bioHeader: {
    flexDirection: "row",
    justifyContent: "space-between",
    alignItems: "flex-start",
    marginBottom: spacing.lg,
  },
  bioDate: {
    fontSize: typography.md,
    fontWeight: typography.semibold,
    color: colors.textPrimary,
  },
  bioLabel: {
    fontSize: typography.xs,
    color: colors.textSecondary,
  },
  bioAddButton: {
    width: 36,
    height: 36,
    borderRadius: borderRadius.full,
    backgroundColor: colors.surfaceSecondary,
    alignItems: "center",
    justifyContent: "center",
  },
  bioGrid: {
    flexDirection: "row",
    justifyContent: "space-between",
  },
  bioItem: {
    alignItems: "center",
  },
  bioValue: {
    fontSize: typography.xl,
    fontWeight: typography.bold,
    color: colors.textPrimary,
  },
  bioUnit: {
    fontSize: typography.xs,
    color: colors.textSecondary,
  },
  bioItemLabel: {
    fontSize: typography.xs,
    color: colors.textLight,
    marginTop: spacing.xs,
  },
  quickActions: {
    flexDirection: "row",
    gap: spacing.md,
    marginTop: spacing.sm,
  },
  quickAction: {
    flex: 1,
    alignItems: "center",
    gap: spacing.sm,
  },
  quickActionGradient: {
    width: 60,
    height: 60,
    borderRadius: borderRadius.lg,
    alignItems: "center",
    justifyContent: "center",
    ...shadows.md,
  },
  quickActionText: {
    fontSize: typography.xs,
    color: colors.textSecondary,
    fontWeight: typography.medium,
  },
  waterCard: {
    padding: spacing.lg,
  },
  waterHeader: {
    flexDirection: "row",
    alignItems: "center",
    marginBottom: spacing.md,
  },
  waterIconContainer: {
    width: 44,
    height: 44,
    borderRadius: borderRadius.md,
    backgroundColor: "#E0F2F1",
    alignItems: "center",
    justifyContent: "center",
  },
  waterSummary: {
    marginLeft: spacing.md,
  },
  waterLabel: {
    fontSize: typography.xs,
    color: colors.textSecondary,
  },
  waterValue: {
    fontSize: typography.xl,
    fontWeight: typography.bold,
    color: colors.textPrimary,
  },
  waterProgressBar: {
    height: 6,
    backgroundColor: colors.surfaceSecondary,
    borderRadius: 3,
    overflow: "hidden",
    marginBottom: spacing.md,
  },
  waterProgressFill: {
    height: "100%",
    backgroundColor: "#0891B2",
    borderRadius: 3,
  },
  waterQuickButtons: {
    flexDirection: "row",
    gap: spacing.sm,
  },
  waterQuickButton: {
    flex: 1,
    paddingVertical: spacing.sm,
    paddingHorizontal: spacing.md,
    backgroundColor: colors.surfaceSecondary,
    borderRadius: borderRadius.md,
    alignItems: "center",
  },
  waterQuickButtonDisabled: {
    opacity: 0.6,
  },
  waterQuickButtonText: {
    fontSize: typography.sm,
    fontWeight: typography.semibold,
    color: colors.textPrimary,
  },
});
