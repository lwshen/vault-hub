package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"log/slog"
	"os"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/lwshen/vault-hub/internal/config"
	"golang.org/x/oauth2"
)

var (
	provider     *oidc.Provider
	verifier     *oidc.IDTokenVerifier
	oauthConfig  *oauth2.Config
	sessionStore *sessions.CookieStore
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
	// Generate a random session key if not set
	sessionKey := config.SessionSecret
	if sessionKey == "" {
		key := make([]byte, 32)
		if _, err := rand.Read(key); err != nil {
			return err
		}
		sessionKey = base64.URLEncoding.EncodeToString(key)
		slog.Warn("OIDC session key not set, generated random key (not recommended for production)")
	}

	sessionStore = sessions.NewCookieStore([]byte(sessionKey))
	sessionStore.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   3600, // 1 hour
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: 2, // http.SameSiteLaxMode (importing net/http would cause cycle)
	}

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

func AuthCodeURL(c echo.Context, baseUrl string) (string, error) {
	state := generateState()
	err := storeInSession(c, "oauth", state)
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

func VerifyState(c echo.Context, state string) error {
	storedState, err := getFromSession(c, "oauth")
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

func storeInSession(c echo.Context, key string, value string) error {
	session, err := sessionStore.Get(c.Request(), "auth_session")
	if err != nil {
		return err
	}
	session.Values[key] = value
	return session.Save(c.Request(), c.Response())
}

func getFromSession(c echo.Context, key string) (string, error) {
	session, err := sessionStore.Get(c.Request(), "auth_session")
	if err != nil {
		return "", err
	}
	value, ok := session.Values[key].(string)
	if !ok {
		return "", errors.New("session value for key " + key + " is not a string")
	}
	return value, nil
}
