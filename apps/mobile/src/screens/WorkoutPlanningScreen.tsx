import React, { useState, useEffect } from "react";
import {
  View,
  Text,
  StyleSheet,
  ScrollView,
  TouchableOpacity,
  Alert,
  ActivityIndicator,
} from "react-native";
import { SafeAreaView } from "react-native-safe-area-context";
import { Ionicons } from "@expo/vector-icons";
import { useNavigation } from "@react-navigation/native";
import { NativeStackNavigationProp } from "@react-navigation/native-stack";

import { useApp } from "../context/AppContext";
import { Button, Card, GoalSelector } from "../components";
import { colors, spacing, borderRadius, typography, shadows } from "../theme";
import { RootStackParamList } from "../navigation/types";
import { WorkoutGoal, MUSCLE_GROUPS } from "../types";
import {
  aiPlanningApi,
  GeneratedPlanResponse,
  GeneratedWorkout,
  AIPlanningStatus,
} from "../services/api";

type WorkoutPlanningScreenNavigationProp = NativeStackNavigationProp<
  RootStackParamList,
  "WorkoutPlanning"
>;

type EquipmentType = "Gym" | "Home" | "Minimal";

const EQUIPMENT_OPTIONS: { key: EquipmentType; label: string; icon: string }[] = [
  { key: "Gym", label: "Academia", icon: "fitness" },
  { key: "Home", label: "Casa", icon: "home" },
  { key: "Minimal", label: "Mínimo", icon: "body" },
];

const DAYS_OPTIONS = [2, 3, 4, 5, 6];

