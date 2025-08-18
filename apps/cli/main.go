package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/lwshen/vault-hub-go-client"
)

var (
	apiKey  string
	baseURL string
	timeout time.Duration
)

var rootCmd = &cobra.Command{
	Use:   "vault-hub",
	Short: "VaultHub CLI - Secure environment variable and API key management",
	Long: `VaultHub CLI is a command-line interface for managing your secure
environment variables and API keys stored in VaultHub.

This CLI allows you to list and retrieve vaults from your VaultHub instance.

Configuration:
  You can set the API key and base URL using:
  - Command line flags: --api-key and --base-url
  - Environment variables: VAULT_HUB_API_KEY and VAULT_HUB_BASE_URL
  - Configuration file: ~/.vault-hub/config.yaml

Examples:
  vault-hub list
  vault-hub get my-api-keys
  vault-hub get abc123-def456-ghi789`,
}

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List all accessible vaults",
	Long: `List all vaults that you have access to.
This command will display basic information about each vault including
name, unique ID, and description.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		vaultClient, err := createClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		vaults, err := vaultClient.ListVaults(ctx)
		if err != nil {
			return fmt.Errorf("failed to list vaults: %w", err)
		}

		if len(vaults) == 0 {
			fmt.Println("No vaults found.")
			return nil
		}

		fmt.Printf("Found %d vault(s):\n\n", len(vaults))
		for _, vault := range vaults {
			fmt.Printf("Name: %s\n", vault.Name)
			fmt.Printf("ID: %s\n", vault.UniqueID)
			if vault.Description != "" {
				fmt.Printf("Description: %s\n", vault.Description)
			}
			fmt.Printf("Created: %s\n", vault.CreatedAt.Format("2006-01-02 15:04:05"))
			fmt.Printf("Updated: %s\n", vault.UpdatedAt.Format("2006-01-02 15:04:05"))
			fmt.Println("---")
		}

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
		vaultClient, err := createClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		var vault *client.Vault
		var err2 error

		// Try to get by unique ID first (UUID format)
		if isUUID(vaultIdentifier) {
			vault, err2 = vaultClient.GetVault(ctx, vaultIdentifier)
		} else {
			// Try by name
			vault, err2 = vaultClient.GetVaultByName(ctx, vaultIdentifier)
		}

		if err2 != nil {
			return fmt.Errorf("failed to get vault: %w", err2)
		}

		fmt.Printf("Vault: %s\n", vault.Name)
		fmt.Printf("ID: %s\n", vault.UniqueID)
		if vault.Description != "" {
			fmt.Printf("Description: %s\n", vault.Description)
		}
		fmt.Printf("Created: %s\n", vault.CreatedAt.Format("2006-01-02 15:04:05"))
		fmt.Printf("Updated: %s\n", vault.UpdatedAt.Format("2006-01-02 15:04:05"))
		fmt.Printf("Value: %s\n", vault.Value)

		return nil
	},
}

var healthCmd = &cobra.Command{
	Use:   "health",
	Short: "Check the health of the VaultHub server",
	Long: `Check if the VaultHub server is running and healthy.
This command will attempt to connect to the server and verify its status.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		vaultClient, err := createClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		if err := vaultClient.Health(ctx); err != nil {
			return fmt.Errorf("server health check failed: %w", err)
		}

		fmt.Println("✅ VaultHub server is healthy")
		return nil
	},
}

func init() {
	// Initialize viper for configuration
	initConfig()

	// Add global flags
	rootCmd.PersistentFlags().StringVar(&apiKey, "api-key", "", "API key for authentication")
	rootCmd.PersistentFlags().StringVar(&baseURL, "base-url", "", "Base URL of VaultHub server")
	rootCmd.PersistentFlags().DurationVar(&timeout, "timeout", 30*time.Second, "Request timeout")

	// Bind flags to viper
	viper.BindPFlag("api_key", rootCmd.PersistentFlags().Lookup("api-key"))
	viper.BindPFlag("base_url", rootCmd.PersistentFlags().Lookup("base-url"))
	viper.BindPFlag("timeout", rootCmd.PersistentFlags().Lookup("timeout"))

	// Add subcommands to root
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(getCmd)
	rootCmd.AddCommand(healthCmd)
}

func initConfig() {
	// Set default values
	viper.SetDefault("base_url", "http://localhost:8080")
	viper.SetDefault("timeout", "30s")

	// Read from environment variables
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.SetEnvPrefix("VAULT_HUB")
	viper.AutomaticEnv()

	// Read from config file
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.vault-hub")
	viper.AddConfigPath(".")

	// Ignore config file errors
	viper.ReadInConfig()
}

func createClient() (*client.Client, error) {
	// Get configuration values
	configAPIKey := viper.GetString("api_key")
	configBaseURL := viper.GetString("base_url")
	configTimeout := viper.GetDuration("timeout")

	// Validate required configuration
	if configAPIKey == "" {
		return nil, fmt.Errorf("API key is required. Set it using --api-key flag, VAULT_HUB_API_KEY environment variable, or config file")
	}

	if configBaseURL == "" {
		return nil, fmt.Errorf("Base URL is required. Set it using --base-url flag, VAULT_HUB_BASE_URL environment variable, or config file")
	}

	// Create client with options
	return client.NewClient(
		configBaseURL,
		configAPIKey,
		client.WithTimeout(configTimeout),
	)
}

func isUUID(str string) bool {
	// Simple UUID validation - check if it contains hyphens and is the right length
	parts := strings.Split(str, "-")
	return len(parts) == 5 && len(str) == 36
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
