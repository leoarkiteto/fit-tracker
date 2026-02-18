import React, {
  createContext,
  useContext,
  useState,
  useEffect,
  ReactNode,
  useCallback,
} from "react";
import AsyncStorage from "@react-native-async-storage/async-storage";
import { authApi, setAuthToken, AuthUser, AuthResult } from "../services/api";

const TOKEN_KEY = "@fittracker:token";
const USER_KEY = "@fittracker:user";

interface AuthContextType {
  user: AuthUser | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  signIn: (email: string, password: string) => Promise<void>;
  signUp: (email: string, password: string, name: string) => Promise<void>;
  signOut: () => Promise<void>;
  refreshAuth: () => Promise<void>;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider: React.FC<{ children: ReactNode }> = ({
  children,
}) => {
  const [user, setUser] = useState<AuthUser | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  const isAuthenticated = user !== null;

  // Carregar token e usuário salvos
  useEffect(() => {
    loadStoredAuth();
  }, []);

  const loadStoredAuth = async () => {
    try {
      const [storedToken, storedUser] = await Promise.all([
        AsyncStorage.getItem(TOKEN_KEY),
        AsyncStorage.getItem(USER_KEY),
      ]);

      if (storedToken && storedUser) {
        setAuthToken(storedToken);
        setUser(JSON.parse(storedUser));

        // Tentar renovar o token
        try {
          const result = await authApi.refreshToken();
          await saveAuth(result);
        } catch {
          // Token expirado, fazer logout
          await clearAuth();
        }
      }
    } catch (error) {
      console.error("Error loading stored auth:", error);
    } finally {
      setIsLoading(false);
    }
  };

  const saveAuth = async (result: AuthResult) => {
    setAuthToken(result.token);
    setUser(result.user);
    await Promise.all([
      AsyncStorage.setItem(TOKEN_KEY, result.token),
      AsyncStorage.setItem(USER_KEY, JSON.stringify(result.user)),
    ]);
  };

  const clearAuth = async () => {
    setAuthToken(null);
    setUser(null);
    await Promise.all([
      AsyncStorage.removeItem(TOKEN_KEY),
      AsyncStorage.removeItem(USER_KEY),
    ]);
  };

  const signIn = useCallback(async (email: string, password: string) => {
    const result = await authApi.login(email, password);
    await saveAuth(result);
  }, []);

  const signUp = useCallback(
    async (email: string, password: string, name: string) => {
      const result = await authApi.register(email, password, name);
      await saveAuth(result);
    },
    []
  );

  const signOut = useCallback(async () => {
    await clearAuth();
  }, []);

  const refreshAuth = useCallback(async () => {
    try {
      const result = await authApi.refreshToken();
      await saveAuth(result);
    } catch {
      await clearAuth();
    }
  }, []);

  return (
    <AuthContext.Provider
      value={{
        user,
        isAuthenticated,
        isLoading,
        signIn,
        signUp,
        signOut,
        refreshAuth,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = (): AuthContextType => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
};
