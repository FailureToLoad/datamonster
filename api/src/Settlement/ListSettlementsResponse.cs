using System.Text.Json.Serialization;
using datamonster.Context;

namespace Settlement;

public class ListSettlementsResponse()
{
    [JsonPropertyName("settlements")]
    public SettlementRecord[]? Settlements { get; set; }
    [JsonPropertyName("count")]
    public int? Count { get; set; }
}