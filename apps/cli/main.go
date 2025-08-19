package main

import (
	"context"
	"encoding/json"
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
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		cfg := openapi.NewConfiguration()
		cfg.Debug = debug
		cfg.Servers = openapi.ServerConfigurations{
			{
				URL: baseURL,
			},
		}
		client = openapi.NewAPIClient(cfg)
		client.GetConfig().DefaultHeader["Authorization"] = "Bearer " + apiKey
	},
}

func init() {
	// Add global flags
	rootCmd.PersistentFlags().StringVar(&apiKey, "api-key", "", "API key for authentication")
	rootCmd.PersistentFlags().StringVar(&baseURL, "base-url", "", "Base URL of VaultHub server")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Enable debug mode")

	// Add subcommands to root
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(getCmd)

	listCmd.Flags().BoolP("json", "j", false, "Output in JSON format")
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
		ctx := context.Background()
		vaults, _, err := client.CliAPI.GetVaultsByAPIKey(ctx).Execute()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		jsonOutput, _ := cmd.Flags().GetBool("json")
		if jsonOutput {
			output, err := json.MarshalIndent(vaults, "", "  ")
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error marshaling JSON: %v\n", err)
				os.Exit(1)
			}
			fmt.Println(string(output))
			return
		}

		if len(vaults) == 0 {
			fmt.Println("No vaults found.")
			return
		}

		fmt.Printf("Found %d vault(s):\n\n", len(vaults))
		for i, vault := range vaults {
			fmt.Printf("  %d. 📦 %s\n", i+1, vault.GetName())
			fmt.Printf("     ID: %s\n", vault.GetUniqueId())
			if vault.Category != nil && *vault.Category != "" {
				fmt.Printf("     Category: %s\n", *vault.Category)
			}
			if vault.Description != nil && *vault.Description != "" {
				fmt.Printf("     Description: %s\n", *vault.Description)
			}
			if i < len(vaults)-1 {
				fmt.Println()
			}
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
