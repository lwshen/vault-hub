package api

import (
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/lwshen/vault-hub/model"
	"github.com/lwshen/vault-hub/packages/api/generated/models"
)

// auditLogFilters holds parsed filter parameters
type auditLogFilters struct {
	pageSize  int
	pageIndex int
	startDate *time.Time
	endDate   *time.Time
	vaultID   *uint
}

// parseAuditLogFilters parses and validates query parameters for audit log filtering
func parseAuditLogFilters(ctx echo.Context, userID uint) (*auditLogFilters, error) {
	filters := &auditLogFilters{}

	// Parse pagination
	filters.pageSize = 20
	if ps := ctx.QueryParam("pageSize"); ps != "" {
		if v, err := strconv.Atoi(ps); err == nil {
			filters.pageSize = v
		}
	}

	filters.pageIndex = 1
	if pi := ctx.QueryParam("pageIndex"); pi != "" {
		if v, err := strconv.Atoi(pi); err == nil {
			filters.pageIndex = v
		}
	}

	// Validate bounds
	if filters.pageSize < 1 || filters.pageSize > 1000 {
		return nil, SendError(ctx, http.StatusBadRequest, "pageSize must be between 1 and 1000")
	}
	if filters.pageIndex < 1 {
		return nil, SendError(ctx, http.StatusBadRequest, "pageIndex must be at least 1")
	}

	// Parse date filters
	if sd := ctx.QueryParam("startDate"); sd != "" {
		parsed, err := time.Parse(time.RFC3339, sd)
		if err != nil {
			return nil, SendError(ctx, http.StatusBadRequest, "invalid startDate format, use ISO 8601")
		}
		filters.startDate = &parsed
	}

	if ed := ctx.QueryParam("endDate"); ed != "" {
		parsed, err := time.Parse(time.RFC3339, ed)
		if err != nil {
			return nil, SendError(ctx, http.StatusBadRequest, "invalid endDate format, use ISO 8601")
		}
		filters.endDate = &parsed
	}

	// Parse vault filter
	if vuid := ctx.QueryParam("vaultUniqueId"); vuid != "" {
		var vault model.Vault
		err := vault.GetByUniqueID(vuid, userID)
		if err != nil {
			return nil, SendError(ctx, http.StatusBadRequest, "invalid vaultUniqueId")
		}
		filters.vaultID = &vault.ID
	}

	return filters, nil
}

// GetAuditLogs handles GET /api/audit-logs with filtering and pagination
func (c *Container) GetAuditLogs(ctx echo.Context) error {
	user, err := getUserFromEchoContext(ctx)
	if err != nil {
		return err
	}

	// Parse and validate filters
	filters, err := parseAuditLogFilters(ctx, user.ID)
	if err != nil {
		return err
	}

	// Calculate offset from page index (1-based to 0-based)
	offset := (filters.pageIndex - 1) * filters.pageSize

	// Get audit logs from database
	params := model.GetAuditLogsWithFiltersParams{
		UserID:    user.ID,
		VaultID:   filters.vaultID,
		StartDate: filters.startDate,
		EndDate:   filters.endDate,
		Limit:     filters.pageSize,
		Offset:    offset,
	}

	logs, err := model.GetAuditLogsWithFilters(params)
	if err != nil {
		slog.Error("Failed to get audit logs", "error", err, "userID", user.ID)
		return SendError(ctx, http.StatusInternalServerError, "failed to retrieve audit logs")
	}

	// Get total count for pagination
	totalCount, err := model.CountAuditLogsWithFilters(params)
	if err != nil {
		slog.Error("Failed to count audit logs", "error", err, "userID", user.ID)
		return SendError(ctx, http.StatusInternalServerError, "failed to count audit logs")
	}

	// Convert to API audit logs
	apiLogs := make([]models.AuditLog, 0, len(logs))
	for i := range logs {
		apiLogs = append(apiLogs, convertToGeneratedAuditLog(&logs[i]))
	}

	response := models.AuditLogsResponse{
		AuditLogs:  apiLogs,
		TotalCount: safeInt64ToInt32(totalCount),
		PageSize:   int32(filters.pageSize),  // #nosec G115 -- validated max 1000
		PageIndex:  int32(filters.pageIndex), // #nosec G115 -- validated >= 1
	}

	return ctx.JSON(http.StatusOK, response)
}

// GetAuditMetrics handles GET /api/audit-logs/metrics
func (c *Container) GetAuditMetrics(ctx echo.Context) error {
	user, err := getUserFromEchoContext(ctx)
	if err != nil {
		return err
	}

	// Get metrics from database
	metrics, err := model.GetAllAuditMetrics(user.ID)
	if err != nil {
		slog.Error("Failed to get audit metrics", "error", err, "userID", user.ID)
		return SendError(ctx, http.StatusInternalServerError, "failed to retrieve audit metrics")
	}

	response := models.AuditMetricsResponse{
		TotalEventsLast30Days:  safeInt64ToInt32(metrics.TotalEventsLast30Days),
		EventsCountLast24Hours: safeInt64ToInt32(metrics.EventsCountLast24Hours),
		VaultEventsLast30Days:  safeInt64ToInt32(metrics.VaultEventsLast30Days),
		ApiKeyEventsLast30Days: safeInt64ToInt32(metrics.APIKeyEventsLast30Days),
	}

	return ctx.JSON(http.StatusOK, response)
}
