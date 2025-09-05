package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	openapi "github.com/lwshen/vault-hub-go-client"
	"github.com/spf13/cobra"
)

// CommandContext holds dependencies for commands
type CommandContext struct {
	GetClient         func() *openapi.APIClient
	DebugLog          func(string, ...any)
	MustGetStringFlag func(*cobra.Command, string) string
}

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

func runListCommand(cmd *cobra.Command, _ []string, ctx *CommandContext) {
	ctx.DebugLog("Executing list command")

	apiCtx := context.Background()
	ctx.DebugLog("Making API request to get vaults by API key")
	vaults, _, err := ctx.GetClient().CliAPI.GetVaultsByAPIKey(apiCtx).Execute()
	if err != nil {
		ctx.DebugLog("API request failed: %v", err)
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	ctx.DebugLog("API request successful, received %d vault(s)", len(vaults))

	jsonOutput, _ := cmd.Flags().GetBool("json")
	ctx.DebugLog("JSON output mode: %v", jsonOutput)

	if jsonOutput {
		ctx.DebugLog("Marshaling vaults to JSON")
		output, err := json.MarshalIndent(vaults, "", "  ")
		if err != nil {
			ctx.DebugLog("JSON marshaling failed: %v", err)
			fmt.Fprintf(os.Stderr, "Error marshaling JSON: %v\n", err)
			os.Exit(1)
		}
		ctx.DebugLog("JSON marshaling successful")
		fmt.Println(string(output))
		return
	}

	if len(vaults) == 0 {
		ctx.DebugLog("No vaults found, displaying empty message")
		fmt.Println("No vaults found.")
		return
	}

	ctx.DebugLog("Displaying %d vault(s) in formatted output", len(vaults))
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
	ctx.DebugLog("List command completed successfully")
}
