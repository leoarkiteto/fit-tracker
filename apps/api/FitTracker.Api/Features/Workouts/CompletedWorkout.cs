namespace FitTracker.Api.Features.Workouts;

public class CompletedWorkout
{
    public Guid Id { get; set; }
    public Guid WorkoutId { get; set; }
    public Workout? Workout { get; set; }
    public DateTime CompletedAt { get; set; } = DateTime.UtcNow;
    public int DurationSeconds { get; set; }

    // Foreign key para UserProfile
    public Guid UserProfileId { get; set; }
}
