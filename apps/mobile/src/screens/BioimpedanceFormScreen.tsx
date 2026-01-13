import React, { useState } from "react";
import {
  View,
  Text,
  StyleSheet,
  ScrollView,
  TouchableOpacity,
  Alert,
  KeyboardAvoidingView,
  Platform,
} from "react-native";
import { SafeAreaView } from "react-native-safe-area-context";
import { Ionicons } from "@expo/vector-icons";
import { useNavigation, useRoute, RouteProp } from "@react-navigation/native";
import { NativeStackNavigationProp } from "@react-navigation/native-stack";

import { useApp } from "../context/AppContext";
import { Button, Input, Card } from "../components";
import { colors, spacing, typography } from "../theme";
import { RootStackParamList } from "../navigation/types";

type BioimpedanceFormScreenNavigationProp =
  NativeStackNavigationProp<RootStackParamList>;
type BioimpedanceFormScreenRouteProp = RouteProp<
  RootStackParamList,
  "BioimpedanceForm"
>;

export const BioimpedanceFormScreen: React.FC = () => {
  const navigation = useNavigation<BioimpedanceFormScreenNavigationProp>();
  const route = useRoute<BioimpedanceFormScreenRouteProp>();
  const { bioimpedanceHistory, addBioimpedance, updateBioimpedance } = useApp();

  const editingData = route.params?.bioimpedanceId
    ? bioimpedanceHistory.find((b) => b.id === route.params.bioimpedanceId)
    : undefined;

  // Form state
  const [weight, setWeight] = useState(editingData?.weight?.toString() || "");
  const [bodyFatPercentage, setBodyFatPercentage] = useState(
    editingData?.bodyFatPercentage?.toString() || "",
  );
  const [muscleMass, setMuscleMass] = useState(
    editingData?.muscleMass?.toString() || "",
  );
  const [boneMass, setBoneMass] = useState(
    editingData?.boneMass?.toString() || "",
  );
  const [waterPercentage, setWaterPercentage] = useState(
    editingData?.waterPercentage?.toString() || "",
  );
  const [visceralFat, setVisceralFat] = useState(
    editingData?.visceralFat?.toString() || "",
  );
  const [bmr, setBmr] = useState(editingData?.bmr?.toString() || "");
  const [metabolicAge, setMetabolicAge] = useState(
    editingData?.metabolicAge?.toString() || "",
  );
  const [notes, setNotes] = useState(editingData?.notes || "");

  const [isSubmitting, setIsSubmitting] = useState(false);

  const handleSubmit = async () => {
    // Validation
    if (!weight) {
      Alert.alert("Erro", "Digite o peso");
      return;
    }
    if (!bodyFatPercentage) {
      Alert.alert("Erro", "Digite o percentual de gordura");
      return;
    }
    if (!muscleMass) {
      Alert.alert("Erro", "Digite a massa muscular");
      return;
    }

    setIsSubmitting(true);

    try {
      const data = {
        date: editingData?.date || new Date().toISOString(),
        weight: parseFloat(weight),
        bodyFatPercentage: parseFloat(bodyFatPercentage),
        muscleMass: parseFloat(muscleMass),
        boneMass: parseFloat(boneMass) || 0,
        waterPercentage: parseFloat(waterPercentage) || 0,
        visceralFat: parseFloat(visceralFat) || 0,
        bmr: parseFloat(bmr) || 0,
        metabolicAge: parseInt(metabolicAge) || 0,
        notes: notes.trim() || undefined,
      };

      if (editingData) {
        await updateBioimpedance({ ...data, id: editingData.id });
      } else {
        await addBioimpedance(data);
      }

      navigation.goBack();
    } catch (error) {
      Alert.alert("Erro", "Não foi possível salvar os dados");
    } finally {
      setIsSubmitting(false);
    }
  };

  const InfoCard: React.FC<{
    icon: string;
    title: string;
    description: string;
  }> = ({ icon, title, description }) => (
    <View style={styles.infoCard}>
      <Ionicons name={icon as any} size={20} color={colors.primary} />
      <View style={styles.infoContent}>
        <Text style={styles.infoTitle}>{title}</Text>
        <Text style={styles.infoDescription}>{description}</Text>
      </View>
    </View>
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
        <Text style={styles.title}>
          {editingData ? "Editar Bioimpedância" : "Nova Avaliação"}
        </Text>
        <View style={{ width: 40 }} />
      </View>

      <KeyboardAvoidingView
        behavior={Platform.OS === "ios" ? "padding" : "height"}
        style={{ flex: 1 }}
      >
        <ScrollView
          style={styles.content}
          contentContainerStyle={styles.contentContainer}
          showsVerticalScrollIndicator={false}
        >
          {/* Info Section */}
          <Card style={styles.infoSection}>
            <Text style={styles.infoSectionTitle}>
              <Ionicons
                name="information-circle"
                size={18}
                color={colors.info}
              />{" "}
              O que é Bioimpedância?
            </Text>
            <Text style={styles.infoSectionText}>
              A bioimpedância é um exame que mede a composição corporal através
              de uma corrente elétrica de baixa intensidade. Registre seus
              resultados para acompanhar sua evolução.
            </Text>
          </Card>

          {/* Main Metrics */}
          <Text style={styles.sectionTitle}>Métricas Principais</Text>

          <Input
            label="Peso (kg)"
            placeholder="Ex: 75.5"
            value={weight}
            onChangeText={setWeight}
            keyboardType="decimal-pad"
            leftIcon="scale-outline"
            required
          />

          <View style={styles.row}>
            <View style={styles.halfWidth}>
              <Input
                label="Gordura Corporal (%)"
                placeholder="Ex: 18.5"
                value={bodyFatPercentage}
                onChangeText={setBodyFatPercentage}
                keyboardType="decimal-pad"
                required
              />
            </View>
            <View style={styles.halfWidth}>
              <Input
                label="Massa Muscular (kg)"
                placeholder="Ex: 35.2"
                value={muscleMass}
                onChangeText={setMuscleMass}
                keyboardType="decimal-pad"
                required
              />
            </View>
          </View>

          {/* Additional Metrics */}
          <Text style={styles.sectionTitle}>Métricas Adicionais</Text>

          <View style={styles.row}>
            <View style={styles.halfWidth}>
              <Input
                label="Massa Óssea (kg)"
                placeholder="Ex: 3.2"
                value={boneMass}
                onChangeText={setBoneMass}
                keyboardType="decimal-pad"
              />
            </View>
            <View style={styles.halfWidth}>
              <Input
                label="Água Corporal (%)"
                placeholder="Ex: 55.0"
                value={waterPercentage}
                onChangeText={setWaterPercentage}
                keyboardType="decimal-pad"
              />
            </View>
          </View>

          <View style={styles.row}>
            <View style={styles.halfWidth}>
              <Input
                label="Gordura Visceral"
                placeholder="Nível (1-12)"
                value={visceralFat}
                onChangeText={setVisceralFat}
                keyboardType="numeric"
                hint="Ideal: 1-9"
              />
            </View>
            <View style={styles.halfWidth}>
              <Input
                label="Taxa Metabólica (kcal)"
                placeholder="Ex: 1800"
                value={bmr}
                onChangeText={setBmr}
                keyboardType="numeric"
              />
            </View>
          </View>

          <Input
            label="Idade Metabólica"
            placeholder="Ex: 28"
            value={metabolicAge}
            onChangeText={setMetabolicAge}
            keyboardType="numeric"
            hint="Idade que seu metabolismo aparenta ter"
          />

          {/* Notes */}
          <Text style={styles.sectionTitle}>Observações</Text>

          <Input
            label="Notas (opcional)"
            placeholder="Anote detalhes sobre esta avaliação, como jejum, horário, etc."
            value={notes}
            onChangeText={setNotes}
            multiline
            numberOfLines={3}
          />

          {/* Reference Values */}
          <Text style={styles.sectionTitle}>Valores de Referência</Text>

          <Card style={styles.referenceCard}>
            <InfoCard
              icon="fitness-outline"
              title="Gordura Corporal"
              description="Homens: 10-20% | Mulheres: 18-28%"
            />
            <InfoCard
              icon="water-outline"
              title="Água Corporal"
              description="Ideal: 50-65%"
            />
            <InfoCard
              icon="body-outline"
              title="Gordura Visceral"
              description="Saudável: 1-9 | Elevado: 10-14 | Alto: 15+"
            />
          </Card>

          <View style={{ height: spacing.xl }} />
        </ScrollView>

        {/* Submit Button */}
        <View style={styles.footer}>
          <Button
            title={editingData ? "Salvar Alterações" : "Registrar Avaliação"}
            onPress={handleSubmit}
            loading={isSubmitting}
            fullWidth
          />
        </View>
      </KeyboardAvoidingView>
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
  content: {
    flex: 1,
  },
  contentContainer: {
    padding: spacing.lg,
  },
  infoSection: {
    backgroundColor: `${colors.info}10`,
    padding: spacing.lg,
    marginBottom: spacing.lg,
  },
  infoSectionTitle: {
    fontSize: typography.md,
    fontWeight: typography.semibold,
    color: colors.info,
    marginBottom: spacing.sm,
  },
  infoSectionText: {
    fontSize: typography.sm,
    color: colors.textSecondary,
    lineHeight: 20,
  },
  sectionTitle: {
    fontSize: typography.md,
    fontWeight: typography.bold,
    color: colors.textPrimary,
    marginTop: spacing.md,
    marginBottom: spacing.md,
  },
  row: {
    flexDirection: "row",
    gap: spacing.md,
  },
  halfWidth: {
    flex: 1,
  },
  referenceCard: {
    padding: spacing.lg,
    gap: spacing.md,
  },
  infoCard: {
    flexDirection: "row",
    alignItems: "flex-start",
    gap: spacing.sm,
  },
  infoContent: {
    flex: 1,
  },
  infoTitle: {
    fontSize: typography.sm,
    fontWeight: typography.semibold,
    color: colors.textPrimary,
  },
  infoDescription: {
    fontSize: typography.xs,
    color: colors.textSecondary,
    marginTop: 2,
  },
  footer: {
    padding: spacing.lg,
    backgroundColor: colors.surface,
    borderTopWidth: 1,
    borderTopColor: colors.border,
  },
});
