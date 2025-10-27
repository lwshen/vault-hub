package models

type CreateVaultRequest struct {

	// Human-readable name
	Name string `json:"name"`

	// Value to be encrypted and stored
	Value string `json:"value"`

	// Human-readable description
	Description string `json:"description,omitempty"`

	// Category/type of vault
	Category string `json:"category,omitempty"`
}
