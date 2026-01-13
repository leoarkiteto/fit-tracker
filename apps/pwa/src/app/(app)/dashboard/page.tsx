"use client";

import { useState, useEffect } from "react";
import { useAuth } from "@/lib/auth-context";
import { workoutsApi, completedWorkoutsApi } from "@fittracker/api-client";
import type { Workout, WorkoutStats } from "@fittracker/types";
import { Card, CardContent, Button } from "@/components/ui";
import Link from "next/link";
import {
  Dumbbell,
  Clock,
  Calendar,
  TrendingUp,
  ChevronRight,
  Play,
  CheckCircle,
  Loader2,
  Plus,
} from "lucide-react";

export default function DashboardPage() {
  const { user } = useAuth();
  const [todayWorkouts, setTodayWorkouts] = useState<Workout[]>([]);
  const [stats, setStats] = useState<WorkoutStats | null>(null);
  const [loading, setLoading] = useState(true);
  const [completedToday, setCompletedToday] = useState<Set<string>>(new Set());

  useEffect(() => {
    const loadData = async () => {
      if (!user?.profileId) {
        setLoading(false);
        return;
      }

      try {
        const [workoutsData, statsData, completedData] = await Promise.all([
          workoutsApi.getToday(user.profileId),
          completedWorkoutsApi.getStats(user.profileId),
          completedWorkoutsApi.getAll(user.profileId),
        ]);

        setTodayWorkouts(workoutsData);
        setStats(statsData);

        // Check which workouts were completed today
        const today = new Date().toISOString().split("T")[0];
        const completedTodayIds = new Set(
          completedData
            .filter((c) => c.completedAt.startsWith(today))
            .map((c) => c.workoutId)
        );
        setCompletedToday(completedTodayIds);
      } catch (error) {
        console.error("Error loading dashboard data:", error);
      } finally {
        setLoading(false);
      }
    };

    loadData();
  }, [user?.profileId]);

  const formatMinutes = (minutes: number) => {
    const hours = Math.floor(minutes / 60);
    const mins = minutes % 60;
    if (hours > 0) {
      return `${hours}h ${mins}m`;
    }
    return `${mins}m`;
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
      {/* Welcome Section */}
      <div className="mb-8 animate-fade-in">
        <h1 className="text-3xl font-bold text-surface-900">
          Welcome back, <span className="gradient-text">{user?.name}</span>!
        </h1>
        <p className="text-surface-600 mt-1">Let&apos;s crush your goals today.</p>
      </div>

      {/* Stats Cards */}
      <div className="grid grid-cols-1 sm:grid-cols-3 gap-4 mb-8">
        <Card className="stat-card animate-slide-up" style={{ animationDelay: "0.1s" }}>
          <CardContent className="p-5">
            <div className="flex items-center gap-4">
              <div className="w-12 h-12 rounded-xl bg-primary-100 flex items-center justify-center">
                <Dumbbell className="w-6 h-6 text-primary-500" />
              </div>
              <div>
                <p className="text-surface-600 text-sm">Total Workouts</p>
                <p className="text-2xl font-bold text-surface-900">
                  {stats?.totalWorkoutsCompleted || 0}
                </p>
              </div>
            </div>
          </CardContent>
        </Card>

        <Card className="stat-card animate-slide-up" style={{ animationDelay: "0.2s" }}>
          <CardContent className="p-5">
            <div className="flex items-center gap-4">
              <div className="w-12 h-12 rounded-xl bg-secondary-100 flex items-center justify-center">
                <Calendar className="w-6 h-6 text-secondary-500" />
              </div>
              <div>
                <p className="text-surface-600 text-sm">This Week</p>
                <p className="text-2xl font-bold text-surface-900">
                  {stats?.workoutsThisWeek || 0}
                </p>
              </div>
            </div>
          </CardContent>
        </Card>

        <Card className="stat-card animate-slide-up" style={{ animationDelay: "0.3s" }}>
          <CardContent className="p-5">
            <div className="flex items-center gap-4">
              <div className="w-12 h-12 rounded-xl bg-accent-100 flex items-center justify-center">
                <Clock className="w-6 h-6 text-accent-400" />
              </div>
              <div>
                <p className="text-surface-600 text-sm">Time Spent</p>
                <p className="text-2xl font-bold text-surface-900">
                  {formatMinutes(stats?.totalMinutesSpent || 0)}
                </p>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Today's Workouts */}
      <div className="animate-slide-up" style={{ animationDelay: "0.4s" }}>
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-xl font-semibold text-surface-900">
            Today&apos;s Workouts
          </h2>
          <Link
            href="/workouts"
            className="text-primary-500 hover:text-primary-600 text-sm font-medium flex items-center gap-1"
          >
            View all
            <ChevronRight className="w-4 h-4" />
          </Link>
        </div>

        {todayWorkouts.length === 0 ? (
          <Card>
            <CardContent className="p-8 text-center">
              <div className="w-16 h-16 rounded-full bg-surface-200 flex items-center justify-center mx-auto mb-4">
                <Calendar className="w-8 h-8 text-surface-500" />
              </div>
              <h3 className="text-lg font-semibold text-surface-900 mb-2">
                No workouts scheduled
              </h3>
              <p className="text-surface-600 mb-4">
                You don&apos;t have any workouts for today. Create one to get started!
              </p>
              <Link href="/workouts/new">
                <Button icon={<Plus className="w-4 h-4" />}>
                  Create Workout
                </Button>
              </Link>
            </CardContent>
          </Card>
        ) : (
          <div className="space-y-3">
            {todayWorkouts.map((workout) => {
              const isCompleted = completedToday.has(workout.id);
              return (
                <Card
                  key={workout.id}
                  hover
                  className={isCompleted ? "opacity-75" : ""}
                >
                  <Link href={`/workouts/${workout.id}`}>
                    <CardContent className="p-5">
                      <div className="flex items-center justify-between">
                        <div className="flex items-center gap-4">
                          <div
                            className={`w-12 h-12 rounded-xl flex items-center justify-center ${
                              isCompleted
                                ? "bg-success/20"
                                : "bg-surface-200"
                            }`}
                          >
                            {isCompleted ? (
                              <CheckCircle className="w-6 h-6 text-success" />
                            ) : (
                              <Dumbbell className="w-6 h-6 text-surface-600" />
                            )}
                          </div>
                          <div>
                            <h3 className="font-semibold text-surface-900">
                              {workout.name}
                            </h3>
                            <p className="text-sm text-surface-600">
                              {workout.exercises.length} exercises
                            </p>
                          </div>
                        </div>
                        {!isCompleted && (
                          <Button size="sm" icon={<Play className="w-4 h-4" />}>
                            Start
                          </Button>
                        )}
                        {isCompleted && (
                          <span className="text-sm text-success font-medium">
                            Completed
                          </span>
                        )}
                      </div>
                    </CardContent>
                  </Link>
                </Card>
              );
            })}
          </div>
        )}
      </div>

      {/* Quick Actions */}
      <div className="mt-8 animate-slide-up" style={{ animationDelay: "0.5s" }}>
        <h2 className="text-xl font-semibold text-surface-900 mb-4">
          Quick Actions
        </h2>
        <div className="grid grid-cols-2 sm:grid-cols-4 gap-3">
          <Link href="/workouts/new">
            <Card hover className="h-full">
              <CardContent className="p-5 text-center">
                <div className="w-12 h-12 rounded-xl bg-primary-100 flex items-center justify-center mx-auto mb-2">
                  <Plus className="w-6 h-6 text-primary-500" />
                </div>
                <p className="text-sm font-medium text-surface-700">New Workout</p>
              </CardContent>
            </Card>
          </Link>
          <Link href="/bioimpedance/new">
            <Card hover className="h-full">
              <CardContent className="p-5 text-center">
                <div className="w-12 h-12 rounded-xl bg-secondary-100 flex items-center justify-center mx-auto mb-2">
                  <TrendingUp className="w-6 h-6 text-secondary-500" />
                </div>
                <p className="text-sm font-medium text-surface-700">Log Body</p>
              </CardContent>
            </Card>
          </Link>
          <Link href="/workouts">
            <Card hover className="h-full">
              <CardContent className="p-5 text-center">
                <div className="w-12 h-12 rounded-xl bg-accent-100 flex items-center justify-center mx-auto mb-2">
                  <Dumbbell className="w-6 h-6 text-accent-400" />
                </div>
                <p className="text-sm font-medium text-surface-700">My Workouts</p>
              </CardContent>
            </Card>
          </Link>
          <Link href="/profile">
            <Card hover className="h-full">
              <CardContent className="p-5 text-center">
                <div className="w-12 h-12 rounded-xl bg-purple-100 flex items-center justify-center mx-auto mb-2">
                  <Calendar className="w-6 h-6 text-purple-500" />
                </div>
                <p className="text-sm font-medium text-surface-700">Profile</p>
              </CardContent>
            </Card>
          </Link>
        </div>
      </div>
    </div>
  );
}
