package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"slices"
	"strings"

	"github.com/failuretoload/datamonster/store/valkey"
)

type contextKey string

const (
	sessionDataKey contextKey = "sessionData"
)

type IntrospectionResponse struct {
	Active bool   `json:"active"`
	Sub    string `json:"sub"`
}

func (p *Provider) introspectToken(ctx context.Context, token string) (*IntrospectionResponse, error) {
	data := url.Values{}
	data.Set("token", token)
	data.Set("client_id", p.ClientID)
	data.Set("client_secret", p.ClientSecret)

	req, err := http.NewRequestWithContext(ctx, "POST", p.IntrospectURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var introspection IntrospectionResponse
	if err := json.NewDecoder(resp.Body).Decode(&introspection); err != nil {
		return nil, err
	}

	return &introspection, nil
}

func AuthMiddleware(provider *Provider, sessions *valkey.SessionStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie(sessionCookieName)
			if err != nil || cookie.Value == "" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			data, err := sessions.Get(r.Context(), cookie.Value)
			if err != nil || data == nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			var session SessionData
			if err := json.Unmarshal(data, &session); err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			introspection, err := provider.introspectToken(r.Context(), session.AccessToken)
			if err != nil || introspection == nil || !introspection.Active {
				_ = sessions.Delete(r.Context(), cookie.Value)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), sessionDataKey, &session)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func Session(r *http.Request) *SessionData {
	if val, ok := r.Context().Value(sessionDataKey).(*SessionData); ok {
		return val
	}
	return nil
}

func UserID(r *http.Request) string {
	if session := Session(r); session != nil {
		return session.UserID
	}
	return ""
}

type Claims struct {
	Subject           string   `json:"sub"`
	Email             string   `json:"email"`
	EmailVerified     bool     `json:"email_verified"`
	PreferredUsername string   `json:"preferred_username"`
	Name              string   `json:"name"`
	GivenName         string   `json:"given_name"`
	FamilyName        string   `json:"family_name"`
	Groups            []string `json:"groups"`
}

func (c *Claims) HasGroup(group string) bool {
	return slices.Contains(c.Groups, group)
}
