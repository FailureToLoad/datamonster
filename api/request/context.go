package request

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid/v5"
)

const (
	userIDKey        contextKey = "userId"
	correlationIDKey contextKey = "correlationID"
	settlementIDKey  contextKey = "settlementID"
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

func SettlementID(ctx context.Context) uuid.UUID {
	if val, ok := ctx.Value(settlementIDKey).(uuid.UUID); ok {
		return val
	}

	return uuid.Nil
}

func SetSettlementID(ctx context.Context, id uuid.UUID) context.Context {
	return context.WithValue(ctx, settlementIDKey, id)
}

func SettlementIDFromURL(r *http.Request) (uuid.UUID, error) {
	rawID := chi.URLParam(r, "id")
	id, err := uuid.FromString(rawID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("unable to parse settlement id from URL: %w", err)
	}

	if id == uuid.Nil {
		return id, fmt.Errorf("settlement id is required")
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
			unique, err := uuid.NewV7()
			if err == nil {
				ctx = context.WithValue(ctx, correlationIDKey, unique.String())
			}

			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		next.ServeHTTP(w, r)
	})
}
