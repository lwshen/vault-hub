package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type rewrite struct {
	old string
	new string
}

const (
	fiberV2Import    = "github.com/gofiber/fiber/" + "v2"
	oldFiberCtx      = "*fiber" + ".Ctx"
	oldUserContext   = ".User" + "Context()"
	oldBodyParser    = ".Body" + "Parser("
	oldSetUserScopes = "c.Context().SetUserValue(ApiKeyAuthScopes, []string{})"
	generatedFile    = "generated.go"
)

var rewrites = []rewrite{
	{old: fiberV2Import, new: "github.com/gofiber/fiber/v3"},
	{old: oldFiberCtx, new: "fiber.Ctx"},
	{old: oldUserContext, new: ".Context()"},
	{old: oldBodyParser, new: ".Bind().Body("},
	{
		old: oldSetUserScopes,
		new: "fiber.StoreInContext(c, ApiKeyAuthScopes, []string{})",
	},
}

var unsupported = []string{
	fiberV2Import,
	oldFiberCtx,
	oldUserContext,
	oldBodyParser,
	oldSetUserScopes,
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: fiber-v3-fix <generated-file>")
		os.Exit(2)
	}

	if filepath.Clean(os.Args[1]) != generatedFile {
		fmt.Fprintf(os.Stderr, "unexpected generated file %q\n", os.Args[1])
		os.Exit(2)
	}

	content, err := os.ReadFile(generatedFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read %s: %v\n", generatedFile, err)
		os.Exit(1)
	}

	updated := string(content)
	for _, rewrite := range rewrites {
		updated = strings.ReplaceAll(updated, rewrite.old, rewrite.new)
	}

	for _, marker := range unsupported {
		if strings.Contains(updated, marker) {
			fmt.Fprintf(os.Stderr, "unsupported Fiber v2 pattern remains in %s: %s\n", generatedFile, marker)
			os.Exit(1)
		}
	}

	//nolint:gosec // go:generate only rewrites generated.go in this directory after validation above.
	if err := os.WriteFile(generatedFile, []byte(updated), 0o600); err != nil {
		fmt.Fprintf(os.Stderr, "write %s: %v\n", generatedFile, err)
		os.Exit(1)
	}
}
