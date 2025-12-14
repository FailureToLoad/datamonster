package auth

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/failuretoload/datamonster/request"
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

type Authorization struct {
	Active bool   `json:"active"`
	Sub    string `json:"sub"`
}

type Authorizer struct {
	clientID      string
	clientSecret  string
	introspectURL string
	client        *http.Client
	sessions      SessionStore
}

func NewAuthorizer(id, secret, introspect string, store SessionStore) (Authorizer, error) {
	var a Authorizer
	if id == "" {
		return a, errors.New("client ID is required")
	}

	if secret == "" {
		return a, errors.New("client secret is required")
	}

	if introspect == "" {
		return a, errors.New("introspect url is required")
	}

	if store == nil {
		return a, errors.New("session storage is required")
	}

	httpClient := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 10,
			IdleConnTimeout:     90 * time.Second,
		},
	}

	a = Authorizer{
		clientID:      id,
		clientSecret:  secret,
		introspectURL: introspect,
		client:        httpClient,
		sessions:      store,
	}

	return a, nil
}

func (a Authorizer) AuthorizeRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(sessionCookieName)
		if err != nil || cookie.Value == "" {
			unauthorized(w, "unable to read cookie", err)
			return
		}

		data, err := a.sessions.Get(r.Context(), cookie.Value)
		if err != nil || data == nil {
			unauthorized(w, "unable to fetch session", err)
			return
		}

		var session SessionData
		if err := json.Unmarshal(data, &session); err != nil {
			unauthorized(w, "unable to unmarshal sessions data", err)
			return
		}

		validation, err := a.checkToken(r.Context(), session.AccessToken)
		if err != nil || validation == nil || !validation.Active {
			_ = a.sessions.Delete(r.Context(), cookie.Value)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := request.SetUserID(r.Context(), session.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (a Authorizer) checkToken(ctx context.Context, token string) (*Authorization, error) {
	data := url.Values{}
	data.Set("token", token)
	data.Set("client_id", a.clientID)
	data.Set("client_secret", a.clientSecret)

	req, err := http.NewRequestWithContext(ctx, "POST", a.introspectURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var authorization Authorization
	if err := json.NewDecoder(resp.Body).Decode(&authorization); err != nil {
		return nil, err
	}

	return &authorization, nil
}

func unauthorized(w http.ResponseWriter, reason string, err error) {
	slog.Error(reason, slog.Any("error", err))
	http.Error(w, "Unauthorized", http.StatusUnauthorized)
}
