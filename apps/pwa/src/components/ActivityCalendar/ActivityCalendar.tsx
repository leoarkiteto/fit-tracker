"use client";

import React, { useMemo, useState, useRef, useEffect } from "react";
import type { Workout, CompletedWorkoutRecord } from "@/lib/types";
import { buildMonthView, getTodayISODate } from "@/lib/calendar-utils";
import { CalendarDay } from "./CalendarDay";
import { DayTooltip } from "./DayTooltip";
import { ChevronLeft, ChevronRight, Loader2 } from "lucide-react";

interface ActivityCalendarProps {
  workouts: Workout[];
  completedRecords: CompletedWorkoutRecord[];
  loading?: boolean;
  error?: string | null;
  onRetry?: () => void;
}

export function ActivityCalendar({
  workouts,
  completedRecords,
  loading = false,
  error = null,
  onRetry,
}: ActivityCalendarProps) {
  const today = getTodayISODate();
  const [viewDate, setViewDate] = useState(() => {
    const d = new Date();
    return { year: d.getFullYear(), month: d.getMonth() + 1 };
  });
  const [activeDay, setActiveDay] = useState<string | null>(null);
  const [tooltipRect, setTooltipRect] = useState<{ left: number; top: number } | null>(null);
  const containerRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const handleClickOutside = (e: MouseEvent) => {
      if (containerRef.current && !containerRef.current.contains(e.target as Node)) {
        setActiveDay(null);
        setTooltipRect(null);
      }
    };
    document.addEventListener("mousedown", handleClickOutside);
    return () => document.removeEventListener("mousedown", handleClickOutside);
  }, []);

  const monthView = useMemo(() => {
    return buildMonthView(
      viewDate.year,
      viewDate.month,
      workouts,
      completedRecords,
      today
    );
  }, [viewDate.year, viewDate.month, workouts, completedRecords, today]);

  const goPrev = () => {
    setViewDate((prev) => {
      if (prev.month === 1) return { year: prev.year - 1, month: 12 };
      return { year: prev.year, month: prev.month - 1 };
    });
    setActiveDay(null);
    setTooltipRect(null);
  };

  const goNext = () => {
    setViewDate((prev) => {
      if (prev.month === 12) return { year: prev.year + 1, month: 1 };
      return { year: prev.year, month: prev.month + 1 };
    });
    setActiveDay(null);
    setTooltipRect(null);
  };

  const handleDayHover = (date: string, el: HTMLElement | null) => {
    if (!date) return;
    setActiveDay(date);
    if (el) {
      const rect = el.getBoundingClientRect();
      setTooltipRect({ left: rect.left, top: rect.bottom + 8 });
    } else {
      setTooltipRect(null);
    }
  };

  if (loading) {
    return (
      <section className="mb-8">
        <h2 className="text-xl font-semibold text-surface-900 mb-4">
          Calendário de Treinos
        </h2>
        <div className="card p-8 flex items-center justify-center min-h-[280px]">
          <Loader2 className="w-8 h-8 text-primary-500 animate-spin" />
        </div>
      </section>
    );
  }

  if (error) {
    return (
      <section className="mb-8">
        <h2 className="text-xl font-semibold text-surface-900 mb-4">
          Calendário de Treinos
        </h2>
        <div className="card p-8 text-center">
          <p className="text-error mb-4">{error}</p>
          {onRetry && (
            <button
              type="button"
              onClick={onRetry}
              className="btn-primary"
            >
              Tentar novamente
            </button>
          )}
        </div>
      </section>
    );
  }

  const activeDayData = activeDay ? monthView.days.find((d) => d.date === activeDay) : null;

  return (
    <section
      className="mb-8"
      ref={containerRef}
      onMouseLeave={() => { setActiveDay(null); setTooltipRect(null); }}
    >
      <h2 className="text-xl font-semibold text-surface-900 mb-4">
        Calendário de Treinos
      </h2>
      <div className="card p-5 relative">
        <div className="flex items-center justify-between mb-4">
          <button
            type="button"
            onClick={goPrev}
            className="p-2 rounded-lg hover:bg-surface-200 text-surface-700 focus:outline-none focus:ring-2 focus:ring-primary-500"
            aria-label="Mês anterior"
          >
            <ChevronLeft className="w-5 h-5" />
          </button>
          <span className="font-semibold text-surface-900">
            {monthView.monthLabel} {monthView.year}
          </span>
          <button
            type="button"
            onClick={goNext}
            className="p-2 rounded-lg hover:bg-surface-200 text-surface-700 focus:outline-none focus:ring-2 focus:ring-primary-500"
            aria-label="Próximo mês"
          >
            <ChevronRight className="w-5 h-5" />
          </button>
        </div>
        <div className="grid grid-cols-7 gap-1 mb-2">
          {monthView.dayOfWeekLabels.map((label) => (
            <div
              key={label}
              className="text-center text-sm text-surface-600 font-medium"
            >
              {label}
            </div>
          ))}
        </div>
        <div className="grid grid-cols-7 gap-1">
          {monthView.days.map((day, index) => (
            <div
              key={day.date || `empty-${index}`}
              className="flex justify-center relative"
              onMouseEnter={(e) => day.date && handleDayHover(day.date, e.currentTarget)}
              onMouseLeave={() => { setActiveDay(null); setTooltipRect(null); }}
            >
              <CalendarDay
                date={day.date}
                state={day.state}
                dayOfMonth={day.dayOfMonth}
                isCurrentMonth={day.isCurrentMonth}
                isActive={activeDay === day.date}
                onSelect={() => {
                  if (!day.date) return;
                  setActiveDay((prev) => (prev === day.date ? null : day.date));
                  setTooltipRect(null);
                }}
                onFocus={() => day.date && setActiveDay(day.date)}
                onBlur={() => { setActiveDay(null); setTooltipRect(null); }}
              />
            </div>
          ))}
        </div>
        {activeDayData && (
          <div
            className={`z-50 pointer-events-none ${tooltipRect ? "fixed" : "mt-4"}`}
            style={tooltipRect ? { left: tooltipRect.left, top: tooltipRect.top } : undefined}
          >
            <DayTooltip activities={activeDayData.activities} />
          </div>
        )}
      </div>
    </section>
  );
}
