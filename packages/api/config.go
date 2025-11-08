package api

import (
	"github.com/lwshen/vault-hub/internal/config"
)

// PublicConfig exposes non-sensitive flags for clients.
func PublicConfig() ConfigResponse {
	return ConfigResponse{
		OidcEnabled:  config.OidcEnabled,
		EmailEnabled: config.EmailEnabled,
	}
}

// Additional HTTP wiring occurs in the Echo router; this package remains framework-agnostic.
