using Microsoft.AspNetCore.Mvc;

namespace FitTracker.Api.Features.AI.WorkoutPlanning;

public static class WorkoutPlanningEndpoints
{
    public static void MapWorkoutPlanningEndpoints(this WebApplication app)
    {
        var group = app.MapGroup("/api/ai/planning")
            .WithTags("AI - Workout Planning")
            .RequireAuthorization();

        // POST /api/ai/planning/generate - Generate a new workout plan
        group
            .MapPost(
                "/generate",
                async (
                    [FromBody] GeneratePlanRequest request,
                    IWorkoutPlanningAgent agent,
                    CancellationToken cancellationToken) =>
                {
                    try
                    {
                        var plan = await agent.GeneratePlanAsync(request, cancellationToken);
                        return Results.Ok(plan);
                    }
                    catch (InvalidOperationException ex)
                    {
                        return Results.BadRequest(new { error = ex.Message });
                    }
                    catch (Exception)
                    {
                        return Results.Problem(
                            detail: "Failed to generate workout plan. Please try again.",
                            statusCode: StatusCodes.Status500InternalServerError);
                    }
                }
            )
            .WithName("GenerateWorkoutPlan")
            .WithDescription("Generate a personalized workout plan using AI based on user profile and preferences")
            .Produces<GeneratedPlanResponse>(StatusCodes.Status200OK)
            .Produces(StatusCodes.Status400BadRequest)
            .Produces(StatusCodes.Status500InternalServerError);

        // GET /api/ai/planning/preview/{planId} - Get a cached plan for preview
        group
            .MapGet(
                "/preview/{planId:guid}",
                async (Guid planId, IWorkoutPlanningAgent agent) =>
                {
                    var plan = await agent.GetCachedPlanAsync(planId);
                    return plan is null
                        ? Results.NotFound(new { error = "Plan not found or has expired" })
                        : Results.Ok(plan);
                }
            )
            .WithName("PreviewWorkoutPlan")
            .WithDescription("Get a previously generated plan for review before accepting")
            .Produces<GeneratedPlanResponse>(StatusCodes.Status200OK)
            .Produces(StatusCodes.Status404NotFound);

        // POST /api/ai/planning/accept - Accept and save a generated plan
        group
            .MapPost(
                "/accept",
                async (
                    [FromBody] AcceptPlanRequest request,
                    IWorkoutPlanningAgent agent,
                    CancellationToken cancellationToken) =>
                {
                    try
                    {
                        var result = await agent.AcceptPlanAsync(request, cancellationToken);
                        return Results.Ok(result);
                    }
                    catch (Exception)
                    {
                        return Results.Problem(
                            detail: "Failed to save workout plan. Please try again.",
                            statusCode: StatusCodes.Status500InternalServerError);
                    }
                }
            )
            .WithName("AcceptWorkoutPlan")
            .WithDescription("Accept a generated plan and save the workouts to the user's account")
            .Produces<AcceptPlanResponse>(StatusCodes.Status200OK)
            .Produces(StatusCodes.Status500InternalServerError);

        // GET /api/ai/planning/status - Check if AI planning is available
        group
            .MapGet(
                "/status",
                (IConfiguration config) =>
                {
                    var ollamaEndpoint = config["AI:Ollama:Endpoint"];
                    var ollamaModel = config["AI:Ollama:Model"];

                    return Results.Ok(new
                    {
                        available = !string.IsNullOrEmpty(ollamaEndpoint),
                        provider = "Ollama",
                        model = ollamaModel ?? "unknown",
                        endpoint = ollamaEndpoint ?? "not configured"
                    });
                }
            )
            .WithName("GetAIPlanningStatus")
            .WithDescription("Check if AI workout planning is available and configured")
            .AllowAnonymous();
    }
}
