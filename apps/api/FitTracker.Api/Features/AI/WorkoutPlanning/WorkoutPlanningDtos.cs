namespace FitTracker.Api.Features.AI.WorkoutPlanning;

/// <summary>
/// Request to generate a new workout plan
/// </summary>
public record GeneratePlanRequest(
    Guid UserProfileId,
    WorkoutGoal Goal,
    int? OverrideDaysPerWeek = null,
    int? OverrideDurationMinutes = null,
    EquipmentType? OverrideEquipment = null,
    string? AdditionalNotes = null
);

/// <summary>
/// Response containing a generated workout plan
/// </summary>
public record GeneratedPlanResponse(
    Guid PlanId,
    string Summary,
    string Rationale,
    List<GeneratedWorkoutDto> Workouts,
    DateTime GeneratedAt
);

/// <summary>
/// A workout in the generated plan
/// </summary>
public record GeneratedWorkoutDto(
    string Name,
    string? Description,
    WorkoutGoal Goal,
    List<DayOfWeekEnum> Days,
    List<GeneratedExerciseDto> Exercises,
    int EstimatedDurationMinutes
);

/// <summary>
/// An exercise in the generated workout
/// </summary>
public record GeneratedExerciseDto(
    string Name,
    MuscleGroup MuscleGroup,
    int Sets,
    int Reps,
    double? SuggestedWeight,
    int RestSeconds,
    string? Notes
);

/// <summary>
/// Request to accept and save a generated plan
/// </summary>
public record AcceptPlanRequest(
    Guid PlanId,
    Guid UserProfileId,
    List<GeneratedWorkoutDto> Workouts
);

/// <summary>
/// Response after accepting a plan
/// </summary>
public record AcceptPlanResponse(List<Guid> CreatedWorkoutIds, string Message);

/// <summary>
/// Context data provided to the AI for planning
/// </summary>
public record UserPlanningContext(
    Guid UserProfileId,
    string Name,
    int? Age,
    double? CurrentWeight,
    double? GoalWeight,
    ExperienceLevel ExperienceLevel,
    int AvailableDaysPerWeek,
    int? PreferredWorkoutDuration,
    EquipmentType EquipmentType,
    WorkoutGoal Goal,
    // Historical data
    int TotalCompletedWorkouts,
    int WorkoutsThisMonth,
    double? LatestBodyFatPercentage,
    double? LatestMuscleMass,
    List<string> CurrentExercises
);
