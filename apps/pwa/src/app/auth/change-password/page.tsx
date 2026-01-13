"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { useAuth } from "@/lib/auth-context";
import { authApi } from "@fittracker/api-client";
import { Button, Input, Card, CardContent, CardHeader } from "@/components/ui";
import { Lock, ArrowLeft, CheckCircle } from "lucide-react";
import Link from "next/link";

export default function ChangePasswordPage() {
  const [currentPassword, setCurrentPassword] = useState("");
  const [newPassword, setNewPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [showPasswords, setShowPasswords] = useState(false);
  const [error, setError] = useState("");
  const [success, setSuccess] = useState(false);
  const [loading, setLoading] = useState(false);
  const { isAuthenticated } = useAuth();
  const router = useRouter();

  if (!isAuthenticated) {
    router.push("/auth/signin");
    return null;
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");

    if (newPassword !== confirmPassword) {
      setError("New passwords do not match");
      return;
    }

    if (newPassword.length < 6) {
      setError("New password must be at least 6 characters");
      return;
    }

    setLoading(true);

    try {
      await authApi.changePassword(currentPassword, newPassword);
      setSuccess(true);
      setCurrentPassword("");
      setNewPassword("");
      setConfirmPassword("");
    } catch (err) {
      setError(
        err instanceof Error ? err.message : "Failed to change password"
      );
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen px-4 py-8 bg-surface-100">
      <div className="max-w-md mx-auto animate-fade-in">
        <Link
          href="/profile"
          className="inline-flex items-center gap-2 text-surface-600 hover:text-surface-900 mb-6"
        >
          <ArrowLeft className="w-5 h-5" />
          Back to Profile
        </Link>

        <Card>
          <CardHeader>
            <h1 className="text-2xl font-bold text-surface-900">Change Password</h1>
            <p className="text-surface-600 mt-1">
              Update your account password
            </p>
          </CardHeader>

          <CardContent className="p-6">
            {success ? (
              <div className="text-center py-8">
                <div className="w-16 h-16 rounded-full bg-green-100 flex items-center justify-center mx-auto mb-4">
                  <CheckCircle className="w-8 h-8 text-success" />
                </div>
                <h2 className="text-xl font-semibold text-surface-900 mb-2">
                  Password Updated!
                </h2>
                <p className="text-surface-600 mb-6">
                  Your password has been changed successfully.
                </p>
                <Button onClick={() => router.push("/profile")}>
                  Back to Profile
                </Button>
              </div>
            ) : (
              <form onSubmit={handleSubmit} className="space-y-5">
                <div className="relative">
                  <Input
                    type={showPasswords ? "text" : "password"}
                    label="Current Password"
                    placeholder="••••••••"
                    value={currentPassword}
                    onChange={(e) => setCurrentPassword(e.target.value)}
                    icon={<Lock className="w-5 h-5" />}
                    required
                    autoComplete="current-password"
                  />
                </div>

                <div className="relative">
                  <Input
                    type={showPasswords ? "text" : "password"}
                    label="New Password"
                    placeholder="••••••••"
                    value={newPassword}
                    onChange={(e) => setNewPassword(e.target.value)}
                    icon={<Lock className="w-5 h-5" />}
                    required
                    autoComplete="new-password"
                  />
                </div>

                <div className="relative">
                  <Input
                    type={showPasswords ? "text" : "password"}
                    label="Confirm New Password"
                    placeholder="••••••••"
                    value={confirmPassword}
                    onChange={(e) => setConfirmPassword(e.target.value)}
                    icon={<Lock className="w-5 h-5" />}
                    required
                    autoComplete="new-password"
                  />
                </div>

                <label className="flex items-center gap-2 text-surface-600 cursor-pointer">
                  <input
                    type="checkbox"
                    checked={showPasswords}
                    onChange={(e) => setShowPasswords(e.target.checked)}
                    className="rounded border-surface-300"
                  />
                  <span className="text-sm">Show passwords</span>
                </label>

                {error && (
                  <div className="p-3 rounded-xl bg-red-50 border border-red-200">
                    <p className="text-error text-sm text-center">{error}</p>
                  </div>
                )}

                <div className="flex gap-3 pt-2">
                  <Button
                    type="button"
                    variant="secondary"
                    className="flex-1"
                    onClick={() => router.push("/profile")}
                  >
                    Cancel
                  </Button>
                  <Button
                    type="submit"
                    className="flex-1"
                    loading={loading}
                  >
                    Update Password
                  </Button>
                </div>
              </form>
            )}
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
