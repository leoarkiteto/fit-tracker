using FitTracker.Api.Data;
using FitTracker.Api.DTOs;
using FitTracker.Api.Models;
using Microsoft.EntityFrameworkCore;

namespace FitTracker.Api.Endpoints;

public static class BioimpedanceEndpoints
{
    public static void MapBioimpedanceEndpoints(this WebApplication app)
    {
        var group = app.MapGroup("/api/profiles/{profileId:guid}/bioimpedance")
            .WithTags("Bioimpedance");

        // GET /api/profiles/{profileId}/bioimpedance - Lista todo o histórico de bioimpedância
        group
            .MapGet(
                "/",
                async (Guid profileId, FitTrackerDbContext db) =>
                {
                    var history = await db
                        .BioimpedanceData.Where(b => b.UserProfileId == profileId)
                        .OrderByDescending(b => b.Date)
                        .ToListAsync();

                    return Results.Ok(history.Select(ToDto));
                }
            )
            .WithName("GetBioimpedanceHistory")
            .Produces<IEnumerable<BioimpedanceDto>>(StatusCodes.Status200OK);

        // GET /api/profiles/{profileId}/bioimpedance/{id} - Obtém um registro por ID
        group
            .MapGet(
                "/{id:guid}",
                async (Guid profileId, Guid id, FitTrackerDbContext db) =>
                {
                    var data = await db.BioimpedanceData.FirstOrDefaultAsync(b =>
                        b.Id == id && b.UserProfileId == profileId
                    );

                    return data is null ? Results.NotFound() : Results.Ok(ToDto(data));
                }
            )
            .WithName("GetBioimpedanceById")
            .Produces<BioimpedanceDto>(StatusCodes.Status200OK)
            .Produces(StatusCodes.Status404NotFound);

        // GET /api/profiles/{profileId}/bioimpedance/latest - Obtém o registro mais recente
        group
            .MapGet(
                "/latest",
                async (Guid profileId, FitTrackerDbContext db) =>
                {
                    var data = await db
                        .BioimpedanceData.Where(b => b.UserProfileId == profileId)
                        .OrderByDescending(b => b.Date)
                        .FirstOrDefaultAsync();

                    return data is null ? Results.NotFound() : Results.Ok(ToDto(data));
                }
            )
            .WithName("GetLatestBioimpedance")
            .Produces<BioimpedanceDto>(StatusCodes.Status200OK)
            .Produces(StatusCodes.Status404NotFound);

        // POST /api/profiles/{profileId}/bioimpedance - Cria um novo registro
        group
            .MapPost(
                "/",
                async (Guid profileId, CreateBioimpedanceRequest request, FitTrackerDbContext db) =>
                {
                    // Verificar se o perfil existe
                    var profileExists = await db.UserProfiles.AnyAsync(p => p.Id == profileId);
                    if (!profileExists)
                        return Results.NotFound("Profile not found");

                    var data = new BioimpedanceData
                    {
                        Id = Guid.NewGuid(),
                        Date = request.Date,
                        Weight = request.Weight,
                        BodyFatPercentage = request.BodyFatPercentage,
                        MuscleMass = request.MuscleMass,
                        BoneMass = request.BoneMass,
                        WaterPercentage = request.WaterPercentage,
                        VisceralFat = request.VisceralFat,
                        Bmr = request.Bmr,
                        MetabolicAge = request.MetabolicAge,
                        Notes = request.Notes,
                        UserProfileId = profileId,
                    };

                    db.BioimpedanceData.Add(data);
                    await db.SaveChangesAsync();

                    return Results.Created(
                        $"/api/profiles/{profileId}/bioimpedance/{data.Id}",
                        ToDto(data)
                    );
                }
            )
            .WithName("CreateBioimpedance")
            .Produces<BioimpedanceDto>(StatusCodes.Status201Created)
            .Produces(StatusCodes.Status404NotFound);

        // PUT /api/profiles/{profileId}/bioimpedance/{id} - Atualiza um registro
        group
            .MapPut(
                "/{id:guid}",
                async (
                    Guid profileId,
                    Guid id,
                    UpdateBioimpedanceRequest request,
                    FitTrackerDbContext db
                ) =>
                {
                    var data = await db.BioimpedanceData.FirstOrDefaultAsync(b =>
                        b.Id == id && b.UserProfileId == profileId
                    );

                    if (data is null)
                        return Results.NotFound();

                    data.Date = request.Date;
                    data.Weight = request.Weight;
                    data.BodyFatPercentage = request.BodyFatPercentage;
                    data.MuscleMass = request.MuscleMass;
                    data.BoneMass = request.BoneMass;
                    data.WaterPercentage = request.WaterPercentage;
                    data.VisceralFat = request.VisceralFat;
                    data.Bmr = request.Bmr;
                    data.MetabolicAge = request.MetabolicAge;
                    data.Notes = request.Notes;

                    await db.SaveChangesAsync();

                    return Results.Ok(ToDto(data));
                }
            )
            .WithName("UpdateBioimpedance")
            .Produces<BioimpedanceDto>(StatusCodes.Status200OK)
            .Produces(StatusCodes.Status404NotFound);

        // DELETE /api/profiles/{profileId}/bioimpedance/{id} - Remove um registro
        group
            .MapDelete(
                "/{id:guid}",
                async (Guid profileId, Guid id, FitTrackerDbContext db) =>
                {
                    var data = await db.BioimpedanceData.FirstOrDefaultAsync(b =>
                        b.Id == id && b.UserProfileId == profileId
                    );

                    if (data is null)
                        return Results.NotFound();

                    db.BioimpedanceData.Remove(data);
                    await db.SaveChangesAsync();

                    return Results.NoContent();
                }
            )
            .WithName("DeleteBioimpedance")
            .Produces(StatusCodes.Status204NoContent)
            .Produces(StatusCodes.Status404NotFound);
    }

    private static BioimpedanceDto ToDto(BioimpedanceData data) =>
        new(
            data.Id,
            data.Date,
            data.Weight,
            data.BodyFatPercentage,
            data.MuscleMass,
            data.BoneMass,
            data.WaterPercentage,
            data.VisceralFat,
            data.Bmr,
            data.MetabolicAge,
            data.Notes
        );
}
