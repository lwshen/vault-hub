package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	openapi "github.com/lwshen/vault-hub-go-client"
	"github.com/spf13/cobra"
)

// NewListCommand creates the list command
func NewListCommand(ctx *CommandContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List all accessible vaults",
		Long: `List all vaults that you have access to.
This command will display basic information about each vault including
name, unique ID, and description.`,
		Run: func(cmd *cobra.Command, args []string) {
			runListCommand(cmd, args, ctx)
		},
	}

	cmd.Flags().BoolP("json", "j", false, "Output in JSON format")
	return cmd
}

// runListCommand executes the vault listing operation
// It fetches vaults from the API and displays them in either JSON or formatted output
func runListCommand(cmd *cobra.Command, _ []string, ctx *CommandContext) {
	ctx.DebugLog("Executing list command")

	// Fetch vaults from API
	vaults, err := fetchVaultsFromAPI(ctx)
	if err != nil {
		ctx.DebugLog("API request failed: %v", err)
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	ctx.DebugLog("API request successful, received %d vault(s)", len(vaults))

	// Check output format preference
	jsonOutput, _ := cmd.Flags().GetBool("json")
	ctx.DebugLog("JSON output mode: %v", jsonOutput)

	if jsonOutput {
		printJSONOutput(vaults, ctx)
	} else {
		printFormattedOutput(vaults, ctx)
	}

	ctx.DebugLog("List command completed successfully")
}

// fetchVaultsFromAPI retrieves all accessible vaults from the API
func fetchVaultsFromAPI(ctx *CommandContext) ([]openapi.VaultLite, error) {
	apiCtx := context.Background()
	ctx.DebugLog("Making API request to get vaults by API key")
	vaults, _, err := ctx.GetClient().CliAPI.GetVaultsByAPIKey(apiCtx).Execute()
	return vaults, err
}

// printJSONOutput marshals and prints vaults in JSON format
func printJSONOutput(vaults []openapi.VaultLite, ctx *CommandContext) {
	ctx.DebugLog("Marshaling vaults to JSON")
	output, err := json.MarshalIndent(vaults, "", "  ")
	if err != nil {
		ctx.DebugLog("JSON marshaling failed: %v", err)
		fmt.Fprintf(os.Stderr, "Error marshaling JSON: %v\n", err)
		os.Exit(1)
	}
	ctx.DebugLog("JSON marshaling successful")
	fmt.Println(string(output))
}

// printFormattedOutput displays vaults in a human-readable format
func printFormattedOutput(vaults []openapi.VaultLite, ctx *CommandContext) {
	if len(vaults) == 0 {
		ctx.DebugLog("No vaults found, displaying empty message")
		fmt.Println("No vaults found.")
		return
	}

	ctx.DebugLog("Displaying %d vault(s) in formatted output", len(vaults))
	fmt.Printf("Found %d vault(s):\n\n", len(vaults))

	for i, vault := range vaults {
		printVaultInfo(vault, i+1)
		if i < len(vaults)-1 {
			fmt.Println()
		}
	}
}

// printVaultInfo displays information for a single vault
func printVaultInfo(vault openapi.VaultLite, index int) {
	fmt.Printf("  %d. 📦 %s\n", index, vault.GetName())
	fmt.Printf("     ID: %s\n", vault.GetUniqueId())

	if vault.Category != nil && *vault.Category != "" {
		fmt.Printf("     Category: %s\n", *vault.Category)
	}

	if vault.Description != nil && *vault.Description != "" {
		fmt.Printf("     Description: %s\n", *vault.Description)
	}
}
