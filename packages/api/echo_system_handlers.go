package api

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/lwshen/vault-hub/internal/config"
	"github.com/lwshen/vault-hub/internal/version"
	"github.com/lwshen/vault-hub/model"
	"github.com/lwshen/vault-hub/packages/api/generated_models"
)

// Health handles GET /api/health
func (c *Container) Health(ctx echo.Context) error {
	status := "ok"
	timestamp := time.Now()
	resp := generated_models.HealthCheckResponse{
		Status:    status,
		Timestamp: timestamp,
	}

	return ctx.JSON(http.StatusOK, resp)
}

// GetStatus handles GET /api/status
func (c *Container) GetStatus(ctx echo.Context) error {
	// Check database status with multiple health indicators
	databaseStatus, _, _ := checkDatabaseHealthEcho()

	// Check system status based on multiple factors
	systemStatus := checkSystemHealthEcho(databaseStatus)

	resp := generated_models.StatusResponse{
		Version:        version.Version,
		Commit:         version.Commit,
		SystemStatus:   systemStatus,
		DatabaseStatus: databaseStatus,
	}

	return ctx.JSON(http.StatusOK, resp)
}

// checkDatabaseHealthEcho performs comprehensive database health checks
func checkDatabaseHealthEcho() (string, int, int64) {
	// Test basic connectivity
	start := time.Now()
	if err := model.DB.Exec("SELECT 1").Error; err != nil {
		return "unavailable", 0, 0
	}
	responseTime := time.Since(start).Milliseconds()

	// Get database connection info
	sqlDB, err := model.DB.DB()
	if err != nil {
		return "degraded", 0, responseTime
	}

	stats := sqlDB.Stats()

	// Determine status based on performance metrics
	if responseTime > 5000 { // >5 seconds is severely degraded
		return "degraded", stats.OpenConnections, responseTime
	}

	if responseTime > 1000 || stats.OpenConnections > 80 { // >1 second or >80% of max connections
		return "degraded", stats.OpenConnections, responseTime
	}

	return "healthy", stats.OpenConnections, responseTime
}

// checkSystemHealthEcho determines overall system status based on various factors
func checkSystemHealthEcho(dbStatus string) string {
	// System is unavailable if database is completely down
	if dbStatus == "unavailable" {
		return "unavailable"
	}

	// System degradation scenarios:
	// 1. Database is degraded
	if dbStatus == "degraded" {
		return "degraded"
	}

	return "healthy"
}

// GetConfig handles GET /api/config
// Returns public configuration that requires no authentication
func (c *Container) GetConfig(ctx echo.Context) error {
	resp := generated_models.ConfigResponse{
		OidcEnabled:  config.OidcEnabled,
		EmailEnabled: config.EmailEnabled,
	}

	return ctx.JSON(http.StatusOK, resp)
}

// GetCurrentUser handles GET /api/user
func (c *Container) GetCurrentUser(ctx echo.Context) error {
	user, err := getUserFromEchoContext(ctx)
	if err != nil {
		return err
	}

	resp := generated_models.GetUserResponse{
		Email: user.Email,
	}

	if user.Name != nil {
		resp.Name = *user.Name
	}

	if user.Avatar != nil {
		resp.Avatar = *user.Avatar
	}

	return ctx.JSON(http.StatusOK, resp)
}
