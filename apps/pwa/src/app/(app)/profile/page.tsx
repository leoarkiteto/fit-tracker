"use client";

import { useState, useEffect } from "react";
import { useAuth } from "@/lib/auth-context";
import { profileApi } from "@fittracker/api-client";
import type { UserProfile } from "@fittracker/types";
import { Card, CardContent, CardHeader, Button, Input } from "@/components/ui";
import Link from "next/link";
import {
  User,
  Mail,
  Key,
  Ruler,
  Scale,
  Target,
  Calendar,
  Save,
  Loader2,
  LogOut,
  ChevronRight,
} from "lucide-react";

export default function ProfilePage() {
  const { user, logout, updateUser } = useAuth();
  const [profile, setProfile] = useState<UserProfile | null>(null);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState("");
  const [success, setSuccess] = useState(false);

  const [formData, setFormData] = useState({
    name: "",
    age: "",
    height: "",
    currentWeight: "",
    goalWeight: "",
  });

  useEffect(() => {
    const loadProfile = async () => {
      if (!user?.profileId) {
        // Set form with user name if no profile exists
        setFormData((prev) => ({ ...prev, name: user?.name || "" }));
        setLoading(false);
        return;
      }

      try {
        const data = await profileApi.getById(user.profileId);
        setProfile(data);
        setFormData({
          name: data.name || "",
          age: data.age?.toString() || "",
          height: data.height?.toString() || "",
          currentWeight: data.currentWeight?.toString() || "",
          goalWeight: data.goalWeight?.toString() || "",
        });
      } catch (error) {
        console.error("Error loading profile:", error);
        // If profile doesn't exist, use user name
        setFormData((prev) => ({ ...prev, name: user?.name || "" }));
      } finally {
        setLoading(false);
      }
    };

    loadProfile();
  }, [user?.profileId, user?.name]);

  const updateField = (field: string, value: string) => {
    setFormData((prev) => ({ ...prev, [field]: value }));
    setSuccess(false);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    setSuccess(false);

    if (!formData.name.trim()) {
      setError("Name is required");
      return;
    }

    setSaving(true);

    try {
      const profileData: UserProfile = {
        name: formData.name,
        age: formData.age ? parseInt(formData.age) : undefined,
        height: formData.height ? parseInt(formData.height) : undefined,
        currentWeight: formData.currentWeight
          ? parseFloat(formData.currentWeight)
          : undefined,
        goalWeight: formData.goalWeight
          ? parseFloat(formData.goalWeight)
          : undefined,
      };

      let savedProfile;
      if (profile?.id) {
        savedProfile = await profileApi.update(profile.id, profileData);
      } else {
        savedProfile = await profileApi.create(profileData);
        // Update the user's profileId in auth context
        if (user) {
          updateUser({ ...user, profileId: savedProfile.id });
        }
      }

      setProfile(savedProfile);
      setSuccess(true);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to save profile");
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
    <div className="page-container max-w-2xl mx-auto">
      <h1 className="text-3xl font-bold text-surface-900 mb-6">Profile</h1>

      {/* Account Info */}
      <Card className="mb-6 animate-slide-up">
        <CardHeader>
          <h2 className="text-lg font-semibold text-surface-900">Account</h2>
        </CardHeader>
        <CardContent className="p-0">
          <div className="flex items-center gap-4 px-6 py-4 border-b border-surface-200">
            <Mail className="w-5 h-5 text-surface-500" />
            <div className="flex-1">
              <p className="text-sm text-surface-600">Email</p>
              <p className="text-surface-900">{user?.email}</p>
            </div>
          </div>
          <Link
            href="/auth/change-password"
            className="flex items-center gap-4 px-6 py-4 hover:bg-surface-100 transition-colors"
          >
            <Key className="w-5 h-5 text-surface-500" />
            <div className="flex-1">
              <p className="text-surface-900">Change Password</p>
              <p className="text-sm text-surface-600">Update your password</p>
            </div>
            <ChevronRight className="w-5 h-5 text-surface-500" />
          </Link>
        </CardContent>
      </Card>

      {/* Profile Form */}
      <Card className="mb-6 animate-slide-up" style={{ animationDelay: "0.1s" }}>
        <CardHeader>
          <h2 className="text-lg font-semibold text-surface-900">Personal Info</h2>
        </CardHeader>
        <CardContent className="p-6">
          <form onSubmit={handleSubmit} className="space-y-4">
            <Input
              label="Name"
              placeholder="Your name"
              value={formData.name}
              onChange={(e) => updateField("name", e.target.value)}
              icon={<User className="w-5 h-5" />}
              required
            />

            <div className="grid grid-cols-2 gap-4">
              <Input
                label="Age"
                type="number"
                placeholder="25"
                value={formData.age}
                onChange={(e) => updateField("age", e.target.value)}
                icon={<Calendar className="w-5 h-5" />}
              />

              <Input
                label="Height (cm)"
                type="number"
                placeholder="175"
                value={formData.height}
                onChange={(e) => updateField("height", e.target.value)}
                icon={<Ruler className="w-5 h-5" />}
              />

              <Input
                label="Current Weight (kg)"
                type="number"
                step="0.1"
                placeholder="70.0"
                value={formData.currentWeight}
                onChange={(e) => updateField("currentWeight", e.target.value)}
                icon={<Scale className="w-5 h-5" />}
              />

              <Input
                label="Goal Weight (kg)"
                type="number"
                step="0.1"
                placeholder="65.0"
                value={formData.goalWeight}
                onChange={(e) => updateField("goalWeight", e.target.value)}
                icon={<Target className="w-5 h-5" />}
              />
            </div>

            {error && (
              <div className="p-3 rounded-xl bg-red-50 border border-red-200">
                <p className="text-error text-sm text-center">{error}</p>
              </div>
            )}

            {success && (
              <div className="p-3 rounded-xl bg-green-50 border border-green-200">
                <p className="text-success text-sm text-center">
                  Profile saved successfully!
                </p>
              </div>
            )}

            <Button
              type="submit"
              className="w-full"
              loading={saving}
              icon={<Save className="w-4 h-4" />}
            >
              Save Changes
            </Button>
          </form>
        </CardContent>
      </Card>

      {/* Sign Out */}
      <Card className="animate-slide-up" style={{ animationDelay: "0.2s" }}>
        <CardContent className="p-4">
          <Button
            variant="danger"
            className="w-full"
            onClick={logout}
            icon={<LogOut className="w-4 h-4" />}
          >
            Sign Out
          </Button>
        </CardContent>
      </Card>
    </div>
  );
}
