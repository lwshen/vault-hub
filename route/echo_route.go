package route

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lwshen/vault-hub/internal/embed"
	"github.com/lwshen/vault-hub/packages/api"
)

// SetupRoutesEcho configures all routes for the Echo server
func SetupRoutesEcho(e *echo.Echo) {
	// Apply authentication middleware
	e.Use(JWTMiddlewareEcho())

	// Create handler container
	container, err := api.NewContainer()
	if err != nil {
		slog.Error("Failed to create handler container", "error", err)
		os.Exit(1)
	}

	// Auth endpoints
	e.POST("/api/auth/login", container.Login)
	e.POST("/api/auth/signup", container.Signup)
	e.GET("/api/auth/logout", container.Logout)
	e.POST("/api/auth/password/reset/request", container.RequestPasswordReset)
	e.POST("/api/auth/password/reset/confirm", container.ConfirmPasswordReset)
	e.POST("/api/auth/magic-link/request", container.RequestMagicLink)
	e.GET("/api/auth/magic-link/token", container.ConsumeMagicLink)

	// User endpoints
	e.GET("/api/user", container.GetCurrentUser)

	// Vault endpoints (web API)
	e.GET("/api/vaults", container.GetVaults)
	e.GET("/api/vaults/:uniqueId", container.GetVault)
	e.POST("/api/vaults", container.CreateVault)
	e.PUT("/api/vaults/:uniqueId", container.UpdateVault)
	e.DELETE("/api/vaults/:uniqueId", container.DeleteVault)

	// API Key endpoints
	e.GET("/api/api-keys", container.GetAPIKeys)
	e.POST("/api/api-keys", container.CreateAPIKey)
	e.PATCH("/api/api-keys/:id", container.UpdateAPIKey)
	e.DELETE("/api/api-keys/:id", container.DeleteAPIKey)

	// CLI vault access endpoints (API key auth)
	e.GET("/api/cli/vaults", container.GetVaultsByAPIKey)
	e.GET("/api/cli/vault/:uniqueId", container.GetVaultByAPIKey)
	e.GET("/api/cli/vault/name/:name", container.GetVaultByNameAPIKey)

	// Audit endpoints
	e.GET("/api/audit-logs", container.GetAuditLogs)
	e.GET("/api/audit-logs/metrics", container.GetAuditMetrics)

	// System endpoints (public)
	e.GET("/api/health", container.Health)
	e.GET("/api/status", container.GetStatus)
	e.GET("/api/config", container.GetConfig)

	// Static web assets
	embedFS, err := embed.GetDistFS()
	if err != nil {
		slog.Error("Failed to initialize embedded filesystem", "error", err)
		os.Exit(1)
	}

	// Serve static files from embedded filesystem
	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root:       "/",
		Filesystem: http.FS(embedFS),
		HTML5:      true,
		Browse:     false,
	}))
}
