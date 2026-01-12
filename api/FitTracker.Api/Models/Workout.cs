namespace FitTracker.Api.Models;

public class Workout
{
    public Guid Id { get; set; }
    public string Name { get; set; } = string.Empty;
    public string? Description { get; set; }
    public WorkoutGoal Goal { get; set; }
    public List<DayOfWeekEnum> Days { get; set; } = new();
    public List<Exercise> Exercises { get; set; } = new();
    public DateTime CreatedAt { get; set; } = DateTime.UtcNow;
    public DateTime UpdatedAt { get; set; } = DateTime.UtcNow;
    public bool IsCompleted { get; set; }
    public DateTime? CompletedAt { get; set; }

    // Foreign key para UserProfile
    public Guid UserProfileId { get; set; }
    public UserProfile? UserProfile { get; set; }
}

public enum WorkoutGoal
{
    Hypertrophy,
    Strength,
    Endurance,
    WeightLoss,
    Maintenance,
}

public enum DayOfWeekEnum
{
    Monday,
    Tuesday,
    Wednesday,
    Thursday,
    Friday,
    Saturday,
    Sunday,
}
