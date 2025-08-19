package main

import (
	"context"
	"fmt"
	"os"

	openapi "github.com/lwshen/vault-hub-go-client"
	"github.com/spf13/cobra"
)

var (
	apiKey  string
	baseURL string
	debug   bool
	client  *openapi.APIClient
)

var rootCmd = &cobra.Command{
	Use:   "vault-hub",
	Short: "VaultHub CLI - Secure environment variable and API key management",
	Long: `VaultHub CLI is a command-line interface for managing your secure
environment variables and API keys stored in VaultHub.

This CLI allows you to list and retrieve vaults from your VaultHub instance.`,
}

func init() {
	// Add global flags
	rootCmd.PersistentFlags().StringVar(&apiKey, "api-key", "", "API key for authentication")
	rootCmd.PersistentFlags().StringVar(&baseURL, "base-url", "", "Base URL of VaultHub server")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Enable debug mode")

	cfg := openapi.NewConfiguration()
	cfg.Debug = debug
	cfg.Servers = openapi.ServerConfigurations{
		{
			URL: baseURL,
		},
	}
	client = openapi.NewAPIClient(cfg)
	client.GetConfig().DefaultHeader["Authorization"] = "Bearer " + apiKey

	// Add subcommands to root
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(getCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List all accessible vaults",
	Long: `List all vaults that you have access to.
This command will display basic information about each vault including
name, unique ID, and description.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Listing vaults...")
		ctx := context.Background()
		vaults, _, err := client.CliAPI.GetVaultsByAPIKey(ctx).Execute()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		for _, vault := range vaults {
			fmt.Printf("Vault: %s\n", vault.Name)
		}
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
	Run: func(cmd *cobra.Command, args []string) {
		vaultIdentifier := args[0]
		fmt.Printf("Getting vault: %s\n", vaultIdentifier)
		// TODO: Implement vault retrieval functionality
		// This should call the /api/cli/vault/{uniqueId} endpoint
		// Need to determine if the identifier is a name or unique ID
		// and handle accordingly
	},
}
