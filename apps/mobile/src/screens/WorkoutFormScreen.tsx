import React, { useState } from "react";
import {
  View,
  Text,
  StyleSheet,
  ScrollView,
  TouchableOpacity,
  Alert,
  Modal,
  KeyboardAvoidingView,
  Platform,
} from "react-native";
import { SafeAreaView } from "react-native-safe-area-context";
import { Ionicons } from "@expo/vector-icons";
import { useNavigation, useRoute, RouteProp } from "@react-navigation/native";
import { NativeStackNavigationProp } from "@react-navigation/native-stack";
import { generateId } from "../utils";

import { useApp } from "../context/AppContext";
import {
  Button,
  Input,
  DaySelector,
  GoalSelector,
  ExerciseCard,
  Card,
} from "../components";
import { colors, spacing, borderRadius, typography } from "../theme";
import { RootStackParamList } from "../navigation/types";
import {
  Exercise,
  DayOfWeek,
  WorkoutGoal,
  MuscleGroup,
  MUSCLE_GROUPS,
} from "../types";

type WorkoutFormScreenNavigationProp = NativeStackNavigationProp<
  RootStackParamList,
  "WorkoutForm"
>;
type WorkoutFormScreenRouteProp = RouteProp<RootStackParamList, "WorkoutForm">;

