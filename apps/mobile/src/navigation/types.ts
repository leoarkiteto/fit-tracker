import { NavigatorScreenParams } from "@react-navigation/native";

export type AuthStackParamList = {
  SignIn: undefined;
  SignUp: undefined;
};

export type TabParamList = {
  Home: undefined;
  Workouts: undefined;
  Profile: undefined;
};

export type RootStackParamList = {
  MainTabs: NavigatorScreenParams<TabParamList>;
  WorkoutForm: { workoutId?: string };
  WorkoutDetail: { workoutId: string };
  ExerciseSession: { workoutId: string; exerciseId: string };
  BioimpedanceForm: { bioimpedanceId?: string };
  BioimpedanceHistory: undefined;
  Profile: undefined;
  ChangePassword: undefined;
  WorkoutPlanning: undefined;
};

declare global {
  namespace ReactNavigation {
    interface RootParamList extends RootStackParamList {}
  }
}
