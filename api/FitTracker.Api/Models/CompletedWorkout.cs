namespace FitTracker.Api.Models;

public class CompletedWorkout
{
    public Guid Id { get; set; }
    public Guid WorkoutId { get; set; }
    public Workout? Workout { get; set; }
    public DateTime CompletedAt { get; set; } = DateTime.UtcNow;

    // Foreign key para UserProfile
    public Guid UserProfileId { get; set; }
    public UserProfile? UserProfile { get; set; }
}
