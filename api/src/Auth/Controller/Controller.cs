using Auth.Service;
using Microsoft.AspNetCore.Mvc;

namespace Auth;
public class AuthController
{
    public async Task<IResult> Authenticate([FromBody] ValidateRequest request,
    [FromServices] IValidatable service,
    [FromServices] ILogger<AuthController> logger,
    HttpContext context)
    {
        ValidateTokenResult result;
        if (string.IsNullOrEmpty(request.Token))
        {
            logger.LogWarning("Token is empty.");
            return Results.Unauthorized();
        }

        result = await service.ValidateToken(request.Token);
        if (!result.Success)
        {
            logger.LogWarning("Token is invalid.");
            return Results.Unauthorized();
        }

        context.Response.Cookies.Append("session", result.SessionCookie!, new CookieOptions()
        {
            HttpOnly = true,
            SameSite = SameSiteMode.None,
            Secure = true,
            Expires = DateTimeOffset.UtcNow.AddDays(7)
        });
        return Results.Ok();
    }
}