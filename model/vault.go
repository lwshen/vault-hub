package model

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/lwshen/vault-hub/internal/encryption"
	"gorm.io/gorm"
)

type Vault struct {
	gorm.Model
	UniqueID    string `gorm:"size:255;not null;unique"` // Unique identifier for the vault
	UserID      uint   `gorm:"index;not null"`           // User who owns this vault
	Name        string `gorm:"size:255"`                 // Human-readable name
	Value       string `gorm:"type:text;not null"`       // Encrypted value
	Description string `gorm:"size:500"`                 // Human-readable description
	Category    string `gorm:"size:100;index"`           // Category/type of vault
}

// CreateVaultParams defines parameters for creating a new vault
type CreateVaultParams struct {
	UniqueID    string
	UserID      uint
	Name        string
	Value       string
	Description string
	Category    string
}

// UpdateVaultParams defines parameters for updating a vault
type UpdateVaultParams struct {
	Name        *string
	Value       *string
	Description *string
	Category    *string
}

// Validate validates the create vault parameters
func (params *CreateVaultParams) Validate() map[string]string {
	errors := map[string]string{}

	if strings.TrimSpace(params.UniqueID) == "" {
		errors["unique_id"] = "unique_id is required"
	} else if len(params.UniqueID) > 255 {
		errors["unique_id"] = "unique_id must be less than 255 characters"
	}

	if strings.TrimSpace(params.Name) == "" {
		errors["name"] = "name is required"
	} else if len(params.Name) > 255 {
		errors["name"] = "name must be less than 255 characters"
	}

	if strings.TrimSpace(params.Value) == "" {
		errors["value"] = "value is required"
	}

	if len(params.Description) > 500 {
		errors["description"] = "description must be less than 500 characters"
	}

	if len(params.Category) > 100 {
		errors["category"] = "category must be less than 100 characters"
	}

	if params.UserID == 0 {
		errors["user_id"] = "user_id is required"
	}

	return errors
}

// Validate validates the update vault parameters
func (params *UpdateVaultParams) Validate() map[string]string {
	errors := map[string]string{}

	if params.Name != nil {
		if strings.TrimSpace(*params.Name) == "" {
			errors["name"] = "name cannot be empty"
		} else if len(*params.Name) > 255 {
			errors["name"] = "name must be less than 255 characters"
		}
	}

	if params.Value != nil && strings.TrimSpace(*params.Value) == "" {
		errors["value"] = "value cannot be empty"
	}

	if params.Description != nil && len(*params.Description) > 500 {
		errors["description"] = "description must be less than 500 characters"
	}

	if params.Category != nil && len(*params.Category) > 100 {
		errors["category"] = "category must be less than 100 characters"
	}

	return errors
}

// Create creates a new vault
func (params *CreateVaultParams) Create() (*Vault, error) {
	// Check if unique_id already exists for this user
	var existing Vault
	err := DB.Where("unique_id = ? AND user_id = ?", params.UniqueID, params.UserID).First(&existing).Error
	if err == nil {
		return nil, fmt.Errorf("vault with unique_id '%s' already exists", params.UniqueID)
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Encrypt the value before storing
	encryptedValue, err := encryption.Encrypt(params.Value)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt value: %w", err)
	}

	vault := Vault{
		UniqueID:    params.UniqueID,
		UserID:      params.UserID,
		Name:        params.Name,
		Value:       encryptedValue,
		Description: params.Description,
		Category:    params.Category,
	}

	err = DB.Create(&vault).Error
	if err != nil {
		return nil, err
	}

	// Decrypt the value for the response
	vault.Value, err = encryption.Decrypt(vault.Value)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt value for response: %w", err)
	}

	return &vault, nil
}

// GetByUniqueID retrieves a vault by unique_id for a specific user
func (v *Vault) GetByUniqueID(uniqueID string, userID uint) error {
	err := DB.Where("unique_id = ? AND user_id = ?", uniqueID, userID).First(v).Error
	if err != nil {
		return err
	}

	// Decrypt the value
	decryptedValue, err := encryption.Decrypt(v.Value)
	if err != nil {
		return fmt.Errorf("failed to decrypt value: %w", err)
	}
	v.Value = decryptedValue

	return nil
}

// GetAllByUser retrieves all vaults for a user, optionally filtered by category
func GetVaultsByUser(userID uint, decrypt bool) ([]Vault, error) {
	var vaults []Vault
	query := DB.Where("user_id = ?", userID)

	err := query.Order("created_at DESC").Find(&vaults).Error
	if err != nil {
		return nil, err
	}

	// Decrypt all values if requested
	if decrypt {
		for i := range vaults {
			decryptedValue, err := encryption.Decrypt(vaults[i].Value)
			if err != nil {
				return nil, fmt.Errorf("failed to decrypt value for vault %d: %w", vaults[i].ID, err)
			}
			vaults[i].Value = decryptedValue
		}
	}

	return vaults, nil
}

// Update updates a vault
func (v *Vault) Update(params *UpdateVaultParams) error {
	updates := map[string]interface{}{}

	if params.Name != nil {
		updates["name"] = *params.Name
	}

	if params.Value != nil {
		// Encrypt the new value before storing
		encryptedValue, err := encryption.Encrypt(*params.Value)
		if err != nil {
			return fmt.Errorf("failed to encrypt value: %w", err)
		}
		updates["value"] = encryptedValue
	}

	if params.Description != nil {
		updates["description"] = *params.Description
	}

	if params.Category != nil {
		updates["category"] = *params.Category
	}

	// Always update the updated_at timestamp
	updates["updated_at"] = time.Now()

	err := DB.Model(v).Updates(updates).Error
	if err != nil {
		return err
	}

	// Reload the vault to get the updated data
	err = DB.Where("id = ?", v.ID).First(v).Error
	if err != nil {
		return err
	}

	// Decrypt the value for the response
	decryptedValue, err := encryption.Decrypt(v.Value)
	if err != nil {
		return fmt.Errorf("failed to decrypt value: %w", err)
	}
	v.Value = decryptedValue

	return nil
}

// Delete deletes a vault
func (v *Vault) Delete() error {
	err := DB.Delete(v).Error
	if err != nil {
		return err
	}
	return nil
}

// CheckVaultOwnership verifies if a vault with the given ID belongs to the specified user
func CheckVaultOwnership(vaultID uint, userID uint) error {
	var count int64
	err := DB.Model(&Vault{}).Where("id = ? AND user_id = ?", vaultID, userID).Count(&count).Error
	if err != nil {
		return err
	}

	if count == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
