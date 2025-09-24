#!/bin/bash

set -e

echo "Building VaultHub with embedded web app..."

# Step 1: Build the web app
echo "Step 1: Building web app..."
cd apps/web
pnpm install
pnpm build
cd ../..

# Step 2: Build the Go binary with embedded web assets
echo "Step 2: Building Go binary with embedded web assets..."
cat > main.go << 'EOF'
package main

import (
	"embed"
	"log"
	"log/slog"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/lwshen/vault-hub/internal/config"
	"github.com/lwshen/vault-hub/internal/version"
	"github.com/lwshen/vault-hub/model"
	"github.com/lwshen/vault-hub/route"
	slogfiber "github.com/samber/slog-fiber"
)

//go:embed all:apps/web/dist
var webDistFS embed.FS

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	logger.Info("Starting VaultHub Server", "version", version.Version, "commit", version.Commit)

	err := model.Open(logger)
	if err != nil {
		logger.Error("Failed to open database", "error", err)
		os.Exit(1)
	}

	app := fiber.New()

	app.Use(slogfiber.New(logger))

	route.SetupRoutes(app, webDistFS)

	log.Fatal(app.Listen(":" + config.AppPort))
}
EOF

go build -o vault-hub-embedded .

# Clean up temporary main.go
rm main.go

echo "Build complete! Binary: vault-hub-embedded"
echo "Binary size: $(du -h vault-hub-embedded | cut -f1)"