import React from "react";
import {
  TouchableOpacity,
  Text,
  StyleSheet,
  ActivityIndicator,
  ViewStyle,
  TextStyle,
} from "react-native";
import { LinearGradient } from "expo-linear-gradient";
import { colors, spacing, borderRadius, typography } from "../theme";

interface ButtonProps {
  title: string;
  onPress: () => void;
  variant?: "primary" | "secondary" | "outline" | "ghost";
  size?: "sm" | "md" | "lg";
  disabled?: boolean;
  loading?: boolean;
  icon?: React.ReactNode;
  style?: ViewStyle;
  textStyle?: TextStyle;
  fullWidth?: boolean;
}

export const Button: React.FC<ButtonProps> = ({
  title,
  onPress,
  variant = "primary",
  size = "md",
  disabled = false,
  loading = false,
  icon,
  style,
  textStyle,
  fullWidth = false,
}) => {
  const getSize = () => {
    switch (size) {
      case "sm":
        return { paddingVertical: spacing.sm, paddingHorizontal: spacing.md };
      case "lg":
        return { paddingVertical: spacing.lg, paddingHorizontal: spacing.xl };
      default:
        return { paddingVertical: spacing.md, paddingHorizontal: spacing.lg };
    }
  };

  const getTextSize = () => {
    switch (size) {
      case "sm":
        return typography.sm;
      case "lg":
        return typography.lg;
      default:
        return typography.md;
    }
  };

  const isDisabled = disabled || loading;

  if (variant === "primary") {
    return (
      <TouchableOpacity
        onPress={onPress}
        disabled={isDisabled}
        style={[fullWidth && styles.fullWidth, style]}
        activeOpacity={0.8}
      >
        <LinearGradient
          colors={
            isDisabled
              ? [colors.textLight, colors.textLight]
              : [colors.primary, colors.primaryDark]
          }
          start={{ x: 0, y: 0 }}
          end={{ x: 1, y: 0 }}
          style={[styles.button, getSize(), styles.gradient]}
        >
          {loading ? (
            <ActivityIndicator color={colors.white} />
          ) : (
            <>
              {icon}
              <Text
                style={[
                  styles.textPrimary,
                  { fontSize: getTextSize() },
                  textStyle,
                ]}
              >
                {title}
              </Text>
            </>
          )}
        </LinearGradient>
      </TouchableOpacity>
    );
  }

  const buttonStyles = [
    styles.button,
    getSize(),
    variant === "secondary" && styles.secondary,
    variant === "outline" && styles.outline,
    variant === "ghost" && styles.ghost,
    isDisabled && styles.disabled,
    fullWidth && styles.fullWidth,
    style,
  ];

  const textStyles = [
    variant === "secondary" && styles.textSecondary,
    variant === "outline" && styles.textOutline,
    variant === "ghost" && styles.textGhost,
    { fontSize: getTextSize() },
    textStyle,
  ];

  return (
    <TouchableOpacity
      onPress={onPress}
      disabled={isDisabled}
      style={buttonStyles}
      activeOpacity={0.7}
    >
      {loading ? (
        <ActivityIndicator
          color={variant === "secondary" ? colors.white : colors.primary}
        />
      ) : (
        <>
          {icon}
          <Text style={textStyles}>{title}</Text>
        </>
      )}
    </TouchableOpacity>
  );
};

const styles = StyleSheet.create({
  button: {
    flexDirection: "row",
    alignItems: "center",
    justifyContent: "center",
    borderRadius: borderRadius.lg,
    gap: spacing.sm,
  },
  gradient: {
    borderRadius: borderRadius.lg,
  },
  secondary: {
    backgroundColor: colors.secondary,
  },
  outline: {
    backgroundColor: "transparent",
    borderWidth: 2,
    borderColor: colors.primary,
  },
  ghost: {
    backgroundColor: "transparent",
  },
  disabled: {
    opacity: 0.5,
  },
  fullWidth: {
    width: "100%",
  },
  textPrimary: {
    color: colors.white,
    fontWeight: typography.semibold,
  },
  textSecondary: {
    color: colors.white,
    fontWeight: typography.semibold,
  },
  textOutline: {
    color: colors.primary,
    fontWeight: typography.semibold,
  },
  textGhost: {
    color: colors.primary,
    fontWeight: typography.medium,
  },
});
