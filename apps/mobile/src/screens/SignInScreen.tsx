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

type SignInScreenNavigationProp = NativeStackNavigationProp<
  AuthStackParamList,
  "SignIn"
>;

export const SignInScreen: React.FC = () => {
  const navigation = useNavigation<SignInScreenNavigationProp>();
  const { signIn } = useAuth();

  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [showPassword, setShowPassword] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const [errors, setErrors] = useState<{ email?: string; password?: string }>(
    {}
  );

  const validate = () => {
    const newErrors: { email?: string; password?: string } = {};

    if (!email.trim()) {
      newErrors.email = "Email é obrigatório";
    } else if (!/\S+@\S+\.\S+/.test(email)) {
      newErrors.email = "Email inválido";
    }

    if (!password) {
      newErrors.password = "Senha é obrigatória";
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSignIn = async () => {
    if (!validate()) return;

    setIsLoading(true);
    try {
      await signIn(email.trim().toLowerCase(), password);
    } catch (error) {
      const message =
        error instanceof Error ? error.message : "Erro ao fazer login";
      if (message.includes("401")) {
        Alert.alert("Erro", "Email ou senha incorretos");
      } else {
        Alert.alert("Erro", "Não foi possível fazer login. Tente novamente.");
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
            <LinearGradient
              colors={[colors.primary, colors.primaryLight]}
              style={styles.logoContainer}
            >
              <Ionicons name="fitness" size={48} color={colors.white} />
            </LinearGradient>
            <Text style={styles.title}>FitTracker</Text>
            <Text style={styles.subtitle}>
              Gerencie seus treinos e alcance seus objetivos
            </Text>
          </View>

          {/* Form */}
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
              error={errors.email}
              leftIcon="mail-outline"
            />

            <Input
              label="Senha"
              placeholder="Sua senha"
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

            <Button
              title={isLoading ? "" : "Entrar"}
              onPress={handleSignIn}
              disabled={isLoading}
              fullWidth
              icon={
                isLoading ? (
                  <ActivityIndicator color={colors.white} />
                ) : (
                  <Ionicons name="log-in-outline" size={20} color={colors.white} />
                )
              }
            />
          </View>

          {/* Footer */}
          <View style={styles.footer}>
            <Text style={styles.footerText}>Não tem uma conta?</Text>
            <TouchableOpacity onPress={() => navigation.navigate("SignUp")}>
              <Text style={styles.footerLink}>Criar conta</Text>
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
    justifyContent: "center",
  },
  header: {
    alignItems: "center",
    marginBottom: spacing.xxl,
  },
  logoContainer: {
    width: 100,
    height: 100,
    borderRadius: borderRadius.xl,
    alignItems: "center",
    justifyContent: "center",
    marginBottom: spacing.lg,
  },
  title: {
    fontSize: typography.xxxl,
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
