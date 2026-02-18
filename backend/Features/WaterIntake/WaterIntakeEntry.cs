namespace FitTracker.Api.Features.WaterIntake;

public class WaterIntakeEntry
{
    public Guid Id { get; set; }
    public Guid UserProfileId { get; set; }
    public int AmountMl { get; set; }
    public DateTime ConsumedAt { get; set; }
    public string? Note { get; set; }
}
