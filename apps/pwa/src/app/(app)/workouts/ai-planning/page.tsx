"use client";

import React, { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import { useAuth } from "../../../../lib/auth-context";
import {
  aiPlanningApi,
  GeneratedPlanResponse,
  GeneratedWorkout,
  AIPlanningStatus,
} from "../../../../lib/api-client";
import type { WorkoutGoal } from "../../../../lib/types";
import { Button } from "../../../../components/ui/Button";
import { Card, CardContent, CardHeader, CardTitle } from "../../../../components/ui/Card";

type EquipmentType = "Gym" | "Home" | "Minimal";

const GOALS: { key: WorkoutGoal; label: string }[] = [
  { key: "hypertrophy", label: "Hipertrofia" },
  { key: "strength", label: "Força" },
  { key: "endurance", label: "Resistência" },
  { key: "weightLoss", label: "Emagrecimento" },
  { key: "maintenance", label: "Manutenção" },
];

const EQUIPMENT_OPTIONS: { key: EquipmentType; label: string; icon: string }[] = [
  { key: "Gym", label: "Academia", icon: "🏋️" },
  { key: "Home", label: "Casa", icon: "🏠" },
  { key: "Minimal", label: "Mínimo", icon: "🧘" },
];

const DAYS_OPTIONS = [2, 3, 4, 5, 6];

const MUSCLE_GROUP_COLORS: Record<string, string> = {
  chest: "#EF4444",
  back: "#3B82F6",
  shoulders: "#8B5CF6",
  biceps: "#F59E0B",
  triceps: "#10B981",
  legs: "#EC4899",
  glutes: "#F97316",
  abs: "#06B6D4",
  cardio: "#84CC16",
};

export default function AIPlanningPage() {
  const router = useRouter();
  const { user } = useAuth();

  // AI Status
  const [aiStatus, setAiStatus] = useState<AIPlanningStatus | null>(null);
  const [checkingStatus, setCheckingStatus] = useState(true);

  // Form state
  const [goal, setGoal] = useState<WorkoutGoal>("hypertrophy");
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
    if (!user?.profileId) {
      alert("Perfil não encontrado");
      return;
    }

    setIsGenerating(true);
    setGeneratedPlan(null);

    try {
      const plan = await aiPlanningApi.generate({
        userProfileId: user.profileId,
        goal,
        overrideDaysPerWeek: daysPerWeek,
        overrideEquipment: equipment,
        additionalNotes: additionalNotes.trim() || undefined,
      });
      setGeneratedPlan(plan);
    } catch (error) {
      console.error("Failed to generate plan:", error);
      alert("Não foi possível gerar o plano. Verifique se o Ollama está rodando.");
    } finally {
      setIsGenerating(false);
    }
  };

  const handleAcceptPlan = async () => {
    if (!user?.profileId || !generatedPlan) return;

    setIsSaving(true);

    try {
      await aiPlanningApi.accept({
        planId: generatedPlan.planId,
        userProfileId: user.profileId,
        workouts: generatedPlan.workouts,
      });

      alert(`${generatedPlan.workouts.length} treinos foram criados com sucesso!`);
      router.push("/workouts");
    } catch (error) {
      console.error("Failed to save plan:", error);
      alert("Não foi possível salvar o plano");
    } finally {
      setIsSaving(false);
    }
  };

  const handleReject = () => {
    if (confirm("Deseja descartar este plano e gerar outro?")) {
      setGeneratedPlan(null);
    }
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
      <div className="flex items-center justify-center min-h-[60vh]">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary mx-auto mb-4"></div>
          <p className="text-text-secondary">Verificando AI...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="p-4 md:p-6 max-w-4xl mx-auto pb-24 md:pb-6">
      {/* Header */}
      <div className="flex items-center gap-3 mb-6">
        <button
          onClick={() => router.back()}
          className="p-2 hover:bg-surface-secondary rounded-lg transition-colors"
        >
          <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
          </svg>
        </button>
        <div className="flex items-center gap-2">
          <span className="text-2xl">✨</span>
          <h1 className="text-xl md:text-2xl font-bold text-text-primary">Planejamento AI</h1>
        </div>
      </div>

      {/* AI Status Warning */}
      {!aiStatus?.available && (
        <Card className="mb-6 bg-amber-50 border-amber-200">
          <CardContent className="flex items-start gap-3 p-4">
            <span className="text-2xl">⚠️</span>
            <div>
              <p className="font-semibold text-text-primary">Ollama não detectado</p>
              <p className="text-sm text-text-secondary mt-1">
                Certifique-se de que o Ollama está rodando em {aiStatus?.endpoint || "localhost:11434"}
              </p>
            </div>
          </CardContent>
        </Card>
      )}

      {/* Generated Plan View */}
      {generatedPlan ? (
        <div className="space-y-6">
          {/* Plan Summary */}
          <Card>
            <CardContent className="p-6">
              <div className="flex items-center gap-3 mb-4">
                <span className="text-3xl">✅</span>
                <h2 className="text-xl font-bold text-text-primary">Plano Gerado!</h2>
              </div>
              <p className="text-text-primary mb-4">{generatedPlan.summary}</p>
              <div className="bg-surface-secondary p-4 rounded-lg">
                <p className="text-xs font-semibold text-text-secondary uppercase tracking-wide mb-1">
                  Por que este plano?
                </p>
                <p className="text-sm text-text-secondary">{generatedPlan.rationale}</p>
              </div>
            </CardContent>
          </Card>

          {/* Workouts Preview */}
          <h3 className="text-lg font-bold text-text-primary">
            Treinos ({generatedPlan.workouts.length})
          </h3>

          {generatedPlan.workouts.map((workout, index) => (
            <WorkoutPreviewCard
              key={index}
              workout={workout}
              formatDays={formatDays}
            />
          ))}

          {/* Action Buttons */}
          <div className="flex gap-3 pt-4">
            <Button
              variant="outline"
              onClick={handleReject}
              className="flex-1"
            >
              Descartar
            </Button>
            <Button
              onClick={handleAcceptPlan}
              disabled={isSaving}
              className="flex-[2]"
            >
              {isSaving ? "Salvando..." : "Aceitar e Criar"}
            </Button>
          </div>
        </div>
      ) : (
        /* Configuration Form */
        <div className="space-y-6">
          {/* Intro Card */}
          <Card>
            <CardContent className="p-6 text-center">
              <span className="text-4xl mb-4 block">💡</span>
              <h2 className="text-lg font-bold text-text-primary mb-2">
                Deixe a IA criar seu plano de treino
              </h2>
              <p className="text-sm text-text-secondary">
                Com base no seu perfil, objetivo e preferências, nossa IA irá criar um plano de treino personalizado para você.
              </p>
            </CardContent>
          </Card>

          {/* Goal Selection */}
          <div>
            <label className="block text-sm font-medium text-text-primary mb-2">
              Objetivo do Treino *
            </label>
            <div className="grid grid-cols-2 md:grid-cols-3 gap-2">
              {GOALS.map((g) => (
                <button
                  key={g.key}
                  onClick={() => setGoal(g.key)}
                  className="p-3 rounded-lg text-sm font-medium transition-all"
                  style={
                    goal === g.key
                      ? { backgroundColor: "#6366F1", color: "#FFFFFF" }
                      : { backgroundColor: "#F1F5F9", color: "#475569" }
                  }
                >
                  {g.label}
                </button>
              ))}
            </div>
          </div>

          {/* Days per Week */}
          <div>
            <label className="block text-sm font-medium text-text-primary mb-2">
              Dias por Semana *
            </label>
            <div className="flex gap-2">
              {DAYS_OPTIONS.map((days) => (
                <button
                  key={days}
                  onClick={() => setDaysPerWeek(days)}
                  className="flex-1 p-3 rounded-lg text-sm font-semibold transition-all"
                  style={
                    daysPerWeek === days
                      ? { backgroundColor: "#6366F1", color: "#FFFFFF" }
                      : { backgroundColor: "#F1F5F9", color: "#475569" }
                  }
                >
                  {days}x
                </button>
              ))}
            </div>
          </div>

          {/* Equipment */}
          <div>
            <label className="block text-sm font-medium text-text-primary mb-2">
              Equipamento Disponível *
            </label>
            <div className="grid grid-cols-3 gap-2">
              {EQUIPMENT_OPTIONS.map((opt) => (
                <button
                  key={opt.key}
                  onClick={() => setEquipment(opt.key)}
                  className="p-4 rounded-lg text-center transition-all"
                  style={
                    equipment === opt.key
                      ? { backgroundColor: "#6366F1", color: "#FFFFFF" }
                      : { backgroundColor: "#F1F5F9", color: "#475569" }
                  }
                >
                  <span className="text-2xl block mb-1">{opt.icon}</span>
                  <span className="text-sm font-medium">{opt.label}</span>
                </button>
              ))}
            </div>
          </div>

          {/* Additional Notes */}
          <div>
            <label className="block text-sm font-medium text-text-primary mb-2">
              Observações (opcional)
            </label>
            <textarea
              value={additionalNotes}
              onChange={(e) => setAdditionalNotes(e.target.value)}
              placeholder="Ex: Tenho uma lesão no ombro, prefiro exercícios compostos..."
              className="w-full p-3 rounded-lg border border-border bg-surface text-text-primary placeholder:text-text-light resize-none focus:outline-none focus:ring-2 focus:ring-primary"
              rows={3}
            />
          </div>

          {/* Generate Button */}
          <Button
            onClick={handleGenerate}
            disabled={!aiStatus?.available || isGenerating}
            className="w-full"
          >
            {isGenerating ? (
              <span className="flex items-center gap-2">
                <div className="animate-spin rounded-full h-4 w-4 border-2 border-white border-t-transparent"></div>
                Gerando...
              </span>
            ) : (
              "✨ Gerar Plano com IA"
            )}
          </Button>

          {isGenerating && (
            <p className="text-center text-sm text-text-secondary">
              Isso pode levar alguns segundos...
            </p>
          )}
        </div>
      )}
    </div>
  );
}

