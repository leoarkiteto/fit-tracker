import React, { useState } from "react";
import {
  View,
  Text,
  StyleSheet,
  TouchableOpacity,
  KeyboardAvoidingView,
  Platform,
  ScrollView,
  ActivityIndicator,
  Alert,
} from "react-native";
import { SafeAreaView } from "react-native-safe-area-context";
import { Ionicons } from "@expo/vector-icons";
import { useNavigation } from "@react-navigation/native";

import { Input, Button } from "../components";
import { colors, spacing, borderRadius, typography } from "../theme";
import { authApi } from "../services/api";

export const ChangePasswordScreen: React.FC = () => {
  const navigation = useNavigation();

  const [currentPassword, setCurrentPassword] = useState("");
  const [newPassword, setNewPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [showPasswords, setShowPasswords] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const [errors, setErrors] = useState<{
    currentPassword?: string;
    newPassword?: string;
    confirmPassword?: string;
  }>({});

  const validate = () => {
    const newErrors: {
      currentPassword?: string;
      newPassword?: string;
      confirmPassword?: string;
    } = {};

    if (!currentPassword) {
      newErrors.currentPassword = "Senha atual é obrigatória";
    }

    if (!newPassword) {
      newErrors.newPassword = "Nova senha é obrigatória";
    } else if (newPassword.length < 6) {
      newErrors.newPassword = "A senha deve ter pelo menos 6 caracteres";
    }

    if (newPassword !== confirmPassword) {
      newErrors.confirmPassword = "As senhas não coincidem";
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleChangePassword = async () => {
    if (!validate()) return;

    setIsLoading(true);
    try {
      await authApi.changePassword(currentPassword, newPassword);
      Alert.alert("Sucesso", "Senha alterada com sucesso!", [
        { text: "OK", onPress: () => navigation.goBack() },
      ]);
    } catch (error) {
      const message =
        error instanceof Error ? error.message : "Erro ao alterar senha";
      if (message.includes("400")) {
        Alert.alert("Erro", "Senha atual incorreta");
      } else {
        Alert.alert("Erro", "Não foi possível alterar a senha. Tente novamente.");
      }
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <SafeAreaView style={styles.container} edges={["top"]}>
      <KeyboardAvoidingView
        behavior={Platform.OS === "ios" ? "padding" : "height"}
        style={styles.keyboardView}
      >
        {/* Header */}
        <View style={styles.header}>
          <TouchableOpacity
            onPress={() => navigation.goBack()}
            style={styles.backButton}
          >
            <Ionicons name="arrow-back" size={24} color={colors.textPrimary} />
          </TouchableOpacity>
          <Text style={styles.headerTitle}>Alterar Senha</Text>
          <View style={styles.headerPlaceholder} />
        </View>

        <ScrollView
          contentContainerStyle={styles.scrollContent}
          showsVerticalScrollIndicator={false}
        >
          {/* Info */}
          <View style={styles.infoContainer}>
            <View style={styles.infoIconContainer}>
              <Ionicons name="lock-closed" size={32} color={colors.primary} />
            </View>
            <Text style={styles.infoText}>
              Para alterar sua senha, insira a senha atual e a nova senha desejada.
            </Text>
          </View>

          {/* Form */}
          <View style={styles.form}>
            <Input
              label="Senha Atual"
              placeholder="Digite sua senha atual"
              value={currentPassword}
              onChangeText={setCurrentPassword}
              secureTextEntry={!showPasswords}
              autoCapitalize="none"
              autoComplete="off"
              textContentType="oneTimeCode"
              error={errors.currentPassword}
              leftIcon="key-outline"
            />

            <Input
              label="Nova Senha"
              placeholder="Mínimo 6 caracteres"
              value={newPassword}
              onChangeText={setNewPassword}
              secureTextEntry={!showPasswords}
              autoCapitalize="none"
              autoComplete="off"
              textContentType="oneTimeCode"
              error={errors.newPassword}
              leftIcon="lock-closed-outline"
            />

            <Input
              label="Confirmar Nova Senha"
              placeholder="Repita a nova senha"
              value={confirmPassword}
              onChangeText={setConfirmPassword}
              secureTextEntry={!showPasswords}
              autoCapitalize="none"
              autoComplete="off"
              textContentType="oneTimeCode"
              error={errors.confirmPassword}
              leftIcon="lock-closed-outline"
            />

            <TouchableOpacity
              onPress={() => setShowPasswords(!showPasswords)}
              style={styles.showPasswordButton}
            >
              <Ionicons
                name={showPasswords ? "eye-off-outline" : "eye-outline"}
                size={20}
                color={colors.primary}
              />
              <Text style={styles.showPasswordText}>
                {showPasswords ? "Ocultar senhas" : "Mostrar senhas"}
              </Text>
            </TouchableOpacity>

            <Button
              title={isLoading ? "" : "Alterar Senha"}
              onPress={handleChangePassword}
              disabled={isLoading}
              fullWidth
              icon={
                isLoading ? (
                  <ActivityIndicator color={colors.white} />
                ) : (
                  <Ionicons name="checkmark-circle" size={20} color={colors.white} />
                )
              }
            />
          </View>
        </ScrollView>
      </KeyboardAvoidingView>
    </SafeAreaView>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: colors.background,
  },
  keyboardView: {
    flex: 1,
  },
  header: {
    flexDirection: "row",
    alignItems: "center",
    justifyContent: "space-between",
    paddingHorizontal: spacing.md,
    paddingVertical: spacing.md,
    borderBottomWidth: 1,
    borderBottomColor: colors.border,
  },
  backButton: {
    padding: spacing.sm,
  },
  headerTitle: {
    fontSize: typography.lg,
    fontWeight: typography.bold,
    color: colors.textPrimary,
  },
  headerPlaceholder: {
    width: 40,
  },
  scrollContent: {
    padding: spacing.lg,
  },
  infoContainer: {
    alignItems: "center",
    marginBottom: spacing.xl,
  },
  infoIconContainer: {
    width: 64,
    height: 64,
    borderRadius: borderRadius.full,
    backgroundColor: `${colors.primary}15`,
    alignItems: "center",
    justifyContent: "center",
    marginBottom: spacing.md,
  },
  infoText: {
    fontSize: typography.md,
    color: colors.textSecondary,
    textAlign: "center",
    lineHeight: 22,
  },
  form: {
    gap: spacing.md,
  },
  showPasswordButton: {
    flexDirection: "row",
    alignItems: "center",
    justifyContent: "center",
    gap: spacing.xs,
    paddingVertical: spacing.sm,
  },
  showPasswordText: {
    fontSize: typography.sm,
    color: colors.primary,
    fontWeight: typography.medium,
  },
});
