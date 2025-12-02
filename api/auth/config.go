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
	keycloakURL := getEnvOrDefault("KEYCLOAK_URL", "https://auth.electriclantern.net")
	realm := getEnvOrDefault("KEYCLOAK_REALM", "data-monster")
	clientID := getEnvOrDefault("KEYCLOAK_CLIENT_ID", "datamonster-app")
	clientSecret := os.Getenv("KEYCLOAK_CLIENT_SECRET")
	redirectURL := getEnvOrDefault("KEYCLOAK_REDIRECT_URL", "http://localhost:8080/auth/callback")

	if clientSecret == "" {
		return nil, fmt.Errorf("KEYCLOAK_CLIENT_SECRET environment variable is required")
	}

	issuerURL := fmt.Sprintf("%s/realms/%s", keycloakURL, realm)

	provider, err := oidc.NewProvider(ctx, issuerURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create OIDC provider: %w", err)
	}

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

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
