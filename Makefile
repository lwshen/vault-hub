.PHONY: build-web build-server build clean

# Build the web application
build-web:
	cd apps/web && pnpm install && pnpm build

# Copy web assets to internal package and build server
build-server: build-web
	cp -r apps/web/dist internal/web/
	go build -o vault-hub ./apps/server

# Build everything (web + server with embedded assets)
build: build-server

# Clean build artifacts
clean:
	rm -f vault-hub vault-hub-embedded
	rm -rf internal/web/dist
	rm -rf apps/web/dist

# Development server (without embedded assets)
dev-server:
	go run ./apps/server

help:
	@echo "Available targets:"
	@echo "  build-web    - Build the web application"
	@echo "  build-server - Build server with embedded web assets"
	@echo "  build        - Build everything (web + server)"
	@echo "  clean        - Clean build artifacts"
	@echo "  dev-server   - Run development server"
	@echo "  help         - Show this help message"