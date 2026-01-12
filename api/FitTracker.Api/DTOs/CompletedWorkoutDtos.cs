namespace FitTracker.Api.DTOs;

public record CompletedWorkoutDto(Guid Id, Guid WorkoutId, DateTime CompletedAt);

public record CompleteWorkoutRequest(Guid WorkoutId);

public record WorkoutStatsDto(
    int TotalWorkoutsCompleted,
    int WorkoutsThisWeek,
    int TotalMinutesSpent
);
