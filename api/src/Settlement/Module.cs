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
        endpoints.MapGet("/settlement", controller.GetSettlementsForUser).RequireAuthorization("default");
        endpoints.MapGet("/settlements/{id}", controller.GetSettlement).RequireAuthorization("default");
        endpoints.MapPost("/settlement", controller.CreateSettlement).RequireAuthorization("default");
        endpoints.MapPut("/settlements/{id}", controller.UpdateSettlement).RequireAuthorization("default");
        endpoints.MapDelete("/settlements/{id}", controller.DeleteSettlement).RequireAuthorization("default");
        return endpoints;
    }
}