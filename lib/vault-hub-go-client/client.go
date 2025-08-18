package vhub

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "time"
)

// Client is a lightweight Vault Hub API client focused on CLI use-cases.
// It supports listing vaults and retrieving single vaults by unique ID or name.
type Client struct {
    BaseURL    string
    APIKey     string
    HTTPClient *http.Client
}

// NewClient constructs a new Client. If httpClient is nil, http.DefaultClient is used.
func NewClient(baseURL, apiKey string) *Client {
    return &Client{
        BaseURL:    baseURL,
        APIKey:     apiKey,
        HTTPClient: http.DefaultClient,
    }
}

// VaultLite mirrors the server's lightweight vault representation.
// Only the fields needed by the CLI are included.
// Timestamps are parsed as RFC3339.
//
// NOTE: Keep in sync with server API surface.
type VaultLite struct {
    Name        string     `json:"name"`
    UniqueId    string     `json:"uniqueId"`
    Category    *string    `json:"category,omitempty"`
    Description *string    `json:"description,omitempty"`
    UpdatedAt   *time.Time `json:"updatedAt,omitempty"`
}

// Vault represents the full vault resource returned by the API.
// It embeds VaultLite and adds the encrypted value and creation metadata.
type Vault struct {
    VaultLite
    Value     string     `json:"value"`
    CreatedAt *time.Time `json:"createdAt,omitempty"`
    UserId    *int64     `json:"userId,omitempty"`
}

// ListVaults returns all vaults accessible with the API key.
func (c *Client) ListVaults(ctx context.Context) ([]VaultLite, error) {
    var vaults []VaultLite
    if err := c.doRequest(ctx, http.MethodGet, "/api/cli/vaults", nil, &vaults); err != nil {
        return nil, err
    }
    return vaults, nil
}

// GetVault fetches a vault by its unique ID.
func (c *Client) GetVault(ctx context.Context, uniqueId string) (*Vault, error) {
    var v Vault
    path := fmt.Sprintf("/api/cli/vault/%s", uniqueId)
    if err := c.doRequest(ctx, http.MethodGet, path, nil, &v); err != nil {
        return nil, err
    }
    return &v, nil
}

// GetVaultByName fetches a vault by its name.
func (c *Client) GetVaultByName(ctx context.Context, name string) (*Vault, error) {
    var v Vault
    path := fmt.Sprintf("/api/cli/vault/name/%s", name)
    if err := c.doRequest(ctx, http.MethodGet, path, nil, &v); err != nil {
        return nil, err
    }
    return &v, nil
}

// doRequest performs the HTTP request and decodes the JSON response into dest if provided.
func (c *Client) doRequest(ctx context.Context, method, path string, body any, dest any) error {
    req, err := http.NewRequestWithContext(ctx, method, c.BaseURL+path, nil)
    if err != nil {
        return err
    }
    req.Header.Set("Accept", "application/json")
    if c.APIKey != "" {
        req.Header.Set("X-API-Key", c.APIKey)
    }

    resp, err := c.HTTPClient.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode >= 300 {
        return fmt.Errorf("request failed: %s", resp.Status)
    }

    if dest != nil {
        decoder := json.NewDecoder(resp.Body)
        if err := decoder.Decode(dest); err != nil {
            return fmt.Errorf("failed to decode response: %w", err)
        }
    }

    return nil
}