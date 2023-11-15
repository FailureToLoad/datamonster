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
	router := chi.NewRouter()
	SetDefaultMiddleware(router)
	SetCorsHandler(router)
	router.Get("/verify", verify)
	router.Post("/authorize", authorize)
	return router
}

func SetCorsHandler(r chi.Router) {
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8090"},
		AllowedMethods:   []string{"HEAD", "GET", "POST", "PATCH", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Origin", "X-Requested-With", "Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           3599, // Maximum value not ignored by any of major browsers
	})
	r.Use(c.Handler)
}

func SetAuthHandler(r chi.Router) {
	r.Use(authHandler)
}

func SetDefaultMiddleware(r chi.Router) {
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	timeout := middleware.Timeout(20 * time.Second)
	r.Use(timeout)
}

func verify(w http.ResponseWriter, r *http.Request) {
	log.Default().Println("Verification successful")
	MakeJsonResponse(w, http.StatusOK, "Success")
}

type AuthorizationRequest struct {
	Token string `json:"token"`
}

func authorize(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var body AuthorizationRequest
	err := DecodeJson(r.Body, &body)
	if err != nil {
		MakeJsonResponse(w, http.StatusBadRequest, "Invalid verification request body")
		return
	}

	// 5 day expiration
	expiresIn := time.Hour * 24 * 5
	cookie, err := client.SessionCookie(r.Context(), body.Token, expiresIn)
	if err != nil {
		log.Default().Printf("Error creating session cookie %s", err.Error())
		MakeJsonResponse(w, http.StatusBadRequest, "Unable to create session cookie")
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    cookie,
		MaxAge:   int(expiresIn.Seconds()),
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
		//SSL required for this guy
		//Secure:   true,
	})
	MakeJsonResponse(w, http.StatusOK, "Success")
}

type ctxUserIdKey string

const UserIdKey ctxUserIdKey = "userId"

func authHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		log.Default().Println("Authenticating request")
		cookie, err := r.Cookie("session")
		if err != nil {
			log.Default().Printf("Error retrieving session cookie %s", err.Error())
			MakeJsonResponse(w, http.StatusUnauthorized, "Unable to retrieve session cookie")
			return
		}
		ctx := r.Context()
		decoded, err := client.VerifySessionCookieAndCheckRevoked(r.Context(), cookie.Value)
		if err != nil {
			log.Default().Printf("Error verifying session cookie %s", err.Error())
			MakeJsonResponse(w, http.StatusUnauthorized, "Unable to verify session cookie")
			return
		}
		ctx = context.WithValue(ctx, UserIdKey, decoded.UID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}
