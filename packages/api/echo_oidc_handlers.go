package api

import (
	"log/slog"
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"
	"github.com/lwshen/vault-hub/internal/auth"
	"github.com/lwshen/vault-hub/model"
)

// LoginOIDC handles GET /api/auth/login/oidc - Initiates OIDC authentication flow
func (c *Container) LoginOIDC(ctx echo.Context) error {
	baseURL := getBaseURL(ctx)
	authURL, err := auth.AuthCodeURLEcho(ctx, baseURL)
	if err != nil {
		slog.Error("Failed to get OIDC URL", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to initiate OIDC login")
	}
	slog.Debug("Login with OIDC", "url", authURL)
	return ctx.Redirect(http.StatusFound, authURL)
}

// LoginOIDCCallback handles GET /api/auth/callback/oidc - Processes OIDC callback
func (c *Container) LoginOIDCCallback(ctx echo.Context) error {
	code := ctx.QueryParam("code")
	state := ctx.QueryParam("state")
	slog.Debug("Login with OIDC callback", "uri", ctx.Request().RequestURI, "code", code, "state", state)

	// Verify state parameter matches cookie
	err := auth.VerifyStateEcho(ctx, state)
	if err != nil {
		slog.Error("OIDC state verification failed", "error", err)
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid OIDC state")
	}

	// Exchange authorization code for token
	token, err := auth.Verify(ctx.Request().Context(), code)
	if err != nil {
		slog.Error("OIDC token exchange failed", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to exchange OIDC token")
	}
	slog.Debug("Login with OIDC callback", "token", token)

	// Get user info from OIDC provider
	userInfo, err := auth.UserInfo(ctx.Request().Context(), token)
	if err != nil {
		slog.Error("Failed to get OIDC user info", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get user info")
	}
	slog.Debug("Login with OIDC callback", "userInfo", userInfo)

	// Extract claims from userInfo
	var claims map[string]interface{}
	if err := userInfo.Claims(&claims); err != nil {
		slog.Error("Failed to extract OIDC claims", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to extract claims")
	}

	// Extract email from claims
	email, ok := claims["email"].(string)
	if !ok || email == "" {
		slog.Error("OIDC userInfo missing email claim", "claims", claims)
		return echo.NewHTTPError(http.StatusBadRequest, "Email claim missing from OIDC response")
	}

	// Look up user by email
	user := model.User{
		Email: email,
	}
	if err := user.GetByEmail(); err != nil {
		// User doesn't exist, create new user from OIDC data
		name := ""
		if nameClaim, ok := claims["name"].(string); ok {
			name = nameClaim
		}

		createParams := model.CreateUserParams{
			Email:    email,
			Password: nil, // OIDC users don't need passwords
			Name:     name,
		}

		newUser, createErr := createParams.Create()
		if createErr != nil {
			slog.Error("Failed to create user from OIDC", "error", createErr, "email", email)
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create user")
		}
		user = *newUser
		slog.Info("User created from OIDC", "email", email, "name", name)
	}

	// Generate JWT token for the user
	jwtToken, err := user.GenerateToken()
	if err != nil {
		slog.Error("Failed to generate token for OIDC user", "error", err, "userID", user.ID)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate token")
	}

	// Record successful login audit log
	clientIP, userAgent := getClientInfoEcho(ctx)
	if err := model.LogUserAction(model.ActionLoginUser, user.ID, model.SourceWeb, clientIP, userAgent); err != nil {
		slog.Error("Failed to create audit log for OIDC login", "error", err, "userID", user.ID)
	}

	// Redirect back to frontend with token in URL fragment (hash) for security
	// URL fragments are never sent to the server, preventing token leakage in logs, Referer headers, and browser history
	redirectURL := "/login#token=" + url.QueryEscape(jwtToken) + "&source=oidc"
	return ctx.Redirect(http.StatusFound, redirectURL)
}

// getBaseURL extracts the base URL from the Echo context
func getBaseURL(ctx echo.Context) string {
	req := ctx.Request()
	scheme := "http"
	if req.TLS != nil || ctx.Request().Header.Get("X-Forwarded-Proto") == "https" {
		scheme = "https"
	}
	host := req.Host
	if host == "" {
		host = req.Header.Get("X-Forwarded-Host")
	}
	if host == "" {
		host = "localhost:3000"
	}
	return scheme + "://" + host
}
