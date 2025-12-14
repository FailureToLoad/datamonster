package request

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
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

func IDParam(r *http.Request) (int, error) {
	id := chi.URLParam(r, "id")
	val, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		return 0, err
	}
	return int(val), nil
}
