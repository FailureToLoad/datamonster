using System.Text.Json.Serialization;

namespace Auth;

public class ValidateRequest
{
    [JsonPropertyName("token")]
    public string? Token { get; set; }
}