using System.ComponentModel.DataAnnotations;

namespace FitTracker.Api.Features.Auth;

public record RegisterRequest(
    [Required] [EmailAddress] string Email,
    [Required] [MinLength(6)] string Password,
    [Required] string Name
);

public record LoginRequest(
    [Required] [EmailAddress] string Email,
    [Required] string Password
);

public record AuthResponse(
    string Token,
    DateTime ExpiresAt,
    UserDto User
);

public record UserDto(
    Guid Id,
    string Email,
    string Name,
    Guid? ProfileId
);

public record ChangePasswordRequest(
    [Required] string CurrentPassword,
    [Required] [MinLength(6)] string NewPassword
);
