namespace FitTracker.Api.Features.Bioimpedance;

public record BioimpedanceDto(
    Guid Id,
    DateTime Date,
    double Weight,
    double BodyFatPercentage,
    double MuscleMass,
    double BoneMass,
    double WaterPercentage,
    int VisceralFat,
    int Bmr,
    int MetabolicAge,
    string? Notes
);

public record CreateBioimpedanceRequest(
    DateTime Date,
    double Weight,
    double BodyFatPercentage,
    double MuscleMass,
    double BoneMass,
    double WaterPercentage,
    int VisceralFat,
    int Bmr,
    int MetabolicAge,
    string? Notes
);

public record UpdateBioimpedanceRequest(
    DateTime Date,
    double Weight,
    double BodyFatPercentage,
    double MuscleMass,
    double BoneMass,
    double WaterPercentage,
    int VisceralFat,
    int Bmr,
    int MetabolicAge,
    string? Notes
);
