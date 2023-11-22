using FirebaseAdmin.Auth;

namespace Auth;
public interface IFirebase
{
    Task<string> VerifySessionCookieAsync(string sessionCookie, bool checkRevoked);
    Task<string> VerifyIdTokenAsync(string idToken, bool checkRevoked);
    Task<string> CreateSessionCookieAsync(string idToken, SessionCookieOptions options);
}