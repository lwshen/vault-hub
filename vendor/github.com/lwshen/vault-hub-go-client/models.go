package openapi

import "time"

// VaultLite matches fields used by CLI output
type VaultLite struct {
	UniqueId    string     `json:"uniqueId"`
	Name        string     `json:"name"`
	Description *string    `json:"description,omitempty"`
	Category    *string    `json:"category,omitempty"`
	UpdatedAt   *time.Time `json:"updatedAt,omitempty"`
}

func (v *VaultLite) GetUniqueId() string { if v == nil { return "" }; return v.UniqueId }
func (v *VaultLite) GetName() string     { if v == nil { return "" }; return v.Name }
func (v *VaultLite) GetDescription() string { if v == nil || v.Description == nil { return "" }; return *v.Description }
func (v *VaultLite) GetDescriptionOk() *string { if v == nil || v.Description == nil { return nil }; return v.Description }

// Vault includes encrypted Value returned by get endpoints
type Vault struct {
	UniqueId    string     `json:"uniqueId"`
	Name        string     `json:"name"`
	Value       string     `json:"value"`
	Description *string    `json:"description,omitempty"`
	Category    *string    `json:"category,omitempty"`
	CreatedAt   *time.Time `json:"createdAt,omitempty"`
	UpdatedAt   *time.Time `json:"updatedAt,omitempty"`
}

func (v *Vault) GetUniqueId() string { if v == nil { return "" }; return v.UniqueId }
func (v *Vault) GetName() string     { if v == nil { return "" }; return v.Name }
func (v *Vault) GetValue() string    { if v == nil { return "" }; return v.Value }
func (v *Vault) GetDescription() string { if v == nil || v.Description == nil { return "" }; return *v.Description }
func (v *Vault) GetDescriptionOk() *string { if v == nil || v.Description == nil { return nil }; return v.Description }
func (v *Vault) GetCategory() string { if v == nil || v.Category == nil { return "" }; return *v.Category }
func (v *Vault) GetCategoryOk() *string { if v == nil || v.Category == nil { return nil }; return v.Category }

