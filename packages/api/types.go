package api

import "time"

// HealthCheckResponse represents the response body for /api/health.
type HealthCheckResponse struct {
	Status    *string    `json:"status,omitempty"`
	Timestamp *time.Time `json:"timestamp,omitempty"`
}

// ConfigResponse exposes non-sensitive configuration flags.
type ConfigResponse struct {
	EmailEnabled bool `json:"emailEnabled"`
	OidcEnabled  bool `json:"oidcEnabled"`
}

// Status enumerations.
type (
	StatusResponseDatabaseStatus string
	StatusResponseSystemStatus   string
)

const (
	StatusResponseDatabaseStatusHealthy     StatusResponseDatabaseStatus = "healthy"
	StatusResponseDatabaseStatusDegraded    StatusResponseDatabaseStatus = "degraded"
	StatusResponseDatabaseStatusUnavailable StatusResponseDatabaseStatus = "unavailable"
)

const (
	StatusResponseSystemStatusHealthy     StatusResponseSystemStatus = "healthy"
	StatusResponseSystemStatusDegraded    StatusResponseSystemStatus = "degraded"
	StatusResponseSystemStatusUnavailable StatusResponseSystemStatus = "unavailable"
)

// StatusResponse aggregates system health information.
type StatusResponse struct {
	Version        string                       `json:"version"`
	Commit         string                       `json:"commit"`
	SystemStatus   StatusResponseSystemStatus   `json:"systemStatus"`
	DatabaseStatus StatusResponseDatabaseStatus `json:"databaseStatus"`
}

// Vault models.
type Vault struct {
	UniqueId    string     `json:"uniqueId"`
	UserId      *int64     `json:"userId,omitempty"`
	Name        string     `json:"name"`
	Value       string     `json:"value"`
	Description *string    `json:"description,omitempty"`
	Category    *string    `json:"category,omitempty"`
	CreatedAt   *time.Time `json:"createdAt,omitempty"`
	UpdatedAt   *time.Time `json:"updatedAt,omitempty"`
}

type VaultLite struct {
	UniqueId    string     `json:"uniqueId"`
	Name        string     `json:"name"`
	Description *string    `json:"description,omitempty"`
	Category    *string    `json:"category,omitempty"`
	UpdatedAt   *time.Time `json:"updatedAt,omitempty"`
}

type VaultsResponse struct {
	Vaults     []VaultLite `json:"vaults"`
	TotalCount int         `json:"totalCount"`
	PageSize   int         `json:"pageSize"`
	PageIndex  int         `json:"pageIndex"`
}

type GetVaultsParams struct {
	PageSize  *int `json:"pageSize,omitempty"`
	PageIndex *int `json:"pageIndex,omitempty"`
}

type CreateVaultRequest struct {
	Name        string  `json:"name"`
	Value       string  `json:"value"`
	Description *string `json:"description,omitempty"`
	Category    *string `json:"category,omitempty"`
}

type UpdateVaultRequest struct {
	Name        *string `json:"name,omitempty"`
	Value       *string `json:"value,omitempty"`
	Description *string `json:"description,omitempty"`
	Category    *string `json:"category,omitempty"`
}

type VaultAPIKey struct {
	Id         int64        `json:"id"`
	Name       string       `json:"name"`
	Vaults     *[]VaultLite `json:"vaults,omitempty"`
	ExpiresAt  *time.Time   `json:"expiresAt,omitempty"`
	LastUsedAt *time.Time   `json:"lastUsedAt,omitempty"`
	IsActive   bool         `json:"isActive"`
	CreatedAt  time.Time    `json:"createdAt"`
	UpdatedAt  *time.Time   `json:"updatedAt,omitempty"`
}

type APIKeysResponse struct {
	ApiKeys    []VaultAPIKey `json:"apiKeys"`
	TotalCount int           `json:"totalCount"`
	PageSize   int           `json:"pageSize"`
	PageIndex  int           `json:"pageIndex"`
}

type CreateAPIKeyRequest struct {
	Name           string     `json:"name"`
	VaultUniqueIds *[]string  `json:"vaultUniqueIds,omitempty"`
	ExpiresAt      *time.Time `json:"expiresAt,omitempty"`
}

