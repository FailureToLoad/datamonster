package web

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func NewRouter() *chi.Mux {
	router := chi.NewRouter()
	SetDefaultMiddleware(router)
	SetCorsHandler(router)
	router.Get("/verify", verify)
	router.Post("/auth", authorize)
	return router
}

func SetCorsHandler(r chi.Router) {
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8090"},
		AllowedMethods:   []string{"HEAD", "GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Origin", "X-Requested-With", "Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           3599, // Maximum value not ignored by any of major browsers
	})
	r.Use(c.Handler)
}

func SetDefaultMiddleware(r chi.Router) {
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	timeout := middleware.Timeout(10 * time.Second)
	r.Use(timeout)
}

func verify(w http.ResponseWriter, r *http.Request) {
	log.Default().Println("Verification successful")
	MakeJsonResponse(w, http.StatusOK, "Success")
}

func authorize(w http.ResponseWriter, r *http.Request) {
	response, err := authorizeRequest(w, r)
	if err != nil {
		log.Default().Printf("Authorization failed: %s", err.Error())
		MakeJsonResponse(w, http.StatusUnauthorized, "authorization failure")
		return
	}

	MakeJsonResponse(w, http.StatusOK, response)
}
