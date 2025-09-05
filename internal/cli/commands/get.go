package commands

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"

	openapi "github.com/lwshen/vault-hub-go-client"
	"github.com/spf13/cobra"
)

// NewGetCommand creates the get command
func NewGetCommand(ctx *CommandContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get --name/--id <vault-name-or-id> --output <file> --exec <command>",
		Short: "Get a specific vault by name or unique ID",
		Long: `Get detailed information about a specific vault including its encrypted value.
You can specify the vault by either its name or unique ID.

Examples:
  vault-hub get --name my-api-keys
  vault-hub get --id abc123-def456-ghi789
  vault-hub get --name my-api-keys --output ./secrets.txt
  vault-hub get --name my-api-keys --output .env --exec "source .env && echo 'Environment loaded'"`,
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

type UpdateResult struct {
	HasUpdates       bool
	TimestampChanged bool
	ContentChanged   bool
	FileExists       bool
	Reason           string
}

func isVaultUpdated(vault *openapi.Vault, filePath string, debugLog func(string, ...any)) (UpdateResult, error) {
	var result UpdateResult
	var fileModTime time.Time

	// Check if file exists and get its modification time
	if fileInfo, err := os.Stat(filePath); err == nil {
		fileModTime = fileInfo.ModTime()
		result.FileExists = true
		debugLog("File exists, last modified: %s", fileModTime.Format(time.RFC3339))

		// Compare file content with vault value
		if existingContent, readErr := os.ReadFile(filePath); readErr == nil {
			result.ContentChanged = string(existingContent) != vault.Value
			debugLog("Content comparison - file differs from vault: %v", result.ContentChanged)
		} else {
			debugLog("Could not read existing file for content comparison: %v", readErr)
			result.ContentChanged = true // Assume content changed if we can't read the file
		}
	} else {
		debugLog("File does not exist")
		result.ContentChanged = true // New file, so content is "changed"
	}

	// Parse vault's updated timestamp
	var vaultUpdatedAt time.Time
	var hasValidTimestamp bool
	if vault.UpdatedAt != nil {
		vaultUpdatedAt = *vault.UpdatedAt
		hasValidTimestamp = true
		debugLog("Vault last updated: %s", vaultUpdatedAt.Format(time.RFC3339))
	} else {
		debugLog("Vault has no timestamp - treating as new vault")
	}

	// Determine if vault has been updated
	// Only proceed with update check if we have valid timestamps
	result.TimestampChanged = hasValidTimestamp && result.FileExists && vaultUpdatedAt.After(fileModTime)
	result.HasUpdates = !result.FileExists || result.TimestampChanged || result.ContentChanged

	// Set update reason
	if !result.FileExists {
		result.Reason = "new file"
	} else if result.TimestampChanged && result.ContentChanged {
		result.Reason = "vault updated and content differs"
	} else if result.TimestampChanged {
		result.Reason = "vault was updated"
	} else if result.ContentChanged {
		result.Reason = "content differs from vault"
	} else {
		result.Reason = "no updates - content matches vault"
	}

	debugLog("Vault has updates: %v (timestamp: %v, content: %v)", result.HasUpdates, result.TimestampChanged, result.ContentChanged)
	return result, nil
}

func handleFileOutput(vault *openapi.Vault, outputFile, followUpCommand string, debugLog func(string, ...any)) {
	debugLog("Writing vault value to file: %s", outputFile)

	updateResult, err := isVaultUpdated(vault, outputFile, debugLog)
	if err != nil {
		debugLog("Error checking vault updates: %v", err)
		fmt.Fprintf(os.Stderr, "Error checking vault updates: %v\n", err)
		os.Exit(1)
	}

	if updateResult.HasUpdates {
		err = os.WriteFile(outputFile, []byte(vault.Value), 0600)
		if err != nil {
			debugLog("File write failed: %v", err)
			fmt.Fprintf(os.Stderr, "Error writing to file: %v\n", err)
			os.Exit(1)
		}
		debugLog("File write successful")
		fmt.Printf("Vault value written to %s (%s)\n", outputFile, updateResult.Reason)

		// Execute follow-up command if specified
		if followUpCommand != "" {
			executeFollowUpCommand(followUpCommand, debugLog)
		}
	} else {
		debugLog("No updates detected, file not modified")
	}
}

func executeFollowUpCommand(followUpCommand string, debugLog func(string, ...any)) {
	debugLog("Executing follow-up command: %s", followUpCommand)
	fmt.Printf("Executing follow-up command: %s\n", followUpCommand)

	// Use shell to handle complex commands properly
	cmd := exec.Command("sh", "-c", followUpCommand)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		debugLog("Follow-up command failed: %v", err)
		fmt.Fprintf(os.Stderr, "Warning: Follow-up command failed: %v\n", err)
	} else {
		debugLog("Follow-up command completed successfully")
	}
}
