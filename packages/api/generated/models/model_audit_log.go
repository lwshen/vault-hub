package models

import (
	"time"
)

type AuditLog struct {

	// When the action occurred
	CreatedAt time.Time `json:"createdAt"`

	Vault VaultLite `json:"vault,omitempty"`

	ApiKey VaultApiKey `json:"apiKey,omitempty"`

	// Type of action performed
	Action string `json:"action"`

	// Source of the request (web interface or CLI)
	Source string `json:"source"`

	// IP address from which the action was performed
	IpAddress string `json:"ipAddress,omitempty"`

	// User agent string from the client
	UserAgent string `json:"userAgent,omitempty"`
}
