package generated_models

import (
	"time"
)

type VaultLite struct {

	// Unique identifier for the vault
	UniqueId string `json:"uniqueId"`

	// Human-readable name
	Name string `json:"name"`

	// Human-readable description
	Description string `json:"description,omitempty"`

	// Category/type of vault
	Category string `json:"category,omitempty"`

	UpdatedAt time.Time `json:"updatedAt,omitempty"`
}
