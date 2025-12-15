package testenv

import (
	"net/http"

	"github.com/failuretoload/datamonster/request"
	"github.com/go-chi/chi/v5"
)

type AuthControllerFake struct{}

func (m AuthControllerFake) RegisterAuthRoutes(_ *chi.Mux) {}

type AuthorizerFake struct {
	authorized bool
	userID     string
}

func NewAuthorizerFake() AuthorizerFake {
	return AuthorizerFake{
		authorized: true,
	}
}

func (m *AuthorizerFake) Authorized() {
	m.authorized = true
}

func (m *AuthorizerFake) Unauthorized() {
	m.authorized = false
}

func (m *AuthorizerFake) ExpectUserID(id string) {
	m.userID = id
}

func (m *AuthorizerFake) AuthorizeRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !m.authorized {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		}
		ctx := request.SetUserID(r.Context(), m.userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
