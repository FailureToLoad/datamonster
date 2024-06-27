using System;
using System.Collections.Generic;

namespace Context;

public partial class SettlementRecord
{
    public int Id { get; set; }

    public string Owner { get; set; } = null!;

    public string Name { get; set; } = null!;

    public short? SurvivalLimit { get; set; }

    public short? DepartingSurvival { get; set; }

    public short? CollectiveCognition { get; set; }

    public short? Year { get; set; }

    public virtual ICollection<SurvivorRecord> Survivors { get; set; } = new List<SurvivorRecord>();
}
