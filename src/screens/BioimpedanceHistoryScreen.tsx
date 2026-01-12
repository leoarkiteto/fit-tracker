import React, { useState, useCallback } from "react";
import {
  View,
  Text,
  StyleSheet,
  ScrollView,
  TouchableOpacity,
  Alert,
  RefreshControl,
} from "react-native";
import { SafeAreaView } from "react-native-safe-area-context";
import { LinearGradient } from "expo-linear-gradient";
import { Ionicons } from "@expo/vector-icons";
import { useNavigation } from "@react-navigation/native";
import { NativeStackNavigationProp } from "@react-navigation/native-stack";

import { useApp } from "../context/AppContext";
import { Card } from "../components";
import { colors, spacing, borderRadius, typography, shadows } from "../theme";
import { RootStackParamList } from "../navigation/types";
import { BioimpedanceData } from "../types";

type BioimpedanceHistoryScreenNavigationProp =
  NativeStackNavigationProp<RootStackParamList>;

export const BioimpedanceHistoryScreen: React.FC = () => {
  const navigation = useNavigation<BioimpedanceHistoryScreenNavigationProp>();
  const { bioimpedanceHistory, deleteBioimpedance, refreshData } = useApp();
  const [refreshing, setRefreshing] = useState(false);

  const onRefresh = useCallback(async () => {
    setRefreshing(true);
    await refreshData();
    setRefreshing(false);
  }, [refreshData]);

  // Sort by date descending
  const sortedHistory = [...bioimpedanceHistory].sort(
    (a, b) => new Date(b.date).getTime() - new Date(a.date).getTime(),
  );

  const latestData = sortedHistory[0];
  const previousData = sortedHistory[1];

  const getChange = (current: number, previous: number | undefined) => {
    if (!previous) return null;
    return current - previous;
  };

  const formatChange = (
    change: number | null,
    unit: string,
    inverse: boolean = false,
  ) => {
    if (change === null) return null;
    const isPositive = inverse ? change < 0 : change > 0;
    const color = isPositive ? colors.success : colors.error;
    const arrow = change > 0 ? "↑" : "↓";
    return (
      <Text style={[styles.changeText, { color }]}>
        {arrow} {Math.abs(change).toFixed(1)}
        {unit}
      </Text>
    );
  };

  const handleDelete = (id: string) => {
    Alert.alert(
      "Excluir Registro",
      "Deseja excluir este registro de bioimpedância?",
      [
        { text: "Cancelar", style: "cancel" },
        {
          text: "Excluir",
          style: "destructive",
          onPress: () => deleteBioimpedance(id),
        },
      ],
    );
  };

  const HistoryCard: React.FC<{
    data: BioimpedanceData;
    isLatest?: boolean;
  }> = ({ data, isLatest }) => (
    <Card
      style={[styles.historyCard, isLatest && styles.latestCard]}
      onPress={() =>
        navigation.navigate("BioimpedanceForm", { bioimpedanceId: data.id })
      }
    >
      <View style={styles.historyHeader}>
        <View>
          <Text style={styles.historyDate}>
            {new Date(data.date).toLocaleDateString("pt-BR", {
              day: "2-digit",
              month: "long",
              year: "numeric",
            })}
          </Text>
          {isLatest && (
            <View style={styles.latestBadge}>
              <Text style={styles.latestBadgeText}>Mais Recente</Text>
            </View>
          )}
        </View>
        <TouchableOpacity
          onPress={() => handleDelete(data.id)}
          style={styles.deleteButton}
        >
          <Ionicons name="trash-outline" size={18} color={colors.error} />
        </TouchableOpacity>
      </View>

      <View style={styles.historyGrid}>
        <View style={styles.historyItem}>
          <Text style={styles.historyValue}>{data.weight}</Text>
          <Text style={styles.historyLabel}>Peso (kg)</Text>
        </View>
        <View style={styles.historyItem}>
          <Text style={styles.historyValue}>{data.bodyFatPercentage}%</Text>
          <Text style={styles.historyLabel}>Gordura</Text>
        </View>
        <View style={styles.historyItem}>
          <Text style={styles.historyValue}>{data.muscleMass}</Text>
          <Text style={styles.historyLabel}>Músculo (kg)</Text>
        </View>
        <View style={styles.historyItem}>
          <Text style={styles.historyValue}>{data.waterPercentage}%</Text>
          <Text style={styles.historyLabel}>Água</Text>
        </View>
      </View>

      {data.notes && (
        <Text style={styles.historyNotes} numberOfLines={2}>
          {data.notes}
        </Text>
      )}
    </Card>
  );

  return (
    <SafeAreaView style={styles.container} edges={["top"]}>
      {/* Header */}
      <View style={styles.header}>
        <TouchableOpacity
          onPress={() => navigation.goBack()}
          style={styles.backButton}
        >
          <Ionicons name="arrow-back" size={24} color={colors.textPrimary} />
        </TouchableOpacity>
        <Text style={styles.title}>Bioimpedância</Text>
        <TouchableOpacity
          onPress={() => navigation.navigate("BioimpedanceForm", {})}
          style={styles.addButton}
        >
          <Ionicons name="add" size={24} color={colors.white} />
        </TouchableOpacity>
      </View>

      <ScrollView
        style={styles.content}
        contentContainerStyle={styles.contentContainer}
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
        {latestData ? (
          <>
            {/* Summary Card */}
            <LinearGradient
              colors={[colors.secondary, colors.secondaryLight]}
              start={{ x: 0, y: 0 }}
              end={{ x: 1, y: 1 }}
              style={styles.summaryCard}
            >
              <Text style={styles.summaryTitle}>Resumo Atual</Text>

              <View style={styles.summaryGrid}>
                <View style={styles.summaryItem}>
                  <Text style={styles.summaryValue}>{latestData.weight}</Text>
                  <Text style={styles.summaryLabel}>kg</Text>
                  {formatChange(
                    getChange(latestData.weight, previousData?.weight),
                    "kg",
                    true,
                  )}
                </View>
                <View style={styles.summaryDivider} />
                <View style={styles.summaryItem}>
                  <Text style={styles.summaryValue}>
                    {latestData.bodyFatPercentage}
                  </Text>
                  <Text style={styles.summaryLabel}>% gordura</Text>
                  {formatChange(
                    getChange(
                      latestData.bodyFatPercentage,
                      previousData?.bodyFatPercentage,
                    ),
                    "%",
                    true,
                  )}
                </View>
                <View style={styles.summaryDivider} />
                <View style={styles.summaryItem}>
                  <Text style={styles.summaryValue}>
                    {latestData.muscleMass}
                  </Text>
                  <Text style={styles.summaryLabel}>kg músculo</Text>
                  {formatChange(
                    getChange(latestData.muscleMass, previousData?.muscleMass),
                    "kg",
                  )}
                </View>
              </View>

              {latestData.metabolicAge > 0 && (
                <View style={styles.metabolicAge}>
                  <Ionicons
                    name="body-outline"
                    size={20}
                    color="rgba(255,255,255,0.8)"
                  />
                  <Text style={styles.metabolicAgeText}>
                    Idade Metabólica: {latestData.metabolicAge} anos
                  </Text>
                </View>
              )}
            </LinearGradient>

            {/* Additional Stats */}
            {(latestData.visceralFat > 0 || latestData.bmr > 0) && (
              <View style={styles.statsRow}>
                {latestData.visceralFat > 0 && (
                  <Card style={styles.statCard}>
                    <Ionicons
                      name="analytics-outline"
                      size={24}
                      color={
                        latestData.visceralFat <= 9
                          ? colors.success
                          : colors.warning
                      }
                    />
                    <Text style={styles.statValue}>
                      {latestData.visceralFat}
                    </Text>
                    <Text style={styles.statLabel}>Gordura Visceral</Text>
                    <Text
                      style={[
                        styles.statStatus,
                        {
                          color:
                            latestData.visceralFat <= 9
                              ? colors.success
                              : colors.warning,
                        },
                      ]}
                    >
                      {latestData.visceralFat <= 9 ? "Saudável" : "Elevado"}
                    </Text>
                  </Card>
                )}
                {latestData.bmr > 0 && (
                  <Card style={styles.statCard}>
                    <Ionicons
                      name="flame-outline"
                      size={24}
                      color={colors.warning}
                    />
                    <Text style={styles.statValue}>{latestData.bmr}</Text>
                    <Text style={styles.statLabel}>TMB (kcal)</Text>
                    <Text style={styles.statStatus}>Taxa Metabólica</Text>
                  </Card>
                )}
              </View>
            )}

            {/* History */}
            <Text style={styles.sectionTitle}>Histórico</Text>
            {sortedHistory.map((data, index) => (
              <HistoryCard key={data.id} data={data} isLatest={index === 0} />
            ))}
          </>
        ) : (
          <Card style={styles.emptyCard}>
            <Ionicons name="body-outline" size={64} color={colors.textLight} />
            <Text style={styles.emptyTitle}>Nenhum registro</Text>
            <Text style={styles.emptyText}>
              Registre sua primeira avaliação de bioimpedância para acompanhar
              sua evolução
            </Text>
            <TouchableOpacity
              onPress={() => navigation.navigate("BioimpedanceForm", {})}
              style={styles.emptyButton}
            >
              <Ionicons name="add-circle" size={20} color={colors.white} />
              <Text style={styles.emptyButtonText}>Registrar Avaliação</Text>
            </TouchableOpacity>
          </Card>
        )}

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
    flexDirection: "row",
    alignItems: "center",
    justifyContent: "space-between",
    paddingHorizontal: spacing.md,
    paddingVertical: spacing.md,
    borderBottomWidth: 1,
    borderBottomColor: colors.border,
    backgroundColor: colors.surface,
  },
  backButton: {
    width: 40,
    height: 40,
    alignItems: "center",
    justifyContent: "center",
  },
  title: {
    fontSize: typography.lg,
    fontWeight: typography.bold,
    color: colors.textPrimary,
  },
  addButton: {
    width: 40,
    height: 40,
    borderRadius: borderRadius.full,
    backgroundColor: colors.primary,
    alignItems: "center",
    justifyContent: "center",
  },
  content: {
    flex: 1,
  },
  contentContainer: {
    padding: spacing.lg,
  },
  summaryCard: {
    borderRadius: borderRadius.xl,
    padding: spacing.lg,
    marginBottom: spacing.lg,
  },
  summaryTitle: {
    fontSize: typography.lg,
    fontWeight: typography.bold,
    color: colors.white,
    marginBottom: spacing.lg,
  },
  summaryGrid: {
    flexDirection: "row",
    justifyContent: "space-around",
  },
  summaryItem: {
    alignItems: "center",
  },
  summaryValue: {
    fontSize: typography.xxxl,
    fontWeight: typography.bold,
    color: colors.white,
  },
  summaryLabel: {
    fontSize: typography.xs,
    color: "rgba(255,255,255,0.7)",
    marginTop: spacing.xs,
  },
  summaryDivider: {
    width: 1,
    backgroundColor: "rgba(255,255,255,0.2)",
  },
  changeText: {
    fontSize: typography.xs,
    fontWeight: typography.semibold,
    marginTop: spacing.xs,
  },
  metabolicAge: {
    flexDirection: "row",
    alignItems: "center",
    justifyContent: "center",
    gap: spacing.sm,
    marginTop: spacing.lg,
    paddingTop: spacing.md,
    borderTopWidth: 1,
    borderTopColor: "rgba(255,255,255,0.2)",
  },
  metabolicAgeText: {
    fontSize: typography.sm,
    color: "rgba(255,255,255,0.9)",
    fontWeight: typography.medium,
  },
  statsRow: {
    flexDirection: "row",
    gap: spacing.md,
    marginBottom: spacing.lg,
  },
  statCard: {
    flex: 1,
    alignItems: "center",
    padding: spacing.lg,
    gap: spacing.xs,
  },
  statValue: {
    fontSize: typography.xxl,
    fontWeight: typography.bold,
    color: colors.textPrimary,
  },
  statLabel: {
    fontSize: typography.xs,
    color: colors.textSecondary,
  },
  statStatus: {
    fontSize: typography.xs,
    fontWeight: typography.medium,
  },
  sectionTitle: {
    fontSize: typography.lg,
    fontWeight: typography.bold,
    color: colors.textPrimary,
    marginBottom: spacing.md,
  },
  historyCard: {
    padding: spacing.lg,
    marginBottom: spacing.md,
  },
  latestCard: {
    borderWidth: 2,
    borderColor: colors.secondary,
  },
  historyHeader: {
    flexDirection: "row",
    justifyContent: "space-between",
    alignItems: "flex-start",
    marginBottom: spacing.md,
  },
  historyDate: {
    fontSize: typography.md,
    fontWeight: typography.semibold,
    color: colors.textPrimary,
  },
  latestBadge: {
    backgroundColor: `${colors.secondary}20`,
    paddingHorizontal: spacing.sm,
    paddingVertical: spacing.xs,
    borderRadius: borderRadius.full,
    marginTop: spacing.xs,
  },
  latestBadgeText: {
    fontSize: typography.xs,
    fontWeight: typography.semibold,
    color: colors.secondary,
  },
  deleteButton: {
    padding: spacing.xs,
  },
  historyGrid: {
    flexDirection: "row",
    justifyContent: "space-between",
  },
  historyItem: {
    alignItems: "center",
  },
  historyValue: {
    fontSize: typography.lg,
    fontWeight: typography.bold,
    color: colors.textPrimary,
  },
  historyLabel: {
    fontSize: typography.xs,
    color: colors.textSecondary,
    marginTop: spacing.xs,
  },
  historyNotes: {
    fontSize: typography.xs,
    color: colors.textLight,
    fontStyle: "italic",
    marginTop: spacing.md,
    paddingTop: spacing.md,
    borderTopWidth: 1,
    borderTopColor: colors.border,
  },
  emptyCard: {
    alignItems: "center",
    padding: spacing.xxl,
    gap: spacing.md,
  },
  emptyTitle: {
    fontSize: typography.lg,
    fontWeight: typography.semibold,
    color: colors.textPrimary,
  },
  emptyText: {
    fontSize: typography.md,
    color: colors.textSecondary,
    textAlign: "center",
  },
  emptyButton: {
    flexDirection: "row",
    alignItems: "center",
    gap: spacing.sm,
    backgroundColor: colors.primary,
    paddingVertical: spacing.md,
    paddingHorizontal: spacing.lg,
    borderRadius: borderRadius.full,
    marginTop: spacing.md,
  },
  emptyButtonText: {
    color: colors.white,
    fontWeight: typography.semibold,
    fontSize: typography.md,
  },
});