export const WorkoutPlanningScreen: React.FC = () => {
  const navigation = useNavigation<WorkoutPlanningScreenNavigationProp>();
  const { profileId, refreshData } = useApp();

  // AI Status
  const [aiStatus, setAiStatus] = useState<AIPlanningStatus | null>(null);
  const [checkingStatus, setCheckingStatus] = useState(true);

  // Form state
  const [goal, setGoal] = useState<WorkoutGoal | null>("hypertrophy");
  const [daysPerWeek, setDaysPerWeek] = useState(3);
  const [equipment, setEquipment] = useState<EquipmentType>("Gym");
  const [additionalNotes, setAdditionalNotes] = useState("");

  // Generation state
  const [isGenerating, setIsGenerating] = useState(false);
  const [generatedPlan, setGeneratedPlan] = useState<GeneratedPlanResponse | null>(null);
  const [isSaving, setIsSaving] = useState(false);

  useEffect(() => {
    checkAIStatus();
  }, []);

  const checkAIStatus = async () => {
    try {
      const status = await aiPlanningApi.getStatus();
      setAiStatus(status);
    } catch (error) {
      console.error("Failed to check AI status:", error);
      setAiStatus({ available: false, provider: "unknown", model: "unknown", endpoint: "" });
    } finally {
      setCheckingStatus(false);
    }
  };

  const handleGenerate = async () => {
    if (!profileId) {
      Alert.alert("Erro", "Perfil não encontrado");
      return;
    }

    if (!goal) {
      Alert.alert("Erro", "Selecione um objetivo");
      return;
    }

    setIsGenerating(true);
    setGeneratedPlan(null);

    try {
      const plan = await aiPlanningApi.generate({
        userProfileId: profileId,
        goal,
        overrideDaysPerWeek: daysPerWeek,
        overrideEquipment: equipment,
        additionalNotes: additionalNotes.trim() || undefined,
      });
      setGeneratedPlan(plan);
    } catch (error) {
      console.error("Failed to generate plan:", error);
      Alert.alert(
        "Erro",
        "Não foi possível gerar o plano. Verifique se o Ollama está rodando."
      );
    } finally {
      setIsGenerating(false);
    }
  };

  const handleAcceptPlan = async () => {
    if (!profileId || !generatedPlan) return;

    setIsSaving(true);

    try {
      await aiPlanningApi.accept({
        planId: generatedPlan.planId,
        userProfileId: profileId,
        workouts: generatedPlan.workouts,
      });

      await refreshData();

      Alert.alert(
        "Sucesso! 🎉",
        `${generatedPlan.workouts.length} treinos foram criados com sucesso!`,
        [{ text: "OK", onPress: () => navigation.goBack() }]
      );
    } catch (error) {
      console.error("Failed to save plan:", error);
      Alert.alert("Erro", "Não foi possível salvar o plano");
    } finally {
      setIsSaving(false);
    }
  };

  const handleReject = () => {
    Alert.alert(
      "Descartar Plano",
      "Deseja descartar este plano e gerar outro?",
      [
        { text: "Cancelar", style: "cancel" },
        {
          text: "Descartar",
          style: "destructive",
          onPress: () => setGeneratedPlan(null),
        },
      ]
    );
  };

  const getMuscleGroupColor = (muscleGroup: string): string => {
    const group = MUSCLE_GROUPS.find(
      (g) => g.key.toLowerCase() === muscleGroup.toLowerCase()
    );
    return group?.color || colors.primary;
  };

  const formatDays = (days: string[]): string => {
    const dayMap: Record<string, string> = {
      monday: "Seg",
      tuesday: "Ter",
      wednesday: "Qua",
      thursday: "Qui",
      friday: "Sex",
      saturday: "Sáb",
      sunday: "Dom",
    };
    return days.map((d) => dayMap[d.toLowerCase()] || d).join(", ");
  };

  if (checkingStatus) {
    return (
      <SafeAreaView style={styles.container} edges={["top"]}>
        <View style={styles.loadingContainer}>
          <ActivityIndicator size="large" color={colors.primary} />
          <Text style={styles.loadingText}>Verificando AI...</Text>
        </View>
      </SafeAreaView>
    );
  }

  return (
    <SafeAreaView style={styles.container} edges={["top"]}>
      {/* Header */}
      <View style={styles.header}>
        <TouchableOpacity
          onPress={() => navigation.goBack()}
          style={styles.backButton}
        >
          <Ionicons name="arrow-back" size={24} color={colors.textPrimary} />
        </TouchableOpacity>
        <View style={styles.headerTitle}>
          <Ionicons name="sparkles" size={20} color={colors.accent} />
          <Text style={styles.title}>Planejamento AI</Text>
        </View>
        <View style={{ width: 40 }} />
      </View>

      <ScrollView
        style={styles.content}
        contentContainerStyle={styles.contentContainer}
        showsVerticalScrollIndicator={false}
      >
        {/* AI Status Banner */}
        {!aiStatus?.available && (
          <Card style={styles.warningBanner}>
            <Ionicons name="warning" size={24} color={colors.warning} />
            <View style={styles.warningContent}>
              <Text style={styles.warningTitle}>Ollama não detectado</Text>
              <Text style={styles.warningText}>
                Certifique-se de que o Ollama está rodando em{" "}
                {aiStatus?.endpoint || "localhost:11434"}
              </Text>
            </View>
          </Card>
        )}

        {/* Generated Plan View */}
        {generatedPlan ? (
          <View style={styles.planContainer}>
            {/* Plan Summary */}
            <Card style={styles.summaryCard}>
              <View style={styles.summaryHeader}>
                <Ionicons name="checkmark-circle" size={32} color={colors.success} />
                <Text style={styles.summaryTitle}>Plano Gerado!</Text>
              </View>
              <Text style={styles.summaryText}>{generatedPlan.summary}</Text>
              <View style={styles.rationaleContainer}>
                <Text style={styles.rationaleLabel}>Por que este plano?</Text>
                <Text style={styles.rationaleText}>{generatedPlan.rationale}</Text>
              </View>
            </Card>

            {/* Workouts Preview */}
            <Text style={styles.sectionTitle}>
              Treinos ({generatedPlan.workouts.length})
            </Text>

            {generatedPlan.workouts.map((workout, index) => (
              <WorkoutPreviewCard
                key={index}
                workout={workout}
                getMuscleGroupColor={getMuscleGroupColor}
                formatDays={formatDays}
              />
            ))}

            {/* Action Buttons */}
            <View style={styles.actionButtons}>
              <Button
                title="Descartar"
                variant="outline"
                onPress={handleReject}
                style={styles.rejectButton}
              />
              <Button
                title="Aceitar e Criar"
                onPress={handleAcceptPlan}
                loading={isSaving}
                style={styles.acceptButton}
              />
            </View>
          </View>
        ) : (
          /* Configuration Form */
          <View style={styles.formContainer}>
            <Card style={styles.introCard}>
              <Ionicons name="bulb" size={32} color={colors.accent} />
              <Text style={styles.introTitle}>
                Deixe a IA criar seu plano de treino
              </Text>
              <Text style={styles.introText}>
                Com base no seu perfil, objetivo e preferências, nossa IA irá
                criar um plano de treino personalizado para você.
              </Text>
            </Card>

            {/* Goal Selection */}
            <GoalSelector
              label="Objetivo do Treino *"
              selectedGoal={goal}
              onSelectGoal={setGoal}
            />

            {/* Days per Week */}
            <View style={styles.formSection}>
              <Text style={styles.sectionLabel}>Dias por Semana *</Text>
              <View style={styles.daysContainer}>
                {DAYS_OPTIONS.map((days) => (
                  <TouchableOpacity
                    key={days}
                    onPress={() => setDaysPerWeek(days)}
                    style={[
                      styles.dayOption,
                      daysPerWeek === days && styles.dayOptionSelected,
                    ]}
                  >
                    <Text
                      style={[
                        styles.dayOptionText,
                        daysPerWeek === days && styles.dayOptionTextSelected,
                      ]}
                    >
                      {days}x
                    </Text>
                  </TouchableOpacity>
                ))}
              </View>
            </View>

            {/* Equipment */}
            <View style={styles.formSection}>
              <Text style={styles.sectionLabel}>Equipamento Disponível *</Text>
              <View style={styles.equipmentContainer}>
                {EQUIPMENT_OPTIONS.map((opt) => (
                  <TouchableOpacity
                    key={opt.key}
                    onPress={() => setEquipment(opt.key)}
                    style={[
                      styles.equipmentOption,
                      equipment === opt.key && styles.equipmentOptionSelected,
                    ]}
                  >
                    <Ionicons
                      name={opt.icon as any}
                      size={24}
                      color={
                        equipment === opt.key ? colors.white : colors.textSecondary
                      }
                    />
                    <Text
                      style={[
                        styles.equipmentText,
                        equipment === opt.key && styles.equipmentTextSelected,
                      ]}
                    >
                      {opt.label}
                    </Text>
                  </TouchableOpacity>
                ))}
              </View>
            </View>

            {/* Generate Button */}
            <Button
              title={isGenerating ? "Gerando..." : "✨ Gerar Plano com IA"}
              onPress={handleGenerate}
              loading={isGenerating}
              disabled={!aiStatus?.available || isGenerating}
              fullWidth
              style={styles.generateButton}
            />

            {isGenerating && (
              <Text style={styles.generatingHint}>
                Isso pode levar alguns segundos...
              </Text>
            )}
          </View>
        )}
      </ScrollView>
    </SafeAreaView>
  );
};

