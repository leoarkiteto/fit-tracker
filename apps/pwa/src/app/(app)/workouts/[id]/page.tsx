"use client";

import { useState, useEffect } from "react";
import { useRouter, useParams } from "next/navigation";
import { useAuth } from "@/lib/auth-context";
import { workoutsApi, completedWorkoutsApi } from "@fittracker/api-client";
import type { Workout } from "@fittracker/types";
import { Card, CardContent, Button } from "@/components/ui";
import Link from "next/link";
import {
  ArrowLeft,
  Dumbbell,
  Play,
  Pencil,
  Trash2,
  Loader2,
  CheckCircle,
  Timer,
  Hash,
  Weight,
} from "lucide-react";

export default function WorkoutDetailPage() {
  const { id } = useParams<{ id: string }>();
  const router = useRouter();
  const { user } = useAuth();
  const [workout, setWorkout] = useState<Workout | null>(null);
  const [loading, setLoading] = useState(true);
  const [deleting, setDeleting] = useState(false);
  const [completing, setCompleting] = useState(false);
  const [isCompleted, setIsCompleted] = useState(false);

  useEffect(() => {
    const loadWorkout = async () => {
      if (!id || !user?.profileId) {
        setLoading(false);
        return;
      }

      try {
        const [workoutData, completedData] = await Promise.all([
          workoutsApi.getById(user.profileId, id),
          completedWorkoutsApi.getAll(user.profileId),
        ]);

        setWorkout(workoutData);

        // Check if completed today
        const today = new Date().toISOString().split("T")[0];
        const completedToday = completedData.some(
          (c) => c.workoutId === id && c.completedAt.startsWith(today)
        );
        setIsCompleted(completedToday);
      } catch (error) {
        console.error("Error loading workout:", error);
        router.push("/workouts");
      } finally {
        setLoading(false);
      }
    };

    loadWorkout();
  }, [id, user?.profileId, router]);

  const handleComplete = async () => {
    if (!workout || !user?.profileId || isCompleted) return;

    setCompleting(true);
    try {
      await completedWorkoutsApi.create(user.profileId, {
        workoutId: workout.id,
        completedAt: new Date().toISOString(),
        durationMinutes: 30, // Default duration
      });
      setIsCompleted(true);
    } catch (error) {
      console.error("Error completing workout:", error);
    } finally {
      setCompleting(false);
    }
  };

  const handleDelete = async () => {
    if (!workout || !user?.profileId) return;
    if (!confirm("Are you sure you want to delete this workout?")) return;

    setDeleting(true);
    try {
      await workoutsApi.delete(user.profileId, workout.id);
      router.push("/workouts");
    } catch (error) {
      console.error("Error deleting workout:", error);
      setDeleting(false);
    }
  };

  const getDaysLabel = (days: string[]) => {
    const dayNames: Record<string, string> = {
      monday: "Mon",
      tuesday: "Tue",
      wednesday: "Wed",
      thursday: "Thu",
      friday: "Fri",
      saturday: "Sat",
      sunday: "Sun",
    };
    return days.map((d) => dayNames[d] || d).join(", ");
  };

  if (loading) {
    return (
      <div className="page-container flex items-center justify-center min-h-[60vh]">
        <Loader2 className="w-8 h-8 text-primary-500 animate-spin" />
      </div>
    );
  }

  if (!workout) {
    return (
      <div className="page-container text-center">
        <p className="text-surface-600">Workout not found</p>
      </div>
    );
  }

  return (
    <div className="page-container max-w-3xl mx-auto">
      {/* Header */}
      <div className="flex items-center gap-4 mb-6">
        <Link
          href="/workouts"
          className="p-2 rounded-xl bg-surface-200 hover:bg-surface-300 transition-colors"
        >
          <ArrowLeft className="w-5 h-5 text-surface-600" />
        </Link>
        <div className="flex-1">
          <h1 className="text-2xl font-bold text-surface-900">{workout.name}</h1>
          <p className="text-surface-600">{getDaysLabel(workout.days)}</p>
        </div>
      </div>

      {/* Status Banner */}
      {isCompleted && (
        <div className="mb-6 p-4 bg-green-50 border border-green-200 rounded-xl flex items-center gap-3">
          <CheckCircle className="w-6 h-6 text-success" />
          <span className="text-success font-medium">
            Workout completed today!
          </span>
        </div>
      )}

      {/* Exercises */}
      <div className="mb-8">
        <h2 className="text-lg font-semibold text-surface-900 mb-4">
          Exercises ({workout.exercises.length})
        </h2>
        <div className="space-y-3">
          {workout.exercises.map((exercise, index) => (
            <Card
              key={exercise.id}
              className="animate-slide-up"
              style={{ animationDelay: `${index * 0.05}s` }}
            >
              <CardContent className="p-5">
                <div className="flex items-start gap-4">
                  <div className="w-10 h-10 rounded-xl bg-primary-100 flex items-center justify-center flex-shrink-0">
                    <span className="text-primary-500 font-semibold">
                      {index + 1}
                    </span>
                  </div>
                  <div className="flex-1">
                    <h3 className="font-semibold text-surface-900">
                      {exercise.name}
                    </h3>
                    <div className="flex flex-wrap gap-3 mt-2">
                      <div className="flex items-center gap-1 text-sm text-surface-600">
                        <Hash className="w-4 h-4" />
                        <span>{exercise.sets} sets</span>
                      </div>
                      <div className="flex items-center gap-1 text-sm text-surface-600">
                        <Dumbbell className="w-4 h-4" />
                        <span>{exercise.reps} reps</span>
                      </div>
                      {exercise.weight && (
                        <div className="flex items-center gap-1 text-sm text-surface-600">
                          <Weight className="w-4 h-4" />
                          <span>{exercise.weight} kg</span>
                        </div>
                      )}
                      {exercise.restSeconds && (
                        <div className="flex items-center gap-1 text-sm text-surface-600">
                          <Timer className="w-4 h-4" />
                          <span>{exercise.restSeconds}s rest</span>
                        </div>
                      )}
                    </div>
                    {exercise.notes && (
                      <p className="text-sm text-surface-500 mt-2">
                        {exercise.notes}
                      </p>
                    )}
                  </div>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      </div>

      {/* Actions */}
      <div className="space-y-3">
        {!isCompleted && (
          <Button
            className="w-full"
            size="lg"
            icon={<Play className="w-5 h-5" />}
            onClick={handleComplete}
            loading={completing}
          >
            Mark as Completed
          </Button>
        )}

        <div className="flex gap-3">
          <Link href={`/workouts/${workout.id}/edit`} className="flex-1">
            <Button
              variant="secondary"
              className="w-full"
              icon={<Pencil className="w-4 h-4" />}
            >
              Edit
            </Button>
          </Link>
          <Button
            variant="danger"
            className="flex-1"
            icon={<Trash2 className="w-4 h-4" />}
            onClick={handleDelete}
            loading={deleting}
          >
            Delete
          </Button>
        </div>
      </div>
    </div>
  );
}
