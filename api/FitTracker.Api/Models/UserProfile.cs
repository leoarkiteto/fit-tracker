namespace FitTracker.Api.Models;

public class UserProfile
{
    public Guid Id { get; set; }
    public string Name { get; set; } = "Atleta";
    public int? Age { get; set; }
    public double? Height { get; set; } // cm
    public double? CurrentWeight { get; set; } // kg
    public double? GoalWeight { get; set; } // kg
    public string? AvatarUrl { get; set; }
    public DateTime CreatedAt { get; set; } = DateTime.UtcNow;
    public DateTime UpdatedAt { get; set; } = DateTime.UtcNow;
}
