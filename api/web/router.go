package web

import (
	"context"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/supertokens/supertokens-golang/supertokens"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	Mux *chi.Mux
}

func (s Server) Start() {
	log.Default().Println("Starting server on port 8080")
	err := http.ListenAndServe(":8080", corsMiddleware(supertokens.Middleware(s.Mux)))
	if err != nil {
		log.Default().Println(err)
	}
}

func NewServer() Server {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Use(middleware.Timeout(10 * time.Second))
	router.Use(userIdExtractor)
	return Server{Mux: router}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, r *http.Request) {
		response.Header().Set("Access-Control-Allow-Origin", "http://localhost:8090")
		response.Header().Set("Access-Control-Allow-Credentials", "true")
		if r.Method == "OPTIONS" {
			// we add content-type + other headers used by SuperTokens
			response.Header().Set("Access-Control-Allow-Headers",
				strings.Join(append([]string{"Content-Type"},
					supertokens.GetAllCORSHeaders()...), ","))
			response.Header().Set("Access-Control-Allow-Methods", "*")
			response.Write([]byte(""))
		} else {
			next.ServeHTTP(response, r)
		}
	})
}

type ctxUserIdKey string

const UserIdKey ctxUserIdKey = "userId"

func userIdExtractor(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionContainer := session.GetSessionFromRequestContext(r.Context())
		userID := sessionContainer.GetUserID()
		if userID != "" {
			ctx := context.WithValue(r.Context(), UserIdKey, userID)
			r = r.WithContext(ctx)
		}
		next.ServeHTTP(w, r)
	})
}
