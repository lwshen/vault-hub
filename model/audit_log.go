package model

import (
	"time"

	"gorm.io/gorm"
)

type ActionType string

const (
	ActionReadConfig   ActionType = "read_config"
	ActionUpdateConfig ActionType = "update_config"
	ActionDeleteConfig ActionType = "delete_config"
	ActionCreateConfig ActionType = "create_config"
	ActionLoginUser    ActionType = "login_user"
	ActionRegisterUser ActionType = "register_user"
	ActionLogoutUser   ActionType = "logout_user"
)

type AuditLog struct {
	gorm.Model
	ConfigID  *uint      `gorm:"index"`
	Action    ActionType `gorm:"size:50;index"`
	UserID    uint       `gorm:"index;constraint:OnDelete:CASCADE"`
	User      User       `gorm:"foreignKey:UserID"`
	IPAddress string     `gorm:"size:45"`
	UserAgent string     `gorm:"size:500"`
	Timestamp time.Time  `gorm:"index"`
}

// CreateAuditLogParams defines parameters for creating an audit log entry
type CreateAuditLogParams struct {
	ConfigID  *uint
	Action    ActionType
	UserID    uint
	IPAddress string
	UserAgent string
}

// CreateAuditLog creates a new audit log entry
func CreateAuditLog(params CreateAuditLogParams) error {
	auditLog := AuditLog{
		ConfigID:  params.ConfigID,
		Action:    params.Action,
		UserID:    params.UserID,
		IPAddress: params.IPAddress,
		UserAgent: params.UserAgent,
		Timestamp: time.Now(),
	}

	err := DB.Create(&auditLog).Error
	if err != nil {
		return err
	}

	return nil
}

// LogConfigurationAction logs a configuration-related action
func LogConfigurationAction(configID uint, action ActionType, userID uint, ipAddress, userAgent string) error {
	return CreateAuditLog(CreateAuditLogParams{
		ConfigID:  &configID,
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
		Order("timestamp DESC")

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

// GetAuditLogsByConfiguration retrieves audit logs for a specific configuration
func GetAuditLogsByConfiguration(configID uint, userID uint) ([]AuditLog, error) {
	var logs []AuditLog
	err := DB.Where("config_id = ? AND user_id = ?", configID, userID).
		Preload("User").
		Order("timestamp DESC").
		Find(&logs).Error

	if err != nil {
		return nil, err
	}

	return logs, nil
}
