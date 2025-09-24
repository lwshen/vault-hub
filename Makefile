.PHONY: build clean build-web build-server run test

# Variables
WEB_DIR = apps/web
EMBED_DIR = internal/embed
DIST_DIR = $(EMBED_DIR)/dist
SERVER_BINARY = vault-hub-server
CLI_BINARY = vault-hub-cli
CRON_BINARY = vault-hub-cron

# Build everything
build: clean build-web build-server build-cli build-cron
	@echo "Build complete!"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(WEB_DIR)/dist
	@rm -rf $(DIST_DIR)
	@rm -f $(SERVER_BINARY) $(CLI_BINARY) $(CRON_BINARY)

# Build the web application
build-web:
	@echo "Building web application..."
	@cd $(WEB_DIR) && pnpm install --silent && pnpm build
	@echo "Copying dist to embed directory..."
	@rm -rf $(DIST_DIR)
	@cp -r $(WEB_DIR)/dist $(EMBED_DIR)/

# Build the server binary with embedded web assets
build-server: build-web
	@echo "Building server binary with embedded web assets..."
	@go build -o $(SERVER_BINARY) ./apps/server

# Build the CLI binary
build-cli:
	@echo "Building CLI binary..."
	@go build -o $(CLI_BINARY) ./apps/cli

# Build the cron binary
build-cron:
	@echo "Building cron binary..."
	@go build -o $(CRON_BINARY) ./apps/cron

# Run the server (for development)
run: build-web
	@echo "Running server..."
	@go run ./apps/server

# Test the build
test-build: build
	@echo "Testing build..."
	@./$(SERVER_BINARY) --version || echo "Server binary test"
	@./$(CLI_BINARY) version || echo "CLI binary test"
	@echo "Build test complete!"

# Install dependencies
deps:
	@echo "Installing Go dependencies..."
	@go mod download
	@echo "Installing web dependencies..."
	@cd $(WEB_DIR) && pnpm install

# Development mode - watch for changes
dev:
	@echo "Starting development mode..."
	@cd $(WEB_DIR) && pnpm dev &
	@go run ./apps/server