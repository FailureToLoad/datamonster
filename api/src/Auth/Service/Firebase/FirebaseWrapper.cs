using FirebaseAdmin;
using FirebaseAdmin.Auth;
using Google.Apis.Auth.OAuth2;

namespace Auth.Service.Firebase;
public class FirebaseWrapper : IFirebase
{
    private readonly FirebaseApp app;
    private readonly FirebaseAuth auth;
    public FirebaseWrapper()
    {
        app = FirebaseApp.Create(new AppOptions()
        {
            Credential = GoogleCredential.GetApplicationDefault(),
            ProjectId = Environment.GetEnvironmentVariable("FIREBASE_PROJECT"),

        });
        auth = FirebaseAuth.GetAuth(app);
    }
    public async Task<string> CreateSessionCookieAsync(string idToken, SessionCookieOptions options)
    {
        try
        {
            return await auth.CreateSessionCookieAsync(idToken, options);
        }
        catch (Exception e)
        {
            throw new CookieCreationException("Error creating session cookie.", e);
        }
    }

    public async Task<string> VerifyIdTokenAsync(string idToken, bool checkRevoked)
    {
        try
        {
            var token = await auth.VerifyIdTokenAsync(idToken, checkRevoked);
            return token.Uid;
        }
        catch (Exception e)
        {
            throw new VerificationException("Error verifying id token.", e);
        }
    }

    public async Task<string> VerifySessionCookieAsync(string sessionCookie, bool checkRevoked)
    {
        try
        {
            var token = await auth.VerifySessionCookieAsync(sessionCookie, checkRevoked);
            return token.Uid;
        }
        catch (Exception e)
        {
            throw new VerificationException("Error verifying session cookie.", e);
        }
    }
}