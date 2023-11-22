using FirebaseAdmin;
using FirebaseAdmin.Auth;
using Google.Apis.Auth.OAuth2;


namespace Auth.Service;

public record ValidateCookieResult
{
    public bool Success { get; }
    public string Id { get; } = "";

    public ValidateCookieResult(bool success, string? id = null)
    {
        Success = success;
        Id = id ?? "";
    }
}

public record ValidateTokenResult
{
    public bool Success { get; }
    public string? SessionCookie { get; } = "";

    public ValidateTokenResult(bool success, string? cookie = null)
    {
        Success = success;
        SessionCookie = cookie ?? "";
    }
}

public class FirebaseAuthService
{
    private readonly FirebaseAuth auth;
    private readonly FirebaseApp fbApp;
    public FirebaseAuthService()
    {
        fbApp = FirebaseApp.Create(new AppOptions()
        {
            Credential = GoogleCredential.GetApplicationDefault(),
            ProjectId = Environment.GetEnvironmentVariable("FIREBASE_PROJECT"),

        });
        auth = FirebaseAuth.GetAuth(fbApp);
    }

    public async Task<ValidateCookieResult> ValidateCookie(string? sessionCookie)
    {
        if (string.IsNullOrEmpty(sessionCookie))
        {
            return new ValidateCookieResult(false);
        }
        var checkRevoked = true;
        FirebaseToken? token;
        try
        {
            token = await auth.VerifySessionCookieAsync(sessionCookie, checkRevoked);
        }
        catch
        {
            return new ValidateCookieResult(false);
        }

        return new ValidateCookieResult(true, token.Uid);
    }

    public async Task<ValidateTokenResult> ValidateToken(string token)
    {
        if (string.IsNullOrEmpty(token))
        {
            return new ValidateTokenResult(false);
        }
        var checkRevoked = true;
        FirebaseToken? decodedToken;
        try
        {
            decodedToken = await auth.VerifyIdTokenAsync(token, checkRevoked);
        }
        catch
        {
            return new ValidateTokenResult(false);
        }

        var options = new SessionCookieOptions()
        {
            ExpiresIn = TimeSpan.FromDays(5),
        };

        string cookie = await auth.CreateSessionCookieAsync(token, options);
        return new ValidateTokenResult(true, cookie);
    }

}