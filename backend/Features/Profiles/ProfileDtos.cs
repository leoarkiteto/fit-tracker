namespace FitTracker.Api.Features.Profiles;

public record UserProfileDto(
    Guid Id,
    string Name,
    int? Age,
    double? Height,
    double? CurrentWeight,
    double? GoalWeight,
    string? AvatarUrl,
    // AI Planning fields
    ExperienceLevel ExperienceLevel,
    int AvailableDaysPerWeek,
    int? PreferredWorkoutDuration,
    EquipmentType EquipmentType
);

public record CreateUserProfileRequest(
    string Name,
    int? Age,
    double? Height,
    double? CurrentWeight,
    double? GoalWeight,
    string? AvatarUrl,
    // AI Planning fields (optional with defaults)
    ExperienceLevel? ExperienceLevel = null,
    int? AvailableDaysPerWeek = null,
    int? PreferredWorkoutDuration = null,
    EquipmentType? EquipmentType = null
);

public record UpdateUserProfileRequest(
    string Name,
    int? Age,
    double? Height,
    double? CurrentWeight,
    double? GoalWeight,
    string? AvatarUrl,
    // AI Planning fields (optional)
    ExperienceLevel? ExperienceLevel = null,
    int? AvailableDaysPerWeek = null,
    int? PreferredWorkoutDuration = null,
    EquipmentType? EquipmentType = null
);
