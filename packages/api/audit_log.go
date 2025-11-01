package api

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/lwshen/vault-hub/handler"
	"github.com/lwshen/vault-hub/model"
	"gorm.io/gorm"
)

// convertToApiAuditLog converts a model.AuditLog to an api.AuditLog
func convertToApiAuditLog(auditLog *model.AuditLog) AuditLog {
	var vault *VaultLite
	var apiKey *VaultAPIKey
	if auditLog.Vault != nil {
		vaultLite := convertToApiVaultLite(auditLog.Vault)
		vault = &vaultLite
	}
	if auditLog.APIKey != nil {
		apiKeyLocal, _ := convertToApiAPIKey(auditLog.APIKey)
		apiKey = apiKeyLocal
	}

	return AuditLog{
		Action:    AuditLogAction(auditLog.Action),
		CreatedAt: auditLog.CreatedAt,
		Vault:     vault,
		ApiKey:    apiKey,
		Source:    AuditLogSource(auditLog.Source),
		IpAddress: &auditLog.IPAddress,
		UserAgent: &auditLog.UserAgent,
	}
}

// GetAuditLogsForUser retrieves filtered and paginated audit logs for the
// authenticated user. It supports filtering by vault, date range, and
// pagination parameters while remaining framework-agnostic for reuse in Echo.
func GetAuditLogsForUser(user *model.User, params GetAuditLogsParams) (AuditLogsResponse, *APIError) {
	if user == nil {
		return AuditLogsResponse{}, newAPIError(http.StatusUnauthorized, "user not found in context")
	}

	if err := validateAuditLogParams(params); err != nil {
		return AuditLogsResponse{}, newAPIError(http.StatusBadRequest, err.Error())
	}

	filterParams, err := buildAuditLogFilterParams(params, user.ID)
	if err != nil {
		return AuditLogsResponse{}, newAPIError(http.StatusBadRequest, err.Error())
	}

	logs, totalCount, err := fetchAuditLogsWithCount(filterParams)
	if err != nil {
		return AuditLogsResponse{}, newAPIError(http.StatusInternalServerError, err.Error())
	}

	response := buildAuditLogsResponse(logs, totalCount, params)
	return response, nil
}

// GetAuditLogs retains the Fiber handler by delegating to GetAuditLogsForUser.
func (Server) GetAuditLogs(c *fiber.Ctx, params GetAuditLogsParams) error {
	user, err := getUserFromContext(c)
	if err != nil {
		return err
	}

	resp, apiErr := GetAuditLogsForUser(user, params)
	if apiErr != nil {
		return handler.SendError(c, apiErr.Status, apiErr.Message)
	}

	return c.Status(http.StatusOK).JSON(resp)
}

// validateAuditLogParams validates pagination and other request parameters
func validateAuditLogParams(params GetAuditLogsParams) error {
	if params.PageSize <= 0 || params.PageSize > 1000 {
		return fmt.Errorf("pageSize must be between 1 and 1000")
	}
	if params.PageIndex < 1 {
		return fmt.Errorf("pageIndex must be greater than 0")
	}
	return nil
}

// buildAuditLogFilterParams constructs filter parameters for audit log queries
func buildAuditLogFilterParams(params GetAuditLogsParams, userID uint) (model.GetAuditLogsWithFiltersParams, error) {
	// Convert pageIndex from 1-based to 0-based for offset calculation
	offset := (params.PageIndex - 1) * params.PageSize

	// Resolve vault ID if vault unique ID is provided
	vaultID, err := resolveVaultID(params.VaultUniqueId, userID)
	if err != nil {
		return model.GetAuditLogsWithFiltersParams{}, err
	}

	return model.GetAuditLogsWithFiltersParams{
		UserID:    userID,
		VaultID:   vaultID,
		StartDate: params.StartDate,
		EndDate:   params.EndDate,
		Limit:     params.PageSize,
		Offset:    offset,
	}, nil
}

// resolveVaultID converts a vault unique ID to a database ID
func resolveVaultID(vaultUniqueId *string, userID uint) (*uint, error) {
	if vaultUniqueId == nil {
		return nil, nil
	}

	var vault model.Vault
	err := vault.GetByUniqueID(*vaultUniqueId, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("vault not found or access denied")
		}
		return nil, fmt.Errorf("failed to resolve vault: %v", err)
	}
	return &vault.ID, nil
}

// fetchAuditLogsWithCount retrieves audit logs and total count in parallel
func fetchAuditLogsWithCount(filterParams model.GetAuditLogsWithFiltersParams) ([]model.AuditLog, int64, error) {
	// Get audit logs with filters
	logs, err := model.GetAuditLogsWithFilters(filterParams)
	if err != nil {
		return nil, 0, err
	}

	// Get total count for pagination
	totalCount, err := model.CountAuditLogsWithFilters(filterParams)
	if err != nil {
		return nil, 0, err
	}

	return logs, totalCount, nil
}

// buildAuditLogsResponse converts audit logs to API format and constructs response
func buildAuditLogsResponse(logs []model.AuditLog, totalCount int64, params GetAuditLogsParams) AuditLogsResponse {
	// Convert to API format
	apiLogs := make([]AuditLog, len(logs))
	for i, log := range logs {
		apiLogs[i] = convertToApiAuditLog(&log)
	}

	return AuditLogsResponse{
		AuditLogs:  apiLogs,
		TotalCount: int(totalCount),
		PageSize:   params.PageSize,
		PageIndex:  params.PageIndex,
	}
}

// GetAuditMetricsForUser provides a framework-neutral implementation for
// retrieving audit metrics.
func GetAuditMetricsForUser(user *model.User) (AuditMetricsResponse, *APIError) {
	if user == nil {
		return AuditMetricsResponse{}, newAPIError(http.StatusUnauthorized, "user not found in context")
	}

	metrics, err := model.GetAllAuditMetrics(user.ID)
	if err != nil {
		return AuditMetricsResponse{}, newAPIError(http.StatusInternalServerError, err.Error())
	}

	response := AuditMetricsResponse{
		TotalEventsLast30Days:  int(metrics.TotalEventsLast30Days),
		EventsCountLast24Hours: int(metrics.EventsCountLast24Hours),
		VaultEventsLast30Days:  int(metrics.VaultEventsLast30Days),
		ApiKeyEventsLast30Days: int(metrics.APIKeyEventsLast30Days),
	}

	return response, nil
}

// GetAuditMetrics keeps the Fiber handler delegating to the helper for reuse.
func (Server) GetAuditMetrics(c *fiber.Ctx) error {
	user, err := getUserFromContext(c)
	if err != nil {
		return err
	}

	resp, apiErr := GetAuditMetricsForUser(user)
	if apiErr != nil {
		return handler.SendError(c, apiErr.Status, apiErr.Message)
	}

	return c.Status(http.StatusOK).JSON(resp)
}
