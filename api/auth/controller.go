package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/failuretoload/datamonster/logger"
	"github.com/failuretoload/datamonster/response"
	"github.com/go-chi/chi/v5"
	"golang.org/x/oauth2"
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
		ctx := r.Context()
		stateCookie, err := r.Cookie(stateCookieName)
		if err != nil {
			response.BadRequest(ctx, w, fmt.Errorf("missing state cookie: %w", err))
			return
		}

		if r.URL.Query().Get("state") != stateCookie.Value {
			response.BadRequest(ctx, w, fmt.Errorf("state mismatch in callback handler"))
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
			response.BadRequest(ctx, w, fmt.Errorf("missing authorization code"))
			return
		}

		token, err := c.oauth2Config.Exchange(r.Context(), code)
		if err != nil {
			response.InternalServerError(ctx, w, fmt.Errorf("token exchange failed: %w", err))
			return
		}

		rawIDToken, ok := token.Extra("id_token").(string)
		if !ok {
			response.InternalServerError(ctx, w, fmt.Errorf("no id token in response"))
			return
		}

		idToken, err := c.verifier.Verify(r.Context(), rawIDToken)
		if err != nil {
			response.InternalServerError(ctx, w, fmt.Errorf("verifying ID token: %w", err))
			return
		}

		var claims Claims
		if err := idToken.Claims(&claims); err != nil {
			response.InternalServerError(ctx, w, fmt.Errorf("parsing token claims: %w", err))
			return
		}

		sessionID, err := generateSessionID()
		if err != nil {
			response.InternalServerError(ctx, w, fmt.Errorf("generating session id: %w", err))
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
			response.InternalServerError(ctx, w, fmt.Errorf("serializing session data: %w", err))
			return
		}

		ttl := time.Until(token.Expiry)
		if err := c.sessions.Set(r.Context(), sessionID, data, ttl); err != nil {
			response.InternalServerError(ctx, w, fmt.Errorf("saving session: %w", err))
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
		ctx := r.Context()
		cookie, err := r.Cookie(sessionCookieName)
		if err != nil {
			response.InternalServerError(ctx, w, err)
		}
		w.WriteHeader(http.StatusNoContent)

		if cookie.Value == "" {
			return
		}
		expireCookie(w)

		sessionBytes, err := c.sessions.Get(r.Context(), cookie.Value)
		if err != nil {
			logger.Warn(ctx, "fetching session for logout", logger.ErrorField(err))
			return
		}
		if len(sessionBytes) == 0 {
			logger.Warn(ctx, "logout: session is empty")
			return
		}

		var sessionData SessionData
		if err := json.Unmarshal(sessionBytes, &sessionData); err != nil {
			logger.Warn(ctx, "unmarshaling session for logout", logger.ErrorField(err))
			return
		}
		if sessionData.AccessToken != "" {
			c.revokeToken(r.Context(), sessionData.AccessToken)
		}

		if err := c.sessions.Delete(ctx, cookie.Value); err != nil {
			logger.Warn(ctx, "logout: failed to delete session", logger.ErrorField(err))
		}
	}
}

func (c *Controller) revokeToken(ctx context.Context, token string) {
	data := url.Values{}
	data.Set("token", token)
	data.Set("client_id", c.clientID)
	data.Set("client_secret", c.clientSecret)

	revocationURL := c.issuerURL + "/api/oidc/revocation"
	req, err := http.NewRequestWithContext(ctx, "POST", revocationURL, strings.NewReader(data.Encode()))
	if err != nil {
		logger.Warn(ctx, "creating token revocation request", logger.ErrorField(err))
		return
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	if _, err := c.httpClient.Do(req); err != nil {
		logger.Warn(ctx, "executing token revocation request", logger.ErrorField(err))
	}
}

func (c *Controller) loginHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		state, err := generateSessionID()
		if err != nil {
			response.InternalServerError(r.Context(), w, fmt.Errorf("failed to generate state: %w", err))
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     stateCookieName,
			Value:    state,
			MaxAge:   300,
			HttpOnly: true,
			Secure:   isSecureCookie(),
			SameSite: http.SameSiteLaxMode,
			Path:     "/",
		})

		authURL := c.oauth2Config.AuthCodeURL(state)
		http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
	}
}

func (c *Controller) checkHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
}
