"use client";

import { useAuth } from "@/lib/auth-context";
import { useRouter } from "next/navigation";
import { useEffect } from "react";
import { Dumbbell, Loader2 } from "lucide-react";

export default function Home() {
  const { isAuthenticated, isLoading } = useAuth();
  const router = useRouter();

  useEffect(() => {
    if (!isLoading) {
      if (isAuthenticated) {
        router.replace("/dashboard");
      } else {
        router.replace("/auth/signin");
      }
    }
  }, [isAuthenticated, isLoading, router]);

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-primary-500 via-primary-600 to-purple-600">
      <div className="text-center animate-fade-in">
        <div className="flex items-center justify-center gap-3 mb-4">
          <div className="w-14 h-14 rounded-xl bg-white/20 backdrop-blur flex items-center justify-center">
            <Dumbbell className="w-8 h-8 text-white" />
          </div>
          <h1 className="text-4xl font-bold text-white">FitTracker</h1>
        </div>
        <Loader2 className="w-8 h-8 text-white animate-spin mx-auto" />
      </div>
    </div>
  );
}
