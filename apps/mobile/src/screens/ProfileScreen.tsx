import React, { useState } from "react";
import {
  View,
  Text,
  StyleSheet,
  ScrollView,
  TouchableOpacity,
  Alert,
} from "react-native";
import { SafeAreaView } from "react-native-safe-area-context";
import { LinearGradient } from "expo-linear-gradient";
import { Ionicons } from "@expo/vector-icons";
import { useNavigation } from "@react-navigation/native";
import { NativeStackNavigationProp } from "@react-navigation/native-stack";

import { useApp } from "../context/AppContext";
import { useAuth } from "../context/AuthContext";
import { Input, Card } from "../components";
import { colors, spacing, borderRadius, typography } from "../theme";
import { RootStackParamList } from "../navigation/types";

type ProfileScreenNavigationProp =
  NativeStackNavigationProp<RootStackParamList>;

export const ProfileScreen: React.FC = () => {
  const navigation = useNavigation<ProfileScreenNavigationProp>();
  const { signOut, user } = useAuth();
  const {
    profile,
    updateProfile,
    totalWorkoutsCompleted,
    totalMinutesSpent,
    workouts,
    bioimpedanceHistory,
  } = useApp();

  const handleLogout = () => {
    Alert.alert(
      "Sair",
      "Tem certeza que deseja sair da sua conta?",
      [
        { text: "Cancelar", style: "cancel" },
        {
          text: "Sair",
          style: "destructive",
          onPress: async () => {
            await signOut();
          },
        },
      ]
    );
  };

  const [isEditing, setIsEditing] = useState(false);
  const [name, setName] = useState(profile.name);
  const [age, setAge] = useState(profile.age?.toString() || "");
  const [height, setHeight] = useState(profile.height?.toString() || "");
  const [currentWeight, setCurrentWeight] = useState(
    profile.currentWeight?.toString() || "",
  );
  const [goalWeight, setGoalWeight] = useState(
    profile.goalWeight?.toString() || "",
  );

  const handleSave = async () => {
    if (!name.trim()) {
      Alert.alert("Erro", "Digite seu nome");
      return;
    }

    await updateProfile({
      name: name.trim(),
      age: age ? parseInt(age) : undefined,
      height: height ? parseFloat(height) : undefined,
      currentWeight: currentWeight ? parseFloat(currentWeight) : undefined,
      goalWeight: goalWeight ? parseFloat(goalWeight) : undefined,
    });

    setIsEditing(false);
  };

  const MenuItem: React.FC<{
    icon: string;
    title: string;
    subtitle?: string;
    onPress?: () => void;
    color?: string;
  }> = ({ icon, title, subtitle, onPress, color = colors.textPrimary }) => (
    <TouchableOpacity onPress={onPress} style={styles.menuItem}>
      <View style={[styles.menuIcon, { backgroundColor: `${color}10` }]}>
        <Ionicons name={icon as any} size={22} color={color} />
      </View>
      <View style={styles.menuContent}>
        <Text style={[styles.menuTitle, { color }]}>{title}</Text>
        {subtitle && <Text style={styles.menuSubtitle}>{subtitle}</Text>}
      </View>
      <Ionicons name="chevron-forward" size={20} color={colors.textLight} />
    </TouchableOpacity>
  );

  return (
    <SafeAreaView style={styles.container} edges={["top"]}>
      <ScrollView showsVerticalScrollIndicator={false}>
        {/* Header */}
        <LinearGradient
          colors={[colors.primary, colors.primaryDark]}
          style={styles.header}
        >
          <View style={styles.avatarContainer}>
            <Text style={styles.avatarText}>
              {profile.name.charAt(0).toUpperCase()}
            </Text>
          </View>

          {isEditing ? (
            <View style={styles.editNameContainer}>
              <Input
                value={name}
                onChangeText={setName}
                placeholder="Seu nome"
                containerStyle={styles.nameInput}
              />
            </View>
          ) : (
            <>
              <Text style={styles.profileName}>{profile.name}</Text>
              {profile.age && (
                <Text style={styles.profileAge}>{profile.age} anos</Text>
              )}
            </>
          )}

          <TouchableOpacity
            onPress={() => (isEditing ? handleSave() : setIsEditing(true))}
            style={styles.editButton}
          >
            <Ionicons
              name={isEditing ? "checkmark" : "pencil"}
              size={20}
              color={colors.white}
            />
          </TouchableOpacity>
        </LinearGradient>

        {/* Stats */}
        <View style={styles.statsContainer}>
          <Card style={styles.statCard}>
            <Text style={styles.statValue}>{workouts.length}</Text>
            <Text style={styles.statLabel}>Treinos</Text>
          </Card>
          <Card style={styles.statCard}>
            <Text style={styles.statValue}>{totalWorkoutsCompleted}</Text>
            <Text style={styles.statLabel}>Concluídos</Text>
          </Card>
          <Card style={styles.statCard}>
            <Text style={styles.statValue}>{bioimpedanceHistory.length}</Text>
            <Text style={styles.statLabel}>Avaliações</Text>
          </Card>
        </View>

        {/* Body Info */}
        {isEditing && (
          <View style={styles.section}>
            <Text style={styles.sectionTitle}>Informações do Corpo</Text>
            <Card style={styles.editCard}>
              <View style={styles.row}>
                <View style={styles.halfWidth}>
                  <Input
                    label="Idade"
                    value={age}
                    onChangeText={setAge}
                    keyboardType="numeric"
                    placeholder="Ex: 28"
                  />
                </View>
                <View style={styles.halfWidth}>
                  <Input
                    label="Altura (cm)"
                    value={height}
                    onChangeText={setHeight}
                    keyboardType="decimal-pad"
                    placeholder="Ex: 175"
                  />
                </View>
              </View>
              <View style={styles.row}>
                <View style={styles.halfWidth}>
                  <Input
                    label="Peso Atual (kg)"
                    value={currentWeight}
                    onChangeText={setCurrentWeight}
                    keyboardType="decimal-pad"
                    placeholder="Ex: 75"
                  />
                </View>
                <View style={styles.halfWidth}>
                  <Input
                    label="Peso Meta (kg)"
                    value={goalWeight}
                    onChangeText={setGoalWeight}
                    keyboardType="decimal-pad"
                    placeholder="Ex: 70"
                  />
                </View>
              </View>
            </Card>
          </View>
        )}

        {/* Body Measurements */}
        {!isEditing &&
          (profile.height || profile.currentWeight || profile.goalWeight) && (
            <View style={styles.section}>
              <Text style={styles.sectionTitle}>Medidas</Text>
              <Card style={styles.measurementsCard}>
                {profile.height && (
                  <View style={styles.measurementItem}>
                    <Ionicons
                      name="resize-outline"
                      size={24}
                      color={colors.primary}
                    />
                    <View>
                      <Text style={styles.measurementValue}>
                        {profile.height} cm
                      </Text>
                      <Text style={styles.measurementLabel}>Altura</Text>
                    </View>
                  </View>
                )}
                {profile.currentWeight && (
                  <View style={styles.measurementItem}>
                    <Ionicons
                      name="scale-outline"
                      size={24}
                      color={colors.secondary}
                    />
                    <View>
                      <Text style={styles.measurementValue}>
                        {profile.currentWeight} kg
                      </Text>
                      <Text style={styles.measurementLabel}>Peso Atual</Text>
                    </View>
                  </View>
                )}
                {profile.goalWeight && (
                  <View style={styles.measurementItem}>
                    <Ionicons
                      name="flag-outline"
                      size={24}
                      color={colors.success}
                    />
                    <View>
                      <Text style={styles.measurementValue}>
                        {profile.goalWeight} kg
                      </Text>
                      <Text style={styles.measurementLabel}>Meta</Text>
                    </View>
                  </View>
                )}
              </Card>
            </View>
          )}

        {/* Menu */}
        <View style={styles.section}>
          <Text style={styles.sectionTitle}>Opções</Text>
          <Card padding="none">
            <MenuItem
              icon="barbell-outline"
              title="Meus Treinos"
              subtitle={`${workouts.length} treinos criados`}
              onPress={() =>
                navigation.navigate("MainTabs", { screen: "Workouts" })
              }
            />
            <View style={styles.menuDivider} />
            <MenuItem
              icon="body-outline"
              title="Bioimpedância"
              subtitle={`${bioimpedanceHistory.length} avaliações`}
              onPress={() => navigation.navigate("BioimpedanceHistory")}
            />
            <View style={styles.menuDivider} />
            <MenuItem
              icon="stats-chart-outline"
              title="Estatísticas"
              subtitle={`${totalMinutesSpent} minutos de treino`}
            />
          </Card>
        </View>

        {/* Account */}
        <View style={styles.section}>
          <Text style={styles.sectionTitle}>Conta</Text>
          <Card padding="none">
            <MenuItem
              icon="person-circle-outline"
              title="Email"
              subtitle={user?.email || ""}
            />
            <View style={styles.menuDivider} />
            <MenuItem
              icon="key-outline"
              title="Alterar Senha"
              onPress={() => navigation.navigate("ChangePassword")}
            />
            <View style={styles.menuDivider} />
            <MenuItem
              icon="log-out-outline"
              title="Sair da Conta"
              onPress={handleLogout}
              color={colors.error}
            />
          </Card>
        </View>

        {/* App Info */}
        <View style={styles.section}>
          <Card style={styles.appInfo}>
            <Ionicons name="fitness" size={40} color={colors.primary} />
            <Text style={styles.appName}>FitTracker</Text>
            <Text style={styles.appVersion}>Versão 1.0.0</Text>
          </Card>
        </View>

        <View style={{ height: spacing.xxl }} />
      </ScrollView>
    </SafeAreaView>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: colors.background,
  },
  header: {
    alignItems: "center",
    paddingTop: spacing.xl,
    paddingBottom: spacing.xxl,
    borderBottomLeftRadius: borderRadius.xl,
    borderBottomRightRadius: borderRadius.xl,
  },
  avatarContainer: {
    width: 100,
    height: 100,
    borderRadius: 50,
    backgroundColor: "rgba(255,255,255,0.2)",
    alignItems: "center",
    justifyContent: "center",
    marginBottom: spacing.md,
    borderWidth: 3,
    borderColor: "rgba(255,255,255,0.3)",
  },
  avatarText: {
    fontSize: 40,
    fontWeight: typography.bold,
    color: colors.white,
  },
  profileName: {
    fontSize: typography.xxl,
    fontWeight: typography.bold,
    color: colors.white,
  },
  profileAge: {
    fontSize: typography.md,
    color: "rgba(255,255,255,0.8)",
    marginTop: spacing.xs,
  },
  editButton: {
    position: "absolute",
    top: spacing.lg,
    right: spacing.lg,
    width: 40,
    height: 40,
    borderRadius: borderRadius.full,
    backgroundColor: "rgba(255,255,255,0.2)",
    alignItems: "center",
    justifyContent: "center",
  },
  editNameContainer: {
    width: "80%",
    marginTop: spacing.md,
  },
  nameInput: {
    marginBottom: 0,
  },
  statsContainer: {
    flexDirection: "row",
    paddingHorizontal: spacing.lg,
    marginTop: -spacing.xl,
    gap: spacing.sm,
  },
  statCard: {
    flex: 1,
    alignItems: "center",
    padding: spacing.md,
  },
  statValue: {
    fontSize: typography.xxl,
    fontWeight: typography.bold,
    color: colors.textPrimary,
  },
  statLabel: {
    fontSize: typography.xs,
    color: colors.textSecondary,
    marginTop: spacing.xs,
  },
  section: {
    paddingHorizontal: spacing.lg,
    marginTop: spacing.lg,
  },
  sectionTitle: {
    fontSize: typography.md,
    fontWeight: typography.bold,
    color: colors.textPrimary,
    marginBottom: spacing.md,
  },
  editCard: {
    padding: spacing.lg,
  },
  row: {
    flexDirection: "row",
    gap: spacing.md,
  },
  halfWidth: {
    flex: 1,
  },
  measurementsCard: {
    flexDirection: "row",
    justifyContent: "space-around",
    padding: spacing.lg,
  },
  measurementItem: {
    alignItems: "center",
    gap: spacing.sm,
  },
  measurementValue: {
    fontSize: typography.lg,
    fontWeight: typography.bold,
    color: colors.textPrimary,
    textAlign: "center",
  },
  measurementLabel: {
    fontSize: typography.xs,
    color: colors.textSecondary,
    textAlign: "center",
  },
  menuItem: {
    flexDirection: "row",
    alignItems: "center",
    padding: spacing.md,
    gap: spacing.md,
  },
  menuIcon: {
    width: 44,
    height: 44,
    borderRadius: borderRadius.md,
    alignItems: "center",
    justifyContent: "center",
  },
  menuContent: {
    flex: 1,
  },
  menuTitle: {
    fontSize: typography.md,
    fontWeight: typography.medium,
  },
  menuSubtitle: {
    fontSize: typography.xs,
    color: colors.textSecondary,
    marginTop: 2,
  },
  menuDivider: {
    height: 1,
    backgroundColor: colors.border,
    marginLeft: spacing.md + 44 + spacing.md,
  },
  appInfo: {
    alignItems: "center",
    padding: spacing.xl,
    gap: spacing.sm,
  },
  appName: {
    fontSize: typography.xl,
    fontWeight: typography.bold,
    color: colors.textPrimary,
  },
  appVersion: {
    fontSize: typography.sm,
    color: colors.textSecondary,
  },
});
