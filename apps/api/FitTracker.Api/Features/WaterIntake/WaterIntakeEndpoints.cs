namespace FitTracker.Api.Features.WaterIntake;

public static class WaterIntakeEndpoints
{
    private const int DefaultGoalMl = 2000;
    private const int MlPerKg = 35;
    private const int MaxAmountMlPerEntry = 2000;

    public static void MapWaterIntakeEndpoints(this WebApplication app)
    {
        var group = app
            .MapGroup("/api/profiles/{profileId:guid}/water")
            .WithTags("WaterIntake");

        // GET /api/profiles/{profileId}/water?date=yyyy-MM-dd - Consumo do dia (total, meta, entradas)
        group
            .MapGet(
                "/",
                async (
                    Guid profileId,
                    [FromQuery] string? date,
                    FitTrackerDbContext db
                ) =>
                {
                    var profile = await db.UserProfiles.FindAsync(profileId);
                    if (profile is null)
                        return Results.NotFound("Profile not found");

                    var dateUtc = ParseDateOrDefault(date);
                    var start = dateUtc.Date;
                    var end = start.AddDays(1);

                    var entries = await db
                        .WaterIntakeEntries
                        .Where(e =>
                            e.UserProfileId == profileId
                            && e.ConsumedAt >= start
                            && e.ConsumedAt < end
                        )
                        .OrderBy(e => e.ConsumedAt)
                        .ToListAsync();

                    var totalMl = entries.Sum(e => e.AmountMl);
                    var goalMl = profile.CurrentWeight.HasValue
                        ? (int)(profile.CurrentWeight.Value * MlPerKg)
                        : DefaultGoalMl;

                    var summary = new DailyWaterSummaryDto(
                        Date: start.ToString("yyyy-MM-dd"),
                        TotalMl: totalMl,
                        GoalMl: goalMl,
                        Entries: entries.Select(ToEntryDto).ToList()
                    );

                    return Results.Ok(summary);
                }
            )
            .WithName("GetDailyWater")
            .Produces<DailyWaterSummaryDto>(StatusCodes.Status200OK)
            .Produces(StatusCodes.Status404NotFound);

        // POST /api/profiles/{profileId}/water - Adicionar consumo
        group
            .MapPost(
                "/",
                async (
                    Guid profileId,
                    CreateWaterIntakeRequest request,
                    FitTrackerDbContext db
                ) =>
                {
                    var profileExists = await db.UserProfiles.AnyAsync(p => p.Id == profileId);
                    if (!profileExists)
                        return Results.NotFound("Profile not found");

                    if (request.AmountMl <= 0 || request.AmountMl > MaxAmountMlPerEntry)
                        return Results.BadRequest(
                            $"AmountMl must be between 1 and {MaxAmountMlPerEntry}."
                        );

                    var consumedAt = request.ConsumedAt ?? DateTime.UtcNow;
                    var entry = new WaterIntakeEntry
                    {
                        Id = Guid.NewGuid(),
                        UserProfileId = profileId,
                        AmountMl = request.AmountMl,
                        ConsumedAt = consumedAt,
                    };

                    db.WaterIntakeEntries.Add(entry);
                    await db.SaveChangesAsync();

                    return Results.Created(
                        $"/api/profiles/{profileId}/water/{entry.Id}",
                        ToEntryDto(entry)
                    );
                }
            )
            .WithName("CreateWaterIntake")
            .Produces<WaterIntakeEntryDto>(StatusCodes.Status201Created)
            .Produces(StatusCodes.Status400BadRequest)
            .Produces(StatusCodes.Status404NotFound);

        // DELETE /api/profiles/{profileId}/water/{entryId} - Remover entrada
        group
            .MapDelete(
                "/{entryId:guid}",
                async (Guid profileId, Guid entryId, FitTrackerDbContext db) =>
                {
                    var entry = await db.WaterIntakeEntries.FirstOrDefaultAsync(e =>
                        e.Id == entryId && e.UserProfileId == profileId
                    );

                    if (entry is null)
                        return Results.NotFound();

                    db.WaterIntakeEntries.Remove(entry);
                    await db.SaveChangesAsync();

                    return Results.NoContent();
                }
            )
            .WithName("DeleteWaterIntake")
            .Produces(StatusCodes.Status204NoContent)
            .Produces(StatusCodes.Status404NotFound);
    }

    private static DateTime ParseDateOrDefault(string? date)
    {
        if (
            !string.IsNullOrWhiteSpace(date)
            && DateTime.TryParse(date, null, System.Globalization.DateTimeStyles.AssumeUniversal, out var parsed)
        )
            return parsed.Kind == DateTimeKind.Unspecified
                ? DateTime.SpecifyKind(parsed, DateTimeKind.Utc)
                : parsed.ToUniversalTime();
        return DateTime.UtcNow;
    }

    private static WaterIntakeEntryDto ToEntryDto(WaterIntakeEntry e) =>
        new(e.Id, e.AmountMl, e.ConsumedAt, e.Note);
}
