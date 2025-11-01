package auth

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/google/uuid"
	"github.com/lwshen/vault-hub/internal/config"
	"golang.org/x/oauth2"
)

var (
	provider    *oidc.Provider
	verifier    *oidc.IDTokenVerifier
	oauthConfig *oauth2.Config
	stateStore  = newStateStore(5 * time.Minute)
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

func AuthCodeURL(baseUrl string) (string, error) {
	if oauthConfig == nil {
		return "", errors.New("OIDC not configured")
	}

	state := generateState()
	stateStore.Save(state)
	oauthConfig.RedirectURL = baseUrl + "/api/auth/callback/oidc"
	return oauthConfig.AuthCodeURL(state), nil
}

func Verify(ctx context.Context, code string) (*oauth2.Token, error) {
	if oauthConfig == nil {
		return nil, errors.New("OIDC not configured")
	}

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

func VerifyState(state string) error {
	if !stateStore.Verify(state) {
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

type stateStoreCache struct {
	data sync.Map
	ttl  time.Duration
}

func newStateStore(ttl time.Duration) *stateStoreCache {
	return &stateStoreCache{ttl: ttl}
}

func (s *stateStoreCache) Save(state string) {
	expires := time.Now().Add(s.ttl)
	s.data.Store(state, expires)
	time.AfterFunc(s.ttl, func() {
		s.data.Delete(state)
	})
}

func (s *stateStoreCache) Verify(state string) bool {
	value, ok := s.data.Load(state)
	if !ok {
		return false
	}
	expires, _ := value.(time.Time)
	if time.Now().After(expires) {
		s.data.Delete(state)
		return false
	}
	s.data.Delete(state)
	return true
}
