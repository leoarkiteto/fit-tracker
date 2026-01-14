using System.Text;
using System.Text.Json;
using System.Text.Json.Serialization;
using FitTracker.Api.Features.AI.WorkoutPlanning;
using FitTracker.Api.Features.Auth;
using FitTracker.Api.Features.Bioimpedance;
using FitTracker.Api.Features.Profiles;
using FitTracker.Api.Features.Workouts;
using FitTracker.Api.Shared.Data;
using Microsoft.AspNetCore.Authentication.JwtBearer;
using Microsoft.EntityFrameworkCore;
using Microsoft.IdentityModel.Tokens;
using Microsoft.OpenApi.Models;
using Microsoft.SemanticKernel;

var builder = WebApplication.CreateBuilder(args);

// ============================================
// SERVICES CONFIGURATION
// ============================================

// Configure JSON serialization to be case-insensitive and handle enums as strings
builder.Services.ConfigureHttpJsonOptions(options =>
{
    options.SerializerOptions.PropertyNameCaseInsensitive = true;
    options.SerializerOptions.Converters.Add(new JsonStringEnumConverter(JsonNamingPolicy.CamelCase));
});

// SQLite Database
builder.Services.AddDbContext<FitTrackerDbContext>(options =>
    options.UseSqlite(
        builder.Configuration.GetConnectionString("DefaultConnection")
            ?? "Data Source=fittracker.db"
    )
);

// Feature Services - Auth
builder.Services.AddScoped<IPasswordService, PasswordService>();
builder.Services.AddScoped<IJwtService, JwtService>();

// Feature Services - AI (Semantic Kernel + Ollama)
var ollamaEndpoint = builder.Configuration["AI:Ollama:Endpoint"] ?? "http://localhost:11434";
var ollamaModel = builder.Configuration["AI:Ollama:Model"] ?? "llama3.2";

var kernelBuilder = Kernel.CreateBuilder();

#pragma warning disable SKEXP0070 // Ollama connector is experimental
kernelBuilder.AddOllamaChatCompletion(
    modelId: ollamaModel,
    endpoint: new Uri(ollamaEndpoint)
);
#pragma warning restore SKEXP0070

var kernel = kernelBuilder.Build();
builder.Services.AddSingleton(kernel);
builder.Services.AddScoped<IWorkoutPlanningAgent, WorkoutPlanningAgent>();

// JWT Authentication
var jwtSecret = builder.Configuration["Jwt:Secret"] ?? "FitTrackerSuperSecretKey2024!@#$%^&*()";
var jwtIssuer = builder.Configuration["Jwt:Issuer"] ?? "FitTracker";
var jwtAudience = builder.Configuration["Jwt:Audience"] ?? "FitTrackerApp";

builder
    .Services.AddAuthentication(options =>
    {
        options.DefaultAuthenticateScheme = JwtBearerDefaults.AuthenticationScheme;
        options.DefaultChallengeScheme = JwtBearerDefaults.AuthenticationScheme;
    })
    .AddJwtBearer(options =>
    {
        options.TokenValidationParameters = new TokenValidationParameters
        {
            ValidateIssuerSigningKey = true,
            IssuerSigningKey = new SymmetricSecurityKey(Encoding.UTF8.GetBytes(jwtSecret)),
            ValidateIssuer = true,
            ValidIssuer = jwtIssuer,
            ValidateAudience = true,
            ValidAudience = jwtAudience,
            ValidateLifetime = true,
            ClockSkew = TimeSpan.Zero,
        };
    });

builder.Services.AddAuthorization();

// CORS - Permitir acesso do app mobile
builder.Services.AddCors(options =>
{
    options.AddPolicy(
        "AllowAll",
        policy =>
        {
            policy.AllowAnyOrigin().AllowAnyMethod().AllowAnyHeader();
        }
    );
});

// Swagger/OpenAPI
builder.Services.AddEndpointsApiExplorer();
builder.Services.AddSwaggerGen(options =>
{
    options.SwaggerDoc(
        "v1",
        new OpenApiInfo
        {
            Title = "FitTracker API",
            Version = "v1",
            Description =
                "API para o aplicativo FitTracker - Gerenciamento de treinos e bioimpedância",
        }
    );

    // Configurar Swagger para suportar JWT
    options.AddSecurityDefinition(
        "Bearer",
        new OpenApiSecurityScheme
        {
            Name = "Authorization",
            Type = SecuritySchemeType.ApiKey,
            Scheme = "Bearer",
            BearerFormat = "JWT",
            In = ParameterLocation.Header,
            Description = "Insira o token JWT no formato: Bearer {seu_token}",
        }
    );

    options.AddSecurityRequirement(
        new OpenApiSecurityRequirement
        {
            {
                new OpenApiSecurityScheme
                {
                    Reference = new OpenApiReference
                    {
                        Type = ReferenceType.SecurityScheme,
                        Id = "Bearer",
                    },
                },
                Array.Empty<string>()
            },
        }
    );
});

var app = builder.Build();

// ============================================
// MIDDLEWARE PIPELINE
// ============================================

app.UseCors("AllowAll");

// Swagger sempre disponível para facilitar desenvolvimento
app.UseSwagger();
app.UseSwaggerUI(options =>
{
    options.SwaggerEndpoint("/swagger/v1/swagger.json", "FitTracker API v1");
    options.RoutePrefix = string.Empty; // Swagger na raiz
});

// Authentication & Authorization
app.UseAuthentication();
app.UseAuthorization();

// Garantir que o banco de dados existe
using (var scope = app.Services.CreateScope())
{
    var db = scope.ServiceProvider.GetRequiredService<FitTrackerDbContext>();
    db.Database.EnsureCreated();
}

// ============================================
// FEATURE ENDPOINTS (Vertical Slice)
// ============================================

// Auth Feature
app.MapAuthEndpoints();

// Profiles Feature
app.MapProfileEndpoints();

// Workouts Feature
app.MapWorkoutEndpoints();
app.MapCompletedWorkoutEndpoints();

// Bioimpedance Feature
app.MapBioimpedanceEndpoints();

// AI Features
app.MapWorkoutPlanningEndpoints();

// Health check endpoint
app.MapGet("/health", () => Results.Ok(new { status = "healthy", timestamp = DateTime.UtcNow }))
    .WithTags("Health")
    .WithName("HealthCheck");

app.Run();
