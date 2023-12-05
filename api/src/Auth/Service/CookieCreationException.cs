namespace Auth.Service;
[Serializable]
public class CookieCreationException : Exception
{
    public CookieCreationException()
    {
    }

    public CookieCreationException(string message)
        : base(message)
    {
    }

    public CookieCreationException(string message, Exception inner)
        : base(message, inner)
    {
    }
}