namespace FitTracker.Api.DTOs;

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
