"use client";

import { useState, useEffect } from "react";
import { useAuth } from "@/lib/auth-context";
import { workoutsApi } from "@fittracker/api-client";
import type { Workout, WorkoutGoal } from "@fittracker/types";
import { WORKOUT_GOALS } from "@fittracker/types";
import { Card, CardContent, Button } from "@/components/ui";
import Link from "next/link";
import {
  Dumbbell,
  Plus,
  ChevronRight,
  Loader2,
  Search,
} from "lucide-react";

export default function WorkoutsPage() {
  const { user } = useAuth();
  const [workouts, setWorkouts] = useState<Workout[]>([]);
  const [loading, setLoading] = useState(true);
  const [searchQuery, setSearchQuery] = useState("");
  const [selectedGoal, setSelectedGoal] = useState<WorkoutGoal | "all">("all");

  useEffect(() => {
    const loadWorkouts = async () => {
      if (!user?.profileId) {
        setLoading(false);
        return;
      }

      try {
        const data = await workoutsApi.getAll(user.profileId);
        setWorkouts(data);
      } catch (error) {
        console.error("Error loading workouts:", error);
      } finally {
        setLoading(false);
      }
    };

    loadWorkouts();
  }, [user?.profileId]);

  const filteredWorkouts = workouts.filter((workout) => {
    const matchesSearch = workout.name
      .toLowerCase()
      .includes(searchQuery.toLowerCase());
    const matchesGoal =
      selectedGoal === "all" || workout.goal === selectedGoal;
    return matchesSearch && matchesGoal;
  });

  const getGoalLabel = (goal: WorkoutGoal) => {
    return WORKOUT_GOALS.find((g) => g.key === goal)?.label || goal;
  };

  const getGoalColor = (goal: WorkoutGoal) => {
    const colors: Record<WorkoutGoal, string> = {
      hypertrophy: "bg-red-100 text-red-600",
      strength: "bg-amber-100 text-amber-600",
      endurance: "bg-blue-100 text-blue-600",
      weight_loss: "bg-green-100 text-green-600",
      maintenance: "bg-purple-100 text-purple-600",
    };
    return colors[goal] || "bg-surface-200 text-surface-600";
  };

  if (loading) {
    return (
      <div className="page-container flex items-center justify-center min-h-[60vh]">
        <Loader2 className="w-8 h-8 text-primary-500 animate-spin" />
      </div>
    );
  }

  return (
    <div className="page-container max-w-4xl mx-auto">
      {/* Header */}
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 mb-6">
        <div>
          <h1 className="text-3xl font-bold text-surface-900">My Workouts</h1>
          <p className="text-surface-600 mt-1">
            {workouts.length} workout{workouts.length !== 1 ? "s" : ""} total
          </p>
        </div>
        <Link href="/workouts/new">
          <Button icon={<Plus className="w-4 h-4" />}>New Workout</Button>
        </Link>
      </div>

      {/* Filters */}
      <div className="flex flex-col sm:flex-row gap-3 mb-6">
        <div className="relative flex-1">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-surface-500" />
          <input
            type="text"
            placeholder="Search workouts..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            className="input pl-10"
          />
        </div>
        <div className="flex gap-2 overflow-x-auto pb-2 sm:pb-0">
          <button
            onClick={() => setSelectedGoal("all")}
            className={`px-4 py-2 rounded-xl text-sm font-medium whitespace-nowrap transition-colors ${
              selectedGoal === "all"
                ? "bg-primary-500 text-white"
                : "bg-surface-200 text-surface-600 hover:bg-surface-300"
            }`}
          >
            All
          </button>
          {WORKOUT_GOALS.map((goal) => (
            <button
              key={goal.key}
              onClick={() => setSelectedGoal(goal.key)}
              className={`px-4 py-2 rounded-xl text-sm font-medium whitespace-nowrap transition-colors ${
                selectedGoal === goal.key
                  ? "bg-primary-500 text-white"
                  : "bg-surface-200 text-surface-600 hover:bg-surface-300"
              }`}
            >
              {goal.label}
            </button>
          ))}
        </div>
      </div>

      {/* Workouts List */}
      {filteredWorkouts.length === 0 ? (
        <Card>
          <CardContent className="p-8 text-center">
            <div className="w-16 h-16 rounded-full bg-surface-200 flex items-center justify-center mx-auto mb-4">
              <Dumbbell className="w-8 h-8 text-surface-500" />
            </div>
            <h3 className="text-lg font-semibold text-surface-900 mb-2">
              {workouts.length === 0
                ? "No workouts yet"
                : "No workouts found"}
            </h3>
            <p className="text-surface-600 mb-4">
              {workouts.length === 0
                ? "Create your first workout to get started!"
                : "Try adjusting your search or filters."}
            </p>
            {workouts.length === 0 && (
              <Link href="/workouts/new">
                <Button icon={<Plus className="w-4 h-4" />}>
                  Create Workout
                </Button>
              </Link>
            )}
          </CardContent>
        </Card>
      ) : (
        <div className="grid gap-4">
          {filteredWorkouts.map((workout, index) => (
            <Link key={workout.id} href={`/workouts/${workout.id}`}>
              <Card
                hover
                className="animate-slide-up"
                style={{ animationDelay: `${index * 0.05}s` }}
              >
                <CardContent className="p-5">
                  <div className="flex items-center justify-between">
                    <div className="flex items-center gap-4">
                      <div className="w-12 h-12 rounded-xl bg-primary-100 flex items-center justify-center">
                        <Dumbbell className="w-6 h-6 text-primary-500" />
                      </div>
                      <div>
                        <h3 className="font-semibold text-surface-900">
                          {workout.name}
                        </h3>
                        <div className="flex items-center gap-2 mt-1">
                          <span
                            className={`px-2 py-0.5 rounded-lg text-xs font-medium ${getGoalColor(
                              workout.goal
                            )}`}
                          >
                            {getGoalLabel(workout.goal)}
                          </span>
                          <span className="text-sm text-surface-600">
                            {workout.exercises.length} exercises
                          </span>
                        </div>
                      </div>
                    </div>
                    <ChevronRight className="w-5 h-5 text-surface-500" />
                  </div>
                </CardContent>
              </Card>
            </Link>
          ))}
        </div>
      )}
    </div>
  );
}
