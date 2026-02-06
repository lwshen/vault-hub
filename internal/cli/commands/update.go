package commands

import (
	"context"
	"fmt"
	"os"

	openapi "github.com/lwshen/vault-hub-go-client"
	"github.com/spf13/cobra"
)

// NewUpdateCommand creates the update command
func NewUpdateCommand(ctx *CommandContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update --name/--id <vault-name-or-id> [flags]",
		Short: "Update a vault by name or unique ID",
		Long: `Update an existing vault's properties. You can update any combination of:
- Value (the encrypted secret)
- Description
- Category
- Favourite flag

You must specify the vault by either its name or unique ID.
Only the fields you provide will be updated; others remain unchanged.

Examples:
  vault-hub update --name my-api-keys --value "new-secret-value"
  vault-hub update --id abc123 --description "Updated description"
  vault-hub update --name my-vault --category "production" --favourite
  vault-hub update --name my-vault --value "new-value" --description "New desc"`,
		Run: func(cmd *cobra.Command, args []string) {
			runUpdateCommand(cmd, args, ctx)
		},
	}

	cmd.Flags().StringP("name", "n", "", "Vault name")
	cmd.Flags().StringP("id", "i", "", "Vault Unique ID")
	cmd.Flags().String("value", "", "New vault value")
	cmd.Flags().String("description", "", "New description")
	cmd.Flags().String("category", "", "New category")
	cmd.Flags().Bool("favourite", false, "Set as favourite")
	cmd.Flags().Bool("unfavourite", false, "Unset as favourite")

	return cmd
}

// updateCommandParams holds the parsed command parameters
type updateCommandParams struct {
	name        string
	id          string
	value       *string
	description *string
	category    *string
	favourite   *bool
}

// runUpdateCommand executes the vault update operation
func runUpdateCommand(cmd *cobra.Command, _ []string, ctx *CommandContext) {
	ctx.DebugLog("Executing update command")

	// Parse command flags
	params := parseUpdateCommandFlags(cmd, ctx)
	ctx.DebugLog("Parameters - name: '%s', id: '%s'", params.name, params.id)

	// Validate required parameters
	if err := validateUpdateParams(params); err != nil {
		ctx.DebugLog("Validation failed: %v", err)
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Update vault via API
	vault, err := updateVault(params, ctx)
	if err != nil {
		ctx.DebugLog("API request failed: %v", err)
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	ctx.DebugLog("API request successful, vault updated")

	// Display success message
	fmt.Printf("Vault '%s' updated successfully\n", vault.Name)
	if params.value != nil {
		fmt.Println("✓ Value updated")
	}
	if params.description != nil {
		fmt.Println("✓ Description updated")
	}
	if params.category != nil {
		fmt.Println("✓ Category updated")
	}
	if params.favourite != nil {
		if *params.favourite {
			fmt.Println("✓ Marked as favourite")
		} else {
			fmt.Println("✓ Unmarked as favourite")
		}
	}
}

// parseUpdateCommandFlags extracts and returns command flags
func parseUpdateCommandFlags(cmd *cobra.Command, ctx *CommandContext) updateCommandParams {
	params := updateCommandParams{
		name: ctx.MustGetStringFlag(cmd, "name"),
		id:   ctx.MustGetStringFlag(cmd, "id"),
	}

	// Parse optional update fields
	if cmd.Flags().Changed("value") {
		value := ctx.MustGetStringFlag(cmd, "value")
		params.value = &value
	}
	if cmd.Flags().Changed("description") {
		description := ctx.MustGetStringFlag(cmd, "description")
		params.description = &description
	}
	if cmd.Flags().Changed("category") {
		category := ctx.MustGetStringFlag(cmd, "category")
		params.category = &category
	}

	// Handle favourite/unfavourite flags
	favourite, _ := cmd.Flags().GetBool("favourite")
	unfavourite, _ := cmd.Flags().GetBool("unfavourite")
	if favourite && unfavourite {
		fmt.Fprintf(os.Stderr, "Error: cannot use both --favourite and --unfavourite\n")
		os.Exit(1)
	}
	if favourite {
		fav := true
		params.favourite = &fav
	} else if unfavourite {
		fav := false
		params.favourite = &fav
	}

	return params
}

// validateUpdateParams ensures required parameters are provided
func validateUpdateParams(params updateCommandParams) error {
	if params.name == "" && params.id == "" {
		return fmt.Errorf("either --name or --id must be provided")
	}
	if params.name != "" && params.id != "" {
		return fmt.Errorf("cannot specify both --name and --id")
	}
	if params.value == nil && params.description == nil && params.category == nil && params.favourite == nil {
		return fmt.Errorf("at least one field must be updated (--value, --description, --category, --favourite, or --unfavourite)")
	}
	return nil
}

// updateVault sends the update request to the API
func updateVault(params updateCommandParams, ctx *CommandContext) (*openapi.Vault, error) {
	apiCtx := context.Background()

	// Build update request - only include fields that should be updated
	updateReq := openapi.UpdateVaultRequest{
		Value:       params.value,
		Description: params.description,
		Category:    params.category,
		Favourite:   params.favourite,
	}

	var vault *openapi.Vault
	var err error

	if params.name != "" {
		ctx.DebugLog("Making API request to update vault by name: %s", params.name)
		vault, _, err = ctx.GetClient().CliAPI.UpdateVaultByNameAPIKey(apiCtx, params.name).
			UpdateVaultRequest(updateReq).Execute()
	} else {
		ctx.DebugLog("Making API request to update vault by ID: %s", params.id)
		vault, _, err = ctx.GetClient().CliAPI.UpdateVaultByAPIKey(apiCtx, params.id).
			UpdateVaultRequest(updateReq).Execute()
	}

	return vault, err
}
