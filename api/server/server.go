package server

import (
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	"github.com/unrolled/secure"
)

type AuthController interface {
	RegisterAuthRoutes(r *chi.Mux)
}

type AuthMiddleware interface {
	AuthorizeRequest(next http.Handler) http.Handler
}

func New(authController AuthController, authMiddleware AuthMiddleware, origins []string) (*http.Server, error) {
	if authController == nil {
		return nil, errors.New("auth controller cannot be nil")
	}

	if authMiddleware == nil {
		return nil, errors.New("auth middleware cannot be nil")
	}

	if len(origins) == 0 {
		return nil, errors.New("allowed origins are required")
	}

	ready := true
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Use(middleware.Timeout(10 * time.Second))

	router.Get("/heartbeat", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	router.Get("/ready", func(w http.ResponseWriter, _ *http.Request) {
		if !ready {
			w.WriteHeader(http.StatusServiceUnavailable)
			_, _ = w.Write([]byte("api is not ready"))
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ready"))
	})

	router.Mount("/auth", authRoutes(authController, origins))
	router.Mount("/api", protectedRoutes(authMiddleware.AuthorizeRequest, origins))

	return &http.Server{
		Addr:           "0.0.0.0:8080",
		Handler:        router,
		ReadTimeout:    15 * time.Second,
		WriteTimeout:   15 * time.Second,
		IdleTimeout:    60 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}, nil
}

func authRoutes(ac AuthController, allowedOrigins []string) *chi.Mux {
	r := baseRouter(allowedOrigins)
	ac.RegisterAuthRoutes(r)
	return r
}

func protectedRoutes(authMiddleware func(next http.Handler) http.Handler, allowedOrigins []string) *chi.Mux {
	r := baseRouter(allowedOrigins)
	r.Use(authMiddleware)

	r.Get("/me", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"authenticated":true}`))
	})

	return r
}

func baseRouter(allowedOrigins []string) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(middleware.Timeout(10 * time.Second))
	r.Use(httprate.LimitByIP(100, time.Minute))
	corsSettings := cors.New(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"HEAD", "GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Origin", "Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           3599,
	})
	r.Use(corsSettings.Handler)

	options := secure.Options{
		STSSeconds:            31536000,
		STSIncludeSubdomains:  true,
		STSPreload:            true,
		FrameDeny:             true,
		ForceSTSHeader:        true,
		ContentTypeNosniff:    true,
		BrowserXssFilter:      true,
		CustomBrowserXssValue: "0",
		ContentSecurityPolicy: "default-src 'self'; frame-ancestors 'none'",
	}
	r.Use(secure.New(options).Handler)

	r.Use(cacheControl)

	return r
}

func cacheControl(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		headers := rw.Header()
		headers.Set("Cache-Control", "no-cache, no-store, max-age=0, must-revalidate")
		headers.Set("Pragma", "no-cache")
		headers.Set("Expires", "0")
		next.ServeHTTP(rw, req)
	})
}
