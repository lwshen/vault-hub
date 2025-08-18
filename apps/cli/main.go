package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

/*
VaultHub CLI Implementation

This CLI application demonstrates how to use the github.com/lwshen/vault-hub-go-client package
to interact with a VaultHub server. The implementation includes:

1. Configuration management for API key and base URL (via flags or environment variables)
2. List command to fetch all accessible vaults
3. Get command to retrieve specific vaults by name or unique ID
4. Proper error handling and user-friendly output

IMPORTANT: The github.com/lwshen/vault-hub-go-client package currently has compilation issues
due to conflicting type definitions (APIKey type is defined in both model_api_key.go and 
configuration.go). This implementation uses mock structures to demonstrate the correct usage
pattern.

To use the actual package once it's fixed, replace the mock structures below with:
  import openapi "github.com/lwshen/vault-hub-go-client"

And update the client initialization to use:
  config := openapi.NewConfiguration()
  client := openapi.NewAPIClient(config)

Configuration Options:
- --api-key: API key for authentication (or set VAULT_HUB_API_KEY env var)
- --base-url: Base URL of VaultHub server (or set VAULT_HUB_BASE_URL env var)

Usage Examples:
  vault-hub --api-key="your-key" --base-url="https://your-server.com" list
  vault-hub --api-key="your-key" --base-url="https://your-server.com" get my-vault
  VAULT_HUB_API_KEY="your-key" VAULT_HUB_BASE_URL="https://your-server.com" vault-hub list
*/

// Temporary structs to demonstrate the API structure
// These would be replaced by the actual openapi package once the compilation issues are resolved
type APIClient struct {
	VaultAPI *VaultAPIService
	config   *Configuration
}

type Configuration struct {
	Servers       []ServerConfiguration
	DefaultHeader map[string]string
	HTTPClient    *http.Client
}

type ServerConfiguration struct {
	URL         string
	Description string
}

type VaultAPIService struct {
	client *APIClient
}

