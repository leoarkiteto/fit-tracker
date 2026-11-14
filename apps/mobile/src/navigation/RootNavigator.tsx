import React from "react";
import { createNativeStackNavigator } from "@react-navigation/native-stack";

import { TabNavigator } from "./TabNavigator";
import {
  WorkoutFormScreen,
  WorkoutDetailScreen,
  ExerciseSessionScreen,
  BioimpedanceFormScreen,
  BioimpedanceHistoryScreen,
  WaterScreen,
  ChangePasswordScreen,
  WorkoutPlanningScreen,
} from "../screens";
import { RootStackParamList } from "./types";

const Stack = createNativeStackNavigator<RootStackParamList>();

export const RootNavigator: React.FC = () => {
  return (
    <Stack.Navigator
      screenOptions={{
        headerShown: false,
        animation: "slide_from_right",
      }}
    >
      <Stack.Screen name="MainTabs" component={TabNavigator} />
      <Stack.Screen
        name="WorkoutForm"
        component={WorkoutFormScreen}
        options={{
          animation: "slide_from_bottom",
          presentation: "modal",
        }}
      />
      <Stack.Screen name="WorkoutDetail" component={WorkoutDetailScreen} />
      <Stack.Screen
        name="ExerciseSession"
        component={ExerciseSessionScreen}
        options={{
          animation: "fade",
          gestureEnabled: false,
        }}
      />
      <Stack.Screen
        name="BioimpedanceForm"
        component={BioimpedanceFormScreen}
        options={{
          animation: "slide_from_bottom",
          presentation: "modal",
        }}
      />
      <Stack.Screen
        name="BioimpedanceHistory"
        component={BioimpedanceHistoryScreen}
      />
      <Stack.Screen name="Water" component={WaterScreen} />
      <Stack.Screen
        name="ChangePassword"
        component={ChangePasswordScreen}
      />
      <Stack.Screen
        name="WorkoutPlanning"
        component={WorkoutPlanningScreen}
        options={{
          animation: "slide_from_bottom",
          presentation: "modal",
        }}
      />
    </Stack.Navigator>
  );
};
