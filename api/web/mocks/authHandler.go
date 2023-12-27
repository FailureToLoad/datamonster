package mocks

import (
	"context"
	"datamonster/web"
	"net/http"
)

const TestUserId = 1

func AuthHandlerFake(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), web.UserIdKey, TestUserId)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}
