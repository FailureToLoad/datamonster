using Auth.Service;
using Microsoft.AspNetCore.Mvc;

namespace Auth;
public class AuthController
{
    public async Task<IResult> Validate([FromBody] ValidateRequest request,
    [FromServices] FirebaseAuthService service,
    [FromServices] ILogger<AuthController> logger,
    HttpContext context)
    {
        ValidateTokenResult result;
        try
        {
            if (request.Token == null)
            {
                logger.LogWarning("Token is null.");
                return Results.Unauthorized();
            }

            result = await service.ValidateToken(request.Token);
        }
        catch
        {
            logger.LogWarning("Unable to get cookie from token.");
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