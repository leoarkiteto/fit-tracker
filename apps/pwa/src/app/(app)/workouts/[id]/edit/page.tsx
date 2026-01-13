"use client";

import { useState, useEffect } from "react";
import { useRouter, useParams } from "next/navigation";
import { useAuth } from "@/lib/auth-context";
import { workoutsApi } from "@fittracker/api-client";
import { WORKOUT_GOALS, DAYS_OF_WEEK } from "@fittracker/types";
import type { Workout, Exercise, DayOfWeek, WorkoutGoal } from "@fittracker/types";
import { Card, CardContent, Button, Input } from "@/components/ui";
import Link from "next/link";
import { ArrowLeft, Plus, Trash2, Save, Loader2 } from "lucide-react";

interface ExerciseForm {
  id: string;
  name: string;
  sets: string;
  reps: string;
  weight: string;
  restSeconds: string;
  notes: string;
}

export default function EditWorkoutPage() {
  const { id } = useParams<{ id: string }>();
  const router = useRouter();
  const { user } = useAuth();
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState("");

  const [name, setName] = useState("");
  const [goal, setGoal] = useState<WorkoutGoal>("hypertrophy");
  const [selectedDays, setSelectedDays] = useState<DayOfWeek[]>([]);
  const [exercises, setExercises] = useState<ExerciseForm[]>([]);

  useEffect(() => {
    const loadWorkout = async () => {
      if (!id || !user?.profileId) {
        setLoading(false);
        return;
      }

      try {
        const workout = await workoutsApi.getById(user.profileId, id);
        setName(workout.name);
        setGoal(workout.goal);
        setSelectedDays(workout.days);
        setExercises(
          workout.exercises.map((ex) => ({
            id: ex.id || `temp-${Date.now()}-${Math.random()}`,
            name: ex.name,
            sets: ex.sets.toString(),
            reps: ex.reps.toString(),
            weight: ex.weight?.toString() || "",
            restSeconds: ex.restSeconds?.toString() || "60",
            notes: ex.notes || "",
          }))
        );
      } catch (error) {
        console.error("Error loading workout:", error);
        router.push("/workouts");
      } finally {
        setLoading(false);
      }
    };

    loadWorkout();
  }, [id, user?.profileId, router]);

  const toggleDay = (day: DayOfWeek) => {
    setSelectedDays((prev) =>
      prev.includes(day)
        ? prev.filter((d) => d !== day)
        : [...prev, day]
    );
  };

  const addExercise = () => {
    setExercises((prev) => [
      ...prev,
      {
        id: `temp-${Date.now()}`,
        name: "",
        sets: "3",
        reps: "12",
        weight: "",
        restSeconds: "60",
        notes: "",
      },
    ]);
  };

  const updateExercise = (id: string, field: keyof ExerciseForm, value: string) => {
    setExercises((prev) =>
      prev.map((ex) => (ex.id === id ? { ...ex, [field]: value } : ex))
    );
  };

  const removeExercise = (id: string) => {
    setExercises((prev) => prev.filter((ex) => ex.id !== id));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");

    if (!name.trim()) {
      setError("Workout name is required");
      return;
    }

    if (selectedDays.length === 0) {
      setError("Select at least one day");
      return;
    }

    if (exercises.length === 0) {
      setError("Add at least one exercise");
      return;
    }

    if (!user?.profileId || !id) {
      setError("Invalid workout");
      return;
    }

    setSaving(true);

    try {
      const exerciseList: Exercise[] = exercises.map((ex, index) => ({
        id: ex.id.startsWith("temp-") ? undefined : ex.id,
        name: ex.name,
        sets: parseInt(ex.sets) || 3,
        reps: parseInt(ex.reps) || 12,
        weight: ex.weight ? parseFloat(ex.weight) : undefined,
        restSeconds: ex.restSeconds ? parseInt(ex.restSeconds) : undefined,
        notes: ex.notes || undefined,
        order: index,
      }));

      await workoutsApi.update(user.profileId, id, {
        name,
        goal,
        days: selectedDays,
        exercises: exerciseList,
      });

      router.push(`/workouts/${id}`);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to update workout");
    } finally {
      setSaving(false);
    }
  };

  if (loading) {
    return (
      <div className="page-container flex items-center justify-center min-h-[60vh]">
        <Loader2 className="w-8 h-8 text-primary-500 animate-spin" />
      </div>
    );
  }

  return (
    <div className="page-container max-w-3xl mx-auto">
      {/* Header */}
      <div className="flex items-center gap-4 mb-6">
        <Link
          href={`/workouts/${id}`}
          className="p-2 rounded-xl bg-surface-200 hover:bg-surface-300 transition-colors"
        >
          <ArrowLeft className="w-5 h-5 text-surface-600" />
        </Link>
        <h1 className="text-2xl font-bold text-surface-900">Edit Workout</h1>
      </div>

      <form onSubmit={handleSubmit}>
        {/* Basic Info */}
        <Card className="mb-6">
          <CardContent className="p-6 space-y-4">
            <Input
              label="Workout Name"
              placeholder="e.g., Push Day, Leg Day..."
              value={name}
              onChange={(e) => setName(e.target.value)}
              required
            />

            <div>
              <label className="label">Goal</label>
              <div className="flex flex-wrap gap-2">
                {WORKOUT_GOALS.map((g) => (
                  <button
                    key={g.key}
                    type="button"
                    onClick={() => setGoal(g.key)}
                    className={`px-4 py-2 rounded-xl text-sm font-medium transition-colors ${
                      goal === g.key
                        ? "bg-primary-500 text-white"
                        : "bg-surface-200 text-surface-600 hover:bg-surface-300"
                    }`}
                  >
                    {g.label}
                  </button>
                ))}
              </div>
            </div>

            <div>
              <label className="label">Days</label>
              <div className="flex flex-wrap gap-2">
                {DAYS_OF_WEEK.map((day) => (
                  <button
                    key={day.key}
                    type="button"
                    onClick={() => toggleDay(day.key)}
                    className={`w-12 h-12 rounded-xl text-sm font-medium transition-colors ${
                      selectedDays.includes(day.key)
                        ? "bg-primary-500 text-white"
                        : "bg-surface-200 text-surface-600 hover:bg-surface-300"
                    }`}
                  >
                    {day.short}
                  </button>
                ))}
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Exercises */}
        <div className="mb-6">
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-lg font-semibold text-surface-900">
              Exercises ({exercises.length})
            </h2>
            <Button
              type="button"
              variant="secondary"
              size="sm"
              icon={<Plus className="w-4 h-4" />}
              onClick={addExercise}
            >
              Add
            </Button>
          </div>

          {exercises.length === 0 ? (
            <Card>
              <CardContent className="p-8 text-center">
                <p className="text-surface-600 mb-4">
                  No exercises added yet
                </p>
                <Button
                  type="button"
                  variant="outline"
                  icon={<Plus className="w-4 h-4" />}
                  onClick={addExercise}
                >
                  Add Exercise
                </Button>
              </CardContent>
            </Card>
          ) : (
            <div className="space-y-4">
              {exercises.map((exercise, index) => (
                <Card key={exercise.id} className="animate-slide-up">
                  <CardContent className="p-5">
                    <div className="flex items-start gap-4">
                      <div className="w-8 h-8 rounded-lg bg-primary-100 flex items-center justify-center flex-shrink-0">
                        <span className="text-primary-500 font-medium text-sm">
                          {index + 1}
                        </span>
                      </div>
                      <div className="flex-1 space-y-4">
                        <Input
                          placeholder="Exercise name"
                          value={exercise.name}
                          onChange={(e) =>
                            updateExercise(exercise.id, "name", e.target.value)
                          }
                          required
                        />
                        <div className="grid grid-cols-2 sm:grid-cols-4 gap-3">
                          <Input
                            label="Sets"
                            type="number"
                            value={exercise.sets}
                            onChange={(e) =>
                              updateExercise(exercise.id, "sets", e.target.value)
                            }
                          />
                          <Input
                            label="Reps"
                            type="number"
                            value={exercise.reps}
                            onChange={(e) =>
                              updateExercise(exercise.id, "reps", e.target.value)
                            }
                          />
                          <Input
                            label="Weight (kg)"
                            type="number"
                            step="0.5"
                            value={exercise.weight}
                            onChange={(e) =>
                              updateExercise(exercise.id, "weight", e.target.value)
                            }
                          />
                          <Input
                            label="Rest (s)"
                            type="number"
                            value={exercise.restSeconds}
                            onChange={(e) =>
                              updateExercise(
                                exercise.id,
                                "restSeconds",
                                e.target.value
                              )
                            }
                          />
                        </div>
                        <Input
                          placeholder="Notes (optional)"
                          value={exercise.notes}
                          onChange={(e) =>
                            updateExercise(exercise.id, "notes", e.target.value)
                          }
                        />
                      </div>
                      <button
                        type="button"
                        onClick={() => removeExercise(exercise.id)}
                        className="p-2 rounded-lg text-error hover:bg-red-50 transition-colors"
                      >
                        <Trash2 className="w-5 h-5" />
                      </button>
                    </div>
                  </CardContent>
                </Card>
              ))}
            </div>
          )}
        </div>

        {error && (
          <div className="p-3 rounded-xl bg-red-50 border border-red-200 mb-6">
            <p className="text-error text-sm text-center">{error}</p>
          </div>
        )}

        <div className="flex gap-3">
          <Link href={`/workouts/${id}`} className="flex-1">
            <Button type="button" variant="secondary" className="w-full">
              Cancel
            </Button>
          </Link>
          <Button
            type="submit"
            className="flex-1"
            loading={saving}
            icon={<Save className="w-4 h-4" />}
          >
            Save Changes
          </Button>
        </div>
      </form>
    </div>
  );
}
