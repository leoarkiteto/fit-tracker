"use client";

import React from "react";
import type { CalendarDayActivity } from "@/lib/calendar-utils";

interface DayTooltipProps {
  activities: CalendarDayActivity[];
  className?: string;
}

export function DayTooltip({ activities, className = "" }: DayTooltipProps) {
  return (
    <div
      className={`bg-surface-900 text-white text-sm rounded-lg shadow-lg px-3 py-2 max-w-[200px] ${className}`}
      role="tooltip"
    >
      {activities.length === 0 ? (
        <p className="text-surface-200">Nenhuma atividade</p>
      ) : (
        <ul className="space-y-2">
          {activities.map((a) => (
            <li key={a.workoutId}>
              <span className="font-medium">{a.name}</span>
              {a.exerciseNames?.length > 0 ? (
                <ul className="list-disc list-inside text-surface-200 text-xs mt-0.5 ml-1">
                  {a.exerciseNames.map((ex, i) => (
                    <li key={i}>{ex}</li>
                  ))}
                </ul>
              ) : null}
            </li>
          ))}
        </ul>
      )}
    </div>
  );
}
