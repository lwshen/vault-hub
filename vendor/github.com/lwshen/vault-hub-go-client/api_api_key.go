package openapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// APIKeyAPIService provides API Key scoped endpoints used by the CLI
type APIKeyAPIService service

// GetVaultsByAPIKey returns all vaults accessible by API key
func (a *APIKeyAPIService) GetVaultsByAPIKey(ctx context.Context) ApiGetVaultsByAPIKeyRequest {
	return ApiGetVaultsByAPIKeyRequest{ApiService: a, ctx: ctx}
}

type ApiGetVaultsByAPIKeyRequest struct {
	ctx        context.Context
	ApiService *APIKeyAPIService
}

func (r ApiGetVaultsByAPIKeyRequest) Execute() ([]VaultLite, *http.Response, error) {
	client := r.ApiService.client
	base, _ := client.cfg.ServerURLWithContext(r.ctx, "APIKeyAPIService.GetVaultsByAPIKey")
	req, _ := http.NewRequest(http.MethodGet, base+"/api/cli/vaults", nil)
	if key, ok := getAPIKeyFromContext(r.ctx); ok {
		req.Header.Set("X-API-Key", key)
	}
	resp, err := client.cfg.HTTPClient.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return nil, resp, fmt.Errorf("unexpected status: %s", resp.Status)
	}
	var out []VaultLite
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, resp, err
	}
	return out, resp, nil
}

// GetVaultByAPIKey retrieves a vault by unique ID
func (a *APIKeyAPIService) GetVaultByAPIKey(ctx context.Context, uniqueId string) ApiGetVaultByAPIKeyRequest {
	return ApiGetVaultByAPIKeyRequest{ApiService: a, ctx: ctx, uniqueId: uniqueId}
}

type ApiGetVaultByAPIKeyRequest struct {
	ctx        context.Context
	ApiService *APIKeyAPIService
	uniqueId   string
}

func (r ApiGetVaultByAPIKeyRequest) Execute() (*Vault, *http.Response, error) {
	client := r.ApiService.client
	base, _ := client.cfg.ServerURLWithContext(r.ctx, "APIKeyAPIService.GetVaultByAPIKey")
	req, _ := http.NewRequest(http.MethodGet, base+"/api/cli/vault/"+r.uniqueId, nil)
	if key, ok := getAPIKeyFromContext(r.ctx); ok {
		req.Header.Set("X-API-Key", key)
	}
	resp, err := client.cfg.HTTPClient.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return nil, resp, fmt.Errorf("unexpected status: %s", resp.Status)
	}
	var out Vault
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, resp, err
	}
	return &out, resp, nil
}

// GetVaultByNameAPIKey retrieves a vault by name
func (a *APIKeyAPIService) GetVaultByNameAPIKey(ctx context.Context, name string) ApiGetVaultByNameAPIKeyRequest {
	return ApiGetVaultByNameAPIKeyRequest{ApiService: a, ctx: ctx, name: name}
}

type ApiGetVaultByNameAPIKeyRequest struct {
	ctx        context.Context
	ApiService *APIKeyAPIService
	name       string
}

func (r ApiGetVaultByNameAPIKeyRequest) Execute() (*Vault, *http.Response, error) {
	client := r.ApiService.client
	base, _ := client.cfg.ServerURLWithContext(r.ctx, "APIKeyAPIService.GetVaultByNameAPIKey")
	req, _ := http.NewRequest(http.MethodGet, base+"/api/cli/vault/name/"+r.name, nil)
	if key, ok := getAPIKeyFromContext(r.ctx); ok {
		req.Header.Set("X-API-Key", key)
	}
	resp, err := client.cfg.HTTPClient.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return nil, resp, fmt.Errorf("unexpected status: %s", resp.Status)
	}
	var out Vault
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, resp, err
	}
	return &out, resp, nil
}

