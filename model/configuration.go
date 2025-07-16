package model

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/lwshen/vault-hub/internal/encryption"
	"gorm.io/gorm"
)

type Configuration struct {
	gorm.Model
	UniqueID    string `gorm:"size:255;not null;unique"` // Unique identifier for the config
	UserID      uint   `gorm:"index;not null"`           // User who owns this configuration
	Name        string `gorm:"size:255"`                 // Human-readable name
	Value       string `gorm:"type:text;not null"`       // Encrypted value
	Description string `gorm:"size:500"`                 // Human-readable description
	Category    string `gorm:"size:100;index"`           // Category/type of config
}

// CreateConfigurationParams defines parameters for creating a new configuration
type CreateConfigurationParams struct {
	UniqueID    string
	UserID      uint
	Name        string
	Value       string
	Description string
	Category    string
}

// UpdateConfigurationParams defines parameters for updating a configuration
type UpdateConfigurationParams struct {
	Name        *string
	Value       *string
	Description *string
	Category    *string
}

// Validate validates the create configuration parameters
func (params *CreateConfigurationParams) Validate() map[string]string {
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

// Validate validates the update configuration parameters
func (params *UpdateConfigurationParams) Validate() map[string]string {
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

// Create creates a new configuration
func (params *CreateConfigurationParams) Create() (*Configuration, error) {
	// Check if unique_id already exists for this user
	var existing Configuration
	err := DB.Where("unique_id = ? AND user_id = ?", params.UniqueID, params.UserID).First(&existing).Error
	if err == nil {
		return nil, fmt.Errorf("configuration with unique_id '%s' already exists", params.UniqueID)
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Encrypt the value before storing
	encryptedValue, err := encryption.Encrypt(params.Value)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt value: %w", err)
	}

	config := Configuration{
		UniqueID:    params.UniqueID,
		UserID:      params.UserID,
		Name:        params.Name,
		Value:       encryptedValue,
		Description: params.Description,
		Category:    params.Category,
	}

	err = DB.Create(&config).Error
	if err != nil {
		return nil, err
	}

	// Decrypt the value for the response
	config.Value, err = encryption.Decrypt(config.Value)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt value for response: %w", err)
	}

	return &config, nil
}

// GetByID retrieves a configuration by ID for a specific user
func (c *Configuration) GetByID(id uint, userID uint) error {
	err := DB.Where("id = ? AND user_id = ?", id, userID).First(c).Error
	if err != nil {
		return err
	}

	// Decrypt the value
	decryptedValue, err := encryption.Decrypt(c.Value)
	if err != nil {
		return fmt.Errorf("failed to decrypt value: %w", err)
	}
	c.Value = decryptedValue

	return nil
}

// GetByUniqueID retrieves a configuration by unique_id for a specific user
func (c *Configuration) GetByUniqueID(uniqueID string, userID uint) error {
	err := DB.Where("unique_id = ? AND user_id = ?", uniqueID, userID).First(c).Error
	if err != nil {
		return err
	}

	// Decrypt the value
	decryptedValue, err := encryption.Decrypt(c.Value)
	if err != nil {
		return fmt.Errorf("failed to decrypt value: %w", err)
	}
	c.Value = decryptedValue

	return nil
}

// GetAllByUser retrieves all configurations for a user, optionally filtered by category
func GetConfigurationsByUser(userID uint, category string) ([]Configuration, error) {
	var configs []Configuration
	query := DB.Where("user_id = ?", userID)

	if category != "" {
		query = query.Where("category = ?", category)
	}

	err := query.Order("created_at DESC").Find(&configs).Error
	if err != nil {
		return nil, err
	}

	// Decrypt all values
	for i := range configs {
		decryptedValue, err := encryption.Decrypt(configs[i].Value)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt value for config %d: %w", configs[i].ID, err)
		}
		configs[i].Value = decryptedValue
	}

	return configs, nil
}

// Update updates a configuration
func (c *Configuration) Update(params *UpdateConfigurationParams) error {
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

	err := DB.Model(c).Updates(updates).Error
	if err != nil {
		return err
	}

	// Reload the configuration to get the updated data
	err = DB.Where("id = ?", c.ID).First(c).Error
	if err != nil {
		return err
	}

	// Decrypt the value for the response
	decryptedValue, err := encryption.Decrypt(c.Value)
	if err != nil {
		return fmt.Errorf("failed to decrypt value: %w", err)
	}
	c.Value = decryptedValue

	return nil
}

// Delete deletes a configuration
func (c *Configuration) Delete() error {
	err := DB.Delete(c).Error
	if err != nil {
		return err
	}
	return nil
}
