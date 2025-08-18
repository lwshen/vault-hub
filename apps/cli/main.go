package main

import (
	"context"
	"fmt"
	"log"
	"os"

	openapi "github.com/lwshen/vault-hub-go-client"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "vault-hub",
	Short: "VaultHub CLI - Secure environment variable and API key management",
	Long: `VaultHub CLI is a command-line interface for managing your secure
environment variables and API keys stored in VaultHub.

This CLI allows you to list and retrieve vaults from your VaultHub instance.`,
}

var (
	apiKey  string
	baseURL string
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func newClientAndAuthCtx() (*openapi.APIClient, context.Context, error) {
	if baseURL == "" {
		return nil, nil, fmt.Errorf("base URL is required; set --base-url or VAULT_HUB_BASE_URL")
	}
	if apiKey == "" {
		return nil, nil, fmt.Errorf("API key is required; set --api-key or VAULT_HUB_API_KEY")
	}

	cfg := openapi.NewConfiguration()
	if len(cfg.Servers) == 0 {
		cfg.Servers = openapi.ServerConfigurations{{URL: baseURL}}
	} else {
		cfg.Servers[0].URL = baseURL
	}

	client := openapi.NewAPIClient(cfg)
	ctx := context.WithValue(context.Background(), openapi.ContextAPIKeys, map[string]openapi.APIKey{
		"ApiKeyAuth": {Key: apiKey},
	})
	return client, ctx, nil
}

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List all accessible vaults",
	Long: `List all vaults that you have access to.
This command will display basic information about each vault including
name, unique ID, and description.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, ctx, err := newClientAndAuthCtx()
		if err != nil {
			return err
		}

		vaults, _, err := client.APIKeyAPI.GetVaultsByAPIKey(ctx).Execute()
		if err != nil {
			return fmt.Errorf("failed to list vaults: %w", err)
		}

		if len(vaults) == 0 {
			fmt.Println("No vaults found.")
			return nil
		}

		for _, v := range vaults {
			desc := ""
			if v.GetDescriptionOk() != nil {
				desc = v.GetDescription()
			}
			fmt.Printf("- %s\n  id: %s\n  desc: %s\n", v.GetName(), v.GetUniqueId(), desc)
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
		client, ctx, err := newClientAndAuthCtx()
		if err != nil {
			return err
		}

		var vault *openapi.Vault
		if _, parseErr := uuid.Parse(vaultIdentifier); parseErr == nil {
			vault, _, err = client.APIKeyAPI.GetVaultByAPIKey(ctx, vaultIdentifier).Execute()
			if err != nil {
				return fmt.Errorf("failed to get vault by id: %w", err)
			}
		} else {
			vault, _, err = client.APIKeyAPI.GetVaultByNameAPIKey(ctx, vaultIdentifier).Execute()
			if err != nil {
				return fmt.Errorf("failed to get vault by name: %w", err)
			}
		}

		fmt.Printf("Name: %s\n", vault.GetName())
		fmt.Printf("ID: %s\n", vault.GetUniqueId())
		if vault.GetDescriptionOk() != nil {
			fmt.Printf("Description: %s\n", vault.GetDescription())
		}
		if vault.GetCategoryOk() != nil {
			fmt.Printf("Category: %s\n", vault.GetCategory())
		}
		fmt.Printf("Value: %s\n", vault.GetValue())
		return nil
	},
}

func init() {
	// Global flags with env var defaults
	rootCmd.PersistentFlags().StringVar(&apiKey, "api-key", os.Getenv("VAULT_HUB_API_KEY"), "API key for authentication (env: VAULT_HUB_API_KEY)")
	rootCmd.PersistentFlags().StringVar(&baseURL, "base-url", os.Getenv("VAULT_HUB_BASE_URL"), "Base URL of VaultHub server, e.g. https://example.com (env: VAULT_HUB_BASE_URL)")

	// Add subcommands to root
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(getCmd)

	if os.Getenv("VAULT_HUB_CLI_DEBUG") == "1" {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
	}
}