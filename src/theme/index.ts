export const colors = {
  // Cores primárias
  primary: "#6366F1", // Indigo vibrante
  primaryLight: "#818CF8",
  primaryDark: "#4F46E5",

  // Cores secundárias
  secondary: "#14B8A6", // Teal
  secondaryLight: "#5EEAD4",
  secondaryDark: "#0D9488",

  // Acentos
  accent: "#F59E0B", // Amber
  accentLight: "#FCD34D",

  // Cores de status
  success: "#10B981",
  warning: "#F59E0B",
  error: "#EF4444",
  info: "#3B82F6",

  // Cores neutras
  white: "#FFFFFF",
  background: "#F8FAFC",
  surface: "#FFFFFF",
  surfaceSecondary: "#F1F5F9",

  // Texto
  textPrimary: "#1E293B",
  textSecondary: "#64748B",
  textLight: "#94A3B8",
  textInverse: "#FFFFFF",

  // Bordas
  border: "#E2E8F0",
  borderLight: "#F1F5F9",

  // Gradientes
  gradientStart: "#6366F1",
  gradientEnd: "#8B5CF6",

  // Cards de treino
  cardYoga: "#E0F2FE",
  cardCardio: "#DCFCE7",
  cardStrength: "#FEF3C7",
  cardStretching: "#FCE7F3",
};

export const spacing = {
  xs: 4,
  sm: 8,
  md: 16,
  lg: 24,
  xl: 32,
  xxl: 48,
};

export const borderRadius = {
  sm: 8,
  md: 12,
  lg: 16,
  xl: 24,
  full: 9999,
};

export const typography = {
  // Font sizes
  xs: 12,
  sm: 14,
  md: 16,
  lg: 18,
  xl: 20,
  xxl: 24,
  xxxl: 32,

  // Font weights
  regular: "400" as const,
  medium: "500" as const,
  semibold: "600" as const,
  bold: "700" as const,
};

export const shadows = {
  sm: {
    shadowColor: "#000",
    shadowOffset: { width: 0, height: 1 },
    shadowOpacity: 0.05,
    shadowRadius: 2,
    elevation: 1,
  },
  md: {
    shadowColor: "#000",
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.08,
    shadowRadius: 4,
    elevation: 3,
  },
  lg: {
    shadowColor: "#000",
    shadowOffset: { width: 0, height: 4 },
    shadowOpacity: 0.1,
    shadowRadius: 8,
    elevation: 5,
  },
};

export default {
  colors,
  spacing,
  borderRadius,
  typography,
  shadows,
};
