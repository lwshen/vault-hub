package generated_models

import (
	"time"
)

type UpdateApiKeyRequest struct {

	// Human-readable name for the API key
	Name string `json:"name,omitempty"`

	// Array of vault unique IDs this key can access (empty = all user's vaults)
	VaultUniqueIds []string `json:"vaultUniqueIds,omitempty"`

	// Optional expiration date
	ExpiresAt *time.Time `json:"expiresAt,omitempty"`
}
