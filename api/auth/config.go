package auth

import (
	"context"
	"fmt"
	"os"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

type KeycloakConfig struct {
	Provider     *oidc.Provider
	OAuth2Config oauth2.Config
	Verifier     *oidc.IDTokenVerifier
	IssuerURL    string
	ClientID     string
}

func NewKeycloakConfig(ctx context.Context) (*KeycloakConfig, error) {
	keycloakURL := os.Getenv("KEYCLOAK_URL")
	if keycloakURL == "" {
		return nil, fmt.Errorf("KEYCLOAK_URL environment variable is required")
	}
	realm := os.Getenv("KEYCLOAK_REALM")
	if realm == "" {
		return nil, fmt.Errorf("KEYCLOAK_REALM environment variable is required")
	}
	clientID := os.Getenv("KEYCLOAK_CLIENT_ID")
	if clientID == "" {
		return nil, fmt.Errorf("KEYCLOAK_CLIENT_ID environment variable is required")
	}
	clientSecret := os.Getenv("KEYCLOAK_CLIENT_SECRET")
	if clientSecret == "" {
		return nil, fmt.Errorf("KEYCLOAK_CLIENT_SECRET environment variable is required")
	}
	redirectURL := os.Getenv("KEYCLOAK_REDIRECT_URL")
	if redirectURL == "" {
		return nil, fmt.Errorf("KEYCLOAK_REDIRECT_URL environment variable is required")
	}

	issuerURL := fmt.Sprintf("%s/realms/%s", keycloakURL, realm)

	providerConfig := &oidc.ProviderConfig{
		IssuerURL:   issuerURL,
		AuthURL:     issuerURL + "/protocol/openid-connect/auth",
		TokenURL:    issuerURL + "/protocol/openid-connect/token",
		JWKSURL:     issuerURL + "/protocol/openid-connect/certs",
		UserInfoURL: issuerURL + "/protocol/openid-connect/userinfo",
	}
	provider := providerConfig.NewProvider(ctx)

	oauth2Config := oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}

	verifier := provider.Verifier(&oidc.Config{
		ClientID:                   clientID,
		SkipClientIDCheck:          false,
		SkipExpiryCheck:            false,
		SkipIssuerCheck:            false,
		InsecureSkipSignatureCheck: false,
	})

	return &KeycloakConfig{
		Provider:     provider,
		OAuth2Config: oauth2Config,
		Verifier:     verifier,
		IssuerURL:    issuerURL,
		ClientID:     clientID,
	}, nil
}
