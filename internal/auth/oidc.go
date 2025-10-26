package auth

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"strings"

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

// AuthCodeURLEcho generates OIDC authorization URL and stores state in signed cookie (Echo version)
func AuthCodeURLEcho(ctx echo.Context, baseURL string) (string, error) {
	state := generateState()
	err := storeStateInCookie(ctx, state)
	if err != nil {
		return "", err
	}
	oauthConfig.RedirectURL = baseURL + "/api/auth/callback/oidc"
	return oauthConfig.AuthCodeURL(state), nil
}

// VerifyStateEcho verifies OAuth state from signed cookie (Echo version)
func VerifyStateEcho(ctx echo.Context, state string) error {
	storedState, err := getStateFromCookie(ctx)
	if err != nil {
		return err
	}
	if storedState != state {
		return errors.New("state mismatch")
	}
	// Delete cookie after verification (one-time use)
	deleteStateCookie(ctx)
	return nil
}

// Verify exchanges authorization code for OAuth token
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

// verifyToken validates the ID token from OAuth response
func verifyToken(ctx context.Context, token *oauth2.Token) (*oidc.IDToken, error) {
	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, errors.New("missing ID token")
	}
	return verifier.Verify(ctx, rawIDToken)
}

// UserInfo fetches user information from OIDC provider
func UserInfo(ctx context.Context, token *oauth2.Token) (*oidc.UserInfo, error) {
	tokenSource := oauthConfig.TokenSource(ctx, token)
	return provider.UserInfo(ctx, tokenSource)
}

// generateState creates a random state string for OAuth CSRF protection
func generateState() string {
	return uuid.New().String()
}

// storeStateInCookie stores OAuth state in a signed cookie
func storeStateInCookie(ctx echo.Context, state string) error {
	signature := signState(state)
	cookieValue := state + "." + signature

	cookie := new(http.Cookie)
	cookie.Name = "oauth_state"
	cookie.Value = cookieValue
	cookie.HttpOnly = true
	cookie.Secure = ctx.Request().TLS != nil || ctx.Request().Header.Get("X-Forwarded-Proto") == "https"
	cookie.SameSite = http.SameSiteLaxMode
	cookie.MaxAge = 600 // 10 minutes
	cookie.Path = "/"

	ctx.SetCookie(cookie)
	return nil
}

// getStateFromCookie retrieves and verifies OAuth state from signed cookie
func getStateFromCookie(ctx echo.Context) (string, error) {
	cookie, err := ctx.Cookie("oauth_state")
	if err != nil {
		return "", errors.New("oauth state cookie not found")
	}

	// Parse state and signature
	parts := strings.Split(cookie.Value, ".")
	if len(parts) != 2 {
		return "", errors.New("invalid oauth state cookie format")
	}

	state := parts[0]
	signature := parts[1]

	// Verify signature
	if !verifyStateSignature(state, signature) {
		return "", errors.New("oauth state signature verification failed")
	}

	return state, nil
}

// deleteStateCookie removes the OAuth state cookie after use
func deleteStateCookie(ctx echo.Context) {
	cookie := new(http.Cookie)
	cookie.Name = "oauth_state"
	cookie.Value = ""
	cookie.HttpOnly = true
	cookie.MaxAge = -1 // Delete cookie
	cookie.Path = "/"
	ctx.SetCookie(cookie)
}

// signState creates HMAC-SHA256 signature for OAuth state
func signState(state string) string {
	h := hmac.New(sha256.New, []byte(config.JwtSecret))
	h.Write([]byte(state))
	return base64.URLEncoding.EncodeToString(h.Sum(nil))
}

// verifyStateSignature verifies HMAC-SHA256 signature of OAuth state
func verifyStateSignature(state, signature string) bool {
	expectedSignature := signState(state)
	return hmac.Equal([]byte(signature), []byte(expectedSignature))
}
