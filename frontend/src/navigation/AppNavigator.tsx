import React from "react";
import { View, ActivityIndicator, StyleSheet } from "react-native";

import { useAuth } from "../context/AuthContext";
import { RootNavigator } from "./RootNavigator";
import { AuthNavigator } from "./AuthNavigator";
import { colors } from "../theme";

export const AppNavigator: React.FC = () => {
  const { isAuthenticated, isLoading } = useAuth();

  // Mostrar loading enquanto verifica autenticação
  if (isLoading) {
    return (
      <View style={styles.loadingContainer}>
        <ActivityIndicator size="large" color={colors.primary} />
      </View>
    );
  }

  // Mostrar telas de auth se não autenticado
  if (!isAuthenticated) {
    return <AuthNavigator />;
  }

  // Mostrar app principal se autenticado
  return <RootNavigator />;
};

const styles = StyleSheet.create({
  loadingContainer: {
    flex: 1,
    alignItems: "center",
    justifyContent: "center",
    backgroundColor: colors.background,
  },
});
