package auth

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"os"
)

const (
	stateCookieName       = "oauth_state"
	stateCookieAge    int = 300
	sessionCookieName     = "dm_session"
)

func isSecureCookie() bool {
	return os.Getenv("COOKIE_SECURE") != "false"
}

func generateSessionID() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func expireCookie(w http.ResponseWriter, r *http.Request, sessions SessionStore) {
	cookie, err := r.Cookie(sessionCookieName)
	if err == nil && cookie.Value != "" {
		_ = sessions.Delete(r.Context(), cookie.Value)
	}
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    "",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   isSecureCookie(),
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	})
}

func setCookie(w http.ResponseWriter, sessionID string, ttl int) {
	if ttl <= 0 {
		ttl = 3600
	}

	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    sessionID,
		MaxAge:   ttl,
		HttpOnly: true,
		Secure:   isSecureCookie(),
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	})
}
