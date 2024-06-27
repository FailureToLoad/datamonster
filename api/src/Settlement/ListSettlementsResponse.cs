using System.Text.Json.Serialization;
using Context;

namespace Settlement;

public class ListSettlementsResponse()
{
    [JsonPropertyName("settlements")]
    public SettlementRecord[]? Settlements { get; set; }
    [JsonPropertyName("count")]
    public int? Count { get; set; }
}