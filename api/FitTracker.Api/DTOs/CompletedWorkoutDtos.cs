namespace FitTracker.Api.DTOs;

public record CompletedWorkoutDto(
    Guid Id,
    Guid WorkoutId,
    DateTime CompletedAt,
    int DurationSeconds
);

public record CompleteWorkoutRequest(Guid WorkoutId, int DurationSeconds);

public record WorkoutStatsDto(
    int TotalWorkoutsCompleted,
    int WorkoutsThisWeek,
    int TotalMinutesSpent
);
