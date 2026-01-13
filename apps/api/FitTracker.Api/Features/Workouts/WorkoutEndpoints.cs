using Microsoft.EntityFrameworkCore;
using FitTracker.Api.Shared.Data;

namespace FitTracker.Api.Features.Workouts;

public static class WorkoutEndpoints
{
    public static void MapWorkoutEndpoints(this WebApplication app)
    {
        var group = app.MapGroup("/api/profiles/{profileId:guid}/workouts").WithTags("Workouts");

        // GET /api/profiles/{profileId}/workouts - Lista todos os treinos de um perfil
        group.MapGet("/", async (Guid profileId, FitTrackerDbContext db) =>
        {
            var workouts = await db.Workouts
                .Include(w => w.Exercises)
                .Where(w => w.UserProfileId == profileId)
                .OrderByDescending(w => w.CreatedAt)
                .ToListAsync();

            return Results.Ok(workouts.Select(ToDto));
        })
        .WithName("GetAllWorkouts")
        .Produces<IEnumerable<WorkoutDto>>(StatusCodes.Status200OK);

        // GET /api/profiles/{profileId}/workouts/{id} - Obtém um treino por ID
        group.MapGet("/{id:guid}", async (Guid profileId, Guid id, FitTrackerDbContext db) =>
        {
            var workout = await db.Workouts
                .Include(w => w.Exercises)
                .FirstOrDefaultAsync(w => w.Id == id && w.UserProfileId == profileId);

            return workout is null ? Results.NotFound() : Results.Ok(ToDto(workout));
        })
        .WithName("GetWorkoutById")
        .Produces<WorkoutDto>(StatusCodes.Status200OK)
        .Produces(StatusCodes.Status404NotFound);

        // POST /api/profiles/{profileId}/workouts - Cria um novo treino
        group.MapPost("/", async (Guid profileId, CreateWorkoutRequest request, FitTrackerDbContext db) =>
        {
            // Verificar se o perfil existe
            var profileExists = await db.UserProfiles.AnyAsync(p => p.Id == profileId);
            if (!profileExists) return Results.NotFound("Profile not found");

            var workout = new Workout
            {
                Id = Guid.NewGuid(),
                Name = request.Name,
                Description = request.Description,
                Goal = ParseWorkoutGoal(request.Goal),
                Days = request.Days.Select(ParseDayOfWeek).ToList(),
                UserProfileId = profileId,
                CreatedAt = DateTime.UtcNow,
                UpdatedAt = DateTime.UtcNow,
                Exercises = request.Exercises.Select(e => new Exercise
                {
                    Id = Guid.NewGuid(),
                    Name = e.Name,
                    MuscleGroup = ParseMuscleGroup(e.MuscleGroup),
                    Sets = e.Sets,
                    Reps = e.Reps,
                    Weight = e.Weight,
                    RestSeconds = e.RestSeconds,
                    Notes = e.Notes
                }).ToList()
            };

            db.Workouts.Add(workout);
            await db.SaveChangesAsync();

            return Results.Created($"/api/profiles/{profileId}/workouts/{workout.Id}", ToDto(workout));
        })
        .WithName("CreateWorkout")
        .Produces<WorkoutDto>(StatusCodes.Status201Created)
        .Produces(StatusCodes.Status404NotFound);

        // PUT /api/profiles/{profileId}/workouts/{id} - Atualiza um treino
        group.MapPut("/{id:guid}", async (Guid profileId, Guid id, UpdateWorkoutRequest request, FitTrackerDbContext db) =>
        {
            // Buscar apenas o workout sem os exercícios para evitar problemas de tracking
            var workout = await db.Workouts
                .FirstOrDefaultAsync(w => w.Id == id && w.UserProfileId == profileId);

            if (workout is null) return Results.NotFound();

            // Atualizar campos do workout
            workout.Name = request.Name;
            workout.Description = request.Description;
            workout.Goal = ParseWorkoutGoal(request.Goal);
            workout.Days = request.Days.Select(ParseDayOfWeek).ToList();
            workout.UpdatedAt = DateTime.UtcNow;

            // Salvar as mudanças do workout primeiro
            await db.SaveChangesAsync();

            // Deletar todos os exercícios antigos diretamente no banco
            await db.Exercises
                .Where(e => e.WorkoutId == id)
                .ExecuteDeleteAsync();

            // Criar novos exercícios
            var newExercises = request.Exercises.Select(e => new Exercise
            {
                Id = Guid.NewGuid(),
                Name = e.Name,
                MuscleGroup = ParseMuscleGroup(e.MuscleGroup),
                Sets = e.Sets,
                Reps = e.Reps,
                Weight = e.Weight,
                RestSeconds = e.RestSeconds,
                Notes = e.Notes,
                WorkoutId = id
            }).ToList();

            await db.Exercises.AddRangeAsync(newExercises);
            await db.SaveChangesAsync();

            // Recarregar o workout com os novos exercícios para retornar
            var updatedWorkout = await db.Workouts
                .Include(w => w.Exercises)
                .FirstAsync(w => w.Id == id);

            return Results.Ok(ToDto(updatedWorkout));
        })
        .WithName("UpdateWorkout")
        .Produces<WorkoutDto>(StatusCodes.Status200OK)
        .Produces(StatusCodes.Status404NotFound);

        // DELETE /api/profiles/{profileId}/workouts/{id} - Remove um treino
        group.MapDelete("/{id:guid}", async (Guid profileId, Guid id, FitTrackerDbContext db) =>
        {
            var workout = await db.Workouts
                .FirstOrDefaultAsync(w => w.Id == id && w.UserProfileId == profileId);

            if (workout is null) return Results.NotFound();

            db.Workouts.Remove(workout);
            await db.SaveChangesAsync();

            return Results.NoContent();
        })
        .WithName("DeleteWorkout")
        .Produces(StatusCodes.Status204NoContent)
        .Produces(StatusCodes.Status404NotFound);

        // GET /api/profiles/{profileId}/workouts/today - Obtém treinos de hoje
        group.MapGet("/today", async (Guid profileId, FitTrackerDbContext db) =>
        {
            var today = DateTime.UtcNow.DayOfWeek;
            var dayEnum = MapSystemDayToEnum(today);

            var workouts = await db.Workouts
                .Include(w => w.Exercises)
                .Where(w => w.UserProfileId == profileId && w.Days.Contains(dayEnum))
                .ToListAsync();

            return Results.Ok(workouts.Select(ToDto));
        })
        .WithName("GetTodayWorkouts")
        .Produces<IEnumerable<WorkoutDto>>(StatusCodes.Status200OK);
    }

