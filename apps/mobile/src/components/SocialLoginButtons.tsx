import React, { useEffect, useRef, useState } from "react";
import { Alert, StyleSheet, Text, View } from "react-native";
import { Ionicons } from "@expo/vector-icons";

import { Button } from "./Button";
import { colors, spacing, typography } from "../theme";
import { useAuth } from "../context/AuthContext";
import { ApiError } from "../services/api";
import {
  isGitHubConfigured,
  isGoogleConfigured,
  signInWithGoogleProvider,
  useGitHubAuthRequest,
} from "../services/socialAuth";

type Provider = "Google" | "GitHub";

const showError = (provider: Provider, error: unknown) => {
  if (error instanceof ApiError) {
    Alert.alert("Erro", error.serverMessage);
  } else {
    Alert.alert(
      "Erro",
      `Não foi possível entrar com ${provider}. Tente novamente.`
    );
  }
};

export const SocialLoginButtons: React.FC = () => {
  const { signInWithGoogle, signInWithGitHub } = useAuth();
  const { request, response, promptAsync, redirectUri } =
    useGitHubAuthRequest();
  const [isLoading, setIsLoading] = useState(false);
  const handledResponse = useRef<unknown>(null);

  // Conclui o login do GitHub quando o navegador devolve o code
  useEffect(() => {
    if (!response || handledResponse.current === response) return;
    handledResponse.current = response;

    if (response.type === "success") {
      if (!request?.codeVerifier || !response.params.code) {
        showError("GitHub", new Error("Missing code or code verifier"));
        setIsLoading(false);
        return;
      }
      setIsLoading(true);
      signInWithGitHub({
        code: response.params.code,
        codeVerifier: request.codeVerifier,
        redirectUri,
      })
        .catch((error) => showError("GitHub", error))
        .finally(() => setIsLoading(false));
    } else if (response.type === "error") {
      // "Cancel" na tela de consentimento do GitHub volta como error=access_denied:
      // é um cancelamento, não uma falha (FR-17).
      const cancelled =
        response.params?.error === "access_denied" ||
        response.error?.code === "access_denied";
      if (!cancelled) {
        showError(
          "GitHub",
          new Error(response.error?.message ?? "GitHub error")
        );
      }
      setIsLoading(false);
    } else {
      // cancel / dismiss / locked: sem mensagem
      setIsLoading(false);
    }
  }, [response, request, redirectUri, signInWithGitHub]);

  const handleGoogle = async () => {
    if (!isGoogleConfigured()) {
      Alert.alert("Erro", "Login com Google indisponível.");
      return;
    }
    setIsLoading(true);
    try {
      const result = await signInWithGoogleProvider();
      if (result) {
        await signInWithGoogle(result.idToken);
      }
    } catch (error) {
      showError("Google", error);
    } finally {
      setIsLoading(false);
    }
  };

  const handleGitHub = async () => {
    if (!isGitHubConfigured()) {
      Alert.alert("Erro", "Login com GitHub indisponível.");
      return;
    }
    if (!request) return;
    setIsLoading(true);
    try {
      // O resultado é tratado no useEffect acima
      await promptAsync();
    } catch (error) {
      showError("GitHub", error);
      setIsLoading(false);
    }
  };

  return (
    <View style={styles.container}>
      <View style={styles.divider}>
        <View style={styles.dividerLine} />
        <Text style={styles.dividerText}>ou</Text>
        <View style={styles.dividerLine} />
      </View>

      <Button
        title="Entrar com Google"
        variant="outline"
        onPress={handleGoogle}
        disabled={isLoading}
        fullWidth
        icon={<Ionicons name="logo-google" size={20} color={colors.primary} />}
      />
      <Button
        title="Entrar com GitHub"
        variant="outline"
        onPress={handleGitHub}
        // Sem client ID o toque mostra "indisponível"; com client ID, espera o pedido PKCE ficar pronto
        disabled={isLoading || (isGitHubConfigured() && !request)}
        fullWidth
        icon={<Ionicons name="logo-github" size={20} color={colors.primary} />}
      />
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    gap: spacing.md,
    marginBottom: spacing.xl,
  },
  divider: {
    flexDirection: "row",
    alignItems: "center",
    gap: spacing.sm,
  },
  dividerLine: {
    flex: 1,
    height: 1,
    backgroundColor: colors.border,
  },
  dividerText: {
    fontSize: typography.sm,
    color: colors.textSecondary,
  },
});