export const WorkoutFormScreen: React.FC = () => {
  const navigation = useNavigation<WorkoutFormScreenNavigationProp>();
  const route = useRoute<WorkoutFormScreenRouteProp>();
  const { workouts, addWorkout, updateWorkout } = useApp();

  const editingWorkout = route.params?.workoutId
    ? workouts.find((w) => w.id === route.params.workoutId)
    : undefined;

  // Form state
  const [name, setName] = useState(editingWorkout?.name || "");
  const [description, setDescription] = useState(
    editingWorkout?.description || "",
  );
  const [goal, setGoal] = useState<WorkoutGoal | null>(
    editingWorkout?.goal || null,
  );
  const [days, setDays] = useState<DayOfWeek[]>(editingWorkout?.days || []);
  const [exercises, setExercises] = useState<Exercise[]>(
    editingWorkout?.exercises || [],
  );

  // Exercise modal state
  const [showExerciseModal, setShowExerciseModal] = useState(false);
  const [editingExercise, setEditingExercise] = useState<Exercise | null>(null);
  const [exerciseName, setExerciseName] = useState("");
  const [exerciseMuscleGroup, setExerciseMuscleGroup] =
    useState<MuscleGroup>("chest");
  const [exerciseSets, setExerciseSets] = useState("3");
  const [exerciseReps, setExerciseReps] = useState("12");
  const [exerciseWeight, setExerciseWeight] = useState("");
  const [exerciseRest, setExerciseRest] = useState("60");
  const [exerciseNotes, setExerciseNotes] = useState("");

  const [isSubmitting, setIsSubmitting] = useState(false);

  const handleToggleDay = (day: DayOfWeek) => {
    setDays((prev) =>
      prev.includes(day) ? prev.filter((d) => d !== day) : [...prev, day],
    );
  };

  const openExerciseModal = (exercise?: Exercise) => {
    if (exercise) {
      setEditingExercise(exercise);
      setExerciseName(exercise.name);
      setExerciseMuscleGroup(exercise.muscleGroup);
      setExerciseSets(String(exercise.sets));
      setExerciseReps(String(exercise.reps));
      setExerciseWeight(
        exercise.weight !== null ? String(exercise.weight) : "",
      );
      setExerciseRest(String(exercise.restSeconds));
      setExerciseNotes(exercise.notes || "");
    } else {
      setEditingExercise(null);
      setExerciseName("");
      setExerciseMuscleGroup("chest");
      setExerciseSets("3");
      setExerciseReps("12");
      setExerciseWeight("");
      setExerciseRest("60");
      setExerciseNotes("");
    }
    setShowExerciseModal(true);
  };

  const handleSaveExercise = () => {
    if (!exerciseName.trim()) {
      Alert.alert("Erro", "Digite o nome do exercício");
      return;
    }

    const exerciseData: Exercise = {
      id: editingExercise?.id || generateId(),
      name: exerciseName.trim(),
      muscleGroup: exerciseMuscleGroup,
      sets: parseInt(exerciseSets) || 3,
      reps: parseInt(exerciseReps) || 12,
      weight: exerciseWeight ? parseFloat(exerciseWeight) : null,
      restSeconds: parseInt(exerciseRest) || 60,
      notes: exerciseNotes.trim() || undefined,
    };

    if (editingExercise) {
      setExercises((prev) =>
        prev.map((e) => (e.id === editingExercise.id ? exerciseData : e)),
      );
    } else {
      setExercises((prev) => [...prev, exerciseData]);
    }

    setShowExerciseModal(false);
  };

  const handleDeleteExercise = (exerciseId: string) => {
    Alert.alert(
      "Remover Exercício",
      "Deseja remover este exercício do treino?",
      [
        { text: "Cancelar", style: "cancel" },
        {
          text: "Remover",
          style: "destructive",
          onPress: () =>
            setExercises((prev) => prev.filter((e) => e.id !== exerciseId)),
        },
      ],
    );
  };

  const handleSubmit = async () => {
    // Validation
    if (!name.trim()) {
      Alert.alert("Erro", "Digite o nome do treino");
      return;
    }
    if (!goal) {
      Alert.alert("Erro", "Selecione um objetivo");
      return;
    }
    if (days.length === 0) {
      Alert.alert("Erro", "Selecione pelo menos um dia da semana");
      return;
    }
    if (exercises.length === 0) {
      Alert.alert("Erro", "Adicione pelo menos um exercício");
      return;
    }

    setIsSubmitting(true);

    try {
      if (editingWorkout) {
        await updateWorkout({
          ...editingWorkout,
          name: name.trim(),
          description: description.trim() || undefined,
          goal,
          days,
          exercises,
        });
      } else {
        await addWorkout({
          name: name.trim(),
          description: description.trim() || undefined,
          goal,
          days,
          exercises,
        });
      }
      navigation.goBack();
    } catch (error) {
      Alert.alert("Erro", "Não foi possível salvar o treino");
    } finally {
      setIsSubmitting(false);
    }
  };

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
        <Text style={styles.title}>
          {editingWorkout ? "Editar Treino" : "Novo Treino"}
        </Text>
        <View style={{ width: 40 }} />
      </View>

      <KeyboardAvoidingView
        behavior={Platform.OS === "ios" ? "padding" : "height"}
        style={{ flex: 1 }}
      >
        <ScrollView
          style={styles.content}
          contentContainerStyle={styles.contentContainer}
          showsVerticalScrollIndicator={false}
        >
          {/* Basic Info */}
          <Input
            label="Nome do Treino"
            placeholder="Ex: Treino A - Peito e Tríceps"
            value={name}
            onChangeText={setName}
            required
          />

          <Input
            label="Descrição"
            placeholder="Descreva seu treino (opcional)"
            value={description}
            onChangeText={setDescription}
            multiline
            numberOfLines={3}
          />

          {/* Goal */}
          <GoalSelector
            label="Objetivo *"
            selectedGoal={goal}
            onSelectGoal={setGoal}
          />

          {/* Days */}
          <DaySelector
            label="Dias da Semana *"
            selectedDays={days}
            onToggleDay={handleToggleDay}
          />

          {/* Exercises */}
          <View style={styles.exercisesSection}>
            <View style={styles.exercisesHeader}>
              <Text style={styles.sectionLabel}>Exercícios *</Text>
              <TouchableOpacity
                onPress={() => openExerciseModal()}
                style={styles.addExerciseButton}
              >
                <Ionicons name="add" size={20} color={colors.primary} />
                <Text style={styles.addExerciseText}>Adicionar</Text>
              </TouchableOpacity>
            </View>

            {exercises.length > 0 ? (
              exercises.map((exercise, index) => (
                <ExerciseCard
                  key={exercise.id}
                  exercise={exercise}
                  index={index}
                  onPress={() => openExerciseModal(exercise)}
                  onDelete={() => handleDeleteExercise(exercise.id)}
                />
              ))
            ) : (
              <Card style={styles.emptyExercises}>
                <Ionicons
                  name="barbell-outline"
                  size={40}
                  color={colors.textLight}
                />
                <Text style={styles.emptyText}>
                  Nenhum exercício adicionado
                </Text>
              </Card>
            )}
          </View>

          <View style={{ height: spacing.xl }} />
        </ScrollView>

        {/* Submit Button */}
        <View style={styles.footer}>
          <Button
            title={editingWorkout ? "Salvar Alterações" : "Criar Treino"}
            onPress={handleSubmit}
            loading={isSubmitting}
            fullWidth
          />
        </View>
      </KeyboardAvoidingView>

      {/* Exercise Modal */}
      <Modal
        visible={showExerciseModal}
        animationType="slide"
        presentationStyle="pageSheet"
        onRequestClose={() => setShowExerciseModal(false)}
      >
        <SafeAreaView style={styles.modalContainer}>
          <View style={styles.modalHeader}>
            <TouchableOpacity onPress={() => setShowExerciseModal(false)}>
              <Text style={styles.modalCancel}>Cancelar</Text>
            </TouchableOpacity>
            <Text style={styles.modalTitle}>
              {editingExercise ? "Editar Exercício" : "Novo Exercício"}
            </Text>
            <TouchableOpacity onPress={handleSaveExercise}>
              <Text style={styles.modalSave}>Salvar</Text>
            </TouchableOpacity>
          </View>

          <ScrollView
            style={styles.modalContent}
            contentContainerStyle={styles.modalContentContainer}
          >
            <Input
              label="Nome do Exercício"
              placeholder="Ex: Supino Reto"
              value={exerciseName}
              onChangeText={setExerciseName}
              required
            />

            <Text style={styles.inputLabel}>Grupo Muscular *</Text>
            <ScrollView
              horizontal
              showsHorizontalScrollIndicator={false}
              style={styles.muscleGroupScroll}
            >
              {MUSCLE_GROUPS.map((group) => (
                <TouchableOpacity
                  key={group.key}
                  onPress={() => setExerciseMuscleGroup(group.key)}
                  style={[
                    styles.muscleGroupChip,
                    exerciseMuscleGroup === group.key && {
                      backgroundColor: group.color,
                    },
                  ]}
                >
                  <Text
                    style={[
                      styles.muscleGroupChipText,
                      exerciseMuscleGroup === group.key &&
                        styles.muscleGroupChipTextActive,
                    ]}
                  >
                    {group.label}
                  </Text>
                </TouchableOpacity>
              ))}
            </ScrollView>

            <View style={styles.row}>
              <View style={styles.halfWidth}>
                <Input
                  label="Séries"
                  placeholder="3"
                  value={exerciseSets}
                  onChangeText={setExerciseSets}
                  keyboardType="numeric"
                />
              </View>
              <View style={styles.halfWidth}>
                <Input
                  label="Repetições"
                  placeholder="12"
                  value={exerciseReps}
                  onChangeText={setExerciseReps}
                  keyboardType="numeric"
                />
              </View>
            </View>

            <View style={styles.row}>
              <View style={styles.halfWidth}>
                <Input
                  label="Peso (kg)"
                  placeholder="Livre"
                  value={exerciseWeight}
                  onChangeText={setExerciseWeight}
                  keyboardType="decimal-pad"
                  hint="Deixe vazio para peso livre"
                />
              </View>
              <View style={styles.halfWidth}>
                <Input
                  label="Descanso (seg)"
                  placeholder="60"
                  value={exerciseRest}
                  onChangeText={setExerciseRest}
                  keyboardType="numeric"
                />
              </View>
            </View>

            <Input
              label="Observações"
              placeholder="Dicas, variações, etc."
              value={exerciseNotes}
              onChangeText={setExerciseNotes}
              multiline
              numberOfLines={2}
            />
          </ScrollView>
        </SafeAreaView>
      </Modal>
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
  },
  sectionLabel: {
    fontSize: typography.sm,
    fontWeight: typography.medium,
    color: colors.textPrimary,
  },
  exercisesSection: {
    marginTop: spacing.md,
  },
  exercisesHeader: {
    flexDirection: "row",
    justifyContent: "space-between",
    alignItems: "center",
    marginBottom: spacing.sm,
  },
  addExerciseButton: {
    flexDirection: "row",
    alignItems: "center",
    gap: spacing.xs,
  },
  addExerciseText: {
    fontSize: typography.sm,
    color: colors.primary,
    fontWeight: typography.medium,
  },
  emptyExercises: {
    alignItems: "center",
    padding: spacing.xl,
    gap: spacing.sm,
  },
  emptyText: {
    fontSize: typography.sm,
    color: colors.textSecondary,
  },
  footer: {
    padding: spacing.lg,
    backgroundColor: colors.surface,
    borderTopWidth: 1,
    borderTopColor: colors.border,
  },

  // Modal styles
  modalContainer: {
    flex: 1,
    backgroundColor: colors.background,
  },
  modalHeader: {
    flexDirection: "row",
    alignItems: "center",
    justifyContent: "space-between",
    paddingHorizontal: spacing.lg,
    paddingVertical: spacing.md,
    borderBottomWidth: 1,
    borderBottomColor: colors.border,
    backgroundColor: colors.surface,
  },
  modalCancel: {
    fontSize: typography.md,
    color: colors.textSecondary,
  },
  modalTitle: {
    fontSize: typography.lg,
    fontWeight: typography.semibold,
    color: colors.textPrimary,
  },
  modalSave: {
    fontSize: typography.md,
    color: colors.primary,
    fontWeight: typography.semibold,
  },
  modalContent: {
    flex: 1,
  },
  modalContentContainer: {
    padding: spacing.lg,
  },
  inputLabel: {
    fontSize: typography.sm,
    fontWeight: typography.medium,
    color: colors.textPrimary,
    marginBottom: spacing.sm,
  },
  muscleGroupScroll: {
    marginBottom: spacing.md,
  },
  muscleGroupChip: {
    paddingHorizontal: spacing.md,
    paddingVertical: spacing.sm,
    borderRadius: borderRadius.full,
    backgroundColor: colors.surfaceSecondary,
    marginRight: spacing.sm,
  },
  muscleGroupChipText: {
    fontSize: typography.sm,
    color: colors.textSecondary,
    fontWeight: typography.medium,
  },
  muscleGroupChipTextActive: {
    color: colors.white,
    fontWeight: typography.bold,
  },
  row: {
    flexDirection: "row",
    gap: spacing.md,
  },
  halfWidth: {
    flex: 1,
  },
});
