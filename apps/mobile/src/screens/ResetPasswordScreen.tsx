import React, { useCallback, useEffect, useRef, useState } from "react";
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
import { useNavigation, useRoute, RouteProp } from "@react-navigation/native";
import { NativeStackNavigationProp } from "@react-navigation/native-stack";

import { Input, Button } from "../components";
import { colors, spacing, typography } from "../theme";
import { authApi } from "../services/api";
import { AuthStackParamList } from "../navigation/types";

type ResetPasswordNavigationProp = NativeStackNavigationProp<
  AuthStackParamList,
  "ResetPassword"
>;

const RESEND_SECONDS = 60;

type FormErrors = {
  code?: string;
  password?: string;
  confirm?: string;
};

export const ResetPasswordScreen: React.FC = () => {
  const navigation = useNavigation<ResetPasswordNavigationProp>();
  const route = useRoute<RouteProp<AuthStackParamList, "ResetPassword">>();
  const { email } = route.params;

  const [code, setCode] = useState("");
  const [password, setPassword] = useState("");
  const [confirm, setConfirm] = useState("");
  const [showPassword, setShowPassword] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const [isResending, setIsResending] = useState(false);
  const [errors, setErrors] = useState<FormErrors>({});
  const [secondsLeft, setSecondsLeft] = useState(RESEND_SECONDS);

  const timerRef = useRef<ReturnType<typeof setInterval> | null>(null);

  const clearTimer = useCallback(() => {
    if (timerRef.current) {
      clearInterval(timerRef.current);
      timerRef.current = null;
    }
  }, []);

  const startCountdown = useCallback(() => {
    clearTimer();
    setSecondsLeft(RESEND_SECONDS);
    timerRef.current = setInterval(() => {
      setSecondsLeft((prev) => {
        if (prev <= 1) {
          clearTimer();
          return 0;
        }
        return prev - 1;
      });
    }, 1000);
  }, [clearTimer]);

  useEffect(() => {
    startCountdown();
    return clearTimer;
  }, [startCountdown, clearTimer]);

  const validate = () => {
    const newErrors: FormErrors = {};

    if (!/^\d{6}$/.test(code)) {
      newErrors.code = "O código deve ter 6 dígitos";
    }
    if (password.length < 6) {
      newErrors.password = "A senha deve ter pelo menos 6 caracteres";
    }
    if (confirm !== password) {
      newErrors.confirm = "As senhas não coincidem";
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleResend = async () => {
    if (secondsLeft > 0 || isResending) return;

    setIsResending(true);
    try {
      await authApi.forgotPassword(email);
      startCountdown();
    } catch {
      Alert.alert("Erro", "Não foi possível reenviar o código. Tente novamente.");
    } finally {
      setIsResending(false);
    }
  };

  const handleReset = async () => {
    if (!validate()) return;

    setIsLoading(true);
    try {
      await authApi.resetPassword(email, code, password);
      Alert.alert("Sucesso", "Senha redefinida. Entre com a nova senha.", [
        {
          text: "OK",
          onPress: () => navigation.popTo("SignIn", { email }),
        },
      ]);
    } catch (error) {
      const message = error instanceof Error ? error.message : "";
      if (message.includes("Código inválido ou expirado")) {
        Alert.alert("Erro", "Código inválido ou expirado");
      } else if (message.includes("A senha deve ter pelo menos 6 caracteres")) {
        Alert.alert("Erro", "A senha deve ter pelo menos 6 caracteres");
      } else {
        Alert.alert("Erro", "Não foi possível redefinir a senha. Tente novamente.");
      }
    } finally {
      setIsLoading(false);
    }
  };

  const resendDisabled = secondsLeft > 0 || isResending;

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
            <Text style={styles.title}>Redefinir senha</Text>
            <Text style={styles.subtitle}>
              Se houver uma conta com este e-mail, enviamos um código
            </Text>
          </View>

          <View style={styles.form}>
            <Input
              label="Email"
              value={email}
              editable={false}
              leftIcon="mail-outline"
            />

            <Input
              label="Código"
              placeholder="000000"
              value={code}
              onChangeText={(text) => setCode(text.replace(/\D/g, "").slice(0, 6))}
              keyboardType="number-pad"
              maxLength={6}
              autoComplete="off"
              textContentType="oneTimeCode"
              error={errors.code}
              leftIcon="keypad-outline"
            />

            <Input
              label="Nova senha"
              placeholder="Mínimo 6 caracteres"
              value={password}
              onChangeText={setPassword}
              secureTextEntry={!showPassword}
              autoCapitalize="none"
              autoComplete="off"
              textContentType="newPassword"
              error={errors.password}
              leftIcon="lock-closed-outline"
              rightIcon={showPassword ? "eye-off-outline" : "eye-outline"}
              onRightIconPress={() => setShowPassword(!showPassword)}
            />

            <Input
              label="Confirmar nova senha"
              placeholder="Repita a nova senha"
              value={confirm}
              onChangeText={setConfirm}
              secureTextEntry={!showPassword}
              autoCapitalize="none"
              autoComplete="off"
              textContentType="newPassword"
              error={errors.confirm}
              leftIcon="lock-closed-outline"
            />

            <Button
              title={isLoading ? "" : "Redefinir senha"}
              onPress={handleReset}
              disabled={isLoading}
              fullWidth
              icon={
                isLoading ? (
                  <ActivityIndicator color={colors.white} />
                ) : (
                  <Ionicons
                    name="checkmark-circle-outline"
                    size={20}
                    color={colors.white}
                  />
                )
              }
            />

            <Button
              title={
                secondsLeft > 0
                  ? `Reenviar código (${secondsLeft}s)`
                  : "Reenviar código"
              }
              onPress={handleResend}
              variant="outline"
              disabled={resendDisabled}
              fullWidth
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
    paddingTop: spacing.xxl * 2,
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
