namespace FitTracker.Api.Features.Profiles;

public record UserProfileDto(
    Guid Id,
    string Name,
    int? Age,
    double? Height,
    double? CurrentWeight,
    double? GoalWeight,
    string? AvatarUrl
);

public record CreateUserProfileRequest(
    string Name,
    int? Age,
    double? Height,
    double? CurrentWeight,
    double? GoalWeight,
    string? AvatarUrl
);

public record UpdateUserProfileRequest(
    string Name,
    int? Age,
    double? Height,
    double? CurrentWeight,
    double? GoalWeight,
    string? AvatarUrl
);
