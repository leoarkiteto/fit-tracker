"use client";

import React, {
  createContext,
  useContext,
  useState,
  useEffect,
  useCallback,
} from "react";
import {
  authApi,
  setAuthToken,
  setApiBaseUrl,
} from "@fittracker/api-client";
import type { AuthUser } from "@fittracker/types";

const TOKEN_KEY = "fittracker_token";
const USER_KEY = "fittracker_user";

interface AuthContextType {
  user: AuthUser | null;
  isLoading: boolean;
  isAuthenticated: boolean;
  login: (email: string, password: string) => Promise<void>;
  register: (email: string, password: string, name: string) => Promise<void>;
  logout: () => void;
  updateUser: (user: AuthUser) => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<AuthUser | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    // Set API base URL for browser
    if (typeof window !== "undefined") {
      const apiUrl =
        process.env.NEXT_PUBLIC_API_URL || "http://localhost:5000";
      setApiBaseUrl(apiUrl);
    }

    // Check for existing session
    const loadSession = async () => {
      try {
        const storedToken = localStorage.getItem(TOKEN_KEY);
        const storedUser = localStorage.getItem(USER_KEY);

        if (storedToken && storedUser) {
          setAuthToken(storedToken);
          const parsedUser = JSON.parse(storedUser);
          setUser(parsedUser);

          // Verify token is still valid
          try {
            const currentUser = await authApi.getCurrentUser();
            setUser(currentUser);
            localStorage.setItem(USER_KEY, JSON.stringify(currentUser));
          } catch {
            // Token invalid, clear session
            localStorage.removeItem(TOKEN_KEY);
            localStorage.removeItem(USER_KEY);
            setAuthToken(null);
            setUser(null);
          }
        }
      } catch (error) {
        console.error("Error loading session:", error);
      } finally {
        setIsLoading(false);
      }
    };

    loadSession();
  }, []);

  const login = useCallback(async (email: string, password: string) => {
    const result = await authApi.login(email, password);
    localStorage.setItem(TOKEN_KEY, result.token);
    localStorage.setItem(USER_KEY, JSON.stringify(result.user));
    setAuthToken(result.token);
    setUser(result.user);
  }, []);

  const register = useCallback(
    async (email: string, password: string, name: string) => {
      const result = await authApi.register(email, password, name);
      localStorage.setItem(TOKEN_KEY, result.token);
      localStorage.setItem(USER_KEY, JSON.stringify(result.user));
      setAuthToken(result.token);
      setUser(result.user);
    },
    []
  );

  const logout = useCallback(() => {
    localStorage.removeItem(TOKEN_KEY);
    localStorage.removeItem(USER_KEY);
    setAuthToken(null);
    setUser(null);
  }, []);

  const updateUser = useCallback((updatedUser: AuthUser) => {
    localStorage.setItem(USER_KEY, JSON.stringify(updatedUser));
    setUser(updatedUser);
  }, []);

  return (
    <AuthContext.Provider
      value={{
        user,
        isLoading,
        isAuthenticated: !!user,
        login,
        register,
        logout,
        updateUser,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
}
