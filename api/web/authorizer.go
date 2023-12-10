package web

import (
	"context"
	"datamonster/token"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v4"
)

type ctxUserIdKey string

type TokenResponse struct {
	UserId      int    `json:"userId"`
	BearerToken string `json:"token"`
}

const UserIdKey ctxUserIdKey = "userId"

func SetAuthHandler(r chi.Router) {
	r.Use(authHandler)
}

func MakeAuthorizedResponse(w http.ResponseWriter, userId int) {
	signedAuthTokenString, signedRefreshTokenString, err := token.GetSignedTokenStrings(userId)
	if err != nil {
		MakeJsonResponse(w, http.StatusInternalServerError, "login attempt failed")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    signedRefreshTokenString,
		Expires:  time.Now().Add(token.RefreshTokenExpirationPeriod),
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	})
	response := TokenResponse{
		UserId:      userId,
		BearerToken: fmt.Sprintf("Bearer %s", signedAuthTokenString),
	}
	MakeJsonResponse(w, http.StatusOK, response)
}

func authorizeRequest(w http.ResponseWriter, r *http.Request) (tokenResponse TokenResponse, err error) {
	cookie, cookieErr := r.Cookie("session")
	if cookieErr != nil {
		return tokenResponse, fmt.Errorf("unable to verify session: %w", err)
	}
	if cookie.Value == "" {
		return tokenResponse, fmt.Errorf("session expired")
	}

	reqToken := r.Header.Get("Authorization")
	if reqToken == "" {
		response, refreshMissingTokenErr := handleMissingAuthToken(w, r, cookie.Value)
		if refreshMissingTokenErr != nil {
			return TokenResponse{}, fmt.Errorf("unable to authorize request: %w", refreshMissingTokenErr)
		}
		return response, nil
	}

	if !strings.Contains(reqToken, "Bearer") {
		return tokenResponse, fmt.Errorf("unable to parse credentials")
	}

	splitToken := strings.Split(reqToken, "Bearer ")
	reqToken = splitToken[1]
	userId, validateTokenErr := token.ParseAndValidateToken(reqToken)
	if validateTokenErr != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			reqToken, err = handleExpiredAuthToken(w, r, cookie.Value)
		} else {
			MakeJsonResponse(w, http.StatusUnauthorized, "invalid token")
			return tokenResponse, fmt.Errorf("invalid token: %w", err)
		}
	}
	response := TokenResponse{
		UserId:      userId,
		BearerToken: reqToken,
	}
	return response, err
}

func handleMissingAuthToken(w http.ResponseWriter, r *http.Request, cookieValue string) (TokenResponse, error) {
	id, token, refreshErr := getNewAuthToken(w, r, cookieValue)
	if refreshErr != nil {
		MakeJsonResponse(w, http.StatusUnauthorized, "session expired")
		return TokenResponse{}, fmt.Errorf("session expired: %w", refreshErr)
	}
	response := TokenResponse{
		BearerToken: token,
		UserId:      id,
	}
	return response, nil
}

func handleExpiredAuthToken(w http.ResponseWriter, r *http.Request, cookieValue string) (string, error) {
	_, token, refreshErr := getNewAuthToken(w, r, cookieValue)
	if refreshErr != nil {
		MakeJsonResponse(w, http.StatusUnauthorized, "session expired")
		return "", fmt.Errorf("session expired: %w", refreshErr)
	}
	return token, nil
}

func authHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		tokenResponse, err := authorizeRequest(w, r)
		if err != nil {
			log.Default().Printf("auth handler failure: %s", err.Error())
			MakeJsonResponse(w, http.StatusUnauthorized, err.Error())
			return
		}

		ctx := context.WithValue(r.Context(), UserIdKey, tokenResponse.UserId)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}

func getNewAuthToken(w http.ResponseWriter, r *http.Request, refreshToken string) (int, string, error) {
	userId, err := token.ParseAndValidateToken(refreshToken)
	if err != nil {
		return 0, "", err
	}
	tokenString, err := token.GetSignedAuthTokenString(userId)
	if err != nil {
		return 0, "", err
	}
	return userId, tokenString, nil
}
