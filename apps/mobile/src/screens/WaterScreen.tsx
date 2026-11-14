import React, { useState, useCallback, useEffect } from "react";
import {
  View,
  Text,
  StyleSheet,
  ScrollView,
  TouchableOpacity,
  ActivityIndicator,
  Alert,
  RefreshControl,
  TextInput,
  KeyboardAvoidingView,
  Platform,
} from "react-native";
import { SafeAreaView } from "react-native-safe-area-context";
import { Ionicons } from "@expo/vector-icons";
import { useNavigation, useFocusEffect } from "@react-navigation/native";
import { NativeStackNavigationProp } from "@react-navigation/native-stack";

import { useApp } from "../context/AppContext";
import { Card, Button } from "../components";
import { colors, spacing, borderRadius, typography } from "../theme";
import { RootStackParamList } from "../navigation/types";
import { waterApi } from "../services";
import type { DailyWaterSummary, WaterIntakeEntry } from "../types";

type WaterScreenNavigationProp = NativeStackNavigationProp<RootStackParamList>;

const todayDateUtc = () => new Date().toISOString().slice(0, 10);

export const WaterScreen: React.FC = () => {
  const navigation = useNavigation<WaterScreenNavigationProp>();
  const { profileId } = useApp();
  const [summary, setSummary] = useState<DailyWaterSummary | null>(null);
  const [loading, setLoading] = useState(true);
  const [adding, setAdding] = useState(false);
  const [customMl, setCustomMl] = useState("");
  const [deletingId, setDeletingId] = useState<string | null>(null);
  const [refreshing, setRefreshing] = useState(false);

  const loadWater = useCallback(async () => {
    if (!profileId) return;
    try {
      const data = await waterApi.getDay(profileId, todayDateUtc());
      setSummary(data);
    } catch (e) {
      console.error("Error loading water:", e);
    } finally {
      setLoading(false);
      setRefreshing(false);
    }
  }, [profileId]);

  useFocusEffect(
    useCallback(() => {
      if (profileId) {
        setLoading(true);
        loadWater();
      }
    }, [profileId, loadWater])
  );

  const onRefresh = useCallback(async () => {
    setRefreshing(true);
    await loadWater();
  }, [loadWater]);

  const addWater = async (amountMl: number) => {
    if (!profileId || adding) return;
    setAdding(true);
    try {
      await waterApi.add(profileId, amountMl);
      setCustomMl("");
      await loadWater();
    } catch (e) {
      console.error("Error adding water:", e);
      Alert.alert(
        "Erro",
        "Não foi possível adicionar. Verifique a conexão e tente novamente."
      );
    } finally {
      setAdding(false);
    }
  };

  const removeEntry = async (entryId: string) => {
    if (!profileId || deletingId) return;
    Alert.alert(
      "Remover registro",
      "Deseja remover este registro de água?",
      [
        { text: "Cancelar", style: "cancel" },
        {
          text: "Remover",
          style: "destructive",
          onPress: async () => {
            setDeletingId(entryId);
            try {
              await waterApi.deleteEntry(profileId, entryId);
              await loadWater();
            } catch (e) {
              console.error("Error removing entry:", e);
            } finally {
              setDeletingId(null);
            }
          },
        },
      ]
    );
  };

  const formatTime = (iso: string) =>
    new Date(iso).toLocaleTimeString("pt-BR", {
      hour: "2-digit",
      minute: "2-digit",
    });

  if (loading && !summary) {
    return (
      <SafeAreaView style={styles.container} edges={["top"]}>
        <View style={styles.header}>
          <TouchableOpacity
            onPress={() => navigation.goBack()}
            style={styles.backButton}
          >
            <Ionicons name="arrow-back" size={24} color={colors.textPrimary} />
          </TouchableOpacity>
          <Text style={styles.headerTitle}>Água</Text>
          <View style={styles.backButton} />
        </View>
        <View style={styles.loadingContainer}>
          <ActivityIndicator size="large" color={colors.primary} />
          <Text style={styles.loadingText}>Carregando...</Text>
        </View>
      </SafeAreaView>
    );
  }

  const totalMl = summary?.totalMl ?? 0;
  const goalMl = summary?.goalMl ?? 2000;
  const progressPct = goalMl > 0 ? Math.min(100, (totalMl / goalMl) * 100) : 0;

  return (
    <SafeAreaView style={styles.container} edges={["top"]}>
      <View style={styles.header}>
        <TouchableOpacity
          onPress={() => navigation.goBack()}
          style={styles.backButton}
        >
          <Ionicons name="arrow-back" size={24} color={colors.textPrimary} />
        </TouchableOpacity>
        <Text style={styles.headerTitle}>Consumo de água</Text>
        <View style={styles.backButton} />
      </View>

      <ScrollView
        style={styles.scroll}
        contentContainerStyle={styles.scrollContent}
        showsVerticalScrollIndicator={false}
        refreshControl={
          <RefreshControl
            refreshing={refreshing}
            onRefresh={onRefresh}
            colors={[colors.primary]}
            tintColor={colors.primary}
          />
        }
      >
        <Card style={styles.summaryCard}>
          <View style={styles.summaryHeader}>
            <View style={styles.waterIconContainer}>
              <Ionicons name="water" size={32} color="#0891B2" />
            </View>
            <View style={styles.summaryText}>
              <Text style={styles.summaryLabel}>Hoje</Text>
              <Text style={styles.summaryValue}>
                {(totalMl / 1000).toFixed(1)} L / {(goalMl / 1000).toFixed(1)} L
              </Text>
            </View>
          </View>
          <View style={styles.progressBar}>
            <View
              style={[styles.progressFill, { width: `${progressPct}%` }]}
            />
          </View>
          <View style={styles.quickButtons}>
            <Button
              title="+250 ml"
              onPress={() => addWater(250)}
              disabled={adding}
              variant="secondary"
              style={styles.quickButton}
            />
            <Button
              title="+500 ml"
              onPress={() => addWater(500)}
              disabled={adding}
              variant="secondary"
              style={styles.quickButton}
            />
          </View>
        </Card>

        <Text style={styles.sectionTitle}>Adicionar quantidade</Text>
        <KeyboardAvoidingView
          behavior={Platform.OS === "ios" ? "padding" : undefined}
        >
          <View style={styles.addRow}>
            <TextInput
              style={styles.input}
              placeholder="ml (ex: 300)"
              placeholderTextColor={colors.textLight}
              value={customMl}
              onChangeText={setCustomMl}
              keyboardType="number-pad"
              maxLength={4}
            />
            <Button
              title="Adicionar"
              onPress={() => {
                const n = parseInt(customMl, 10);
                if (!isNaN(n) && n > 0 && n <= 2000) addWater(n);
              }}
              disabled={adding || !customMl.trim()}
              style={styles.addButton}
            />
          </View>
        </KeyboardAvoidingView>

        <Text style={styles.sectionTitle}>Registros de hoje</Text>
        {!summary?.entries?.length ? (
          <Card style={styles.emptyCard}>
            <Ionicons
              name="water-outline"
              size={48}
              color={colors.textLight}
            />
            <Text style={styles.emptyText}>
              Nenhum registro ainda. Use os botões acima para adicionar.
            </Text>
          </Card>
        ) : (
          <View style={styles.entriesList}>
            {([...(summary?.entries ?? [])] as WaterIntakeEntry[])
              .reverse()
              .map((entry) => (
                <Card key={entry.id} style={styles.entryCard}>
                  <View style={styles.entryRow}>
                    <View>
                      <Text style={styles.entryAmount}>{entry.amountMl} ml</Text>
                      <Text style={styles.entryTime}>
                        {formatTime(entry.consumedAt)}
                      </Text>
                    </View>
                    <TouchableOpacity
                      onPress={() => removeEntry(entry.id)}
                      disabled={deletingId === entry.id}
                      style={styles.deleteButton}
                    >
                      {deletingId === entry.id ? (
                        <ActivityIndicator size="small" color={colors.error} />
                      ) : (
                        <Ionicons
                          name="trash-outline"
                          size={22}
                          color={colors.error}
                        />
                      )}
                    </TouchableOpacity>
                  </View>
                </Card>
              ))}
          </View>
        )}
        <View style={{ height: 100 }} />
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
    flexDirection: "row",
    alignItems: "center",
    justifyContent: "space-between",
    paddingHorizontal: spacing.md,
    paddingVertical: spacing.md,
    borderBottomWidth: 1,
    borderBottomColor: colors.border,
  },
  backButton: {
    width: 40,
    height: 40,
    alignItems: "center",
    justifyContent: "center",
  },
  headerTitle: {
    fontSize: typography.lg,
    fontWeight: typography.bold,
    color: colors.textPrimary,
  },
  loadingContainer: {
    flex: 1,
    alignItems: "center",
    justifyContent: "center",
    gap: spacing.md,
  },
  loadingText: {
    fontSize: typography.md,
    color: colors.textSecondary,
  },
  scroll: {
    flex: 1,
  },
  scrollContent: {
    padding: spacing.lg,
  },
  summaryCard: {
    padding: spacing.lg,
    marginBottom: spacing.lg,
  },
  summaryHeader: {
    flexDirection: "row",
    alignItems: "center",
    marginBottom: spacing.md,
  },
  waterIconContainer: {
    width: 56,
    height: 56,
    borderRadius: borderRadius.lg,
    backgroundColor: "#E0F2F1",
    alignItems: "center",
    justifyContent: "center",
  },
  summaryText: {
    marginLeft: spacing.md,
  },
  summaryLabel: {
    fontSize: typography.sm,
    color: colors.textSecondary,
  },
  summaryValue: {
    fontSize: typography.xxl,
    fontWeight: typography.bold,
    color: colors.textPrimary,
  },
  progressBar: {
    height: 8,
    backgroundColor: colors.surfaceSecondary,
    borderRadius: 4,
    overflow: "hidden",
    marginBottom: spacing.lg,
  },
  progressFill: {
    height: "100%",
    backgroundColor: "#0891B2",
    borderRadius: 4,
  },
  quickButtons: {
    flexDirection: "row",
    gap: spacing.sm,
  },
  quickButton: {
    flex: 1,
  },
  sectionTitle: {
    fontSize: typography.lg,
    fontWeight: typography.bold,
    color: colors.textPrimary,
    marginBottom: spacing.md,
  },
  addRow: {
    flexDirection: "row",
    gap: spacing.sm,
    marginBottom: spacing.xl,
  },
  input: {
    flex: 1,
    height: 48,
    backgroundColor: colors.surface,
    borderWidth: 1,
    borderColor: colors.border,
    borderRadius: borderRadius.md,
    paddingHorizontal: spacing.md,
    fontSize: typography.md,
    color: colors.textPrimary,
  },
  addButton: {
    minWidth: 120,
  },
  emptyCard: {
    alignItems: "center",
    padding: spacing.xl,
    gap: spacing.md,
  },
  emptyText: {
    fontSize: typography.md,
    color: colors.textSecondary,
    textAlign: "center",
  },
  entriesList: {
    gap: spacing.sm,
  },
  entryCard: {
    padding: spacing.md,
  },
  entryRow: {
    flexDirection: "row",
    alignItems: "center",
    justifyContent: "space-between",
  },
  entryAmount: {
    fontSize: typography.md,
    fontWeight: typography.semibold,
    color: colors.textPrimary,
  },
  entryTime: {
    fontSize: typography.sm,
    color: colors.textSecondary,
    marginTop: spacing.xs,
  },
  deleteButton: {
    padding: spacing.sm,
  },
});
