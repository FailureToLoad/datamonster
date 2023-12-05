using System;
using System.Collections.Generic;

namespace Context;

public partial class SurvivorRecord
{
    public int Id { get; set; }

    public int Settlement { get; set; }

    public string Name { get; set; } = null!;

    public short Birth { get; set; }

    public short Huntxp { get; set; }

    public short Survival { get; set; }

    public short Movement { get; set; }

    public short Accuracy { get; set; }

    public short Strength { get; set; }

    public short Evasion { get; set; }

    public short Luck { get; set; }

    public short Speed { get; set; }

    public short Insanity { get; set; }

    public short SystemicPressure { get; set; }

    public short Torment { get; set; }

    public short Lumi { get; set; }

    public char? Gender { get; set; }

    public virtual SettlementRecord SettlementNavigation { get; set; } = null!;
}
