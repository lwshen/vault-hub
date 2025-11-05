package models

import (
	"time"
)

type Vault struct {

	// Unique identifier for the vault
	UniqueId string `json:"uniqueId"`

	// ID of the user who owns this vault
	UserId int64 `json:"userId,omitempty"`

	// Human-readable name
	Name string `json:"name"`

	// Encrypted value
	Value string `json:"value"`

	// Human-readable description
	Description string `json:"description,omitempty"`

	// Category/type of vault
	Category string `json:"category,omitempty"`

	CreatedAt *time.Time `json:"createdAt,omitempty"`

	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
}
