package model

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql/driver"
	"encoding/hex"
	"encoding/json"
	"errors"
	"log/slog"
	"strings"
	"time"

	"gorm.io/gorm"
)

// VaultIDs represents a custom type for storing vault IDs as JSON
type VaultIDs []uint

// Value implements the driver.Valuer interface for storing as JSON in database
func (v VaultIDs) Value() (driver.Value, error) {
	if v == nil {
		return nil, nil
	}
	return json.Marshal(v)
}

// Scan implements the sql.Scanner interface for reading JSON from database
func (v *VaultIDs) Scan(value interface{}) error {
	if value == nil {
		*v = nil
		return nil
	}

	switch s := value.(type) {
	case []byte:
		return json.Unmarshal(s, v)
	case string:
		return json.Unmarshal([]byte(s), v)
	default:
		return errors.New("cannot scan VaultIDs from this type")
	}
}

// APIKey represents an API key for accessing vaults
type APIKey struct {
	gorm.Model
	UserID     uint       `gorm:"not null;index"`          // User who owns this API key
	Name       string     `gorm:"size:255;not null"`       // Human-readable name for the API key
	KeyHash    string     `gorm:"size:64;not null;unique"` // SHA-256 hash of the API key
	VaultIDs   VaultIDs   `gorm:"type:json"`               // JSON array of vault IDs (null = all user's vaults)
	ExpiresAt  *time.Time `gorm:"index"`                   // Optional expiration date
	LastUsedAt *time.Time // Track when it was last used
	IsActive   bool       `gorm:"default:true;index"` // Enable/disable the key
}

// CreateAPIKeyParams defines parameters for creating a new API key
type CreateAPIKeyParams struct {
	UserID    uint
	Name      string
	VaultIDs  []uint // Empty slice or nil means all user's vaults
	ExpiresAt *time.Time
}

// Validate validates the create API key parameters
func (params *CreateAPIKeyParams) Validate() map[string]string {
	errors := map[string]string{}

	if params.UserID == 0 {
		errors["user_id"] = "user ID is required"
	}

	name := strings.TrimSpace(params.Name)
	if name == "" {
		errors["name"] = "name is required"
	}

	if len(params.Name) > 255 {
		errors["name"] = "name must be less than 255 characters"
	}

	// Allow duplicate names – uniqueness check removed as per new requirements

	if params.ExpiresAt != nil && params.ExpiresAt.Before(time.Now()) {
		errors["expires_at"] = "expiration date must be in the future"
	}

	return errors
}

// GenerateAPIKey creates a new cryptographically secure API key
func GenerateAPIKey() (string, error) {
	// Generate 32 bytes of random data
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	// Convert to hex string with a prefix
	return "vhub_" + hex.EncodeToString(bytes), nil
}

// HashAPIKey creates a SHA-256 hash of the API key
func HashAPIKey(key string) string {
	hash := sha256.Sum256([]byte(key))
	return hex.EncodeToString(hash[:])
}

// Create creates a new API key and returns it with the plain text key
func (params *CreateAPIKeyParams) Create() (*APIKey, string, error) {
	// Generate the API key
	plainKey, err := GenerateAPIKey()
	if err != nil {
		return nil, "", err
	}

	// Hash the key for storage
	keyHash := HashAPIKey(plainKey)

	// Convert vault IDs
	var vaultIDs VaultIDs
	if len(params.VaultIDs) > 0 {
		vaultIDs = VaultIDs(params.VaultIDs)
	}

	apiKey := APIKey{
		UserID:    params.UserID,
		Name:      strings.TrimSpace(params.Name),
		KeyHash:   keyHash,
		VaultIDs:  vaultIDs,
		ExpiresAt: params.ExpiresAt,
		IsActive:  true,
	}

	err = DB.Create(&apiKey).Error
	if err != nil {
		return nil, "", err
	}

	return &apiKey, plainKey, nil
}

// GetByKeyHash finds an API key by its hash
func GetAPIKeyByHash(keyHash string) (*APIKey, error) {
	var apiKey APIKey
	err := DB.Where("key_hash = ? AND is_active = ? AND (expires_at IS NULL OR expires_at > ?)",
		keyHash, true, time.Now()).
		Preload("User").
		First(&apiKey).Error

	if err != nil {
		return nil, err
	}

	return &apiKey, nil
}

// ValidateAPIKey validates an API key and returns the associated user
func ValidateAPIKey(key string) (*APIKey, error) {
	if !strings.HasPrefix(key, "vhub_") {
		return nil, errors.New("invalid API key format")
	}

	keyHash := HashAPIKey(key)
	apiKey, err := GetAPIKeyByHash(keyHash)
	if err != nil {
		return nil, err
	}

	// Update last used timestamp
	now := time.Now()
	apiKey.LastUsedAt = &now
	if err := DB.Save(&apiKey).Error; err != nil {
		// Log the error but don't fail the validation - usage tracking is not critical
		// for API key validation functionality
		slog.Error("Failed to update API key last used timestamp",
			"api_key_id", apiKey.ID,
			"error", err)
	}

	return apiKey, nil
}

