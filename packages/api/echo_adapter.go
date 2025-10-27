package api

import (
	"time"
	"github.com/lwshen/vault-hub/model"
)

// EchoAdapter provides conversion between current implementation models and value-based models
type EchoAdapter struct{}

// NewEchoAdapter creates a new adapter instance
func NewEchoAdapter() *EchoAdapter {
	return &EchoAdapter{}
}

// EchoEchoGetUserResponse represents user data for Echo API responses (value-based)
type EchoEchoGetUserResponse struct {
	Id        int64     `json:"id"`
	Email     string     `json:"email"`
	Name      string     `json:"name"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
}

// EchoVault represents vault data with decrypted value for Echo (value-based)
type EchoVault struct {
	UniqueId    string    `json:"uniqueId"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Value       string    `json:"value"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	IsActive    bool      `json:"isActive"`
}

// EchoVaultLite represents vault data without decrypted value for Echo (value-based)
type EchoVaultLite struct {
	UniqueId  string    `json:"uniqueId"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	IsActive  bool      `json:"isActive"`
}

// EchoAPIKey represents API key data for Echo (value-based)
type EchoAPIKey struct {
	Id         int64      `json:"id"`
	Name       string      `json:"name"`
	Key        string      `json:"key"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  time.Time  `json:"updatedAt"`
	LastUsedAt interface{} `json:"lastUsedAt,omitempty"`
	IsActive   bool       `json:"isActive"`
}

// EchoAuditLog represents audit log data for Echo (value-based)
type EchoAuditLog struct {
	Id       int64  `json:"id"`
	UserId   int64   `json:"userId"`
	Action   string  `json:"action"`
	Resource string  `json:"resource,omitempty"`
	Source   string  `json:"source"`
	ClientIp string  `json:"clientIp,omitempty"`
	UserAgent string  `json:"userAgent,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
}

// EchoStatusResponse represents system status for Echo (value-based)
type EchoStatusResponse struct {
	Version  string                  `json:"version"`
	Uptime   string                  `json:"uptime"`
	Database EchoStatusResponseDatabase  `json:"database"`
	System   EchoStatusResponseSystem   `json:"system"`
}

// EchoStatusResponseDatabase represents database status for Echo
type EchoStatusResponseDatabase struct {
	Status      string                        `json:"status"`
	ResponseTime float64                       `json:"responseTime"`
	Connections EchoStatusResponseDatabaseConnections `json:"connections"`
}

// EchoStatusResponseDatabaseConnections represents database connection info for Echo
type EchoStatusResponseDatabaseConnections struct {
	Active int `json:"active"`
	Idle   int `json:"idle"`
}

// EchoStatusResponseSystem represents system status for Echo
type EchoStatusResponseSystem struct {
	Status     string `json:"status"`
	MemoryUsed string `json:"memoryUsed"`
	DiskFree   string `json:"diskFree"`
}

