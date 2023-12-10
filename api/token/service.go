package token

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const (
	AuthTokenExpirationPeriod    = time.Hour * 1
	RefreshTokenExpirationPeriod = time.Hour * 24
)

func GetSignedAuthTokenString(userId int) (string, error) {
	key, err := getSigningKey()
	if err != nil {
		return "", err
	}
	authTokenClaims := getClaims(fmt.Sprintf("%d", userId), time.Now().Add(AuthTokenExpirationPeriod))
	authToken := jwt.NewWithClaims(jwt.SigningMethodHS256, authTokenClaims)
	signedAuthTokenString, authSignErr := authToken.SignedString(key)
	if authSignErr != nil {
		return "", authSignErr
	}
	return signedAuthTokenString, nil
}

func GetSignedRefreshTokenString(userId int) (string, error) {
	key, err := getSigningKey()
	if err != nil {
		return "", err
	}
	refreshTokenClaims := getClaims(fmt.Sprintf("%d", userId), time.Now().Add(RefreshTokenExpirationPeriod))
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	signedRefreshTokenString, refreshSignErr := refreshToken.SignedString(key)
	if refreshSignErr != nil {
		return "", refreshSignErr
	}
	return signedRefreshTokenString, nil
}

func GetSignedTokenStrings(userId int) (string, string, error) {
	key, err := getSigningKey()
	if err != nil {
		return "", "", err
	}

	authTokenClaims := getClaims(fmt.Sprintf("%d", userId), time.Now().Add(AuthTokenExpirationPeriod))
	authToken := jwt.NewWithClaims(jwt.SigningMethodHS256, authTokenClaims)
	signedAuthTokenString, authSignErr := authToken.SignedString(key)
	if authSignErr != nil {
		return "", "", authSignErr
	}

	refreshTokenClaims := getClaims(fmt.Sprintf("%d", userId), time.Now().Add(RefreshTokenExpirationPeriod))
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	signedRefreshTokenString, refreshSignErr := refreshToken.SignedString(key)
	if refreshSignErr != nil {
		return "", "", refreshSignErr
	}

	return signedAuthTokenString, signedRefreshTokenString, nil
}

func ParseAndValidateToken(tokenString string) (user int, err error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		secret, err := getSigningKey()
		if err != nil {
			return nil, err
		}
		return secret, nil
	})

	if err != nil {
		return 0, err
	}

	if !token.Valid {
		return 0, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("error viewing claims")
	}
	idString := claims["sub"].(string)
	user, convErr := strconv.Atoi(idString)
	if convErr != nil {
		return 0, convErr
	}
	return user, nil
}

func getClaims(userId string, expires time.Time) jwt.RegisteredClaims {
	return jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(expires),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Issuer:    "datamonster.api",
		Subject:   userId,
	}
}

func getSigningKey() ([]byte, error) {
	s := os.Getenv("SECRET")
	if s == "" {
		return nil, errors.New("no secret found")
	}
	key := []byte(s)
	return key, nil
}
