package request

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

type (
	contextKey string
)

const (
	userIDKey     contextKey = "userId"
	userGroupsKey contextKey = "userGroups"
)

func AutheliaMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := r.Header.Get("Remote-User")

		if user == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, userIDKey, user)
		if groups := r.Header.Get("Remote-Groups"); groups != "" {
			ctx = context.WithValue(ctx, userGroupsKey, strings.Split(groups, ","))
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func UserID(r *http.Request) string {
	if val, ok := r.Context().Value(userIDKey).(string); ok {
		return val
	}
	return ""
}

func UserGroups(r *http.Request) []string {
	if val, ok := r.Context().Value(userGroupsKey).([]string); ok {
		return val
	}
	return nil
}

func IDParam(r *http.Request) (int, error) {
	id := chi.URLParam(r, "id")
	val, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		return 0, err
	}
	return int(val), nil
}
