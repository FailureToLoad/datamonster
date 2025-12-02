package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"slices"

	"github.com/failuretoload/datamonster/store/valkey"
)

type contextKey string

const (
	sessionDataKey contextKey = "sessionData"
)

func AuthMiddleware(sessions *valkey.SessionStore) func(http.Handler) http.Handler {
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
	Subject           string `json:"sub"`
	Email             string `json:"email"`
	EmailVerified     bool   `json:"email_verified"`
	PreferredUsername string `json:"preferred_username"`
	Name              string `json:"name"`
	GivenName         string `json:"given_name"`
	FamilyName        string `json:"family_name"`
	RealmAccess       struct {
		Roles []string `json:"roles"`
	} `json:"realm_access"`
	ResourceAccess map[string]struct {
		Roles []string `json:"roles"`
	} `json:"resource_access"`
}

func (c *Claims) HasRealmRole(role string) bool {
	return slices.Contains(c.RealmAccess.Roles, role)
}

func (c *Claims) HasClientRole(clientID, role string) bool {
	if access, ok := c.ResourceAccess[clientID]; ok {
		return slices.Contains(access.Roles, role)
	}
	return false
}
