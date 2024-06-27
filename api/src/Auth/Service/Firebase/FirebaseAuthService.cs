using FirebaseAdmin.Auth;


namespace Auth.Service.Firebase;

public class FirebaseAuthService(IFirebase firebase) : IValidatable
{
    private readonly IFirebase auth = firebase;

    public async Task<ValidateCookieResult> ValidateCookie(string? sessionCookie)
    {
        if (string.IsNullOrEmpty(sessionCookie))
        {
            return new ValidateCookieResult(false);
        }
        var checkRevoked = true;
        try
        {
            var uid = await auth.VerifySessionCookieAsync(sessionCookie, checkRevoked);
            return new ValidateCookieResult(true, uid);
        }
        catch
        {
            return new ValidateCookieResult(false);
        }
    }

    public async Task<ValidateTokenResult> ValidateToken(string token)
    {
        if (string.IsNullOrEmpty(token))
        {
            return new ValidateTokenResult(false);
        }
        var checkRevoked = true;

        try
        {
            _ = await auth.VerifyIdTokenAsync(token, checkRevoked);
            var options = new SessionCookieOptions()
            {
                ExpiresIn = TimeSpan.FromDays(5),
            };

            var cookie = await auth.CreateSessionCookieAsync(token, options);
            return new ValidateTokenResult(true, cookie);
        }
        catch
        {
            return new ValidateTokenResult(false);
        }
    }

}