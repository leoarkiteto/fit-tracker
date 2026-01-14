using System.Text.Json;
using FitTracker.Api.Features.Bioimpedance;
using FitTracker.Api.Features.Profiles;
using FitTracker.Api.Features.Workouts;
using FitTracker.Api.Shared.Data;
using Microsoft.EntityFrameworkCore;
using Microsoft.SemanticKernel;
using Microsoft.SemanticKernel.ChatCompletion;

namespace FitTracker.Api.Features.AI.WorkoutPlanning;

/// <summary>
/// AI Agent for generating personalized workout plans using Semantic Kernel + Ollama
/// </summary>
public class WorkoutPlanningAgent : IWorkoutPlanningAgent
{
    private readonly Kernel _kernel;
    private readonly FitTrackerDbContext _db;
    private readonly ILogger<WorkoutPlanningAgent> _logger;

    // Cache for generated plans (in production, use Redis or similar)
    private static readonly Dictionary<Guid, GeneratedPlanResponse> _planCache = new();

    public WorkoutPlanningAgent(
        Kernel kernel,
        FitTrackerDbContext db,
        ILogger<WorkoutPlanningAgent> logger)
    {
        _kernel = kernel;
        _db = db;
        _logger = logger;
    }

    /// <summary>
    /// Generate a personalized workout plan based on user profile and preferences
    /// </summary>
    public async Task<GeneratedPlanResponse> GeneratePlanAsync(
        GeneratePlanRequest request,
        CancellationToken cancellationToken = default)
    {
        _logger.LogInformation("Generating workout plan for user {UserId}", request.UserProfileId);

        // 1. Gather user context
        var context = await BuildUserContextAsync(request, cancellationToken);
        if (context is null)
        {
            throw new InvalidOperationException($"User profile {request.UserProfileId} not found");
        }

        // 2. Build the prompt with context
        var prompt = BuildPrompt(context, request.AdditionalNotes);

        // 3. Call the LLM
        var chatService = _kernel.GetRequiredService<IChatCompletionService>();
        
        var chatHistory = new ChatHistory();
        chatHistory.AddSystemMessage(WorkoutPlanningPrompts.SystemPrompt);
        chatHistory.AddUserMessage(prompt);

        _logger.LogDebug("Sending prompt to LLM: {Prompt}", prompt);

        var response = await chatService.GetChatMessageContentAsync(
            chatHistory,
            cancellationToken: cancellationToken);

        var responseText = response.Content ?? "";
        _logger.LogDebug("LLM Response: {Response}", responseText);

        // 4. Parse the response
        var plan = ParsePlanResponse(responseText, context.Goal);

        // 5. Cache the plan for later acceptance
        _planCache[plan.PlanId] = plan;

        _logger.LogInformation("Generated plan {PlanId} with {WorkoutCount} workouts",
            plan.PlanId, plan.Workouts.Count);

        return plan;
    }

    /// <summary>
    /// Accept a generated plan and save workouts to the database
    /// </summary>
    public async Task<AcceptPlanResponse> AcceptPlanAsync(
        AcceptPlanRequest request,
        CancellationToken cancellationToken = default)
    {
        _logger.LogInformation("Accepting plan {PlanId} for user {UserId}",
            request.PlanId, request.UserProfileId);

        var createdIds = new List<Guid>();

        foreach (var workoutDto in request.Workouts)
        {
            var workout = new Workout
            {
                Id = Guid.NewGuid(),
                Name = workoutDto.Name,
                Description = workoutDto.Description,
                Goal = workoutDto.Goal,
                Days = workoutDto.Days,
                UserProfileId = request.UserProfileId,
                CreatedAt = DateTime.UtcNow,
                UpdatedAt = DateTime.UtcNow,
                Exercises = workoutDto.Exercises.Select(e => new Exercise
                {
                    Id = Guid.NewGuid(),
                    Name = e.Name,
                    MuscleGroup = e.MuscleGroup,
                    Sets = e.Sets,
                    Reps = e.Reps,
                    Weight = e.SuggestedWeight,
                    RestSeconds = e.RestSeconds,
                    Notes = e.Notes
                }).ToList()
            };

            _db.Workouts.Add(workout);
            createdIds.Add(workout.Id);
        }

        await _db.SaveChangesAsync(cancellationToken);

        // Remove from cache
        _planCache.Remove(request.PlanId);

        _logger.LogInformation("Created {Count} workouts from plan {PlanId}",
            createdIds.Count, request.PlanId);

        return new AcceptPlanResponse(
            createdIds,
            $"Successfully created {createdIds.Count} workouts from your AI-generated plan!"
        );
    }

    /// <summary>
    /// Get a cached plan by ID (for preview)
    /// </summary>
    public Task<GeneratedPlanResponse?> GetCachedPlanAsync(Guid planId)
    {
        _planCache.TryGetValue(planId, out var plan);
        return Task.FromResult(plan);
    }

