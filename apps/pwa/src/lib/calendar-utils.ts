import type { Workout, CompletedWorkoutRecord } from "@/lib/types";

export type CalendarDayState = "completed" | "today" | "default";

export interface CalendarDayActivity {
  name: string;
  workoutId: string;
  exerciseNames: string[];
}

export interface CalendarDay {
  date: string; // YYYY-MM-DD
  state: CalendarDayState;
  activities: CalendarDayActivity[];
  dayOfMonth: number;
  isCurrentMonth: boolean;
}

export interface MonthView {
  year: number;
  month: number;
  monthLabel: string;
  dayOfWeekLabels: string[];
  days: CalendarDay[];
}

const PT_MONTHS: string[] = [
  "Janeiro", "Fevereiro", "Março", "Abril", "Maio", "Junho",
  "Julho", "Agosto", "Setembro", "Outubro", "Novembro", "Dezembro",
];

const PT_WEEKDAYS_SHORT: string[] = ["Dom", "Seg", "Ter", "Qua", "Qui", "Sex", "Sáb"];

const DAY_OF_WEEK_KEYS = ["sunday", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday"] as const;

/** Dia da semana de hoje (data local), no formato usado por Workout.days */
export function getTodayWeekdayKey(): (typeof DAY_OF_WEEK_KEYS)[number] {
  return DAY_OF_WEEK_KEYS[new Date().getDay()];
}

function toISODate(year: number, month: number, day: number): string {
  const m = String(month).padStart(2, "0");
  const d = String(day).padStart(2, "0");
  return `${year}-${m}-${d}`;
}

function getWeekdayKey(date: Date): (typeof DAY_OF_WEEK_KEYS)[number] {
  return DAY_OF_WEEK_KEYS[date.getDay()];
}

function isSameDay(a: string, b: string): boolean {
  return a === b;
}

export function buildMonthView(
  year: number,
  month: number,
  workouts: Workout[],
  completedRecords: CompletedWorkoutRecord[],
  todayISODate: string
): MonthView {
  const first = new Date(year, month - 1, 1);
  const last = new Date(year, month - 1 + 1, 0);
  const startPadding = first.getDay();
  const daysInMonth = last.getDate();

  const completedByDate = new Map<string, CompletedWorkoutRecord[]>();
  for (const r of completedRecords) {
    const date = r.completedAt.slice(0, 10);
    if (!completedByDate.has(date)) completedByDate.set(date, []);
    completedByDate.get(date)!.push(r);
  }

  const workoutById = new Map(workouts.map((w) => [w.id, w]));

  const days: CalendarDay[] = [];
  // Células vazias no início (antes do dia 1)
  for (let i = 0; i < startPadding; i++) {
    days.push({
      date: "",
      state: "default",
      activities: [],
      dayOfMonth: 0,
      isCurrentMonth: false,
    });
  }

  for (let d = 1; d <= daysInMonth; d++) {
    const date = toISODate(year, month, d);
    const completedOnDay = completedByDate.get(date) ?? [];
    const hasCompleted = completedOnDay.length > 0;
    const isToday = isSameDay(date, todayISODate);
    const state: CalendarDayState = isToday ? "today" : hasCompleted ? "completed" : "default";

    const dateObj = new Date(year, month - 1, d);
    const weekdayKey = getWeekdayKey(dateObj);
    const scheduled = workouts.filter((w) => w.days?.includes(weekdayKey as "monday"));
    const completedWorkoutIds = new Set(completedOnDay.map((r) => r.workoutId));
    const activityNames = new Set<string>();
    const activities: CalendarDayActivity[] = [];

    for (const w of scheduled) {
      if (!activityNames.has(w.name)) {
        activityNames.add(w.name);
        activities.push({
          name: w.name,
          workoutId: w.id,
          exerciseNames: (w.exercises ?? []).map((e) => e.name),
        });
      }
    }
    for (const r of completedOnDay) {
      const w = workoutById.get(r.workoutId);
      if (w && !activityNames.has(w.name)) {
        activityNames.add(w.name);
        activities.push({
          name: w.name,
          workoutId: w.id,
          exerciseNames: (w.exercises ?? []).map((e) => e.name),
        });
      }
    }

    days.push({
      date,
      state,
      activities,
      dayOfMonth: d,
      isCurrentMonth: true,
    });
  }

  // Células vazias no final (apenas dias do mês corrente são exibidos)
  const totalCells = 42;
  const remaining = totalCells - days.length;
  for (let i = 0; i < remaining; i++) {
    days.push({
      date: "",
      state: "default",
      activities: [],
      dayOfMonth: 0,
      isCurrentMonth: false,
    });
  }

  return {
    year,
    month,
    monthLabel: PT_MONTHS[month - 1] ?? String(month),
    dayOfWeekLabels: PT_WEEKDAYS_SHORT,
    days,
  };
}

export function getTodayISODate(): string {
  const d = new Date();
  const y = d.getFullYear();
  const m = String(d.getMonth() + 1).padStart(2, "0");
  const day = String(d.getDate()).padStart(2, "0");
  return `${y}-${m}-${day}`;
}
