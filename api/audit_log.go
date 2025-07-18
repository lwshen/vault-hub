package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lwshen/vault-hub/handler"
	"github.com/lwshen/vault-hub/model"
)

// convertToApiAuditLog converts a model.AuditLog to an api.AuditLog
func convertToApiAuditLog(auditLog *model.AuditLog) AuditLog {
	var vaultID *int64
	if auditLog.VaultID != nil {
		// #nosec G115
		vaultIDValue := int64(*auditLog.VaultID)
		vaultID = &vaultIDValue
	}

	return AuditLog{
		Action:    AuditLogAction(auditLog.Action),
		CreatedAt: auditLog.CreatedAt,
		VaultId:   vaultID,
		IpAddress: &auditLog.IPAddress,
		UserAgent: &auditLog.UserAgent,
	}
}

func (Server) GetAuditLogs(c *fiber.Ctx, params GetAuditLogsParams) error {
	// Get authenticated user
	user, err := getUserFromContext(c)
	if err != nil {
		return err
	}

	// Validate pagination parameters
	if params.PageSize <= 0 || params.PageSize > 1000 {
		return handler.SendError(c, fiber.StatusBadRequest, "pageSize must be between 1 and 1000")
	}
	if params.PageIndex < 1 {
		return handler.SendError(c, fiber.StatusBadRequest, "pageIndex must be greater than 0")
	}

	// Convert pageIndex from 1-based to 0-based for offset calculation
	offset := (params.PageIndex - 1) * params.PageSize

	// Prepare filter parameters
	var vaultID *uint
	if params.VaultId != nil {
		// #nosec G115
		vaultIDValue := uint(*params.VaultId)
		vaultID = &vaultIDValue
	}

	filterParams := model.GetAuditLogsWithFiltersParams{
		UserID:    user.ID,
		VaultID:   vaultID,
		StartDate: params.StartDate,
		EndDate:   params.EndDate,
		Limit:     params.PageSize,
		Offset:    offset,
	}

	// Get audit logs with filters
	logs, err := model.GetAuditLogsWithFilters(filterParams)
	if err != nil {
		return handler.SendError(c, fiber.StatusInternalServerError, err.Error())
	}

	// Get total count for pagination
	totalCount, err := model.CountAuditLogsWithFilters(filterParams)
	if err != nil {
		return handler.SendError(c, fiber.StatusInternalServerError, err.Error())
	}

	// Convert to API format
	apiLogs := make([]AuditLog, len(logs))
	for i, log := range logs {
		apiLogs[i] = convertToApiAuditLog(&log)
	}

	// Prepare response
	response := AuditLogsResponse{
		AuditLogs:  apiLogs,
		TotalCount: int(totalCount),
		PageSize:   params.PageSize,
		PageIndex:  params.PageIndex,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
