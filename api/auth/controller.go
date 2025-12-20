package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/go-chi/chi/v5"
	"golang.org/x/oauth2"
)

const (
	stateCookieName   = "oauth_state"
	stateCookieAge    = 300
	sessionCookieName = "dm_session"
)

type Claims struct {
	Subject string `json:"sub"`
}

type Config struct {
	ClientID      string
	ClientSecret  string
	IssuerURL     string
	RedirectURL   string
	IntrospectURL string
	ClientURL     string
	TokenURL      string
	Sessions      SessionStore
}

func (c Config) Validate() error {
	if c.ClientID == "" {
		return ErrFieldMissing("clientID")
	}

	if c.ClientSecret == "" {
		return ErrFieldMissing("clientSecret")
	}

	if c.IssuerURL == "" {
		return ErrFieldMissing("issuerURL")
	}

	if c.RedirectURL == "" {
		return ErrFieldMissing("redirectURL")
	}

	if c.IntrospectURL == "" {
		return ErrFieldMissing("introspectURL")
	}

	if c.ClientURL == "" {
		return ErrFieldMissing("clientURL")
	}

	if c.Sessions == nil {
		return ErrFieldMissing("sessions")
	}

	if c.TokenURL == "" {
		return ErrFieldMissing("tokenURL")
	}

	return nil
}

func (c Config) Authorizer() (Authorizer, error) {
	return newAuthorizer(
		c.ClientID,
		c.ClientSecret,
		c.IntrospectURL,
		c.TokenURL,
		c.Sessions,
	)
}

type Controller struct {
	clientID     string
	clientSecret string
	issuerURL    string
	clientURL    string
	oauth2Config *oauth2.Config
	verifier     *oidc.IDTokenVerifier
	httpClient   *http.Client
	sessions     SessionStore
}

func (c Config) Controller() (*Controller, error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}

	ctx := context.Background()
	provider, err := oidc.NewProvider(ctx, c.IssuerURL)
	if err != nil {
		return nil, err
	}

	oauth2Config := &oauth2.Config{
		ClientID:     c.ClientID,
		ClientSecret: c.ClientSecret,
		RedirectURL:  c.RedirectURL,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email", "groups"},
	}

	verifier := provider.Verifier(&oidc.Config{ClientID: c.ClientID})

	return &Controller{
		clientID:     c.ClientID,
		clientSecret: c.ClientSecret,
		issuerURL:    c.IssuerURL,
		clientURL:    c.ClientURL,
		oauth2Config: oauth2Config,
		verifier:     verifier,
		sessions:     c.Sessions,
	}, nil
}

func (c Controller) RegisterAuthRoutes(r *chi.Mux) {
	r.Get("/login", c.loginHandler())
	r.Get("/callback", c.callbackHandler())
	r.Get("/logout", c.logoutHandler())
	r.Get("/check", c.checkHandler())
}

func (c *Controller) callbackHandler() http.HandlerFunc {
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
		token, err := c.oauth2Config.Exchange(r.Context(), code)
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

		idToken, err := c.verifier.Verify(r.Context(), rawIDToken)
		if err != nil {
			slog.Error("failed to verify ID token", "error", err)
			http.Error(w, "Failed to verify ID token", http.StatusInternalServerError)
			return
		}

		var claims Claims
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
		if err := c.sessions.Set(r.Context(), sessionID, data, ttl); err != nil {
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

		http.Redirect(w, r, c.clientURL, http.StatusTemporaryRedirect)
	}
}

func (c *Controller) logoutHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		revokeCookie(w, r, c.sessions)

		postLogoutRedirect := r.URL.Query().Get("redirect_uri")
		if postLogoutRedirect == "" {
			postLogoutRedirect = c.clientURL
		}

		logoutURL := c.issuerURL + "/api/oidc/revocation"
		fullLogoutURL := logoutURL + "?post_logout_redirect_uri=" + postLogoutRedirect
		http.Redirect(w, r, fullLogoutURL, http.StatusTemporaryRedirect)
	}
}

func (c *Controller) loginHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		authURL := c.oauth2Config.AuthCodeURL(state)
		http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
	}
}

func (c *Controller) checkHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(sessionCookieName)
		if err != nil || cookie.Value == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		exists, err := c.sessions.Exists(r.Context(), cookie.Value)
		if err != nil || !exists {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
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
