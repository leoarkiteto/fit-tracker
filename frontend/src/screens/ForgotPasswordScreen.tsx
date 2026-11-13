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
import { NativeStackNavigationProp } from "@react-navigation/native-stack";

import { Input, Button } from "../components";
import { colors, spacing, typography } from "../theme";
import { authApi } from "../services/api";
import { AuthStackParamList } from "../navigation/types";

type ForgotPasswordNavigationProp = NativeStackNavigationProp<
  AuthStackParamList,
  "ForgotPassword"
>;

export const ForgotPasswordScreen: React.FC = () => {
  const navigation = useNavigation<ForgotPasswordNavigationProp>();

  const [email, setEmail] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | undefined>();

  const validate = () => {
    let message: string | undefined;
    if (!email.trim()) {
      message = "Email é obrigatório";
    } else if (!/\S+@\S+\.\S+/.test(email)) {
      message = "Email inválido";
    }
    setError(message);
    return !message;
  };

  const handleSubmit = async () => {
    if (!validate()) return;

    const normalized = email.trim().toLowerCase();
    setIsLoading(true);
    try {
      await authApi.forgotPassword(normalized);
      navigation.navigate("ResetPassword", { email: normalized });
    } catch {
      Alert.alert("Erro", "Não foi possível enviar o código. Tente novamente.");
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <SafeAreaView style={styles.container}>
      <KeyboardAvoidingView
        behavior={Platform.OS === "ios" ? "padding" : "height"}
        style={styles.keyboardView}
      >
        <ScrollView
          contentContainerStyle={styles.scrollContent}
          showsVerticalScrollIndicator={false}
          keyboardShouldPersistTaps="handled"
        >
          <TouchableOpacity
            style={styles.backButton}
            onPress={() => navigation.goBack()}
            accessibilityRole="button"
            accessibilityLabel="Voltar"
          >
            <Ionicons name="arrow-back" size={24} color={colors.textPrimary} />
          </TouchableOpacity>

          <View style={styles.header}>
            <Text style={styles.title}>Esqueci minha senha</Text>
            <Text style={styles.subtitle}>
              Informe seu email para receber um código de recuperação
            </Text>
          </View>

          <View style={styles.form}>
            <Input
              label="Email"
              placeholder="seu@email.com"
              value={email}
              onChangeText={setEmail}
              keyboardType="email-address"
              autoCapitalize="none"
              autoComplete="email"
              textContentType="emailAddress"
              error={error}
              leftIcon="mail-outline"
            />

            <Button
              title={isLoading ? "" : "Enviar código"}
              onPress={handleSubmit}
              disabled={isLoading}
              fullWidth
              icon={
                isLoading ? (
                  <ActivityIndicator color={colors.white} />
                ) : (
                  <Ionicons name="send-outline" size={20} color={colors.white} />
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
  scrollContent: {
    flexGrow: 1,
    padding: spacing.lg,
    justifyContent: "center",
  },
  backButton: {
    position: "absolute",
    top: spacing.lg,
    left: spacing.lg,
    padding: spacing.xs,
  },
  header: {
    alignItems: "center",
    marginBottom: spacing.xl,
  },
  title: {
    fontSize: typography.xxxl,
    fontWeight: typography.bold,
    color: colors.textPrimary,
    marginBottom: spacing.xs,
    textAlign: "center",
  },
  subtitle: {
    fontSize: typography.md,
    color: colors.textSecondary,
    textAlign: "center",
  },
  form: {
    gap: spacing.md,
  },
});
