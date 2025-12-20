package request

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid/v5"
)

const (
	userIDKey        contextKey = "userId"
	correlationIDKey contextKey = "correlationID"
)

type (
	contextKey string
)

func UserID(ctx context.Context) string {
	if val, ok := ctx.Value(userIDKey).(string); ok {
		return val
	}
	return ""
}

func SetUserID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, userIDKey, id)
}

func IDParam(r *http.Request) (uuid.UUID, error) {
	rawID := chi.URLParam(r, "id")
	id, err := uuid.FromString(rawID)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func CorrelationID(ctx context.Context) string {
	if val, ok := ctx.Value(correlationIDKey).(string); ok {
		return val
	}
	return ""
}

func CorrelationIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if CorrelationID(ctx) == "" {
			unique, err := uuid.NewV4()
			if err == nil {
				ctx = context.WithValue(ctx, correlationIDKey, unique.String())
			}

			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		next.ServeHTTP(w, r)
	})
}
