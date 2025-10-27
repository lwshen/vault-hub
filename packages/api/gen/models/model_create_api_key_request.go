package models

import (
	"time"
)

type CreateApiKeyRequest struct {

	// Human-readable name for the API key
	Name string `json:"name"`

	// Array of vault unique IDs this key can access (empty = all user's vaults)
	VaultUniqueIds []string `json:"vaultUniqueIds,omitempty"`

	// Optional expiration date
	ExpiresAt time.Time `json:"expiresAt,omitempty"`
}
