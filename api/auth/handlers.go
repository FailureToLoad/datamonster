package auth

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/failuretoload/datamonster/store/valkey"
)

const (
	stateCookieName   = "oauth_state"
	stateCookieAge    = 300
	sessionCookieName = "dm_session"
)

type SessionData struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	IDToken      string `json:"id_token"`
	UserID       string `json:"user_id"`
}

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

func (p *Provider) LoginHandler(w http.ResponseWriter, r *http.Request) {
	state, err := generateSessionID()
	if err != nil {
		http.Error(w, "Failed to generate state", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     stateCookieName,
		Value:    state,
		MaxAge:   stateCookieAge,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	})

	authURL := p.OAuth2Config.AuthCodeURL(state)
	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

func (p *Provider) CallbackHandler(sessions *valkey.SessionStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Debug("callback handler invoked", "url", r.URL.String())

		stateCookie, err := r.Cookie(stateCookieName)
		if err != nil {
			slog.Error("missing state cookie", "error", err)
			http.Error(w, "Missing state cookie", http.StatusBadRequest)
			return
		}

		if r.URL.Query().Get("state") != stateCookie.Value {
			http.Error(w, "Invalid state parameter", http.StatusBadRequest)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     stateCookieName,
			Value:    "",
			MaxAge:   -1,
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteLaxMode,
			Path:     "/",
		})

		code := r.URL.Query().Get("code")
		if code == "" {
			http.Error(w, "Missing authorization code", http.StatusBadRequest)
			return
		}

		slog.Debug("exchanging code for token", "code_prefix", code[:10])
		token, err := p.OAuth2Config.Exchange(r.Context(), code)
		slog.Debug("exchange complete", "error", err)
		if err != nil {
			slog.Error("token exchange failed", "error", err)
			http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
			return
		}

		rawIDToken, ok := token.Extra("id_token").(string)
		if !ok {
			http.Error(w, "No id_token in response", http.StatusInternalServerError)
			return
		}

		idToken, err := p.Verifier.Verify(r.Context(), rawIDToken)
		if err != nil {
			slog.Error("failed to verify ID token", "error", err)
			http.Error(w, "Failed to verify ID token", http.StatusInternalServerError)
			return
		}

		var claims struct {
			Subject string `json:"sub"`
		}
		if err := idToken.Claims(&claims); err != nil {
			http.Error(w, "Failed to parse claims", http.StatusInternalServerError)
			return
		}

		sessionID, err := generateSessionID()
		if err != nil {
			http.Error(w, "Failed to generate session", http.StatusInternalServerError)
			return
		}

		sessionData := SessionData{
			AccessToken:  token.AccessToken,
			RefreshToken: token.RefreshToken,
			IDToken:      rawIDToken,
			UserID:       claims.Subject,
		}

		data, err := json.Marshal(sessionData)
		if err != nil {
			http.Error(w, "Failed to serialize session", http.StatusInternalServerError)
			return
		}

		ttl := time.Until(token.Expiry)
		if err := sessions.Set(r.Context(), sessionID, data, ttl); err != nil {
			http.Error(w, "Failed to store session", http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     sessionCookieName,
			Value:    sessionID,
			MaxAge:   int(ttl.Seconds()),
			HttpOnly: true,
			Secure:   isSecureCookie(),
			SameSite: http.SameSiteLaxMode,
			Path:     "/",
		})

		frontendURL := os.Getenv("FRONTEND_URL")
		if frontendURL == "" {
			frontendURL = "/"
		}
		http.Redirect(w, r, frontendURL, http.StatusTemporaryRedirect)
	}
}

func (p *Provider) LogoutHandler(sessions *valkey.SessionStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		postLogoutRedirect := r.URL.Query().Get("redirect_uri")
		if postLogoutRedirect == "" {
			postLogoutRedirect = os.Getenv("FRONTEND_URL")
			if postLogoutRedirect == "" {
				postLogoutRedirect = "/"
			}
		}

		logoutURL := p.IssuerURL + "/api/oidc/revocation"
		fullLogoutURL := logoutURL + "?post_logout_redirect_uri=" + postLogoutRedirect
		http.Redirect(w, r, fullLogoutURL, http.StatusTemporaryRedirect)
	}
}

func (p *Provider) RegisterRoutes(sessions *valkey.SessionStore, r interface {
	Get(string, http.HandlerFunc)
	Post(string, http.HandlerFunc)
},
) {
	r.Get("/login", p.LoginHandler)
	r.Get("/callback", p.CallbackHandler(sessions))
	r.Get("/logout", p.LogoutHandler(sessions))
	r.Get("/check", CheckHandler(sessions))
}

func CheckHandler(sessions *valkey.SessionStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(sessionCookieName)
		if err != nil || cookie.Value == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		exists, err := sessions.Exists(r.Context(), cookie.Value)
		if err != nil || !exists {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
