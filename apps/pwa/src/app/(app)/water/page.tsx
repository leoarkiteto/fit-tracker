"use client";

import { useState, useEffect } from "react";
import { useAuth } from "@/lib/auth-context";
import { waterApi } from "../../../lib/api-client";
import type { DailyWaterSummary, WaterIntakeEntry } from "../../../lib/types";
import { Card, CardContent, Button, Input } from "@/components/ui";
import {
  Droplets,
  Loader2,
  Plus,
  Trash2,
} from "lucide-react";

export default function WaterPage() {
  const { user } = useAuth();
  const [summary, setSummary] = useState<DailyWaterSummary | null>(null);
  const [loading, setLoading] = useState(true);
  const [adding, setAdding] = useState(false);
  const [customMl, setCustomMl] = useState("");
  const [deletingId, setDeletingId] = useState<string | null>(null);

  const todayDateUtc = () => new Date().toISOString().slice(0, 10);

  const loadWater = async () => {
    if (!user?.profileId) return;
    try {
      const data = await waterApi.getDay(user.profileId, todayDateUtc());
      setSummary(data);
    } catch (e) {
      console.error("Error loading water:", e);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    if (!user?.profileId) {
      setLoading(false);
      return;
    }
    loadWater();
  }, [user?.profileId]);

  const addWater = async (amountMl: number) => {
    if (!user?.profileId || adding) return;
    setAdding(true);
    try {
      await waterApi.add(user.profileId, amountMl);
      setCustomMl("");
      await loadWater();
    } catch (e) {
      console.error("Error adding water:", e);
    } finally {
      setAdding(false);
    }
  };

  const removeEntry = async (entryId: string) => {
    if (!user?.profileId || deletingId) return;
    setDeletingId(entryId);
    try {
      await waterApi.deleteEntry(user.profileId, entryId);
      await loadWater();
    } catch (e) {
      console.error("Error removing entry:", e);
    } finally {
      setDeletingId(null);
    }
  };

  const formatTime = (iso: string) =>
    new Date(iso).toLocaleTimeString(undefined, {
      hour: "2-digit",
      minute: "2-digit",
    });

  if (loading) {
    return (
      <div className="page-container flex items-center justify-center min-h-[60vh]">
        <Loader2 className="w-8 h-8 text-primary-500 animate-spin" />
      </div>
    );
  }

  const totalMl = summary?.totalMl ?? 0;
  const goalMl = summary?.goalMl ?? 2000;
  const progressPct = goalMl > 0 ? Math.min(100, (totalMl / goalMl) * 100) : 0;

  return (
    <div className="page-container max-w-2xl mx-auto">
      <div className="mb-8 animate-fade-in">
        <h1 className="text-3xl font-bold text-surface-900">
          Consumo de água
        </h1>
        <p className="text-surface-600 mt-1">
          Acompanhe sua meta diária e registre o que bebeu.
        </p>
      </div>

      {/* Summary card */}
      <Card className="mb-8 animate-slide-up">
        <CardContent className="p-6">
          <div className="flex items-center gap-4 mb-4">
            <div className="w-14 h-14 rounded-xl bg-cyan-100 flex items-center justify-center">
              <Droplets className="w-8 h-8 text-cyan-600" />
            </div>
            <div className="flex-1 min-w-0">
              <p className="text-surface-600 text-sm">Hoje</p>
              <p className="text-2xl font-bold text-surface-900">
                {(totalMl / 1000).toFixed(1)} L / {(goalMl / 1000).toFixed(1)} L
              </p>
            </div>
          </div>
          <div className="w-full h-3 bg-surface-200 rounded-full overflow-hidden mb-6">
            <div
              className="h-full bg-cyan-500 rounded-full transition-all duration-300"
              style={{ width: `${progressPct}%` }}
            />
          </div>
          <div className="flex gap-2">
            <Button
              size="sm"
              variant="secondary"
              onClick={() => addWater(250)}
              disabled={adding}
              className="flex-1"
            >
              +250 ml
            </Button>
            <Button
              size="sm"
              variant="secondary"
              onClick={() => addWater(500)}
              disabled={adding}
              className="flex-1"
            >
              +500 ml
            </Button>
          </div>
        </CardContent>
      </Card>

      {/* Add custom amount */}
      <div className="mb-8 animate-slide-up" style={{ animationDelay: "0.1s" }}>
        <h2 className="text-lg font-semibold text-surface-900 mb-3">
          Adicionar quantidade
        </h2>
        <div className="flex gap-2">
          <Input
            type="number"
            placeholder="ml (ex: 300)"
            value={customMl}
            onChange={(e) => setCustomMl(e.target.value)}
            min={1}
            max={2000}
            className="flex-1"
          />
          <Button
            onClick={() => {
              const n = parseInt(customMl, 10);
              if (!isNaN(n) && n > 0 && n <= 2000) addWater(n);
            }}
            disabled={adding || !customMl.trim()}
            icon={<Plus className="w-4 h-4" />}
          >
            Adicionar
          </Button>
        </div>
      </div>

      {/* Entries list */}
      <div className="animate-slide-up" style={{ animationDelay: "0.2s" }}>
        <h2 className="text-lg font-semibold text-surface-900 mb-3">
          Registros de hoje
        </h2>
        {!summary?.entries?.length ? (
          <Card>
            <CardContent className="p-8 text-center">
              <Droplets className="w-12 h-12 text-surface-300 mx-auto mb-3" />
              <p className="text-surface-600">
                Nenhum registro ainda. Use os botões acima para adicionar.
              </p>
            </CardContent>
          </Card>
        ) : (
          <div className="space-y-2">
            {(summary.entries as WaterIntakeEntry[])
              .slice()
              .reverse()
              .map((entry) => (
                <Card key={entry.id} className="flex items-center">
                  <CardContent className="p-4 flex items-center justify-between w-full">
                    <div className="flex items-center gap-4">
                      <span className="text-surface-900 font-medium">
                        {entry.amountMl} ml
                      </span>
                      <span className="text-surface-500 text-sm">
                        {formatTime(entry.consumedAt)}
                      </span>
                    </div>
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={() => removeEntry(entry.id)}
                      disabled={deletingId === entry.id}
                      className="text-surface-500 hover:text-error"
                    >
                      {deletingId === entry.id ? (
                        <Loader2 className="w-4 h-4 animate-spin" />
                      ) : (
                        <Trash2 className="w-4 h-4" />
                      )}
                    </Button>
                  </CardContent>
                </Card>
              ))}
          </div>
        )}
      </div>
    </div>
  );
}
