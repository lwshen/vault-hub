package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// Client represents a VaultHub API client
type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

// ClientOption is a function that configures a Client
type ClientOption func(*Client)

// WithHTTPClient sets a custom HTTP client
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// WithTimeout sets the HTTP client timeout
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		if c.httpClient == nil {
			c.httpClient = &http.Client{}
		}
		c.httpClient.Timeout = timeout
	}
}

// NewClient creates a new VaultHub client
func NewClient(baseURL, apiKey string, opts ...ClientOption) (*Client, error) {
	// Validate base URL
	if baseURL == "" {
		return nil, fmt.Errorf("base URL cannot be empty")
	}
	
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %w", err)
	}
	if parsedURL.Scheme == "" {
		parsedURL.Scheme = "https"
	}

	client := &Client{
		baseURL: parsedURL.String(),
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}

	// Apply options
	for _, opt := range opts {
		opt(client)
	}

	return client, nil
}

// VaultLite represents a lightweight vault object
type VaultLite struct {
	ID          int64     `json:"id"`
	UniqueID    string    `json:"uniqueId"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// Vault represents a full vault object with encrypted value
type Vault struct {
	VaultLite
	Value string `json:"value"`
}

// ListVaults retrieves all accessible vaults for the API key
func (c *Client) ListVaults(ctx context.Context) ([]VaultLite, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/api/cli/vaults", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var vaults []VaultLite
	if err := json.NewDecoder(resp.Body).Decode(&vaults); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return vaults, nil
}

// GetVault retrieves a vault by its unique ID
func (c *Client) GetVault(ctx context.Context, uniqueID string) (*Vault, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/api/cli/vault/"+uniqueID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("vault not found")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var vault Vault
	if err := json.NewDecoder(resp.Body).Decode(&vault); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &vault, nil
}

// GetVaultByName retrieves a vault by its name
func (c *Client) GetVaultByName(ctx context.Context, name string) (*Vault, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/api/cli/vault/name/"+name, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("vault not found")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var vault Vault
	if err := json.NewDecoder(resp.Body).Decode(&vault); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &vault, nil
}

// Health checks the health of the VaultHub server
func (c *Client) Health(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/health", nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}