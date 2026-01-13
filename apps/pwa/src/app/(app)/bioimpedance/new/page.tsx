"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { useAuth } from "@/lib/auth-context";
import { bioimpedanceApi } from "@fittracker/api-client";
import { Card, CardContent, Button, Input } from "@/components/ui";
import Link from "next/link";
import { ArrowLeft, Save } from "lucide-react";

export default function NewBioimpedancePage() {
  const router = useRouter();
  const { user } = useAuth();
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState("");

  const today = new Date().toISOString().split("T")[0];
  const [formData, setFormData] = useState({
    date: today,
    weight: "",
    bodyFatPercentage: "",
    muscleMass: "",
    boneMass: "",
    waterPercentage: "",
    visceralFat: "",
    bmr: "",
    metabolicAge: "",
    notes: "",
  });

  const updateField = (field: string, value: string) => {
    setFormData((prev) => ({ ...prev, [field]: value }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");

    if (!formData.weight) {
      setError("Weight is required");
      return;
    }

    if (!user?.profileId) {
      setError("Please complete your profile first");
      return;
    }

    setSaving(true);

    try {
      await bioimpedanceApi.create(user.profileId, {
        date: formData.date,
        weight: parseFloat(formData.weight) || 0,
        bodyFatPercentage: parseFloat(formData.bodyFatPercentage) || 0,
        muscleMass: parseFloat(formData.muscleMass) || 0,
        boneMass: parseFloat(formData.boneMass) || 0,
        waterPercentage: parseFloat(formData.waterPercentage) || 0,
        visceralFat: parseFloat(formData.visceralFat) || 0,
        bmr: parseFloat(formData.bmr) || 0,
        metabolicAge: parseFloat(formData.metabolicAge) || 0,
        notes: formData.notes || undefined,
      });

      router.push("/bioimpedance");
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to save entry");
    } finally {
      setSaving(false);
    }
  };

  return (
    <div className="page-container max-w-2xl mx-auto">
      {/* Header */}
      <div className="flex items-center gap-4 mb-6">
        <Link
          href="/bioimpedance"
          className="p-2 rounded-xl bg-surface-200 hover:bg-surface-300 transition-colors"
        >
          <ArrowLeft className="w-5 h-5 text-surface-600" />
        </Link>
        <h1 className="text-2xl font-bold text-surface-900">New Measurement</h1>
      </div>

      <form onSubmit={handleSubmit}>
        <Card className="mb-6">
          <CardContent className="p-6 space-y-4">
            <Input
              label="Date"
              type="date"
              value={formData.date}
              onChange={(e) => updateField("date", e.target.value)}
              required
            />

            <div className="grid grid-cols-2 gap-4">
              <Input
                label="Weight (kg)"
                type="number"
                step="0.1"
                placeholder="70.5"
                value={formData.weight}
                onChange={(e) => updateField("weight", e.target.value)}
                required
              />

              <Input
                label="Body Fat (%)"
                type="number"
                step="0.1"
                placeholder="15.0"
                value={formData.bodyFatPercentage}
                onChange={(e) => updateField("bodyFatPercentage", e.target.value)}
              />

              <Input
                label="Muscle Mass (kg)"
                type="number"
                step="0.1"
                placeholder="35.0"
                value={formData.muscleMass}
                onChange={(e) => updateField("muscleMass", e.target.value)}
              />

              <Input
                label="Bone Mass (kg)"
                type="number"
                step="0.1"
                placeholder="3.0"
                value={formData.boneMass}
                onChange={(e) => updateField("boneMass", e.target.value)}
              />

              <Input
                label="Water (%)"
                type="number"
                step="0.1"
                placeholder="55.0"
                value={formData.waterPercentage}
                onChange={(e) => updateField("waterPercentage", e.target.value)}
              />

              <Input
                label="Visceral Fat"
                type="number"
                placeholder="5"
                value={formData.visceralFat}
                onChange={(e) => updateField("visceralFat", e.target.value)}
              />

              <Input
                label="BMR (kcal)"
                type="number"
                placeholder="1700"
                value={formData.bmr}
                onChange={(e) => updateField("bmr", e.target.value)}
              />

              <Input
                label="Metabolic Age"
                type="number"
                placeholder="25"
                value={formData.metabolicAge}
                onChange={(e) => updateField("metabolicAge", e.target.value)}
              />
            </div>

            <div>
              <label className="label">Notes (optional)</label>
              <textarea
                className="input min-h-[100px] resize-none"
                placeholder="Any additional notes..."
                value={formData.notes}
                onChange={(e) => updateField("notes", e.target.value)}
              />
            </div>
          </CardContent>
        </Card>

        {error && (
          <div className="p-3 rounded-xl bg-red-50 border border-red-200 mb-6">
            <p className="text-error text-sm text-center">{error}</p>
          </div>
        )}

        <div className="flex gap-3">
          <Link href="/bioimpedance" className="flex-1">
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
            Save Entry
          </Button>
        </div>
      </form>
    </div>
  );
}
