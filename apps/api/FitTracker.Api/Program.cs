var builder = WebApplication.CreateBuilder(args);

// ============================================
// SERVICES CONFIGURATION
// ============================================

// Configure JSON serialization to be case-insensitive and handle enums as strings
builder.Services.ConfigureHttpJsonOptions(options =>
{
    options.SerializerOptions.PropertyNameCaseInsensitive = true;
    options.SerializerOptions.Converters.Add(
        new JsonStringEnumConverter(JsonNamingPolicy.CamelCase)
    );
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
kernelBuilder.AddOllamaChatCompletion(modelId: ollamaModel, endpoint: new Uri(ollamaEndpoint));
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
        policy => policy.AllowAnyOrigin().AllowAnyMethod().AllowAnyHeader()
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

    options.AddSecurityRequirement(_ =>
    {
        var schemeRef = new OpenApiSecuritySchemeReference("Bearer");
        return new OpenApiSecurityRequirement { [schemeRef] = [] };
    });
});

var app = builder.Build();

// ============================================
// MIDDLEWARE PIPELINE
// ============================================

app.UseCors("AllowAll");

// OpenAPI: expor o documento JSON (Swashbuckle) e documentação com Scalar na raiz
app.UseSwagger();
app.MapScalarApiReference(
    "/",
    options =>
    {
        options.WithTitle("FitTracker API");
        options.WithOpenApiRoutePattern("/swagger/{documentName}/swagger.json");
    }
);

// Authentication & Authorization
app.UseAuthentication();
app.UseAuthorization();

// Garantir que o banco de dados existe e que tabelas novas existem
using (var scope = app.Services.CreateScope())
{
    var db = scope.ServiceProvider.GetRequiredService<FitTrackerDbContext>();
    db.Database.EnsureCreated();

    // EnsureCreated() não adiciona tabelas a um DB já existente.
    // Garantir que WaterIntakeEntries existe (para DBs criados antes desta feature).
    if (db.Database.IsSqlite())
    {
        var conn = db.Database.GetDbConnection();
        conn.Open();
        try
        {
            using var checkCmd = conn.CreateCommand();
            checkCmd.CommandText =
                "SELECT name FROM sqlite_master WHERE type='table' AND name='WaterIntakeEntries'";
            var exists = checkCmd.ExecuteScalar();
            if (exists == null || exists == DBNull.Value)
            {
                using var createCmd = conn.CreateCommand();
                createCmd.CommandText =
                    """
                    CREATE TABLE WaterIntakeEntries (
                        Id TEXT NOT NULL PRIMARY KEY,
                        UserProfileId TEXT NOT NULL,
                        AmountMl INTEGER NOT NULL,
                        ConsumedAt TEXT NOT NULL,
                        Note TEXT
                    );
                    """;
                createCmd.ExecuteNonQuery();

                using var indexCmd = conn.CreateCommand();
                indexCmd.CommandText =
                    "CREATE INDEX IX_WaterIntakeEntries_UserProfileId_ConsumedAt ON WaterIntakeEntries (UserProfileId, ConsumedAt);";
                indexCmd.ExecuteNonQuery();
            }
        }
        finally
        {
            conn.Close();
        }
    }
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

// WaterIntake Feature
app.MapWaterIntakeEndpoints();

// AI Features
app.MapWorkoutPlanningEndpoints();

// Health check endpoint
app.MapGet("/health", () => Results.Ok(new { status = "healthy", timestamp = DateTime.UtcNow }))
    .WithTags("Health")
    .WithName("HealthCheck");

app.Run();
