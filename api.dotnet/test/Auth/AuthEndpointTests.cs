using System.Runtime.CompilerServices;
using Auth;
using Auth.Service;
using Auth.Service.Firebase;
using FirebaseAdmin.Auth;
using Microsoft.AspNetCore.Http;
using Microsoft.AspNetCore.Http.HttpResults;
using Microsoft.Extensions.Logging;
using Moq;

namespace datamonster.test;

public class AuthEnpointTests
{
    internal ValidateTokenResult authenticated;
    internal ValidateTokenResult unauthorized;
    internal ValidateRequest validRequest;
    internal ValidateRequest invalidRequest;
    internal Mock<IFirebase> firebase;

    [SetUp]
    public void Setup()
    {

        authenticated = new ValidateTokenResult(true, "sessionCookie");
        unauthorized = new ValidateTokenResult(false);
        validRequest = new ValidateRequest()
        {
            Token = "validtoken"
        };
        invalidRequest = new ValidateRequest()
        {
            Token = ""
        };
        firebase = new Mock<IFirebase>(); ;
    }

    [Test]
    public async Task ValidToken_Set_SessionCookie()
    {
        firebase.Setup(x => x.VerifyIdTokenAsync(validRequest.Token, true))
            .ReturnsAsync("uid");
        firebase.Setup(x => x.CreateSessionCookieAsync(validRequest.Token, It.IsAny<SessionCookieOptions>()))
            .ReturnsAsync("sessionCookie");

        var authService = new FirebaseAuthService(firebase.Object);
        var controller = new AuthController();
        var context = new DefaultHttpContext();
        var logger = new Logger<AuthController>(new LoggerFactory());

        var result = await controller.Authenticate(validRequest, authService, logger, context);
        Assert.That(result, Is.InstanceOf<Ok>());
        var setCookieHeader = context.Response.Headers.SetCookie.ToString();
        Assert.That(setCookieHeader, Does.Contain("sessionCookie"));
    }

    [Test]
    public async Task InvalidToken_Return_ForbiddenStatus()
    {
        firebase.Setup(x => x.VerifyIdTokenAsync(validRequest.Token, true))
            .ThrowsAsync(new VerificationException());
        var authService = new FirebaseAuthService(firebase.Object);
        var controller = new AuthController();
        var context = new DefaultHttpContext();
        var logger = new Logger<AuthController>(new LoggerFactory());

        var result = await controller.Authenticate(validRequest, authService, logger, context);

        Assert.That(result, Is.InstanceOf<UnauthorizedHttpResult>());
    }

    [Test]
    public async Task NoToken_Return_ForbiddenResult()
    {
        var authService = new FirebaseAuthService(firebase.Object);
        var controller = new AuthController();
        var context = new DefaultHttpContext();
        var logger = new Logger<AuthController>(new LoggerFactory());

        var result = await controller.Authenticate(invalidRequest, authService, logger, context);

        Assert.That(result, Is.InstanceOf<UnauthorizedHttpResult>());
    }
}