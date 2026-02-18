import React from "react";
import {
  View,
  StyleSheet,
  StyleProp,
  ViewStyle,
  TouchableOpacity,
} from "react-native";
import { colors, spacing, borderRadius, shadows } from "../theme";

interface CardProps {
  children: React.ReactNode;
  style?: StyleProp<ViewStyle>;
  onPress?: () => void;
  variant?: "elevated" | "outlined" | "filled";
  padding?: "none" | "sm" | "md" | "lg";
}

export const Card: React.FC<CardProps> = ({
  children,
  style,
  onPress,
  variant = "elevated",
  padding = "md",
}) => {
  const getPadding = () => {
    switch (padding) {
      case "none":
        return 0;
      case "sm":
        return spacing.sm;
      case "lg":
        return spacing.lg;
      default:
        return spacing.md;
    }
  };

  const cardStyles = [
    styles.card,
    variant === "elevated" && styles.elevated,
    variant === "outlined" && styles.outlined,
    variant === "filled" && styles.filled,
    { padding: getPadding() },
    style,
  ];

  if (onPress) {
    return (
      <TouchableOpacity
        onPress={onPress}
        activeOpacity={0.7}
        style={cardStyles}
      >
        {children}
      </TouchableOpacity>
    );
  }

  return <View style={cardStyles}>{children}</View>;
};

const styles = StyleSheet.create({
  card: {
    borderRadius: borderRadius.lg,
    backgroundColor: colors.surface,
  },
  elevated: {
    ...shadows.md,
  },
  outlined: {
    borderWidth: 1,
    borderColor: colors.border,
  },
  filled: {
    backgroundColor: colors.surfaceSecondary,
  },
});
