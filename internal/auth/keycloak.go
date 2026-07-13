package auth

import (
	"context"
	"log"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"

	"github.com/mic615/chill-crate-api/internal/config"
)

type Client struct {
	Provider *oidc.Provider
	OIDC     *oidc.IDTokenVerifier
	OAuth    oauth2.Config
}

var KCClient *Client

func NewClient(cfg *config.Config) {
	providerURL := cfg.KeycloakURL + "/realms/" + cfg.Realm
	provider, err := oidc.NewProvider(context.Background(), providerURL)
	if err != nil {
		log.Fatalf("failed to create OIDC provider: %v", err)
	}

	verifier := provider.Verifier(&oidc.Config{ClientID: cfg.KCClientID})
	oauth := oauth2.Config{
		ClientID:     cfg.KCClientID,
		ClientSecret: cfg.KCClientSecret,
		Endpoint:     provider.Endpoint(),
		RedirectURL:  cfg.RedirectURL,
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}

	KCClient = &Client{
		Provider: provider,
		OIDC:     verifier,
		OAuth:    oauth,
	}
}
