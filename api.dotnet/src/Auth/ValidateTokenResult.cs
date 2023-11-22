namespace Auth;
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