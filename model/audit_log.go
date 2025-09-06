package model

import (
	"time"

	"gorm.io/gorm"
)

type ActionType string

const (
	ActionLoginUser    ActionType = "login_user"
	ActionRegisterUser ActionType = "register_user"
	ActionLogoutUser   ActionType = "logout_user"
	ActionReadVault    ActionType = "read_vault"
	ActionUpdateVault  ActionType = "update_vault"
	ActionDeleteVault  ActionType = "delete_vault"
	ActionCreateVault  ActionType = "create_vault"
	ActionCreateAPIKey ActionType = "create_api_key"
	//nolint:gosec // G101 here is the enum name
	ActionUpdateAPIKey ActionType = "update_api_key"
	ActionDeleteAPIKey ActionType = "delete_api_key"
)

type AuditLog struct {
	gorm.Model
	VaultID   *uint      `gorm:"index"`
	Vault     *Vault     `gorm:"foreignKey:VaultID"`
	APIKeyID  *uint      `gorm:"index"`
	APIKey    *APIKey    `gorm:"foreignKey:APIKeyID"`
	Action    ActionType `gorm:"size:50;index"`
	UserID    uint       `gorm:"index;constraint:OnDelete:CASCADE"`
	User      User       `gorm:"foreignKey:UserID"`
	IPAddress string     `gorm:"size:45"`
	UserAgent string     `gorm:"size:500"`
}

// CreateAuditLogParams defines parameters for creating an audit log entry
type CreateAuditLogParams struct {
	VaultID   *uint
	APIKeyID  *uint
	Action    ActionType
	UserID    uint
	IPAddress string
	UserAgent string
}

// CreateAuditLog creates a new audit log entry
func CreateAuditLog(params CreateAuditLogParams) error {
	auditLog := AuditLog{
		VaultID:   params.VaultID,
		APIKeyID:  params.APIKeyID,
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

// LogAPIKeyAction logs an API key-related action
func LogAPIKeyAction(apiKeyID uint, action ActionType, userID uint, ipAddress, userAgent string) error {
	return CreateAuditLog(CreateAuditLogParams{
		APIKeyID:  &apiKeyID,
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
		Preload("Vault", func(db *gorm.DB) *gorm.DB {
			return db.Unscoped()
		}).
		Preload("APIKey", func(db *gorm.DB) *gorm.DB {
			return db.Unscoped()
		}).
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
		Preload("Vault", func(db *gorm.DB) *gorm.DB {
			return db.Unscoped()
		}).
		Preload("APIKey", func(db *gorm.DB) *gorm.DB {
			return db.Unscoped()
		}).
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
		Preload("Vault", func(db *gorm.DB) *gorm.DB {
			return db.Unscoped()
		}).
		Preload("APIKey", func(db *gorm.DB) *gorm.DB {
			return db.Unscoped()
		}).
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

// AuditMetrics holds all audit metrics for efficient single-query retrieval
type AuditMetrics struct {
	TotalEventsLast30Days  int64
	EventsCountLast24Hours int64
	VaultEventsLast30Days  int64
	APIKeyEventsLast30Days int64
}

// GetAllAuditMetrics retrieves all audit metrics in a single optimized query
func GetAllAuditMetrics(userID uint) (*AuditMetrics, error) {
	now := time.Now()
	thirtyDaysAgo := now.AddDate(0, 0, -30)
	twentyFourHoursAgo := now.Add(-24 * time.Hour)

	var result struct {
		TotalEventsLast30Days  int64
		EventsCountLast24Hours int64
		VaultEventsLast30Days  int64
		APIKeyEventsLast30Days int64
	}

	vaultActions := []string{
		string(ActionReadVault),
		string(ActionUpdateVault),
		string(ActionDeleteVault),
		string(ActionCreateVault),
	}

	apiKeyActions := []string{
		string(ActionCreateAPIKey),
		string(ActionUpdateAPIKey),
		string(ActionDeleteAPIKey),
	}

	err := DB.Model(&AuditLog{}).
		Select(`
			COUNT(CASE WHEN created_at >= ? THEN 1 END) as total_events_last30_days,
			COUNT(CASE WHEN created_at >= ? THEN 1 END) as events_count_last24_hours,
			COUNT(CASE WHEN created_at >= ? AND action IN ? THEN 1 END) as vault_events_last30_days,
			COUNT(CASE WHEN created_at >= ? AND action IN ? THEN 1 END) as api_key_events_last30_days
		`, thirtyDaysAgo, twentyFourHoursAgo, thirtyDaysAgo, vaultActions, thirtyDaysAgo, apiKeyActions).
		Where("user_id = ?", userID).
		Scan(&result).Error

	if err != nil {
		return nil, err
	}

	return &AuditMetrics{
		TotalEventsLast30Days:  result.TotalEventsLast30Days,
		EventsCountLast24Hours: result.EventsCountLast24Hours,
		VaultEventsLast30Days:  result.VaultEventsLast30Days,
		APIKeyEventsLast30Days: result.APIKeyEventsLast30Days,
	}, nil
}
