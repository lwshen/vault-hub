package main

import (
	"github.com/GIT_USER_ID/GIT_REPO_ID/handlers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	//todo: handle the error!
	c, _ := handlers.NewContainer()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())


	// CreateAPIKey - 
	e.POST("/api/api-keys", c.CreateAPIKey)

	// DeleteAPIKey - 
	e.DELETE("/api/api-keys/:id", c.DeleteAPIKey)

	// GetAPIKeys - 
	e.GET("/api/api-keys", c.GetAPIKeys)

	// UpdateAPIKey - 
	e.PATCH("/api/api-keys/:id", c.UpdateAPIKey)

	// GetAuditLogs - 
	e.GET("/api/audit-logs", c.GetAuditLogs)

	// GetAuditMetrics - 
	e.GET("/api/audit-logs/metrics", c.GetAuditMetrics)

	// ConfirmPasswordReset - 
	e.POST("/api/auth/password/reset/confirm", c.ConfirmPasswordReset)

	// ConsumeMagicLink - 
	e.GET("/api/auth/magic-link/token", c.ConsumeMagicLink)

	// Login - 
	e.POST("/api/auth/login", c.Login)

	// Logout - 
	e.GET("/api/auth/logout", c.Logout)

	// RequestMagicLink - 
	e.POST("/api/auth/magic-link/request", c.RequestMagicLink)

	// RequestPasswordReset - 
	e.POST("/api/auth/password/reset/request", c.RequestPasswordReset)

	// Signup - 
	e.POST("/api/auth/signup", c.Signup)

	// GetVaultByAPIKey - 
	e.GET("/api/cli/vault/:uniqueId", c.GetVaultByAPIKey)

	// GetVaultByNameAPIKey - 
	e.GET("/api/cli/vault/name/:name", c.GetVaultByNameAPIKey)

	// GetVaultsByAPIKey - 
	e.GET("/api/cli/vaults", c.GetVaultsByAPIKey)

	// GetConfig - Get public configuration
	e.GET("/api/config", c.GetConfig)

	// Health - 
	e.GET("/api/health", c.Health)

	// GetStatus - Get system status
	e.GET("/api/status", c.GetStatus)

	// GetCurrentUser - 
	e.GET("/api/user", c.GetCurrentUser)

	// CreateVault - 
	e.POST("/api/vaults", c.CreateVault)

	// DeleteVault - 
	e.DELETE("/api/vaults/:uniqueId", c.DeleteVault)

	// GetVault - 
	e.GET("/api/vaults/:uniqueId", c.GetVault)

	// GetVaults - 
	e.GET("/api/vaults", c.GetVaults)

	// UpdateVault - 
	e.PUT("/api/vaults/:uniqueId", c.UpdateVault)


	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