type VaultLite struct {
	UniqueId    string  `json:"uniqueId"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
}

type Vault struct {
	UniqueId    string     `json:"uniqueId"`
	Name        string     `json:"name"`
	Value       string     `json:"value"`
	Description *string    `json:"description,omitempty"`
	Category    *string    `json:"category,omitempty"`
	CreatedAt   *time.Time `json:"createdAt,omitempty"`
	UpdatedAt   *time.Time `json:"updatedAt,omitempty"`
}

// Mock API methods - these would be replaced by the actual openapi package methods
func (v *VaultAPIService) GetVaults(ctx context.Context) *GetVaultsRequest {
	return &GetVaultsRequest{service: v}
}

func (v *VaultAPIService) GetVault(ctx context.Context, uniqueId string) *GetVaultRequest {
	return &GetVaultRequest{service: v, uniqueId: uniqueId}
}

type GetVaultsRequest struct {
	service *VaultAPIService
}

func (r *GetVaultsRequest) Execute() ([]VaultLite, *http.Response, error) {
	// This would make the actual HTTP request to the VaultHub API
	// For now, return a mock response to demonstrate the structure
	return nil, nil, fmt.Errorf("API client package has compilation issues - this is a demonstration of the correct structure")
}

type GetVaultRequest struct {
	service  *VaultAPIService
	uniqueId string
}

func (r *GetVaultRequest) Execute() (*Vault, *http.Response, error) {
	// This would make the actual HTTP request to the VaultHub API
	// For now, return a mock response to demonstrate the structure
	return nil, nil, fmt.Errorf("API client package has compilation issues - this is a demonstration of the correct structure")
}

// Global configuration variables
var (
	apiKey  string
	baseURL string
	client  *APIClient
)

var rootCmd = &cobra.Command{
	Use:   "vault-hub",
	Short: "VaultHub CLI - Secure environment variable and API key management",
	Long: `VaultHub CLI is a command-line interface for managing your secure
environment variables and API keys stored in VaultHub.

This CLI allows you to list and retrieve vaults from your VaultHub instance.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return initializeClient()
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// initializeClient sets up the API client with configuration
func initializeClient() error {
	// Check for required configuration
	if apiKey == "" {
		if envKey := os.Getenv("VAULT_HUB_API_KEY"); envKey != "" {
			apiKey = envKey
		} else {
			return fmt.Errorf("API key is required. Use --api-key flag or set VAULT_HUB_API_KEY environment variable")
		}
	}

	if baseURL == "" {
		if envURL := os.Getenv("VAULT_HUB_BASE_URL"); envURL != "" {
			baseURL = envURL
		} else {
			return fmt.Errorf("base URL is required. Use --base-url flag or set VAULT_HUB_BASE_URL environment variable")
		}
	}

	// Create configuration (this would use openapi.NewConfiguration() once the package is fixed)
	config := &Configuration{
		Servers: []ServerConfiguration{
			{
				URL:         strings.TrimSuffix(baseURL, "/"),
				Description: "VaultHub Server",
			},
		},
		DefaultHeader: make(map[string]string),
		HTTPClient:    http.DefaultClient,
	}

	// Set API key authentication
	config.DefaultHeader["Authorization"] = "Bearer " + apiKey

	// Create client (this would use openapi.NewAPIClient(config) once the package is fixed)
	client = &APIClient{
		config: config,
	}
	client.VaultAPI = &VaultAPIService{client: client}

	fmt.Printf("✓ Initialized VaultHub client\n")
	fmt.Printf("  Base URL: %s\n", baseURL)
	fmt.Printf("  API Key:  %s...%s\n", apiKey[:min(4, len(apiKey))], apiKey[max(0, len(apiKey)-4):])

	return nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List all accessible vaults",
	Long: `List all vaults that you have access to.
This command will display basic information about each vault including
name, unique ID, and description.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Fetching vaults...")

		// Call the API to get vaults
		resp, httpResp, err := client.VaultAPI.GetVaults(context.Background()).Execute()
		if err != nil {
			if httpResp != nil {
				return fmt.Errorf("failed to fetch vaults (HTTP %d): %v", httpResp.StatusCode, err)
			}
			return fmt.Errorf("failed to fetch vaults: %v", err)
		}

		if len(resp) == 0 {
			fmt.Println("No vaults found.")
			return nil
		}

		// Display vaults in a table format
		fmt.Printf("%-36s %-20s %s\n", "UNIQUE ID", "NAME", "DESCRIPTION")
		fmt.Printf("%-36s %-20s %s\n", strings.Repeat("-", 36), strings.Repeat("-", 20), strings.Repeat("-", 20))

		for _, vault := range resp {
			description := ""
			if vault.Description != nil {
				description = *vault.Description
			}
			fmt.Printf("%-36s %-20s %s\n", vault.UniqueId, vault.Name, description)
		}

		fmt.Printf("\nTotal: %d vault(s)\n", len(resp))
		return nil
	},
}

var getCmd = &cobra.Command{
	Use:   "get <vault-name-or-id>",
	Short: "Get a specific vault by name or unique ID",
	Long: `Get detailed information about a specific vault including its encrypted value.
You can specify the vault by either its name or unique ID.

Examples:
  vault-hub get my-api-keys
  vault-hub get abc123-def456-ghi789`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		vaultIdentifier := args[0]

		// First, try to get the vault directly by unique ID
		vault, httpResp, err := client.VaultAPI.GetVault(context.Background(), vaultIdentifier).Execute()
		if err == nil {
			displayVault(vault)
			return nil
		}

		// If direct lookup failed, it might be a name. List all vaults and find by name
		if httpResp != nil && httpResp.StatusCode == 404 {
			fmt.Printf("Vault with ID '%s' not found. Searching by name...\n", vaultIdentifier)

			vaults, _, listErr := client.VaultAPI.GetVaults(context.Background()).Execute()
			if listErr != nil {
				return fmt.Errorf("failed to search vaults by name: %v", listErr)
			}

			// Find vault by name (case-insensitive)
			for _, vaultLite := range vaults {
				if strings.EqualFold(vaultLite.Name, vaultIdentifier) {
					// Get full vault details
					fullVault, _, getErr := client.VaultAPI.GetVault(context.Background(), vaultLite.UniqueId).Execute()
					if getErr != nil {
						return fmt.Errorf("failed to get vault details: %v", getErr)
					}
					displayVault(fullVault)
					return nil
				}
			}

			return fmt.Errorf("vault with name '%s' not found", vaultIdentifier)
		}

		// Other error occurred
		if httpResp != nil {
			return fmt.Errorf("failed to get vault (HTTP %d): %v", httpResp.StatusCode, err)
		}
		return fmt.Errorf("failed to get vault: %v", err)
	},
}

// displayVault formats and displays vault information
func displayVault(vault *Vault) {
	fmt.Printf("Vault Details:\n")
	fmt.Printf("  Unique ID:   %s\n", vault.UniqueId)
	fmt.Printf("  Name:        %s\n", vault.Name)

	if vault.Description != nil {
		fmt.Printf("  Description: %s\n", *vault.Description)
	}

	if vault.Category != nil {
		fmt.Printf("  Category:    %s\n", *vault.Category)
	}

	fmt.Printf("  Value:       %s\n", vault.Value)

	if vault.CreatedAt != nil {
		fmt.Printf("  Created:     %s\n", vault.CreatedAt.Format("2006-01-02 15:04:05"))
	}

	if vault.UpdatedAt != nil {
		fmt.Printf("  Updated:     %s\n", vault.UpdatedAt.Format("2006-01-02 15:04:05"))
	}
}

func init() {
	// Add global flags for API key and base URL
	rootCmd.PersistentFlags().StringVar(&apiKey, "api-key", "", "API key for authentication (or set VAULT_HUB_API_KEY env var)")
	rootCmd.PersistentFlags().StringVar(&baseURL, "base-url", "", "Base URL of VaultHub server (or set VAULT_HUB_BASE_URL env var)")

	// Don't mark flags as required since we allow env vars as alternatives
	// The validation is handled in initializeClient()

	// Add subcommands to root
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(getCmd)
}
