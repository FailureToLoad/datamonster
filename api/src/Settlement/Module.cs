namespace Settlement;

public static class SettlementModule
{
    private static readonly SettlementController controller = new();
    public static IServiceCollection RegisterSettlementModule(this IServiceCollection services)
    {
        services.AddScoped<SettlementRepository>();
        return services;
    }

    public static IEndpointRouteBuilder MapSettlementEndpoints(this IEndpointRouteBuilder endpoints)
    {
        endpoints.MapGet("/settlement", controller.GetSettlementsForUser).RequireAuthorization("read:settlements");
        endpoints.MapGet("/settlements/{id}", controller.GetSettlement).RequireAuthorization("read:settlements");
        endpoints.MapPost("/settlement", controller.CreateSettlement).RequireAuthorization("create:settlements");
        endpoints.MapPut("/settlements/{id}", controller.UpdateSettlement).RequireAuthorization("update:settlements");
        return endpoints;
    }
}