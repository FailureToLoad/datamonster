using System.Security.Claims;
using System.Text.Encodings.Web;
using Auth.Service;
using Microsoft.AspNetCore.Authentication;
using Microsoft.AspNetCore.Authentication.Cookies;
using Microsoft.Extensions.Options;

namespace Auth.Handler;

public class FirebaseAuthHandler : AuthenticationHandler<AuthenticationSchemeOptions>
{
    private readonly FirebaseAuthService authService;
    private readonly string sessionCookieName = "session";
    public FirebaseAuthHandler(IOptionsMonitor<AuthenticationSchemeOptions> options, ILoggerFactory logger, UrlEncoder encoder, FirebaseAuthService service) : base(options, logger, encoder)
    {
        authService = service;
    }

    protected override async Task HandleChallengeAsync(AuthenticationProperties properties)
    {
        Response.StatusCode = 401;
        await Response.WriteAsync("Unauthorized");
    }
    protected override async Task<AuthenticateResult> HandleAuthenticateAsync()
    {
        if (!Request.Cookies.ContainsKey(sessionCookieName))
        {
            Logger.LogInformation("No session cookie provided.");
            AuthenticationProperties props = new()
            {
                RedirectUri = "/login",
            };

            return AuthenticateResult.Fail("No session cookie provided.", props);
        }
        string? sessionCookie = Request.Cookies[sessionCookieName];
        ValidateCookieResult result = await authService.ValidateCookie(sessionCookie);
        if (!result.Success)
        {
            Logger.LogError("Error verifying session cookie");
            return AuthenticateResult.Fail("Error verifying session cookie");
        }
        var claims = new Claim[]
        {
            new("id", result.Id)
        };
        var identity = new ClaimsIdentity(claims, "firebase");
        var principle = new ClaimsPrincipal(identity);
        var ticket = new AuthenticationTicket(principle, Scheme.Name);
        return AuthenticateResult.Success(ticket);
    }
}