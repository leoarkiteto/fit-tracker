namespace FitTracker.Api.Features.WaterIntake;

public record WaterIntakeEntryDto(Guid Id, int AmountMl, DateTime ConsumedAt, string? Note);

public record CreateWaterIntakeRequest(int AmountMl, DateTime? ConsumedAt);

public record DailyWaterSummaryDto(
    string Date,
    int TotalMl,
    int GoalMl,
    IReadOnlyList<WaterIntakeEntryDto> Entries
);
