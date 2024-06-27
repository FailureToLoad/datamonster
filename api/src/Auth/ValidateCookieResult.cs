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