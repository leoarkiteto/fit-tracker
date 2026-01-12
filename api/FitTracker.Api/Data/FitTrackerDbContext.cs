using FitTracker.Api.Models;
using Microsoft.EntityFrameworkCore;

namespace FitTracker.Api.Data;

public class FitTrackerDbContext : DbContext
{
    public FitTrackerDbContext(DbContextOptions<FitTrackerDbContext> options)
        : base(options) { }

    public DbSet<User> Users => Set<User>();
    public DbSet<UserProfile> UserProfiles => Set<UserProfile>();
    public DbSet<Workout> Workouts => Set<Workout>();
    public DbSet<Exercise> Exercises => Set<Exercise>();
    public DbSet<BioimpedanceData> BioimpedanceData => Set<BioimpedanceData>();
    public DbSet<CompletedWorkout> CompletedWorkouts => Set<CompletedWorkout>();

    protected override void OnModelCreating(ModelBuilder modelBuilder)
    {
        base.OnModelCreating(modelBuilder);

        // User
        modelBuilder.Entity<User>(entity =>
        {
            entity.HasKey(e => e.Id);
            entity.Property(e => e.Email).IsRequired().HasMaxLength(256);
            entity.HasIndex(e => e.Email).IsUnique();
            entity.Property(e => e.PasswordHash).IsRequired();
            entity.Property(e => e.Name).IsRequired().HasMaxLength(100);

            entity
                .HasOne(e => e.UserProfile)
                .WithOne()
                .HasForeignKey<User>(e => e.UserProfileId)
                .OnDelete(DeleteBehavior.SetNull);
        });

        // UserProfile
        modelBuilder.Entity<UserProfile>(entity =>
        {
            entity.HasKey(e => e.Id);
            entity.Property(e => e.Name).IsRequired().HasMaxLength(100);
            entity.Property(e => e.AvatarUrl).HasMaxLength(500);
        });

        // Workout
        modelBuilder.Entity<Workout>(entity =>
        {
            entity.HasKey(e => e.Id);
            entity.Property(e => e.Name).IsRequired().HasMaxLength(100);
            entity.Property(e => e.Description).HasMaxLength(500);
            entity.Property(e => e.Goal).HasConversion<string>();

            // Converter lista de dias para string JSON
            entity
                .Property(e => e.Days)
                .HasConversion(
                    v => string.Join(',', v.Select(d => d.ToString())),
                    v =>
                        v.Split(',', StringSplitOptions.RemoveEmptyEntries)
                            .Select(d => Enum.Parse<DayOfWeekEnum>(d))
                            .ToList()
                );

            entity
                .HasOne(e => e.UserProfile)
                .WithMany()
                .HasForeignKey(e => e.UserProfileId)
                .OnDelete(DeleteBehavior.Cascade);

            entity
                .HasMany(e => e.Exercises)
                .WithOne(e => e.Workout)
                .HasForeignKey(e => e.WorkoutId)
                .OnDelete(DeleteBehavior.Cascade);
        });

        // Exercise
        modelBuilder.Entity<Exercise>(entity =>
        {
            entity.HasKey(e => e.Id);
            entity.Property(e => e.Name).IsRequired().HasMaxLength(100);
            entity.Property(e => e.MuscleGroup).HasConversion<string>();
            entity.Property(e => e.Notes).HasMaxLength(500);
        });

        // BioimpedanceData
        modelBuilder.Entity<BioimpedanceData>(entity =>
        {
            entity.HasKey(e => e.Id);
            entity.Property(e => e.Notes).HasMaxLength(500);

            entity
                .HasOne(e => e.UserProfile)
                .WithMany()
                .HasForeignKey(e => e.UserProfileId)
                .OnDelete(DeleteBehavior.Cascade);
        });

        // CompletedWorkout
        modelBuilder.Entity<CompletedWorkout>(entity =>
        {
            entity.HasKey(e => e.Id);

            entity
                .HasOne(e => e.Workout)
                .WithMany()
                .HasForeignKey(e => e.WorkoutId)
                .OnDelete(DeleteBehavior.Cascade);

            entity
                .HasOne(e => e.UserProfile)
                .WithMany()
                .HasForeignKey(e => e.UserProfileId)
                .OnDelete(DeleteBehavior.Cascade);
        });
    }
}
