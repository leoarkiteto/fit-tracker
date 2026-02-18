namespace FitTracker.Api.Features.Bioimpedance;

public class BioimpedanceData
{
    public Guid Id { get; set; }
    public DateTime Date { get; set; }
    public double Weight { get; set; } // kg
    public double BodyFatPercentage { get; set; } // %
    public double MuscleMass { get; set; } // kg
    public double BoneMass { get; set; } // kg
    public double WaterPercentage { get; set; } // %
    public int VisceralFat { get; set; } // nível
    public int Bmr { get; set; } // Taxa Metabólica Basal (kcal)
    public int MetabolicAge { get; set; } // anos
    public string? Notes { get; set; }

    // Foreign key para UserProfile
    public Guid UserProfileId { get; set; }
}