// Sub-component for workout preview
const WorkoutPreviewCard: React.FC<{
  workout: GeneratedWorkout;
  getMuscleGroupColor: (group: string) => string;
  formatDays: (days: string[]) => string;
}> = ({ workout, getMuscleGroupColor, formatDays }) => {
  const [expanded, setExpanded] = useState(false);

  return (
    <Card style={styles.workoutCard}>
      <TouchableOpacity
        onPress={() => setExpanded(!expanded)}
        style={styles.workoutHeader}
      >
        <View style={styles.workoutInfo}>
          <Text style={styles.workoutName}>{workout.name}</Text>
          <Text style={styles.workoutDays}>{formatDays(workout.days)}</Text>
        </View>
        <View style={styles.workoutMeta}>
          <Text style={styles.workoutDuration}>
            ~{workout.estimatedDurationMinutes} min
          </Text>
          <Ionicons
            name={expanded ? "chevron-up" : "chevron-down"}
            size={20}
            color={colors.textSecondary}
          />
        </View>
      </TouchableOpacity>

      {expanded && (
        <View style={styles.exercisesList}>
          {workout.exercises.map((exercise, index) => (
            <View key={index} style={styles.exerciseItem}>
              <View
                style={[
                  styles.muscleIndicator,
                  { backgroundColor: getMuscleGroupColor(exercise.muscleGroup) },
                ]}
              />
              <View style={styles.exerciseInfo}>
                <Text style={styles.exerciseName}>{exercise.name}</Text>
                <Text style={styles.exerciseDetails}>
                  {exercise.sets}x{exercise.reps}
                  {exercise.suggestedWeight ? ` • ${exercise.suggestedWeight}kg` : ""}
                  {" • "}{exercise.restSeconds}s descanso
                </Text>
                {exercise.notes && (
                  <Text style={styles.exerciseNotes}>{exercise.notes}</Text>
                )}
              </View>
            </View>
          ))}
        </View>
      )}
    </Card>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: colors.background,
  },
  loadingContainer: {
    flex: 1,
    justifyContent: "center",
    alignItems: "center",
    gap: spacing.md,
  },
  loadingText: {
    fontSize: typography.md,
    color: colors.textSecondary,
  },
  header: {
    flexDirection: "row",
    alignItems: "center",
    justifyContent: "space-between",
    paddingHorizontal: spacing.md,
    paddingVertical: spacing.md,
    borderBottomWidth: 1,
    borderBottomColor: colors.border,
    backgroundColor: colors.surface,
  },
  backButton: {
    width: 40,
    height: 40,
    alignItems: "center",
    justifyContent: "center",
  },
  headerTitle: {
    flexDirection: "row",
    alignItems: "center",
    gap: spacing.xs,
  },
  title: {
    fontSize: typography.lg,
    fontWeight: typography.bold,
    color: colors.textPrimary,
  },
  content: {
    flex: 1,
  },
  contentContainer: {
    padding: spacing.lg,
    paddingBottom: spacing.xxl,
  },
  warningBanner: {
    flexDirection: "row",
    alignItems: "flex-start",
    gap: spacing.md,
    backgroundColor: "#FEF3C7",
    marginBottom: spacing.lg,
  },
  warningContent: {
    flex: 1,
  },
  warningTitle: {
    fontSize: typography.md,
    fontWeight: typography.semibold,
    color: colors.textPrimary,
  },
  warningText: {
    fontSize: typography.sm,
    color: colors.textSecondary,
    marginTop: spacing.xs,
  },
  formContainer: {
    gap: spacing.lg,
  },
  introCard: {
    alignItems: "center",
    padding: spacing.xl,
    gap: spacing.md,
  },
  introTitle: {
    fontSize: typography.lg,
    fontWeight: typography.bold,
    color: colors.textPrimary,
    textAlign: "center",
  },
  introText: {
    fontSize: typography.sm,
    color: colors.textSecondary,
    textAlign: "center",
    lineHeight: 20,
  },
  formSection: {
    marginTop: spacing.md,
  },
  sectionLabel: {
    fontSize: typography.sm,
    fontWeight: typography.medium,
    color: colors.textPrimary,
    marginBottom: spacing.sm,
  },
  daysContainer: {
    flexDirection: "row",
    gap: spacing.sm,
  },
  dayOption: {
    flex: 1,
    paddingVertical: spacing.md,
    borderRadius: borderRadius.md,
    backgroundColor: colors.surfaceSecondary,
    alignItems: "center",
  },
  dayOptionSelected: {
    backgroundColor: colors.primary,
  },
  dayOptionText: {
    fontSize: typography.md,
    fontWeight: typography.semibold,
    color: colors.textSecondary,
  },
  dayOptionTextSelected: {
    color: colors.white,
  },
  equipmentContainer: {
    flexDirection: "row",
    gap: spacing.sm,
  },
  equipmentOption: {
    flex: 1,
    paddingVertical: spacing.lg,
    borderRadius: borderRadius.md,
    backgroundColor: colors.surfaceSecondary,
    alignItems: "center",
    gap: spacing.xs,
  },
  equipmentOptionSelected: {
    backgroundColor: colors.primary,
  },
  equipmentText: {
    fontSize: typography.sm,
    fontWeight: typography.medium,
    color: colors.textSecondary,
  },
  equipmentTextSelected: {
    color: colors.white,
  },
  generateButton: {
    marginTop: spacing.xl,
  },
  generatingHint: {
    fontSize: typography.sm,
    color: colors.textSecondary,
    textAlign: "center",
    marginTop: spacing.sm,
  },

  // Plan View
  planContainer: {
    gap: spacing.lg,
  },
  summaryCard: {
    gap: spacing.md,
  },
  summaryHeader: {
    flexDirection: "row",
    alignItems: "center",
    gap: spacing.sm,
  },
  summaryTitle: {
    fontSize: typography.xl,
    fontWeight: typography.bold,
    color: colors.textPrimary,
  },
  summaryText: {
    fontSize: typography.md,
    color: colors.textPrimary,
    lineHeight: 24,
  },
  rationaleContainer: {
    backgroundColor: colors.surfaceSecondary,
    padding: spacing.md,
    borderRadius: borderRadius.md,
    marginTop: spacing.sm,
  },
  rationaleLabel: {
    fontSize: typography.xs,
    fontWeight: typography.semibold,
    color: colors.textSecondary,
    textTransform: "uppercase",
    letterSpacing: 0.5,
  },
  rationaleText: {
    fontSize: typography.sm,
    color: colors.textSecondary,
    marginTop: spacing.xs,
    lineHeight: 20,
  },
  sectionTitle: {
    fontSize: typography.lg,
    fontWeight: typography.bold,
    color: colors.textPrimary,
  },
  workoutCard: {
    padding: 0,
    overflow: "hidden",
  },
  workoutHeader: {
    flexDirection: "row",
    justifyContent: "space-between",
    alignItems: "center",
    padding: spacing.md,
  },
  workoutInfo: {
    flex: 1,
  },
  workoutName: {
    fontSize: typography.md,
    fontWeight: typography.semibold,
    color: colors.textPrimary,
  },
  workoutDays: {
    fontSize: typography.sm,
    color: colors.textSecondary,
    marginTop: 2,
  },
  workoutMeta: {
    flexDirection: "row",
    alignItems: "center",
    gap: spacing.sm,
  },
  workoutDuration: {
    fontSize: typography.sm,
    color: colors.textSecondary,
  },
  exercisesList: {
    borderTopWidth: 1,
    borderTopColor: colors.border,
    paddingVertical: spacing.sm,
  },
  exerciseItem: {
    flexDirection: "row",
    paddingHorizontal: spacing.md,
    paddingVertical: spacing.sm,
    gap: spacing.sm,
  },
  muscleIndicator: {
    width: 4,
    borderRadius: 2,
  },
  exerciseInfo: {
    flex: 1,
  },
  exerciseName: {
    fontSize: typography.sm,
    fontWeight: typography.medium,
    color: colors.textPrimary,
  },
  exerciseDetails: {
    fontSize: typography.xs,
    color: colors.textSecondary,
    marginTop: 2,
  },
  exerciseNotes: {
    fontSize: typography.xs,
    color: colors.textLight,
    fontStyle: "italic",
    marginTop: 2,
  },
  actionButtons: {
    flexDirection: "row",
    gap: spacing.md,
    marginTop: spacing.lg,
  },
  rejectButton: {
    flex: 1,
  },
  acceptButton: {
    flex: 2,
  },
});
