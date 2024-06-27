using System;
using System.Collections.Generic;
using Microsoft.EntityFrameworkCore;

namespace Context;

public partial class RecordsContext : DbContext
{
    public RecordsContext()
    {
    }

    public RecordsContext(DbContextOptions<RecordsContext> options)
        : base(options)
    {
    }

    public virtual DbSet<SettlementRecord> Settlements { get; set; }

    public virtual DbSet<SurvivorRecord> Survivors { get; set; }

    protected override void OnConfiguring(DbContextOptionsBuilder optionsBuilder)
    {
        string? conn = Environment.GetEnvironmentVariable("EF_STRING");
        if (!optionsBuilder.IsConfigured)
        {
            optionsBuilder.UseNpgsql(conn);
        }
    }

    protected override void OnModelCreating(ModelBuilder modelBuilder)
    {
        modelBuilder.Entity<SettlementRecord>(entity =>
        {
            entity.HasKey(e => e.Id).HasName("settlement_pkey");

            entity.ToTable("settlement");

            entity.Property(e => e.Id).HasColumnName("id");
            entity.Property(e => e.CollectiveCognition)
                .HasDefaultValue((short)0)
                .HasColumnName("collective_cognition");
            entity.Property(e => e.DepartingSurvival)
                .HasDefaultValue((short)0)
                .HasColumnName("departing_survival");
            entity.Property(e => e.Name)
                .HasMaxLength(50)
                .HasColumnName("name");
            entity.Property(e => e.Owner)
                .HasMaxLength(128)
                .HasColumnName("owner");
            entity.Property(e => e.SurvivalLimit)
                .HasDefaultValue((short)0)
                .HasColumnName("survival_limit");
            entity.Property(e => e.Year)
                .HasDefaultValue((short)0)
                .HasColumnName("year");
        });

        modelBuilder.Entity<SurvivorRecord>(entity =>
        {
            entity.HasKey(e => e.Id).HasName("survivor_pkey");

            entity.ToTable("survivor");

            entity.Property(e => e.Id).HasColumnName("id");
            entity.Property(e => e.Accuracy)
                .HasDefaultValue((short)0)
                .HasColumnName("accuracy");
            entity.Property(e => e.Birth)
                .HasDefaultValue((short)0)
                .HasColumnName("birth");
            entity.Property(e => e.Evasion)
                .HasDefaultValue((short)0)
                .HasColumnName("evasion");
            entity.Property(e => e.Gender)
                .HasMaxLength(1)
                .HasColumnName("gender");
            entity.Property(e => e.Huntxp)
                .HasDefaultValue((short)0)
                .HasColumnName("huntxp");
            entity.Property(e => e.Insanity)
                .HasDefaultValue((short)0)
                .HasColumnName("insanity");
            entity.Property(e => e.Luck)
                .HasDefaultValue((short)0)
                .HasColumnName("luck");
            entity.Property(e => e.Lumi)
                .HasDefaultValue((short)0)
                .HasColumnName("lumi");
            entity.Property(e => e.Movement)
                .HasDefaultValue((short)5)
                .HasColumnName("movement");
            entity.Property(e => e.Name)
                .HasMaxLength(50)
                .HasColumnName("name");
            entity.Property(e => e.Settlement).HasColumnName("settlement");
            entity.Property(e => e.Speed)
                .HasDefaultValue((short)0)
                .HasColumnName("speed");
            entity.Property(e => e.Strength)
                .HasDefaultValue((short)0)
                .HasColumnName("strength");
            entity.Property(e => e.Survival)
                .HasDefaultValue((short)1)
                .HasColumnName("survival");
            entity.Property(e => e.SystemicPressure)
                .HasDefaultValue((short)0)
                .HasColumnName("systemic_pressure");
            entity.Property(e => e.Torment)
                .HasDefaultValue((short)0)
                .HasColumnName("torment");

            entity.HasOne(d => d.SettlementNavigation).WithMany(p => p.Survivors)
                .HasForeignKey(d => d.Settlement)
                .HasConstraintName("fk_settlement_id");
        });

        OnModelCreatingPartial(modelBuilder);
    }

    partial void OnModelCreatingPartial(ModelBuilder modelBuilder);
}
