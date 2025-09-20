package api

import (
	"net/http"
	"os"
	"runtime"
	"syscall"
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
	// Check critical system resources
	memStats := &runtime.MemStats{}
	runtime.ReadMemStats(memStats)
	
	// Check disk space
	diskUsage := checkDiskSpace()
	
	// Memory usage check - if using >90% of allocated memory, system is degraded
	memUsagePercent := float64(memStats.Alloc) / float64(memStats.Sys) * 100
	
	// Critical system failures that make system unavailable regardless of database
	if memUsagePercent > 95 || diskUsage > 98 {
		return StatusResponseSystemStatusUnavailable
	}
	
	// System is unavailable if database is completely down
	if dbStatus == StatusResponseDatabaseStatusUnavailable {
		return StatusResponseSystemStatusUnavailable
	}

	// System degradation scenarios:
	// 1. Database is degraded
	// 2. High memory usage (>80%)
	// 3. High disk usage (>90%)
	// 4. High database response time (>2 seconds)
	if dbStatus == StatusResponseDatabaseStatusDegraded ||
		memUsagePercent > 80 ||
		diskUsage > 90 ||
		dbResponseTime > 2000 {
		return StatusResponseSystemStatusDegraded
	}

	return StatusResponseSystemStatusHealthy
}

// checkDiskSpace returns disk usage percentage
func checkDiskSpace() float64 {
	var stat syscall.Statfs_t
	wd, err := os.Getwd()
	if err != nil {
		return 0 // If we can't check, assume it's fine
	}
	
	err = syscall.Statfs(wd, &stat)
	if err != nil {
		return 0 // If we can't check, assume it's fine
	}

	// Calculate disk usage percentage
	total := stat.Blocks * uint64(stat.Bsize)
	available := stat.Bavail * uint64(stat.Bsize)
	used := total - available
	
	if total == 0 {
		return 0
	}
	
	return float64(used) / float64(total) * 100
}
