package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	vhub "github.com/lwshen/vault-hub-go-client"
	"github.com/spf13/cobra"
)

// global CLI flags
var (
	apiKey  string
	baseURL string
)

var rootCmd = &cobra.Command{
	Use:   "vault-hub",
	Short: "VaultHub CLI - Secure environment variable and API key management",
	Long: `VaultHub CLI is a command-line interface for managing your secure
environment variables and API keys stored in VaultHub.

This CLI allows you to list and retrieve vaults from your VaultHub instance.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Allow environment variables as fallback
		if apiKey == "" {
			apiKey = os.Getenv("VAULTHUB_API_KEY")
		}
		if baseURL == "" {
			baseURL = os.Getenv("VAULTHUB_BASE_URL")
		}

		if apiKey == "" {
			return fmt.Errorf("api-key flag or VAULTHUB_API_KEY env var is required")
		}
		if baseURL == "" {
			return fmt.Errorf("base-url flag or VAULTHUB_BASE_URL env var is required")
		}
		return nil
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// helper to build a configured API client and auth context
func buildClient() (*vhub.Client, context.Context) {
	ctx := context.Background()
	client := vhub.NewClient(strings.TrimRight(baseURL, "/"), apiKey)
	return client, ctx
}

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List all accessible vaults",
	Long: `List all vaults that you have access to.
This command will display basic information about each vault including
name, unique ID, and description.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, ctx := buildClient()

		vaults, err := client.ListVaults(ctx)
		if err != nil {
			return err
		}

		if len(vaults) == 0 {
			fmt.Println("No vaults found")
			return nil
		}

		// Pretty print list as table-style output
		fmt.Printf("%-40s %-40s %-20s\n", "NAME", "UNIQUE ID", "CATEGORY")
		fmt.Println(strings.Repeat("-", 110))
		for _, v := range vaults {
			category := "-"
			if v.Category != nil {
				category = *v.Category
			}
			fmt.Printf("%-40s %-40s %-20s\n", v.Name, v.UniqueId, category)
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
		ident := args[0]

		client, ctx := buildClient()

		var vault *vhub.Vault
		var err error

		if strings.Contains(ident, "-") {
			vault, err = client.GetVault(ctx, ident)
			if err != nil {
				// fallback to name
				vault, err = client.GetVaultByName(ctx, ident)
			}
		} else {
			vault, err = client.GetVaultByName(ctx, ident)
			if err != nil {
				// fallback to uniqueId
				vault, err = client.GetVault(ctx, ident)
			}
		}
		if err != nil {
			return err
		}

		// Pretty print vault as JSON for now
		data, _ := json.MarshalIndent(vault, "", "  ")
		fmt.Println(string(data))
		return nil
	},
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().StringVar(&apiKey, "api-key", "", "API key for authentication (env: VAULTHUB_API_KEY)")
	rootCmd.PersistentFlags().StringVar(&baseURL, "base-url", "", "Base URL of VaultHub server (env: VAULTHUB_BASE_URL)")

	// Add subcommands to root
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(getCmd)
}
