namespace FitTracker.Api.Features.Profiles;

public class UserProfile
{
    public Guid Id { get; set; }
    public string Name { get; set; } = "John Doe";
    public int? Age { get; set; }
    public double? Height { get; set; } // cm
    public double? CurrentWeight { get; set; } // kg
    public double? GoalWeight { get; set; } // kg
    public string? AvatarUrl { get; set; }
    public DateTime CreatedAt { get; set; } = DateTime.UtcNow;
    public DateTime UpdatedAt { get; set; } = DateTime.UtcNow;

    // AI Planning fields
    public ExperienceLevel ExperienceLevel { get; set; } = ExperienceLevel.Beginner;
    public int AvailableDaysPerWeek { get; set; } = 3;
    public int? PreferredWorkoutDuration { get; set; } // minutes
    public EquipmentType EquipmentType { get; set; } = EquipmentType.Gym;
}

public enum ExperienceLevel
{
    Beginner,
    Intermediate,
    Advanced,
}

public enum EquipmentType
{
    Gym, // Full gym equipment
    Home, // Basic home equipment (dumbbells, bands)
    Minimal, // Bodyweight only
}
