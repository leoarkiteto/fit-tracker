using FitTracker.Api.Data;
using FitTracker.Api.DTOs;
using FitTracker.Api.Models;
using Microsoft.EntityFrameworkCore;

namespace FitTracker.Api.Endpoints;

public static class ProfileEndpoints
{
    public static void MapProfileEndpoints(this WebApplication app)
    {
        var group = app.MapGroup("/api/profiles").WithTags("Profiles");

        // GET /api/profiles - Lista todos os perfis
        group
            .MapGet(
                "/",
                async (FitTrackerDbContext db) =>
                {
                    var profiles = await db.UserProfiles.ToListAsync();
                    return Results.Ok(profiles.Select(ToDto));
                }
            )
            .WithName("GetAllProfiles")
            .Produces<IEnumerable<UserProfileDto>>(StatusCodes.Status200OK);

        // GET /api/profiles/{id} - Obtém um perfil por ID
        group
            .MapGet(
                "/{id:guid}",
                async (Guid id, FitTrackerDbContext db) =>
                {
                    var profile = await db.UserProfiles.FindAsync(id);
                    return profile is null ? Results.NotFound() : Results.Ok(ToDto(profile));
                }
            )
            .WithName("GetProfileById")
            .Produces<UserProfileDto>(StatusCodes.Status200OK)
            .Produces(StatusCodes.Status404NotFound);

        // POST /api/profiles - Cria um novo perfil
        group
            .MapPost(
                "/",
                async (CreateUserProfileRequest request, FitTrackerDbContext db) =>
                {
                    var profile = new UserProfile
                    {
                        Id = Guid.NewGuid(),
                        Name = request.Name,
                        Age = request.Age,
                        Height = request.Height,
                        CurrentWeight = request.CurrentWeight,
                        GoalWeight = request.GoalWeight,
                        AvatarUrl = request.AvatarUrl,
                        CreatedAt = DateTime.UtcNow,
                        UpdatedAt = DateTime.UtcNow,
                    };

                    db.UserProfiles.Add(profile);
                    await db.SaveChangesAsync();

                    return Results.Created($"/api/profiles/{profile.Id}", ToDto(profile));
                }
            )
            .WithName("CreateProfile")
            .Produces<UserProfileDto>(StatusCodes.Status201Created);

        // PUT /api/profiles/{id} - Atualiza um perfil
        group
            .MapPut(
                "/{id:guid}",
                async (Guid id, UpdateUserProfileRequest request, FitTrackerDbContext db) =>
                {
                    var profile = await db.UserProfiles.FindAsync(id);
                    if (profile is null)
                        return Results.NotFound();

                    profile.Name = request.Name;
                    profile.Age = request.Age;
                    profile.Height = request.Height;
                    profile.CurrentWeight = request.CurrentWeight;
                    profile.GoalWeight = request.GoalWeight;
                    profile.AvatarUrl = request.AvatarUrl;
                    profile.UpdatedAt = DateTime.UtcNow;

                    await db.SaveChangesAsync();

                    return Results.Ok(ToDto(profile));
                }
            )
            .WithName("UpdateProfile")
            .Produces<UserProfileDto>(StatusCodes.Status200OK)
            .Produces(StatusCodes.Status404NotFound);

        // DELETE /api/profiles/{id} - Remove um perfil
        group
            .MapDelete(
                "/{id:guid}",
                async (Guid id, FitTrackerDbContext db) =>
                {
                    var profile = await db.UserProfiles.FindAsync(id);
                    if (profile is null)
                        return Results.NotFound();

                    db.UserProfiles.Remove(profile);
                    await db.SaveChangesAsync();

                    return Results.NoContent();
                }
            )
            .WithName("DeleteProfile")
            .Produces(StatusCodes.Status204NoContent)
            .Produces(StatusCodes.Status404NotFound);
    }

    private static UserProfileDto ToDto(UserProfile profile) =>
        new(
            profile.Id,
            profile.Name,
            profile.Age,
            profile.Height,
            profile.CurrentWeight,
            profile.GoalWeight,
            profile.AvatarUrl
        );
}
