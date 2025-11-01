package generated_models

type UpdateVaultRequest struct {

	// Human-readable name
	Name string `json:"name,omitempty"`

	// Value to be encrypted and stored
	Value string `json:"value,omitempty"`

	// Human-readable description
	Description string `json:"description,omitempty"`

	// Category/type of vault
	Category string `json:"category,omitempty"`
}
