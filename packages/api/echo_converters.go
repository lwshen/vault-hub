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
		CreatedAt:   &vault.CreatedAt,
		UpdatedAt:   &vault.UpdatedAt,
	}
}

// convertToGeneratedVaultLite converts a model.Vault to a generated.VaultLite
func convertToGeneratedVaultLite(vault *model.Vault) generated_models.VaultLite {
	return generated_models.VaultLite{
		UniqueId:    vault.UniqueID,
		Name:        vault.Name,
		Description: vault.Description,
		Category:    vault.Category,
		UpdatedAt:   &vault.UpdatedAt,
	}
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
		result.ExpiresAt = apiKey.ExpiresAt
	}

	if apiKey.LastUsedAt != nil {
		result.LastUsedAt = apiKey.LastUsedAt
	}

	result.UpdatedAt = &apiKey.UpdatedAt

	return &result, nil
}

// convertToGeneratedAuditLog converts a model.AuditLog to a generated.AuditLog
func convertToGeneratedAuditLog(auditLog *model.AuditLog) generated_models.AuditLog {
	result := generated_models.AuditLog{
		CreatedAt: auditLog.CreatedAt,
		Action:    string(auditLog.Action),
		Source:    string(auditLog.Source),
		IpAddress: auditLog.IPAddress,
		UserAgent: auditLog.UserAgent,
	}

	// Convert vault if present
	if auditLog.Vault != nil {
		vaultLite := convertToGeneratedVaultLite(auditLog.Vault)
		result.Vault = vaultLite
	}

	// Convert API key if present
	if auditLog.APIKey != nil {
		// For audit logs, we don't need the full vault access list
		// Just basic API key info
		// #nosec G115
		id := int64(auditLog.APIKey.ID)
		apiKeyLite := generated_models.VaultApiKey{
			Id:        id,
			Name:      auditLog.APIKey.Name,
			Vaults:    []generated_models.VaultLite{}, // Empty for audit log display
			IsActive:  !auditLog.APIKey.DeletedAt.Valid,
			CreatedAt: auditLog.APIKey.CreatedAt,
			UpdatedAt: &auditLog.APIKey.UpdatedAt,
		}
		if auditLog.APIKey.ExpiresAt != nil {
			apiKeyLite.ExpiresAt = auditLog.APIKey.ExpiresAt
		}
		if auditLog.APIKey.LastUsedAt != nil {
			apiKeyLite.LastUsedAt = auditLog.APIKey.LastUsedAt
		}
		result.ApiKey = apiKeyLite
	}

	return result
}
