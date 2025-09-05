package commands

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	openapi "github.com/lwshen/vault-hub-go-client"
	"github.com/spf13/cobra"
)

// NewGetCommand creates the get command
func NewGetCommand(ctx *CommandContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <vault-name-or-id>",
		Short: "Get a specific vault by name or unique ID",
		Long: `Get detailed information about a specific vault including its encrypted value.
You can specify the vault by either its name or unique ID.

Examples:
  vault-hub get --name my-api-keys
  vault-hub get --id abc123-def456-ghi789
  vault-hub get --name my-api-keys --output ./secrets.txt
  vault-hub get --name my-api-keys --output ./secrets.txt --exec "source ./secrets.txt && npm start"`,
		Run: func(cmd *cobra.Command, args []string) {
			runGetCommand(cmd, args, ctx)
		},
	}

	cmd.Flags().StringP("name", "n", "", "Vault name")
	cmd.Flags().StringP("id", "i", "", "Vault Unique ID")
	cmd.Flags().StringP("output", "o", "", "Output to file instead of stdout")
	cmd.Flags().StringP("exec", "e", "", "Command to execute if vault has been updated")

	return cmd
}

func runGetCommand(cmd *cobra.Command, _ []string, ctx *CommandContext) {
	ctx.DebugLog("Executing get command")

	name := ctx.MustGetStringFlag(cmd, "name")
	id := ctx.MustGetStringFlag(cmd, "id")
	outputFile := ctx.MustGetStringFlag(cmd, "output")
	followUpCommand := ctx.MustGetStringFlag(cmd, "exec")

	ctx.DebugLog("Parameters - name: '%s', id: '%s', output: '%s', exec: '%s'", name, id, outputFile, followUpCommand)

	if name == "" && id == "" {
		ctx.DebugLog("Validation failed: neither name nor id provided")
		fmt.Fprintf(os.Stderr, "Error: either name or id must be provided\n")
		os.Exit(1)
	}

	apiCtx := context.Background()
	var vault *openapi.Vault
	var err error
	if name != "" {
		ctx.DebugLog("Making API request to get vault by name: %s", name)
		vault, _, err = ctx.GetClient().CliAPI.GetVaultByNameAPIKey(apiCtx, name).Execute()
	} else {
		ctx.DebugLog("Making API request to get vault by ID: %s", id)
		vault, _, err = ctx.GetClient().CliAPI.GetVaultByAPIKey(apiCtx, id).Execute()
	}

	if err != nil {
		ctx.DebugLog("API request failed: %v", err)
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	ctx.DebugLog("API request successful, vault retrieved")

	if outputFile != "" {
		handleFileOutput(vault, outputFile, followUpCommand, ctx.DebugLog)
	} else {
		ctx.DebugLog("Outputting vault value to stdout")
		fmt.Println(vault.Value)
	}
	ctx.DebugLog("Get command completed successfully")
}

func handleFileOutput(vault *openapi.Vault, outputFile, followUpCommand string, debugLog func(string, ...any)) {
	debugLog("Writing vault value to file: %s", outputFile)

	// Check if file exists and get its modification time
	var fileModTime time.Time
	var fileExists bool
	if fileInfo, err := os.Stat(outputFile); err == nil {
		fileModTime = fileInfo.ModTime()
		fileExists = true
		debugLog("File exists, last modified: %s", fileModTime.Format(time.RFC3339))
	} else {
		debugLog("File does not exist")
	}

	// Parse vault's updated timestamp
	var vaultUpdatedAt time.Time
	if vault.UpdatedAt != nil {
		vaultUpdatedAt = *vault.UpdatedAt
	} else {
		vaultUpdatedAt = time.Now() // Fallback to current time if no timestamp
	}
	debugLog("Vault last updated: %s", vaultUpdatedAt.Format(time.RFC3339))

	// Determine if vault has been updated
	vaultHasUpdates := !fileExists || vaultUpdatedAt.After(fileModTime)
	debugLog("Vault has updates: %v", vaultHasUpdates)

	err := os.WriteFile(outputFile, []byte(vault.Value), 0600)
	if err != nil {
		debugLog("File write failed: %v", err)
		fmt.Fprintf(os.Stderr, "Error writing to file: %v\n", err)
		os.Exit(1)
	}
	debugLog("File write successful")

	if vaultHasUpdates {
		fmt.Printf("Vault value written to %s (vault was updated)\n", outputFile)

		// Execute follow-up command if specified
		if followUpCommand != "" {
			executeFollowUpCommand(followUpCommand, debugLog)
		}
	} else {
		fmt.Printf("Vault value written to %s (no updates since last fetch)\n", outputFile)
	}
}

func executeFollowUpCommand(followUpCommand string, debugLog func(string, ...any)) {
	debugLog("Executing follow-up command: %s", followUpCommand)
	fmt.Printf("Executing follow-up command: %s\n", followUpCommand)

	cmdParts := strings.Fields(followUpCommand)
	if len(cmdParts) > 0 {
		// #nosec G204 - This is intentional: user-provided command execution is the expected behavior
		cmd := exec.Command(cmdParts[0], cmdParts[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			debugLog("Follow-up command failed: %v", err)
			fmt.Fprintf(os.Stderr, "Warning: Follow-up command failed: %v\n", err)
		} else {
			debugLog("Follow-up command completed successfully")
		}
	}
}
