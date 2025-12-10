package auth

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

type Config struct {
	ClientID      string
	ClientSecret  string
	IssuerURL     string
	RedirectURL   string
	IntrospectURL string
}

func (c Config) Validate() error {
	if c.ClientID == "" {
		return errors.New("clientID is required")
	}

	if c.ClientSecret == "" {
		return errors.New("clientSecret is required")
	}

	if c.IssuerURL == "" {
		return errors.New("issuerURL is required")
	}

	if c.RedirectURL == "" {
		return errors.New("redirectURL is required")
	}

	if c.IntrospectURL == "" {
		return errors.New("introspectURL is required")
	}

	return nil
}

type Provider struct {
	ClientID      string
	ClientSecret  string
	IssuerURL     string
	IntrospectURL string
	OAuth2Config  *oauth2.Config
	Verifier      *oidc.IDTokenVerifier
	httpClient    *http.Client
}

func NewAuthConfig(c Config) (*Provider, error) {
	if validateErr := c.Validate(); validateErr != nil {
		return nil, validateErr
	}

	ctx := context.Background()
	provider, err := oidc.NewProvider(ctx, c.IssuerURL)
	if err != nil {
		return nil, err
	}

	oauth2Config := &oauth2.Config{
		ClientID:     c.ClientID,
		ClientSecret: c.ClientSecret,
		RedirectURL:  c.RedirectURL,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email", "groups"},
	}

	verifier := provider.Verifier(&oidc.Config{ClientID: c.ClientID})

	httpClient := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 10,
			IdleConnTimeout:     90 * time.Second,
		},
	}

	return &Provider{
		ClientID:      c.ClientID,
		ClientSecret:  c.ClientSecret,
		IssuerURL:     c.IssuerURL,
		IntrospectURL: c.IntrospectURL,
		OAuth2Config:  oauth2Config,
		Verifier:      verifier,
		httpClient:    httpClient,
	}, nil
}