    private static WorkoutDto ToDto(Workout workout) => new(
        workout.Id,
        workout.Name,
        workout.Description,
        workout.Goal.ToString().ToLower(),
        workout.Days.Select(d => d.ToString().ToLower()).ToList(),
        workout.Exercises.Select(ToExerciseDto).ToList(),
        workout.CreatedAt,
        workout.UpdatedAt,
        workout.IsCompleted,
        workout.CompletedAt
    );

    private static ExerciseDto ToExerciseDto(Exercise exercise) => new(
        exercise.Id,
        exercise.Name,
        exercise.MuscleGroup.ToString().ToLower(),
        exercise.Sets,
        exercise.Reps,
        exercise.Weight,
        exercise.RestSeconds,
        exercise.Notes
    );

    private static WorkoutGoal ParseWorkoutGoal(string goal) => goal.ToLower() switch
    {
        "hypertrophy" => WorkoutGoal.Hypertrophy,
        "strength" => WorkoutGoal.Strength,
        "endurance" => WorkoutGoal.Endurance,
        "weight_loss" => WorkoutGoal.WeightLoss,
        "maintenance" => WorkoutGoal.Maintenance,
        _ => WorkoutGoal.Hypertrophy
    };

    private static DayOfWeekEnum ParseDayOfWeek(string day) => day.ToLower() switch
    {
        "monday" => DayOfWeekEnum.Monday,
        "tuesday" => DayOfWeekEnum.Tuesday,
        "wednesday" => DayOfWeekEnum.Wednesday,
        "thursday" => DayOfWeekEnum.Thursday,
        "friday" => DayOfWeekEnum.Friday,
        "saturday" => DayOfWeekEnum.Saturday,
        "sunday" => DayOfWeekEnum.Sunday,
        _ => DayOfWeekEnum.Monday
    };

    private static MuscleGroup ParseMuscleGroup(string group) => group.ToLower() switch
    {
        "chest" => MuscleGroup.Chest,
        "back" => MuscleGroup.Back,
        "shoulders" => MuscleGroup.Shoulders,
        "biceps" => MuscleGroup.Biceps,
        "triceps" => MuscleGroup.Triceps,
        "legs" => MuscleGroup.Legs,
        "glutes" => MuscleGroup.Glutes,
        "abs" => MuscleGroup.Abs,
        "cardio" => MuscleGroup.Cardio,
        _ => MuscleGroup.Chest
    };

    private static DayOfWeekEnum MapSystemDayToEnum(DayOfWeek day) => day switch
    {
        DayOfWeek.Monday => DayOfWeekEnum.Monday,
        DayOfWeek.Tuesday => DayOfWeekEnum.Tuesday,
        DayOfWeek.Wednesday => DayOfWeekEnum.Wednesday,
        DayOfWeek.Thursday => DayOfWeekEnum.Thursday,
        DayOfWeek.Friday => DayOfWeekEnum.Friday,
        DayOfWeek.Saturday => DayOfWeekEnum.Saturday,
        DayOfWeek.Sunday => DayOfWeekEnum.Sunday,
        _ => DayOfWeekEnum.Monday
    };
}
