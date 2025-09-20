package api

import (
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/lwshen/vault-hub/internal/version"
	"github.com/lwshen/vault-hub/model"
)

func (s Server) GetStatus(ctx *fiber.Ctx) error {
	// Check database status with multiple health indicators
	databaseStatus, dbConnections, dbResponseTime := checkDatabaseHealth()

	// Check system status based on multiple factors
	systemStatus := checkSystemHealth(databaseStatus, dbConnections, dbResponseTime)

	resp := StatusResponse{
		Version:        version.Version,
		Commit:         version.Commit,
		SystemStatus:   systemStatus,
		DatabaseStatus: databaseStatus,
	}

	return ctx.
		Status(http.StatusOK).
		JSON(resp)
}

// checkDatabaseHealth performs comprehensive database health checks
func checkDatabaseHealth() (StatusResponseDatabaseStatus, int, int64) {
	// Test basic connectivity
	start := time.Now()
	if err := model.DB.Exec("SELECT 1").Error; err != nil {
		return StatusResponseDatabaseStatusUnavailable, 0, 0
	}
	responseTime := time.Since(start).Milliseconds()

	// Get database connection info
	sqlDB, err := model.DB.DB()
	if err != nil {
		return StatusResponseDatabaseStatusDegraded, 0, responseTime
	}

	stats := sqlDB.Stats()

	// Determine status based on performance metrics
	if responseTime > 5000 { // >5 seconds is severely degraded
		return StatusResponseDatabaseStatusDegraded, stats.OpenConnections, responseTime
	}

	if responseTime > 1000 || stats.OpenConnections > 80 { // >1 second or >80% of max connections
		return StatusResponseDatabaseStatusDegraded, stats.OpenConnections, responseTime
	}

	return StatusResponseDatabaseStatusHealthy, stats.OpenConnections, responseTime
}

// checkSystemHealth determines overall system status based on various factors
func checkSystemHealth(dbStatus StatusResponseDatabaseStatus, dbConnections int, dbResponseTime int64) StatusResponseSystemStatus {
	// System is unavailable if database is completely down
	if dbStatus == StatusResponseDatabaseStatusUnavailable {
		return StatusResponseSystemStatusUnavailable
	}

	// System degradation scenarios:
	// 1. Database is degraded
	// 2. High database response time (>2 seconds)
	if dbStatus == StatusResponseDatabaseStatusDegraded ||
		dbResponseTime > 2000 {
		return StatusResponseSystemStatusDegraded
	}

	return StatusResponseSystemStatusHealthy
}
