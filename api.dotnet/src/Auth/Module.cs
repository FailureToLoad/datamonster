using Auth.Handler;
using Auth.Service;
using Microsoft.AspNetCore.Authentication;

namespace Auth;
public static class AuthModule
{
    private static readonly AuthController controller = new();
    public static IServiceCollection RegisterAuthModule(this IServiceCollection services)
    {
        services.AddCors(options =>
        {
            options.AddDefaultPolicy(builder =>
            {
                builder.WithOrigins("http://localhost:8090");
                builder.WithMethods("HEAD", "GET", "POST", "PATCH", "PUT", "DELETE");
                builder.WithHeaders("Origin", "X-Requested-With", "Accept", "Authorization", "Content-Type", "X-CSRF-Token");
                builder.WithExposedHeaders("Link");
                builder.AllowCredentials();
                builder.SetPreflightMaxAge(TimeSpan.FromSeconds(2520));
            });
        });
        services.AddSingleton<FirebaseAuthService>();
        services.AddAuthentication("firebase").AddScheme<AuthenticationSchemeOptions, FirebaseAuthHandler>("firebase", null);
        services.AddAuthorizationBuilder().AddPolicy("default", policy =>
        {
            policy.RequireClaim("id");
        });
        return services;
    }

    public static IEndpointRouteBuilder MapAuthEndpoints(this IEndpointRouteBuilder endpoints)
    {
        endpoints.MapPost("/auth", controller.Validate);
        return endpoints;
    }

}