// EchoHealthCheckResponse represents health check response for Echo (value-based)
type EchoHealthCheckResponse struct {
	Status    string    `json:"status,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}

// EchoLoginResponse represents login response for Echo (value-based)
type EchoLoginResponse struct {
	Token string `json:"token"`
}

// EchoConfigResponse represents application config for Echo (value-based)
type EchoConfigResponse struct {
	IsOidcEnabled      bool `json:"isOidcEnabled"`
	IsEmailEnabled      bool `json:"isEmailEnabled"`
	PasswordMinLength   int64 `json:"passwordMinLength"`
	IsRegistrationOpen  bool `json:"isRegistrationOpen"`
}

// Conversion methods

func (a *EchoAdapter) ConvertUser(user *model.User) EchoGetUserResponse {
	if user == nil {
		return EchoGetUserResponse{}
	}

	return EchoGetUserResponse{
		Id:        int64(user.ID),
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func (a *EchoAdapter) ConvertEchoVault(vault *model.EchoVault, decryptedValue string) EchoVault {
	if vault == nil {
		return EchoVault{}
	}

	return EchoVault{
		UniqueId:    vault.UniqueID,
		Name:        vault.Name,
		Description: a.stringFromPtr(vault.Description),
		Value:       decryptedValue,
		CreatedAt:   vault.CreatedAt,
		UpdatedAt:   vault.UpdatedAt,
		IsActive:    vault.IsActive,
	}
}

func (a *EchoAdapter) ConvertEchoVaultLite(vault *model.EchoVault) EchoVaultLite {
	if vault == nil {
		return EchoVaultLite{}
	}

	return EchoVaultLite{
		UniqueId:  vault.UniqueID,
		Name:      vault.Name,
		CreatedAt: vault.CreatedAt,
		UpdatedAt: vault.UpdatedAt,
		IsActive:  vault.IsActive,
	}
}

func (a *EchoAdapter) ConvertEchoAPIKey(apiKey *model.EchoAPIKey) EchoAPIKey {
	if apiKey == nil {
		return EchoAPIKey{}
	}

	return EchoAPIKey{
		Id:         int64(apiKey.ID),
		Name:       apiKey.Name,
		Key:        apiKey.Key,
		CreatedAt:  apiKey.CreatedAt,
		UpdatedAt:  apiKey.UpdatedAt,
		LastUsedAt: a.timeFromPtr(apiKey.LastUsedAt),
		IsActive:   apiKey.IsActive,
	}
}

func (a *EchoAdapter) ConvertEchoAuditLog(auditLog *model.EchoAuditLog) EchoAuditLog {
	if auditLog == nil {
		return EchoAuditLog{}
	}

	return EchoAuditLog{
		Id:        int64(auditLog.ID),
		UserId:    a.int64FromPtr(auditLog.UserID),
		Action:    string(auditLog.Action),
		Resource:  a.stringFromPtr(auditLog.Resource),
		Source:    string(auditLog.Source),
		ClientIp:  a.stringFromPtr(auditLog.ClientIP),
		UserAgent:  a.stringFromPtr(auditLog.UserAgent),
		CreatedAt: auditLog.CreatedAt,
	}
}

func (a *EchoAdapter) ConvertEchoStatusResponse(status *model.EchoStatusResponse) EchoStatusResponse {
	if status == nil {
		return EchoStatusResponse{}
	}

	return EchoStatusResponse{
		Version: status.Version,
		Uptime:  status.Uptime,
		Database: EchoStatusResponseDatabase{
			Status:      status.Database.Status,
			ResponseTime: status.Database.ResponseTime,
			Connections: EchoStatusResponseDatabaseConnections{
				Active: status.Database.Connections.Active,
				Idle:   status.Database.Connections.Idle,
			},
		},
		System: EchoStatusResponseSystem{
			Status:     status.System.Status,
			MemoryUsed: status.System.MemoryUsed,
			DiskFree:   status.System.DiskFree,
		},
	}
}

// Helper functions for pointer to value conversion

func (a *EchoAdapter) stringFromPtr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func (a *EchoAdapter) int64FromPtr(i *uint) int64 {
	if i == nil {
		return 0
	}
	return int64(*i)
}

func (a *EchoAdapter) timeFromPtr(t *time.Time) interface{} {
	if t == nil {
		return nil
	}
	return *t
}

func (a *EchoAdapter) ConvertEchoHealthCheckResponse(status string, timestamp time.Time) EchoHealthCheckResponse {
	return EchoHealthCheckResponse{
		Status:    status,
		Timestamp: timestamp,
	}
}

func (a *EchoAdapter) ConvertEchoLoginResponse(token string) EchoLoginResponse {
	return EchoLoginResponse{
		Token: token,
	}
}

func (a *EchoAdapter) ConvertEchoConfigResponse() EchoConfigResponse {
	return EchoConfigResponse{
		IsOidcEnabled:    model.IsOidcEnabled(),
		IsEmailEnabled:     model.IsEmailEnabled(),
		PasswordMinLength:   int64(model.ConfigPasswordMinLength),
		IsRegistrationOpen: model.IsRegistrationOpen(),
	}
}