namespace FitTracker.Api.Features.Workouts;

// Workout DTOs
public record WorkoutDto(
    Guid Id,
    string Name,
    string? Description,
    string Goal,
    List<string> Days,
    List<ExerciseDto> Exercises,
    DateTime CreatedAt,
    DateTime UpdatedAt,
    bool IsCompleted,
    DateTime? CompletedAt
);

public record CreateWorkoutRequest(
    string Name,
    string? Description,
    string Goal,
    List<string> Days,
    List<CreateExerciseRequest> Exercises
);

public record UpdateWorkoutRequest(
    string Name,
    string? Description,
    string Goal,
    List<string> Days,
    List<UpdateExerciseRequest> Exercises
);

// Exercise DTOs
public record ExerciseDto(
    Guid Id,
    string Name,
    string MuscleGroup,
    int Sets,
    int Reps,
    double? Weight,
    int RestSeconds,
    string? Notes
);

public record CreateExerciseRequest(
    string Name,
    string MuscleGroup,
    int Sets,
    int Reps,
    double? Weight,
    int RestSeconds,
    string? Notes
);

public record UpdateExerciseRequest(
    Guid? Id, // null = novo exercício
    string Name,
    string MuscleGroup,
    int Sets,
    int Reps,
    double? Weight,
    int RestSeconds,
    string? Notes
);

// Completed Workout DTOs
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
