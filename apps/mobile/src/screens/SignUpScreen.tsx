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
import { LinearGradient } from "expo-linear-gradient";
import { Ionicons } from "@expo/vector-icons";
import { useNavigation } from "@react-navigation/native";
import { NativeStackNavigationProp } from "@react-navigation/native-stack";

import { Input, Button } from "../components";
import { colors, spacing, borderRadius, typography } from "../theme";
import { useAuth } from "../context/AuthContext";

type AuthStackParamList = {
  SignIn: undefined;
  SignUp: undefined;
};

type SignUpScreenNavigationProp = NativeStackNavigationProp<
  AuthStackParamList,
  "SignUp"
>;

export const SignUpScreen: React.FC = () => {
  const navigation = useNavigation<SignUpScreenNavigationProp>();
  const { signUp } = useAuth();

  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [showPassword, setShowPassword] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const [errors, setErrors] = useState<{
    name?: string;
    email?: string;
    password?: string;
    confirmPassword?: string;
  }>({});

  const validate = () => {
    const newErrors: {
      name?: string;
      email?: string;
      password?: string;
      confirmPassword?: string;
    } = {};

    if (!name.trim()) {
      newErrors.name = "Nome é obrigatório";
    }

    if (!email.trim()) {
      newErrors.email = "Email é obrigatório";
    } else if (!/\S+@\S+\.\S+/.test(email)) {
      newErrors.email = "Email inválido";
    }

    if (!password) {
      newErrors.password = "Senha é obrigatória";
    } else if (password.length < 6) {
      newErrors.password = "Senha deve ter pelo menos 6 caracteres";
    }

    if (password !== confirmPassword) {
      newErrors.confirmPassword = "As senhas não coincidem";
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSignUp = async () => {
    if (!validate()) return;

    setIsLoading(true);
    try {
      await signUp(email.trim().toLowerCase(), password, name.trim());
    } catch (error) {
      const message =
        error instanceof Error ? error.message : "Erro ao criar conta";
      if (message.includes("409")) {
        Alert.alert("Erro", "Este email já está em uso");
      } else {
        Alert.alert("Erro", "Não foi possível criar a conta. Tente novamente.");
      }
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
        >
          {/* Header */}
          <View style={styles.header}>
            <TouchableOpacity
              onPress={() => navigation.goBack()}
              style={styles.backButton}
            >
              <Ionicons
                name="arrow-back"
                size={24}
                color={colors.textPrimary}
              />
            </TouchableOpacity>
            <LinearGradient
              colors={[colors.primary, colors.primaryLight]}
              style={styles.logoContainer}
            >
              <Ionicons name="person-add" size={36} color={colors.white} />
            </LinearGradient>
            <Text style={styles.title}>Criar Conta</Text>
            <Text style={styles.subtitle}>
              Preencha os dados para começar
            </Text>
          </View>

          {/* Form */}
          <View style={styles.form}>
            <Input
              label="Nome"
              placeholder="Seu nome"
              value={name}
              onChangeText={setName}
              autoCapitalize="words"
              textContentType="name"
              error={errors.name}
              leftIcon="person-outline"
            />

            <Input
              label="Email"
              placeholder="seu@email.com"
              value={email}
              onChangeText={setEmail}
              keyboardType="email-address"
              autoCapitalize="none"
              autoComplete="email"
              textContentType="emailAddress"
              error={errors.email}
              leftIcon="mail-outline"
            />

            <Input
              label="Senha"
              placeholder="Mínimo 6 caracteres"
              value={password}
              onChangeText={setPassword}
              secureTextEntry={!showPassword}
              autoCapitalize="none"
              autoComplete="off"
              textContentType="oneTimeCode"
              error={errors.password}
              leftIcon="lock-closed-outline"
              rightIcon={showPassword ? "eye-off-outline" : "eye-outline"}
              onRightIconPress={() => setShowPassword(!showPassword)}
            />

            <Input
              label="Confirmar Senha"
              placeholder="Repita a senha"
              value={confirmPassword}
              onChangeText={setConfirmPassword}
              secureTextEntry={!showPassword}
              autoCapitalize="none"
              autoComplete="off"
              textContentType="oneTimeCode"
              error={errors.confirmPassword}
              leftIcon="lock-closed-outline"
            />

            <Button
              title={isLoading ? "" : "Criar Conta"}
              onPress={handleSignUp}
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
          </View>

          {/* Footer */}
          <View style={styles.footer}>
            <Text style={styles.footerText}>Já tem uma conta?</Text>
            <TouchableOpacity onPress={() => navigation.goBack()}>
              <Text style={styles.footerLink}>Fazer login</Text>
            </TouchableOpacity>
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
  },
  header: {
    alignItems: "center",
    marginBottom: spacing.xl,
  },
  backButton: {
    position: "absolute",
    left: 0,
    top: 0,
    padding: spacing.sm,
  },
  logoContainer: {
    width: 80,
    height: 80,
    borderRadius: borderRadius.xl,
    alignItems: "center",
    justifyContent: "center",
    marginBottom: spacing.md,
    marginTop: spacing.xl,
  },
  title: {
    fontSize: typography.xxl,
    fontWeight: typography.bold,
    color: colors.textPrimary,
    marginBottom: spacing.xs,
  },
  subtitle: {
    fontSize: typography.md,
    color: colors.textSecondary,
    textAlign: "center",
  },
  form: {
    gap: spacing.md,
    marginBottom: spacing.xl,
  },
  footer: {
    flexDirection: "row",
    justifyContent: "center",
    alignItems: "center",
    gap: spacing.xs,
    marginTop: "auto",
    paddingBottom: spacing.lg,
  },
  footerText: {
    fontSize: typography.md,
    color: colors.textSecondary,
  },
  footerLink: {
    fontSize: typography.md,
    fontWeight: typography.semibold,
    color: colors.primary,
  },
});
