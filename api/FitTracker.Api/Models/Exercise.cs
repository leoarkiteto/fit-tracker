namespace FitTracker.Api.Models;

public class Exercise
{
    public Guid Id { get; set; }
    public string Name { get; set; } = string.Empty;
    public MuscleGroup MuscleGroup { get; set; }
    public int Sets { get; set; }
    public int Reps { get; set; }
    public double? Weight { get; set; } // null = peso corporal
    public int RestSeconds { get; set; }
    public string? Notes { get; set; }

    // Foreign key para Workout
    public Guid WorkoutId { get; set; }
    public Workout? Workout { get; set; }
}

public enum MuscleGroup
{
    Chest,
    Back,
    Shoulders,
    Biceps,
    Triceps,
    Legs,
    Glutes,
    Abs,
    Cardio,
}
