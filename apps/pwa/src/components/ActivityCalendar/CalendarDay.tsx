"use client";

import React from "react";
import type { CalendarDayState } from "@/lib/calendar-utils";

interface CalendarDayProps {
  date: string;
  state: CalendarDayState;
  dayOfMonth: number;
  isCurrentMonth: boolean;
  isActive?: boolean;
  onSelect?: () => void;
  onFocus?: () => void;
  onBlur?: () => void;
}

const stateStyles: Record<CalendarDayState, string> = {
  completed: "bg-emerald-100 text-emerald-800 border-emerald-400",
  today: "border-primary-500 bg-primary-100 text-primary-700 font-semibold",
  default: "border-surface-300 bg-surface-100 text-surface-700",
};

export function CalendarDay({
  date,
  state,
  dayOfMonth,
  isCurrentMonth,
  isActive,
  onSelect,
  onFocus,
  onBlur,
}: CalendarDayProps) {
  const isEmpty = dayOfMonth === 0 || !date;
  if (isEmpty) {
    return <div className="w-9 h-9 flex items-center justify-center" aria-hidden />;
  }

  const base =
    "w-9 h-9 flex items-center justify-center rounded-full border text-sm transition-colors cursor-pointer focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-1";
  const opacity = isCurrentMonth ? "opacity-100" : "opacity-40";
  const style = stateStyles[state];

  return (
    <button
      type="button"
      data-date={date}
      className={`${base} ${style} ${opacity} ${isActive ? "ring-2 ring-primary-500 ring-offset-1" : ""}`}
      onClick={onSelect}
      onFocus={onFocus}
      onBlur={onBlur}
      tabIndex={0}
      aria-label={`Dia ${dayOfMonth}${state === "today" ? ", hoje" : ""}`}
    >
      {dayOfMonth}
    </button>
  );
}
