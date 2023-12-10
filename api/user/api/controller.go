package api

import (
	"datamonster/user"
	"datamonster/web"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

type Controller struct {
	svc *user.UserService
}

func NewController(s *user.UserService) *Controller {
	return &Controller{svc: s}
}

func (c Controller) RegisterRoutes(r chi.Router) {
	r.Group(func(r chi.Router) {
		web.SetDefaultMiddleware(r)
		web.SetCorsHandler(r)
		r.Post("/login", c.login)
		r.Post("/logout", c.logout)
		r.Post("/register", c.register)
	})
}

func (c Controller) login(w http.ResponseWriter, r *http.Request) {
	var creds LoginRequest
	err := web.DecodeJson(r.Body, &creds)
	if err != nil || !creds.Validate() {
		web.MakeJsonResponse(w, http.StatusBadRequest, "invalid request")
		return
	}

	userId, validationErr := c.svc.ValidateCredentials(r.Context(), creds.Username, creds.Password)
	if validationErr != nil {
		web.MakeJsonResponse(w, http.StatusUnauthorized, "login attempt failed")
		return
	}

	web.MakeAuthorizedResponse(w, userId)
}

func (c Controller) logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    "",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	})
	web.MakeJsonResponse(w, http.StatusNoContent, nil)
}

func (c Controller) register(w http.ResponseWriter, r *http.Request) {
	var creds RegisterRequest
	err := web.DecodeJson(r.Body, &creds)
	if err != nil || !creds.Validate() {
		web.MakeJsonResponse(w, http.StatusBadRequest, "invalid request")
		return
	}
	userId, err := c.svc.Register(r.Context(), creds.Username, creds.Password)
	if err != nil {
		log.Default().Printf("Error registering user %s", err.Error())
		web.MakeJsonResponse(w, http.StatusInternalServerError, "registration failed")
		return
	}

	web.MakeAuthorizedResponse(w, userId)
}