    private async Task<UserPlanningContext?> BuildUserContextAsync(
        GeneratePlanRequest request,
        CancellationToken cancellationToken)
    {
        var profile = await _db.UserProfiles.FindAsync(
            new object[] { request.UserProfileId },
            cancellationToken);

        if (profile is null) return null;

        // Get completed workouts count
        var totalCompleted = await _db.CompletedWorkouts
            .CountAsync(c => c.UserProfileId == request.UserProfileId, cancellationToken);

        var monthStart = new DateTime(DateTime.UtcNow.Year, DateTime.UtcNow.Month, 1);
        var monthlyCompleted = await _db.CompletedWorkouts
            .CountAsync(c => c.UserProfileId == request.UserProfileId
                && c.CompletedAt >= monthStart, cancellationToken);

        // Get latest bioimpedance
        var latestBio = await _db.BioimpedanceData
            .Where(b => b.UserProfileId == request.UserProfileId)
            .OrderByDescending(b => b.Date)
            .FirstOrDefaultAsync(cancellationToken);

        // Get current exercises the user does
        var currentExercises = await _db.Workouts
            .Where(w => w.UserProfileId == request.UserProfileId)
            .SelectMany(w => w.Exercises)
            .Select(e => e.Name)
            .Distinct()
            .ToListAsync(cancellationToken);

        return new UserPlanningContext(
            UserProfileId: profile.Id,
            Name: profile.Name,
            Age: profile.Age,
            CurrentWeight: profile.CurrentWeight,
            GoalWeight: profile.GoalWeight,
            ExperienceLevel: request.OverrideEquipment.HasValue
                ? profile.ExperienceLevel
                : profile.ExperienceLevel,
            AvailableDaysPerWeek: request.OverrideDaysPerWeek ?? profile.AvailableDaysPerWeek,
            PreferredWorkoutDuration: request.OverrideDurationMinutes ?? profile.PreferredWorkoutDuration,
            EquipmentType: request.OverrideEquipment ?? profile.EquipmentType,
            Goal: request.Goal,
            TotalCompletedWorkouts: totalCompleted,
            WorkoutsThisMonth: monthlyCompleted,
            LatestBodyFatPercentage: latestBio?.BodyFatPercentage,
            LatestMuscleMass: latestBio?.MuscleMass,
            CurrentExercises: currentExercises
        );
    }

    private string BuildPrompt(UserPlanningContext context, string? additionalNotes)
    {
        return WorkoutPlanningPrompts.GenerateWeeklyPlanPrompt
            .Replace("{{$userName}}", context.Name)
            .Replace("{{$userAge}}", context.Age?.ToString() ?? "Not specified")
            .Replace("{{$currentWeight}}", context.CurrentWeight?.ToString("F1") ?? "Not specified")
            .Replace("{{$goalWeight}}", context.GoalWeight?.ToString("F1") ?? "Not specified")
            .Replace("{{$experienceLevel}}", WorkoutPlanningPrompts.GetExperienceDescription(context.ExperienceLevel))
            .Replace("{{$availableDays}}", context.AvailableDaysPerWeek.ToString())
            .Replace("{{$preferredDuration}}", context.PreferredWorkoutDuration?.ToString() ?? "45-60")
            .Replace("{{$equipment}}", WorkoutPlanningPrompts.GetEquipmentDescription(context.EquipmentType))
            .Replace("{{$goal}}", WorkoutPlanningPrompts.GetGoalDescription(context.Goal))
            .Replace("{{$totalWorkouts}}", context.TotalCompletedWorkouts.ToString())
            .Replace("{{$monthlyWorkouts}}", context.WorkoutsThisMonth.ToString())
            .Replace("{{$bodyFat}}", context.LatestBodyFatPercentage?.ToString("F1") ?? "Not available")
            .Replace("{{$muscleMass}}", context.LatestMuscleMass?.ToString("F1") ?? "Not available")
            .Replace("{{$currentExercises}}", context.CurrentExercises.Count > 0
                ? string.Join(", ", context.CurrentExercises)
                : "None recorded yet")
            .Replace("{{$additionalNotes}}", additionalNotes ?? "No additional notes");
    }

