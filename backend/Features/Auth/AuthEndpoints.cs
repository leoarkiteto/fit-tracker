namespace FitTracker.Api.Features.Auth;

public static class AuthEndpoints
{
    public static void MapAuthEndpoints(this WebApplication app)
    {
        var group = app.MapGroup("/api/auth").WithTags("Authentication");

        // POST /api/auth/register - Registrar novo usuário
        group
            .MapPost(
                "/register",
                async (
                    RegisterRequest request,
                    FitTrackerDbContext db,
                    IPasswordService passwordService,
                    IJwtService jwtService
                ) =>
                {
                    // Verificar se email já existe (comparação case-insensitive via normalização)
                    var emailLower = request.Email.Trim().ToLowerInvariant();
                    var existingUser = await db.Users.FirstOrDefaultAsync(u =>
                        u.Email == emailLower
                    );

                    if (existingUser is not null)
                        return Results.Conflict(new { message = "Email já está em uso" });

                    // Criar perfil de usuário
                    var profile = new UserProfile
                    {
                        Id = Guid.NewGuid(),
                        Name = request.Name,
                        CreatedAt = DateTime.UtcNow,
                        UpdatedAt = DateTime.UtcNow,
                    };

                    // Criar usuário
                    var user = new User
                    {
                        Id = Guid.NewGuid(),
                        Email = emailLower,
                        PasswordHash = passwordService.HashPassword(request.Password),
                        Name = request.Name,
                        CreatedAt = DateTime.UtcNow,
                        UserProfileId = profile.Id,
                    };

                    db.UserProfiles.Add(profile);
                    db.Users.Add(user);
                    await db.SaveChangesAsync();

                    // Gerar token
                    var token = jwtService.GenerateToken(user);
                    var expiresAt = DateTime.UtcNow.AddHours(24);

                    return Results.Created(
                        "/api/auth/me",
                        new AuthResponse(
                            token,
                            expiresAt,
                            new UserDto(user.Id, user.Email, user.Name, user.UserProfileId)
                        )
                    );
                }
            )
            .WithName("Register")
            .Produces<AuthResponse>(StatusCodes.Status201Created)
            .Produces(StatusCodes.Status409Conflict)
            .AllowAnonymous();

        // POST /api/auth/login - Login
        group
            .MapPost(
                "/login",
                async (
                    LoginRequest request,
                    FitTrackerDbContext db,
                    IPasswordService passwordService,
                    IJwtService jwtService
                ) =>
                {
                    var emailLower = request.Email.Trim().ToLowerInvariant();
                    var user = await db.Users.FirstOrDefaultAsync(u =>
                        u.Email == emailLower
                    );

                    if (user is null)
                        return Results.Unauthorized();

                    if (!passwordService.VerifyPassword(request.Password, user.PasswordHash))
                        return Results.Unauthorized();

                    // Atualizar último login
                    user.LastLoginAt = DateTime.UtcNow;
                    await db.SaveChangesAsync();

                    // Gerar token
                    var token = jwtService.GenerateToken(user);
                    var expiresAt = DateTime.UtcNow.AddHours(24);

                    return Results.Ok(
                        new AuthResponse(
                            token,
                            expiresAt,
                            new UserDto(user.Id, user.Email, user.Name, user.UserProfileId)
                        )
                    );
                }
            )
            .WithName("Login")
            .Produces<AuthResponse>(StatusCodes.Status200OK)
            .Produces(StatusCodes.Status401Unauthorized)
            .AllowAnonymous();

        // GET /api/auth/me - Obter usuário atual
        group
            .MapGet(
                "/me",
                async (ClaimsPrincipal user, FitTrackerDbContext db) =>
                {
                    var userIdClaim = user.FindFirst(ClaimTypes.NameIdentifier)?.Value;
                    if (userIdClaim is null || !Guid.TryParse(userIdClaim, out var userId))
                        return Results.Unauthorized();

                    var dbUser = await db.Users.FindAsync(userId);
                    if (dbUser is null)
                        return Results.NotFound();

                    return Results.Ok(
                        new UserDto(dbUser.Id, dbUser.Email, dbUser.Name, dbUser.UserProfileId)
                    );
                }
            )
            .WithName("GetCurrentUser")
            .Produces<UserDto>(StatusCodes.Status200OK)
            .Produces(StatusCodes.Status401Unauthorized)
            .RequireAuthorization();

        // POST /api/auth/change-password - Alterar senha
        group
            .MapPost(
                "/change-password",
                async (
                    ChangePasswordRequest request,
                    ClaimsPrincipal user,
                    FitTrackerDbContext db,
                    IPasswordService passwordService
                ) =>
                {
                    var userIdClaim = user.FindFirst(ClaimTypes.NameIdentifier)?.Value;
                    if (userIdClaim is null || !Guid.TryParse(userIdClaim, out var userId))
                        return Results.Unauthorized();

                    var dbUser = await db.Users.FindAsync(userId);
                    if (dbUser is null)
                        return Results.NotFound();

                    if (
                        !passwordService.VerifyPassword(
                            request.CurrentPassword,
                            dbUser.PasswordHash
                        )
                    )
                    {
                        return Results.BadRequest(new { message = "Senha atual incorreta" });
                    }

                    dbUser.PasswordHash = passwordService.HashPassword(request.NewPassword);
                    await db.SaveChangesAsync();

                    return Results.Ok(new { message = "Senha alterada com sucesso" });
                }
            )
            .WithName("ChangePassword")
            .Produces(StatusCodes.Status200OK)
            .Produces(StatusCodes.Status400BadRequest)
            .RequireAuthorization();

        // POST /api/auth/refresh - Renovar token
        group
            .MapPost(
                "/refresh",
                async (ClaimsPrincipal user, FitTrackerDbContext db, IJwtService jwtService) =>
                {
                    var userIdClaim = user.FindFirst(ClaimTypes.NameIdentifier)?.Value;
                    if (userIdClaim is null || !Guid.TryParse(userIdClaim, out var userId))
                        return Results.Unauthorized();

                    var dbUser = await db.Users.FindAsync(userId);
                    if (dbUser is null)
                        return Results.NotFound();

                    var token = jwtService.GenerateToken(dbUser);
                    var expiresAt = DateTime.UtcNow.AddHours(24);

                    return Results.Ok(
                        new AuthResponse(
                            token,
                            expiresAt,
                            new UserDto(dbUser.Id, dbUser.Email, dbUser.Name, dbUser.UserProfileId)
                        )
                    );
                }
            )
            .WithName("RefreshToken")
            .Produces<AuthResponse>(StatusCodes.Status200OK)
            .RequireAuthorization();
    }
}
