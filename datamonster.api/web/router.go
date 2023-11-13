package web

import (
	"context"
	"log"
	"net/http"
	"time"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

var (
	client *auth.Client
)

func init() {
	routerInitCtx := context.Background()
	a, err := firebase.NewApp(routerInitCtx, nil)
	if err != nil {
		log.Fatal("Unable to create firebase app")
	}
	var authError error
	client, authError = a.Auth(routerInitCtx)
	if authError != nil {
		log.Fatal("Unable to create firebase auth client")
	}
}

func NewRouter() *chi.Mux {
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"HEAD", "GET", "POST", "PATCH", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           3599, // Maximum value not ignored by any of major browsers
	})
	router := chi.NewRouter()
	router.Use(c.Handler)

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	timeout := middleware.Timeout(20 * time.Second)
	router.Use(timeout)
	router.Use(authHandler)
	return router
}

type ctxUserIdKey string

const UserIdKey ctxUserIdKey = "userId"

func authHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		authHeader := r.Header.Get("Authorization")
		validToken, verifyErr := client.VerifyIDTokenAndCheckRevoked(ctx, authHeader)
		if verifyErr != nil {
			MakeJsonResponse(w, http.StatusUnauthorized, "Unauthorized")
			return
		} else {
			ctx = context.WithValue(ctx, UserIdKey, validToken.UID)
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}
