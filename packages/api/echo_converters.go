package api

import (
	"github.com/lwshen/vault-hub/model"
	"github.com/lwshen/vault-hub/packages/api/generated_models"
)

// convertToGeneratedVault converts a model.Vault to a generated.Vault
func convertToGeneratedVault(vault *model.Vault) generated_models.Vault {
	// #nosec G115
	userID := int64(vault.UserID)
	return generated_models.Vault{
		UniqueId:    vault.UniqueID,
		UserId:      userID,
		Name:        vault.Name,
		Value:       vault.Value,
		Description: vault.Description,
		Category:    vault.Category,
		CreatedAt:   vault.CreatedAt,
		UpdatedAt:   vault.UpdatedAt,
	}
}

// convertToGeneratedVaultLite converts a model.Vault to a generated.VaultLite
func convertToGeneratedVaultLite(vault *model.Vault) generated_models.VaultLite {
	return generated_models.VaultLite{
		UniqueId:    vault.UniqueID,
		Name:        vault.Name,
		Description: vault.Description,
		Category:    vault.Category,
		UpdatedAt:   vault.UpdatedAt,
	}
}

// convertToGeneratedAPIKey converts a model.APIKey to a generated.VaultApiKey (without vaults)
func convertToGeneratedAPIKey(key *model.APIKey) generated_models.VaultApiKey {
	// #nosec G115
	id := int64(key.ID)

	result := generated_models.VaultApiKey{
		Id:        id,
		Name:      key.Name,
		CreatedAt: key.CreatedAt,
		IsActive:  !key.DeletedAt.Valid,
	}

	if key.ExpiresAt != nil {
		result.ExpiresAt = *key.ExpiresAt
	}

	if key.LastUsedAt != nil {
		result.LastUsedAt = *key.LastUsedAt
	}

	result.UpdatedAt = key.UpdatedAt

	return result
}

// convertToGeneratedAPIKeyWithVaults converts a model.APIKey to a generated.VaultApiKey with vaults
func convertToGeneratedAPIKeyWithVaults(apiKey *model.APIKey) (*generated_models.VaultApiKey, error) {
	// Get accessible vaults for this API key
	vaults, err := apiKey.GetAccessibleVaults()
	if err != nil {
		return nil, err
	}

	// Convert vaults to VaultLite
	apiVaults := make([]generated_models.VaultLite, 0)
	for _, vault := range vaults {
		apiVaults = append(apiVaults, convertToGeneratedVaultLite(&vault))
	}

	// #nosec G115
	id := int64(apiKey.ID)
	result := generated_models.VaultApiKey{
		Id:        id,
		Name:      apiKey.Name,
		Vaults:    apiVaults,
		IsActive:  !apiKey.DeletedAt.Valid,
		CreatedAt: apiKey.CreatedAt,
	}

	if apiKey.ExpiresAt != nil {
		result.ExpiresAt = *apiKey.ExpiresAt
	}

	if apiKey.LastUsedAt != nil {
		result.LastUsedAt = *apiKey.LastUsedAt
	}

	result.UpdatedAt = apiKey.UpdatedAt

	return &result, nil
}

// convertToGeneratedAuditLog converts a model.AuditLog to generated.AuditLog
func convertToGeneratedAuditLog(log *model.AuditLog) generated_models.AuditLog {
	result := generated_models.AuditLog{
		Action:    string(log.Action),
		CreatedAt: log.CreatedAt,
		Source:    string(log.Source),
		IpAddress: log.IPAddress,
		UserAgent: log.UserAgent,
	}

	// Convert related entities if present
	if log.Vault != nil {
		vaultLite := convertToGeneratedVaultLite(log.Vault)
		result.Vault = vaultLite
	}

	if log.APIKey != nil {
		apiKey := convertToGeneratedAPIKey(log.APIKey)
		result.ApiKey = apiKey
	}

	return result
}