    private GeneratedPlanResponse ParsePlanResponse(string responseText, WorkoutGoal defaultGoal)
    {
        try
        {
            // Try to extract JSON from the response (LLM might add extra text)
            var jsonStart = responseText.IndexOf('{');
            var jsonEnd = responseText.LastIndexOf('}');

            if (jsonStart == -1 || jsonEnd == -1)
            {
                throw new InvalidOperationException("No valid JSON found in LLM response");
            }

            var jsonText = responseText.Substring(jsonStart, jsonEnd - jsonStart + 1);

            var options = new JsonSerializerOptions
            {
                PropertyNameCaseInsensitive = true
            };

            var parsed = JsonSerializer.Deserialize<LlmPlanResponse>(jsonText, options)
                ?? throw new InvalidOperationException("Failed to parse LLM response");

            var workouts = parsed.Workouts.Select(w => new GeneratedWorkoutDto(
                Name: w.Name,
                Description: w.Description,
                Goal: ParseGoal(w.Goal, defaultGoal),
                Days: w.Days.Select(ParseDay).ToList(),
                Exercises: w.Exercises.Select(e => new GeneratedExerciseDto(
                    Name: e.Name,
                    MuscleGroup: ParseMuscleGroup(e.MuscleGroup),
                    Sets: e.Sets,
                    Reps: e.Reps,
                    SuggestedWeight: e.SuggestedWeight,
                    RestSeconds: e.RestSeconds,
                    Notes: e.Notes
                )).ToList(),
                EstimatedDurationMinutes: w.EstimatedDurationMinutes
            )).ToList();

            return new GeneratedPlanResponse(
                PlanId: Guid.NewGuid(),
                Summary: parsed.Summary,
                Rationale: parsed.Rationale,
                Workouts: workouts,
                GeneratedAt: DateTime.UtcNow
            );
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to parse LLM response: {Response}", responseText);

            // Return a fallback plan
            return CreateFallbackPlan(defaultGoal);
        }
    }

    private GeneratedPlanResponse CreateFallbackPlan(WorkoutGoal goal)
    {
        _logger.LogWarning("Using fallback plan due to LLM parsing failure");

        return new GeneratedPlanResponse(
            PlanId: Guid.NewGuid(),
            Summary: "Basic workout plan (AI response could not be parsed)",
            Rationale: "This is a default plan. Please try generating again or contact support.",
            Workouts: new List<GeneratedWorkoutDto>
            {
                new GeneratedWorkoutDto(
                    Name: "Full Body Workout",
                    Description: "A balanced full-body workout",
                    Goal: goal,
                    Days: new List<DayOfWeekEnum> { DayOfWeekEnum.Monday, DayOfWeekEnum.Wednesday, DayOfWeekEnum.Friday },
                    Exercises: new List<GeneratedExerciseDto>
                    {
                        new("Squat", MuscleGroup.Legs, 3, 10, null, 90, "Focus on form"),
                        new("Bench Press", MuscleGroup.Chest, 3, 10, null, 90, "Control the descent"),
                        new("Bent Over Row", MuscleGroup.Back, 3, 10, null, 90, "Keep back straight"),
                        new("Shoulder Press", MuscleGroup.Shoulders, 3, 10, null, 60, "Don't lock elbows"),
                        new("Plank", MuscleGroup.Abs, 3, 30, null, 60, "30 seconds hold")
                    },
                    EstimatedDurationMinutes: 45
                )
            },
            GeneratedAt: DateTime.UtcNow
        );
    }

    private static WorkoutGoal ParseGoal(string goal, WorkoutGoal defaultGoal)
    {
        return goal.ToLowerInvariant() switch
        {
            "hypertrophy" => WorkoutGoal.Hypertrophy,
            "strength" => WorkoutGoal.Strength,
            "endurance" => WorkoutGoal.Endurance,
            "weightloss" or "weight_loss" or "weight loss" => WorkoutGoal.WeightLoss,
            "maintenance" => WorkoutGoal.Maintenance,
            _ => defaultGoal
        };
    }

    private static DayOfWeekEnum ParseDay(string day)
    {
        return day.ToLowerInvariant() switch
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
    }

    private static MuscleGroup ParseMuscleGroup(string group)
    {
        return group.ToLowerInvariant() switch
        {
            "chest" => MuscleGroup.Chest,
            "back" => MuscleGroup.Back,
            "shoulders" => MuscleGroup.Shoulders,
            "biceps" => MuscleGroup.Biceps,
            "triceps" => MuscleGroup.Triceps,
            "legs" => MuscleGroup.Legs,
            "glutes" => MuscleGroup.Glutes,
            "abs" or "core" => MuscleGroup.Abs,
            "cardio" => MuscleGroup.Cardio,
            _ => MuscleGroup.Chest
        };
    }

    // Internal classes for JSON parsing
    private record LlmPlanResponse(
        string Summary,
        string Rationale,
        List<LlmWorkout> Workouts
    );

    private record LlmWorkout(
        string Name,
        string? Description,
        string Goal,
        List<string> Days,
        int EstimatedDurationMinutes,
        List<LlmExercise> Exercises
    );

    private record LlmExercise(
        string Name,
        string MuscleGroup,
        int Sets,
        int Reps,
        double? SuggestedWeight,
        int RestSeconds,
        string? Notes
    );
}

/// <summary>
/// Interface for the Workout Planning Agent
/// </summary>
public interface IWorkoutPlanningAgent
{
    Task<GeneratedPlanResponse> GeneratePlanAsync(
        GeneratePlanRequest request,
        CancellationToken cancellationToken = default);

    Task<AcceptPlanResponse> AcceptPlanAsync(
        AcceptPlanRequest request,
        CancellationToken cancellationToken = default);

    Task<GeneratedPlanResponse?> GetCachedPlanAsync(Guid planId);
}