type CreateAPIKeyResponse struct {
	ApiKey VaultAPIKey `json:"apiKey"`
	Key    string      `json:"key"`
}

type UpdateAPIKeyRequest struct {
	Name           *string    `json:"name,omitempty"`
	VaultUniqueIds *[]string  `json:"vaultUniqueIds,omitempty"`
	ExpiresAt      *time.Time `json:"expiresAt,omitempty"`
}

type GetAPIKeysParams struct {
	PageSize  int `json:"pageSize"`
	PageIndex int `json:"pageIndex"`
}

// Auth payloads.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type SignupRequest struct {
	Email    string  `json:"email"`
	Password string  `json:"password"`
	Name     *string `json:"name,omitempty"`
}

type SignupResponse struct {
	Token string `json:"token"`
}

type PasswordResetRequest struct {
	Email string `json:"email"`
}

type PasswordResetConfirmRequest struct {
	Token       string `json:"token"`
	NewPassword string `json:"newPassword"`
}

type MagicLinkRequest struct {
	Email string `json:"email"`
}

type EmailTokenResponse struct {
	Success bool   `json:"success"`
	Code    string `json:"code"`
}

type ConsumeMagicLinkParams struct {
	Token string `json:"token"`
}

// Audit log types.
type (
	AuditLogAction string
	AuditLogSource string
)

const (
	// #nosec G101 -- These values are audit action identifiers, not credentials.
	AuditLogActionCreateApiKey AuditLogAction = "create_api_key"
	AuditLogActionCreateVault  AuditLogAction = "create_vault"
	AuditLogActionDeleteApiKey AuditLogAction = "delete_api_key"
	AuditLogActionDeleteVault  AuditLogAction = "delete_vault"
	AuditLogActionLoginUser    AuditLogAction = "login_user"
	AuditLogActionLogoutUser   AuditLogAction = "logout_user"
	AuditLogActionReadVault    AuditLogAction = "read_vault"
	AuditLogActionRegisterUser AuditLogAction = "register_user"
	AuditLogActionUpdateApiKey AuditLogAction = "update_api_key" // #nosec G101 -- action identifier, not a credential
	AuditLogActionUpdateVault  AuditLogAction = "update_vault"
)

const (
	AuditLogSourceWeb AuditLogSource = "web"
	AuditLogSourceCLI AuditLogSource = "cli"
)

type AuditLog struct {
	Action    AuditLogAction `json:"action"`
	CreatedAt time.Time      `json:"createdAt"`
	Vault     *VaultLite     `json:"vault,omitempty"`
	ApiKey    *VaultAPIKey   `json:"apiKey,omitempty"`
	Source    AuditLogSource `json:"source"`
	IpAddress *string        `json:"ipAddress,omitempty"`
	UserAgent *string        `json:"userAgent,omitempty"`
}

type AuditLogsResponse struct {
	AuditLogs  []AuditLog `json:"auditLogs"`
	TotalCount int        `json:"totalCount"`
	PageSize   int        `json:"pageSize"`
	PageIndex  int        `json:"pageIndex"`
}

type AuditMetricsResponse struct {
	ApiKeyEventsLast30Days int `json:"apiKeyEventsLast30Days"`
	EventsCountLast24Hours int `json:"eventsCountLast24Hours"`
	TotalEventsLast30Days  int `json:"totalEventsLast30Days"`
	VaultEventsLast30Days  int `json:"vaultEventsLast30Days"`
}

type GetAuditLogsParams struct {
	StartDate     *time.Time `json:"startDate,omitempty"`
	EndDate       *time.Time `json:"endDate,omitempty"`
	VaultUniqueId *string    `json:"vaultUniqueId,omitempty"`
	PageSize      int        `json:"pageSize"`
	PageIndex     int        `json:"pageIndex"`
}

// CLI parameter helper structs.
type GetVaultByAPIKeyParams struct {
	XEnableClientEncryption *string `json:"X-Enable-Client-Encryption,omitempty"`
}

type GetVaultByNameAPIKeyParams struct {
	XEnableClientEncryption *string `json:"X-Enable-Client-Encryption,omitempty"`
}

// User response.
type GetUserResponse struct {
	Email  string  `json:"email"`
	Avatar *string `json:"avatar,omitempty"`
	Name   *string `json:"name,omitempty"`
}
