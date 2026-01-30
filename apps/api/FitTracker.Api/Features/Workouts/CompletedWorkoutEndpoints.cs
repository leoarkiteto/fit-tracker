namespace FitTracker.Api.Features.Workouts;

public static class CompletedWorkoutEndpoints
{
    public static void MapCompletedWorkoutEndpoints(this WebApplication app)
    {
        var group = app.MapGroup("/api/profiles/{profileId:guid}/completed-workouts")
            .WithTags("Completed Workouts");

        // GET /api/profiles/{profileId}/completed-workouts - Lista todos os treinos concluídos
        group
            .MapGet(
                "/",
                async (Guid profileId, FitTrackerDbContext db) =>
                {
                    var completed = await db
                        .CompletedWorkouts.Where(c => c.UserProfileId == profileId)
                        .OrderByDescending(c => c.CompletedAt)
                        .ToListAsync();

                    return Results.Ok(completed.Select(ToDto));
                }
            )
            .WithName("GetCompletedWorkouts")
            .Produces<IEnumerable<CompletedWorkoutDto>>(StatusCodes.Status200OK);

        // POST /api/profiles/{profileId}/completed-workouts - Marca um treino como concluído
        group
            .MapPost(
                "/",
                async (Guid profileId, CompleteWorkoutRequest request, FitTrackerDbContext db) =>
                {
                    // Verificar se o perfil existe
                    var profileExists = await db.UserProfiles.AnyAsync(p => p.Id == profileId);
                    if (!profileExists)
                        return Results.NotFound("Profile not found");

                    // Verificar se o treino existe
                    var workoutExists = await db.Workouts.AnyAsync(w => w.Id == request.WorkoutId);
                    if (!workoutExists)
                        return Results.NotFound("Workout not found");

                    var completed = new CompletedWorkout
                    {
                        Id = Guid.NewGuid(),
                        WorkoutId = request.WorkoutId,
                        UserProfileId = profileId,
                        CompletedAt = DateTime.UtcNow,
                        DurationSeconds = request.DurationSeconds,
                    };

                    db.CompletedWorkouts.Add(completed);
                    await db.SaveChangesAsync();

                    return Results.Created(
                        $"/api/profiles/{profileId}/completed-workouts/{completed.Id}",
                        ToDto(completed)
                    );
                }
            )
            .WithName("CompleteWorkout")
            .Produces<CompletedWorkoutDto>(StatusCodes.Status201Created)
            .Produces(StatusCodes.Status404NotFound);

        // GET /api/profiles/{profileId}/completed-workouts/stats - Obtém estatísticas dos treinos
        group
            .MapGet(
                "/stats",
                async (Guid profileId, FitTrackerDbContext db) =>
                {
                    var completed = await db
                        .CompletedWorkouts.Where(c => c.UserProfileId == profileId)
                        .ToListAsync();

                    var weekAgo = DateTime.UtcNow.AddDays(-7);
                    var workoutsThisWeek = completed.Count(c => c.CompletedAt >= weekAgo);

                    // Calcular minutos totais baseado na duração real
                    var totalSeconds = completed.Sum(c => c.DurationSeconds);
                    var totalMinutes = totalSeconds / 60;

                    var stats = new WorkoutStatsDto(
                        TotalWorkoutsCompleted: completed.Count,
                        WorkoutsThisWeek: workoutsThisWeek,
                        TotalMinutesSpent: totalMinutes
                    );

                    return Results.Ok(stats);
                }
            )
            .WithName("GetWorkoutStats")
            .Produces<WorkoutStatsDto>(StatusCodes.Status200OK);

        // DELETE /api/profiles/{profileId}/completed-workouts/{id} - Remove um registro de treino concluído
        group
            .MapDelete(
                "/{id:guid}",
                async (Guid profileId, Guid id, FitTrackerDbContext db) =>
                {
                    var completed = await db.CompletedWorkouts.FirstOrDefaultAsync(c =>
                        c.Id == id && c.UserProfileId == profileId
                    );

                    if (completed is null)
                        return Results.NotFound();

                    db.CompletedWorkouts.Remove(completed);
                    await db.SaveChangesAsync();

                    return Results.NoContent();
                }
            )
            .WithName("DeleteCompletedWorkout")
            .Produces(StatusCodes.Status204NoContent)
            .Produces(StatusCodes.Status404NotFound);
    }

    private static CompletedWorkoutDto ToDto(CompletedWorkout completed) =>
        new(completed.Id, completed.WorkoutId, completed.CompletedAt, completed.DurationSeconds);
}
