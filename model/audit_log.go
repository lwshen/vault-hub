package model

import (
	"time"

	"gorm.io/gorm"
)

type ActionType string

const (
	ActionReadVault    ActionType = "read_vault"
	ActionUpdateVault  ActionType = "update_vault"
	ActionDeleteVault  ActionType = "delete_vault"
	ActionCreateVault  ActionType = "create_vault"
	ActionLoginUser    ActionType = "login_user"
	ActionRegisterUser ActionType = "register_user"
	ActionLogoutUser   ActionType = "logout_user"
	ActionCreateAPIKey ActionType = "create_api_key"
	ActionUpdateAPIKey ActionType = "update_api_key"
	ActionDeleteAPIKey ActionType = "delete_api_key"
)

type AuditLog struct {
	gorm.Model
	VaultID   *uint      `gorm:"index"`
	Vault     *Vault     `gorm:"foreignKey:VaultID"`
	Action    ActionType `gorm:"size:50;index"`
	UserID    uint       `gorm:"index;constraint:OnDelete:CASCADE"`
	User      User       `gorm:"foreignKey:UserID"`
	IPAddress string     `gorm:"size:45"`
	UserAgent string     `gorm:"size:500"`
}

// CreateAuditLogParams defines parameters for creating an audit log entry
type CreateAuditLogParams struct {
	VaultID   *uint
	Action    ActionType
	UserID    uint
	IPAddress string
	UserAgent string
}

// CreateAuditLog creates a new audit log entry
func CreateAuditLog(params CreateAuditLogParams) error {
	auditLog := AuditLog{
		VaultID:   params.VaultID,
		Action:    params.Action,
		UserID:    params.UserID,
		IPAddress: params.IPAddress,
		UserAgent: params.UserAgent,
	}

	err := DB.Create(&auditLog).Error
	if err != nil {
		return err
	}

	return nil
}

// LogVaultAction logs a vault-related action
func LogVaultAction(vaultID uint, action ActionType, userID uint, ipAddress, userAgent string) error {
	return CreateAuditLog(CreateAuditLogParams{
		VaultID:   &vaultID,
		Action:    action,
		UserID:    userID,
		IPAddress: ipAddress,
		UserAgent: userAgent,
	})
}

// LogUserAction logs a user-related action (login, register, logout)
func LogUserAction(action ActionType, userID uint, ipAddress, userAgent string) error {
	return CreateAuditLog(CreateAuditLogParams{
		Action:    action,
		UserID:    userID,
		IPAddress: ipAddress,
		UserAgent: userAgent,
	})
}

// GetAuditLogsByUser retrieves audit logs for a specific user
func GetAuditLogsByUser(userID uint, limit int, offset int) ([]AuditLog, error) {
	var logs []AuditLog
	query := DB.Where("user_id = ?", userID).
		Preload("User").
		Preload("Vault").
		Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Find(&logs).Error
	if err != nil {
		return nil, err
	}

	return logs, nil
}

// GetAuditLogsByVault retrieves audit logs for a specific vault
func GetAuditLogsByVault(vaultID uint, userID uint) ([]AuditLog, error) {
	var logs []AuditLog
	err := DB.Where("vault_id = ? AND user_id = ?", vaultID, userID).
		Preload("User").
		Preload("Vault").
		Order("created_at DESC").
		Find(&logs).Error

	if err != nil {
		return nil, err
	}

	return logs, nil
}

// GetAuditLogsWithFiltersParams defines parameters for filtering audit logs
type GetAuditLogsWithFiltersParams struct {
	UserID    uint
	VaultID   *uint
	StartDate *time.Time
	EndDate   *time.Time
	Limit     int
	Offset    int
}

// GetAuditLogsWithFilters retrieves audit logs with optional filtering and pagination
func GetAuditLogsWithFilters(params GetAuditLogsWithFiltersParams) ([]AuditLog, error) {
	var logs []AuditLog
	query := DB.Where("user_id = ?", params.UserID).
		Preload("User").
		Preload("Vault").
		Order("created_at DESC")

	// Add vault filter if specified
	if params.VaultID != nil {
		query = query.Where("vault_id = ?", *params.VaultID)
	}

	// Add date range filter if specified
	if params.StartDate != nil {
		query = query.Where("created_at >= ?", *params.StartDate)
	}
	if params.EndDate != nil {
		query = query.Where("created_at <= ?", *params.EndDate)
	}

	// Add pagination
	if params.Limit > 0 {
		query = query.Limit(params.Limit)
	}
	if params.Offset > 0 {
		query = query.Offset(params.Offset)
	}

	err := query.Find(&logs).Error
	if err != nil {
		return nil, err
	}

	return logs, nil
}

// CountAuditLogsWithFilters counts total audit logs matching the filter criteria
func CountAuditLogsWithFilters(params GetAuditLogsWithFiltersParams) (int64, error) {
	var count int64
	query := DB.Model(&AuditLog{}).Where("user_id = ?", params.UserID)

	// Add vault filter if specified
	if params.VaultID != nil {
		query = query.Where("vault_id = ?", *params.VaultID)
	}

	// Add date range filter if specified
	if params.StartDate != nil {
		query = query.Where("created_at >= ?", *params.StartDate)
	}
	if params.EndDate != nil {
		query = query.Where("created_at <= ?", *params.EndDate)
	}

	err := query.Count(&count).Error
	if err != nil {
		return 0, err
	}

	return count, nil
}
