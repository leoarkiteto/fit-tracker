import { Platform } from "react-native";
import * as AuthSession from "expo-auth-session";
import * as WebBrowser from "expo-web-browser";
import {
  GoogleSignin,
  isCancelledResponse,
  isErrorWithCode,
  isSuccessResponse,
  statusCodes,
} from "@react-native-google-signin/google-signin";
import { SOCIAL_AUTH } from "./config";

// Finaliza a sessão de autenticação pendente quando o app é reaberto pelo redirect
WebBrowser.maybeCompleteAuthSession();

// ============= Disponibilidade =============

export const isGoogleConfigured = (): boolean => {
  if (!SOCIAL_AUTH.googleWebClientId) return false;
  if (Platform.OS === "ios" && !SOCIAL_AUTH.googleIosClientId) return false;
  return true;
};

export const isGitHubConfigured = (): boolean => !!SOCIAL_AUTH.githubClientId;

// ============= Google =============

let googleConfigured = false;

const configureGoogle = () => {
  if (googleConfigured) return;
  GoogleSignin.configure({
    webClientId: SOCIAL_AUTH.googleWebClientId,
    iosClientId: SOCIAL_AUTH.googleIosClientId,
  });
  googleConfigured = true;
};

/**
 * Abre o login nativo do Google. Retorna `null` se a pessoa cancelar;
 * lança erro nos demais casos.
 */
export const signInWithGoogleProvider = async (): Promise<{
  idToken: string;
} | null> => {
  configureGoogle();
  try {
    if (Platform.OS === "android") {
      await GoogleSignin.hasPlayServices();
    }
    const response = await GoogleSignin.signIn();
    if (isCancelledResponse(response)) return null;
    if (!isSuccessResponse(response)) {
      throw new Error("Google sign-in failed");
    }
    const idToken = response.data.idToken;
    if (!idToken) {
      throw new Error("Google did not return an ID token");
    }
    return { idToken };
  } catch (error) {
    if (isErrorWithCode(error) && error.code === statusCodes.SIGN_IN_CANCELLED) {
      return null;
    }
    throw error;
  }
};

// ============= GitHub =============

const githubDiscovery: AuthSession.DiscoveryDocument = {
  authorizationEndpoint: "https://github.com/login/oauth/authorize",
  tokenEndpoint: "https://github.com/login/oauth/access_token",
};

/**
 * Prepara o pedido de autorização do GitHub (PKCE). O code NÃO é trocado no app:
 * o backend o troca pelo token, pois só ele conhece o client secret.
 */
export const useGitHubAuthRequest = () => {
  const redirectUri = AuthSession.makeRedirectUri({
    scheme: "fittracker",
    path: "oauth/github",
  });

  const [request, response, promptAsync] = AuthSession.useAuthRequest(
    {
      clientId: SOCIAL_AUTH.githubClientId ?? "",
      scopes: ["read:user", "user:email"],
      usePKCE: true,
      redirectUri,
    },
    githubDiscovery
  );

  return { request, response, promptAsync, redirectUri };
};
