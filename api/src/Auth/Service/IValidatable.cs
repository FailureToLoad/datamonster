using Auth;
namespace Auth;
public interface IValidatable
{
    Task<ValidateCookieResult> ValidateCookie(string? sessionCookie);
    Task<ValidateTokenResult> ValidateToken(string token);
}