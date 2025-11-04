package generated_models

import (
	"time"
)

type VaultApiKey struct {

	// Unique API key ID
	Id int64 `json:"id"`

	// Human-readable name for the API key
	Name string `json:"name"`

	// Array of vaults this key can access (null/empty = all user's vaults)
	Vaults []VaultLite `json:"vaults,omitempty"`

	// Optional expiration date
	ExpiresAt *time.Time `json:"expiresAt,omitempty"`

	// When the key was last used
	LastUsedAt *time.Time `json:"lastUsedAt,omitempty"`

	// Whether the key is currently active
	IsActive bool `json:"isActive"`

	// When the key was created
	CreatedAt time.Time `json:"createdAt"`

	// When the key was last updated
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
}
