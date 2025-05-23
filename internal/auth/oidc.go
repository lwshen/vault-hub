package auth

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/google/uuid"
	"github.com/lwshen/vault-hub/internal/config"
	"golang.org/x/oauth2"
)

var (
	provider     *oidc.Provider
	verifier     *oidc.IDTokenVerifier
	oauthConfig  *oauth2.Config
	sessionStore *session.Store
)

func init() {
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
	sessionStore = session.New(session.Config{
		KeyLookup:  "cookie:auth_session",
		Expiration: time.Hour * 1,
	})

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

func AuthCodeURL(ctx *fiber.Ctx, baseUrl string) (string, error) {
	state := generateState()
	err := storeInSession(ctx, "oauth", state)
	if err != nil {
		return "", err
	}
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

func VerifyState(ctx *fiber.Ctx, state string) error {
	storedState, err := getFromSession(ctx, "oauth")
	if err != nil {
		return err
	}
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

func storeInSession(ctx *fiber.Ctx, key string, value string) error {
	session, err := sessionStore.Get(ctx)
	if err != nil {
		return err
	}
	session.Set(key, value)
	return session.Save()
}

func getFromSession(ctx *fiber.Ctx, key string) (string, error) {
	session, err := sessionStore.Get(ctx)
	if err != nil {
		return "", err
	}
	value, ok := session.Get(key).(string)
	if !ok {
		return "", errors.New("session value for key " + key + " is not a string")
	}
	return value, nil
}
