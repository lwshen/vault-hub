package auth

import (
	"context"
	"errors"
	"log/slog"
	"os"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/lwshen/vault-hub/internal/config"
	"golang.org/x/oauth2"
)

var (
	provider    *oidc.Provider
	verifier    *oidc.IDTokenVerifier
	oauthConfig *oauth2.Config
	// TODO: Implement Echo session storage
	sessionStates map[string]string
)

func init() {
	sessionStates = make(map[string]string)
	enabled := config.OidcEnabled
	slog.Info("OIDC", "enabled", enabled)
	if enabled {
		err := SetupOIDC()
		if err != nil {
			slog.Error("Failed to setup OIDC", "error", err)
			os.Exit(1)
		}
	}
}

func SetupOIDC() error {
	// TODO: Implement proper Echo session storage
	// For now using simple in-memory state storage
	ctx := context.Background()
	oidcProvider, err := oidc.NewProvider(ctx, config.OidcIssuer)
	if err != nil {
		return err
	}
	provider = oidcProvider
	oauthConfig = &oauth2.Config{
		ClientID:     config.OidcClientId,
		ClientSecret: config.OidcClientSecret,
		Scopes:       []string{oidc.ScopeOpenID, "email", "profile"},
		Endpoint:     oidcProvider.Endpoint(),
	}
	verifier = oidcProvider.Verifier(&oidc.Config{ClientID: oauthConfig.ClientID})
	return nil
}

func AuthCodeURL(ctx echo.Context, baseUrl string) (string, error) {
	state := generateState()
	sessionStates[state] = state // Simple in-memory storage
	oauthConfig.RedirectURL = baseUrl + "/api/auth/callback/oidc"
	return oauthConfig.AuthCodeURL(state), nil
}

func Verify(ctx context.Context, code string) (*oauth2.Token, error) {
	token, err := oauthConfig.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}

	_, err = verifyToken(ctx, token)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func verifyToken(ctx context.Context, token *oauth2.Token) (*oidc.IDToken, error) {
	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, errors.New("missing ID token")
	}
	return verifier.Verify(ctx, rawIDToken)
}

func VerifyState(ctx echo.Context, state string) error {
	storedState, exists := sessionStates[state]
	if !exists {
		return errors.New("state not found in session")
	}
	delete(sessionStates, state) // Clean up
	if storedState != state {
		return errors.New("cannot verify state")
	}
	return nil
}

func UserInfo(ctx context.Context, token *oauth2.Token) (*oidc.UserInfo, error) {
	tokenSource := oauthConfig.TokenSource(ctx, token)
	return provider.UserInfo(ctx, tokenSource)
}

func generateState() string {
	return uuid.New().String()
}

// TODO: Implement proper Echo session storage
// For now using simple in-memory state storage above
