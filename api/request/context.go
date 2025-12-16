package request

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid/v5"
)

const (
	userIDKey contextKey = "userId"
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
