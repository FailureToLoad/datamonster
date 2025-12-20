package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/failuretoload/datamonster/request"
	"github.com/failuretoload/datamonster/response"
)

type SessionData struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	IDToken      string `json:"id_token"`
	UserID       string `json:"user_id"`
}

type SessionStore interface {
	Delete(ctx context.Context, sessionID string) error
	Get(ctx context.Context, sessionID string) ([]byte, error)
	Set(ctx context.Context, sessionID string, data []byte, ttl time.Duration) error
	Exists(ctx context.Context, sessionID string) (bool, error)
}

type authorizationResponse struct {
	Active bool   `json:"active"`
	Sub    string `json:"sub"`
}

type refreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	IDToken      string `json:"id_token"`
}

type Authorizer struct {
	clientID      string
	clientSecret  string
	introspectURL string
	client        *http.Client
	sessions      SessionStore
	tokenURL      string
}

func newAuthorizer(id, secret, introspectURL, tokenURL string, sessions SessionStore) (Authorizer, error) {
	if id == "" {
		return Authorizer{}, ErrFieldMissing("clientID")
	}

	if secret == "" {
		return Authorizer{}, ErrFieldMissing("clientSecret")
	}

	if introspectURL == "" {
		return Authorizer{}, ErrFieldMissing("introspectURL")
	}

	if tokenURL == "" {
		return Authorizer{}, ErrFieldMissing("tokenURL")
	}

	if sessions == nil {
		return Authorizer{}, ErrFieldMissing("sessions")
	}

	httpClient := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 10,
			IdleConnTimeout:     90 * time.Second,
		},
	}

	return Authorizer{
		clientID:      id,
		clientSecret:  secret,
		introspectURL: introspectURL,
		tokenURL:      tokenURL,
		client:        httpClient,
		sessions:      sessions,
	}, nil
}

func (a Authorizer) AuthorizeRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		cookie, err := r.Cookie(sessionCookieName)
		if err != nil {
			response.Unauthorized(ctx, w, fmt.Errorf("unable to read cookie: %w", err))
			return
		}

		if cookie.Value == "" {
			response.Unauthorized(ctx, w, fmt.Errorf("missing session ID"))
			return
		}
		sessionID := cookie.Value
		data, err := a.sessions.Get(ctx, sessionID)
		if err != nil {
			response.Unauthorized(ctx, w, fmt.Errorf("unable to fetch session: %w", err))
			return
		}

		if data == nil {
			response.Unauthorized(ctx, w, fmt.Errorf("session not found: %s", sessionID))
			return
		}

		var session SessionData
		if err := json.Unmarshal(data, &session); err != nil {
			response.Unauthorized(ctx, w, fmt.Errorf("unable to unmarshal session data: %w", err))
			return
		}

		if !a.isActiveToken(ctx, session.AccessToken) {
			_ = a.sessions.Delete(ctx, cookie.Value)
			response.Unauthorized(ctx, w, fmt.Errorf("session expired"))
			return
		}

		age, err := a.refreshSession(ctx, session, sessionID)
		if err != nil {
			response.InternalServerError(ctx, w, fmt.Errorf("unable to refresh access token: %w", err))
			return
		}

		setCookie(w, sessionID, age)

		ctx = request.SetUserID(ctx, session.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (a Authorizer) isActiveToken(ctx context.Context, token string) bool {
	data := url.Values{}
	data.Set("token", token)
	data.Set("client_id", a.clientID)
	data.Set("client_secret", a.clientSecret)

	req, err := http.NewRequestWithContext(ctx, "POST", a.introspectURL, strings.NewReader(data.Encode()))
	if err != nil {
		slog.Error("creating introspection request", slog.Any("error", err))
		return false
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := a.client.Do(req)
	if err != nil {
		slog.Error("performing token introspection", slog.Any("error", err))
		return false
	}
	defer resp.Body.Close()

	var authorization authorizationResponse
	if err := json.NewDecoder(resp.Body).Decode(&authorization); err != nil {
		slog.Error("decoding authorization response", slog.Any("error", err))
		return false
	}

	return authorization.Active
}

func (a Authorizer) refreshSession(ctx context.Context, sessionData SessionData, sessionID string) (int, error) {
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", sessionData.RefreshToken)
	data.Set("client_id", a.clientID)
	data.Set("client_secret", a.clientSecret)

	req, err := http.NewRequestWithContext(ctx, "POST", a.tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return 0, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := a.client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return 0, fmt.Errorf("refresh token request failed: status=%d", resp.StatusCode)
	}

	var refreshed refreshTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&refreshed); err != nil {
		return 0, err
	}

	expiresInSeconds := refreshed.ExpiresIn
	if expiresInSeconds <= 0 {
		expiresInSeconds = 3600
	}

	ttl := time.Duration(expiresInSeconds) * time.Second

	if refreshed.AccessToken == "" {
		return 0, errors.New("refresh did not return an access token")
	}

	sessionData.AccessToken = refreshed.AccessToken
	if refreshed.RefreshToken != "" {
		sessionData.RefreshToken = refreshed.RefreshToken
	}
	if refreshed.IDToken != "" {
		sessionData.IDToken = refreshed.IDToken
	}

	refreshedSession, err := json.Marshal(sessionData)
	if err != nil {
		return 0, err
	}

	if err := a.sessions.Set(ctx, sessionID, refreshedSession, ttl); err != nil {
		return 0, err
	}

	return expiresInSeconds, nil
}
