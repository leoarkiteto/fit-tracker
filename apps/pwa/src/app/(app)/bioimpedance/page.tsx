"use client";

import { useState, useEffect } from "react";
import { useAuth } from "@/lib/auth-context";
import { bioimpedanceApi } from "@fittracker/api-client";
import type { BioimpedanceData } from "@fittracker/types";
import { Card, CardContent, Button } from "@/components/ui";
import Link from "next/link";
import {
  Activity,
  Plus,
  TrendingUp,
  TrendingDown,
  Loader2,
  Calendar,
  Scale,
  Droplet,
  Flame,
  ChevronRight,
} from "lucide-react";

export default function BioimpedancePage() {
  const { user } = useAuth();
  const [history, setHistory] = useState<BioimpedanceData[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const loadData = async () => {
      if (!user?.profileId) {
        setLoading(false);
        return;
      }

      try {
        const data = await bioimpedanceApi.getAll(user.profileId);
        // Sort by date descending
        setHistory(data.sort((a, b) => new Date(b.date).getTime() - new Date(a.date).getTime()));
      } catch (error) {
        console.error("Error loading bioimpedance:", error);
      } finally {
        setLoading(false);
      }
    };

    loadData();
  }, [user?.profileId]);

  const latest = history[0];
  const previous = history[1];

  const getTrend = (current: number, prev: number) => {
    const diff = current - prev;
    if (Math.abs(diff) < 0.1) return { icon: null, color: "text-surface-500" };
    return diff > 0
      ? { icon: TrendingUp, color: "text-success" }
      : { icon: TrendingDown, color: "text-error" };
  };

  const formatDate = (date: string) => {
    return new Date(date).toLocaleDateString("en-US", {
      month: "short",
      day: "numeric",
      year: "numeric",
    });
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
          <h1 className="text-3xl font-bold text-surface-900">Body Metrics</h1>
          <p className="text-surface-600 mt-1">Track your body composition</p>
        </div>
        <Link href="/bioimpedance/new">
          <Button icon={<Plus className="w-4 h-4" />}>New Entry</Button>
        </Link>
      </div>

      {history.length === 0 ? (
        <Card>
          <CardContent className="p-8 text-center">
            <div className="w-16 h-16 rounded-full bg-surface-200 flex items-center justify-center mx-auto mb-4">
              <Activity className="w-8 h-8 text-surface-500" />
            </div>
            <h3 className="text-lg font-semibold text-surface-900 mb-2">
              No measurements yet
            </h3>
            <p className="text-surface-600 mb-4">
              Start tracking your body composition to see your progress
            </p>
            <Link href="/bioimpedance/new">
              <Button icon={<Plus className="w-4 h-4" />}>
                Add First Entry
              </Button>
            </Link>
          </CardContent>
        </Card>
      ) : (
        <>
          {/* Latest Stats */}
          <div className="grid grid-cols-2 sm:grid-cols-4 gap-4 mb-8">
            <Card className="stat-card animate-slide-up" style={{ animationDelay: "0.1s" }}>
              <CardContent className="p-4">
                <div className="flex items-center gap-3 mb-2">
                  <Scale className="w-5 h-5 text-info" />
                  <span className="text-sm text-surface-600">Weight</span>
                </div>
                <p className="text-2xl font-bold text-surface-900">
                  {latest.weight} kg
                </p>
                {previous && (
                  <div className="flex items-center gap-1 mt-1">
                    {(() => {
                      const trend = getTrend(latest.weight, previous.weight);
                      return trend.icon ? (
                        <trend.icon className={`w-4 h-4 ${trend.color}`} />
                      ) : null;
                    })()}
                    <span className="text-xs text-surface-600">
                      {(latest.weight - previous.weight).toFixed(1)} kg
                    </span>
                  </div>
                )}
              </CardContent>
            </Card>

            <Card className="stat-card animate-slide-up" style={{ animationDelay: "0.2s" }}>
              <CardContent className="p-4">
                <div className="flex items-center gap-3 mb-2">
                  <Activity className="w-5 h-5 text-warning" />
                  <span className="text-sm text-surface-600">Body Fat</span>
                </div>
                <p className="text-2xl font-bold text-surface-900">
                  {latest.bodyFatPercentage}%
                </p>
                {previous && (
                  <div className="flex items-center gap-1 mt-1">
                    {(() => {
                      const trend = getTrend(previous.bodyFatPercentage, latest.bodyFatPercentage);
                      return trend.icon ? (
                        <trend.icon className={`w-4 h-4 ${trend.color}`} />
                      ) : null;
                    })()}
                    <span className="text-xs text-surface-600">
                      {(latest.bodyFatPercentage - previous.bodyFatPercentage).toFixed(1)}%
                    </span>
                  </div>
                )}
              </CardContent>
            </Card>

            <Card className="stat-card animate-slide-up" style={{ animationDelay: "0.3s" }}>
              <CardContent className="p-4">
                <div className="flex items-center gap-3 mb-2">
                  <Flame className="w-5 h-5 text-error" />
                  <span className="text-sm text-surface-600">Muscle</span>
                </div>
                <p className="text-2xl font-bold text-surface-900">
                  {latest.muscleMass} kg
                </p>
                {previous && (
                  <div className="flex items-center gap-1 mt-1">
                    {(() => {
                      const trend = getTrend(latest.muscleMass, previous.muscleMass);
                      return trend.icon ? (
                        <trend.icon className={`w-4 h-4 ${trend.color}`} />
                      ) : null;
                    })()}
                    <span className="text-xs text-surface-600">
                      {(latest.muscleMass - previous.muscleMass).toFixed(1)} kg
                    </span>
                  </div>
                )}
              </CardContent>
            </Card>

            <Card className="stat-card animate-slide-up" style={{ animationDelay: "0.4s" }}>
              <CardContent className="p-4">
                <div className="flex items-center gap-3 mb-2">
                  <Droplet className="w-5 h-5 text-secondary-500" />
                  <span className="text-sm text-surface-600">Water</span>
                </div>
                <p className="text-2xl font-bold text-surface-900">
                  {latest.waterPercentage}%
                </p>
              </CardContent>
            </Card>
          </div>

          {/* History */}
          <div>
            <h2 className="text-xl font-semibold text-surface-900 mb-4">History</h2>
            <div className="space-y-3">
              {history.map((entry, index) => (
                <Card
                  key={entry.id}
                  hover
                  className="animate-slide-up"
                  style={{ animationDelay: `${0.5 + index * 0.05}s` }}
                >
                  <CardContent className="p-4">
                    <div className="flex items-center justify-between">
                      <div className="flex items-center gap-4">
                        <div className="w-10 h-10 rounded-xl bg-primary-100 flex items-center justify-center">
                          <Calendar className="w-5 h-5 text-primary-500" />
                        </div>
                        <div>
                          <p className="font-medium text-surface-900">
                            {formatDate(entry.date)}
                          </p>
                          <p className="text-sm text-surface-600">
                            {entry.weight} kg · {entry.bodyFatPercentage}% fat ·{" "}
                            {entry.muscleMass} kg muscle
                          </p>
                        </div>
                      </div>
                      <ChevronRight className="w-5 h-5 text-surface-500" />
                    </div>
                  </CardContent>
                </Card>
              ))}
            </div>
          </div>
        </>
      )}
    </div>
  );
}
