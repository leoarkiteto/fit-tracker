import AsyncStorage from "@react-native-async-storage/async-storage";

// Armazena apenas o ID do perfil localmente
// Todos os outros dados são sincronizados com a API

const KEYS = {
  PROFILE_ID: "@fittracker:profile_id",
};

// Profile ID - armazenado localmente para identificar o usuário
export const saveProfileId = async (profileId: string): Promise<void> => {
  try {
    await AsyncStorage.setItem(KEYS.PROFILE_ID, profileId);
  } catch (error) {
    console.error("Error saving profile ID:", error);
    throw error;
  }
};

export const loadProfileId = async (): Promise<string | null> => {
  try {
    return await AsyncStorage.getItem(KEYS.PROFILE_ID);
  } catch (error) {
    console.error("Error loading profile ID:", error);
    return null;
  }
};

export const clearProfileId = async (): Promise<void> => {
  try {
    await AsyncStorage.removeItem(KEYS.PROFILE_ID);
  } catch (error) {
    console.error("Error clearing profile ID:", error);
    throw error;
  }
};

// Clear all data (logout)
export const clearAllData = async (): Promise<void> => {
  try {
    await AsyncStorage.multiRemove(Object.values(KEYS));
  } catch (error) {
    console.error("Error clearing data:", error);
    throw error;
  }
};