// HasVaultAccess checks if the API key has access to a specific vault
func (k *APIKey) HasVaultAccess(vaultID uint) bool {
	// First, verify that the vault belongs to the user who owns this API key
	var vault Vault
	err := DB.Where("id = ? AND user_id = ?", vaultID, k.UserID).First(&vault).Error
	if err != nil {
		// Vault doesn't exist or doesn't belong to the user
		return false
	}

	// If VaultIDs is empty, it means access to all vaults belonging to the user
	if len(k.VaultIDs) == 0 {
		return true
	}

	// Check if the vault ID is in the allowed list
	for _, id := range k.VaultIDs {
		if id == vaultID {
			return true
		}
	}

	return false
}

// GetAccessibleVaults returns the vaults this API key can access
func (k *APIKey) GetAccessibleVaults() ([]Vault, error) {
	var vaults []Vault

	query := DB.Where("user_id = ?", k.UserID)

	// If VaultIDs is specified, filter by those IDs
	if len(k.VaultIDs) > 0 {
		query = query.Where("id IN ?", []uint(k.VaultIDs))
	}

	err := query.Find(&vaults).Error
	return vaults, err
}

// UpdateAPIKeyParams defines parameters for updating an API key
type UpdateAPIKeyParams struct {
	Name      *string
	VaultIDs  *[]uint
	ExpiresAt *time.Time
	IsActive  *bool
}

// Validate validates the update API key parameters
func (params *UpdateAPIKeyParams) Validate() map[string]string {
	errors := map[string]string{}

	if params.Name != nil {
		name := strings.TrimSpace(*params.Name)
		if name == "" {
			errors["name"] = "name cannot be empty"
		}
		if len(name) > 255 {
			errors["name"] = "name must be less than 255 characters"
		}
	}

	if params.ExpiresAt != nil && params.ExpiresAt.Before(time.Now()) {
		errors["expires_at"] = "expiration date must be in the future"
	}

	return errors
}

// ValidateForUpdate validates the update API key parameters for a specific API key
func (params *UpdateAPIKeyParams) ValidateForUpdate(userID uint, currentAPIKeyID uint) map[string]string {
	errors := params.Validate()

	if params.Name != nil {
		name := strings.TrimSpace(*params.Name)
		if name != "" {
			// Check name uniqueness for the user, excluding the current API key
			var count int64
			err := DB.Model(&APIKey{}).Where("user_id = ? AND name = ? AND id != ?", userID, name, currentAPIKeyID).Count(&count).Error
			if err != nil {
				errors["name"] = "failed to validate name uniqueness"
			} else if count > 0 {
				errors["name"] = "API key name already exists"
			}
		}
	}

	return errors
}

// Update updates an API key
func (k *APIKey) Update(params UpdateAPIKeyParams) error {
	if params.Name != nil {
		k.Name = strings.TrimSpace(*params.Name)
	}

	if params.VaultIDs != nil {
		if len(*params.VaultIDs) > 0 {
			k.VaultIDs = VaultIDs(*params.VaultIDs)
		} else {
			k.VaultIDs = nil
		}
	}

	if params.ExpiresAt != nil {
		k.ExpiresAt = params.ExpiresAt
	}

	if params.IsActive != nil {
		k.IsActive = *params.IsActive
	}

	return DB.Save(k).Error
}

// GetUserAPIKeys returns all API keys for a user
func GetUserAPIKeys(userID uint) ([]APIKey, error) {
	var apiKeys []APIKey
	err := DB.Where("user_id = ?", userID).Find(&apiKeys).Error
	return apiKeys, err
}

// GetUserAPIKeysWithPagination returns API keys for a user with pagination
func GetUserAPIKeysWithPagination(userID uint, pageSize, pageIndex int) ([]APIKey, int64, error) {
	var apiKeys []APIKey
	var totalCount int64

	// Get total count
	err := DB.Model(&APIKey{}).Where("user_id = ?", userID).Count(&totalCount).Error
	if err != nil {
		return nil, 0, err
	}

	// Calculate offset (pageIndex is 1-based)
	offset := (pageIndex - 1) * pageSize

	// Get paginated results
	err = DB.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&apiKeys).Error

	return apiKeys, totalCount, err
}

// Delete soft deletes an API key
func (k *APIKey) Delete() error {
	return DB.Delete(k).Error
}
