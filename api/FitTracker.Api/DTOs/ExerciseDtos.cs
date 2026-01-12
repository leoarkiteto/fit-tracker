namespace FitTracker.Api.DTOs;

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