// Sub-component for workout preview
function WorkoutPreviewCard({
  workout,
  formatDays,
}: {
  workout: GeneratedWorkout;
  formatDays: (days: string[]) => string;
}) {
  const [expanded, setExpanded] = useState(false);

  const getMuscleGroupColor = (muscleGroup: string): string => {
    return MUSCLE_GROUP_COLORS[muscleGroup.toLowerCase()] || "#6366F1";
  };

  return (
    <Card>
      <button
        onClick={() => setExpanded(!expanded)}
        className="w-full p-4 flex justify-between items-center text-left"
      >
        <div>
          <h4 className="font-semibold text-text-primary">{workout.name}</h4>
          <p className="text-sm text-text-secondary">{formatDays(workout.days)}</p>
        </div>
        <div className="flex items-center gap-3">
          <span className="text-sm text-text-secondary">
            ~{workout.estimatedDurationMinutes} min
          </span>
          <svg
            className={`w-5 h-5 text-text-secondary transition-transform ${
              expanded ? "rotate-180" : ""
            }`}
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
          </svg>
        </div>
      </button>

      {expanded && (
        <div className="border-t border-border px-4 py-2">
          {workout.exercises.map((exercise, index) => (
            <div key={index} className="flex gap-3 py-2">
              <div
                className="w-1 rounded-full"
                style={{ backgroundColor: getMuscleGroupColor(exercise.muscleGroup) }}
              />
              <div className="flex-1">
                <p className="font-medium text-text-primary text-sm">{exercise.name}</p>
                <p className="text-xs text-text-secondary">
                  {exercise.sets}x{exercise.reps}
                  {exercise.suggestedWeight ? ` • ${exercise.suggestedWeight}kg` : ""}
                  {" • "}{exercise.restSeconds}s descanso
                </p>
                {exercise.notes && (
                  <p className="text-xs text-text-light italic mt-0.5">{exercise.notes}</p>
                )}
              </div>
            </div>
          ))}
        </div>
      )}
    </Card>
  );
}
