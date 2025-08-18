package openapi

import (
	"context"
)

type service struct{ client *APIClient }

// APIClient holds services
type APIClient struct {
	cfg      *Configuration
	common   service

	APIKeyAPI *APIKeyAPIService
}

// NewAPIClient creates a new API client
func NewAPIClient(cfg *Configuration) *APIClient {
	c := &APIClient{cfg: cfg}
	c.common.client = c
	c.APIKeyAPI = (*APIKeyAPIService)(&c.common)
	return c
}

// Helper to attach API key header; in this shim callers manage context value
func getAPIKeyFromContext(ctx context.Context) (string, bool) {
	if ctx == nil {
		return "", false
	}
	if auth, ok := ctx.Value(ContextAPIKeys).(map[string]APIKey); ok {
		if apiKey, ok := auth["ApiKeyAuth"]; ok {
			if apiKey.Prefix != "" {
				return apiKey.Prefix + " " + apiKey.Key, true
			}
			return apiKey.Key, true
		}
	}
	return "", false
}

